package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sharat789/zamazon-be/configs"
	"github.com/sharat789/zamazon-be/internal/api/rest"
	"github.com/sharat789/zamazon-be/internal/api/rest/handlers"
	"github.com/sharat789/zamazon-be/internal/domain"
	"github.com/sharat789/zamazon-be/internal/helper"
	"github.com/sharat789/zamazon-be/pkg/payment"
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
	err = db.AutoMigrate(&domain.User{},
		&domain.Address{},
		&domain.BankDetails{},
		&domain.Category{},
		&domain.Product{},
		&domain.Cart{},
		&domain.Order{},
		&domain.OrderItem{},
		&domain.Payment{},
	)

	if err != nil {
		log.Fatalf("error on migration %v", err)
	}

	log.Println("migration successful")
	auth := helper.Auth{}
	paymentClient := payment.NewPaymentClient(cfg.StripeSecret, cfg.SuccessURL, cfg.CancelURL)
	rh := &rest.RestHandler{
		app,
		db,
		auth,
		paymentClient,
	}

	SetupRoutes(rh)
	app.Listen(cfg.Port)
}

func SetupRoutes(rh *rest.RestHandler) {
	//user route handler
	handlers.SetupUserRoutes(rh)
	//transaction route handler
	handlers.SetupTransactionRoutes(rh)
	//catalog route handler
	handlers.SetupCatalogRoutes(rh)
}
