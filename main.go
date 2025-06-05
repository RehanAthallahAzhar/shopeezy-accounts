package main

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	"github.com/rehanazhar/shopeezy-account/databases"
	"github.com/rehanazhar/shopeezy-account/handlers"
	"github.com/rehanazhar/shopeezy-account/models"
	"github.com/rehanazhar/shopeezy-account/repositories"
	"github.com/rehanazhar/shopeezy-account/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file: " + err.Error())
	}

	portStr := os.Getenv("DB_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic("Invalid DB_PORT in .env file or not set: " + err.Error())
	}

	dbCredential := models.Credential{
		Host:         os.Getenv("DB_HOST"),
		Username:     os.Getenv("DB_USER"),
		Password:     os.Getenv("DB_PASSWORD"),
		DatabaseName: os.Getenv("DB_NAME"),
		Port:         port,
	}

	dbInstance := databases.NewDB()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := dbInstance.Connect(ctx, &dbCredential)
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	err = conn.AutoMigrate(&models.User{}) // &models.User{}, &models.Session{}
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	e := echo.New()

	usersRepo := repositories.NewProductRepository(conn)
	handler := handlers.NewHandler(usersRepo)

	routes.InitRoutes(e, handler)
	// usersRepo := repositories.NewUserRepository(conn)
	// sessionsRepo := repositories.NewSessionsRepository(conn)
	// productsRepo := repositories.NewProductRepository(conn)
	// cartsRepo := repositories.NewCartRepository(conn, productsRepo)

	// handler := handlers.NewHandler(usersRepo, sessionsRepo, productsRepo, cartsRepo, accountGRPCClient) // <--- TERUSKAN gRPC CLIENT DI SINI

	// routes.InitRoutes(e, &handler)

	e.Logger.Fatal(e.Start(":1324"))
}
