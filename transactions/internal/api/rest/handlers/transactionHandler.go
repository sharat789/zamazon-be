package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sharat789/zamazon-be-ms/common/auth"
	"github.com/sharat789/zamazon-be-ms/transactions/configs"
	"github.com/sharat789/zamazon-be-ms/transactions/internal/api/rest"
	"github.com/sharat789/zamazon-be-ms/transactions/internal/dto"
	"github.com/sharat789/zamazon-be-ms/transactions/internal/repository"
	"github.com/sharat789/zamazon-be-ms/transactions/internal/service"
	"github.com/sharat789/zamazon-be-ms/transactions/pkg/payment"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
)

type TransactionHandler struct {
	transactionService service.TransactionService
	paymentClient      payment.PaymentClient
	Config             configs.AppConfig
	userServiceURL     string
}

func initialiseTransactionService(db *gorm.DB, auth auth.Auth) service.TransactionService {
	return service.TransactionService{
		Repo: repository.NewTransactionRepository(db),
		Auth: auth,
	}
}

func SetupTransactionRoutes(rh *rest.RestHandler) {
	app := rh.App
	svc := initialiseTransactionService(rh.DB, rh.Auth)
	handler := TransactionHandler{
		svc,
		rh.PaymentClient,
		rh.Config,
		rh.Config.UserServiceURL,
	}
	pubRoutes := app.Group("/buyer")
	pubRoutes.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
	secRoute := pubRoutes.Group("/", rh.Auth.AuthorizeUser)
	secRoute.Get("/verify", handler.VerifyPayment)
	secRoute.Get("/checkout", handler.CreateCheckoutSession)
	secRoute.Get("/orders", handler.GetOrders)
	secRoute.Get("/order/:id", handler.GetOrder)
}

// Helper method to call user service APIs
func (h *TransactionHandler) callUserService(method, endpoint string, body interface{}, token string) (*http.Response, error) {
	var req *http.Request
	var err error

	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest(method, h.userServiceURL+endpoint, bytes.NewBuffer(jsonData))
	} else {
		req, err = http.NewRequest(method, h.userServiceURL+endpoint, nil)
	}

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", token)
	}

	client := &http.Client{}
	return client.Do(req)
}

// FindCart calls the user service to get cart items
func (h *TransactionHandler) findCart(token string) ([]dto.CartItem, float64, error) {
	resp, err := h.callUserService("GET", fmt.Sprintf("/users/cart"), nil, token)
	if err != nil {
		return nil, 0, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error while closing response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, 0, errors.New("failed to get cart from user service")
	}

	var response struct {
		Data    []dto.CartItem `json:"data"`
		Message string         `json:"message"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, 0, err
	}

	totalAmount := 0.0
	for _, item := range response.Data {
		totalAmount += float64(item.Quantity) * item.Price
	}

	return response.Data, totalAmount, nil
}

func (h *TransactionHandler) createOrder(request dto.CreateOrderRequest, token string) error {
	// Validate request before sending
	if request.UserID == 0 || request.OrderRefNumber == "" || request.Amount <= 0 {
		return fmt.Errorf("incomplete order request: userID=%d, orderRef=%s, amount=%.2f",
			request.UserID, request.OrderRefNumber, request.Amount)
	}

	endpoint := "/users/order"

	// Debug log to see what we're sending
	requestJSON, _ := json.Marshal(request)
	log.Printf("Sending order request: %s", string(requestJSON))

	resp, err := h.callUserService("POST", endpoint, request, token)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Log the actual status and response body for debugging
	bodyBytes, _ := io.ReadAll(resp.Body)
	log.Printf("Create order response: %d %s - Body: %s",
		resp.StatusCode, resp.Status, string(bodyBytes))

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create order in user service: %s", resp.Status)
	}

	return nil
}

func (h *TransactionHandler) CreateCheckoutSession(ctx *fiber.Ctx) error {
	user := h.transactionService.Auth.GetCurrentUser(ctx)
	token := ctx.Get("Authorization")

	// Get cart items and total amount using user service HTTP call
	cartItems, totalAmount, err := h.findCart(token)
	if err != nil {
		log.Printf("Error while fetching cart items: %v", err)
		return rest.InternalErrorResponse(ctx, errors.New("unable to fetch cart items"))
	}

	if len(cartItems) == 0 {
		return rest.ErrorResponse(ctx, 400, errors.New("cart is empty"))
	}

	// Generate order ID
	orderId, err := auth.GenerateRandom(8)
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
		ClientSecret: "",
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
	token := ctx.Get("Authorization")
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
			request := dto.CreateOrderRequest{
				UserID:         user.ID,
				OrderRefNumber: payment.OrderId,
				PaymentID:      payment.PaymentId,
				Amount:         payment.Amount,
			}
			log.Printf("Creating order: UserID=%d, OrderRef=%s, PaymentID=%s, Amount=%.2f",
				user.ID, payment.OrderId, payment.PaymentId, payment.Amount)
			err = h.createOrder(request, token)
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
	return nil
}

func (h *TransactionHandler) GetOrders(ctx *fiber.Ctx) error {
	user := h.transactionService.Auth.GetCurrentUser(ctx)
	token := ctx.Get("Authorization")
	log.Println(user)
	resp, err := h.callUserService("GET", "/orders", nil, token)
	if err != nil {
		return rest.InternalErrorResponse(ctx, errors.New("failed to fetch orders"))
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	var response struct {
		Data []dto.OrderResponse `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return rest.InternalErrorResponse(ctx, errors.New("invalid response from user service"))
	}

	return rest.SuccessResponse(ctx, "orders", response.Data)
}

func (h *TransactionHandler) GetOrder(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	token := ctx.Get("Authorization")

	// Call user service to get order details
	resp, err := h.callUserService("GET", "/orders/"+id, nil, token)
	if err != nil {
		return rest.InternalErrorResponse(ctx, errors.New("failed to fetch order details"))
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	// Parse the response
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return rest.InternalErrorResponse(ctx, errors.New("invalid response from user service"))
	}

	return rest.SuccessResponse(ctx, "order", response["data"])
}
