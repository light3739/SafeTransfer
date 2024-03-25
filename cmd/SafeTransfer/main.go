package main

import (
	"SafeTransfer/internal/api"
	"SafeTransfer/internal/config"
	"SafeTransfer/internal/db"
	"SafeTransfer/internal/model"
	"SafeTransfer/internal/repository"
	"SafeTransfer/internal/service"
	"SafeTransfer/internal/storage"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/swaggo/http-swagger"
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
	host := getEnvOrDefault("DB_HOST", "localhost")
	port := getEnvOrDefault("DB_PORT", "5432")
	dbname := getEnvOrDefault("DB_NAME", "postgres")
	user := getEnvOrDefault("DB_USER", "postgres")
	password := getEnvOrDefault("DB_PASSWORD", "postgres")
	sslmode := getEnvOrDefault("SSL_MODE", "disable")

	dataSourceName := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", host, port, dbname, user, password, sslmode)

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
	ipfsAddress := getEnvOrDefault("IPFS_ADDRESS", "/ip4/127.0.0.1/tcp/5001")
	return storage.NewIPFSStorage(ipfsAddress)
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
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "EthereumAddress"},
		ExposedHeaders:   []string{"Link", "X-File-Hash"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler
}

func startServer(router *chi.Mux) {
	port := getEnvOrDefault("PORT", defaultPort)
	addr := fmt.Sprintf(":%s", port)
	fmt.Printf("Starting SafeTransfer server on %s...\n", addr)

	server := &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 30 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	gracefulShutdown(server)
}

func gracefulShutdown(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
