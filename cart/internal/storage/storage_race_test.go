package storage

import (
	"sync"
	"testing"

	"hw2/cart/internal/domain"

	"github.com/stretchr/testify/require"
)

func TestStorage_ConcurrentAccess(t *testing.T) {
	s := NewStorage()
	const (
		userID     = domain.UserID(42)
		goroutines = 50
		itemsPerG  = 20
	)

	var wg sync.WaitGroup
	wg.Add(goroutines * 2)

	for g := 0; g < goroutines; g++ {
		go func(gID int) {
			defer wg.Done()
			for i := 0; i < itemsPerG; i++ {
				sku := int64((gID % 5) + 1)
				s.AddCartItem(userID, domain.CartItem{
					SKU:   sku,
					Name:  "concurrent_item",
					Price: 100,
					Count: 1,
				})
			}
		}(g)
	}

	for g := 0; g < goroutines; g++ {
		go func() {
			defer wg.Done()
			for i := 0; i < itemsPerG; i++ {
				_, _ = s.GetCart(userID)
			}
		}()
	}

	wg.Wait()

	cart, err := s.GetCart(userID)
	require.NoError(t, err)

	totalCount := uint32(0)
	for _, item := range cart.Items {
		totalCount += item.Count
	}
	require.Equal(t, uint32(goroutines*itemsPerG), totalCount)
}
