package main

import (
	_ "SafeTransfer/docs"
	"SafeTransfer/internal/api"
	"SafeTransfer/internal/config"
	"SafeTransfer/internal/db"
	"SafeTransfer/internal/model"
	"SafeTransfer/internal/repository"
	"SafeTransfer/internal/service"
	"SafeTransfer/internal/storage"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
)

const defaultPort = "8083"

// @title SafeTransfer API
// @description This is a sample server for SafeTransfer.
// @version 1.0
// @host localhost:8083
func main() {
	cfg := config.LoadConfig()
	database := setupDatabase()
	defer database.Close()

	ipfsStorage := setupIPFSStorage()

	fileRepo := repository.NewFileRepository(database)

	fileService := service.NewFileService(ipfsStorage, fileRepo)
	downloadService := service.NewDownloadService(ipfsStorage, fileRepo)
	userRepo := repository.NewUserRepository(database)
	userService := service.NewUserService(userRepo, cfg.JWTSecretKey)

	apiHandler := api.NewAPIHandler(fileService, downloadService, userService)
	router := setupRouter(apiHandler)

	startServer(router)
}

func setupDatabase() *db.Database {
	dataSourceName := "user=postgres dbname=postgres password=postgres host=localhost sslmode=disable"
	database, err := db.NewDatabase(dataSourceName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = database.AutoMigrate(&model.File{}, &model.User{})
	if err != nil {
		log.Fatalf("Failed to migrate schema: %v", err)
	}

	return database
}

func setupIPFSStorage() *storage.IPFSStorage {
	return storage.NewIPFSStorage("/ip4/127.0.0.1/tcp/5001")
}

func setupRouter(apiHandler *api.Handler) *chi.Mux {
	router := chi.NewRouter()
	router.Use(corsHandler())
	router.Mount("/swagger", httpSwagger.WrapHandler)
	apiHandler.RegisterRoutes(router)
	return router
}

func corsHandler() func(http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler
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
