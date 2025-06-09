package main

import (
	"context"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	// Import package-package Echo API Anda
	"github.com/RehanAthallahAzhar/shopeezy-accounts/databases"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/handlers"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/models"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/repositories"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/routes"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/services"

	// Penting: Import path sekarang menunjuk ke repositori proto yang terpisah
	authpb "github.com/RehanAthallahAzhar/shopeezy-protos/proto/auth"
)

// server struct akan mengimplementasikan antarmuka AuthServiceServer
type authServer struct {
	authpb.UnimplementedAuthServiceServer
	TokenSvc services.TokenService // Tambahkan TokenService di sini
}

// ValidateToken adalah implementasi metode RPC ValidateToken
// ValidateToken adalah implementasi metode RPC ValidateToken
func (s *authServer) ValidateToken(ctx context.Context, req *authpb.ValidateTokenRequest) (*authpb.ValidateTokenResponse, error) {
	log.Printf("Menerima permintaan ValidateToken: %s", req.GetToken())

	token := req.GetToken()
	// KOREKSI: Dapatkan 'userRole' dari TokenService
	isValid, userId, username, userRole, errMsg, err := s.TokenSvc.Validate(ctx, token)
	if err != nil {
		log.Printf("Internal error during token validation: %v", err)
		return &authpb.ValidateTokenResponse{
			IsValid:      false,
			ErrorMessage: "Kesalahan internal server",
		}, status.Errorf(codes.Internal, "internal server error")
	}

	if !isValid {
		return &authpb.ValidateTokenResponse{
			IsValid:      false,
			ErrorMessage: errMsg,
		}, status.Errorf(codes.Unauthenticated, errMsg)
	}

	return &authpb.ValidateTokenResponse{
		IsValid:  true,
		UserId:   userId,
		Username: username,
		Role:     userRole, // Teruskan Role di respons gRPC
	}, nil
}
func main() {
	// --- Inisialisasi dan Konfigurasi REST API Echo ---
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file: " + err.Error())
	}

	dbPortStr := os.Getenv("DB_PORT")
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		panic("Invalid DB_PORT in .env file or not set: " + err.Error())
	}

	dbCredential := models.Credential{
		Host:         os.Getenv("DB_HOST"),
		Username:     os.Getenv("DB_USER"),
		Password:     os.Getenv("DB_PASSWORD"),
		DatabaseName: os.Getenv("DB_NAME"),
		Port:         dbPort,
	}

	dbInstance := databases.NewDB()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := dbInstance.Connect(ctx, &dbCredential)
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// Perbarui automigrate untuk model User
	err = conn.AutoMigrate(&models.User{})
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	// --- Inisialisasi Repositories dan Services ---
	userRepo := repositories.NewUserRepository(conn)   // Pastikan ini NewUserRepository
	tokenService := services.NewTokenService(userRepo) // Inisialisasi TokenService

	// --- Inisialisasi Handlers untuk Echo API ---
	apiHandler := handlers.NewHandler(userRepo, tokenService) // Teruskan UserRepo dan TokenService

	e := echo.New()

	// --- Inisialisasi Rute Echo API ---
	routes.InitRoutes(e, apiHandler, tokenService) // Passing TokenService to InitRoutes

	// --- Inisialisasi dan Konfigurasi gRPC Server ---
	grpcPort := ":50051" // Port terpisah untuk gRPC

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("Gagal mendengarkan di port gRPC %s: %v", grpcPort, err)
	}

	// Daftarkan implementasi gRPC Anda, teruskan TokenService
	s := grpc.NewServer()
	authpb.RegisterAuthServiceServer(s, &authServer{TokenSvc: tokenService})

	log.Printf("Server gRPC AuthService berjalan di %s", grpcPort)

	// Jalankan gRPC server dalam goroutine terpisah
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Gagal menjalankan server gRPC: %v", err)
		}
	}()

	// --- Mulai Server REST API Echo (Blokir goroutine utama) ---
	echoPort := ":1324"
	e.Logger.Fatal(e.Start(echoPort))
}
