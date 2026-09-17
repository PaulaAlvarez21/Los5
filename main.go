package main

import (
	sqlc "Los5/db/sqlc"
	"database/sql" //Interfaz estándar de Go para bases de datos SQL
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var queries *sqlc.Queries

func main() {

	// 1. Conectarse a PostgreSQL

	connStr := "user=user password=password dbname=mydb host=localhost port=5432 sslmode=disable"

	connection, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal("no se pudo conectar a la base de datos:", err)
	}

	if err := connection.Ping(); err != nil {
		log.Fatal("no se pudo hacer ping a la base de datos:", err)
	}

	defer connection.Close()

	// 2. Crear la instancia de sqlc
	queries = sqlc.New(connection) //pasar a el repository (tienen que tener el puntero a las queries)

	// 3. Crear/configurar los handlers

	// 4. Configurar las rutas

	// 5. Iniciar el servidor HTTP

}
