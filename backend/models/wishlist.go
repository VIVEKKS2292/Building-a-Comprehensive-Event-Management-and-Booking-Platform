package models

type Wishlist struct {
	ID      uint `gorm:"primaryKey"`
	UserID  uint `gorm:"foreignKey:UserID"`
	EventID uint `gorm:"foreignKey:EventID"`
}
