package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sharat789/zamazon-be-ms/users/configs"
	"github.com/sharat789/zamazon-be-ms/users/internal/api/metrics"
	"github.com/sharat789/zamazon-be-ms/users/internal/api/rest"
	"github.com/sharat789/zamazon-be-ms/users/internal/api/rest/handlers"
	"github.com/sharat789/zamazon-be-ms/users/internal/client"
	"github.com/sharat789/zamazon-be-ms/users/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func StartServer(cfg configs.AppConfig) {
	app := fiber.New()
	app.Use(metrics.PrometheusMiddleware())
	db, err := gorm.Open(postgres.Open(cfg.DataSourceName), &gorm.Config{})
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
	if err != nil {
		log.Fatalf("db conn error %v", err)
	}

	log.Println("db connected...")
	err = db.AutoMigrate(&domain.User{},
		&domain.Address{},
		&domain.Cart{},
		&domain.Order{},
		&domain.OrderItem{},
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

	catalogClient := client.NewCatalogClient(cfg.CatalogURL)
	authClient := client.NewAuthClient(cfg.AuthURL)
	rh := &rest.RestHandler{
		app,
		db,
		cfg,
	}

	SetupRoutes(rh, catalogClient, authClient)
	app.Listen(cfg.Port)
}

func SetupRoutes(rh *rest.RestHandler, catalogClient *client.CatalogClient, authClient *client.AuthClient) {
	handlers.SetupUserRoutes(rh, catalogClient, authClient)
}
