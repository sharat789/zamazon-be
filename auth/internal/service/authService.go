// File: auth/internal/service/authService.go
package service

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strings"
	"time"
)

const (
	ROLE_SELLER = "seller"
	ROLE_BUYER  = "buyer"
)

type TokenUser struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	UserRole string `json:"user_role"`
}

type AuthService struct {
	Secret string
}

func NewAuthService(secret string) *AuthService {
	return &AuthService{
		Secret: secret,
	}
}

func (a *AuthService) CreateHashPassword(password string) (string, error) {
	if len(password) < 8 {
		return "", errors.New("length of the password must be at least 8 characters")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("hashing failed")
	}

	return string(hashPassword), nil
}

func (a *AuthService) GenerateToken(id uint, email string, role string) (string, error) {
	if id == 0 || email == "" || role == "" {
		return "", errors.New("invalid user information for token generation")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"email":   email,
		"role":    role,
		"expiry":  time.Now().Add(time.Hour * 24 * 15).Unix(),
	})

	tokenString, err := token.SignedString([]byte(a.Secret))
	if err != nil {
		return "", errors.New("unable to get signed token")
	}
	return tokenString, nil
}

func (a *AuthService) VerifyPassword(plainPassword string, hashedPassword string) error {
	if len(plainPassword) < 8 {
		return errors.New("password length should be at least 8 characters long")
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		return errors.New("password does not match")
	}

	return nil
}

func (a *AuthService) VerifyToken(tokenString string) (TokenUser, error) {
	tokenArray := strings.Split(tokenString, " ")
	if len(tokenArray) != 2 {
		return TokenUser{}, errors.New("invalid token format")
	}

	if tokenArray[0] != "Bearer" {
		return TokenUser{}, errors.New("invalid token type")
	}

	t, err := jwt.Parse(tokenArray[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.Secret), nil
	})

	if err != nil {
		return TokenUser{}, fmt.Errorf("token parsing error: %v", err)
	}

	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		expiry, ok := claims["expiry"].(float64)
		if !ok {
			return TokenUser{}, errors.New("invalid expiry claim")
		}
		if time.Now().Unix() > int64(expiry) {
			return TokenUser{}, errors.New("token has expired")
		}
		user := TokenUser{
			ID:       uint(claims["user_id"].(float64)),
			Email:    claims["email"].(string),
			UserRole: claims["role"].(string),
		}
		return user, nil
	}

	return TokenUser{}, errors.New("token verification failed")
}

func (a *AuthService) GenerateCode() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 8

	result := make([]byte, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := range result {
		result[i] = charset[r.Intn(len(charset))]
	}

	return string(result), nil
}

func (a *AuthService) AuthorizeByRole(user TokenUser, requiredRole string) error {
	if user.ID == 0 {
		return errors.New("user not authenticated")
	}

	if user.UserRole != requiredRole {
		return errors.New("insufficient permissions")
	}

	return nil
}
