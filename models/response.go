package models

import "errors"

var ErrProductNotFound = errors.New("product not found")
var ErrCartItemNotFound = errors.New("cart item not found")
var ErrInsufficientStock = errors.New("insufficient stock for this quantity")

type ErrorResponse struct {
	Error any `json:"error"`
}

// SuccessResponse untuk response sukses standar (tanpa data)
type SuccessResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}
