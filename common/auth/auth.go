package auth

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

// UserRole defines available user roles
const (
	ROLE_SELLER = "seller"
	ROLE_BUYER  = "buyer"
)

// TokenUser represents minimal user information needed for authentication
type TokenUser struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	UserRole string `json:"user_role"`
}

type Auth struct {
	Secret string
}

func SetupAuth(secret string) Auth {
	return Auth{
		secret,
	}
}

func (a Auth) CreateHashPassword(p string) (string, error) {
	if len(p) < 8 {
		return "", errors.New("length of the password must be at least 8 characters")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("hashing failed")
	}

	return string(hashPassword), nil
}

func (a Auth) GenerateToken(id uint, email string, role string) (string, error) {
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

func (a Auth) VerifyPassword(plainPassword string, hashedPassword string) error {
	if len(plainPassword) < 8 {
		return errors.New("password length should be at least 8 characters long")
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		return errors.New("password does not match")
	}

	return nil
}

func (a Auth) VerifyToken(token string) (TokenUser, error) {
	tokenArray := strings.Split(token, " ")
	if len(tokenArray) != 2 {
		return TokenUser{}, errors.New("invalid token format")
	}

	if tokenArray[0] != "Bearer" {
		return TokenUser{}, errors.New("invalid token type")
	}

	tokenString := tokenArray[1]

	t, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
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

func (a Auth) AuthorizeUser(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")

	user, err := a.VerifyToken(authHeader)
	if err == nil && user.ID > 0 {
		ctx.Locals("user", user)
		return ctx.Next()
	} else {
		return ctx.Status(401).JSON(&fiber.Map{
			"message": "authorization failed",
			"reason":  err,
		})
	}
}

func (a Auth) GetCurrentUser(ctx *fiber.Ctx) TokenUser {
	user := ctx.Locals("user")
	return user.(TokenUser)
}

func (a Auth) GenerateCode() (string, error) {
	// Assuming GenerateRandom is implemented elsewhere in the package
	return GenerateRandom(8)
}

func (a Auth) AuthorizeByRole(ctx *fiber.Ctx, requiredRole string) error {
	authHeader := ctx.Get("Authorization")

	user, err := a.VerifyToken(authHeader)

	if err != nil {
		return ctx.Status(401).JSON(&fiber.Map{
			"message": "authorization failed",
			"reason":  err,
		})
	} else if user.ID > 0 && user.UserRole == requiredRole {
		ctx.Locals("user", user)
		return ctx.Next()
	} else {
		return ctx.Status(401).JSON(&fiber.Map{
			"message": "authorization failed",
			"reason":  errors.New("insufficient permissions"),
		})
	}
}
