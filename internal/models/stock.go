package models

type Stock struct {
	ID           string  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Symbol       string  `json:"symbol" gorm:"unique;not null"`
	Description  string  `json:"description" gorm:"not null"`
	InitialPrice float64 `json:"initial_price" gorm:"not null"`
	CurrentPrice float64 `json:"current_price" gorm:"not null"`
}
