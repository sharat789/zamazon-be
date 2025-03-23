package domain

import "time"

type Category struct {
	ID           uint      `json:"id" gorm:"PrimaryKey"`
	Name         string    `json:"name" gorm:"index;"`
	ParentID     int       `json:"parent_id"`
	ImageURL     string    `json:"image_url"`
	Products     []Product `json:"products"`
	DisplayOrder int       `json:"display_order"`
	CreatedAt    time.Time `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"default:current_timestamp"`
}
