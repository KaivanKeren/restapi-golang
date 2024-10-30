package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func ConnectDB() {
	var err error
	DB, err = sql.Open("pgx", "postgres://myuser:mypassword@localhost:5432/restapi-golang?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Cek koneksi
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	// Inisialisasi tabel users
	err = createUsersTable()
	if err != nil {
		log.Fatalf("Error creating users table: %v", err)
	}

	fmt.Println("Connected to the database and ensured users table exists!")
}

func createUsersTable() error {
	// SQL untuk membuat tabel users jika belum ada
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(50) NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL
	);`

	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create users table: %v", err)
	}

	return nil
}
