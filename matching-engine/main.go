package main

import (
	"container/heap"
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/mann-som/zerodha/internal/models"
	"github.com/mann-som/zerodha/internal/repositories"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type OrderHeap []models.Order

func (h *OrderHeap) Len() int           { return len(*h) }
func (h *OrderHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }
func (h *OrderHeap) Push(x interface{}) { *h = append(*h, x.(models.Order)) }
func (h *OrderHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (h OrderHeap) Less(i, j int) bool { return false }

type BuyHeap struct{ OrderHeap }

func (h BuyHeap) Less(i, j int) bool {
	return h.OrderHeap[i].Price > h.OrderHeap[j].Price || (h.OrderHeap[i].Price == h.OrderHeap[j].Price && h.OrderHeap[i].CreatedAt.Before(h.OrderHeap[j].CreatedAt))
} // Max price, earliest time

type SellHeap struct{ OrderHeap }

func (h SellHeap) Less(i, j int) bool {
	return h.OrderHeap[i].Price < h.OrderHeap[j].Price || (h.OrderHeap[i].Price == h.OrderHeap[j].Price && h.OrderHeap[i].CreatedAt.Before(h.OrderHeap[j].CreatedAt))
} // Min price, earliest time

type MatchingEngine struct {
	db          *gorm.DB
	redisClient *redis.Client
	userRepo    *repositories.UserRepository
	stockRepo   *repositories.StockRepository
	orderBooks  map[string]*OrderBook
	mu          sync.Mutex
}

type OrderBook struct {
	buyHeap  *BuyHeap
	sellHeap *SellHeap
	mu       sync.Mutex
}

func NewMatchingEngine(db *gorm.DB, redisClient *redis.Client, userRepo *repositories.UserRepository, stockRepo *repositories.StockRepository) *MatchingEngine {
	return &MatchingEngine{
		db:          db,
		redisClient: redisClient,
		userRepo:    userRepo,
		stockRepo:   stockRepo,
		orderBooks:  make(map[string]*OrderBook),
	}
}

func (me *MatchingEngine) Run() {
	for {
		msg, err := me.redisClient.BLPop(context.Background(), 1*time.Second, "order_queue").Result()
		if err != nil {
			if err == redis.Nil {
				continue
			}
			log.Printf("Error popping from queue: %v", err)
			continue
		}
		var order models.Order
		if err := json.Unmarshal([]byte(msg[1]), &order); err != nil {
			log.Printf("Error unmarshaling order: %v", err)
			continue
		}

		me.mu.Lock()
		if _, ok := me.orderBooks[order.Symbol]; !ok {
			me.orderBooks[order.Symbol] = &OrderBook{
				buyHeap:  &BuyHeap{},
				sellHeap: &SellHeap{},
			}
		}
		me.mu.Unlock()

		me.orderBooks[order.Symbol].mu.Lock()
		me.matchOrder(order)
		me.orderBooks[order.Symbol].mu.Unlock()
	}
}

func (me *MatchingEngine) matchOrder(order models.Order) {
	orderBook := me.orderBooks[order.Symbol]
	if order.Side == "buy" {
		heap.Push(orderBook.buyHeap, order)
		for orderBook.sellHeap.Len() > 0 && orderBook.sellHeap.OrderHeap[0].Price <= order.Price {
			match := heap.Pop(orderBook.sellHeap).(models.Order)
			me.executeTrade(order, match)
		}
	} else {
		heap.Push(orderBook.sellHeap, order)
		for orderBook.buyHeap.Len() > 0 && orderBook.buyHeap.OrderHeap[0].Price >= order.Price {
			match := heap.Pop(orderBook.buyHeap).(models.Order)
			me.executeTrade(match, order)
		}
	}
}

func (me *MatchingEngine) executeTrade(buy models.Order, sell models.Order) {
	tradeQuantity := min(buy.Quantity, sell.Quantity)
	tradePrice := (buy.Price + sell.Price) / 2

	err := me.db.Transaction(func(tx *gorm.DB) error {
		buyer, err := me.userRepo.GetWithTx(tx, buy.UserID)
		if err != nil {
			return err
		}
		seller, err := me.userRepo.GetWithTx(tx, sell.UserID)
		if err != nil {
			return err
		}

		totalCost := float64(tradeQuantity) * tradePrice
		buyer.Balance -= totalCost
		seller.Balance += totalCost

		if err := tx.Save(&buyer).Error; err != nil {
			return err
		}
		if err := tx.Save(&seller).Error; err != nil {
			return err
		}

		buy.Quantity -= tradeQuantity
		sell.Quantity -= tradeQuantity
		if buy.Quantity == 0 {
			buy.Status = "filled"
		} else {
			buy.Status = "partially_filled"
		}
		if sell.Quantity == 0 {
			sell.Status = "filled"
		} else {
			sell.Status = "partially_filled"
		}

		if err := tx.Save(&buy).Error; err != nil {
			return err
		}
		if err := tx.Save(&sell).Error; err != nil {
			return err
		}

		// Update stock price
		stock, err := me.stockRepo.GetBySymbolWithTx(tx, buy.Symbol)
		if err != nil {
			return err
		}
		stock.CurrentPrice = tradePrice
		if err := tx.Save(&stock).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Printf("Error executing trade: %v", err)
	}

	// Push update via WebSocket (implement in main backend)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file, using defaults")
	}

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN not set")
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	redisClient := redis.NewClient(&redis.Options{Addr: redisAddr})

	engine := NewMatchingEngine(db, redisClient, repositories.NewUserRepository(db), repositories.NewStockRepository(db))

	log.Println("Matching engine starting...")
	engine.Run()
}
