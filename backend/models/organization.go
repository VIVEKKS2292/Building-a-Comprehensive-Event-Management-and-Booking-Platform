package models

import (
	"gorm.io/gorm"
)

// Organization - Organization model
type Organization struct {
	gorm.Model
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Website     string `json:"website"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	Contact     string `json:"contact"`
	Industry    string `json:"industry"`
	UserID      uint   `json:"user_id" gorm:"unique"` // Enforce unique constraint
}
