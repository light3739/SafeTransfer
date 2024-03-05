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
	"os"
)

const defaultPort = "8083"

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
	dataSourceName := "user=test2 dbname=test2 password=test2 host=localhost sslmode=disable"
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
