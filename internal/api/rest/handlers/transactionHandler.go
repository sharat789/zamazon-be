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
	secRoute.Get("/checkout", handler.CreateCheckoutSession)
	secRoute.Get("/orders", handler.GetOrders)
	secRoute.Get("/order/:id", handler.GetOrder)

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

func (h *TransactionHandler) CreateCheckoutSession(ctx *fiber.Ctx) error {
	user := h.transactionService.Auth.GetCurrentUser(ctx)

	// Get cart items and total amount
	cartItems, totalAmount, err := h.userService.FindCart(user.ID)
	if err != nil {
		return rest.InternalErrorResponse(ctx, errors.New("unable to fetch cart items"))
	}

	if len(cartItems) == 0 {
		return rest.ErrorResponse(ctx, 400, errors.New("cart is empty"))
	}

	// Generate order ID
	orderId, err := helper.GenerateRandom(8)
	if err != nil {
		return rest.InternalErrorResponse(ctx, errors.New("could not generate order id"))
	}

	// Create checkout session with Stripe
	checkoutSession, err := h.paymentClient.CreateCheckoutSession(totalAmount, user.ID, orderId)
	if err != nil {
		return rest.ErrorResponse(ctx, 400, err)
	}

	// Store payment details in a pending state
	err = h.transactionService.StoreCreatedPayment(dto.CreatePaymentRequest{
		UserId:       user.ID,
		Amount:       totalAmount,
		OrderId:      orderId,
		ClientSecret: "", // Not applicable for Checkout
		PaymentId:    checkoutSession.ID,
		PaymentType:  "checkout",
	})

	if err != nil {
		return rest.ErrorResponse(ctx, 400, err)
	}

	return ctx.Status(200).JSON(&fiber.Map{
		"message":      "Checkout session created",
		"session_id":   checkoutSession.ID,
		"checkout_url": checkoutSession.URL,
	})
}

func (h *TransactionHandler) VerifyPayment(ctx *fiber.Ctx) error {
	user := h.transactionService.Auth.GetCurrentUser(ctx)
	sessionId := ctx.Query("session_id")

	// If session_id is provided, verify checkout session
	if sessionId != "" {
		session, err := h.paymentClient.GetCheckoutSession(sessionId)
		if err != nil {
			return rest.ErrorResponse(ctx, 400, errors.New("invalid checkout session"))
		}

		// Check if session is paid
		if session.PaymentStatus == "paid" {
			// Get payment from database using session ID as payment ID
			payment, err := h.transactionService.GetPaymentByID(sessionId)
			if err != nil || payment.ID == 0 {
				return rest.ErrorResponse(ctx, 404, errors.New("payment not found"))
			}

			// Check if the payment has already been verified
			if payment.Status == "success" {
				return ctx.Status(200).JSON(&fiber.Map{
					"message":  "Payment already verified and order created",
					"order_id": payment.OrderId,
				})
			}

			// Create order from cart
			err = h.userService.CreateOrder(user.ID, payment.OrderId, payment.PaymentId, payment.Amount)
			if err != nil {
				return rest.InternalErrorResponse(ctx, err)
			}

			// Update payment status
			sessionJSON, _ := json.Marshal(session)
			sessionLogs := string(sessionJSON)
			err = h.transactionService.UpdatePayment(user.ID, "success", sessionLogs)
			if err != nil {
				return rest.InternalErrorResponse(ctx, err)
			}

			return ctx.Status(200).JSON(&fiber.Map{
				"message":  "Payment successful and order created",
				"order_id": payment.OrderId,
			})
		} else {
			return rest.ErrorResponse(ctx, 400, errors.New("payment not completed"))
		}
	}

	// For direct payment intents (existing flow)
	activePayment, err := h.transactionService.GetActivePayment(user.ID)
	if err != nil || activePayment.ID == 0 {
		return rest.ErrorResponse(ctx, 404, errors.New("no active payment found"))
	}

	paymentResponse, _ := h.paymentClient.GetPaymentStatus(activePayment.PaymentId)
	paymentJSON, _ := json.Marshal(paymentResponse)
	paymentLogs := string(paymentJSON)

	paymentStatus := "failed"
	if paymentResponse.Status == "succeeded" {
		err = h.userService.CreateOrder(user.ID, activePayment.OrderId, activePayment.PaymentId, activePayment.Amount)
		paymentStatus = "success"
	}
	if err != nil {
		return rest.InternalErrorResponse(ctx, err)
	}

	err = h.transactionService.UpdatePayment(user.ID, paymentStatus, paymentLogs)
	if err != nil {
		return err
	}

	return ctx.Status(200).JSON(&fiber.Map{
		"message": "Payment verified",
		"status":  paymentStatus,
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
