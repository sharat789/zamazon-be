package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sharat789/zamazon-be-ms/common/auth"
	"github.com/sharat789/zamazon-be-ms/transactions/configs"
	"github.com/sharat789/zamazon-be-ms/transactions/pkg/payment"
	"gorm.io/gorm"
)

type RestHandler struct {
	App           *fiber.App
	DB            *gorm.DB
	PaymentClient payment.PaymentClient
	Config        configs.AppConfig
	Auth          auth.Auth
}
