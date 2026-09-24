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