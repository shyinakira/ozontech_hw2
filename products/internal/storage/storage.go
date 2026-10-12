package storage

import (
	"encoding/json"
	"fmt"
	"hw2/products/internal/domain"
	"os"
)

const PRODUCTS_DATA_PATH = "../../products.json"

type Product struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type Storage struct {
	products map[int64]Product
}

func readProductsData() (map[int64]Product, error) {
	file, err := os.Open(PRODUCTS_DATA_PATH)
	if err != nil {
		return nil, fmt.Errorf("failed to read products data: %w", err)
	}

	storage := make(map[int64]Product)
	if err := json.NewDecoder(file).Decode(&storage); err != nil {
		return nil, fmt.Errorf("failed to decode products data: %w", err)
	}

	return storage, nil
}

func NewStorage() (*Storage, error) {
	storage, err := readProductsData()
	if err != nil {
		return nil, err
	}
	return &Storage{
		products: storage,
	}, nil
}

func (s *Storage) GetProduct(id int64) (Product, error) {
	product, ok := s.products[id]
	if !ok {
		return Product{}, fmt.Errorf("unknown id %d: %w", id, domain.ErrNotFound)
	}

	return product, nil
}