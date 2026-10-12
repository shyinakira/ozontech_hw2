package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"hw2/cart/internal/domain"
	"hw2/cart/internal/storage"
)

type handlerErrorResponse struct {
	Status int    `json:"status"`
	URL    string `json:"url"`
	Msg    string `json:"msg"`
}

func sendError(r handlerErrorResponse, w http.ResponseWriter) {
	w.WriteHeader(r.Status)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(r)
	if err != nil {
		log.Printf("failed to send structured response: %w\n", err)
	}
}

type addCartItemBody struct {
	Count uint32 `json:"count"`
}

func AddCartItem(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		sendError(handlerErrorResponse{
			Status: http.StatusBadRequest,
			URL:    r.URL.Path,
			Msg:    "User ID must be an int",
		}, w)
		return
	}

	skuStr := r.PathValue("sku")
	sku, err := strconv.Atoi(skuStr)
	if err != nil {
		sendError(handlerErrorResponse{
			Status: http.StatusBadRequest,
			URL:    r.URL.Path,
			Msg:    "SKU must be an int",
		}, w)
		return
	}

	var body addCartItemBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sendError(handlerErrorResponse{
			Status: http.StatusBadRequest,
			URL:    r.URL.Path,
			Msg:    "Invalid body request",
		}, w)
		return
	}

	count := body.Count
	//TODO: gRPC Interact
	cartItem := domain.CartItem{
		SKU:   sku,
		Name:  name,
		Price: price,
		Count: count,
	} //TODO: validate
	err = storage.AddCartItem(storage.UserID(userID), cartItem)

	if err != nil {
		if errors.Is(err, domain.ErrCartNotFound) {
			sendError(handlerErrorResponse{
				Status: http.StatusNotFound,
				URL:    r.URL.Path,
				Msg:    "User not found",
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

	w.WriteHeader(http.StatusCreated)
}

type cartItem struct {
	SKU   int64  `json:"sku"`
	Name  string `json:"name"`
	Count uint32 `json:"count"`
	Price uint32 `json:"price"`
}

type cartResponse struct {
	Items      []cartItem `json:"items"`
	TotalPrice uint32     `json:"total_price"`
}

func buildCartRespone(cart domain.Cart) cartResponse {
	res := cartResponse{}
	for _, item := range cart.Items {
		res.Items = append(res.Items, cartItem{
			SKU: item.SKU,
			Name: item.Name,
			Count: item.Count,
			Price: item.Price,
		})
		res.TotalPrice += item.Price
	}
}

func GetCart(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		sendError(handlerErrorResponse{
			Status: http.StatusBadRequest,
			URL:    r.URL.Path,
			Msg:    "User ID must be an int",
		}, w)
		return
	}
	
	cart, err := storage.GetCart(storage.UserID(userID))
	if err != nil {
		if errors.Is(err, domain.ErrCartNotFound) {
			sendError(handlerErrorResponse{
				Status: http.StatusNotFound,
				URL:    r.URL.Path,
				Msg:    "User not found",
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

	res := buildCartRespone(cart)
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
