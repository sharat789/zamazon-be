package domain

import "time"

type Product struct {
	ID          uint    `json:"id" gorm:"PrimaryKey"`
	Name        string  `json:"name" gorm:"index;"`
	Description string  `json:"description"`
	CategoryID  uint    `json:"category_id"`
	ImageURL    string  `json:"image_url"`
	Price       float64 `json:"price"`
	//UserID      uint      `json:"user_id"`  Will be refactored when User microservice would be used
	Stock     uint      `json:"stock"`
	CreatedAt time.Time `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:current_timestamp"`
}
