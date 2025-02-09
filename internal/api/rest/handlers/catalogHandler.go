package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sharat789/zamazon-be/internal/api/rest"
	"github.com/sharat789/zamazon-be/internal/repository"
	"github.com/sharat789/zamazon-be/internal/service"
	"log"
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
	//publicRoutes := app.Group("/users")
	//public endpoints for buyers
	app.Get("/products")
	app.Get("/products/:id")
	app.Get("/categories")
	app.Get("/categories/:id")

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

func (h CatalogHandler) CreateCategories(ctx *fiber.Ctx) error {
	user := h.userService.Auth.GetCurrentUser(ctx)

	log.Printf("current user %v", user)
	return rest.SuccessResponse(ctx, "create category", nil)
}

func (h CatalogHandler) EditCategory(ctx *fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "edit category", nil)
}

func (h CatalogHandler) DeleteCategory(ctx *fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "delete category", nil)
}

func (h CatalogHandler) CreateProducts(ctx *fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "create product", nil)
}

func (h CatalogHandler) EditProduct(ctx *fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "edit product", nil)
}

func (h CatalogHandler) DeleteProduct(ctx *fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "delete product", nil)
}

func (h CatalogHandler) UpdateStock(ctx *fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "update stock", nil)
}

func (h CatalogHandler) GetProductByID(ctx *fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "get product by id", nil)
}

func (h CatalogHandler) GetProducts(ctx *fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "get products", nil)
}
