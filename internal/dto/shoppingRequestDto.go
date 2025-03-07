package dto

type CreateCartRequest struct {
	ProductID uint `json:"product_id"`
	Qty       uint `json:"qty"`
}

type CreatePaymentRequest struct {
	UserId       uint    `json:"user_id"`
	Amount       float64 `json:"amount"`
	OrderId      string  `json:"order_id"`
	ClientSecret string  `json:"client_secret"`
	PaymentId    string  `json:"payment_id"`
	PaymentType  string  `json:"payment_type,omitempty"` // "intent" or "checkout"
}
