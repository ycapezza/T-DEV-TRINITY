package models

import (
	"time"

	"github.com/lib/pq"
)

type Product struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	Name            string         `json:"name"`
	Price           float64        `json:"price"`
	Brand           string         `json:"brand"`
	Picture         string         `json:"picture"`
	Categories      pq.StringArray `json:"categories" gorm:"type:text[]"`
	NutritionalInfo string         `json:"nutritional_info"`
	StockQuantity   int            `json:"stock_quantity"`
	OpenFoodFactsID string         `json:"open_food_facts_id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}
