package service

import (
	"errors"
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

func (s TransactionService) GetActivePayment(userId uint) (*domain.Payment, error) {
	return s.Repo.FindExistingPayment(userId)
}

func (s TransactionService) StoreCreatedPayment(input dto.CreatePaymentRequest) error {
	payment := domain.Payment{
		UserId:       input.UserId,
		Amount:       input.Amount,
		OrderId:      input.OrderId,
		Status:       string(domain.PaymentStatusInitial),
		PaymentId:    input.PaymentId,
		ClientSecret: input.ClientSecret,
	}

	return s.Repo.CreatePayment(&payment)
}

func (s TransactionService) UpdatePayment(userId uint, status string, paymentLog string) error {
	p, err := s.GetActivePayment(userId)
	if err != nil {
		return err
	}
	p.Status = string(domain.PaymentStatus(status))
	p.Response = paymentLog
	return s.Repo.UpdatePayment(p)
}

func (s TransactionService) GetPaymentByID(paymentId string) (domain.Payment, error) {
	payment, err := s.Repo.FindPaymentByID(paymentId)
	if err != nil {
		return domain.Payment{}, errors.New("payment not found")
	}
	return payment, nil
}
func NewTransactionService(repo repository.TransactionRepository, auth helper.Auth) TransactionService {
	return TransactionService{
		Repo: repo,
		Auth: auth,
	}
}
