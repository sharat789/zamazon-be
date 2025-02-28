package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sharat789/zamazon-be/internal/api/rest"
	"github.com/sharat789/zamazon-be/internal/dto"
	"github.com/sharat789/zamazon-be/internal/repository"
	"github.com/sharat789/zamazon-be/internal/service"
	"log"
	"net/http"
	"strconv"
)

type UserHandler struct {
	userService service.UserService
}

func SetupUserRoutes(rh *rest.RestHandler) {
	app := rh.App
	svc := service.UserService{
		Repo:        repository.NewUserRepository(rh.DB),
		CatalogRepo: repository.NewCatalogRepository(rh.DB),
		Auth:        rh.Auth,
	}
	handler := UserHandler{
		svc,
	}
	publicRoutes := app.Group("/users")
	//public endpoints
	publicRoutes.Post("/register", handler.RegisterUser)
	publicRoutes.Post("/login", handler.Login)

	privateRoutes := publicRoutes.Group("/", rh.Auth.AuthorizeUser)
	//private endpoints
	privateRoutes.Post("/verifyUser", handler.VerifyUser)
	privateRoutes.Get("/verify", handler.GetVerificationCode)

	privateRoutes.Post("/userProfile", handler.CreateUserProfile)
	privateRoutes.Get("/userProfile", handler.GetUserProfile)
	privateRoutes.Patch("/userProfile", handler.UpdateUserProfile)

	privateRoutes.Post("/cart", handler.AddToCart)
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
	user := h.userService.Auth.GetCurrentUser(ctx)

	var req dto.VerificationCodeInput

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "invalid input",
		})
	}

	err := h.userService.VerifyCode(user.ID, req.Code)

	if err != nil {
		log.Printf("%v", err)
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "verified successfully",
	})
}
func (h *UserHandler) GetVerificationCode(ctx *fiber.Ctx) error {
	user := h.userService.Auth.GetCurrentUser(ctx)

	code, err := h.userService.GetVerificationCode(user)

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "unable to generate verification code",
			"error":   err.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Get verification code",
		"data":    code,
	})
}
func (h *UserHandler) GetUserProfile(ctx *fiber.Ctx) error {
	user := h.userService.Auth.GetCurrentUser(ctx)
	log.Println(user)

	profile, err := h.userService.GetUserProfile(user.ID)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "unable to fetch user's profile",
		})
	}
	return rest.SuccessResponse(ctx, "user profile found", profile)
}

func (h *UserHandler) UpdateUserProfile(ctx *fiber.Ctx) error {
	user := h.userService.Auth.GetCurrentUser(ctx)
	log.Println(user)

	req := dto.ProfileInput{}

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Please provide valid inputs",
		})
	}

	err := h.userService.UpdateProfile(user.ID, req)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "unable to update user's profile",
		})
	}

	return rest.SuccessResponse(ctx, "user profile updated", nil)
}

func (h *UserHandler) CreateUserProfile(ctx *fiber.Ctx) error {
	user := h.userService.Auth.GetCurrentUser(ctx)
	req := dto.ProfileInput{}

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Please provide valid inputs",
			"err":     err.Error(),
		})
	}
	log.Println(user)

	err := h.userService.CreateUserProfile(user.ID, req)

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "unable to create user's profile",
		})
	}
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Create user's profile",
	})
}
func (h *UserHandler) GetCart(ctx *fiber.Ctx) error {
	user := h.userService.Auth.GetCurrentUser(ctx)
	cart, _, err := h.userService.FindCart(user.ID)
	if err != nil {
		return rest.InternalErrorResponse(ctx, errors.New("unable to fetch cart"))
	}
	return rest.SuccessResponse(ctx, "cart found for user", cart)
}
func (h *UserHandler) AddToCart(ctx *fiber.Ctx) error {
	req := dto.CreateCartRequest{}
	err := ctx.BodyParser(&req)

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Please provide valid product and quantity",
		})
	}

	user := h.userService.Auth.GetCurrentUser(ctx)
	log.Println(user)

	cartItems, err := h.userService.CreateCart(req, user)

	if err != nil {
		return rest.InternalErrorResponse(ctx, err)
	}

	return rest.SuccessResponse(ctx, "cart created", cartItems)
}
func (h *UserHandler) GetOrders(ctx *fiber.Ctx) error {
	user := h.userService.Auth.GetCurrentUser(ctx)
	orders, err := h.userService.GetOrders(user)
	if err != nil {
		return rest.InternalErrorResponse(ctx, errors.New("unable to fetch orders"))
	}

	return rest.SuccessResponse(ctx, "orders found for user", orders)
}
func (h *UserHandler) GetOrderByID(ctx *fiber.Ctx) error {
	orderId, _ := strconv.Atoi(ctx.Params("id"))
	user := h.userService.Auth.GetCurrentUser(ctx)
	order, err := h.userService.GetOrderByID(uint(orderId), user.ID)
	if err != nil {
		return rest.InternalErrorResponse(ctx, errors.New("unable to fetch orders"))
	}

	return rest.SuccessResponse(ctx, "orders found for user", order)
}
func (h *UserHandler) becomeSeller(ctx *fiber.Ctx) error {
	user := h.userService.Auth.GetCurrentUser(ctx)
	req := dto.SellerInput{}

	err := ctx.BodyParser(&req)

	if err != nil {
		return ctx.Status(400).JSON(&fiber.Map{
			"message": "request parameters are not valid",
		})
	}

	token, err := h.userService.BecomeSeller(user.ID, req)

	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"message": "seller signup failed",
		})
	}
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Become seller",
		"token":   token,
	})
}
