# Acción de borrado en las listas de entidades

Fecha: 2026-11-10

## 1. Objetivo

Agregar la opción de **eliminar** cada registro desde las tablas de las páginas
`/departamentos`, `/huespedes` y `/reservas`, manteniendo la interfaz
**server-rendered con templ y sin JavaScript**.

La capa de lógica y la de persistencia ya tenían el borrado implementado, así que
el trabajo se limitó a la capa web: handlers, rutas, templates y estilos.

## 2. Lo que ya existía

Antes de este cambio, el borrado ya estaba resuelto de punta a punta salvo por la
interfaz web:

| Capa | Método existente |
|---|---|
| Lógica | `DepartamentoService.EliminarDepartamento` |
| Lógica | `HuespedService.EliminarHuesped` |
| Lógica | `ReservaService.EliminarReserva` |
| Repositorio | `BDRepository.DeleteDepartamento` |
| Repositorio | `BDRepository.DeleteHuesped` |
| Repositorio | `BDRepository.DeleteReserva` |
| SQL (sqlc) | `DeleteDepartamento`, `DeleteHuesped`, `DeleteReserva` |

También existían los handlers JSON `EliminarDepartamento`, `EliminarHuesped` y
`EliminarReserva` en los archivos `*HandlersAPI.go`, pero **no se usaron**: esos
handlers pertenecen a la API, no están registrados en `main.go` y responden JSON.
La solución web es independiente y no toca esa capa.

## 3. Decisión de diseño: POST en lugar de DELETE

Se usa `POST` y no `DELETE` porque un `<form>` HTML puro solo puede enviar `GET`
y `POST`. Para enviar un `DELETE` real desde el navegador haría falta JavaScript
(`fetch`) o un middleware de *method override*, que igualmente viaja como `POST`.

Además, un borrado nunca debe hacerse con un link `GET`:

- `GET` debe ser seguro (no modificar estado).
- El navegador puede precargar links, los crawlers seguirlos y el botón "atrás"
  re-disparar la acción.
- Un refresh accidental podría borrar datos.

Por eso el flujo es el clásico **Post/Redirect/Get (PRG)**:

```
Formulario (POST /entidad/{id}/eliminar)
        │
        ▼
Handler: parsea el ID, llama al service, borra
        │
        ▼
http.Redirect(..., http.StatusSeeOther)  →  GET /entidad
        │
        ▼
La lista se vuelve a renderizar sin el registro eliminado
```

No hay confirmación previa porque `confirm()` es JavaScript. Si se quisiera una
confirmación sin JS, habría que agregar una página intermedia con un formulario
que haga el POST definitivo.

## 4. Rutas agregadas

En `main.go`, junto a las rutas POST existentes:

```go
http.HandleFunc("POST /departamentos/{id}/eliminar", deptosHandler.EliminarDepartamentoWeb)
http.HandleFunc("POST /reservas/{id}/eliminar", reservasHandler.EliminarReservaWeb)
http.HandleFunc("POST /huespedes/{id}/eliminar", huespedesHandler.EliminarHuespedWeb)
```

Todas las rutas usan el patrón de `net/http` con `{id}`, que el handler lee con
`parseID(r)` (helper existente en `pkg/handlers/helpers.go`).

## 5. Handlers web agregados

Se agregó un método por entidad en su archivo `*HandlersWEB.go`. Todos siguen el
mismo patrón que `CambiarHigieneWeb`: parsear el ID, delegar en el service,
manejar el error renderizando la página correspondiente y redirigir con
`303 See Other` si todo salió bien.

### `pkg/handlers/departamentoHandlersWEB.go`

```go
func (h *DepartamentoWebHandler) EliminarDepartamentoWeb(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		h.renderDepartamentos(w, r, statusDeError(err), mensajeDeError(err))
		return
	}

	if err := h.departamentoService.EliminarDepartamento(r.Context(), id); err != nil {
		h.renderDepartamentos(w, r, statusDeError(err), mensajeDeError(err))
		return
	}

	http.Redirect(w, r, "/departamentos", http.StatusSeeOther)
}
```

### `pkg/handlers/huespedHandlersWEB.go`

```go
func (h *HuespedWebHandler) EliminarHuespedWeb(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		h.renderHuespedes(w, r, statusDeError(err), mensajeDeError(err))
		return
	}

	if err := h.huespedService.EliminarHuesped(r.Context(), id); err != nil {
		h.renderHuespedes(w, r, statusDeError(err), mensajeDeError(err))
		return
	}

	http.Redirect(w, r, "/huespedes", http.StatusSeeOther)
}
```

