package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/rehanazhar/account-cashier-app/handlers"
	"github.com/rehanazhar/account-cashier-app/middlewares"
)

func InitRoutes(e *echo.Echo, api *handlers.API) {
	e.Static("/static", "template")

	publicGroup := e.Group("/user")
	{
		publicGroup.POST("/create", api.Register())
		publicGroup.POST("/login", api.Login())
	}

	protectedUserGroup := e.Group("/user-utils", middlewares.AuthMiddleware)
	{
		protectedUserGroup.PUT("/update/:id", api.UpdateUser)
		protectedUserGroup.DELETE("/delete/:id", api.DeleteUser)
		protectedUserGroup.GET("/findall", api.FindAllUsers)
		protectedUserGroup.GET("/find/:id", api.FindUserById)
	}
}
