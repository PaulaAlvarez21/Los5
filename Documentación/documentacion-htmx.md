# HTMX en Los 5 Alquileres

Este documento describe la implementación de HTMX para crear y eliminar
departamentos, reservas y huéspedes sin recargar la página.

La idea central es:

- El formulario de alta se abre en un modal (`<dialog>`) cargado por HTMX.
- Al guardar, el servidor devuelve **solamente la fila nueva** (un `<tr>`) y
  HTMX la agrega al final del `<tbody>` con `hx-swap="beforeend"`.
- Al borrar, el botón envía un `DELETE` real y HTMX elimina la fila usando
  `hx-target="closest tr"` + `hx-swap="outerHTML"`.
- Se mantiene el comportamiento viejo sin HTMX (formularios con `action` y
  `method="post"`, redirecciones `303`) para no romper enlaces directos,
  pruebas existentes ni navegación sin JavaScript.

No se agregó JavaScript propio: todo el comportamiento cliente lo aporta HTMX
y el cierre del modal usa el `<form method="dialog">` nativo del elemento
`<dialog>`.

---

## 1. Script de HTMX y contenedores globales

Archivo: `views/layout.templ`

Se agregó el script en el `<head>` usando el CDN:

```html
<script src="https://cdn.jsdelivr.net/npm/htmx.org@2.0.11/dist/htmx.min.js"></script>
```

No se agregó `htmx-config` ni ninguna configuración extra.

Dentro del layout se agregaron dos contenedores globales:

```html
<main class="container main-content">
    <div id="page-errors"></div>
    { children... }
</main>

<div id="modal-root"></div>
```

- `#modal-root`: recibe el modal con el formulario de alta.
- `#page-errors`: recibe los errores de borrado (para no reemplazar la fila).

---

## 2. Modal reutilizable

Archivo nuevo: `views/modal.templ`

Se usa un `<dialog open>` para que aparezca abierto ni bien HTMX lo inserta en
`#modal-root`. El cierre se hace con un `<form method="dialog">`:

```templ
templ Modal(titulo string, mensaje string) {
    <dialog open class="app-modal">
        <div class="modal-box">
            <div class="modal-header">
                <h3>{ titulo }</h3>
                <form method="dialog">
                    <button type="submit" class="modal-close" aria-label="Cerrar">×</button>
                </form>
            </div>
            @ErrorBanner(mensaje)
            { children... }
        </div>
    </dialog>
}
```

Además se crearon tres wrappers que combinan el modal con cada formulario:

- `DepartamentoFormModal(mensaje string)`
- `ReservaFormModal(deptos, huespedes, mensaje)`
- `HuespedFormModal(reservas, mensaje)`

El parámetro `mensaje` se usa para volver a mostrar el formulario con el error
de validación adentro cuando el alta falla.

---

## 3. Formularios preparados para HTMX

Archivo: `views/entity_form.templ`

A los tres formularios se les agregaron atributos HTMX, conservando `action` y
`method` como fallback sin JavaScript.

Departamento:

```html
<form action="/departamentos" method="post" class="form-grid"
    hx-post="/departamentos"
    hx-target="#departamento-rows"
    hx-swap="beforeend">
```

Reserva:

```html
<form action="/reservas" method="post" class="form-grid"
    hx-post="/reservas"
    hx-target="#reserva-rows"
    hx-swap="beforeend">
```

Huésped (caso especial, ver §9):

```html
<form action="/huespedes" method="post" class="form-grid"
    hx-post="/huespedes"
    hx-target="#modal-root"
    hx-swap="innerHTML">
```

En departamentos y reservas, `hx-target` apunta al `<tbody>` y
`hx-swap="beforeend"` inserta la respuesta al final. Por eso el servidor debe
responder **una fila `<tr>`**, no la tabla completa. En huéspedes el formulario
apunta al modal y la fila se agrega mediante un swap OOB.

---

## 4. Filas reutilizables y tablas con `tbody` identificado

Archivo: `views/entity_list.templ`

