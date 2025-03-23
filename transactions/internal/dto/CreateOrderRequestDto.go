package dto

type CreateOrderRequest struct {
	UserID         uint    `json:"user_id"`
	OrderRefNumber string  `json:"order_ref_number"` // Changed from "order_id"
	PaymentID      string  `json:"payment_id"`
	Amount         float64 `json:"amount"`
}
