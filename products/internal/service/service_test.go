package service

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	pb "hw2/products/gen/go/proto"
	"hw2/products/internal/domain"
)

type MockStorage struct {
	storage map[int64]domain.Product
}

func (s *MockStorage) GetProduct(id int64) (domain.Product, error) {
	p, ok := s.storage[id]
	if !ok {
		return domain.Product{}, domain.ErrNotFound
	}

	return p, nil
}

func TestService(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		storageMap       map[int64]domain.Product
		req              *pb.GetProductRequest
		expectedResponse *pb.GetProductResponse
		expectedCode     codes.Code
	}{
		{
			name: "happy path",
			storageMap: map[int64]domain.Product{
				1: {
					Id:    1,
					Name:  "product_1",
					Price: 100,
				},
			},
			req: &pb.GetProductRequest{
				Sku: 1,
			},
			expectedResponse: &pb.GetProductResponse{
				Sku:   1,
				Name:  "product_1",
				Price: 100,
			},
			expectedCode: codes.OK,
		},
		{
			name: "not found",
			storageMap: map[int64]domain.Product{
				1: {
					Id:    1,
					Name:  "product_1",
					Price: 100,
				},
			},
			req: &pb.GetProductRequest{
				Sku: 2,
			},
			expectedCode: codes.NotFound,
		},
		{
			name: "invalid sku zero",
			storageMap: map[int64]domain.Product{
				1: {
					Id:    1,
					Name:  "product_1",
					Price: 100,
				},
			},
			req: &pb.GetProductRequest{
				Sku: 0,
			},
			expectedCode: codes.InvalidArgument,
		},
		{
			name: "invalid sku negative",
			storageMap: map[int64]domain.Product{
				1: {
					Id:    1,
					Name:  "product_1",
					Price: 100,
				},
			},
			req: &pb.GetProductRequest{
				Sku: -1,
			},
			expectedCode: codes.InvalidArgument,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			mockStorage := &MockStorage{storage: test.storageMap}
			service, err := NewProductService(mockStorage)
			require.NoError(t, err)

			resp, err := service.GetProduct(t.Context(), test.req)

			require.Equal(t, test.expectedCode, status.Code(err))

			if test.expectedCode == codes.OK {
				require.True(t, proto.Equal(test.expectedResponse, resp))
			}
		})
	}
}
