package service

import (
	"github.com/sharat789/zamazon-be/internal/domain"
	"github.com/sharat789/zamazon-be/internal/dto"
	"github.com/sharat789/zamazon-be/internal/helper"
	"github.com/sharat789/zamazon-be/internal/repository"
)

type TransactionService struct {
	Repo repository.TransactionRepository
	Auth helper.Auth
}

func (s TransactionService) GetOrders(u domain.User) ([]domain.OrderItem, error) {
	orders, err := s.Repo.FindOrders(u.ID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (s TransactionService) GetOrderDetails(id uint, u domain.User) (dto.SellerOrderDetails, error) {
	order, err := s.Repo.FindOrderById(id, u.ID)
	if err != nil {
		return dto.SellerOrderDetails{}, err
	}
	return order, nil
}
func NewTransactionService(repo repository.TransactionRepository, auth helper.Auth) TransactionService {
	return TransactionService{
		Repo: repo,
		Auth: auth,
	}
}
