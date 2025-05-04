package dto

type ProductResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"image_url"`
	Description string  `json:"description"`
	UserID      uint    `json:"user_id"`
	CategoryID  uint    `json:"category_id"`
	Stock       uint    `json:"stock"`
	CreatedAt   string  `json:"created_at,omitempty"`
	UpdatedAt   string  `json:"updated_at,omitempty"`
}
