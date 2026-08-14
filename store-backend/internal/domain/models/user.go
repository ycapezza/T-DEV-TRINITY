package models

import "time"

type User struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    FirstName   string    `json:"first_name"`
    LastName    string    `json:"last_name"`
    Email       string    `json:"email" gorm:"unique"`
    Password    string    `json:"-"`
    PhoneNumber string    `json:"phone_number"`
    Address     string    `json:"address"`
    ZipCode     string    `json:"zip_code"`
    City        string    `json:"city"`
    Country     string    `json:"country"`
    IsAdmin     bool      `json:"is_admin" gorm:"default:false"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}