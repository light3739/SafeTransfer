package main

import (
	_ "SafeTransfer/docs"
	"SafeTransfer/internal/api"
	"SafeTransfer/internal/db"
	"SafeTransfer/internal/model"
	"SafeTransfer/internal/storage"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
)

const defaultPort = "8083"

// @title SafeTransfer API
// @description This is a sample server for SafeTransfer.
// @version 1.0
// @host localhost:8080
func main() {
	database := setupDatabase()
	defer database.Close()

	ipfsStorage := setupIPFSStorage()
	privateKey := generatePrivateKey()

	apiHandler := api.NewAPIHandler(ipfsStorage, database, privateKey)
	router := setupRouter(apiHandler)

	startServer(router)
}

func setupDatabase() *db.Database {
	dataSourceName := "user=postgres dbname=postgres password=postgres host=localhost sslmode=disable"
	database, err := db.NewDatabase(dataSourceName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = database.AutoMigrate(&model.File{})
	if err != nil {
		log.Fatalf("Failed to migrate schema: %v", err)
	}

	return database
}

func setupIPFSStorage() *storage.IPFSStorage {
	return storage.NewIPFSStorage("/ip4/127.0.0.1/tcp/5001")
}

func generatePrivateKey() *rsa.PrivateKey {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}
	return privateKey
}

func setupRouter(apiHandler *api.Handler) *chi.Mux {
	router := chi.NewRouter()
	router.Mount("/swagger", httpSwagger.WrapHandler)
	apiHandler.RegisterRoutes(router)
	return router
}

func startServer(router *chi.Mux) {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	addr := fmt.Sprintf(":%s", port)
	fmt.Printf("Starting SafeTransfer server on %s...\n", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
