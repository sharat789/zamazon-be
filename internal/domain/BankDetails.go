package domain

import "time"

type BankDetails struct {
	ID          uint      `json:"id" gorm:"PrimaryKey"`
	UserId      uint      `json:"user_id"`
	AccountNo   uint      `json:"account_no" gorm:"index;unique;not null"`
	SwiftCode   string    `json:"swift_code"`
	PaymentType string    `json:"payment_type"`
	CreatedAt   time.Time `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"default:current_timestamp"`
}
