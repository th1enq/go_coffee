package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/th1enq/go_coffee/config"
	"github.com/th1enq/go_coffee/gateway"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to loading config: %v", err)
	}
	gw, err := gateway.NewGateWay(cfg)
	if err != nil {
		log.Fatalf("failed to create gateway: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/api/v1/", gw) 
	
	fmt.Printf("Starting HTTP Server on: %s\n", cfg.Server.GatewayPort)
	if err := http.ListenAndServe(":"+cfg.Server.GatewayPort, mux); err != nil {
		log.Fatalf("failed to serve HTTP Server: %v", err)
	}
}
