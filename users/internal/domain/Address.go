package domain

import "time"

type Address struct {
	ID           uint      `json:"id" gorm:"PrimaryKey"`
	AddressLine1 string    `json:"address_line1"`
	AddressLine2 string    `json:"address_line2"`
	City         string    `json:"city"`
	PostCode     uint      `json:"post_code"`
	Country      string    `json:"country"`
	UserID       uint      `json:"user_id"`
	CreatedAt    time.Time `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"default:current_timestamp"`
}
