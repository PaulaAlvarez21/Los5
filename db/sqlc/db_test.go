package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	connStr := "user=user password=password dbname=mydb host=localhost port=5432 sslmode=disable"

	connection, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal("no se pudo conectar a la base de datos:", err)
	}

	if err := connection.Ping(); err != nil {
		log.Fatal("no se pudo hacer ping a la base de datos:", err)
	}

	testQueries = New(connection)
	defer connection.Close()

	os.Exit(m.Run())
}
