package test
import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)


type CartClient struct {
	URL    string
	httpClient *http.Client
}

func NewCartClient() *CartClient {
	return &CartClient{
		URL: "http://localhost:8080",
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type AddItemRequestBody struct {
	Count int64 `json:"count"`
}

type CartItemResponse struct {
	SKU   int64  `json:"sku"`
	Name  string `json:"name"`
	Count uint32 `json:"count"`
	Price uint32 `json:"price"`
}

type CartResponse struct {
	Items      []CartItemResponse `json:"items"`
	TotalPrice uint32             `json:"total_price"`
}

func (c *CartClient) AddItem(ctx context.Context, userID, sku int64, body any) (*http.Response, error) {
	url := fmt.Sprintf("%s/user/%d/cart/%d", c.URL, userID, sku)
	
	var reqBody io.Reader
	if body != nil {
		switch v := body.(type) {
		case string:
			reqBody = bytes.NewBufferString(v)
		default:
			data, err := json.Marshal(body)
			if err != nil {
				return nil,  fmt.Errorf("failed to marshal request: %w", err)
			}
			reqBody = bytes.NewReader(data)
		}
	}
	
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	
	return resp, nil
}

func (c *CartClient) GetCart(ctx context.Context, userID int64) (*http.Response, *CartResponse, error) {
	url := fmt.Sprintf("%s/user/%d/cart", c.URL, userID)
	
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create new request: %w", err)
	}
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to execute request: %w", err)
	}
	
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return resp, nil, nil
	}
	
	var cart CartResponse
	if err := json.Unmarshal(respBytes, &cart); err != nil {
		return resp, nil, fmt.Errorf("failed to unmarshal cart response: %w", err)
	}
	
	return resp, &cart, nil
}
func generateUserID() int64 {
	return time.Now().UnixNano()
}