package storage

import (
	"fmt"
	"hw2/cart/internal/domain"
)

type User struct {
	ID domain.UserID
}

type Storage struct {
	carts map[domain.UserID]*domain.Cart
}

func NewStorage() *Storage {
	return &Storage{
		carts: make(map[domain.UserID]*domain.Cart),
	}
}

func (s *Storage) AddCartItem(userID domain.UserID, item domain.CartItem) {
	cart, ok := s.carts[userID]
	if !ok {
		s.carts[userID] = &domain.Cart{
			Items: []domain.CartItem{item},
			TotalPrice: item.Price,
		}
		return
	}
	
	cart.AddItem(item)
}

func (s *Storage) GetCart(userID domain.UserID) (domain.Cart, error) {
	cart, ok := s.carts[userID]
	if !ok {
		return domain.Cart{}, fmt.Errorf("failed to get cart: %w", domain.ErrCartNotFound)
	}

	return *cart, nil
}
