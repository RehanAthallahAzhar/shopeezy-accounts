package handlers

import (
	stdErrors "errors"
	"fmt"
	"log"
	"net/http"

	"github.com/RehanAthallahAzhar/shopeezy-accounts/models"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/pkg/errors"
	"github.com/labstack/echo/v4"
)

const (
	MsgUserRetrieved  = "User retrieved successfully"
	MsgUserCreated    = "User created successfully"
	MsgUserUpdated    = "User updated successfully"
	MsgUserDeleted    = "User deleted successfully"
	MsgUsersRetrieved = "Users retrieved successfully"
	MsgLogin          = "Login successful"
	MsgLogout         = "Logout successful"
)

func respondSuccess(c echo.Context, status int, message string, data interface{}) error {
	return c.JSON(status, models.SuccessResponse{
		Message: message,
		Data:    data,
	})
}

func respondError(c echo.Context, status int, err error) error {
	return c.JSON(status, models.ErrorResponse{
		Error: err.Error(),
	})
}

func respondWithUser(c echo.Context, user *models.UserResponse, msg string, err error) error {
	if err != nil {
		return handleServiceError(c, err)
	}
	return respondSuccess(c, http.StatusOK, MsgUserRetrieved, user)
}

func respondWithUsers(c echo.Context, users []models.UserResponse, err error) error {
	if err != nil {
		return handleServiceError(c, err)
	}
	if len(users) == 0 {
		return respondError(c, http.StatusNotFound, errors.ErrUserNotFound)
	}
	return respondSuccess(c, http.StatusOK, MsgUsersRetrieved, users)
}

func handleServiceError(c echo.Context, err error) error {
	switch {
	case stdErrors.Is(err, errors.ErrUserNotFound):
		return respondError(c, http.StatusNotFound, err)
	case stdErrors.Is(err, errors.ErrInternalServerError):
		return respondError(c, http.StatusForbidden, err)
	default:
		log.Printf("Service error: %v", err)
		return respondError(c, http.StatusInternalServerError, fmt.Errorf("%w: %s", errors.ErrInternalServerError, err))
	}
}
