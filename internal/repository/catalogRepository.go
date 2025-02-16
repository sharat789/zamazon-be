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

	CreateProduct(e *domain.Product) error
	FindProducts() ([]*domain.Product, error)
	FindProductByID(id uint) (*domain.Product, error)
	FindSellerProducts(id uint) ([]*domain.Product, error)
	EditProduct(e *domain.Product) (*domain.Product, error)
	DeleteProduct(e *domain.Product) error
}

type catalogRepository struct {
	db *gorm.DB
}

func (c catalogRepository) CreateProduct(e *domain.Product) error {
	err := c.db.Model(&domain.Product{}).Create(&e).Error

	if err != nil {
		log.Printf("Error while creating product %v", err)
		return errors.New("could not create product")
	}
	return nil
}

func (c catalogRepository) FindProducts() ([]*domain.Product, error) {
	var products []*domain.Product
	err := c.db.Find(&products).Error

	if err != nil {
		return nil, err
	}

	return products, nil
}

func (c catalogRepository) FindProductByID(id uint) (*domain.Product, error) {
	var product *domain.Product

	err := c.db.First(&product, id).Error

	if err != nil {
		log.Printf("Could not find product with the id %d: %v", id, err)
		return &domain.Product{}, errors.New("could not find product")
	}
	return product, nil
}

func (c catalogRepository) FindSellerProducts(id uint) ([]*domain.Product, error) {
	var products []*domain.Product

	err := c.db.Where("user_id=?", id).Find(&products).Error

	if err != nil {
		log.Printf("Could not find products for seller with the id %d: %v", id, err)
		return nil, errors.New("could not find products")
	}
	return products, nil
}

func (c catalogRepository) EditProduct(e *domain.Product) (*domain.Product, error) {
	err := c.db.Save(&e).Error

	if err != nil {
		log.Printf("db error: %v", err)
		return nil, errors.New("fail to update product")
	}
	return e, nil
}

func (c catalogRepository) DeleteProduct(e *domain.Product) error {
	err := c.db.Delete(&e).Error

	if err != nil {
		log.Printf("db error: %v", err)
		return errors.New("fail to delete product")
	}
	return nil
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
