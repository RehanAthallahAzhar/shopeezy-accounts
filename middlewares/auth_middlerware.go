package middlewares

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/RehanAthallahAzhar/shopeezy-accounts/models"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/services/token"

	"github.com/labstack/echo/v4"
)

type AuthMiddlewareOptions struct {
	TokenService token.TokenService
}

// AuthMiddleware is an Echo middleware function to validate JWT tokens for REST API.
func AuthMiddleware(opts AuthMiddlewareOptions) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" || len(authHeader) < 7 || !strings.HasPrefix(authHeader, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Authentication token missing or invalid format"})
			}
			token := authHeader[7:]

			isValid, userID, username, userRole, errMsg, err := opts.TokenService.ValidateToken(context.Background(), token)
			if err != nil {
				log.Printf("Token validation error: %v", err)
				// If the error is due to an expired token or cryptographic invalidity,
				// TokenService already returns isValid=false and an appropriate errMsg.
				return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Server error while validating token"})
			}

			if !isValid {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token: " + errMsg})
			}

			c.Set("userID", userID)
			c.Set("username", username)
			c.Set("role", userRole)

			// Continue to the next handler
			return next(c)
		}
	}
}

func RequireRoles(allowedRoles ...string) echo.MiddlewareFunc {
	roleSet := make(map[string]struct{})
	for _, r := range allowedRoles {
		roleSet[r] = struct{}{}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role, ok := c.Get("role").(string)
			if !ok {
				return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Unauthorized"})
			}

			if _, allowed := roleSet[role]; !allowed {
				return c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Access denied"})
			}

			return next(c)
		}
	}
}
