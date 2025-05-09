package user

import (
	"context"

	"github.com/th1enq/go_coffee/db"
	"github.com/th1enq/go_coffee/internal/auth"
	gen "github.com/th1enq/go_coffee/proto"
)

type UserService struct {
	gen.UnimplementedUserServiceServer
	repo       *UserRepository
	jwtManager *auth.JWTManager
}

func NewUserService(db *db.DB, jwtManager *auth.JWTManager) *UserService {
	return &UserService{
		repo:       NewUserRepository(db),
		jwtManager: jwtManager,
	}
}

func (s *UserService) Register(ctx context.Context, req *gen.RegisterRequest) (*gen.RegisterResponse, error) {
	return nil, nil
}

func (s *UserService) Login(ctx context.Context, req *gen.LoginRequest) (*gen.LoginResponse, error) {
	return nil, nil
}
