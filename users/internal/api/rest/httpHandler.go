package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sharat789/zamazon-be-ms/users/configs"
	"gorm.io/gorm"
)

type RestHandler struct {
	App    *fiber.App
	DB     *gorm.DB
	Config configs.AppConfig
}
