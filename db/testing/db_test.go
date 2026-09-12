package db_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	sqlc "Los5/db/sqlc"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var testQueries *sqlc.Queries

func TestMain(m *testing.M) {
	connStr := "user=user password=password dbname=mydb host=localhost port=5432 sslmode=disable"

	connection, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal("no se pudo conectar a la base de datos:", err)
	}

	if err := connection.Ping(); err != nil {
		log.Fatal("no se pudo hacer ping a la base de datos:", err)
	}

	testQueries = sqlc.New(connection)
	defer connection.Close()

	os.Exit(m.Run())
}
