package main

import (
	sqlc "Los5/db/sqlc"
	handlers "Los5/pkg/handlers"
	logica "Los5/pkg/logica"
	repo "Los5/pkg/repositorios"
	"Los5/views"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/a-h/templ"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {

	// 1. Conectarse a PostgreSQL

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=localhost port=5432 sslmode=disable", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))

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

	// 3. Crear el repositorio común a todas las entidades y casos de uso
	BDrepository := repo.NewBDRepository(queries)

	// 4. Crear los servicios / lógica de negocio

	deptosServ := logica.NewDepartamentoService(BDrepository)
	reservasServ := logica.NewReservaService(BDrepository)
	huespedesServ := logica.NewHuespedService(BDrepository)
	combinadoServ := logica.NewCombinadoService(BDrepository)

	// 5. Crear los handlers web.
	// Los handlers API (*API.go) se conservan compilando, pero no se registran.

	deptosHandler := handlers.NewDepartamentoWebHandler(deptosServ)
	reservasHandler := handlers.NewReservaWebHandler(reservasServ, huespedesServ, combinadoServ)
	huespedesHandler := handlers.NewHuespedWebHandler(huespedesServ, reservasServ)
	combinadoHandler := handlers.NewCombinadoWebHandler(combinadoServ)

	// 6. Configurar las rutas web

	http.Handle("GET /{$}", templ.Handler(views.IndexPage()))

	// Archivos estáticos embebidos (Pico.css y estilos propios)
	http.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(views.StaticFS))))

	// Departamentos
	http.HandleFunc("GET /departamentos", deptosHandler.DepartamentosWeb)
	http.HandleFunc("POST /departamentos", deptosHandler.CrearDepartamentoWeb)
	http.HandleFunc("POST /departamentos/{id}/limpieza", deptosHandler.CambiarHigieneWeb)
	http.HandleFunc("POST /departamentos/{id}/eliminar", deptosHandler.EliminarDepartamentoWeb)

	// Reservas
	http.HandleFunc("GET /reservas", reservasHandler.ReservasWeb)
	http.HandleFunc("POST /reservas", reservasHandler.CrearReservaWeb)
	http.HandleFunc("POST /reservas/{id}/eliminar", reservasHandler.EliminarReservaWeb)

	// Huéspedes
	http.HandleFunc("GET /huespedes", huespedesHandler.HuespedesWeb)
	http.HandleFunc("POST /huespedes", huespedesHandler.CrearHuespedWeb)
	http.HandleFunc("POST /huespedes/{id}/eliminar", huespedesHandler.EliminarHuespedWeb)

	// Casos de uso combinados
	http.HandleFunc("GET /departamentos/{id}/reservas", combinadoHandler.ReservasPorDepartamentoWeb)
	http.HandleFunc("GET /departamentos/{id}/huespedes", combinadoHandler.HuespedesPorDepartamentoWeb)
	http.HandleFunc("GET /departamentos/{id}/disponibilidad", combinadoHandler.DisponibilidadWeb)
	http.HandleFunc("GET /reservas/{id}/total", combinadoHandler.TotalReservaWeb)

	// 7. Iniciar el servidor HTTP

	log.Println("Servidor escuchando en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
