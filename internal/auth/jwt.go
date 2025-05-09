package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/th1enq/go_coffee/config"
	"google.golang.org/grpc/metadata"
)

type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserID   int64  `json:"user_id"`
	Username string `json:"username,omitempty"`
}

func NewJWTManager(cfg *config.Config, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:     cfg.Auth.JWTSecretKey,
		tokenDuration: tokenDuration,
	}
}

// Generate creates a new JWT token for a user ID
func (m *JWTManager) Generate(userID int64) (string, error) {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: userID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}

// Verify validates the JWT token and returns the claims
func (m *JWTManager) Verify(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.New("unexpected token sign in method")
			}
			return []byte(m.secretKey), nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}

// GetUserIDFromContext extracts the user ID from the context
func GetUserIDFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("no metadata found in context")
	}

	values := md.Get("authorization")
	if len(values) == 0 {
		return "", errors.New("no authorization token found")
	}

	accessToken := values[0]
	// Normally you would parse and verify the token here
	// This is a simplified example
	return accessToken, nil
}
