package storage

import (
	"fmt"
	"hw2/cart/internal/domain"
)
type UserID int64

type User struct {
	ID UserID
}

type Storage struct {
	carts map[UserID]domain.Cart
}

func (s *Storage) addCartItem(userID UserID, item domain.CartItem) error {
	cart, ok := s.carts[userID]
	if !ok {
		return fmt.Errorf("failed to add cart item: %w", domain.ErrCartNotFound)
	}

	cart.AddItem(item)
	return nil
}

func (s *Storage) getCart(userID UserID) (domain.Cart, error) {
	cart, ok := s.carts[userID]
	if !ok {
		return domain.Cart{}, fmt.Errorf("failed to get cart: %w", domain.ErrCartNotFound)
	}

	return cart, nil
}

var storage Storage

func NewStorage() {
	storage = Storage{
		carts: make(map[UserID]domain.Cart),
	}
}

func AddCartItem(userID UserID, item domain.CartItem) error {
	err := storage.addCartItem(userID, item)
	if err != nil {
		return err
	}

	return nil
}

func GetCart(userID UserID) (domain.Cart, error) {
	cart, err := storage.getCart(userID)
	if err != nil {
		return domain.Cart{}, err
	}

	return cart, nil
}