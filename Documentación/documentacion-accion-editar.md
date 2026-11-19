# Acción de editar en las listas de entidades

Fecha: 2026-11-19

## 1. Objetivo

Agregar la opción de **editar** cada registro desde las tablas de
`/departamentos`, `/huespedes` y `/reservas`, manteniendo la interfaz
**server-rendered con templ y HTMX**, sin JavaScript propio.

La idea es reutilizar exactamente el modal que ya se usa para el alta, pero con
los datos actuales del registro, y reutilizar también la capa de lógica y de
persistencia, que ya tenía el update implementado.

## 2. Lo que ya existía (no se tocó)

Antes de este cambio, la actualización ya estaba resuelta de punta a punta salvo
por la interfaz web:

| Capa | Método existente |
|---|---|
| SQL (sqlc) | `UpdateDepartamento`, `UpdateHuesped`, `UpdateReserva` |
| Repositorio | `BDRepository.UpdateDepartamento`, `UpdateHuesped`, `UpdateReserva` |
| Lógica | `DepartamentoService.ActualizarDepartamento`, `HuespedService.ActualizarHuesped`, `ReservaService.ActualizarReserva` |
| API | `DepartamentoAPIHandler.ActualizarDepartamento`, etc. (no registrados en `main.go`) |

Tampoco se modificaron:

- `db/` (queries, schema, migraciones y código generado por sqlc).
- `pkg/repositorios/`.
- `pkg/logica/`.
- Los handlers API (`*HandlersAPI.go`).
- Los flujos de alta, borrado, cambio de higiene, filtros y vistas combinadas.
- `pkg/handlers/helpers.go` y `web.hurl`.

## 3. Decisiones de diseño

### 3.1. Un botón por fila que abre el modal

Cada fila tiene un botón `Editar` con:

```html
<button type="button" class="btn btn-secondary"
	hx-get="/departamentos/7/editar"
	hx-target="#modal-root"
	hx-swap="innerHTML">
	Editar
</button>
```

Es el mismo mecanismo que el botón `+ Nuevo`, pero con el ID del registro en la
URL. El handler `GET` devuelve el modal con el formulario precargado.

### 3.2. El ID siempre sale de la URL

El ID **nunca** se toma del formulario. Se lee con `parseID(r)`, que usa
`r.PathValue("id")`. El handler arma el struct de dominio con ese ID y llama al
service correspondiente. Así se reutilizan las validaciones que ya existen:

- Departamento: ID, nombre y dirección obligatorios.
- Huésped: ID, nombre y teléfono obligatorios (el apellido es opcional desde el
  commit del atajo para crear huéspedes).
- Reserva: ID, departamento, cantidad de noches y rango de fechas válidos.

### 3.3. Respuesta sin recargar la página

Los formularios de edición apuntan a la fila concreta:

```html
<form action="/departamentos/7/editar" method="post"
	hx-post="/departamentos/7/editar"
	hx-target="#departamento-7"
	hx-swap="outerHTML">
```

Por eso a las filas se les agregó un ID:

- `departamento-{id}`
- `reserva-{id}`
- `huesped-{id}` (tanto en `HuespedRow` como en `HuespedRowOOB`, para que las
  filas insertadas por HTMX también se puedan editar sin recargar).

El handler responde un componente `...Editado`, que es la fila ya modificada más
un swap OOB que vacía `#modal-root`:

```templ
templ DepartamentoEditado(d dominio.Departamento) {
	@DepartamentoRow(d)
	<div id="modal-root" hx-swap-oob="innerHTML"></div>
}
```

HTMX hace `outerHTML` sobre la fila y, en paralelo, cierra el modal. No
hace falta recargar ni reconsultar el listado completo excepto en reservas, donde
la fila necesita datos combinados.

### 3.4. Errores simples

No se volvió a renderizar el modal con el error adentro. En su lugar, cualquier
fallo (ID inválido, `ParseForm`, validación del service) se maneja con lo que ya
existía:

