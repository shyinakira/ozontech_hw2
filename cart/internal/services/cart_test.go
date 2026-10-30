package services

import (
	"context"
	"testing"

	"hw2/cart/internal/domain"
	pb "hw2/products/gen/go/proto"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MockStorage struct {
	carts map[domain.UserID]*domain.Cart
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		carts: make(map[domain.UserID]*domain.Cart),
	}
}

func (s *MockStorage) AddCartItem(userID domain.UserID, item domain.CartItem) {
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

func (s *MockStorage) GetCart(userID domain.UserID) (domain.Cart, error) {
	cart, ok := s.carts[userID]
	if !ok || len(cart.Items) == 0 {
		return domain.Cart{}, domain.ErrCartNotFound
	}

	itemsCopy := make([]domain.CartItem, len(cart.Items))
	copy(itemsCopy, cart.Items)

	return domain.Cart{
		Items:      itemsCopy,
		TotalPrice: cart.TotalPrice,
	}, nil
}

type MockProductServiceClient struct {
	products map[int64]*pb.GetProductResponse
}

func NewMockProductServiceClient() *MockProductServiceClient {
	return &MockProductServiceClient{
		products: map[int64]*pb.GetProductResponse{
			1: {
				Sku:   1,
				Name:  "test_product_1",
				Price: 100,
			},
			2: {
				Sku:   2,
				Name:  "test_product_2",
				Price: 200,
			},
		},
	}
}

func (c *MockProductServiceClient) GetProduct(ctx context.Context, in *pb.GetProductRequest, opts ...grpc.CallOption) (*pb.GetProductResponse, error) {
	sku := in.GetSku()
	prod, ok := c.products[sku]
	if !ok {
		return nil, status.Error(codes.NotFound, "product not found")
	}

	return prod, nil
}

func TestAddCartItem(t *testing.T) {
	t.Parallel()

	type input struct {
		userIDStr string
		skuIDStr  string
		count     uint32
	}

	tests := []struct {
		name          string
		in            input
		expectedError error
	}{
		{
			name: "happy path",
			in: input{
				userIDStr: "1",
				skuIDStr:  "1",
				count:     2,
			},
			expectedError: nil,
		},
		{
			name: "invalid user id negative",
			in: input{
				userIDStr: "-100",
				skuIDStr:  "1",
				count:     2,
			},
			expectedError: domain.ErrInvalidUserID,
		},
		{
			name: "invalid user id zero",
			in: input{
				userIDStr: "0",
				skuIDStr:  "1",
				count:     2,
			},
			expectedError: domain.ErrInvalidUserID,
		},
		{
			name: "invalid sku id negative",
			in: input{
				userIDStr: "1",
				skuIDStr:  "-5",
				count:     2,
			},
			expectedError: domain.ErrInvalidSkuID,
		},
		{
			name: "invalid sku id zero",
			in: input{
				userIDStr: "1",
				skuIDStr:  "0",
				count:     2,
			},
			expectedError: domain.ErrInvalidSkuID,
		},
		{
			name: "invalid count zero",
			in: input{
				userIDStr: "1",
				skuIDStr:  "1",
				count:     0,
			},
			expectedError: domain.ErrInvalidSkuCount,
		},
		{
			name: "product not found in product service",
			in: input{
				userIDStr: "1",
				skuIDStr:  "9999",
				count:     1,
			},
			expectedError: domain.ErrProductNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			mockStorage := NewMockStorage()
			mockClient := NewMockProductServiceClient()
			service := NewCartService(mockStorage, mockClient)

			err := service.AddCartItem(t.Context(), test.in.userIDStr, test.in.skuIDStr, test.in.count)
			require.ErrorIs(t, err, test.expectedError)
		})
	}
}

func TestGetCart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                 string
		setupCart            func(s *MockStorage)
		userIDStr            string
		expectedCartResponse domain.Cart
		expectedError        error
	}{
		{
			name: "happy path with parallel fetch",
			setupCart: func(s *MockStorage) {
				s.AddCartItem(1, domain.CartItem{SKU: 1, Count: 2})
				s.AddCartItem(1, domain.CartItem{SKU: 2, Count: 1})
			},
			userIDStr: "1",
			expectedCartResponse: domain.Cart{
				Items: []domain.CartItem{
					{SKU: 1, Name: "test_product_1", Price: 100, Count: 2},
					{SKU: 2, Name: "test_product_2", Price: 200, Count: 1},
				},
				TotalPrice: 400,
			},
			expectedError: nil,
		},
		{
			name: "invalid user id negative",
			setupCart: func(s *MockStorage) {},
			userIDStr: "-1",
			expectedError: domain.ErrInvalidUserID,
		},
		{
			name: "cart not found",
			setupCart: func(s *MockStorage) {},
			userIDStr: "999",
			expectedError: domain.ErrCartNotFound,
		},
		{
			name: "product removed from product service returns ErrProductNotFound",
			setupCart: func(s *MockStorage) {
				s.AddCartItem(1, domain.CartItem{SKU: 999, Count: 1})
			},
			userIDStr:     "1",
			expectedError: domain.ErrProductNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			mockStorage := NewMockStorage()
			test.setupCart(mockStorage)
			mockClient := NewMockProductServiceClient()
			service := NewCartService(mockStorage, mockClient)

			cart, err := service.GetCart(test.userIDStr)
			require.ErrorIs(t, err, test.expectedError)
			if test.expectedError == nil {
				require.Equal(t, test.expectedCartResponse, cart)
			}
		})
	}
}
