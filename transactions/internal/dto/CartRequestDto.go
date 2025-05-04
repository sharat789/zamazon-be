package dto

import "time"

type CartItem struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	ProductID uint      `json:"product_id"`
	Name      string    `json:"name"`
	ImageURL  string    `json:"image_url"`
	SellerID  uint      `json:"seller_id"`
	Price     float64   `json:"price"`
	Quantity  int       `json:"qty"` // Note: JSON field is "qty"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CartResponse struct {
	Items []CartItem `json:"items"`
}
