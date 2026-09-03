package main

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	db "Los5.com/ServidorWeb/db/sqlc"
)

var testQueries *db.Queries

func TestMain(m *testing.M) {
	connStr := "user=user password=password dbname=mydb host=localhost port=5432 sslmode=disable"

	connection, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal("no se pudo conectar a la base de datos:", err)
	}

	if err := connection.Ping(); err != nil {
		log.Fatal("no se pudo hacer ping a la base de datos:", err)
	}

	testQueries = db.New(connection)

	os.Exit(m.Run())
}
