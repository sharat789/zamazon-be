package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sharat789/zamazon-be/internal/api/rest"
	"github.com/sharat789/zamazon-be/internal/dto"
	"github.com/sharat789/zamazon-be/internal/service"
	"net/http"
)

type UserHandler struct {
	userService service.UserService
}

func SetupUserRoutes(rh *rest.RestHandler) {
	app := rh.App
	svc := service.UserService{}
	handler := UserHandler{
		svc,
	}
	//public endpoints
	app.Post("/registerUser", handler.RegisterUser)
	app.Post("/login", handler.Login)

	//private endpoints
	app.Post("/verifyUser", handler.VerifyUser)
	app.Get("/verify", handler.GetVerificationCode)

	app.Post("/userProfile", handler.CreateUserProfile)
	app.Get("/userProfile", handler.GetUserProfile)

	app.Post("/cart", handler.CreateCart)
	app.Get("/cart", handler.GetCart)

	app.Get("/order", handler.GetOrders)
	app.Get("/order/:id", handler.GetOrderByID)

	app.Post("become-seller", handler.becomeSeller)
}
func (h *UserHandler) RegisterUser(ctx *fiber.Ctx) error {
	user := dto.UserSignup{}
	err := ctx.BodyParser(&user)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Please provide valid inputs",
		})
	}
	token, err := h.userService.UserSignup(user)

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "User signup failed",
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": token,
	})
}
func (h *UserHandler) Login(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Login",
	})
}
func (h *UserHandler) VerifyUser(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Verify user",
	})
}
func (h *UserHandler) GetVerificationCode(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Get verification code",
	})
}
func (h *UserHandler) GetUserProfile(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Get user's profile",
	})
}
func (h *UserHandler) CreateUserProfile(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Create user's profile",
	})
}
func (h *UserHandler) GetCart(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Get user's cart",
	})
}
func (h *UserHandler) CreateCart(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Create new cart",
	})
}
func (h *UserHandler) GetOrders(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Get user's orders",
	})
}
func (h *UserHandler) GetOrderByID(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Get order by id",
	})
}
func (h *UserHandler) becomeSeller(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Become seller",
	})
}
