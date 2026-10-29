package domain

import "errors"

var (
	ErrInvalidUserID = errors.New("invalid user id. Must be int64, > 0")
	ErrInvalidSkuID  = errors.New("invalid sku id. Must be int64, > 0")
	ErrInvalidSkuCount      = errors.New("invalid sku count. Must be > 0")
	ErrCartNotFound  = errors.New("user not found")
)
