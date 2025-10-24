package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/RehanAthallahAzhar/shopeezy-accounts/db"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/internal/configs"
	dbGenerated "github.com/RehanAthallahAzhar/shopeezy-accounts/internal/db"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/internal/handlers"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/internal/helpers"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/internal/models"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/internal/pkg/logger"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/internal/pkg/redisclient"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/internal/repositories"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/internal/routes"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/internal/services"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/internal/services/token"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	grpcServer "github.com/RehanAthallahAzhar/shopeezy-accounts/internal/grpc"
	accountpb "github.com/RehanAthallahAzhar/shopeezy-protos/pb/account"
	authpb "github.com/RehanAthallahAzhar/shopeezy-protos/pb/auth"
)

func main() {
	log := logger.NewLogger()

	cfg, err := configs.LoadConfig(log)
	if err != nil {
		log.Fatalf("FATAL: Gagal memuat konfigurasi: %v", err)
	}

	dbCredential := models.Credential{
		Host:         cfg.Database.Host,
		Username:     cfg.Database.User,
		Password:     cfg.Database.Password,
		DatabaseName: cfg.Database.Name,
		Port:         cfg.Database.Port,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Setup DB
	conn, err := db.Connect(ctx, &dbCredential)
	if err != nil {
		log.Fatalf("DB connection error: %v", err)
	}

	// Migrations
	log.Println("Running database migrations...")

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		dbCredential.Username,
		dbCredential.Password,
		dbCredential.Host,
		dbCredential.Port,
		dbCredential.DatabaseName,
	)

	m, err := migrate.New(
		// "file://../../db/migrations", // Local DB
		"file://db/migrations", // Container DB
		connectionString,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Database migrations ran successfully.")

	defer conn.Close()

	// Init SQLC
	sqlcQueries := dbGenerated.New(conn)

	//jwt
	_ = helpers.NewJWTHelper(cfg.Server.JWTSecret)

	// Setup Redis
	redisClient, err := redisclient.NewRedisClient()
	if err != nil {
		log.Fatalf("Failed to Inilialization redis client : %v", err)
	}
	defer redisClient.Close() // Make sure the Redis connection is closed

	// Setup Repo
	usersRepo := repositories.NewUserRepository(sqlcQueries)
	jwtBlacklistRepo := repositories.NewJWTBlacklistRepository(redisClient)

	validate := validator.New()

	// Setup Service
	tokenService := token.NewJWTTokenService(cfg.Server.JWTSecret, jwtBlacklistRepo)
	userService := services.NewUserService(usersRepo, validate, tokenService, jwtBlacklistRepo)

	// Setup gRPC
	lis, err := net.Listen("tcp", ":"+cfg.Server.GRPCPort)
	if err != nil {
		log.Fatalf("Failed to listen for gRPC server: %s: %v", cfg.Server.GRPCPort, err)
	}

	s := grpc.NewServer()
	authpb.RegisterAuthServiceServer(s, grpcServer.NewAuthServer(tokenService))
	accountpb.RegisterAccountServiceServer(s, grpcServer.NewAccountServer(usersRepo))
	reflection.Register(s)

	log.Printf("gRPC server for Account service is listening on port %s", lis.Addr())

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	e := echo.New()

	// Setup Route
	handler := handlers.NewHandler(usersRepo, userService, tokenService, jwtBlacklistRepo)
	routes.InitRoutes(e, handler, tokenService)

	// Start Echo API REST Server (Block main goroutine)
	log.Printf("Server REST API Echo is listening on port %s", cfg.Server.Port)
	e.Logger.Fatal(e.Start(":" + cfg.Server.Port))

	/*
		e.Start(echoPort) -> Ini adalah fungsi pemblokir (blocking function).
			Begitu Anda memanggilnya, fungsi ini akan mengambil alih main() dan akan terus berjalan tanpa henti
			untuk mendengarkan permintaan HTTP.
	*/
}
