package domain

import "time"

const (
	SELLER = "seller"
	BUYER  = "buyer"
)

type User struct {
	ID               uint      `json:"id" gorm:"PrimaryKey"`
	FName            string    `json:"f_name"`
	LName            string    `json:"l_name"`
	Email            string    `json:"email" gorm:"index;unique;not null"`
	Phone            string    `json:"phone"`
	Password         string    `json:"password"`
	VerificationCode string    `json:"verificationCode"`
	Expiry           time.Time `json:"expiry"`
	Address          Address   `json:"address"`
	Payments         []Payment `json:"payment"`
	IsVerified       bool      `json:"isVerified" gorm:"default:false"`
	UserType         string    `json:"user_type" gorm:"default:buyer"`
	CreatedAt        time.Time `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"default:current_timestamp"`
}
