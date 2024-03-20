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
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost" // Default value if DB_HOST is not set
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432" // Default value if DB_PORT is not set
	}

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "postgres" // Default value if DB_NAME is not set
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres" // Default value if DB_USER is not set
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "postgres" // Default value if DB_PASSWORD is not set
	}

	sslmode := os.Getenv("SSL_MODE")
	if sslmode == "" {
		sslmode = "disable" // Default value if SSL_MODE is not set
	}

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
	ipfsAddress := os.Getenv("IPFS_ADDRESS")
	if ipfsAddress == "" {
		ipfsAddress = "/ip4/127.0.0.1/tcp/5001" // Default to localhost if not specified
	}

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
		AllowedOrigins: []string{
			"http://localhost:3000",                       // Allow requests from local development server
			"https://fdf7-213-156-110-145.ngrok-free.app", // Allow requests from the ngrok URL
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "EthereumAddress"},
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
