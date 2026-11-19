# TP5 — Interfaz web server-rendered con templ

Fecha: 2026-11-04

## 1. Objetivo

Reemplazar la interfaz JavaScript por una interfaz **server-rendered** con
[`a-h/templ`](https://templ.guide/), sin JavaScript, manteniendo la arquitectura
en capas del proyecto (`handlers → logica → repositorios → db/sqlc`) y
conservando los handlers JSON de la API aunque ya no se registren.

Además se agrega el caso de uso **reservas por departamento** con su cadena
completa: ruta, handler, lógica, repositorio, query SQL y template.

## 2. Alcance de los cambios

- La interfaz principal pasa a ser HTML generado en el servidor.
- No se usa JavaScript: no hay `fetch`, `onclick`, `app.js` ni frameworks de
  frontend.
- Los formularios usan `GET`/`POST` normales y el patrón
  **Post/Redirect/Get (PRG)**.
- Los handlers API se conservan y siguen compilando, pero **no se registran**
  en `main.go`.
- Se separan los handlers por entidad y por casos combinados.
- Se agrega `GET /departamentos/{id}/reservas`.
- Se agregan pruebas de flujo HTML con `web.hurl` y el comando `make test-web`.
- Se documenta todo en este archivo.

## 3. Arquitectura

```
┌─────────────────────────────────────────────────────────────┐
│ Navegador                                                   │
│  - HTML, formularios GET/POST, links                        │
│  - Sin JavaScript                                           │
└───────────────────────────┬─────────────────────────────────┘
                            │ HTTP
┌───────────────────────────▼─────────────────────────────────┐
│ main.go                                                     │
│  - conecta PostgreSQL                                       │
│  - crea sqlc, repositorio y services                        │
│  - instancia SOLO handlers web                              │
│  - registra rutas web y /static/                            │
└───────────────────────────┬─────────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────┐
│ pkg/handlers                                                │
│  - *WEB.go: leen formularios/query params, renderizan HTML   │
│  - *API.go: handlers JSON conservados, sin rutas             │
│  - helpers.go: render, errores, parseo de IDs/fechas         │
└───────────────────────────┬─────────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────┐
│ pkg/logica                                                  │
│  - validaciones y reglas de negocio                         │
│  - combina datos para los casos de uso                      │
└───────────────────────────┬─────────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────┐
│ pkg/repositorios                                            │
│  - traduce dominio ↔ sqlc                                   │
│  - traduce errores SQL a AppError                           │
└───────────────────────────┬─────────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────┐
│ db/queries + db/sqlc                                        │
│  - SQL versionado y código Go generado                      │
└─────────────────────────────────────────────────────────────┘
```

### Capas

| Capa | Responsabilidad | No conoce |
|---|---|---|
| `pkg/handlers` | HTTP, formularios, render HTML, códigos de estado | SQL ni sqlc |
| `pkg/logica` | Reglas de negocio y validaciones | HTTP, HTML ni sqlc |
| `pkg/repositorios` | Acceso a datos y conversión de tipos | HTTP ni HTML |
| `db/sqlc` | Código generado a partir de SQL | Reglas de negocio |

Los structs de `pkg/dominio` son el contrato entre capas.

## 4. Estructura de archivos relevante

```
.
├── main.go                          # wiring y registro de rutas web
├── web.hurl                         # pruebas HURL del flujo HTML
├── tp5_documentacion.md             # este documento
├── Makefile                         # build, generate, check, test, test-web
├── pkg/
│   ├── dominio/                     # entidades, DTOs combinados y AppError
│   ├── handlers/
│   │   ├── departamentoHandlersAPI.go
│   │   ├── huespedHandlersAPI.go
│   │   ├── reservaHandlersAPI.go
│   │   ├── departamentoHandlersWEB.go
│   │   ├── huespedHandlersWEB.go
│   │   ├── reservaHandlersWEB.go
│   │   ├── combinadoHandlersWEB.go
│   │   └── helpers.go
│   ├── logica/
│   │   ├── departamentoService.go
│   │   ├── huespedService.go
│   │   ├── reservaService.go
│   │   └── combinadoService.go
│   └── repositorios/
│       └── BDRepository.go
├── db/
│   ├── queries/
│   │   ├── Departamento.sql
│   │   ├── Huesped.sql
│   │   ├── Reserva.sql
│   │   └── Combinado.sql
│   └── schema/schema.sql
└── views/
    ├── layout.templ                 # layout base, header y nav
    ├── index.templ                  # dashboard
    ├── entity_form.templ            # formularios reutilizables
    ├── entity_list.templ            # listados reutilizables
    ├── departamentos.templ
    ├── reservas.templ
    ├── huespedes.templ
    ├── reservas_por_departamento.templ
    ├── huespedes_por_departamento.templ
    ├── disponibilidad.templ
    ├── total_reserva.templ
    ├── error.templ
    ├── viewmodels.go                # datos que reciben las páginas
    ├── helpers.go                   # helpers de templates
    ├── assets.go                    # go:embed de CSS
    ├── pico.min.css                 # Pico.css local
    └── styles.css                   # estilos propios
```

## 5. Separación de handlers

### Handlers API (conservados, sin rutas)

| Archivo | Tipo | Responsabilidad |
|---|---|---|
| `departamentoHandlersAPI.go` | `DepartamentoAPIHandler` | CRUD JSON de departamentos |
| `huespedHandlersAPI.go` | `HuespedAPIHandler` | CRUD JSON de huéspedes |
| `reservaHandlersAPI.go` | `ReservaAPIHandler` | CRUD JSON de reservas y helper `responderError` |

Estos archivos siguen compilando y pueden reutilizarse si en el futuro se
quiere volver a exponer una API. En esta entrega **no se instancian ni se
registran** en `main.go`.

### Handlers web (registrados)

| Archivo | Tipo | Responsabilidad |
|---|---|---|
| `departamentoHandlersWEB.go` | `DepartamentoWebHandler` | Listado, alta y cambio de higiene |
| `huespedHandlersWEB.go` | `HuespedWebHandler` | Listado y alta de huéspedes |
| `reservaHandlersWEB.go` | `ReservaWebHandler` | Listado con filtros y alta de reservas |
| `combinadoHandlersWEB.go` | `CombinadoWebHandler` | Casos de uso combinados, incluido reservas por departamento |

### Helpers compartidos (`pkg/handlers/helpers.go`)

| Helper | Uso |
|---|---|
| `renderHTML` | Escribe un componente templ con el status indicado |
| `responderErrorWeb` | Traduce `AppError` a página HTML de error |
| `statusDeError` / `mensajeDeError` | Extraen status y mensaje del error |
| `parseID` | Lee y valida `{id}` de la ruta |
| `parseFecha` | Parsea `YYYY-MM-DD` |
| `parseRango` | Valida un rango `desde < hasta` |
| `parseFiltroReservas` | Lee `depto_id`, `desde` y `hasta` del query string |
| `filtroVista` | Conserva los valores del filtro para re-renderizar el formulario |
| `stringOpcional` | Convierte `""` en `nil` |
| `parseBoolForm` | Interpreta checkboxes (`true`, `on`, `1`, `si`) |

## 6. Rutas registradas

Solo se registran rutas web. Las rutas API no existen en esta entrega.

| Método | Ruta | Handler | Descripción |
|---|---|---|---|
| GET | `/` | `views.IndexPage` | Dashboard con links a las secciones |
| GET | `/departamentos` | `DepartamentoWebHandler.DepartamentosWeb` | Formulario + listado |
| POST | `/departamentos` | `DepartamentoWebHandler.CrearDepartamentoWeb` | Alta de departamento |
| POST | `/departamentos/{id}/limpieza` | `DepartamentoWebHandler.CambiarHigieneWeb` | Alterna higiene |
| GET | `/reservas` | `ReservaWebHandler.ReservasWeb` | Formulario + filtros + listado |
| POST | `/reservas` | `ReservaWebHandler.CrearReservaWeb` | Alta de reserva |
| GET | `/huespedes` | `HuespedWebHandler.HuespedesWeb` | Formulario + listado |
| POST | `/huespedes` | `HuespedWebHandler.CrearHuespedWeb` | Alta de huésped |
| GET | `/departamentos/{id}/reservas` | `CombinadoWebHandler.ReservasPorDepartamentoWeb` | **Caso de uso nuevo** |
| GET | `/departamentos/{id}/huespedes` | `CombinadoWebHandler.HuespedesPorDepartamentoWeb` | Huéspedes de un departamento |
| GET | `/departamentos/{id}/disponibilidad` | `CombinadoWebHandler.DisponibilidadWeb` | Formulario o resultado |
| GET | `/reservas/{id}/total` | `CombinadoWebHandler.TotalReservaWeb` | Total calculado |
| GET | `/static/...` | `http.FileServer` | CSS embebido |

### Parámetros

- `GET /reservas` acepta `depto_id`, `desde` y `hasta`.
- `GET /departamentos/{id}/disponibilidad` acepta `desde` y `hasta`. Sin fechas
  muestra el formulario; con fechas muestra el resultado.
- Todos los IDs y fechas inválidas responden `400`.
- Un departamento o reserva inexistente responde `404`.

## 7. Cadena de datos por request web

No se duplicaron queries ni lógica: cada request reutiliza la cadena existente.

| Request | Handler | Service | Repositorio | Query SQL |
|---|---|---|---|---|
| `GET /departamentos` | `DepartamentosWeb` | `DepartamentoService.ObtenerDepartamentos` | `ListDepartamentos` | `ListDepartamentos` |
| `POST /departamentos` | `CrearDepartamentoWeb` | `DepartamentoService.CrearDepartamento` | `CreateDepartamento` | `CreateDepartamento` |
| `POST /departamentos/{id}/limpieza` | `CambiarHigieneWeb` | `ObtenerDepartamento` + `ActualizarDepartamento` | `GetDepartamento` + `UpdateDepartamento` | `GetDepartamento` + `UpdateDepartamento` |
| `GET /reservas` | `ReservasWeb` | `CombinadoService.ObtenerPaginaReservas` | `ListReservas`, `ListReservasByDepto`, `ListHuespedes`, `ListHuespedesByDepto`, `ListDepartamentos` | mismas queries |
| `POST /reservas` | `CrearReservaWeb` | `ReservaService.CrearReserva` (+ `HuespedService.ActualizarHuesped` si se asigna huésped) | `CreateReserva` (+ `GetHuesped` + `UpdateHuesped`) | `CreateReserva` (+ `GetHuesped` + `UpdateHuesped`) |
| `GET /huespedes` | `HuespedesWeb` | `HuespedService.ObtenerHuespedes` + `ReservaService.ObtenerReservas` | `ListHuespedes` + `ListReservas` | `ListHuespedes` + `ListReservas` |
| `POST /huespedes` | `CrearHuespedWeb` | `HuespedService.CrearHuesped` | `CreateHuesped` | `CreateHuesped` |
| `GET /departamentos/{id}/reservas` | `ReservasPorDepartamentoWeb` | `CombinadoService.ObtenerReservasPorDepartamento` | `GetDepartamento` + `ListReservasByDepto` | `GetDepartamento` + `ListReservasByDepto` |
| `GET /departamentos/{id}/huespedes` | `HuespedesPorDepartamentoWeb` | `CombinadoService.ObtenerHuespedesPorDepartamento` | `GetDepartamento` + `ListHuespedesByDepto` | `GetDepartamento` + `ListHuespedesByDepto` |
| `GET /departamentos/{id}/disponibilidad` | `DisponibilidadWeb` | `CombinadoService.ConsultarDisponibilidad` | `GetDepartamento` + `CountReservasSolapadas` | `GetDepartamento` + `CountReservasSolapadas` |
| `GET /reservas/{id}/total` | `TotalReservaWeb` | `CombinadoService.CalcularTotalReserva` | `GetReserva` | `GetReserva` |

### Caso de uso nuevo: reservas por departamento

Cadena completa:

1. **Ruta:** `GET /departamentos/{id}/reservas`.
2. **Handler:** `CombinadoWebHandler.ReservasPorDepartamentoWeb` valida el ID y
   llama al service.
3. **Lógica:** `CombinadoService.ObtenerReservasPorDepartamento` valida el ID,
   obtiene el departamento y lista sus reservas.
4. **Repositorio:** `BDRepository.GetDepartamento` y
   `BDRepository.ListReservasByDepto`.
5. **SQL:** `GetDepartamento` y `ListReservasByDepto` (ya existían en
   `db/queries/Combinado.sql` y `db/queries/Departamento.sql`).
6. **Template:** `views.ReservasPorDepartamentoPage` con
   `ReservasPorDepartamentoPageData`.

El listado incluye ID, período, noches, precio base, descuento y link al total.

## 8. Templates y view models

### Layout

`views/layout.templ` define:

- `Layout(title, activeTab)`: HTML base, Pico.css, `styles.css`, header y main.
- `Header(activeTab)`: navegación a `/`, `/departamentos`, `/reservas`,
  `/huespedes`.
- `ErrorBanner(mensaje)`: alerta HTML cuando hay un error de validación.

No se incluye ningún `<script>`.

### Componentes reutilizables

| Archivo | Componentes |
|---|---|
| `entity_form.templ` | `DepartamentoForm`, `ReservaForm`, `HuespedForm` |
| `entity_list.templ` | `DepartamentoList`, `ReservaList`, `HuespedList`, `ReservaSimpleList` |

### Páginas

| Archivo | Componente | Datos |
|---|---|---|
| `index.templ` | `IndexPage` | Dashboard con links |
| `departamentos.templ` | `DepartamentosPage` | `DepartamentosPageData` |
| `reservas.templ` | `ReservasPage` | `ReservasPageData` |
| `huespedes.templ` | `HuespedesPage` | `HuespedesPageData` |
| `reservas_por_departamento.templ` | `ReservasPorDepartamentoPage` | `ReservasPorDepartamentoPageData` |
| `huespedes_por_departamento.templ` | `HuespedesPorDepartamentoPage` | `HuespedesPorDepartamentoPageData` |
| `disponibilidad.templ` | `DisponibilidadPage` | `dominio.Disponibilidad` |
| `total_reserva.templ` | `TotalReservaPage` | `dominio.TotalReserva` |
| `error.templ` | `ErrorPage` | status y mensaje |

Los view models están en `views/viewmodels.go` para no pasar muchos parámetros
sueltos y conservar los valores de los filtros al re-renderizar.

### CSS

- `views/pico.min.css`: Pico.css descargado localmente. No se usa CDN.
- `views/styles.css`: paleta y componentes propios (paneles, dashboard,
  formularios, tablas, badges, alertas).
- `views/assets.go`: embebe ambos archivos con `//go:embed`.
- `main.go` los sirve en `/static/` con `http.FileServer`.

Esto permite que la imagen Docker no dependa de copiar la carpeta `views/`.

## 9. Formularios, PRG y errores

### Post/Redirect/Get

- `POST` correcto → `303 See Other` hacia el listado:
  - `/departamentos`
  - `/reservas`
  - `/huespedes`
- `POST` con error de validación → se vuelve a renderizar la página con
  `ErrorBanner` y status `400`.
- ID inválido o recurso inexistente → `views.ErrorPage` con `400` o `404`.

Esto evita reenvíos al refrescar y mantiene URLs limpias.

### Validaciones

| Formulario | Reglas |
|---|---|
| Departamento | `nombre` y `direccion` obligatorios |
| Reserva | `id_depto > 0`, `cant_noches > 0`, `fecha_inicio < fecha_fin`, `precio_base` presente |
| Huésped | `nombre` y `apellido` obligatorios; `id_reserva` opcional |
| Filtros de reservas | `depto_id > 0`, fechas válidas, `desde < hasta` |
| Disponibilidad | `desde` y `hasta` obligatorios y `desde < hasta` |

## 10. Wiring en `main.go`

1. Conecta a PostgreSQL con `pgx`.
2. Crea `sqlc.New(connection)`.
3. Crea `repo.NewBDRepository(queries)`.
4. Crea los cuatro services.
5. Instancia **solo** los handlers web:
   - `NewDepartamentoWebHandler`
   - `NewReservaWebHandler`
   - `NewHuespedWebHandler`
   - `NewCombinadoWebHandler`
6. Registra las rutas web y `/static/`.
7. Inicia `http.ListenAndServe(":8080", nil)`.

Los constructores `NewDepartamentoAPIHandler`, `NewHuespedAPIHandler` y
`NewReservaAPIHandler` existen pero no se llaman.

## 11. Makefile

| Comando | Qué hace |
|---|---|
| `make generate` | `sqlc generate` + `templ generate` |
| `make sqlc` | Regenera `db/sqlc` |
| `make templ` | Regenera `views/*_templ.go` |
| `make build` | Genera y compila `tmp/Los5` |
| `make check` | Genera, compila todos los paquetes y corre `go vet` |
| `make run` | Air con hot reload |
| `make run-app` | Ciclo completo: build, DB limpia, ejecuta y limpia al salir |
| `make test` | Tests de acceso a datos |
| `make test-web` | Levanta DB + app y corre `hurl --test web.hurl` |
| `make test-api` | Legado: corre `requests.hurl` (la API ya no tiene rutas) |

> `db/sqlc/*` y `views/*_templ.go` son generados y están en `.gitignore`.
> Siempre hay que correr `make generate` (o `make build`) después de cambiar
> queries o templates.

## 12. Pruebas web (`web.hurl`)

`web.hurl` prueba el flujo HTML completo. Se ejecuta con:

```bash
make test-web
```

Cobertura:

1. **Dashboard y estáticos**
   - `GET /` responde 200 y contiene links a las tres secciones.
   - `/static/pico.min.css` y `/static/styles.css` responden 200 con
     `Content-Type: text/css`.
2. **Departamentos**
   - Listado vacío.
   - Alta por formulario → `303` a `/departamentos`.
   - El listado muestra el departamento creado.
   - Cambio de higiene → `303` y el listado muestra `Mantenimiento`.
   - Alta inválida → `400` con mensaje de validación.
3. **Reservas**
   - Listado vacío.
   - Alta por formulario → `303` a `/reservas`.
   - El listado muestra departamento, fechas, noches y precio.
   - Filtro por departamento y rango de fechas.
   - Fechas inválidas y `depto_id` inválido → `400`.
4. **Huéspedes**
   - Listado vacío.
   - Alta por formulario → `303` a `/huespedes`.
   - El listado muestra nombre, apellido y email.
   - Alta inválida → `400`.
5. **Casos combinados**
   - Reservas por departamento.
   - Huéspedes por departamento.
   - Disponibilidad con y sin solapamiento.
   - Formulario de disponibilidad sin fechas.
   - Total calculado (`60000.00` subtotal, `10.00%`, `54000.00` total).
   - IDs y fechas inválidas → `400`/`404`.

`requests.hurl` se conserva como referencia histórica de la API JSON, pero no
se ejecuta desde `make test-web` porque esas rutas ya no están registradas.

## 13. Decisiones de diseño

- **Server-side rendering con templ**: cumple el requisito de no usar
  JavaScript y mantiene el HTML tipado y testeable.
- **Rutas en la raíz**: `/departamentos`, `/reservas`, `/huespedes` en lugar
  de `/web/...`, para que la interfaz web sea la interfaz principal.
- **API conservada pero no registrada**: los handlers JSON quedan como
  referencia y base para una futura API, sin afectar la seguridad ni el
  enrutado actual.
- **Handlers separados por entidad y combinado**: cada archivo tiene una
  responsabilidad clara y evita un `webHandlers.go` monolítico.
- **Helpers compartidos**: evita duplicar render, manejo de errores y parseo.
- **View models**: las páginas reciben un struct con todo lo que necesitan, en
  lugar de muchos parámetros sueltos.
- **CSS local embebido**: no hay CDN ni dependencia de archivos externos en
  runtime; la imagen Docker es autocontenida.
- **PRG**: evita reenvíos de formularios y mantiene URLs limpias.
- **Reutilización de queries**: no se agregaron queries duplicadas; el caso de
  reservas por departamento usa `ListReservasByDepto`, que ya existía.
- **Sin JavaScript**: no hay `app.js`, `fetch`, `onclick` ni modales
  controlados por JS; los formularios y links son HTML estándar.

## 14. Cómo correr

### Ciclo completo

```bash
make run-app
```

Luego abrir <http://localhost:8080>.

### Desarrollo con hot reload

```bash
make run
```

### Verificación rápida

```bash
make check
```

### Pruebas web

```bash
make test-web
```

Requiere Docker para PostgreSQL y `hurl` instalado.

## 15. Limitaciones y trabajo futuro

- No hay autenticación ni autorización.
- No hay paginación en los listados.
- No hay edición ni baja desde la interfaz web; solo alta, listado y cambio de
  higiene. La API JSON conservada sí tiene CRUD completo.
- El total de reserva se calcula en el service a partir de `precio_base`,
  `cant_noches` y `descuento`; no incluye impuestos ni tarifas por temporada.
- Los handlers API podrían registrarse en el futuro bajo un prefijo `/api/`
  si se necesita exponerlos sin mezclarlos con las rutas web.

## 16. Verificación realizada

- `templ generate` sin errores.
- `go build ./...` sin errores.
- `go vet ./...` sin errores.
- `make test-web` ejecutado con Docker: **32 requests, 0 fallos**.
