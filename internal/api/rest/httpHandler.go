package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sharat789/zamazon-be/configs"
	"github.com/sharat789/zamazon-be/internal/helper"
	"github.com/sharat789/zamazon-be/pkg/payment"
	"gorm.io/gorm"
)

type RestHandler struct {
	App           *fiber.App
	DB            *gorm.DB
	Auth          helper.Auth
	PaymentClient payment.PaymentClient
	Config        configs.AppConfig
}
