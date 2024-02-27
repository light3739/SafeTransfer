// main.go

package main

import (
	_ "SafeTransfer/docs"
	"SafeTransfer/pkg/api"
	"SafeTransfer/pkg/db"
	"SafeTransfer/pkg/model"
	"SafeTransfer/pkg/storage"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/swaggo/http-swagger"
	"log"
	"net/http"
)

func main() {
	dataSourceName := "user=test2 dbname=test2 password=test2 host=localhost sslmode=disable"
	database, err := db.NewDatabase(dataSourceName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Migrate the schema
	err = database.AutoMigrate(&model.File{}) // Corrected line
	if err != nil {
		log.Fatalf("Failed to migrate schema: %v", err)
	}

	// Set up router
	ipfsStorage := storage.NewIPFSStorage("/ip4/127.0.0.1/tcp/5001")
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}
	apiHandler := api.NewAPIHandler(ipfsStorage, database, privateKey)

	// Initialize Chi router
	router := chi.NewRouter()

	// Register Swagger documentation route
	router.Mount("/swagger", httpSwagger.WrapHandler)

	apiHandler.RegisterRoutes(router)

	// Start HTTP server
	port := 8083 // Change to your desired port
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Starting SafeTransfer server on %s...\n", addr)

	err = http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
