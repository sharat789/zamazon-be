package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sharat789/zamazon-be-ms/users/internal/api/middleware"
	"github.com/sharat789/zamazon-be-ms/users/internal/api/rest"
	"github.com/sharat789/zamazon-be-ms/users/internal/client"
	"github.com/sharat789/zamazon-be-ms/users/internal/dto"
	"github.com/sharat789/zamazon-be-ms/users/internal/repository"
	"github.com/sharat789/zamazon-be-ms/users/internal/service"
	"log"
	"net/http"
	"strconv"
)

type UserHandler struct {
	userService service.UserService
}

func SetupUserRoutes(rh *rest.RestHandler, catalogClient *client.CatalogClient, authClient *client.AuthClient) {
	app := rh.App
	svc := service.UserService{
		Repo:          repository.NewUserRepository(rh.DB),
		CatalogClient: catalogClient,
		AuthClient:    authClient,
	}
	handler := UserHandler{
		svc,
	}
	publicRoutes := app.Group("/users")
	//public endpoints
	publicRoutes.Post("/register", handler.RegisterUser)
	publicRoutes.Post("/login", handler.Login)
	publicRoutes.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
	privateRoutes := publicRoutes.Group("/", middleware.AuthorizeUser(authClient))
	//private endpoints
	privateRoutes.Post("/verifyUser", handler.VerifyUser)
	privateRoutes.Get("/verify", handler.GetVerificationCode)

	privateRoutes.Post("/profile", handler.CreateUserProfile)
	privateRoutes.Get("/profile", handler.GetUserProfile)
	privateRoutes.Patch("/profile", handler.UpdateUserProfile)

	privateRoutes.Post("/cart", handler.AddToCart)
	privateRoutes.Get("/cart", handler.GetCart)
	privateRoutes.Put("/cart/:productId", handler.UpdateProductQtyInCart)
	privateRoutes.Delete("/cart/:productID", handler.RemoveProductFromCart)
	privateRoutes.Delete("/cart", handler.ClearCart)

	privateRoutes.Get("/order", handler.GetOrders)
	privateRoutes.Get("/order/:id", handler.GetOrderByID)
	privateRoutes.Post("/order", handler.CreateOrder)
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
	user := h.userService.GetCurrentUser(ctx)

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
	tokenUser := h.userService.GetCurrentUser(ctx)

	code, err := h.userService.GetVerificationCode(tokenUser.ID)

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
	user := h.userService.GetCurrentUser(ctx)
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
	user := h.userService.GetCurrentUser(ctx)
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
	user := h.userService.GetCurrentUser(ctx)
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
	user := h.userService.GetCurrentUser(ctx)
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

	user := h.userService.GetCurrentUser(ctx)
	log.Println(user)

	cartItems, err := h.userService.CreateCart(req, user)

	if err != nil {
		return rest.InternalErrorResponse(ctx, err)
	}

	return rest.SuccessResponse(ctx, "cart created", cartItems)
}

func (h *UserHandler) UpdateProductQtyInCart(ctx *fiber.Ctx) error {
	productID, _ := strconv.Atoi(ctx.Params("productID"))
	req := dto.UpdateCartRequest{}

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "Please provide valid product and quantity",
		})
	}

	user := h.userService.GetCurrentUser(ctx)
	err := h.userService.UpdateProductQtyInCart(user.ID, uint(productID), req.Qty)

	if err != nil {
		return rest.InternalErrorResponse(ctx, err)
	}

	return rest.SuccessResponse(ctx, "cart updated", nil)
}

func (h *UserHandler) RemoveProductFromCart(ctx *fiber.Ctx) error {
	productID, _ := strconv.Atoi(ctx.Params("productID"))
	user := h.userService.GetCurrentUser(ctx)

	err := h.userService.RemoveProductFromCart(user.ID, uint(productID))
	if err != nil {
		return rest.InternalErrorResponse(ctx, err)
	}

	return rest.SuccessResponse(ctx, "product removed from cart", nil)
}

func (h *UserHandler) ClearCart(ctx *fiber.Ctx) error {
	user := h.userService.GetCurrentUser(ctx)

	err := h.userService.ClearCart(user.ID)
	if err != nil {
		return rest.InternalErrorResponse(ctx, err)
	}

	return rest.SuccessResponse(ctx, "cart cleared", nil)
}
func (h *UserHandler) GetOrders(ctx *fiber.Ctx) error {
	tokenUser := h.userService.GetCurrentUser(ctx)
	orders, err := h.userService.GetOrders(tokenUser.ID)
	if err != nil {
		return rest.InternalErrorResponse(ctx, errors.New("unable to fetch orders"))
	}

	return rest.SuccessResponse(ctx, "orders found for user", orders)
}

func (h *UserHandler) GetOrderByID(ctx *fiber.Ctx) error {
	orderId, _ := strconv.Atoi(ctx.Params("id"))
	user := h.userService.GetCurrentUser(ctx)
	order, err := h.userService.GetOrderByID(uint(orderId), user.ID)
	if err != nil {
		return rest.InternalErrorResponse(ctx, errors.New("unable to fetch orders"))
	}

	return rest.SuccessResponse(ctx, "orders found for user", order)
}

func (h *UserHandler) CreateOrder(ctx *fiber.Ctx) error {
	var request dto.CreateOrderRequest
	if err := ctx.BodyParser(&request); err != nil {
		return rest.ErrorResponse(ctx, http.StatusBadRequest, errors.New("invalid request format"))
	}

	// Validate required fields
	if request.UserID == 0 || request.OrderRefNumber == "" || request.Amount <= 0 {
		return rest.ErrorResponse(ctx, http.StatusBadRequest, errors.New("missing required order information"))
	}

	// Use the token user for authorization check
	tokenUser := h.userService.GetCurrentUser(ctx)

	// Additional security check - ensure the order is for the authorized user
	if tokenUser.ID != request.UserID && tokenUser.UserRole != "seller" {
		return rest.ErrorResponse(ctx, http.StatusForbidden, errors.New("unauthorized to create order for this user"))
	}

	// Create the order through service layer
	err := h.userService.CreateOrder(request)
	if err != nil {
		log.Printf("Failed to create order: %v", err)
		return rest.InternalErrorResponse(ctx, errors.New("failed to create order"))
	}

	// Clear the user's cart after successful order
	if err = h.userService.ClearCart(request.UserID); err != nil {
		log.Printf("Failed to clear cart after order creation: %v", err)
		// We don't want to fail the whole request if just cart clearing fails
	}

	return rest.SuccessResponse(ctx, "order created successfully", fiber.Map{
		"reference": request.OrderRefNumber,
	})
}

//func (h *UserHandler) becomeSeller(ctx *fiber.Ctx) error {
//	user := h.userService.GetCurrentUser(ctx)
//	req := dto.SellerInput{}
//
//	err := ctx.BodyParser(&req)
//
//	if err != nil {
//		return ctx.Status(400).JSON(&fiber.Map{
//			"message": "request parameters are not valid",
//		})
//	}
//
//	token, err := h.userService.BecomeSeller(user.ID, req)
//
//	if err != nil {
//		return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
//			"message": "seller signup failed",
//		})
//	}
//	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
//		"message": "Become seller",
//		"token":   token,
//	})
//}
