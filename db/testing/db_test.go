package db_test

/* db_test aplica pruebas de caja negra (External Testing):
Evalúa el paquete desde la perspectiva de un usuario externo.
Previene el acoplamiento a detalles de implementación privada.
Propio de go, opción de testeo */

import (
	"database/sql" //Interfaz estándar de Go para bases de datos SQL
	"log"
	"os"
	"testing" //Paquete de go para pruebas

	sqlc "Los5/db/sqlc"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var testQueries *sqlc.Queries

//Variable que almacena la instancia de consultas de sqlc para que cualquier otro test del paquete pueda usarla.

func TestMain(m *testing.M) {
	connStr := "user=user password=password dbname=mydb host=localhost port=5432 sslmode=disable"

	//conexion con la base de datos
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
