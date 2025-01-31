package service

import (
	"github.com/sharat789/zamazon-be/internal/domain"
	"github.com/sharat789/zamazon-be/internal/dto"
	"log"
)

type UserService struct {
}

func (s UserService) UserSignup(input dto.UserSignup) (string, error) {
	log.Println(input)
	return "someTokenCreated", nil
}
func (s UserService) findUserByEmail(email string) (*domain.User, error) {
	return nil, nil
}

func (s UserService) Login(input any) (string, error) {
	return "", nil
}

func (s UserService) GetVerificationCode(e domain.User) (int, error) {
	return 0, nil
}

func (s UserService) VerifyCode(verificationCode int) error {
	return nil
}

func (s UserService) CreateUserProfile(id uint, input any) error {
	return nil
}

func (s UserService) GetUserProfile(id uint) (*domain.User, error) {
	return nil, nil
}

func (s UserService) UpdateProfile(id uint, input any) error {
	return nil
}

func (s UserService) BecomeSeller(id uint, input any) (string, error) {
	return "", nil
}

func (s UserService) FindCart(id uint) ([]interface{}, error) {
	return nil, nil
}

func (s UserService) CreateCart(input any, u domain.User) ([]interface{}, error) {
	return nil, nil
}

func (s UserService) CreateOrder(u domain.User) (int, error) {
	return 0, nil
}

func (s UserService) GetOrders(u domain.User) ([]interface{}, error) {
	return nil, nil
}

func (s UserService) GetOrderByID(id uint, uid uint) (interface{}, error) {
	return nil, nil
}