### `pkg/handlers/reservaHandlersWEB.go`

```go
func (h *ReservaWebHandler) EliminarReservaWeb(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		h.renderReservas(w, r, statusDeError(err), mensajeDeError(err))
		return
	}

	if err := h.reservaService.EliminarReserva(r.Context(), id); err != nil {
		h.renderReservas(w, r, statusDeError(err), mensajeDeError(err))
		return
	}

	http.Redirect(w, r, "/reservas", http.StatusSeeOther)
}
```

No se agregaron imports nuevos: los archivos ya importaban `net/http` y los
helpers `parseID`, `statusDeError` y `mensajeDeError` viven en `helpers.go`.

## 6. Cambios en los templates

Archivo: `views/entity_list.templ`.

El botón de borrado es un formulario POST con un `<button type="submit">`, igual
que el resto de las acciones que modifican estado:

```templ
<form action={ templ.URL(fmt.Sprintf("/entidad/%d/eliminar", id)) } method="post">
	<button type="submit" class="btn btn-danger">Eliminar</button>
</form>
```

### `DepartamentoList`

Ya tenía la columna **Acciones** y un `div.row-actions`; solo se agregó el
formulario de borrado junto a los botones de higiene, reservas, huéspedes y
disponibilidad.

### `HuespedList`

No tenía columna de acciones. Se agregó:

- `<th>Acciones</th>` al final del encabezado.
- Una celda nueva por fila con el formulario de borrado.
- El `colspan` de la fila vacía pasó de `4` a `5`.

### `ReservaList`

Se agregó:

- `<th>Acciones</th>` al final del encabezado.
- Un `div.row-actions` que agrupa el link existente **Ver total** y el nuevo
  formulario de borrado.
- El `colspan` de la fila vacía pasó de `7` a `8`.

`ReservaSimpleList` (usada en la página de reservas por departamento) se dejó
sin cambios para no alterar el flujo de esa vista.

## 7. Estilos

En `views/styles.css` se agregó la clase `.btn-danger`, reutilizando las
variables de color de peligro que ya existían:

```css
.btn-danger {
	background-color: var(--danger-text);
	color: #ffffff;
}

.btn-danger:hover {
	background-color: #991b1b;
}
```

## 8. Generación y verificación

Los archivos `*_templ.go` se regeneran con templ y están ignorados por Git
(`.gitignore`), por lo que no se editan a mano:

```bash
make templ    # templ generate
make check    # go build ./... + go vet ./...
```

Ambos comandos pasan sin errores.

## 9. Cómo probarlo

```bash
make run
```

Luego, en `/departamentos`, `/huespedes` o `/reservas`, cada fila muestra el
botón **Eliminar**. Al presionarlo:

1. El navegador hace `POST` a la ruta `.../{id}/eliminar`.
2. El handler borra el registro y responde `303 See Other`.
3. El navegador hace `GET` de la lista y el registro ya no aparece.

## 10. Limitaciones conocidas

### Restricciones de claves foráneas

El esquema no usa `ON DELETE CASCADE`:

- Borrar un **departamento con reservas** falla por la FK
  `fk_reserva_depto`.
- Borrar una **reserva con huéspedes asignados** falla por la FK
  `fk_huesped_reserva`.
- Borrar un **huésped** no tiene problemas, porque nadie lo referencia.

En los dos primeros casos, el error de PostgreSQL no se traduce a un error de
dominio: el handler muestra la página genérica `500 - error interno`. Una mejora
posible es detectar el error de FK (`pgconn.PgError`, código `23503`) en el
servicio o repositorio y devolver un `dominio.ValidationError` con un mensaje
claro del tipo "no se puede eliminar: tiene registros asociados".

### Sin confirmación

No hay diálogo de confirmación porque requeriría JavaScript. La alternativa
clásica sería una página intermedia `GET .../{id}/eliminar` que muestre un
formulario de confirmación cuyo POST ejecute el borrado.

## 11. Resumen de archivos modificados

| Archivo | Cambio |
|---|---|
| `main.go` | 3 rutas POST de borrado |
| `pkg/handlers/departamentoHandlersWEB.go` | `EliminarDepartamentoWeb` |
| `pkg/handlers/huespedHandlersWEB.go` | `EliminarHuespedWeb` |
| `pkg/handlers/reservaHandlersWEB.go` | `EliminarReservaWeb` |
| `views/entity_list.templ` | Botones y columnas de acciones en las 3 listas |
| `views/styles.css` | Clase `.btn-danger` |

No se modificaron servicios, repositorios, queries, sqlc, viewmodels,
formularios ni los handlers de la API.
