package domain

import "time"

type User struct {
	ID               uint      `json:"id"`
	FName            string    `json:"f_name"`
	LName            string    `json:"l_name"`
	Email            string    `json:"email"`
	Phone            string    `json:"phone"`
	Password         string    `json:"password"`
	VerificationCode int       `json:"verificationCode"`
	Expiry           time.Time `json:"expiry"`
	IsVerified       bool      `json:"isVerified"`
	UserType         string    `json:"user_type"`
}
