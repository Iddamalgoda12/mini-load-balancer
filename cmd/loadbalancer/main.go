package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"mini-load-balancer/internal/balancer"
	"mini-load-balancer/internal/config"
)

func main() {
	cfg, err := config.Parse(os.Args[1:], os.LookupEnv)
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	handler, err := balancer.New(cfg.Backends)
	if err != nil {
		log.Fatalf("load balancer error: %v", err)
	}

	server := &http.Server{
		Addr:              cfg.ListenAddr,
		Handler:           handler,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
	}

	log.Printf("listening on %s; backends: %v", cfg.ListenAddr, cfg.Backends)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}
}
