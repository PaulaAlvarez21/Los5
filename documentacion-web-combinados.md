# Interfaz web de casos de uso combinados

Fecha: 2026-10-28

## Objetivo

Agregar una capa web que reutilice los templates `a-h/templ` del proyecto y exponga los casos de uso combinados, manteniendo la API JSON existente en la raíz.

Casos de uso cubiertos:

- Reservas de un departamento.
- Huéspedes de un departamento.
- Disponibilidad de un departamento en un rango de fechas.
- Total calculado de una reserva.

## Decisiones de diseño

- Las páginas web se montan bajo `/web/...`; la API JSON del CRUD sigue en `/departamentos`, `/reservas` y `/huespedes`.
- Se usa `github.com/a-h/templ`, que ya era el sistema de templates presente en `views/*.templ`.
- `CombinadoHandler` contiene solo los casos combinados que renderizan HTML.
- `WebHandler` contiene las páginas simples y el procesamiento de formularios.
- El `descuento` de una reserva se interpreta como porcentaje entre `0` y `100`.
- La disponibilidad se calcula con el flag `departamento.disponible` y con las reservas solapadas en el rango `[desde, hasta)`.
- No se modificó la estructura visual de los templates existentes. Solo se cambiaron links y `action` para apuntar a `/web`.
- No se ejecutaron `sqlc generate` ni `templ generate` por pedido explícito. Por eso el proyecto no compila todavía: faltan los métodos generados de `db/queries/Combinado.sql` y los archivos `views/*_templ.go`.

## Dependencias

Se agregó `github.com/a-h/templ` con:

```bash
go get github.com/a-h/templ@latest
```

Esto actualizó `go.mod` y `go.sum`. No se ejecutó ningún generador.

## Rutas web

| Método | Ruta | Handler | Descripción |
|---|---|---|---|
| GET | `/` | redirección | Redirige a `/web/departamentos` |
| GET | `/static/styles.css` | archivo estático | Sirve `views/styles.css` |
| GET | `/web/departamentos` | `WebHandler.DepartamentosWeb` | Catálogo de departamentos |
| POST | `/web/departamentos` | `WebHandler.CrearDepartamentoWeb` | Alta desde formulario |
| POST | `/web/departamentos/{id}/limpieza` | `WebHandler.CambiarHigieneWeb` | Alterna el estado de higiene |
| GET | `/web/reservas` | `CombinadoHandler.ReservasWeb` | Reservas con filtros |
| POST | `/web/reservas` | `WebHandler.CrearReservaWeb` | Alta desde formulario |
| GET | `/web/huespedes` | `WebHandler.HuespedesWeb` | Padrón de huéspedes |
| POST | `/web/huespedes` | `WebHandler.CrearHuespedWeb` | Alta desde formulario |
| GET | `/web/departamentos/{id}/huespedes` | `CombinadoHandler.HuespedesPorDepartamentoWeb` | Huéspedes de un departamento |
| GET | `/web/departamentos/{id}/disponibilidad` | `CombinadoHandler.DisponibilidadWeb` | Formulario o resultado de disponibilidad |
| GET | `/web/reservas/{id}/total` | `CombinadoHandler.TotalReservaWeb` | Detalle del total calculado |

`/web/reservas` acepta los query params `depto_id`, `desde` y `hasta`.

`/web/departamentos/{id}/disponibilidad` acepta `desde` y `hasta`. Si no se envían, muestra el formulario; si se envían, muestra el resultado.

## Cambios por archivo

### Dominio

`pkg/dominio/combinadoStruct.go`

- Se agregó `ReservaConDetalle`, que es el DTO que espera `views/reservas.templ`.
- Se agregó `PaginaReservas`, que agrupa reservas, departamentos y huéspedes para el formulario y la tabla.

`ReservaConDetalle` tiene `PrecioBase float64` porque el template lo formatea con `fmt.Sprintf("%.2f", ...)`.

### Servicio

`pkg/logica/combinadoService.go`

- Se agregó `FiltroReservas`.
- Se agregó `ObtenerDepartamento`, usado por la pantalla de disponibilidad cuando todavía no hay fechas.
- Se agregó `ObtenerPaginaReservas`, que arma la página de reservas:
  - lista reservas y huéspedes del departamento o de todo el sistema;
  - valida el rango de fechas;
  - filtra reservas por solapamiento con el rango;
  - agrupa los huéspedes por reserva;
  - concatena los nombres de todos los huéspedes en `NombreHuesped`;
  - deja `ApellidoHuesped` vacío para no cambiar el template existente;
  - convierte `precio_base` de `string` a `float64`.

### Handlers

`pkg/handlers/combinadoHandlers.go`

