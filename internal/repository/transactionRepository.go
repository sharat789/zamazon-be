package repository

import (
	"github.com/sharat789/zamazon-be/internal/domain"
	"github.com/sharat789/zamazon-be/internal/dto"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreatePayment(payment *domain.Payment) error
	FindExistingPayment(userId uint) (*domain.Payment, error)
	FindOrders(userId uint) ([]domain.OrderItem, error)
	FindOrderById(orderId uint, userId uint) (dto.SellerOrderDetails, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func (t transactionRepository) FindExistingPayment(userId uint) (*domain.Payment, error) {
	var payment *domain.Payment
	err := t.db.Where("user_id = ? AND status = ?", userId, "initial").Order("created_at desc").First(&payment).Error
	return payment, err
}

func (t transactionRepository) CreatePayment(payment *domain.Payment) error {
	return t.db.Create(payment).Error
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
