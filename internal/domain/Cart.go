package domain

import "time"

type Cart struct {
	ID        uint      `json:"id" gorm:"PrimaryKey"`
	UserID    uint      `json:"user_id"`
	ProductID uint      `json:"product_id"`
	Name      string    `json:"name" gorm:"index;"`
	ImageURL  string    `json:"image_url"`
	SellerId  uint      `json:"seller_id"`
	Price     float64   `json:"price"`
	Qty       uint      `json:"qty"`
	CreatedAt time.Time `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:current_timestamp"`
}
