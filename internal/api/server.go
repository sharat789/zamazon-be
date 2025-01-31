package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sharat789/zamazon-be/configs"
	"github.com/sharat789/zamazon-be/internal/api/rest"
	"github.com/sharat789/zamazon-be/internal/api/rest/handlers"
)

func StartServer(cfg configs.AppConfig) {
	app := fiber.New()
	rh := &rest.RestHandler{
		app,
	}
	SetupRoutes(rh)
	app.Listen(cfg.Port)
}

func SetupRoutes(rh *rest.RestHandler) {
	//user route handler
	handlers.SetupUserRoutes(rh)
	//transaction route handler
	//catalog route handler
}
