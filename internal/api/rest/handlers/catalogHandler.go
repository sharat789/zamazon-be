package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sharat789/zamazon-be/internal/api/rest"
	"github.com/sharat789/zamazon-be/internal/domain"
	"github.com/sharat789/zamazon-be/internal/dto"
	"github.com/sharat789/zamazon-be/internal/repository"
	"github.com/sharat789/zamazon-be/internal/service"
	"strconv"
)

type CatalogHandler struct {
	userService service.CatalogService
}

func SetupCatalogRoutes(rh *rest.RestHandler) {
	app := rh.App
	svc := service.CatalogService{
		Repo: repository.NewCatalogRepository(rh.DB),
		Auth: rh.Auth,
	}
	handler := CatalogHandler{
		svc,
	}
	//public endpoints for buyers
	app.Get("/products", handler.GetProducts)
	app.Get("/products/:id", handler.GetProductByID)
	app.Get("/categories", handler.GetCategories)
	app.Get("/categories/:id", handler.GetCategoryByID)

	//private endpoints
	sellerRoutes := app.Group("/seller", rh.Auth.AuthorizeSeller)
	sellerRoutes.Post("/categories", handler.CreateCategories)
	sellerRoutes.Patch("/categories/:id", handler.EditCategory)
	sellerRoutes.Delete("/categories/:id", handler.DeleteCategory)

	sellerRoutes.Get("/products", handler.GetProducts)
	sellerRoutes.Get("/products/:id", handler.GetProductByID)
	sellerRoutes.Post("/products", handler.CreateProducts)
	sellerRoutes.Put("/products/:id", handler.EditProduct)
	sellerRoutes.Patch("/products/:id", handler.UpdateStock)
	sellerRoutes.Delete("/products/:id", handler.DeleteProduct)

}

func (h CatalogHandler) GetCategories(ctx *fiber.Ctx) error {
	categories, err := h.userService.GetCategories()

	if err != nil {
		return rest.ErrorResponse(ctx, 404, err)
	}

	return rest.SuccessResponse(ctx, "categories", categories)
}

func (h CatalogHandler) GetCategoryByID(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))

	category, err := h.userService.GetCategory(uint(id))

	if err != nil {
		return rest.ErrorResponse(ctx, 404, err)
	}
	return rest.SuccessResponse(ctx, "category", category)
}

func (h CatalogHandler) CreateCategories(ctx *fiber.Ctx) error {
	req := dto.CreateCategoryRequest{}

	err := ctx.BodyParser(&req)

	if err != nil {
		return rest.BadRequestErrorResponse(ctx, "category request is invalid")
	}

	err = h.userService.CreateCategory(req)

	if err != nil {
		return rest.InternalErrorResponse(ctx, err)
	}
	return rest.SuccessResponse(ctx, "create category", nil)
}

func (h CatalogHandler) EditCategory(ctx *fiber.Ctx) error {
	req := dto.CreateCategoryRequest{}

	err := ctx.BodyParser(&req)

	if err != nil {
		return rest.BadRequestErrorResponse(ctx, "update category request is invalid")
	}

	id, _ := strconv.Atoi(ctx.Params("id"))
	updatedCategory, err := h.userService.EditCategory(uint(id), req)

	if err != nil {
		return rest.InternalErrorResponse(ctx, err)
	}
	return rest.SuccessResponse(ctx, "create category", updatedCategory)
}

func (h CatalogHandler) DeleteCategory(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	err := h.userService.DeleteCategory(uint(id))

	if err != nil {
		return rest.InternalErrorResponse(ctx, err)
	}
	return rest.SuccessResponse(ctx, "delete category success", nil)
}

func (h CatalogHandler) CreateProducts(ctx *fiber.Ctx) error {
	req := dto.CreateProductRequest{}
	err := ctx.BodyParser(&req)

	if err != nil {
		return rest.BadRequestErrorResponse(ctx, "product request is invalid")
	}

	user := h.userService.Auth.GetCurrentUser(ctx)
	err = h.userService.CreateProduct(req, user)

	if err != nil {
		return rest.InternalErrorResponse(ctx, err)
	}

	return rest.SuccessResponse(ctx, "create product", nil)
}

func (h CatalogHandler) EditProduct(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	req := dto.CreateProductRequest{}
	err := ctx.BodyParser(&req)

	if err != nil {
		return rest.BadRequestErrorResponse(ctx, "product request is invalid")
	}
	user := h.userService.Auth.GetCurrentUser(ctx)
	product, err := h.userService.EditProduct(uint(id), req, user)

	if err != nil {
		return rest.InternalErrorResponse(ctx, err)
	}
	return rest.SuccessResponse(ctx, "edit product", product)
}

func (h CatalogHandler) DeleteProduct(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	err := h.userService.DeleteProduct(id)
	return rest.SuccessResponse(ctx, "delete product", err)
}

func (h CatalogHandler) UpdateStock(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	req := dto.UpdateStockRequest{}
	err := ctx.BodyParser(&req)

	if err != nil {
		return rest.BadRequestErrorResponse(ctx, "update stock request is invalid")
	}

	user := h.userService.Auth.GetCurrentUser(ctx)
	product := domain.Product{
		ID:     uint(id),
		Stock:  uint(req.Stock),
		UserID: user.ID,
	}

	updatedProduct, err := h.userService.UpdateProductStock(product)

	return rest.SuccessResponse(ctx, "update stock", updatedProduct)
}

func (h CatalogHandler) GetProductByID(ctx *fiber.Ctx) error {

	id, _ := strconv.Atoi(ctx.Params("id"))

	product, err := h.userService.GetProductByID(uint(id))

	if err != nil {
		return rest.ErrorResponse(ctx, 404, err)
	}
	return rest.SuccessResponse(ctx, "get product by id", product)
}

func (h CatalogHandler) GetProducts(ctx *fiber.Ctx) error {

	products, err := h.userService.GetProducts()

	if err != nil {
		return rest.ErrorResponse(ctx, 404, err)
	}

	return rest.SuccessResponse(ctx, "get products", products)
}
