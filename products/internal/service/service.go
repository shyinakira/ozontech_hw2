package service

import (
	"context"
	"errors"
	"fmt"

	pb "hw2/products/gen/go/proto"
	"hw2/products/internal/domain"
	"hw2/products/internal/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductsService struct {
	pb.UnimplementedProductServiceServer
	products *storage.Storage
}

func NewProductService() (*ProductsService, error) {
	s, err := storage.NewStorage()
	if err != nil {
		return nil, fmt.Errorf("failed to create storage: %w", err)
	}
	return &ProductsService{
		products: s,
	}, nil
}

func (s *ProductsService) GetProduct(ctx context.Context, r *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	sku := r.GetSku()
	product, err := s.products.GetProduct(sku)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "product not found")
		}
	}

	return &pb.GetProductResponse{
		Sku: sku,
		Name: product.Name,
		Price: int64(product.Price),
	}, nil
}
