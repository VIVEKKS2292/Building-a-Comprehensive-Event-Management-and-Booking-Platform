package models

type Event struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Location    string  `json:"location"`
	Date        string  `json:"date"`
	Time        string  `json:"time"`
	Category    string  `json:"category"`
	TicketPrice float64 `json:"ticket_price"`
	Organizer   string  `json:"organizer"`
	UserID      uint    `json:"user_id"` // Foreign key reference to the User model
}
