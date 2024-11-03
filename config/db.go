package config

import (
    "database/sql"
    "fmt"
    _ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func ConnectDB() {
    var err error
    dsn := "postgres://myuser:mypassword@localhost:5432/restapi-golang"
    DB, err = sql.Open("pgx", dsn)
    if err != nil {
        panic(err)
    }

    if err = DB.Ping(); err != nil {
        panic(err)
    }

    fmt.Println("Connected to the database")
}