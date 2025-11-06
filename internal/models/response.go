package models

// ErrorResponse untuk response error standar

type ErrorResponse struct {
	Error any `json:"error"`
}

// SuccessResponse untuk response sukses standar (tanpa data)
type SuccessResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}
