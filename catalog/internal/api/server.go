package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sharat789/zamazon-be-ms/catalog/configs"
	"github.com/sharat789/zamazon-be-ms/catalog/internal/api/rest"
	"github.com/sharat789/zamazon-be-ms/catalog/internal/api/rest/handlers"
	"github.com/sharat789/zamazon-be-ms/catalog/internal/domain"
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
	err = db.AutoMigrate(
		&domain.Category{},
		&domain.Product{},
	)

	if err != nil {
		log.Fatalf("error on migration %v", err)
	}

	log.Println("migration successful")

	c := cors.New(cors.Config{
		AllowOrigins: "http://localhost:4200, http://localhost:3030/",
		AllowHeaders: "Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	})

	app.Use(c)
	//auth := helper.Auth{}
	rh := &rest.RestHandler{
		app,
		db,
		cfg,
	}

	SetupRoutes(rh)
	app.Listen(cfg.Port)
}

func SetupRoutes(rh *rest.RestHandler) {
	handlers.SetupCatalogRoutes(rh)
}
