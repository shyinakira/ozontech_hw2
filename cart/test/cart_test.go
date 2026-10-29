package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCartService_HappyPath(t *testing.T) {
	client := NewCartClient()
	userID := generateUserID()

	res, _, err := client.GetCart(t.Context(), userID)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, res.StatusCode)

	res, err = client.AddItem(t.Context(), userID, 1, AddItemRequestBody{Count: 2})
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode)

	res, err = client.AddItem(t.Context(), userID, 2, AddItemRequestBody{Count: 1})
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode)

	res, cart, err := client.GetCart(t.Context(), userID)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	require.Equal(t, int64(1), cart.Items[0].SKU)
	require.Equal(t, "Product_1", cart.Items[0].Name)
	require.Equal(t, uint32(2), cart.Items[0].Count)
	require.Equal(t, uint32(1000), cart.Items[0].Price)

	require.Equal(t, int64(2), cart.Items[1].SKU)
	require.Equal(t, "Product_2", cart.Items[1].Name)
	require.Equal(t, uint32(1), cart.Items[1].Count)
	require.Equal(t, uint32(2000), cart.Items[1].Price)

	require.Equal(t, uint32(4000), cart.TotalPrice)
}

func TestCartService_InvalidData(t *testing.T) {
	type reqData struct {
		userID int64
		skuID  int64
		body   AddItemRequestBody
	}

	tests := []struct {
		name                string
		req                 reqData
		expectedErrorStatus int
	}{
		{
			name: "invalid user id",
			req: reqData{
				userID: -100,
				skuID: 1,
				body: AddItemRequestBody{Count: 1},
			},
			expectedErrorStatus: http.StatusBadRequest,
		},
		{
			name: "invalid sku id",
			req: reqData{
				userID: generateUserID(),
				skuID: -100,
				body: AddItemRequestBody{Count: 1},
			},
			expectedErrorStatus: http.StatusBadRequest,
		},
		{
			name: "invalid sku count",
			req: reqData{
				userID: generateUserID(),
				skuID: 1,
				body: AddItemRequestBody{Count: 0},
			},
			expectedErrorStatus: http.StatusBadRequest,
		},
	}

	client := NewCartClient()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := client.AddItem(t.Context(), test.req.userID, test.req.skuID, test.req.body)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			require.Equal(t, res.StatusCode, test.expectedErrorStatus)
		})
	}
}

func TestCartService_NonExistedProduct(t *testing.T) {
	client := NewCartClient()
	
	userID := generateUserID()

	res, err := client.AddItem(t.Context(), userID, 5, AddItemRequestBody{Count: 2})
	require.NoError(t, err)
	require.Equal(t, http.StatusPreconditionFailed, res.StatusCode)
}