Se extrajeron las filas a componentes reutilizables:

- `DepartamentoRow(d dominio.Departamento)`
- `ReservaRow(r dominio.ReservaConDetalle)`
- `HuespedRow(h dominio.Huesped)`

Y filas vacías con ID:

- `DepartamentoEmptyRow()` → `<tr id="departamento-empty">`
- `ReservaEmptyRow()` → `<tr id="reserva-empty">`
- `HuespedEmptyRow()` → `<tr id="huesped-empty">`

Cada listado quedó con un `tbody` identificado:

```html
<tbody id="departamento-rows">
<tbody id="reserva-rows">
<tbody id="huesped-rows">
```

Esto permite que HTMX agregue filas con `beforeend` sin recargar la página.

---

## 5. Botones de borrado

En cada fila se reemplazó el formulario viejo de eliminación por un botón que
usa `hx-delete`:

```html
<button
    type="button"
    class="btn btn-danger"
    hx-delete="/departamentos/7"
    hx-target="closest tr"
    hx-swap="outerHTML"
    hx-confirm="¿Borrar este departamento?">
    Borrar
</button>
```

Equivalentes para reservas y huéspedes:

```html
hx-delete="/reservas/{id}"
hx-delete="/huespedes/{id}"
```

Cómo funciona:

1. `hx-delete` envía un `DELETE` HTTP real.
2. `hx-target="closest tr"` selecciona la fila del botón.
3. El handler responde `200 OK` con body vacío.
4. `hx-swap="outerHTML"` reemplaza la fila por nada, es decir, la elimina.

`hx-confirm` muestra la confirmación del navegador antes de enviar.

---

## 6. Apertura de los modales

Archivos: `views/departamentos.templ`, `views/reservas.templ`,
`views/huespedes.templ`

Se eliminaron los bloques `<details>` que mostraban el formulario siempre
visible. En su lugar hay un botón en la cabecera de cada sección:

```html
<button
    class="btn btn-primary"
    hx-get="/departamentos/nuevo"
    hx-target="#modal-root"
    hx-swap="innerHTML">
    + Nuevo departamento
</button>
```

Equivalentes:

- `hx-get="/reservas/nuevo"`
- `hx-get="/huespedes/nuevo"`

Cada handler GET devuelve únicamente el modal, que HTMX inserta en
`#modal-root`.

---

## 7. Handlers de alta

Archivos:

- `pkg/handlers/departamentoHandlersWEB.go`
- `pkg/handlers/reservaHandlersWEB.go`
- `pkg/handlers/huespedHandlersWEB.go`

Los handlers POST existentes (`CrearDepartamentoWeb`, `CrearReservaWeb`,
`CrearHuespedWeb`) ahora distinguen requests de HTMX con `esHTMX(r)`.

### Alta exitosa con HTMX

1. Se crea la entidad con el service correspondiente.
2. Se vuelve a consultar el listado.
3. Se responde `200 OK` con el componente `...Creado`.
4. HTMX inserta la fila nueva al final del `tbody`.
5. Se cierran la fila vacía (si corresponde) y el modal con swaps OOB.

Ejemplo conceptual:

```go
if esHTMX(r) {
    departamentos, err := h.departamentoService.ObtenerDepartamentos(r.Context())
    if err != nil {
        h.responderErrorAlta(w, r, statusDeError(err), mensajeDeError(err))
        return
    }

    renderHTML(w, r, http.StatusOK, views.DepartamentoCreado(
        creado,
        len(departamentos) == 1,
    ))
    return
}

http.Redirect(w, r, "/departamentos", http.StatusSeeOther)
```

En reservas, después de crearla se consulta la página completa de reservas y
se busca el detalle de la reserva creada (`buscarReservaDetalle`) para poder
renderizar `ReservaRow`. La reserva necesita datos combinados
(departamento/huésped) que no vienen en el alta.

### Alta con error y HTMX

Se responde `200 OK` con el modal otra vez, pero indicando a HTMX que el swap
debe ir a `#modal-root`:

