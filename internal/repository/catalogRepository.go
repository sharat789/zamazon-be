package repository

import (
	"errors"
	"github.com/sharat789/zamazon-be/internal/domain"
	"gorm.io/gorm"
	"log"
)

type CatalogRepository interface {
	CreateCategory(e *domain.Category) error
	FindCategories() ([]*domain.Category, error)
	FindCategoryByID(id uint) (*domain.Category, error)
	EditCategory(e *domain.Category) (*domain.Category, error)
	DeleteCategory(id uint) error
}

type catalogRepository struct {
	db *gorm.DB
}

func (c catalogRepository) CreateCategory(e *domain.Category) error {
	err := c.db.Create(&e).Error

	if err != nil {
		log.Printf("Error while creating category %v", err)
		return errors.New("could not create category")
	}
	return nil
}

func (c catalogRepository) FindCategories() ([]*domain.Category, error) {
	var categories []*domain.Category
	err := c.db.Find(&categories).Error

	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (c catalogRepository) FindCategoryByID(id uint) (category *domain.Category, err error) {
	var foundCategory *domain.Category

	err = c.db.First(&category, id).Error

	if err != nil {
		log.Printf("Could not find category with the id %d: %v", id, err)
		return &domain.Category{}, errors.New("could not find category")
	}
	return foundCategory, nil
}

func (c catalogRepository) EditCategory(e *domain.Category) (*domain.Category, error) {
	err := c.db.Save(&e).Error

	if err != nil {
		log.Printf("db error: %v", err)
		return nil, errors.New("fail to update category")
	}
	return e, nil
}

func (c catalogRepository) DeleteCategory(id uint) error {
	err := c.db.Delete(&domain.Category{}, id).Error

	if err != nil {
		log.Printf("db error: %v", err)
		return errors.New("fail to delete category")
	}
	return nil
}

func NewCatalogRepository(db *gorm.DB) CatalogRepository {
	return &catalogRepository{
		db: db,
	}
}
