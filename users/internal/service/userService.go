package service

import (
	"errors"
	"github.com/sharat789/zamazon-be-ms/common/auth"
	"github.com/sharat789/zamazon-be-ms/users/internal/client"
	"github.com/sharat789/zamazon-be-ms/users/internal/domain"
	"github.com/sharat789/zamazon-be-ms/users/internal/dto"
	"github.com/sharat789/zamazon-be-ms/users/internal/repository"
	"log"
	"time"
)

type UserService struct {
	Repo          repository.UserRepository
	Auth          auth.Auth
	CatalogClient *client.CatalogClient
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

func (s UserService) GetVerificationCode(userID uint) (string, error) {
	if s.isVerifiedUser(userID) {
		return "", errors.New("user already verified")
	}

	code, err := s.Auth.GenerateCode()
	if err != nil {
		return "", err
	}

	user := domain.User{
		Expiry:           time.Now().Add(30 * time.Minute),
		VerificationCode: code,
	}
	_, err = s.Repo.UpdateUser(userID, user)

	if err != nil {
		return "", errors.New("unable to update verification code")
	}
	return code, nil
}

func (s UserService) VerifyCode(id uint, verificationCode string) error {
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

func (s UserService) CreateUserProfile(id uint, input dto.ProfileInput) error {
	user, err := s.Repo.FindUserByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	if input.FirstName != "" {
		user.FName = input.FirstName
	}
	if input.LastName != "" {
		user.LName = input.LastName
	}
	_, err = s.Repo.UpdateUser(id, user)

	if err != nil {
		return err
	}

	address := domain.Address{
		AddressLine1: input.AddressInput.AddressLine1,
		AddressLine2: input.AddressInput.AddressLine2,
		City:         input.AddressInput.City,
		PostCode:     input.AddressInput.PostCode,
		Country:      input.AddressInput.Country,
		UserID:       id,
	}

	err = s.Repo.CreateProfile(address)
	if err != nil {
		return err
	}
	return nil
}

func (s UserService) GetUserProfile(id uint) (*domain.User, error) {
	user, err := s.Repo.FindUserByID(id)

	if err != nil {
		log.Printf("Error while fetching user profile %v", err)
	}

	return &user, nil
}

func (s UserService) UpdateProfile(id uint, input dto.ProfileInput) error {
	user, err := s.Repo.FindUserByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	if input.FirstName != "" {
		user.FName = input.FirstName
	}
	if input.LastName != "" {
		user.LName = input.LastName
	}
	_, err = s.Repo.UpdateUser(id, user)

	if err != nil {
		return err
	}

	address := domain.Address{
		AddressLine1: input.AddressInput.AddressLine1,
		AddressLine2: input.AddressInput.AddressLine2,
		City:         input.AddressInput.City,
		PostCode:     input.AddressInput.PostCode,
		Country:      input.AddressInput.Country,
		UserID:       id,
	}

	err = s.Repo.UpdateProfile(address)

	if err != nil {
		return err
	}
	return nil
}

func (s UserService) FindCart(id uint) ([]domain.Cart, float64, error) {
	cartItems, err := s.Repo.FindCartItems(id)

	if err != nil {
		log.Printf("Error while fetching cart %v", err)
		return nil, 0, errors.New("unable to fetch cart items")
	}
	var totalAmount float64
	for _, item := range cartItems {
		totalAmount += float64(item.Qty) * item.Price
	}
	return cartItems, totalAmount, nil
}

func (s UserService) CreateCart(input dto.CreateCartRequest, u auth.TokenUser) ([]domain.Cart, error) {
	// check if cart exists
	cart, _ := s.Repo.FindCartItem(u.ID, input.ProductID)
	if cart.ID != 0 {
		if input.ProductID == 0 {
			return nil, errors.New("product id is required")
		}
		if input.Qty < 1 {
			err := s.Repo.DeleteCartById(cart.ID)
			if err != nil {
				log.Printf("Error while deleting cart %v", err)
				return nil, errors.New("unable to delete cart")
			}
		} else {
			cart.Qty = input.Qty
			err := s.Repo.UpdateCartItem(cart)
			if err != nil {
				log.Printf("Error while updating cart %v", err)
				return nil, errors.New("unable to update cart")
			}
		}
	} else {
		// Use the catalog client to get the product
		product, err := s.CatalogClient.GetProductByID(input.ProductID)
		log.Println(product)
		if err != nil {
			return nil, errors.New("product not found")
		}

		err = s.Repo.CreateCart(domain.Cart{
			UserID:    u.ID,
			ProductID: input.ProductID,
			Name:      product.Name,
			ImageURL:  product.ImageURL,
			SellerId:  product.UserID,
			Price:     product.Price,
			Qty:       input.Qty,
		})
		if err != nil {
			log.Printf("Error while creating cart %v", err)
			return nil, errors.New("unable to create cart")
		}
	}
	return s.Repo.FindCartItems(u.ID)
}

func (s UserService) UpdateProductQtyInCart(userID uint, productID uint, qty int) error {
	cart, err := s.Repo.FindCartItem(userID, productID)
	if err != nil {
		return errors.New("cart item not found")
	}
	cart.Qty = uint(qty)
	err = s.Repo.UpdateCartItem(cart)
	if err != nil {
		return errors.New("unable to update cart item")
	}
	return nil
}

func (s UserService) RemoveProductFromCart(userID uint, productID uint) error {
	err := s.Repo.DeleteCartItem(userID, productID)
	if err != nil {
		return errors.New("unable to remove product from cart")
	}
	return nil
}
func (s UserService) CreateOrder(request dto.CreateOrderRequest) error {
	cartItems, _, err := s.FindCart(request.UserID)

	if err != nil {
		return errors.New("unable to fetch cart items")
	}

	if len(cartItems) == 0 {
		return errors.New("cart is empty")
	}

	var orderItems []domain.OrderItem

	for _, item := range cartItems {
		orderItems = append(orderItems, domain.OrderItem{
			ProductID: item.ProductID,
			Name:      item.Name,
			ImageURL:  item.ImageURL,
			SellerId:  item.SellerId,
			Price:     uint(item.Price),
			Qty:       item.Qty,
		})
	}

	order := domain.Order{
		UserID:         request.UserID,
		PaymentId:      request.PaymentId,
		OrderRefNumber: request.OrderRefNumber,
		Amount:         request.Amount,
		Items:          orderItems,
		Status:         "completed",
	}

	err = s.Repo.CreateOrder(order)

	if err != nil {
		return err
	}

	err = s.Repo.DeleteCartItems(request.UserID)
	log.Printf("Error while deleting cart items %v", err)

	return err
}

func (s UserService) ClearCart(userID uint) error {
	err := s.Repo.DeleteCartItems(userID)
	if err != nil {
		return errors.New("unable to clear cart")
	}
	return nil
}

func (s UserService) GetOrders(userID uint) ([]domain.Order, error) {
	orders, err := s.Repo.FindOrders(userID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (s UserService) GetOrderByID(id uint, userId uint) (domain.Order, error) {
	order, err := s.Repo.FindOrderByID(id, userId)

	if err != nil {
		return order, err
	}
	return order, nil
}
