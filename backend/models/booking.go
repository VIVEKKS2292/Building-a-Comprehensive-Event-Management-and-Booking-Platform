package models

import (
	"time"
)

type Booking struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	EventID   uint      `json:"event_id"`
	UserID    uint      `json:"user_id"`
	Tickets   int       `json:"tickets"`
	Status    string    `json:"status"` // "booked", "canceled", "pending"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TicketAvailability struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	EventID   uint `json:"event_id"`
	Available int  `json:"available"`
}
