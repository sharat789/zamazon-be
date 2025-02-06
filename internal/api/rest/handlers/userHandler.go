package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sharat789/zamazon-be/internal/api/rest"
	"github.com/sharat789/zamazon-be/internal/dto"
	"github.com/sharat789/zamazon-be/internal/repository"
	"github.com/sharat789/zamazon-be/internal/service"
	"net/http"
)

type UserHandler struct {
	userService service.UserService
}

func SetupUserRoutes(rh *rest.RestHandler) {
	app := rh.App
	svc := service.UserService{
		Repo: repository.NewUserRepository(rh.DB),
		Auth: rh.Auth,
	}
	handler := UserHandler{
		svc,
	}
	publicRoutes := app.Group("/users")
	//public endpoints
	publicRoutes.Post("/registerUser", handler.RegisterUser)
	publicRoutes.Post("/login", handler.Login)

	privateRoutes := publicRoutes.Group("/", rh.Auth.AuthorizeUser)
	//private endpoints
	privateRoutes.Post("/verifyUser", handler.VerifyUser)
	privateRoutes.Get("/verify", handler.GetVerificationCode)

	privateRoutes.Post("/userProfile", handler.CreateUserProfile)
	privateRoutes.Get("/userProfile", handler.GetUserProfile)

	privateRoutes.Post("/cart", handler.CreateCart)
	privateRoutes.Get("/cart", handler.GetCart)

	privateRoutes.Get("/order", handler.GetOrders)
	privateRoutes.Get("/order/:id", handler.GetOrderByID)

	privateRoutes.Post("become-seller", handler.becomeSeller)
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
		"message": user.Email,
		"token":   token,
	})
}
func (h *UserHandler) Login(ctx *fiber.Ctx) error {
	loginInput := dto.UserLogin{}
	err := ctx.BodyParser(&loginInput)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Please provide valid inputs",
		})
	}
	token, err := h.userService.Login(loginInput.Email, loginInput.Password)

	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": loginInput.Email,
		"token":   token,
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
