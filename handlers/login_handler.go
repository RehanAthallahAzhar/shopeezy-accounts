package handlers

import (
	"net/http"

	"gorm.io/gorm"

	"github.com/RehanAthallahAzhar/shopeezy-accounts/models"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (api *API) LoginUser(c echo.Context) error {
	req := new(models.UserAuthRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	user, err := api.UserRepo.FindUserByUsername(c.Request().Context(), req.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Username atau password salah"})
		}
		c.Logger().Errorf("Gagal mencari user: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "user not found"})
	}

	// password verification
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Username atau password salah"})
	}

	// Generate JWT
	tokenString, err := api.TokenSvc.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.Logger().Errorf("Gagal generate token JWT: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Gagal memproses login"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Login berhasil",
		"token":   tokenString,
	})
}
