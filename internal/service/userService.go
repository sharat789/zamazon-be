package service

import (
	"errors"
	"github.com/sharat789/zamazon-be/internal/domain"
	"github.com/sharat789/zamazon-be/internal/dto"
	"github.com/sharat789/zamazon-be/internal/helper"
	"github.com/sharat789/zamazon-be/internal/repository"
	"log"
)

type UserService struct {
	Repo repository.UserRepository
	Auth helper.Auth
}

func (s UserService) UserSignup(input dto.UserSignup) (string, error) {
	hashPassword, err := s.Auth.CreateHashPassword(input.Password)

	if err != nil {
		return "", err
	}

	user, err := s.Repo.CreateUser(domain.User{
		Email:    input.Email,
		Password: hashPassword,
		Phone:    input.Phone,
	})

	if err != nil {
		return "", err
	}

	return s.Auth.GenerateToken(user.ID, user.Email, user.UserType)
}
func (s UserService) findUserByEmail(email string) (*domain.User, error) {
	user, err := s.Repo.FindUser(email)
	log.Println(user)
	return &user, err
}

func (s UserService) Login(email string, password string) (string, error) {
	log.Println(email, password)

	user, err := s.findUserByEmail(email)

	if err != nil {
		return "", errors.New("user does not exist with the provided email")
	}
	log.Println(user.Password)
	err = s.Auth.VerifyPassword(password, user.Password)

	if err != nil {
		return "", err
	}
	return s.Auth.GenerateToken(user.ID, user.Email, user.UserType)
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
