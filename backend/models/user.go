package models

type User struct {
    ID       uint   `gorm:"primaryKey" json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email" gorm:"unique"`
    Password string `json:"password"` // Store hashed passwords
    Role     string `json:"role"`
}
