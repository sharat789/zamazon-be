package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sharat789/zamazon-be/internal/api/rest"
	"github.com/sharat789/zamazon-be/internal/helper"
	"github.com/sharat789/zamazon-be/internal/repository"
	"github.com/sharat789/zamazon-be/internal/service"
	"gorm.io/gorm"
	"strconv"
)

type TransactionHandler struct {
	transactionService service.TransactionService
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
	handler := TransactionHandler{
		svc,
	}

	secRoute := app.Group("/", rh.Auth.AuthorizeUser)
	secRoute.Get("/payment", handler.makePayment)

	sellerRoute := app.Group("/seller", rh.Auth.AuthorizeSeller)
	sellerRoute.Get("/orders", handler.getOrders)
	sellerRoute.Get("/order/:id", handler.getOrder)
}

func (h *TransactionHandler) makePayment(ctx *fiber.Ctx) error {
	payload := struct {
		Message string `json:"message"`
	}{
		Message: "Payment successful",
	}
	return ctx.Status(200).JSON(payload)
}

func (h *TransactionHandler) getOrders(ctx *fiber.Ctx) error {
	user := h.transactionService.Auth.GetCurrentUser(ctx)
	orders, err := h.transactionService.GetOrders(user)
	if err != nil {
		return rest.ErrorResponse(ctx, 404, err)
	}
	return rest.SuccessResponse(ctx, "orders", orders)
}

func (h *TransactionHandler) getOrder(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	user := h.transactionService.Auth.GetCurrentUser(ctx)
	order, err := h.transactionService.GetOrderDetails(uint(id), user)
	if err != nil {
		return rest.ErrorResponse(ctx, 404, err)
	}
	return rest.SuccessResponse(ctx, "order", order)
}