```go
func retarget(w http.ResponseWriter, target string) {
    w.Header().Set("HX-Retarget", target)
    w.Header().Set("HX-Reswap", "innerHTML")
}
```

El formulario se vuelve a mostrar con `ErrorBanner` adentro del modal. Se usa
`200` para que HTMX haga el swap sin necesitar configuración especial.

### Alta sin HTMX

Se conserva el comportamiento anterior: se renderiza la página con el error o
se redirige con `http.StatusSeeOther` (`303`) al listado.

Esto es lo que siguen usando `web.hurl` y cualquier navegación sin JavaScript.

---

## 8. Handlers de borrado

Se agregaron métodos nuevos:

- `EliminarDepartamentoWeb`
- `EliminarReservaWeb`
- `EliminarHuespedWeb`

Cada uno:

1. Lee el ID con `parseID(r)` (usa `r.PathValue("id")`).
2. Llama al service (`EliminarDepartamento`, `EliminarReserva`,
   `EliminarHuesped`), que termina en el `Delete...` de sqlc.
3. Si la request es de HTMX, responde `200 OK` vacío.
4. Si no es de HTMX, redirige `303` al listado.

```go
if esHTMX(r) {
    w.WriteHeader(http.StatusOK)
    return
}

http.Redirect(w, r, "/departamentos", http.StatusSeeOther)
```

### Error de borrado con HTMX

Se redirige el swap a `#page-errors` y se responde un `ErrorBanner`:

```go
func responderErrorHTMX(w http.ResponseWriter, r *http.Request, err error, target string) {
    retarget(w, target)
    renderHTML(w, r, http.StatusOK, views.ErrorBanner(mensajeDeError(err)))
}
```

Así el error se ve en la página sin reemplazar la fila que se intentó borrar.

---

## 9. Fragmentos de respuesta con swaps OOB

Archivo: `views/entity_list.templ`

Cuando se crea la primera entidad, la tabla muestra una fila
“No hay ... registrados”. Esa fila hay que quitarla al insertar la primera fila
real. Para eso se usan *out-of-band swaps*.

Componentes:

- `DepartamentoCreado(d, quitarVacio)`
- `ReservaCreada(r, quitarVacio)`
- `HuespedCreado(h, quitarVacio)`

Ejemplo:

```templ
templ DepartamentoCreado(d dominio.Departamento, quitarVacio bool) {
    @DepartamentoRow(d)
    if quitarVacio {
        <tr id="departamento-empty" hx-swap-oob="delete"></tr>
    }
    <div id="modal-root" hx-swap-oob="innerHTML"></div>
}
```

- La primera `<tr>` es la respuesta normal: HTMX la agrega al `tbody`.
- La `<tr id="departamento-empty" hx-swap-oob="delete">` elimina la fila vacía.
- El `<div id="modal-root" hx-swap-oob="innerHTML">` vacía el contenedor y
  cierra el modal.

`quitarVacio` es `true` cuando el listado recién consultado tiene exactamente
una entidad (la recién creada).

### Detalle importante: no envolver la fila OOB en `<template>`

En la primera versión se envolvió la fila OOB en `<template>` para que el HTML
fuera más “válido”. Eso **no funciona**: HTMX busca `[hx-swap-oob]` con
`querySelectorAll` sobre el fragmento de respuesta y no entra al contenido de
un `<template>`.

Por eso las filas OOB se emiten como hermanas directas de la fila nueva:

```html
<tr>...fila nueva...</tr>
<tr id="departamento-empty" hx-swap-oob="delete"></tr>
<div id="modal-root" hx-swap-oob="innerHTML"></div>
```

El parser HTML del navegador conserva los `<tr>` en el fragmento que arma HTMX,
y el procesamiento OOB ocurre antes del swap normal, por lo que la fila OOB se
procesa y se quita del fragmento antes de insertar la fila nueva.

### Caso especial: alta de huéspedes (OOB con `beforeend`)

