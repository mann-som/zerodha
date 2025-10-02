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
	*h = old[:n-1]
	return x
}
func (h *OrderHeap) Peek() *models.Order {
	if h.Len() == 0 {
		return nil
	}
	return &(*h)[0]
}

type BuyHeap struct{ OrderHeap }

func (h *BuyHeap) Less(i, j int) bool {
	a := h.OrderHeap[i]
	b := h.OrderHeap[j]
	if a.Price != b.Price {
		return a.Price > b.Price
	}
	return a.CreatedAt.Before(b.CreatedAt)
}

type SellHeap struct{ OrderHeap }

func (h *SellHeap) Less(i, j int) bool {
	a := h.OrderHeap[i]
	b := h.OrderHeap[j]
	if a.Price != b.Price {
		return a.Price < b.Price
	}
	return a.CreatedAt.Before(b.CreatedAt)
}

func (h *BuyHeap) Len() int           { return h.OrderHeap.Len() }
func (h *BuyHeap) Swap(i, j int)      { h.OrderHeap.Swap(i, j) }
func (h *BuyHeap) Push(x interface{}) { h.OrderHeap.Push(x) }
func (h *BuyHeap) Pop() interface{}   { return h.OrderHeap.Pop() }
func (h *BuyHeap) Peek() *models.Order {
	return h.OrderHeap.Peek()
}

func (h *SellHeap) Len() int           { return h.OrderHeap.Len() }
func (h *SellHeap) Swap(i, j int)      { h.OrderHeap.Swap(i, j) }
func (h *SellHeap) Push(x interface{}) { h.OrderHeap.Push(x) }
func (h *SellHeap) Pop() interface{}   { return h.OrderHeap.Pop() }
func (h *SellHeap) Peek() *models.Order {
	return h.OrderHeap.Peek()
}

type MatchingEngine struct {
	db          *gorm.DB
	redisClient *redis.Client
	userRepo    *repositories.UserRepository
	stockRepo   *repositories.StockRepository
	orderBooks  map[string]*OrderBook // per symbol
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
		log.Printf("Popped message from queue: %v", msg[1])

		var order models.Order
		if err := json.Unmarshal([]byte(msg[1]), &order); err != nil {
			log.Printf("Error unmarshaling order: %v", err)
			continue
		}
		log.Printf("Unmarshaled order: ID %v, Symbol %v, Side %v, Quantity %v, Price %v", order.ID, order.Symbol, order.Side, order.Quantity, order.Price)

		me.mu.Lock()
		if _, ok := me.orderBooks[order.Symbol]; !ok {
			me.orderBooks[order.Symbol] = &OrderBook{
				buyHeap:  &BuyHeap{},
				sellHeap: &SellHeap{},
			}
			heap.Init(me.orderBooks[order.Symbol].buyHeap)
			heap.Init(me.orderBooks[order.Symbol].sellHeap)
			log.Printf("Created new order book for symbol: %v", order.Symbol)
		}
		me.mu.Unlock()

		ob := me.orderBooks[order.Symbol]
		ob.mu.Lock()
		me.matchOrder(order, ob)
		ob.mu.Unlock()
	}
}

// matchOrder: new incoming order is matched against opposite side first.
// Any leftover (unfilled quantity) is inserted into its side heap.
func (me *MatchingEngine) matchOrder(incoming models.Order, ob *OrderBook) {
	log.Printf("Matching order: ID %v, Side %v, Quantity %v, Price %v for symbol %v", incoming.ID, incoming.Side, incoming.Quantity, incoming.Price, incoming.Symbol)
	remaining := incoming.Quantity

	if incoming.Side == "buy" {
		// match against sell heap (best sells: lowest price)
		for remaining > 0 && ob.sellHeap.Len() > 0 {
			bestSell := ob.sellHeap.Peek()
			if bestSell == nil {
				break
			}
			// only match if best sell price <= buy price
			if bestSell.Price > incoming.Price {
				break
			}
			// pop the resting sell
			resting := heap.Pop(ob.sellHeap).(models.Order)

			tradeQty := min(remaining, resting.Quantity)
			tradePrice := resting.Price // resting order price

			// Execute trade for tradeQty at tradePrice
			if err := me.executeTradeDB(incoming, resting, tradeQty, tradePrice); err != nil {
				// If transaction fails, we should push resting back and abort matching to avoid inconsistencies
				log.Printf("Trade DB transaction failed: %v. Pushing resting back and aborting matching.", err)
				heap.Push(ob.sellHeap, resting) // push back whole resting
				return
			}

			remaining -= tradeQty
			resting.Quantity -= tradeQty

			// if resting still has qty, push back into sell heap
			if resting.Quantity > 0 {
				heap.Push(ob.sellHeap, resting)
			}
			// continue the loop to try match more
		}

		// if any remaining, insert the leftover buy into buy heap
		if remaining > 0 {
			incoming.Quantity = remaining
			incoming.Status = "open"
			heap.Push(ob.buyHeap, incoming)
			log.Printf("Inserted leftover buy order ID %v qty %v into buy heap", incoming.ID, remaining)
		} else {
			log.Printf("Incoming buy order ID %v fully filled", incoming.ID)
		}
	} else { // incoming.Side == "sell"
		// match against buy heap (best buys: highest price)
		for remaining > 0 && ob.buyHeap.Len() > 0 {
			bestBuy := ob.buyHeap.Peek()
			if bestBuy == nil {
				break
			}
			// only match if best buy price >= sell price
			if bestBuy.Price < incoming.Price {
				break
			}
			// pop the resting buy
			resting := heap.Pop(ob.buyHeap).(models.Order)

			tradeQty := min(remaining, resting.Quantity)
			tradePrice := resting.Price // resting order price (resting buy price)

			// Execute trade
			if err := me.executeTradeDB(resting, incoming, tradeQty, tradePrice); err != nil {
				log.Printf("Trade DB transaction failed: %v. Pushing resting back and aborting matching.", err)
				heap.Push(ob.buyHeap, resting)
				return
			}

			remaining -= tradeQty
			resting.Quantity -= tradeQty

			// if resting still has qty, push back into buy heap
			if resting.Quantity > 0 {
				heap.Push(ob.buyHeap, resting)
			}
		}

		// if any remaining, insert leftover sell into sell heap
		if remaining > 0 {
			incoming.Quantity = remaining
			incoming.Status = "open"
			heap.Push(ob.sellHeap, incoming)
			log.Printf("Inserted leftover sell order ID %v qty %v into sell heap", incoming.ID, remaining)
		} else {
			log.Printf("Incoming sell order ID %v fully filled", incoming.ID)
		}
	}
}

