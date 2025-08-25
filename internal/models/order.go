package models

import "time"

type Order struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserID    string    `json:"user_id" gorm:"type:uuid;not null"`
	Symbol    string    `json:"symbol" gorm:"not null"`
	Side      string    `json:"side" gorm:"not null"`
	Quantity  int       `json:"quantity" gorm:"not null"`
	Price     float64   `json:"price" gorm:"not null"`
	Status    string    `json:"status" gorm:"not null;default:pending"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