El formulario de huésped apunta a `#modal-root` con `hx-swap="innerHTML"` para
poder cerrar el modal y, en el atajo desde "Nueva reserva", reemplazarlo por el
disparador que reabre el formulario de reserva. Por eso la fila nueva **no** es
la respuesta normal: se inserta con un swap OOB desde `HuespedCreado`.

Acá el OOB **no** puede ir directamente sobre el `<tr>`:

```html
<!-- Incorrecto -->
<tr hx-swap-oob="beforeend:#huesped-rows"> ... </tr>
```

htmx inserta el **contenido** del elemento marcado cuando el swap no es
`outerHTML`. Con el atributo en el `<tr>`, los `<td>` caen como hijos directos
del `<tbody>` sin su fila contenedora. El navegador agrupa todas esas celdas
sueltas en una única fila anónima, así que cada huésped nuevo aparece a la
derecha del anterior en lugar de en una línea nueva (se notaba al agregar el
tercer huésped).

La solución es envolver la fila en un `<tbody>` que lleve el OOB:

```templ
templ HuespedRowOOB(h dominio.Huesped) {
	<tbody hx-swap-oob="beforeend:#huesped-rows">
		<tr>...</tr>
	</tbody>
}
```

htmx mueve los hijos del `<tbody>` —es decir, el `<tr>` completo— al final de
`#huesped-rows`. El wrapper no se inserta: es solo el contenedor del contenido
que se quiere agregar. El `<tr id="huesped-empty" hx-swap-oob="delete">` sigue
siendo un hermano directo, igual que en los otros listados.

---

## 10. Helpers

Archivo: `pkg/handlers/helpers.go`

Se agregaron tres helpers:

```go
func esHTMX(r *http.Request) bool {
    return r.Header.Get("HX-Request") == "true"
}

func retarget(w http.ResponseWriter, target string) {
    w.Header().Set("HX-Retarget", target)
    w.Header().Set("HX-Reswap", "innerHTML")
}

func responderErrorHTMX(w http.ResponseWriter, r *http.Request, err error, target string) {
    retarget(w, target)
    renderHTML(w, r, http.StatusOK, views.ErrorBanner(mensajeDeError(err)))
}
```

No se modificó `renderHTML` ni el manejo de errores de dominio existente:
`statusDeError` y `mensajeDeError` se siguen usando.

---

## 11. Rutas nuevas

Archivo: `main.go`

Se agregaron:

| Método   | Ruta                    | Handler                     |
|----------|-------------------------|-----------------------------|
| `GET`    | `/departamentos/nuevo`  | `NuevoDepartamentoWeb`      |
| `GET`    | `/reservas/nuevo`       | `NuevaReservaWeb`           |
| `GET`    | `/huespedes/nuevo`      | `NuevoHuespedWeb`           |
| `DELETE` | `/departamentos/{id}`   | `EliminarDepartamentoWeb`   |
| `DELETE` | `/reservas/{id}`        | `EliminarReservaWeb`        |
| `DELETE` | `/huespedes/{id}`       | `EliminarHuespedWeb`        |

Se conservaron las rutas viejas como fallback:

- `POST /departamentos/{id}/eliminar`
- `POST /reservas/{id}/eliminar`
- `POST /huespedes/{id}/eliminar`

No hay conflicto con las rutas existentes porque Go 1.22+ enruta por método
HTTP (`GET`, `POST`, `DELETE`) además del patrón de path.

---

## 12. Estilos del modal

Archivo: `views/styles.css`

Se agregaron estilos para:

- `.app-modal`: `<dialog>` a pantalla completa con fondo oscuro y blur.
- `.app-modal .modal-box`: alto máximo y scroll interno.
- `.modal-header`: título y botón de cierre en la misma línea.
- `.modal-close`: botón “×” sin fondo.
- `.modal-close:hover`: color de peligro.

El modal reutiliza la clase `.modal-box` que ya existía en la hoja de estilos.

---

## 13. Flujo completo

### Crear

