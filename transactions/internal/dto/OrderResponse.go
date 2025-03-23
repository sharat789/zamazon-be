package dto

type OrderItem struct {
	ID        uint    `json:"id"`
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
	ImageURL  string  `json:"image_url"`
}

type OrderResponse struct {
	ID          string      `json:"id"`
	Date        string      `json:"date"`
	Status      string      `json:"status"`
	Total       float64     `json:"total"`
	Items       []OrderItem `json:"items"`
	PaymentID   string      `json:"payment_id"`
	PaymentType string      `json:"payment_type"`
}
