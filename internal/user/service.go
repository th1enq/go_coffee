package user

import (
	"context"
	"strconv"

	"github.com/th1enq/go_coffee/db"
	"github.com/th1enq/go_coffee/internal/auth"
	"github.com/th1enq/go_coffee/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	proto.UnimplementedUserServiceServer
	repo       *UserRepository
	jwtManager *auth.JWTManager
}

func NewUserService(db *db.DB, jwtManager *auth.JWTManager) *UserService {
	return &UserService{
		repo:       NewUserRepository(db),
		jwtManager: jwtManager,
	}
}

func (s *UserService) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	// Create user in the repository
	userID, err := s.repo.CreateUser(ctx, req.Username, req.Password, req.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	// Convert string ID to int64 for JWT generation
	userIDInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "invalid user ID format: %v", err)
	}

	// Generate JWT token for the new user
	token, err := s.jwtManager.Generate(userIDInt)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token: %v", err)
	}

	return &proto.RegisterResponse{
		UserId: userID,
		Token:  token,
	}, nil
}

func (s *UserService) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	// Get user by username
	user, err := s.repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	// Check password
	if !user.CheckPassword(req.Password) {
		return nil, status.Errorf(codes.Unauthenticated, "incorrect password")
	}

	// Generate JWT token
	userID := int64(user.ID)
	token, err := s.jwtManager.Generate(userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token: %v", err)
	}

	return &proto.LoginResponse{
		Token: token,
		User: &proto.UserInfo{
			Id:        int64(user.ID),
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.String(),
			UpdatedAt: user.UpdatedAt.String(),
		},
	}, nil
}

func (s *UserService) GetUserInfo(ctx context.Context, req *proto.GetUserInfoRequest) (*proto.GetUserInfoResponse, error) {
	// Extract user ID from context
	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated: %v", err)
	}

	// Convert string ID to uint
	id, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "invalid user ID: %v", err)
	}

	// Get user by ID
	user, err := s.repo.GetUserByID(ctx, uint(id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	return &proto.GetUserInfoResponse{
		User: &proto.UserInfo{
			Id:        int64(user.ID),
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.String(),
			UpdatedAt: user.UpdatedAt.String(),
		},
	}, nil
}

func (s *UserService) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	// Create user in the repository
	userID, err := s.repo.CreateUser(ctx, req.Username, req.Password, req.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	// Convert string ID to int64 for JWT generation
	userIDInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "invalid user ID format: %v", err)
	}

	// Generate JWT token for the new user
	token, err := s.jwtManager.Generate(userIDInt)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token: %v", err)
	}

	return &proto.CreateUserResponse{
		UserId: userID,
		Token:  token,
	}, nil
}
