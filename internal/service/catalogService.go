package service

import (
	"errors"
	"github.com/sharat789/zamazon-be/internal/domain"
	"github.com/sharat789/zamazon-be/internal/dto"
	"github.com/sharat789/zamazon-be/internal/helper"
	"github.com/sharat789/zamazon-be/internal/repository"
)

type CatalogService struct {
	Repo repository.CatalogRepository
	Auth helper.Auth
}

func (s CatalogService) CreateCategory(input dto.CreateCategoryRequest) error {

	err := s.Repo.CreateCategory(&domain.Category{
		Name:         input.Name,
		ImageURL:     input.ImageUrl,
		DisplayOrder: input.DisplayOrder,
	})
	return err
}

func (s CatalogService) EditCategory(id uint, input dto.CreateCategoryRequest) (*domain.Category, error) {
	existingCat, err := s.Repo.FindCategoryByID(id)

	if err != nil {
		return nil, errors.New("category does not exist")
	}

	if len(input.Name) > 0 {
		existingCat.Name = input.Name
	}
	if len(input.ImageUrl) > 0 {
		existingCat.ImageURL = input.ImageUrl
	}
	if input.ParentId > 0 {
		existingCat.ParentID = int(input.ParentId)
	}
	if input.DisplayOrder > 0 {
		existingCat.DisplayOrder = input.DisplayOrder
	}

	updatedCat, err := s.Repo.EditCategory(existingCat)

	return updatedCat, err
}

func (s CatalogService) DeleteCategory(id uint) error {
	err := s.Repo.DeleteCategory(id)

	if err != nil {
		return errors.New("category does not exist for deletion")
	}
	return nil
}

func (s CatalogService) GetCategory(id uint) (*domain.Category, error) {
	category, err := s.Repo.FindCategoryByID(id)

	if err != nil {
		return nil, err
	}
	return category, nil
}

func (s CatalogService) GetCategories() ([]*domain.Category, error) {
	categories, err := s.Repo.FindCategories()

	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (s CatalogService) CreateProduct(input dto.CreateProductRequest, user domain.User) error {

	err := s.Repo.CreateProduct(&domain.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		CategoryID:  input.CategoryID,
		UserID:      user.ID,
		ImageURL:    input.ImageURL,
		Stock:       uint(input.Stock),
	})
	return err
}

func (s CatalogService) EditProduct(id uint, input dto.CreateProductRequest, user domain.User) (*domain.Product, error) {
	existingProd, err := s.Repo.FindProductByID(id)

	if err != nil {
		return nil, errors.New("product does not exist")
	}

	// verify ownership
	if existingProd.UserID != user.ID {
		return nil, errors.New("product does not belong to the user")
	}

	if len(input.Name) > 0 {
		existingProd.Name = input.Name
	}
	if len(input.Description) > 0 {
		existingProd.Description = input.Description
	}
	if input.Price > 0 {
		existingProd.Price = input.Price
	}
	if input.CategoryID > 0 {
		existingProd.CategoryID = input.CategoryID
	}
	if len(input.ImageURL) > 0 {
		existingProd.ImageURL = input.ImageURL
	}
	if input.Stock > 0 {
		existingProd.Stock = uint(input.Stock)
	}

	updatedProd, err := s.Repo.EditProduct(existingProd)

	return updatedProd, err
}

func (s CatalogService) DeleteProduct(id int) error {
	deletedProduct, err := s.Repo.FindProductByID(uint(id))

	if err != nil {
		return errors.New("product does not exist for deletion")
	}

	err = s.Repo.DeleteProduct(deletedProduct)

	if err != nil {
		return errors.New("could not delete product")
	}

	return nil
}

func (s CatalogService) GetProductByID(id uint) (*domain.Product, error) {
	product, err := s.Repo.FindProductByID(id)

	if err != nil {
		return nil, err
	}
	return product, nil

}

func (s CatalogService) GetProducts() ([]*domain.Product, error) {
	products, err := s.Repo.FindProducts()

	if err != nil {
		return nil, err
	}

	return products, nil
}

func (s CatalogService) GetSellerProducts(id uint) ([]*domain.Product, error) {
	products, err := s.Repo.FindSellerProducts(id)

	if err != nil {
		return nil, err
	}

	return products, nil
}

func (s CatalogService) UpdateProductStock(e domain.Product) (*domain.Product, error) {
	product, err := s.Repo.FindProductByID(e.ID)

	if err != nil {
		return nil, errors.New("product does not exist")
	}

	if product.UserID != e.UserID {
		return nil, errors.New("product does not belong to the user")
	}

	product.Stock = e.Stock
	editProduct, err := s.Repo.EditProduct(product)
	if err != nil {
		return nil, errors.New("could not update stock")
	}
	return editProduct, nil
}
