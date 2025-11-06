package routes

import (
	"github.com/RehanAthallahAzhar/shopeezy-accounts/internal/handlers"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/internal/middlewares"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/internal/services/token"
	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo, api *handlers.UserHandler, tokenService token.TokenService) {
	e.Static("/static", "template")

	// without token
	e.POST("/api/v1/accounts/register", api.RegisterUser)
	e.POST("/api/v1/accounts/login", api.Login)

	// Logout Endpoint (requires token to be blacklisted, but not validated by this middleware)
	// JWT parsing and blacklist logic is handled within the handler.Logout
	e.POST("/api/v1/accounts/logout", api.Logout)

	// Requires JWT Authentication
	// Create JWT authentication middleware
	jwtAuthMiddleware := middlewares.AuthMiddleware(middlewares.AuthMiddlewareOptions{
		TokenService: tokenService,
	})

	accountProtectedGroup := e.Group("/api/v1/accounts")
	accountProtectedGroup.Use(jwtAuthMiddleware) // Apply JWT middleware
	{
		// all users
		accountProtectedGroup.GET("/profile", api.GetUserProfile)
		accountProtectedGroup.PUT("/update", api.UpdateUser)
		accountProtectedGroup.DELETE("/delete/:id", api.DeleteUser)

		// admin
		accountProtectedGroup.GET("/list", api.GetAllUsers, middlewares.RequireRoles("admin"))
		accountProtectedGroup.GET("/:id", api.GetUserById, middlewares.RequireRoles("admin"))
	}
}
