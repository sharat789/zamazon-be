package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sharat789/zamazon-be/configs"
	"github.com/sharat789/zamazon-be/internal/api/rest"
	"github.com/sharat789/zamazon-be/internal/api/rest/handlers"
	"github.com/sharat789/zamazon-be/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func StartServer(cfg configs.AppConfig) {
	app := fiber.New()

	db, err := gorm.Open(postgres.Open(cfg.DataSourceName), &gorm.Config{})

	if err != nil {
		log.Fatalf("db conn error %v", err)
	}

	log.Println("db connected...")
	db.AutoMigrate(&domain.User{})
	rh := &rest.RestHandler{
		app,
		db,
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
