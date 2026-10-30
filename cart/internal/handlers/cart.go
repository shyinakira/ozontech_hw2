package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sort"

	"hw2/cart/internal/domain"
	"hw2/cart/internal/services"
)

type CartHandler struct {
	service *services.CartService
}

func NewCartHandler(s *services.CartService) *CartHandler {
	return &CartHandler{
		service: s,
	}
}

type cartItem struct {
	SKU   int64  `json:"sku" validate:"required,min=1"`
	Name  string `json:"name" validate:"required"`
	Count int64  `json:"count" validate:"required,min=1"`
	Price int32  `json:"price" validate:"required,min=0"`
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

	skuStr := r.PathValue("sku")

	var body addCartItemBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sendError(handlerErrorResponse{
			Status: http.StatusBadRequest,
			URL:    r.URL.Path,
			Msg:    "Invalid request body",
		}, w)
		return
	}

	err := h.service.AddCartItem(r.Context(), userIDStr, skuStr, uint32(body.Count))
	if err != nil {
		if errors.Is(err, domain.ErrInvalidUserID) {
			sendError(handlerErrorResponse{
				Status: http.StatusBadRequest,
				URL:    r.URL.Path,
				Msg:    "Invalid user id",
			}, w)
		} else if errors.Is(err, domain.ErrInvalidSkuID) {
			sendError(handlerErrorResponse{
				Status: http.StatusBadRequest,
				URL:    r.URL.Path,
				Msg:    "Invalid sku id",
			}, w)
		} else if errors.Is(err, domain.ErrInvalidSkuCount) {
			sendError(handlerErrorResponse{
				Status: http.StatusBadRequest,
				URL:    r.URL.Path,
				Msg:    "Invalid sku count",
			}, w)
		} else if errors.Is(err, domain.ErrProductNotFound) {
			sendError(handlerErrorResponse{
				Status: http.StatusPreconditionFailed,
				URL:    r.URL.Path,
				Msg:    "Unknown product",
			}, w)
		} else {
			sendError(handlerErrorResponse{
				Status: http.StatusInternalServerError,
				URL:    r.URL.Path,
				Msg:    "Internal error",
			}, w)
		}
		return
	}

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
			Price: int32(item.Price),
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

	cart, err := h.service.GetCart(userIDStr)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidUserID) {
			sendError(handlerErrorResponse{
				Status: http.StatusBadRequest,
				URL:    r.URL.Path,
				Msg:    "Invalid user id",
			}, w)
		} else if errors.Is(err, domain.ErrCartNotFound) {
			sendError(handlerErrorResponse{
				Status: http.StatusNotFound,
				URL:    r.URL.Path,
				Msg:    "Cart not found",
			}, w)
		} else if errors.Is(err, domain.ErrProductNotFound) {
			sendError(handlerErrorResponse{
				Status: http.StatusPreconditionFailed,
				URL:    r.URL.Path,
				Msg:    "Product not found",
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

	res := buildCartResponse(cart)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf("unexpected handler error: %v", err)
		sendError(handlerErrorResponse{
			Status: http.StatusInternalServerError,
			URL:    r.URL.Path,
			Msg:    "Internal Server Error",
		}, w)
	}
}
