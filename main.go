package main

import (
	sqlc "Los5/db/sqlc"
	handlers "Los5/pkg/handlers"
	logica "Los5/pkg/logica"
	repo "Los5/pkg/repositorios"
	"database/sql"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {

	// 1. Conectarse a PostgreSQL

	connStr := "user=user password=password dbname=mydb host=localhost port=5432 sslmode=disable" //quitar hardcodeado, poner variables de entorno

	connection, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal("no se pudo conectar a la base de datos:", err)
	}
	defer connection.Close()

	if err := connection.Ping(); err != nil {
		log.Fatal("no se pudo hacer ping a la base de datos:", err)
	}

	// 2. Crear la instancia de sqlc

	queries := sqlc.New(connection)

	// 3. Crear los repositorios

	deptosRepo := repo.NewDepartamentosRepository(queries)
	reservasRepo := repo.NewReservasRepository(queries)
	huespedesRepo := repo.NewHuespedRepository(queries)

	// 4. Crear los servicios / lógica de negocio

	deptosServ := logica.NewDepartamentoService(deptosRepo)
	reservasServ := logica.NewReservaService(reservasRepo)
	huespedesServ := logica.NewHuespedService(huespedesRepo)

	// 5. Crear los handlers

	deptosHandler := handlers.NewDepartamentoHandler(deptosServ)
	reservasHandler := handlers.NewReservaHandler(reservasServ)
	huespedesHandler := handlers.NewHuespedHandler(huespedesServ)

	// 6. Configurar las rutas

	// Departamentos
	http.HandleFunc("GET /departamentos", deptosHandler.GetDepartamentos)
	http.HandleFunc("GET /departamentos/{id}", deptosHandler.GetDepartamento)
	http.HandleFunc("POST /departamentos", deptosHandler.CreateDepartamento)
	http.HandleFunc("PUT /departamentos/{id}", deptosHandler.UpdateDepartamento)
	http.HandleFunc("DELETE /departamentos/{id}", deptosHandler.DeleteDepartamento)

	// Huéspedes
	http.HandleFunc("GET /huespedes", huespedesHandler.GetHuespedes)
	http.HandleFunc("GET /huespedes/{id}", huespedesHandler.GetHuesped)
	http.HandleFunc("POST /huespedes", huespedesHandler.CreateHuesped)
	http.HandleFunc("PUT /huespedes/{id}", huespedesHandler.UpdateHuesped)
	http.HandleFunc("DELETE /huespedes/{id}", huespedesHandler.DeleteHuesped)

	// Reservas
	http.HandleFunc("GET /reservas", reservasHandler.ObtenerReservas)
	http.HandleFunc("GET /reservas/{id}", reservasHandler.ObtenerReserva)
	http.HandleFunc("POST /reservas", reservasHandler.CrearReserva)
	http.HandleFunc("PUT /reservas/{id}", reservasHandler.ActualizarReserva)
	http.HandleFunc("DELETE /reservas/{id}", reservasHandler.EliminarReserva)

	// 7. Iniciar el servidor HTTP

	log.Println("Servidor escuchando en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
