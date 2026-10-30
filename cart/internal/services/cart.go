package services

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"hw2/cart/internal/domain"
	pb "hw2/products/gen/go/proto"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	product, err := h.productsClient.GetProduct(reqCtx, &pb.GetProductRequest{Sku: int64(sku)})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return domain.ErrProductNotFound
		}
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

func (h *CartService) checkCartItem(ctx context.Context, item domain.CartItem) (domain.CartItem, error) {
	product, err := h.productsClient.GetProduct(ctx, &pb.GetProductRequest{Sku: int64(item.SKU)})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return domain.CartItem{}, domain.ErrProductNotFound
		}
		return domain.CartItem{}, fmt.Errorf("failed to get product: %w", err)
	}

	actualProduct := domain.CartItem{
		SKU:   product.Sku,
		Name:  product.Name,
		Price: uint32(product.Price),
		Count: item.Count,
	}

	return actualProduct, nil
}

func (h *CartService) GetCart(userIDStr string) (domain.Cart, error) {
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		return domain.Cart{}, domain.ErrInvalidUserID
	}

	localCart, err := h.storage.GetCart(domain.UserID(userID))
	if err != nil {
		return domain.Cart{}, fmt.Errorf("failed to get user cart from storage: %w", err)
	}

	if len(localCart.Items) == 0 {
		return domain.Cart{}, domain.ErrCartNotFound
	}

	actualItems := make([]domain.CartItem, len(localCart.Items))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)
	for i := 0; i < len(localCart.Items); i++ {
		idx := i
		item := localCart.Items[i]
		eg.Go(func() error {
			actualProduct, err := h.checkCartItem(ctx, item)
			if err != nil {
				return err
			}

			actualItems[idx] = actualProduct
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return domain.Cart{}, err
	}

	totalPrice := uint32(0)
	for _, item := range actualItems {
		totalPrice += item.Count * item.Price
	}

	return domain.Cart{
		Items:      actualItems,
		TotalPrice: totalPrice,
	}, nil
}
