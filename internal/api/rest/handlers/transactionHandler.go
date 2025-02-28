package handlers

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sharat789/zamazon-be/configs"
	"github.com/sharat789/zamazon-be/internal/api/rest"
	"github.com/sharat789/zamazon-be/internal/dto"
	"github.com/sharat789/zamazon-be/internal/helper"
	"github.com/sharat789/zamazon-be/internal/repository"
	"github.com/sharat789/zamazon-be/internal/service"
	"github.com/sharat789/zamazon-be/pkg/payment"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type TransactionHandler struct {
	transactionService service.TransactionService
	userService        service.UserService
	paymentClient      payment.PaymentClient
	Config             configs.AppConfig
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
		rh.Config,
	}

	secRoute := app.Group("/buyer", rh.Auth.AuthorizeUser)
	secRoute.Get("/payment", handler.MakePayment)
	secRoute.Get("/verify", handler.VerifyPayment)

	sellerRoute := app.Group("/seller", rh.Auth.AuthorizeSeller)
	sellerRoute.Get("/orders", handler.GetOrders)
	sellerRoute.Get("/order/:id", handler.GetOrder)
}

func (h *TransactionHandler) MakePayment(ctx *fiber.Ctx) error {
	user := h.transactionService.Auth.GetCurrentUser(ctx)
	pubKey := h.Config.PubKey

	// Check if user has an active payment
	activePayment, err := h.transactionService.GetActivePayment(user.ID)
	if activePayment.ID > 0 {
		return ctx.Status(http.StatusOK).JSON(&fiber.Map{
			"message": "Payment already initiated",
			"pub_key": pubKey,
			"secret":  activePayment.ClientSecret,
		})
	}
	_, amount, err := h.userService.FindCart(user.ID)

	orderId, err := helper.GenerateRandom(8)
	if err != nil {
		return rest.InternalErrorResponse(ctx, errors.New("could not generate order id"))
	}

	// Create payment session
	paymentResult, err := h.paymentClient.CreatePayment(amount, user.ID, orderId)
	if err != nil {
		return rest.ErrorResponse(ctx, 400, err)
	}
	// Store payment details
	err = h.transactionService.StoreCreatedPayment(dto.CreatePaymentRequest{
		UserId:       user.ID,
		Amount:       amount,
		OrderId:      orderId,
		ClientSecret: paymentResult.ClientSecret,
		PaymentId:    paymentResult.ID,
	})
	if err != nil {
		return rest.ErrorResponse(ctx, 400, err)
	}

	return ctx.Status(200).JSON(&fiber.Map{
		"message": "Payment initiated",
		"pub_key": pubKey,
		"secret":  paymentResult.ClientSecret,
	})
}

func (h *TransactionHandler) VerifyPayment(ctx *fiber.Ctx) error {
	user := h.transactionService.Auth.GetCurrentUser(ctx)

	activePayment, err := h.transactionService.GetActivePayment(user.ID)
	if err != nil || activePayment.ID == 0 {
		return rest.ErrorResponse(ctx, 404, errors.New("no active payment found"))
	}

	paymentResponse, _ := h.paymentClient.GetPaymentStatus(activePayment.PaymentId)
	paymentJSON, _ := json.Marshal(paymentResponse)
	paymentLogs := string(paymentJSON)

	paymentStatus := "failed"
	if paymentResponse.Status == "succeeded" {

		paymentStatus = "success"
	}

	err = h.transactionService.UpdatePayment(user.ID, paymentStatus, paymentLogs)
	if err != nil {
		return err
	}

	return ctx.Status(200).JSON(&fiber.Map{
		"message":  "Payment initiated",
		"response": paymentResponse,
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
