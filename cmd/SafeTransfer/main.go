// main.go

package main

import (
	_ "SafeTransfer/docs"
	"SafeTransfer/pkg/api"
	"SafeTransfer/pkg/storage"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/swaggo/http-swagger"
	"log"
	"net/http"
)

// @title SafeTransfer API
// @version 1.0
// @description A secure, decentralized file sharing platform with end-to-end encryption, digital signatures, temporary links, and blockchain audit trails.

func main() {
	// Set up router
	ipfsStorage := storage.NewIPFSStorage("/ip4/127.0.0.1/tcp/5001")
	apiHandler := api.NewAPIHandler(ipfsStorage)

	// Initialize Chi router
	router := chi.NewRouter()

	// Register Swagger documentation route
	router.Mount("/swagger", httpSwagger.WrapHandler)

	apiHandler.RegisterRoutes(router)

	// Start HTTP server
	port := 8083 // Change to your desired port
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Starting SafeTransfer server on %s...\n", addr)

	err := http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
