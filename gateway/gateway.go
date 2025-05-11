package gateway

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/th1enq/go_coffee/config"
	"github.com/th1enq/go_coffee/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func authMiddleware(ctx context.Context, req *http.Request) metadata.MD {
	md := metadata.MD{}

	authHeader := req.Header.Get("Authorization")
	if authHeader != "" {
		md.Set("authorization", authHeader)
		if strings.HasPrefix(authHeader, "Bearer ") {
			token := strings.TrimPrefix(authHeader, "Bearer ")
			md.Set("token", token)
		}
	}

	return md
}

func NewGateWay(cfg *config.Config) (http.Handler, error) {
	gwmux := runtime.NewServeMux(
		runtime.WithMetadata(authMiddleware),
	)

	ctx := context.Background()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	userServiceURL := fmt.Sprintf("%s:%s", "user", cfg.Server.UserPort)

	if err := proto.RegisterUserServiceHandlerFromEndpoint(ctx, gwmux, userServiceURL, opts); err != nil {
		log.Fatalf("failed to register gateway for user service: %v", err)
	}

	characterServiceURL := fmt.Sprintf("%s:%s", "character", cfg.Server.CharacterPort)

	if err := proto.RegisterCharacterServiceHandlerFromEndpoint(ctx, gwmux, characterServiceURL, opts); err != nil {
		log.Fatalf("failed to register gateway for character service: %v", err)
	}

	return gwmux, nil
}
