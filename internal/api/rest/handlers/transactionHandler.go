package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sharat789/zamazon-be/internal/api/rest"
	"github.com/sharat789/zamazon-be/internal/helper"
	"github.com/sharat789/zamazon-be/internal/repository"
	"github.com/sharat789/zamazon-be/internal/service"
	"github.com/sharat789/zamazon-be/pkg/payment"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

type TransactionHandler struct {
	transactionService service.TransactionService
	userService        service.UserService
	paymentClient      payment.PaymentClient
}

func initialiseTransactionService(db *gorm.DB, auth helper.Auth) service.TransactionService {
	return service.TransactionService{
		Repo: repository.NewTransactionRepository(db),
		Auth: auth,
	}
}

func SetupTransactionRoutes(rh *rest.RestHandler) {
	app := rh.App
	svc := initialiseTransactionService(rh.DB, rh.Auth)
	userService := service.UserService{
		Repo:        repository.NewUserRepository(rh.DB),
		CatalogRepo: repository.NewCatalogRepository(rh.DB),
		Auth:        rh.Auth,
	}
	handler := TransactionHandler{
		svc,
		userService,
		rh.PaymentClient,
	}

	secRoute := app.Group("/", rh.Auth.AuthorizeUser)
	secRoute.Get("/payment", handler.MakePayment)

	sellerRoute := app.Group("/seller", rh.Auth.AuthorizeSeller)
	sellerRoute.Get("/orders", handler.GetOrders)
	sellerRoute.Get("/order/:id", handler.GetOrder)
}

func (h *TransactionHandler) MakePayment(ctx *fiber.Ctx) error {
	user := h.transactionService.Auth.GetCurrentUser(ctx)

	// Check if user has an active payment
	activePayment, err := h.transactionService.GetActivePayment(user.ID)
	log.Println(activePayment)
	if activePayment != nil && activePayment.ID > 0 {
		return ctx.Status(http.StatusOK).JSON(&fiber.Map{
			"message":     "Payment already initiated",
			"payment_url": activePayment.PaymentUrl,
		})
	}
	_, amount, err := h.userService.FindCart(user.ID)

	orderId, err := helper.GenerateRandom(8)

	// Create payment session
	sessionResult, err := h.paymentClient.CreatePayment(amount, 123, orderId)

	// Store payment details
	err = h.transactionService.StoreCreatedPayment(user.ID, sessionResult, amount, orderId)
	if err != nil {
		return rest.ErrorResponse(ctx, 400, err)
	}

	return ctx.Status(200).JSON(&fiber.Map{
		"message":    "Payment initiated",
		"session":    sessionResult,
		"paymentUrl": sessionResult.URL,
	})
}

func (h *TransactionHandler) GetOrders(ctx *fiber.Ctx) error {
	user := h.transactionService.Auth.GetCurrentUser(ctx)
	orders, err := h.transactionService.GetOrders(user)
	if err != nil {
		return rest.ErrorResponse(ctx, 404, err)
	}
	return rest.SuccessResponse(ctx, "orders", orders)
}

func (h *TransactionHandler) GetOrder(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	user := h.transactionService.Auth.GetCurrentUser(ctx)
	order, err := h.transactionService.GetOrderDetails(uint(id), user)
	if err != nil {
		return rest.ErrorResponse(ctx, 404, err)
	}
	return rest.SuccessResponse(ctx, "order", order)
}
