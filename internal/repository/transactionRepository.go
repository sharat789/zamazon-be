package repository

import (
	"github.com/sharat789/zamazon-be/internal/domain"
	"github.com/sharat789/zamazon-be/internal/dto"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreatePayment(payment *domain.Payment) error
	FindOrders(userId uint) ([]domain.OrderItem, error)
	FindOrderById(orderId uint, userId uint) (dto.SellerOrderDetails, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func (t transactionRepository) CreatePayment(payment *domain.Payment) error {
	//TODO implement me
	panic("implement me")
}

func (t transactionRepository) FindOrders(userId uint) ([]domain.OrderItem, error) {
	//TODO implement me
	panic("implement me")
}

func (t transactionRepository) FindOrderById(orderId uint, userId uint) (dto.SellerOrderDetails, error) {
	//TODO implement me
	panic("implement me")
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{
		db,
	}
}
