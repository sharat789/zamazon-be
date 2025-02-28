package domain

import "time"

type Payment struct {
	ID            uint      `json:"id" gorm:"PrimaryKey"`
	UserId        uint      `json:"user_id"`
	CaptureMethod string    `json:"capture_method"`
	Amount        float64   `json:"amount"`
	OrderId       string    `json:"order_id"`
	CustomerId    string    `json:"customer_id"`
	PaymentId     string    `json:"payment_id"`
	ClientSecret  string    `json:"client_secret"`
	Status        string    `json:"status" gorm:"default:initial"`
	Response      string    `json:"response"`
	CreatedAt     time.Time `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"default:current_timestamp"`
}

type PaymentStatus string

const (
	PaymentStatusInitial PaymentStatus = "initial"
	PaymentStatusSuccess PaymentStatus = "success"
	PaymentStatusFailed  PaymentStatus = "failed"
	PaymentStatusPending PaymentStatus = "pending"
)
