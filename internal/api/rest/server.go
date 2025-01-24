package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sharat789/zamazon-be/configs"
)

func StartServer(cfg configs.AppConfig) {
	app := fiber.New()
	app.Listen(cfg.Port)
}
