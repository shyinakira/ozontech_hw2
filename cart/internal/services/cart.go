package services

import (
	"context"
	"fmt"
	"strconv"

	"hw2/cart/internal/domain"
	pb "hw2/products/gen/go/proto"
)

type Storage interface {
	AddCartItem(userID domain.UserID, item domain.CartItem)
	GetCart(userID domain.UserID) (domain.Cart, error)
}

type CartService struct {
	storage        Storage
	productsClient pb.ProductServiceClient
}

func NewCartService(s Storage, pc pb.ProductServiceClient) *CartService {
	return &CartService{
		storage:        s,
		productsClient: pc,
	}
}

func (h *CartService) AddCartItem(ctx context.Context, userIDStr string, skuIDStr string, count uint32) error {
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		return domain.ErrInvalidUserID
	}

	sku, err := strconv.Atoi(skuIDStr)
	if err != nil || sku <= 0 {
		return domain.ErrInvalidSkuID
	}

	if count <= 0 {
		return domain.ErrInvalidSkuCount
	}

	product, err := h.productsClient.GetProduct(ctx, &pb.GetProductRequest{Sku: int64(sku)})
	if err != nil {
		return fmt.Errorf("failed to check product: %w", err)
	}

	cartItem := domain.CartItem{
		SKU:   int64(sku),
		Name:  product.Name,
		Price: uint32(product.Price),
		Count: count,
	}

	h.storage.AddCartItem(domain.UserID(userID), cartItem)

	return nil
}

func (h *CartService) GetCart(userIDStr string) (domain.Cart, error) {
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		return domain.Cart{}, domain.ErrInvalidUserID
	}

	cart, err := h.storage.GetCart(domain.UserID(userID))
	if err != nil {
		return domain.Cart{}, fmt.Errorf("failed to get user cart from storage: %w", err)
	}

	return cart, nil
}
