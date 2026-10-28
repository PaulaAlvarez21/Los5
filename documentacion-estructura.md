# Documentación de estructura de capas

Este proyecto separa la aplicación en capas para respetar la modularización: cada capa tiene una responsabilidad única y se comunica con la siguiente a través de los **structs de dominio** (`pkg/dominio/`), evitando que tipos de sqlc (o de la base de datos) se filtren hacia arriba.

---

## 1. Las capas y cómo se conectan

```
HTTP request
   │
   ▼
┌──────────────────────────┐
│  handlers (net/http)      │  lee/decodifica JSON, responde JSON + códigos HTTP
└──────────────────────────┘
   │  dominio.Huesped / dominio.Departamento / dominio.Reserva
   ▼
┌──────────────────────────┐
│  logica (services)        │  reglas de negocio y validaciones
└──────────────────────────┘
   │  dominio.Huesped / dominio.Departamento / dominio.Reserva
   ▼
┌──────────────────────────┐
│  repositorios             │  convierte dominio ↔ sqlc y ejecuta queries
└──────────────────────────┘
   │  *sqlc.Queries (generado por sqlc)
   ▼
┌──────────────────────────┐
│  db/sqlc  →  base de datos│
└──────────────────────────┘
```

### Cableado en `main.go`

Las capas se arman de abajo hacia arriba, inyectando cada instancia en la siguiente:

1. Se abre la conexión (`sql.Open`) y se crea `sqlc.New(connection)`.
2. Se crean los **repositorios** con esa instancia de queries:
   ```go
   reservasRepo := repo.NewReservasRepository(queries)
   ```
3. Se crean los **services** (lógica) inyectando el repositorio:
   ```go
   reservasServ := logica.NewReservaService(reservasRepo)
   ```
4. Se crean los **handlers** inyectando el service:
   ```go
   reservasHandler := handlers.NewReservaHandler(reservasServ)
   ```
5. Se registran las rutas HTTP apuntando a cada método del handler:
   ```go
   http.HandleFunc("POST /reservas", reservasHandler.CrearReserva)
   ```

Como cada capa guarda la dependencia en un struct y la recibe por el constructor, **nadie instancia dependencias internas**: el handler no sabe que existe sqlc, la lógica no sabe que existe la base, y el repositorio es el único que conoce sqlc.

---

## 2. El tipo que cruza las capas: `dominio`

En `pkg/dominio/` hay un struct por entidad (`Huesped`, `Departamento`, `Reserva`). Esos structs **son el contrato entre capas**:

- **Handlers** decodifican/encodifican JSON directamente en `dominio.X`.
- **Lógica** recibe y devuelve `dominio.X`.
- **Repositorios** reciben y devuelven `dominio.X`.

El repositorio es el encargado de traducir entre el dominio y lo que espera sqlc:

| Campo nullable en dominio | sqlc |
|---|---|
| `*string` (`nil` = ausencia de valor) | `sql.NullString` (`Valid=false` = NULL) |

- Al **escribir**: `toNullString(*string) sql.NullString` convierte `nil → NullString{}` y `&s → {String: s, Valid: true}`.
- Al **leer**: `fromNullString(sql.NullString) *string` convierte `Valid=false → nil`.
- Los queries `Create...`/`Update...` requieren params propios de sqlc, así que el repo construye el `CreateXParams`/`UpdateXParams` campo por campo.
- Las respuestas de `:one`/`:many` vienen como `sqlc.X`; el repo los transforma con `toDominioX(...)`.

Esto logra que, si cambiara sqlc por otra herramienta de acceso a datos, solo habría que tocar la carpeta `pkg/repositorios/`.

---

## 3. Tratamiento de errores

### `pkg/dominio/errores.go`

Define un tipo `AppError` que **lleva el código HTTP incorporado**:

```go
type AppError struct {
    Status int
    Msg    string
}
func (e *AppError) Error() string { return e.Msg }
```

Con dos constructores:
- `NotFoundError()` → `Status: 404`, mensaje genérico "recurso no encontrado".
- `ValidationError(msg)` → `Status: 400`, con el mensaje específico.

La idea es que cada capa cree el error en el lugar donde sabe qué significa, y el handler solo tenga que leer el `Status` para responder. El error "sube" capa por capa sin que ninguna tenga que reinterpretarlo.

### Mapeo en los repositorios

El repo traduce errores de la base a errores de dominio. En `GetX`, un registro que no existe llega como `sql.ErrNoRows`, y se convierte:

```go
if errors.Is(err, sql.ErrNoRows) {
    return dominio.X{}, dominio.NotFoundError()   //404
}
return dominio.X{}, err   //cualquier otro error es ya un error real (500)
```

