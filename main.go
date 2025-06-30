package main

import (
	"context"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/RehanAthallahAzhar/shopeezy-accounts/databases"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/handlers"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/models"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/pkg/redisclient"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/repositories"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/routes"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/services"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/services/token"

	grpcServer "github.com/RehanAthallahAzhar/shopeezy-accounts/grpc"
	authpb "github.com/RehanAthallahAzhar/shopeezy-protos/proto/auth"
)

func main() {
	// REST API
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file: " + err.Error())
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable not set. Please set it for JWT signing.")
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

	conn, err := dbInstance.NewDB(ctx, &dbCredential)
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}
	defer func() {
		sqlDB, err := conn.DB()
		if err != nil {
			log.Printf("Error getting underlying DB: %v", err)
			return
		}
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	err = conn.AutoMigrate(&models.User{})
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	// Redis
	redisClient, err := redisclient.NewRedisClient()
	if err != nil {
		log.Fatalf("Failed to Inilialization redis client : %v", err)
	}
	defer redisClient.Close() // Make sure the Redis connection is closed

	// repo
	usersRepo := repositories.NewUserRepository(conn)
	jwtBlacklistRepo := repositories.NewJWTBlacklistRepository(redisClient)

	validate := validator.New()

	// svc
	tokenService := token.NewJWTTokenService(jwtSecret, jwtBlacklistRepo)
	userService := services.NewUserService(usersRepo, validate, tokenService, jwtBlacklistRepo)

	// gRPC
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen for gRPC server: %s: %v", grpcPort, err)
	}

	s := grpc.NewServer()
	authpb.RegisterAuthServiceServer(s, grpcServer.NewAuthServer(tokenService))
	reflection.Register(s)

	log.Printf("gRPC server for Account service is listening on port %s", lis.Addr())

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	e := echo.New()

	// route
	handler := handlers.NewHandler(usersRepo, userService, tokenService, jwtBlacklistRepo)
	routes.InitRoutes(e, handler, tokenService)

	// Start Echo API REST Server (Block main goroutine)
	log.Printf("Server REST API Echo mendengarkan di port 1324")
	e.Logger.Fatal(e.Start(":1324"))

	/*
		e.Start(echoPort) -> Ini adalah fungsi pemblokir (blocking function).
			Begitu Anda memanggilnya, fungsi ini akan mengambil alih main() dan akan terus berjalan tanpa henti
			untuk mendengarkan permintaan HTTP.
	*/
}
