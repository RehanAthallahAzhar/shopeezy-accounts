package handlers

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/RehanAthallahAzhar/shopeezy-accounts/helpers"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/models"
	"github.com/labstack/echo/v4"
)

func (api *API) RegisterUser(c echo.Context) error {
	req := new(models.UserAuthRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Permintaan tidak valid"})
	}

	if req.Role != "user" && req.TokenRole != "secret password from HRD or IT Manager" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Role token salah"})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Logger().Errorf("Gagal hash password: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Gagal memproses registrasi"})
	}

	newUserID := helpers.GenerateNewUserID()
	newUser := &models.User{
		ID:       newUserID, // Set ID
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "user", // role default
	}

	// Simpan pengguna ke database
	if err := api.UserRepo.CreateUser(c.Request().Context(), newUser); err != nil {
		if err == gorm.ErrDuplicatedKey { // Contoh penanganan jika username sudah ada
			return c.JSON(http.StatusConflict, map[string]string{"message": "Username sudah terdaftar"})
		}
		c.Logger().Errorf("Gagal menyimpan user baru: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Gagal menyimpan user"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Registrasi berhasil"})
}
