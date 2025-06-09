package routes

import (
	"github.com/RehanAthallahAzhar/shopeezy-accounts/handlers"    // Import handler Anda
	"github.com/RehanAthallahAzhar/shopeezy-accounts/middlewares" // Import middleware Anda
	"github.com/RehanAthallahAzhar/shopeezy-accounts/services"    // Import services Anda
	"github.com/labstack/echo/v4"
)

// InitRoutes menginisialisasi semua rute API untuk account-app
func InitRoutes(e *echo.Echo, api *handlers.API, tokenService services.TokenService) { // Tambahkan tokenService sebagai parameter
	// Static files (jika ada)
	e.Static("/static", "template")

	// Buat opsi middleware autentikasi untuk API REST account-app
	authOpts := middlewares.AuthMiddlewareOptions{
		TokenService: tokenService,
	}

	// Contoh grup rute yang memerlukan autentikasi
	// Asumsikan Anda memiliki rute untuk profil pengguna, update, dll.
	userGroup := e.Group("/users")
	userGroup.Use(middlewares.AuthMiddleware(authOpts)) // Terapkan middleware
	{
		// Contoh: rute yang memerlukan token
		userGroup.GET("/userlist", api.FindAllUsers)
		userGroup.GET("/user/:id", api.FindUserById)
		userGroup.PUT("/update", api.UpdateUser)
		userGroup.DELETE("/delete/:id", api.DeleteUser)
	}

	// Rute untuk autentikasi (login, register) TIDAK memerlukan middleware
	e.POST("/auth/register", api.RegisterUser)
	e.POST("/auth/login", api.LoginUser)

	// Tambahkan rute-rute lain sesuai kebutuhan account-app Anda
}
