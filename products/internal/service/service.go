package service

import (
	"context"
	"errors"

	pb "hw2/products/gen/go/proto"
	"hw2/products/internal/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Storage interface {
	GetProduct(id int64) (domain.Product, error)
}

type ProductsService struct {
	pb.UnimplementedProductServiceServer
	products Storage
}

func NewProductService(s Storage) (*ProductsService, error) {
	return &ProductsService{
		products: s,
	}, nil
}

func (s *ProductsService) GetProduct(ctx context.Context, r *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	sku := r.GetSku()
	if sku <= 0 {
		return nil, status.Error(codes.InvalidArgument, "sku must be greater than 0")
	}
	
	product, err := s.products.GetProduct(sku)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "product not found")
		}
		return nil, status.Error(codes.Unknown, "internal error")
	}

	return &pb.GetProductResponse{
		Sku: sku,
		Name: product.Name,
		Price: int32(product.Price),
	}, nil
}
