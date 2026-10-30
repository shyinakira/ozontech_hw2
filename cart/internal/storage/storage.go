package storage

import (
	"fmt"
	"hw2/cart/internal/domain"
	"sync"
)

type User struct {
	ID domain.UserID
}

type Storage struct {
	mu    sync.RWMutex
	carts map[domain.UserID]*domain.Cart
}

func NewStorage() *Storage {
	return &Storage{
		carts: make(map[domain.UserID]*domain.Cart),
	}
}

func (s *Storage) AddCartItem(userID domain.UserID, item domain.CartItem) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cart, ok := s.carts[userID]
	if !ok {
		s.carts[userID] = &domain.Cart{
			Items:      []domain.CartItem{item},
			TotalPrice: item.Price,
		}
		return
	}

	cart.AddItem(item)
}

func (s *Storage) GetCart(userID domain.UserID) (domain.Cart, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cart, ok := s.carts[userID]
	if !ok || len(cart.Items) == 0 {
		return domain.Cart{}, fmt.Errorf("failed to get cart: %w", domain.ErrCartNotFound)
	}

	itemsCopy := make([]domain.CartItem, len(cart.Items))
	copy(itemsCopy, cart.Items)

	return domain.Cart{
		Items:      itemsCopy,
		TotalPrice: cart.TotalPrice,
	}, nil
}
