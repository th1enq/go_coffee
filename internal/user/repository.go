package user

import (
	"context"
	"fmt"
	"strconv"

	"github.com/th1enq/go_coffee/db"
	"github.com/th1enq/go_coffee/internal/model"
)

type UserRepository struct {
	db *db.DB
}

func NewUserRepository(db *db.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// CreateUser creates a new user in the database and returns the user ID
func (r *UserRepository) CreateUser(ctx context.Context, username, password, email string) (string, error) {
	// Create user model
	user := &model.User{
		Username: username,
		Email:    email,
		Password: password,
	}

	// Validate user data
	if err := user.Validate(); err != nil {
		return "", fmt.Errorf("invalid user data: %w", err)
	}

	// Hash the password
	if err := user.HashPassword(); err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	// Insert the user into the database using GORM
	result := r.db.Create(user)
	if result.Error != nil {
		return "", fmt.Errorf("failed to create user: %w", result.Error)
	}

	// Return the user ID as string
	return strconv.FormatUint(uint64(user.ID), 10), nil
}

// GetUserByUsername retrieves a user by username
func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User

	result := r.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find user: %w", result.Error)
	}

	return &user, nil
}

// GetUserByID retrieves a user by ID
func (r *UserRepository) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User

	result := r.db.First(&user, id)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find user: %w", result.Error)
	}

	return &user, nil
}