El `DeleteX` **no** hace este mapeo a propósito: borrar algo que no existe no es un error, simplemente no hace nada y devuelve 204.

### Respuesta en los handlers

En `pkg/handlers/` hay un único helper compartido por todos los handlers:

```go
func responderError(w http.ResponseWriter, err error) {
    var appErr *dominio.AppError
    if errors.As(err, &appErr) {
        http.Error(w, appErr.Msg, appErr.Status)   //usa el Status que trae el error
        return
    }
    http.Error(w, "error interno", http.StatusInternalServerError) //500 si no es AppError
}
```

Todos los handlers terminan `if err != nil { responderError(w, err); return }`. Si el error vino como `AppError` (404 o 400) se responde con su status; si es cualquier otra cosa, se responde 500 sin filtrar detalles internos.

Además, cada handler valida el ID de la URL con `strconv.Atoi(r.PathValue("id"))` y responde 400 si no es numérico.

---

## 4. Validaciones en la capa de lógica

La lógica valida **solo `Crear...` y `Actualizar...`**, porque son los datos que ingresa el usuario y pueden violar reglas de negocio. Los `Obtener...`/`Eliminar...` solo reciben IDs que ya validó el handler, así que no hay reglas que aplicar. Cuando una regla se rompe, devuelve `dominio.ValidationError(...)` y el handler responde 400 automáticamente.

### Reserva (`reservaService.go`)

```go
// Crear
reserva.IDDepto > 0             // la reserva debe referenciar un departamento existente
&& reserva.CantNocches > 0      // al menos 1 noche
&& reserva.FechaInicio.Before(reserva.FechaFin)  // el check-in antes que el check-out

// Actualizar: las mismas + reserva.IDReserva > 0  (debe identificarse el registro a modificar)
```

### Huesped (`huespedService.go`)

```go
// Crear
huesped.IDReserva > 0          // un huésped no puede existir sin una reserva (FK)
&& huesped.Nombre != ""        // campo obligatorio (NOT NULL)
&& huesped.Apellido != ""      // campo obligatorio (NOT NULL)

// Actualizar: las mismas + huesped.IDHuesped > 0
```

### Departamento (`departamentoService.go`)

```go
// Crear
departamento.Nombre != ""      // campo obligatorio (NOT NULL)
&& departamento.Direccion != ""// campo obligatorio (NOT NULL)

// Actualizar: las mismas + departamento.IDDepto > 0
```

Notas:
- Campos opcionales (`*string` en dominio, NULL en la DB): `Descripcion`, `Telefono`, `Email`, `Observaciones`, `Descuento` no se validan porque la ausencia de valor es válida.
- No se valida estado de `bool` (`Disponible`/`Limpio`) porque siempre tienen un valor.
- Las validaciones deben replicarse en crear y actualizar, porque una actualización también puede traer datos inválidos.

---

## 5. Resumen de responsabilidades

| Capa | Responsabilidad | No hace |
|---|---|---|
| `handlers` | Decodificar/encodificar JSON, responder códigos HTTP, propagar `r.Context()` | No conoce la base de datos ni reglas de negocio |
| `logica` | Validaciones y reglas de negocio, pasar `ctx` al repo | No conoce HTTP ni sqlc |
| `repositorios` | Traducir `dominio ↔ sqlc`, ejecutar queries, mapear errores de DB | No conoce HTTP ni reglas de negocio |
| `db/sqlc` | Queries generadas contra la base | Código generado, no se toca a mano |
| `pkg/dominio` | Structs de entidades y errores compartidos | Reglas de negocio (aunque puede pasar DTOs) |

---

## 6. API REST y pruebas

Sobre las capas anteriores, `main.go` registra las rutas con el router de `net/http` (Go 1.22+), que permite declarar el método en el patrón (`"POST /reservas"`) y leer parámetros de ruta con `r.PathValue("id")`.

Cada ruta apunta directo al método del handler correspondiente:

