package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func ConnectDB() {
    var err error
    
    // Get environment variables with fallbacks
    dbHost := getEnv("DB_HOST", "localhost")
    dbUser := getEnv("DB_USER", "myuser")
    dbPassword := getEnv("DB_PASSWORD", "mypassword")
    dbName := getEnv("DB_NAME", "restapi-golang")
    dbPort := getEnv("DB_PORT", "5432")

    // Build connection string
    dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", 
        dbUser, dbPassword, dbHost, dbPort, dbName)
    
    // Connect with retry logic
    for i := 0; i < 5; i++ {
        DB, err = sql.Open("pgx", dsn)
        if err != nil {
            log.Printf("Failed to open DB: %v", err)
            time.Sleep(time.Second * 5)
            continue
        }

        if err = DB.Ping(); err != nil {
            log.Printf("Failed to ping DB: %v", err)
            time.Sleep(time.Second * 5)
            continue
        }

        fmt.Println("Connected to the database")
        return
    }
    
    panic(fmt.Sprintf("Failed to connect to database after 5 attempts: %v", err))
}

// Helper function to get environment variables
func getEnv(key, fallback string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return fallback
}