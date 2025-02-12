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
