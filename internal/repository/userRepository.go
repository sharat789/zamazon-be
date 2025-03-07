package repository

import (
	"errors"
	"github.com/sharat789/zamazon-be/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
)

type UserRepository interface {
	CreateUser(u domain.User) (domain.User, error)
	FindUser(email string) (domain.User, error)
	FindUserByID(id uint) (domain.User, error)
	UpdateUser(id uint, u domain.User) (domain.User, error)
	CreateBankAccount(e domain.BankDetails) error

	//cart operations
	FindCartItems(userId uint) ([]domain.Cart, error)
	FindCartItem(userId, productId uint) (domain.Cart, error)
	UpdateCartItem(cart domain.Cart) error
	CreateCart(cart domain.Cart) error
	DeleteCartById(id uint) error
	DeleteCartItems(userId uint) error
	DeleteCartItem(userID, productID uint) error

	//order operations
	FindOrders(userId uint) ([]domain.Order, error)
	CreateOrder(order domain.Order) error
	FindOrderByID(orderId uint, userId uint) (domain.Order, error)

	// profile operations
	CreateProfile(e domain.Address) error
	UpdateProfile(e domain.Address) error
}

type userRepository struct {
	db *gorm.DB
}

func (r userRepository) CreateUser(user domain.User) (domain.User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		log.Printf("Error while creating user %v", err)
		return domain.User{}, errors.New("could not create user")
	}
	return user, nil
}

func (r userRepository) FindUser(email string) (domain.User, error) {
	var user domain.User

	err := r.db.Preload("Address").First(&user, "email=?", email).Error
	if err != nil {
		log.Printf("Could not find user with the email %s: %v", email, err)
		return domain.User{}, errors.New("could not find user")
	}

	return user, nil
}

func (r userRepository) FindUserByID(id uint) (domain.User, error) {
	var user domain.User

	err := r.db.Preload("Address").First(&user, id).Error
	if err != nil {
		log.Printf("Could not find user with the id %d: %v", id, err)
		return domain.User{}, errors.New("could not find user")
	}

	return user, nil
}

func (r userRepository) UpdateUser(id uint, u domain.User) (domain.User, error) {
	var user domain.User

	err := r.db.Model(&user).Clauses(clause.Returning{}).Where("id=?", id).Updates(u).Error

	if err != nil {
		log.Printf("error on update %v", err)
		return domain.User{}, errors.New("failed update user")
	}

	return user, nil
}

func (r userRepository) CreateBankAccount(e domain.BankDetails) error {
	return r.db.Create(&e).Error
}

func (r userRepository) FindOrderByID(orderId uint, userId uint) (domain.Order, error) {
	order := domain.Order{}
	err := r.db.Preload("Items").Where("id=? AND user_id=?", orderId, userId).First(&order).Error
	if err != nil {
		log.Printf("Error while fetching order %v", err)
		return domain.Order{}, errors.New("could not fetch order")
	}
	return order, nil
}

func (r userRepository) FindOrders(userId uint) ([]domain.Order, error) {
	var orders []domain.Order
	err := r.db.Where("user_id=?", userId).Find(&orders).Error

	if err != nil {
		log.Printf("Error while fetching orders %v", err)
		return nil, errors.New("could not fetch orders")
	}
	return orders, nil
}

func (r userRepository) CreateOrder(order domain.Order) error {
	err := r.db.Create(&order).Error
	if err != nil {
		log.Printf("Error while creating order %v", err)
		return errors.New("could not create order")
	}
	return nil
}

func (r userRepository) CreateProfile(e domain.Address) error {
	err := r.db.Create(&e).Error
	if err != nil {
		log.Printf("Error while creating profile with address%v", err)
		return errors.New("could not create profile")
	}
	return nil
}

func (r userRepository) UpdateProfile(e domain.Address) error {
	err := r.db.Where("user_id=?", e.UserID).Updates(e).Error
	if err != nil {
		log.Printf("Error while updating profile with address%v", err)
		return errors.New("could not update profile")
	}
	return nil
}

func (r userRepository) FindCartItems(userId uint) ([]domain.Cart, error) {
	var carts []domain.Cart
	err := r.db.Where("user_id=?", userId).Find(&carts).Error
	return carts, err
}

func (r userRepository) FindCartItem(userId, productId uint) (domain.Cart, error) {
	cartItem := domain.Cart{}
	err := r.db.First(&cartItem, "user_id=? AND product_id=?", userId, productId).Error
	return cartItem, err
}

func (r userRepository) UpdateCartItem(cart domain.Cart) error {
	var c domain.Cart
	err := r.db.Model(&c).Clauses(clause.Returning{}).Where("id=?", cart.ID).Updates(cart).Error
	return err
}

func (r userRepository) DeleteCartItem(userID, productID uint) error {
	err := r.db.Where("user_id = ? AND product_id = ?", userID, productID).Delete(&domain.Cart{}).Error
	return err
}

func (r userRepository) CreateCart(cart domain.Cart) error {
	return r.db.Create(&cart).Error
}

func (r userRepository) DeleteCartById(id uint) error {
	err := r.db.Delete(&domain.Cart{}, id).Error
	return err
}

func (r userRepository) DeleteCartItems(userId uint) error {
	err := r.db.Where("user_id=?", userId).Delete(&domain.Cart{}).Error
	return err
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db,
	}
}
