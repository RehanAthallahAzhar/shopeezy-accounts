package middlewares

import (
	"context"
	"net/http"

	"github.com/RehanAthallahAzhar/shopeezy-accounts/services" // Impor services Anda
	"github.com/labstack/echo/v4"
)

// AuthMiddlewareOptions berisi dependensi untuk middleware autentikasi
type AuthMiddlewareOptions struct {
	TokenService services.TokenService
}

// AuthMiddleware adalah fungsi middleware untuk memvalidasi token untuk API REST.
func AuthMiddleware(opts AuthMiddlewareOptions) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Token otentikasi tidak ditemukan atau format salah"})
			}
			token := authHeader[7:]

			// Panggil TokenService untuk memvalidasi token
			isValid, userId, username, userRole, errMsg, err := opts.TokenService.Validate(context.Background(), token)
			if err != nil {
				//Todo: tambahkan kasus dmn token expired
				c.Logger().Errorf("Kesalahan validasi token: %v", err)
				return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Kesalahan server saat memvalidasi token"})
			}

			if !isValid {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Token tidak valid: " + errMsg})
			}

			// Jika token valid, Anda bisa menyimpan informasi pengguna di Echo Context
			c.Set("userID", userId)
			c.Set("username", username)
			c.Set("role", userRole)

			// Lanjutkan ke handler berikutnya
			return next(c)
		}
	}
}
