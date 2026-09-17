// main.go
db, _ := sql.Open("pgx", dsn)
queries := sqlc.New(db)                    // UNA instancia
reservas := NewReservasRepository(queries) // el mismo puntero
huespedes := NewHuespedRepository(queries)
deptos := NewDepartamentosRepository(queries)

Y cada repo:

type ReservasRepository struct {
	query *sqlc.Queries
}

func NewReservasRepository(q *sqlc.Queries) *ReservasRepository {
	return &ReservasRepository{query: q}
}


Ventajas de pasar queries en vez de db:
- Garantizás una sola instancia de *sqlc.Queries compartida (los 3 repos tienen la misma).
- Los repos no saben nada de la conexión: es más desacoplado, más fácil de testear (podés inyectar un Queries armado con un mock o con otra conexión).
- Menos repetición: el sqlc.New se hace una vez, no en cada constructor.
Detalle a tener en cuenta: los repos ya no pueden hacer WithTx (que sale del *sql.DB) a menos que les pases también una forma de obtener la transacción. Si en tu sistema las reservas van a necesitar transacciones (ej. reservar = insertar reserva + actualizar huésped), vas a necesitar el db en algún lado. Para operaciones simples, pasar queries es la opción más limpia.