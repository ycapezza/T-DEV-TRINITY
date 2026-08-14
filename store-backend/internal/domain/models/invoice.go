package models

import "time"

type Invoice struct {
    ID        uint        `json:"id" gorm:"primaryKey"`
    UserID    uint        `json:"user_id"`
    User      User        `json:"user" gorm:"foreignKey:UserID"`
    Total     float64     `json:"total"`
    Status    string      `json:"status"`
    PaypalID  string      `json:"paypal_id"`
    Items     []OrderItem `json:"items" gorm:"foreignKey:InvoiceID"`
    CreatedAt time.Time   `json:"created_at"`
    UpdatedAt time.Time   `json:"updated_at"`
}