```go
if esHTMX(r) {
	responderErrorHTMX(w, r, err, "#page-errors")
	return
}
h.renderDepartamentos(w, r, statusDeError(err), mensajeDeError(err))
```

El modal queda abierto y el mensaje aparece en el banner de `#page-errors`. Se
cierra con el botón `×`. Es menos elegante, pero evita reconsultar datos y
duplicar plantillas de error. Para errores de formulario se usa
`dominio.ValidationError("mensaje")`.

En los handlers esto se implementó con una closure local `fallo` que concentra
ese bloque y se llama desde cada validación.

## 4. Rutas agregadas

En `main.go`:

| Método | Ruta | Handler |
|---|---|---|
| `GET` | `/departamentos/{id}/editar` | `EditarDepartamentoWeb` |
| `POST` | `/departamentos/{id}/editar` | `ActualizarDepartamentoWeb` |
| `GET` | `/reservas/{id}/editar` | `EditarReservaWeb` |
| `POST` | `/reservas/{id}/editar` | `ActualizarReservaWeb` |
| `GET` | `/huespedes/{id}/editar` | `EditarHuespedWeb` |
| `POST` | `/huespedes/{id}/editar` | `ActualizarHuespedWeb` |

No hay conflicto con rutas existentes porque Go 1.22+ enruta también por método
HTTP.

## 5. Cambios en los templates

### `views/entity_list.templ`

- ID en `DepartamentoRow`, `ReservaRow`, `HuespedRow` y `HuespedRowOOB`.
- Botón `Editar` en esos mismos cuatro componentes.
- Componentes nuevos `DepartamentoEditado`, `ReservaEditado` y `HuespedEditado`,
  que devuelven la fila modificada + el OOB que vacía `#modal-root`.
- No se tocó `ReservaSimpleList`, igual que se hizo con el borrado.

### `views/entity_form.templ`

Se agregaron tres formularios hermanos de los de alta:

- `DepartamentoEditarForm(d)`: campos con `value`, textarea con `valorTexto` y
  los checkboxes con `if` para poner `checked` cuando corresponde.
- `HuespedEditarForm(h, reservas)`: reserva actual preseleccionada, teléfono
  obligatorio y apellido opcional.
- `ReservaEditarForm(r, primerHuesped, deptos, huespedes)`: departamento y
  huésped preseleccionados, fechas en formato `2006-01-02`,
  `@CampoCantNoches(fmt.Sprint(r.CantNoches))` readonly, precio base, descuento y
  observaciones. No incluye el atajo `+ Crear huésped`, porque reemplazaría el
  modal que se está editando.

### `views/modal.templ`

Tres wrappers nuevos, al lado de los de alta:

```templ
templ DepartamentoEditarModal(d dominio.Departamento) {
	@Modal("Editar departamento", "") {
		@DepartamentoEditarForm(d)
	}
}
```

Y los equivalentes `HuespedEditarModal` y `ReservaEditarModal`. El segundo
parámetro de `Modal` es `""` porque los errores van a `#page-errors`.

Un *wrapper* es simplemente una plantilla pequeña que combina el modal genérico
con un formulario concreto. Existe porque el handler necesita devolver un único
componente templ y desde Go es incómodo armar los hijos de `@Modal`.

## 6. Handlers web agregados

### Departamentos

`pkg/handlers/departamentoHandlersWEB.go`:

- `EditarDepartamentoWeb`: `parseID` → `ObtenerDepartamento` → modal.
- `ActualizarDepartamentoWeb`: arma `dominio.Departamento` con el ID de la URL,
  llama a `ActualizarDepartamento` y responde `DepartamentoEditado` (HTMX) o
  redirige `303` al listado.

### Huéspedes

`pkg/handlers/huespedHandlersWEB.go`:

- `EditarHuespedWeb`: `parseID` → `ObtenerHuesped` → `ObtenerReservas` → modal.
- `ActualizarHuespedWeb`: parsea el `id_reserva` opcional igual que en el alta,
  exige teléfono no vacío, llama a `ActualizarHuesped` y responde
  `HuespedEditado`.

