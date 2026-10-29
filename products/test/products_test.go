package test

import (
	"testing"

	pb "hw2/products/gen/go/proto"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func TestProductService(t *testing.T) {
	tests := []struct {
		name         string
		req          *pb.GetProductRequest
		expectedRes  *pb.GetProductResponse
		expectedCode codes.Code
	}{
		{
			name: "happy path",
			req: &pb.GetProductRequest{
				Sku: 1,
			},
			expectedRes: &pb.GetProductResponse{
				Sku:   1,
				Name:  "Product_1",
				Price: 1000,
			},
			expectedCode: codes.OK,
		},
		{
			name: "not existing product",
			req: &pb.GetProductRequest{
				Sku: 4,
			},
			expectedCode: codes.NotFound,
		},
		{
			name: "invalid sku",
			req: &pb.GetProductRequest{
				Sku: -12,
			},
			expectedCode: codes.InvalidArgument,
		},
	}

	connGrpc, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to create grpc conn :50051, %v\n", err)
	}
	defer connGrpc.Close()
	productsClient := pb.NewProductServiceClient(connGrpc)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			productRes, err := productsClient.GetProduct(t.Context(), test.req)

			require.Equal(t, test.expectedCode, status.Code(err))

			if test.expectedCode == codes.OK {
				require.True(t, proto.Equal(test.expectedRes, productRes), "protobuf responses are not equal")
			}
		})
	}
}