// executeTradeDB: does DB transaction for the given tradeQuantity at tradePrice.
// buy and sell passed are the *original* buy and sell order structs (incoming or resting).
// This function updates buyer/seller balances, reduces order quantities appropriately, sets statuses,
// and updates stock current price. It assumes caller will reinsert any leftovers into heaps.
func (me *MatchingEngine) executeTradeDB(buy models.Order, sell models.Order, tradeQty int, tradePrice float64) error {
	totalCost := float64(tradeQty) * tradePrice

	return me.db.Transaction(func(tx *gorm.DB) error {
		// fetch buyer and seller with lock (GetWithTx should use SELECT ... FOR UPDATE if supported)
		buyer, err := me.userRepo.GetWithTx(tx, buy.UserID)
		if err != nil {
			log.Printf("Error fetching buyer %v: %v", buy.UserID, err)
			return err
		}
		seller, err := me.userRepo.GetWithTx(tx, sell.UserID)
		if err != nil {
			log.Printf("Error fetching seller %v: %v", sell.UserID, err)
			return err
		}

		// Basic balance check (optional - for safety)
		// You might want to ensure buyer has sufficient funds; if not, abort transaction.
		if buyer.Balance < totalCost {
			log.Printf("Buyer %v has insufficient balance: %v required %v", buyer.ID, buyer.Balance, totalCost)
			return gorm.ErrInvalidTransaction
		}

		// Update balances
		buyer.Balance -= totalCost
		seller.Balance += totalCost

		if err := tx.Save(&buyer).Error; err != nil {
			return err
		}
		if err := tx.Save(&seller).Error; err != nil {
			return err
		}

		// Update buy order record in DB: reduce quantity and update status
		var buyOrder models.Order
		if err := tx.Where("id = ?", buy.ID).First(&buyOrder).Error; err != nil {
			return err
		}
		buyOrder.Quantity -= tradeQty
		if buyOrder.Quantity <= 0 {
			buyOrder.Status = "filled"
			buyOrder.Quantity = 0
		} else {
			buyOrder.Status = "partially_filled"
		}
		if err := tx.Save(&buyOrder).Error; err != nil {
			return err
		}

		// Update sell order record in DB: reduce quantity and update status
		var sellOrder models.Order
		if err := tx.Where("id = ?", sell.ID).First(&sellOrder).Error; err != nil {
			return err
		}
		sellOrder.Quantity -= tradeQty
		if sellOrder.Quantity <= 0 {
			sellOrder.Status = "filled"
			sellOrder.Quantity = 0
		} else {
			sellOrder.Status = "partially_filled"
		}
		if err := tx.Save(&sellOrder).Error; err != nil {
			return err
		}

		// Update stock price/current price
		stock, err := me.stockRepo.GetBySymbolWithTx(tx, buy.Symbol)
		if err != nil {
			return err
		}
		stock.CurrentPrice = tradePrice
		if err := tx.Save(&stock).Error; err != nil {
			return err
		}

		// You can also persist a trade record here (not currently implemented)
		// e.g., tx.Create(&models.Trade{BuyOrderID: buy.ID, SellOrderID: sell.ID, Quantity: tradeQty, Price: tradePrice, ...})

		return nil
	})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
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
	log.Println("Connected to Redis")

	engine := NewMatchingEngine(db, redisClient, repositories.NewUserRepository(db), repositories.NewStockRepository(db))

	log.Println("Matching engine starting...")
	engine.Run()
}
