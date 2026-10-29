package service

import (
	"testing"

	"github.com/stretchr/testify/require"

	pb "hw2/products/gen/go/proto"
	"hw2/products/internal/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MockStorage struct {
	storage map[int64]domain.Product
}

func (s *MockStorage) GetProduct(id int64) (domain.Product, error) {
	p, ok := s.storage[id]
	if ok != true {
		return domain.Product{}, domain.ErrNotFound
	}

	return p, nil
}

type ServiceTestCase struct {
	name             string
	storageState     *MockStorage
	req              *pb.GetProductRequest
	expectedResponse *pb.GetProductResponse
	expectedError    error
}

func TestService(t *testing.T) {
	tests := []ServiceTestCase{
		{
			name: "default",
			storageState: &MockStorage{
				storage: map[int64]domain.Product{
					1: domain.Product{
						Id:    1,
						Name:  "product_1",
						Price: 100,
					},
				},
			},
			req: &pb.GetProductRequest{
				Sku: 1,
			},
			expectedResponse: &pb.GetProductResponse{
				Sku: 1,
				Name:  "product_1",
				Price: 100,
			},
		},
		{
			name: "not found",
			storageState: &MockStorage{
				storage: map[int64]domain.Product{
					1: domain.Product{
						Id:    1,
						Name:  "produc1_1",
						Price: 100,
					},
				},
			},
			req: &pb.GetProductRequest{
				Sku: 2,
			},
			expectedError: status.Error(codes.NotFound, "product not found"),
		},
		{
			name: "wrong sku",
			storageState: &MockStorage{
				storage: map[int64]domain.Product{
					1: domain.Product{
						Id:    1,
						Name:  "produc1_1",
						Price: 100,
					},
				},
			},
			req: &pb.GetProductRequest{
				Sku: -1,
			},
			expectedError: status.Error(codes.InvalidArgument, "sku must be greater than 0"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service, err := NewProductService(test.storageState)
			if err != nil {
				require.Fail(t, "failde to create service")
			}

			resp, err := service.GetProduct(t.Context(), test.req)
			require.Equal(t, test.expectedResponse, resp)
			require.Equal(t, test.expectedError, err)
		})
	}
}
