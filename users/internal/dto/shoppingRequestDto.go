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

type OrderItem struct {
	ID        uint   `json:"id"`
	OrderID   uint   `json:"order_id"`
	ProductID uint   `json:"product_id"`
	Name      string `json:"name"`
	ImageURL  string `json:"image_url"`
	SellerId  uint   `json:"seller_id"`
	Price     uint   `json:"price"`
	Qty       uint   `json:"qty"`
}

type CreateOrderRequest struct {
	UserID         uint    `json:"user_id"`
	Amount         float64 `json:"amount"`
	OrderRefNumber string  `json:"order_ref_number"`
	PaymentId      string  `json:"payment_id"`
}
