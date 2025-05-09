package main

import (
	"log"
	"net"
	"time"

	"github.com/th1enq/go_coffee/config"
	"github.com/th1enq/go_coffee/db"
	"github.com/th1enq/go_coffee/internal/auth"
	"github.com/th1enq/go_coffee/internal/user"
	gen "github.com/th1enq/go_coffee/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("User Service: failed loading config: %v", err)
	}

	db, err := db.LoadDB(cfg)
	if err != nil {
		log.Fatalf("User Service: failed connect db: %v", err)
	}

	jwtManager := auth.NewJWTManager(cfg, 24*time.Hour)

	server := grpc.NewServer()
	service := user.NewUserService(db, jwtManager)

	gen.RegisterUserServiceServer(server, service)

	reflection.Register(server)

	lis, err := net.Listen("tcp", ":"+cfg.Server.UserPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Starting user service on: %s", cfg.Server.UserPort)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
