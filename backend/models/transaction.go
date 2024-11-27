package models

import (
	"time"
)

type Transaction struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	EventID   uint      `json:"event_id"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"` // "success", "failure", "pending"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
