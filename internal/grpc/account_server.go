package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/RehanAthallahAzhar/shopeezy-accounts/internal/helpers"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/internal/services"

	accountpb "github.com/RehanAthallahAzhar/shopeezy-protos/pb/account"
)

type AccountServer struct {
	accountpb.UnimplementedAccountServiceServer
	UserService services.UserService
}

// NewAccountServer creates a new AccountServer instance.
func NewAccountServer(userService services.UserService) *AccountServer {
	return &AccountServer{UserService: userService}
}

func (s *AccountServer) GetUser(ctx context.Context, req *accountpb.GetUserRequest) (*accountpb.User, error) {
	userID := req.GetId()
	if userID == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user ID cannot be empty")
	}

	uuid, err := helpers.StringToUUID(userID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID format")
	}

	user, err := s.UserService.GetUserById(ctx, uuid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	// 4. Kemas hasilnya ke dalam format Protobuf dan kirim kembali
	return &accountpb.User{
		Id:          user.Id.String(),
		Name:        user.Name,
		Username:    user.Username,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Address:     user.Addres,
		// Tambahkan field lain jika ada
	}, nil
}
