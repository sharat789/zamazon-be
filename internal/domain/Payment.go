package domain

import "time"

type Payment struct {
	ID            uint      `json:"id" gorm:"PrimaryKey"`
	UserId        uint      `json:"user_id"`
	CaptureMethod string    `json:"capture_method"`
	Amount        float64   `json:"amount"`
	TransactionId string    `json:"transaction_id"`
	CustomerId    string    `json:"customer_id"`
	PaymentId     string    `json:"payment_id"`
	Status        string    `json:"status"`
	Response      string    `json:"response"`
	CreatedAt     time.Time `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"default:current_timestamp"`
}
