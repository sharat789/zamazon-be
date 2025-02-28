package service

import (
	"github.com/sharat789/zamazon-be/internal/domain"
	"github.com/sharat789/zamazon-be/internal/dto"
	"github.com/sharat789/zamazon-be/internal/helper"
	"github.com/sharat789/zamazon-be/internal/repository"
	"github.com/stripe/stripe-go/v78"
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

func (s TransactionService) GetActivePayment(userId uint) (*domain.Payment, error) {
	return s.Repo.FindExistingPayment(userId)
}

func (s TransactionService) StoreCreatedPayment(userId uint, ps *stripe.CheckoutSession, amount float64, orderId string) error {
	payment := domain.Payment{
		UserId:     userId,
		Amount:     amount,
		OrderId:    orderId,
		Status:     string(domain.PaymentStatusInitial),
		PaymentUrl: ps.URL,
		PaymentId:  ps.ID,
	}

	return s.Repo.CreatePayment(&payment)
}
func NewTransactionService(repo repository.TransactionRepository, auth helper.Auth) TransactionService {
	return TransactionService{
		Repo: repo,
		Auth: auth,
	}
}