1. El usuario hace clic en “+ Nuevo departamento”.
2. HTMX hace `GET /departamentos/nuevo`.
3. El handler devuelve `DepartamentoFormModal`.
4. HTMX inserta el modal en `#modal-root`.
5. El usuario completa y envía el formulario.
6. HTMX hace `POST /departamentos` con `HX-Request: true`.
7. El handler crea la entidad y consulta el listado.
8. Devuelve la fila nueva + OOB para quitar la fila vacía (si aplica) + OOB
   para vaciar `#modal-root`.
9. HTMX inserta la fila al final del `tbody` y cierra el modal.
10. No hay recarga completa de la página.

### Borrar

1. El usuario hace clic en “Borrar” en una fila.
2. `hx-confirm` pide confirmación.
3. HTMX hace `DELETE /departamentos/{id}`.
4. El handler elimina y responde `200 OK` vacío.
5. HTMX reemplaza `closest tr` por nada.
6. La fila desaparece sin recargar la página.

### Errores

- Alta inválida: se vuelve a mostrar el modal con el `ErrorBanner` adentro
  (usando `HX-Retarget: #modal-root`).
- Borrado inválido o fallido: se muestra el `ErrorBanner` en `#page-errors`
  (usando `HX-Retarget: #page-errors`).
- Sin HTMX: se mantienen las páginas de error y las redirecciones `303`.

---

## 14. Verificación realizada

Se ejecutó:

```bash
templ generate
go build ./...
go vet ./...
make test-web
```

Resultados:

- `templ generate` completó correctamente.
- `go build ./...` y `go vet ./...` no reportaron errores.
- `make test-web` pasó los 32 requests de `web.hurl`, incluyendo altas por
  formulario con redirección `303` y las validaciones existentes.

También se verificó el fragmento generado por `DepartamentoCreado`: el servidor
emite una `<tr>` nueva, la fila OOB de la tabla vacía y el `div` OOB para
`#modal-root`, en ese orden.

Además se levantó una instancia del binario compilado en un puerto alternativo
y se probaron los tres flujos nuevos:

- `GET /departamentos/nuevo` devuelve el `<dialog>` con el formulario.
- `POST /departamentos`, `/huespedes` y `/reservas` con `HX-Request: true`
  devuelven la fila nueva, la fila OOB de la tabla vacía (cuando el listado
  estaba vacío) y el OOB de `#modal-root`.
- `DELETE` con `HX-Request: true` devuelve `200 OK` con body vacío.
- El mismo `POST` sin `HX-Request` sigue devolviendo `303` al listado.

### Prueba manual con curl

El servidor debe estar corriendo y la base levantada (`make db`).

```bash
# Abrir el modal (devuelve el HTML del dialog)
curl -s http://localhost:8080/departamentos/nuevo

# Crear con HTMX (devuelve la fila nueva + OOB)
curl -s -X POST http://localhost:8080/departamentos \
  -H 'HX-Request: true' \
  --data-urlencode 'nombre=Depto HTMX' \
  --data-urlencode 'direccion=Calle 1' \
  --data-urlencode 'disponible=true' \
  --data-urlencode 'limpio=true'

# Borrar con HTMX (devuelve 200 con body vacío)
curl -i -X DELETE http://localhost:8080/departamentos/1 \
  -H 'HX-Request: true'
```

---

## 15. Notas de mantenimiento

- `views/*_templ.go` lo genera templ y está en `.gitignore`. Después de tocar
  cualquier `.templ` hay que ejecutar `make templ` (o `templ generate`). El
  `Dockerfile` ya lo hace en la etapa builder.
- `db/sqlc/` también es generado. Las llamadas a `Delete...` no cambiaron; los
  handlers siguen usando los métodos de service existentes.
- Si se agrega una entidad nueva con el mismo patrón, hay que replicar:
  1. formulario con `hx-post`, `hx-target` y `hx-swap="beforeend"`;
  2. `tbody` con ID y componente de fila;
  3. botón con `hx-delete`, `closest tr` y `outerHTML`;
  4. handler `GET .../nuevo`, handler `DELETE /.../{id}` y respuesta de fila
     creada con OOB;
  5. rutas en `main.go`.
