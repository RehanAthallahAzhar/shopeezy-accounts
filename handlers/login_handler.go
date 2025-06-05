package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rehanazhar/account-cashier-app/helpers"
	"github.com/rehanazhar/account-cashier-app/models"
	"golang.org/x/crypto/bcrypt"
)

func (api *API) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		// Inisialisasi struct untuk menampung data dari request
		var req = models.UserLoginRequest{}

		// Validasi input dari request body menggunakan ShouldBindJSON
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Bad Request: Invalid JSON format"})
		}

		user, err := api.UserRepo.FindUserByUsername(ctx, req.Username)
		if err != nil {
			if errors.Is(err, models.ErrProductNotFound) {
				return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
			}
			return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to retrieve user"})
		}

		// Bandingkan password yang dimasukkan dengan password yang sudah di-hash di database
		// Jika tidak cocok, kirimkan respons error Unauthorized
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: helpers.TranslateErrorMessage(err),
			})
		}

		// Jika login berhasil, generate token untuk user
		token := helpers.GenerateToken(user.Username)

		// Kirimkan response sukses dengan status OK dan data user serta token
		return c.JSON(http.StatusOK, models.SuccessResponse{
			Message: "Login Success",
			Data: models.UserResponse{
				Id:        user.Id,
				Name:      user.Name,
				Username:  user.Username,
				Email:     user.Email,
				CreatedAt: user.CreatedAt.String(),
				UpdatedAt: user.UpdatedAt.String(),
				Token:     &token,
			},
		})
	}
}
