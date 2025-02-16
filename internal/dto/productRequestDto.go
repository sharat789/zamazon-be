package dto

type CreateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CategoryID  uint    `json:"category_id"`
	ImageURL    string  `json:"image_url"`
	Stock       int     `json:"stock"`
}

type UpdateStockRequest struct {
	Stock int `json:"stock"`
}
