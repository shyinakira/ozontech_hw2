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

var (
	SUCCESS_CART = domain.Cart{
		Items: []domain.CartItem{
			{
				SKU:   1,
				Name:  "test_product",
				Count: 2,
				Price: 100,
			},
		},
		TotalPrice: 200,
	}
)

type MockStorage struct{}

func (s *MockStorage) AddCartItem(userID domain.UserID, item domain.CartItem) {}

func (s *MockStorage) GetCart(userID domain.UserID) (domain.Cart, error) {
	if userID != 1 {
		return domain.Cart{}, domain.ErrCartNotFound
	}

	return SUCCESS_CART, nil
}

type MockProductServiceClient struct{}

func (c *MockProductServiceClient) GetProduct(ctx context.Context, in *pb.GetProductRequest, opts ...grpc.CallOption) (*pb.GetProductResponse, error) {
	sku := in.GetSku()
	if sku != 1 {
		return nil, status.Error(codes.NotFound, "product not found")
	}

	return &pb.GetProductResponse{
		Sku:   sku,
		Name:  "test_product",
		Price: 100,
	}, nil
}

func TestAddCartItem(t *testing.T) {
	t.Parallel()
	
	type testAddCartItemInput struct {
		userIDStr string
		skuIDStr  string
		count     uint32
	}

	tests := []struct {
		name          string
		cartStorage   *MockStorage
		productClient *MockProductServiceClient
		in            testAddCartItemInput
		expectedError error
	}{
		{
			name:          "default",
			cartStorage:   &MockStorage{},
			productClient: &MockProductServiceClient{},
			in: testAddCartItemInput{
				userIDStr: "1",
				skuIDStr:  "1",
				count:     2,
			},
		},
		{
			name:          "invalid user id",
			cartStorage:   &MockStorage{},
			productClient: &MockProductServiceClient{},
			in: testAddCartItemInput{
				userIDStr: "-100",
				skuIDStr:  "1",
				count:     2,
			},
			expectedError: domain.ErrInvalidUserID,
		},
		{
			name:          "invalid sku count",
			cartStorage:   &MockStorage{},
			productClient: &MockProductServiceClient{},
			in: testAddCartItemInput{
				userIDStr: "1",
				skuIDStr:  "1",
				count:     0,
			},
			expectedError: domain.ErrInvalidSkuCount,
		},
		{
			name:          "products service AddCartItem error",
			cartStorage:   &MockStorage{},
			productClient: &MockProductServiceClient{},
			in: testAddCartItemInput{
				userIDStr: "1",
				skuIDStr:  "2",
				count:     0,
			},
			expectedError: domain.ErrInvalidSkuCount,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			service := NewCartService(test.cartStorage, test.productClient)
			err := service.AddCartItem(t.Context(), test.in.userIDStr, test.in.skuIDStr, test.in.count)
			require.Equal(t, err, test.expectedError)
		})
	}
}

func TestGetCart(t *testing.T) {
	t.Parallel()
	
	tests := []struct {
		name                 string
		cartStorage          *MockStorage
		productClient        *MockProductServiceClient
		userIDStr            string
		expectedCartResponse domain.Cart
		expectedError        error
	}{
		{
			name:          "default",
			cartStorage:   &MockStorage{},
			productClient: &MockProductServiceClient{},
			userIDStr:     "1",
			expectedCartResponse: SUCCESS_CART,

		},
		{
			name:          "invalid user id",
			cartStorage:   &MockStorage{},
			productClient: &MockProductServiceClient{},
			expectedError: domain.ErrInvalidUserID,
		},
		{
			name:          "products service GetCart error",
			cartStorage:   &MockStorage{},
			productClient: &MockProductServiceClient{},
			userIDStr:     "2",
			expectedError: domain.ErrCartNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			service := NewCartService(test.cartStorage, test.productClient)
			cart, err := service.GetCart(test.userIDStr)
			require.ErrorIs(t, err, test.expectedError)
			require.Equal(t, cart, test.expectedCartResponse)
		})
	}
}
