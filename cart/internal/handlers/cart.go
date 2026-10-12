package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sort"
	"strconv"

	"hw2/cart/internal/domain"
	"hw2/cart/internal/storage"
	pb "hw2/products/gen/go/proto"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CartHandler struct {
	validator      *validator.Validate
	storage        *storage.Storage
	productsClient pb.ProductServiceClient
}

func NewCartHandler(s *storage.Storage, pc pb.ProductServiceClient) *CartHandler {
	return &CartHandler{
		validator: validator.New(),
		storage:        s,
		productsClient: pc,
	}
}

type cartItem struct {
	SKU   int64  `json:"sku" validate:"required,min=1"`
	Name  string `json:"name" validate:"required"`
	Count int64  `json:"count" validate:"required,min=1"`
	Price int64  `json:"price" validate:"required,min=0"`
}

func (i cartItem) ToDomain() domain.CartItem {
	return domain.CartItem{
		SKU:   i.SKU,
		Name:  i.Name,
		Count: uint32(i.Count),
		Price: uint32(i.Price),
	}
}

type addCartItemBody struct {
	Count int64 `json:"count" validate:"min=1"`
}

func (h *CartHandler) AddCartItem(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		sendError(handlerErrorResponse{
			Status: http.StatusBadRequest,
			URL:    r.URL.Path,
			Msg:    "Invalid user ID",
		}, w)
		return
	}

	skuStr := r.PathValue("sku")
	sku, err := strconv.Atoi(skuStr)
	if err != nil || sku <= 0  {
		sendError(handlerErrorResponse{
			Status: http.StatusBadRequest,
			URL:    r.URL.Path,
			Msg:    "Invalid SKU",
		}, w)
		return
	}

	var body addCartItemBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sendError(handlerErrorResponse{
			Status: http.StatusBadRequest,
			URL:    r.URL.Path,
			Msg:    "Invalid request body",
		}, w)
		return
	}

	count := body.Count

	product, err := h.productsClient.GetProduct(r.Context(), &pb.GetProductRequest{Sku: int64(sku)})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			sendError(handlerErrorResponse{
				Status: http.StatusPreconditionFailed,
				URL:    r.URL.Path,
				Msg:    "Product doesn't exists",
			}, w)
		} else {
			log.Printf("unexpected handler error: %v", err)
			sendError(handlerErrorResponse{
				Status: http.StatusInternalServerError,
				URL:    r.URL.Path,
				Msg:    "Internal Server Error",
			}, w)
		}
		return
	}

	cartItem := cartItem{
		SKU:   int64(sku),
		Name:  product.Name,
		Price: product.Price,
		Count: count,
	}
	err = h.validator.Struct(cartItem)
	if err != nil {
		sendError(handlerErrorResponse{
			Status: http.StatusBadRequest,
			URL:    r.URL.Path,
			Msg:    "Invalid item data",
		}, w)
		return
	}
	
	h.storage.AddCartItem(storage.UserID(userID), cartItem.ToDomain())

	w.WriteHeader(http.StatusCreated)
}

type cartResponse struct {
	Items      []cartItem `json:"items"`
	TotalPrice uint32     `json:"total_price"`
}

func buildCartResponse(cart domain.Cart) cartResponse {
	res := cartResponse{}
	for _, item := range cart.Items {
		res.Items = append(res.Items, cartItem{
			SKU:   item.SKU,
			Name:  item.Name,
			Count: int64(item.Count),
			Price: int64(item.Price),
		})
		res.TotalPrice += item.Price * item.Count
	}

	sort.Slice(res.Items, func(i, j int) bool {
		return res.Items[i].SKU < res.Items[j].SKU
	})
	return res
}

func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		sendError(handlerErrorResponse{
			Status: http.StatusBadRequest,
			URL:    r.URL.Path,
			Msg:    "Invalid user ID",
		}, w)
		return
	}

	if userID <= 0 {
		sendError(handlerErrorResponse{
			Status: http.StatusBadRequest,
			URL:    r.URL.Path,
			Msg:    "User ID must be greater than 0",
		}, w)
		return
	}

	cart, err := h.storage.GetCart(storage.UserID(userID))
	if err != nil {
		if errors.Is(err, domain.ErrCartNotFound) {
			sendError(handlerErrorResponse{
				Status: http.StatusNotFound,
				URL:    r.URL.Path,
				Msg:    "Cart not found",
			}, w)
		} else {
			log.Printf("unexpected handler error: %v", err)
			sendError(handlerErrorResponse{
				Status: http.StatusInternalServerError,
				URL:    r.URL.Path,
				Msg:    "Internal Server Error",
			}, w)
		}
		return
	}

	if len(cart.Items) == 0 {
		sendError(handlerErrorResponse{
			Status: http.StatusNotFound,
			URL:    r.URL.Path,
			Msg:    "Cart not found",
		}, w)
	}

	res := buildCartResponse(cart)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf("unexpected handler error: %v", err)
		sendError(handlerErrorResponse{
			Status: http.StatusInternalServerError,
			URL:    r.URL.Path,
			Msg:    "Internal Server Error",
		}, w)
	}

	w.WriteHeader(http.StatusOK)
}