| Método | Ruta | Handler | Service (logica) | Repositorio |
|---|---|---|---|---|
| GET | `/departamentos` | `DepartamentoHandler.ObtenerDepartamentos` | `DepartamentoService.ObtenerDepartamentos` | `DepartamentosRepository.ListDepartamentos` |
| GET | `/departamentos/{id}` | `DepartamentoHandler.ObtenerDepartamento` | `DepartamentoService.ObtenerDepartamento` | `DepartamentosRepository.GetDepartamento` |
| POST | `/departamentos` | `DepartamentoHandler.CrearDepartamento` | `DepartamentoService.CrearDepartamento` | `DepartamentosRepository.CreateDepartamento` |
| PUT | `/departamentos/{id}` | `DepartamentoHandler.ActualizarDepartamento` | `DepartamentoService.ActualizarDepartamento` | `DepartamentosRepository.UpdateDepartamento` |
| DELETE | `/departamentos/{id}` | `DepartamentoHandler.EliminarDepartamento` | `DepartamentoService.EliminarDepartamento` | `DepartamentosRepository.DeleteDepartamento` |
| GET | `/huespedes` | `HuespedHandler.ObtenerHuespedes` | `HuespedService.ObtenerHuespedes` | `HuespedRepository.ListHuespedes` |
| GET | `/huespedes/{id}` | `HuespedHandler.ObtenerHuesped` | `HuespedService.ObtenerHuesped` | `HuespedRepository.GetHuesped` |
| POST | `/huespedes` | `HuespedHandler.CrearHuesped` | `HuespedService.CrearHuesped` | `HuespedRepository.CreateHuesped` |
| PUT | `/huespedes/{id}` | `HuespedHandler.ActualizarHuesped` | `HuespedService.ActualizarHuesped` | `HuespedRepository.UpdateHuesped` |
| DELETE | `/huespedes/{id}` | `HuespedHandler.EliminarHuesped` | `HuespedService.EliminarHuesped` | `HuespedRepository.DeleteHuesped` |
| GET | `/reservas` | `ReservaHandler.ObtenerReservas` | `ReservaService.ObtenerReservas` | `ReservasRepository.ListReservas` |
| GET | `/reservas/{id}` | `ReservaHandler.ObtenerReserva` | `ReservaService.ObtenerReserva` | `ReservasRepository.GetReserva` |
| POST | `/reservas` | `ReservaHandler.CrearReserva` | `ReservaService.CrearReserva` | `ReservasRepository.CreateReserva` |
| PUT | `/reservas/{id}` | `ReservaHandler.ActualizarReserva` | `ReservaService.ActualizarReserva` | `ReservasRepository.UpdateReserva` |
| DELETE | `/reservas/{id}` | `ReservaHandler.EliminarReserva` | `ReservaService.EliminarReserva` | `ReservasRepository.DeleteReserva` |

El recorrido de una request es siempre el mismo:

```
HTTP request
   → handler (decode JSON + PathValue)   → service (validaciones)   → repositorio (dominio ↔ sqlc)   → db/sqlc
   ← handler (encode JSON + status)      ← service (AppError si falla) ← repositorio (mapea sql.ErrNoRows)
```

### Pruebas de endpoints

Los endpoints se prueban con **HURL** desde el archivo `requests.hurl`, que declara cada request con su código HTTP esperado y `[Asserts]` sobre el JSON de respuesta. `make test-api` orquesta todo el ciclo:

1. **Previas**: `.env`, `sqlc generate` + compilación, `docker compose down -v`, `docker compose up -d --wait database` (espera el healthcheck) y levantar la app local.
2. **Prueba**: `hurl --test requests.hurl` recorre los 15 endpoints.
3. **Posteriores**: matar la app y `docker compose down -v`.

`hurl --test` devuelve un exit code distinto de 0 si algún assert falla, por lo que el comando es apto para CI. El detalle de uso está en el `README.md`.

---

## 7. Casos de uso combinados (HTML)

Además del CRUD JSON por entidad, existen cuatro casos de uso que necesitan más de una entidad:

- Reservas de un departamento.
- Huéspedes de un departamento.
- Disponibilidad de un departamento en un rango de fechas.
- Total calculado de una reserva.

Estos casos viven en `CombinadoService` y `CombinadoHandler`. El service valida las reglas de negocio y combina datos del repositorio; el handler no devuelve JSON, sino que renderiza HTML con `a-h/templ` bajo las rutas `/web/...`.

Los DTOs de salida están en `pkg/dominio/combinadoStruct.go` (`ReservasPorDepartamento`, `HuespedesPorDepartamento`, `Disponibilidad`, `TotalReserva`) y también cruzan las capas como el resto de los structs de dominio.

Las queries propias de estos casos están separadas en `db/queries/Combinado.sql` para no mezclarlas con los CRUD de cada entidad.

Reglas principales:

- Un ID inválido se responde como `400`.
- Un departamento o reserva inexistente se responde como `404`.
- La disponibilidad requiere `desde < hasta` y considera tanto el flag `disponible` del departamento como las reservas solapadas.
- El total interpreta `descuento` como porcentaje entre `0` y `100`.