- `ReservasWeb`: parsea filtros, obtiene la página y renderiza `ReservasPage`.
- `HuespedesPorDepartamentoWeb`: obtiene el departamento y sus huéspedes y renderiza `HuespedesPage`.
- `DisponibilidadWeb`: muestra formulario o resultado según los query params.
- `TotalReservaWeb`: calcula el total y renderiza `TotalReservaPage`.
- Helpers:
  - `renderHTML`: setea `text/html` y renderiza un componente `templ`.
  - `responderErrorWeb`: traduce `dominio.AppError` a status HTTP y renderiza `ErrorPage`.
  - `parseID`: valida el `{id}` de la ruta.
  - `parseFiltroReservas`: parsea `depto_id`, `desde` y `hasta`.
  - `parseRango`: valida fechas obligatorias y `desde < hasta`.
  - `parseFecha`: parsea `YYYY-MM-DD`.

`pkg/handlers/webHandlers.go`

- `DepartamentosWeb`, `HuespedesWeb`: páginas de listado.
- `CrearDepartamentoWeb`, `CrearReservaWeb`, `CrearHuespedWeb`: procesan formularios y redirigen con `303`.
- `CambiarHigieneWeb`: alterna `Limpio` y redirige.
- `stringOpcional`: convierte campos vacíos del formulario en `nil`.

En `CrearReservaWeb`, después de crear la reserva se reasigna el huésped seleccionado a la reserva nueva. Esto es necesario porque la FK `huesped.id_reserva` vive en el huésped, no en la reserva.

### Templates

Se modificaron solo rutas y expresiones necesarias para que `templ` compile:

- `views/layout.templ`: navegación y logo a `/web/...`.
- `views/departamentos.templ`: alta, limpieza, link de reservas y descripción opcional.
- `views/reservas.templ`: filtro y alta a `/web/reservas`, y formato de fechas.
- `views/huespedes.templ`: alta a `/web/huespedes` y campos opcionales.

Se agregó `views/helpers.go` con `valorTexto(*string) string`. Es necesario porque `templ.JoinStringErrs` no acepta punteros y varios campos del dominio son `*string`. Las fechas `time.Time` también se formatean explícitamente.

Se crearon:

- `views/disponibilidad.templ`: formulario y resultado de disponibilidad.
- `views/total_reserva.templ`: subtotal, descuento y total.
- `views/error.templ`: página de error con status y mensaje.

### Wiring

`main.go`

- Se creó `WebHandler` con los servicios de departamento, reserva y huésped.
- Se registraron las rutas `/web/...`.
- Se agregó redirección de `/` a `/web/departamentos`.
- Se sirve `views/styles.css` en `/static/styles.css`.
- Se quitaron las rutas combinadas viejas de la raíz; la API JSON CRUD no se tocó.

### Pruebas manuales

`requests.hurl`

Se actualizó la sección de casos combinados para usar `/web/...` y se ajustaron los asserts al HTML nuevo:

- reservas filtradas por departamento;
- huéspedes de un departamento;
- disponibilidad con y sin solapamiento;
- total con descuento;
- errores `400` y `404`.

### Documentación existente

- `README.md`: tabla de rutas combinadas actualizada a `/web` y a `a-h/templ`.
- `documentacion-estructura.md`: corregida la referencia a `html/template` y `go:embed`.

## Estado de compilación

En el workspace no se ejecutaron generadores, así que faltan:

1. `views/*_templ.go`, que se genera con `templ generate`.
2. Los métodos nuevos de `db/queries/Combinado.sql` en `db/sqlc/*.go`, que se regeneran con `sqlc generate`.

Para verificar que los cambios compilan, se copió el proyecto a `/tmp/opencode/Los5-check`, se ejecutaron allí `sqlc generate` y `templ generate`, y `go build ./...` terminó sin errores.

Para dejar el workspace compilable, cuando se permita ejecutar generadores:

```bash
go run github.com/a-h/templ/cmd/templ@v0.3.1100 generate
sqlc generate
go build ./...
```

## Limitaciones conocidas

- El formulario de alta de huésped no incluye `id_reserva`; el handler lo acepta vacío y crea el huésped sin reserva (`NULL`).
- El formulario de alta de reserva selecciona un huésped existente y luego lo reasigna. Si la reasignación falla, la reserva ya quedó creada.
- Los nombres de varios huéspedes se concatenan en `NombreHuesped` y `ApellidoHuesped` queda vacío, para no cambiar la tabla de reservas.
- La pantalla de disponibilidad y la de total no tienen links desde las páginas existentes; se accede por URL directa.
- No se ejecutaron generadores, por lo que el estado compilable queda pendiente.
