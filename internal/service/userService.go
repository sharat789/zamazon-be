package service

import (
	"errors"
	"github.com/sharat789/zamazon-be/internal/domain"
	"github.com/sharat789/zamazon-be/internal/dto"
	"github.com/sharat789/zamazon-be/internal/helper"
	"github.com/sharat789/zamazon-be/internal/repository"
	"log"
	"time"
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

func (s UserService) isVerifiedUser(id uint) bool {
	currentUser, err := s.Repo.FindUserByID(id)

	return err == nil && currentUser.IsVerified
}
func (s UserService) GetVerificationCode(e domain.User) (int, error) {
	if s.isVerifiedUser(e.ID) {
		return 0, errors.New("user already verified")
	}

	code, err := s.Auth.GenerateCode()
	if err != nil {
		return 0, err
	}

	user := domain.User{
		Expiry:           time.Now().Add(30 * time.Minute),
		VerificationCode: code,
	}
	_, err = s.Repo.UpdateUser(e.ID, user)

	if err != nil {
		return 0, errors.New("unable to update verification code")
	}
	return code, nil
}

func (s UserService) VerifyCode(id uint, verificationCode int) error {
	if s.isVerifiedUser(id) {
		log.Println("is verified")
		return errors.New("user already verified")
	}
	user, err := s.Repo.FindUserByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	if user.VerificationCode != verificationCode {
		return errors.New("verification code doesn't match")
	}

	if time.Now().After(user.Expiry) {
		return errors.New("verification code expired")
	}

	updatedUser := domain.User{
		IsVerified: true,
	}

	_, err = s.Repo.UpdateUser(id, updatedUser)

	if err != nil {
		return errors.New("unable to verify user")
	}
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

func (s UserService) BecomeSeller(id uint, input dto.SellerInput) (string, error) {
	user, _ := s.Repo.FindUserByID(id)

	if user.UserType == domain.SELLER {
		return "", errors.New("you are already in the seller program")
	}

	seller, err := s.Repo.UpdateUser(id, domain.User{
		FName:    input.FirstName,
		LName:    input.LastName,
		Phone:    input.PhoneNumber,
		UserType: domain.SELLER,
	})

	if err != nil {
		return "", err
	}

	token, err := s.Auth.GenerateToken(user.ID, user.Email, seller.UserType)

	err = s.Repo.CreateBankAccount(domain.BankDetails{
		AccountNo:   input.BankAccountNumber,
		SwiftCode:   input.SwiftCode,
		PaymentType: input.PaymentType,
		UserId:      id,
	})
	return token, err
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
