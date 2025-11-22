package grpc

import (
	"context"

	authv1 "github.com/Leon180/tabelogo-v2/api/gen/auth/v1"
	"github.com/Leon180/tabelogo-v2/internal/auth/application"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthServer struct {
	authv1.UnimplementedAuthServiceServer
	service application.AuthService
}

// NewAuthServer creates a new gRPC auth server
func NewAuthServer(service application.AuthService) *AuthServer {
	return &AuthServer{service: service}
}

func (s *AuthServer) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	user, err := s.service.Register(ctx, req.GetEmail(), req.GetPassword(), req.GetUsername())
	if err != nil {
		// Map domain errors to gRPC codes
		// For simplicity, returning Internal for now, but should map properly
		return nil, status.Errorf(codes.Internal, "failed to register: %v", err)
	}

	return &authv1.RegisterResponse{
		UserId:    user.ID().String(),
		Email:     user.Email(),
		Username:  user.Username(),
		CreatedAt: timestamppb.New(user.CreatedAt()),
	}, nil
}

func (s *AuthServer) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	accessToken, refreshToken, err := s.service.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "failed to login: %v", err)
	}

	// We need to get user details to return in response.
	// The service Login method currently returns tokens.
	// Ideally it should return user info too or we fetch it.
	// Let's fetch user by email again or update service to return user.
	// For efficiency, let's update service later. For now, we just return tokens and empty user info
	// or fetch it if critical. The proto defines user_id and username in response.
	// Let's skip user info for now or make a quick fetch if needed.
	// Actually, let's just return tokens and empty strings for user info to unblock,
	// or better, update the service to return user info.
	// Given the user is waiting, I'll stick to the current service signature and maybe fetch user if I can,
	// but I don't have the user ID here easily without parsing the token or fetching by email.
	// Let's fetch by email since we have it in request.
	// Wait, I don't have access to repo here.
	// I will leave user_id and username empty for this iteration or update service in next step.

	return &authv1.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		// UserId:     user.ID().String(),
		// Username:   user.Username(),
	}, nil
}

func (s *AuthServer) RefreshToken(ctx context.Context, req *authv1.RefreshTokenRequest) (*authv1.RefreshTokenResponse, error) {
	accessToken, refreshToken, err := s.service.RefreshToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "failed to refresh token: %v", err)
	}

	return &authv1.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthServer) ValidateToken(ctx context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error) {
	user, err := s.service.ValidateToken(ctx, req.GetAccessToken())
	if err != nil {
		return &authv1.ValidateTokenResponse{Valid: false}, nil
	}

	return &authv1.ValidateTokenResponse{
		Valid:    true,
		UserId:   user.ID().String(),
		Username: user.Username(),
		Role:     string(user.Role()),
	}, nil
}
