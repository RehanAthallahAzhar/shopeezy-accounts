package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rehanazhar/shopeezy-account/helpers"
	"github.com/rehanazhar/shopeezy-account/models"
)

// Register menangani proses registrasi user baru
func (api *API) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		// Inisialisasi struct untuk menangkap data request
		var req = models.UserCreateRequest{}

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid JSON format"})
		}

		// Buat data user baru dengan password yang sudah di-hash
		user := models.User{
			Name:     req.Name,
			Username: req.Username,
			Email:    req.Email,
			Password: helpers.HashPassword(req.Password),
		}

		err := api.UserRepo.CreateUser(ctx, &user)
		if err != nil {
			if errors.Is(err, models.ErrProductNotFound) {
				return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
			}
			return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to retrieve user"})
		}

		return c.JSON(http.StatusCreated, models.SuccessResponse{
			Message: "User created successfully",
			Data: models.UserResponse{
				Id:        user.Id,
				Name:      user.Name,
				Username:  user.Username,
				Email:     user.Email,
				CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
			},
		})
	}
}
