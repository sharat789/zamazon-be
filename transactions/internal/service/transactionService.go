package service

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sharat789/zamazon-be-ms/transactions/internal/client"
	"github.com/sharat789/zamazon-be-ms/transactions/internal/domain"
	"github.com/sharat789/zamazon-be-ms/transactions/internal/dto"
	"github.com/sharat789/zamazon-be-ms/transactions/internal/repository"
)

type TransactionService struct {
	Repo       repository.TransactionRepository
	AuthClient *client.AuthClient
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

func (s TransactionService) GetCurrentUser(c *fiber.Ctx) *client.TokenUser {
	user, ok := c.Locals("user").(*client.TokenUser)
	if !ok {
		return nil
	}
	return user
}
func NewTransactionService(repo repository.TransactionRepository, authClient *client.AuthClient) TransactionService {
	return TransactionService{
		Repo:       repo,
		AuthClient: authClient,
	}
}