### Reservas

`pkg/handlers/reservaHandlersWEB.go`:

- Helper privado `primerHuespedDeReserva(huespedes, idReserva)`: devuelve el
  huésped asociado con **menor `IDHuesped`** (el primero cargado, porque
  `id_huesped` es `serial`). No alcanza con tomar el primero del slice, porque
  `ListHuespedes` ordena por apellido y nombre.
- `EditarReservaWeb`: `parseID` → `ObtenerReserva` → `ObtenerPaginaReservas` →
  calcula el primer huésped → modal.
- `ActualizarReservaWeb`:
  1. Valida departamento y huésped seleccionado.
  2. Calcula las noches con `fechaFin.Sub(fechaInicio).Hours() / 24`, igual que
     en el alta; el campo es `readonly` y se recalcula al cambiar fechas mediante
     `GET /reservas/noches`.
  3. Llama a `ActualizarReserva`.
  4. Lee el hidden `id_huesped_actual`. Si es distinto del seleccionado y sigue
     asociado a esta reserva, desasigna **solo ese** huésped
     (`IDReserva = nil`). Los demás huéspedes de la reserva no se tocan.
  5. Asigna el huésped elegido a la reserva.
  6. Con HTMX reconsulta la página con `ObtenerPaginaReservas`, usa
     `buscarReservaDetalle(id)` (helper existente) y responde
     `ReservaEditado`. Sin HTMX redirige `303` a `/reservas`.

## 7. Flujo completo

### Editar un departamento

1. El usuario hace clic en `Editar` en una fila.
2. HTMX hace `GET /departamentos/7/editar` y mete el modal en `#modal-root`.
3. El usuario cambia los datos y envía.
4. HTMX hace `POST /departamentos/7/editar`.
5. El service valida y actualiza.
6. La respuesta es la fila nueva + OOB que vacía `#modal-root`.
7. HTMX reemplaza `#departamento-7` y cierra el modal. Sin recarga.

### Editar una reserva

1. El modal muestra preseleccionados el departamento y el **primer** huésped de
   la reserva (menor ID).
2. Si el usuario elige otro huésped, el hidden `id_huesped_actual` permite
   desasignar solo al anterior y asignar al nuevo.
3. La fila se reemplaza con los datos combinados actualizados
   (`ReservaConDetalle`).

### Sin HTMX

Los formularios tienen `action` y `method="post"`, así que funcionan también sin
JavaScript: el handler guarda y responde `303 See Other` al listado.

## 8. Generación y verificación

Los archivos `*_templ.go` son generados y están en `.gitignore`, así que nunca se
editan a mano:

```bash
make templ   # templ generate
make check   # go build ./... + go vet ./...
```

Ambos pasan sin errores.

Nota: `web.hurl` crea primero el huésped principal (con `sin_reserva: 1`) y
después la reserva enviando `id_huesped`, porque desde `72b5a55` ese campo es
obligatorio. El alta de reserva asigna ese huésped a la reserva recién creada,
así que los casos combinados siguen encontrando al huésped de la reserva.

## 9. Resumen de archivos modificados

| Archivo | Cambio |
|---|---|
| `main.go` | 6 rutas `GET/POST .../editar` |
| `pkg/handlers/departamentoHandlersWEB.go` | `EditarDepartamentoWeb`, `ActualizarDepartamentoWeb` |
| `pkg/handlers/huespedHandlersWEB.go` | `EditarHuespedWeb`, `ActualizarHuespedWeb` |
| `pkg/handlers/reservaHandlersWEB.go` | `primerHuespedDeReserva`, `EditarReservaWeb`, `ActualizarReservaWeb` |
| `views/entity_list.templ` | IDs, botones `Editar` y componentes `...Editado` |
| `views/entity_form.templ` | Tres formularios de edición |
| `views/modal.templ` | Tres wrappers de edición |
| `Documentación/documentacion-accion-editar.md` | Este documento |
