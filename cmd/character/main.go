package main

import (
	"log"
	"net"

	"github.com/th1enq/go_coffee/config"
	"github.com/th1enq/go_coffee/db"
	"github.com/th1enq/go_coffee/internal/character"
	"github.com/th1enq/go_coffee/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Character Service: failed loading config: %v", err)
	}

	// Set service type for character service
	db, err := db.LoadDB(cfg, config.CharacterService)
	if err != nil {
		log.Fatalf("Character Service: failed connect db: %v", err)
	}

	server := grpc.NewServer()
	service := character.NewCharacterService(db)

	proto.RegisterCharacterServiceServer(server, service)

	reflection.Register(server)

	lis, err := net.Listen("tcp", ":"+cfg.Server.CharacterPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Starting character service on: %s", cfg.Server.CharacterPort)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
