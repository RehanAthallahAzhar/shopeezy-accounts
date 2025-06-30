package grpc

import (
	"context"
	"log"

	authpb "github.com/RehanAthallahAzhar/shopeezy-protos/proto/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/RehanAthallahAzhar/shopeezy-accounts/helpers"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/services/token"
)

type AuthServer struct {
	authpb.UnimplementedAuthServiceServer
	TokenService token.TokenService
}

// NewAuthServer creates a new AuthServer instance.
func NewAuthServer(tokenService token.TokenService) *AuthServer {
	return &AuthServer{TokenService: tokenService}
}

// ValidateToken mengimplementasikan metode ValidateToken gRPC.
func (s *AuthServer) ValidateToken(ctx context.Context, req *authpb.ValidateTokenRequest) (*authpb.ValidateTokenResponse, error) {
	tokenString := req.GetToken()
	log.Printf("%s[INFO]%s Received gRPC ValidateToken request with token: %s", helpers.ColorYellow, helpers.ColorReset, tokenString)

	//todo: bukannya pakai jwt token service?
	isValid, userID, username, role, errMsg, err := s.TokenService.ValidateToken(ctx, tokenString)
	if err != nil {
		log.Printf("Error validating token in TokenService: %v", err)
		return nil, status.Errorf(codes.Internal, "Internal server error during token validation: %v", err)
	}

	if !isValid {
		return &authpb.ValidateTokenResponse{
			IsValid:      false,
			ErrorMessage: errMsg,
		}, status.Errorf(codes.Unauthenticated, errMsg)
	}

	return &authpb.ValidateTokenResponse{
		IsValid:      true,
		UserId:       userID,
		Username:     username,
		Role:         role,
		ErrorMessage: "",
	}, nil
}
