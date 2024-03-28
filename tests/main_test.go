package tests

import (
	"SafeTransfer/internal/db"
	"SafeTransfer/internal/model"
	"SafeTransfer/utils"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabaseMigration(t *testing.T) {
	// Set up a test database connection
	testDB := setupTestDatabase(t)
	defer testDB.Close()

	// Check if the migration was successful
	assertTableExists(t, testDB, "files")
	assertTableExists(t, testDB, "users")
}

func setupTestDatabase(t *testing.T) *db.Database {
	// Set up a test database connection to PostgreSQL
	host := utils.GetEnvOrDefault("DB_HOST", "localhost")
	port := utils.GetEnvOrDefault("DB_PORT", "5432")
	dbname := utils.GetEnvOrDefault("DB_NAME", "postgres")
	user := utils.GetEnvOrDefault("DB_USER", "postgres")
	sslmode := utils.GetEnvOrDefault("SSL_MODE", "disable")
	password, err := os.ReadFile("/run/secrets/db_password")
	if err != nil {
		log.Fatalf("Failed to read database password: %v", err)
	}

	dataSourceName := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", host, port, dbname, user, password, sslmode)

	testDB, err := db.NewDatabase(dataSourceName)
	require.NoError(t, err, "failed to create test database")

	err = testDB.AutoMigrate(&model.File{}, &model.User{})
	require.NoError(t, err, "failed to migrate test database")

	return testDB
}

func assertTableExists(t *testing.T, db *db.Database, tableName string) {
	// Check if the specified table exists in the database
	exists := db.Migrator().HasTable(tableName)
	assert.True(t, exists, "table '%s' does not exist", tableName)
}
