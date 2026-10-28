# Los5

El proyecto es un gestor de reservas de alojamientos temporales: los administradores de la empresa "Los 5" podrán usarlo para gestionar el alquiler de sus departamentos.

En esta entrega (**TP3**) el proyecto expone una **API REST** sobre esa base de datos. El servidor está escrito en Go con `net/http` de la librería estándar y separa la aplicación en capas (`handlers → logica → repositorios → db/sqlc`), usando **structs de dominio** como contrato entre ellas. La capa de acceso a datos sigue generándose con **sqlc** a partir de queries SQL, y los endpoints se prueban de punta a punta con **HURL**.

La arquitectura en capas está explicada en detalle en [`documentacion-estructura.md`](./documentacion-estructura.md).

## Stack

- **Go** (backend, `net/http` de la STL)
- **PostgreSQL** (base de datos)
- **sqlc** (genera el Go tipado a partir de las queries SQL)
- **Atlas** (migraciones de esquema)
- **Air** (recompilación automática en desarrollo)
- **Docker Compose** (levanta la base de datos)
- **HURL** (prueba los endpoints de la API)
- **make** (orquesta todo el flujo)

## Estructura

```
.
├── main.go                 # cableado de capas y registro de rutas HTTP
├── pkg/
│   ├── dominio/            # structs de entidades y AppError (contrato entre capas)
│   ├── handlers/           # lee/decodifica JSON y responde códigos HTTP
│   ├── logica/             # reglas de negocio y validaciones (services)
│   └── repositorios/       # traduce dominio ↔ sqlc y ejecuta las queries
├── db/
│   ├── queries/            # queries SQL que sqlc convierte en funciones Go
│   ├── schema/             # esquema deseado de la base (source de verdad)
│   ├── sqlc/               # código Go generado por sqlc (no se edita ni versiona)
│   └── testing/            # tests de acceso a datos (se versionan)
│       └── *_test.go       # paquete externo db_test que importa Los5/db/sqlc
├── requests.hurl           # pruebas de cada endpoint de la API (HURL)
├── documentacion-estructura.md  # explicación de las capas
├── Makefile                # automatiza build, test, test-api, migraciones, etc.
├── docker-compose.yaml     # integra base de datos PostgreSQL
├── sqlc.yaml               # config de sqlc
├── Dockerfile              # build de la app en dos etapas
└── air.toml                # config de Air (hot reload)
```

## Cómo correr

> Casi todo está automatizado con **make**, así que con un par de comandos alcanza.

### Primera vez

Instalar las herramientas de desarrollo (solo una vez por máquina):

```bash
# instalacion docker compose
sudo apt install docker-compose-plugin #V2

# instalacion go
sudo apt install golang-go

#instalacion de sqlc
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

#instalacion de la herramienta make
sudo apt install make

#instalacion de HURL (para probar los endpoints)
sudo apt install hurl            # Debian/Ubuntu
# alternativas: brew install hurl  /  snap install hurl  /  cargo install hurl

opcionales:
#atlas
curl -sSf https://atlasgo.sh | sh  # o: go install ariga.io/atlas/cmd/atlas@latest
#air
go install github.com/air-verse/air@latest
```

Asegurate de tener `$(go env GOPATH)/bin` en el `PATH`.

### Probar la API (TP3)

El objetivo de la consigna es: se clona el repo, se entra a la branch `tp3` y con **un solo comando** se levanta todo, se prueban los endpoints y se limpia. Ese comando es:

```bash
git clone <URL-del-repo>
cd Los5
git checkout tp3
make test-api
```

`make test-api` encadena todas las tareas necesarias, en este orden:

1. **`.env`** → si no existe, copia `.env.example` (credenciales genéricas) a `.env`.
2. **Tareas previas**
   - `build` → ejecuta `sqlc generate` (regenera el código tipado) y compila el binario en `tmp/Los5`.
   - `borrar` → `docker compose down -v` (borra contenedores **y volúmenes**, para arrancar de cero).
   - `db` → `docker compose up -d --wait database` (levanta la base y **espera** a que el healthcheck pase).
   - Levanta la aplicación (`./tmp/Los5`) con las variables del `.env`.
3. **Prueba de los endpoints** → `hurl --test requests.hurl` recorre cada endpoint y valida método, código HTTP y cuerpo.
4. **Tareas posteriores** → mata la app, borra `tmp/` y hace `docker compose down -v` (borra contenedores y volúmenes).

Si HURL reporta algún fallo, el comando termina con error (`exit status` distinto de 0), así sirve tal cual para CI.

### Ejecutar los tests de acceso a datos

Para correr solo los tests de la base (paquete `db/testing`):

```bash
make test
```

### Otros comandos útiles

```bash
make build      # genera código con sqlc y compila el binario
make generate   # solo regenera el código sqlc
make run        # ejecuta el servidor con Air (hot reload)
make db         # solo levanta la base de datos en Docker
make migrate name=cambios    # genera una migración nueva con Atlas
make apply      # aplica las migraciones a la base
make status     # muestra el estado de las migraciones
make clean      # borra artefactos de build
make borrar     # borra volumenes del contenedor
make .env       # copia las credenciales genericas de la plantilla al .env local
```

## API REST

El servidor escucha en `http://localhost:8080`. Todas las entidades exponen el mismo CRUD.

| Entidad        | Método | Ruta                  | Descripción                       | Éxito |
|----------------|--------|-----------------------|-----------------------------------|-------|
| Departamentos  | GET    | `/departamentos`      | Lista todos los departamentos     | 200   |
| Departamentos  | GET    | `/departamentos/{id}` | Obtiene un departamento por ID    | 200   |
| Departamentos  | POST   | `/departamentos`      | Crea un departamento              | 201   |
| Departamentos  | PUT    | `/departamentos/{id}` | Actualiza un departamento         | 200   |
| Departamentos  | DELETE | `/departamentos/{id}` | Elimina un departamento           | 204   |
| Huéspedes      | GET    | `/huespedes`          | Lista todos los huéspedes         | 200   |
| Huéspedes      | GET    | `/huespedes/{id}`     | Obtiene un huésped por ID         | 200   |
| Huéspedes      | POST   | `/huespedes`          | Crea un huésped                   | 201   |
| Huéspedes      | PUT    | `/huespedes/{id}`     | Actualiza un huésped              | 200   |
| Huéspedes      | DELETE | `/huespedes/{id}`     | Elimina un huésped                | 204   |
| Reservas       | GET    | `/reservas`           | Lista todas las reservas          | 200   |
| Reservas       | GET    | `/reservas/{id}`      | Obtiene una reserva por ID        | 200   |
| Reservas       | POST   | `/reservas`           | Crea una reserva                  | 201   |
| Reservas       | PUT    | `/reservas/{id}`      | Actualiza una reserva             | 200   |
| Reservas       | DELETE | `/reservas/{id}`      | Elimina una reserva               | 204   |

### Casos de uso combinados (HTML)

Estos endpoints consultan más de una entidad y responden con HTML usando `a-h/templ`.

| Caso de uso | Método | Ruta | Descripción | Éxito |
|---|---|---|---|---|
| Reservas por departamento | GET | `/web/reservas?depto_id={id}` | Lista las reservas de un departamento | 200 |
| Huéspedes por departamento | GET | `/web/departamentos/{id}/huespedes` | Lista los huéspedes de las reservas de un departamento | 200 |
| Disponibilidad | GET | `/web/departamentos/{id}/disponibilidad?desde=YYYY-MM-DD&hasta=YYYY-MM-DD` | Indica si el departamento está disponible en el rango | 200 |
| Total de reserva | GET | `/web/reservas/{id}/total` | Calcula subtotal, descuento y total de la reserva | 200 |

Códigos de error: `400` (JSON, ID o fechas inválidas, o validación de negocio), `404` (recurso inexistente), `500` (error interno).

### Ejemplo con curl

Con la base levantada (`make db`) y la app corriendo (`make run`), por ejemplo:

```bash
# Crear un departamento
curl -i -X POST http://localhost:8080/departamentos \
  -H 'Content-Type: application/json' \
  -d '{"nombre":"Departamento Centro","direccion":"Av. España 123","disponible":true,"limpio":true,"descripcion":"De dos ambientes"}'

# Listar departamentos
curl -i http://localhost:8080/departamentos
```

La prueba automatizada de **todos** los endpoints (caso feliz de cada uno) está en [`requests.hurl`](./requests.hurl), que es lo que corre `make test-api`.

## Decisiones que fuimos tomando

- **Arquitectura en capas con `dominio` como contrato**: los handlers no conocen la base, la lógica no conoce HTTP ni sqlc, y el repositorio es el único que habla con sqlc. Las capas se inyectan por constructor en `main.go`. El detalle está en `documentacion-estructura.md`.

- **Error con código HTTP incorporado (`AppError`)**: cada capa crea el error donde sabe qué significa (`NotFoundError` → 404, `ValidationError` → 400) y el handler solo lee el `Status`. Cualquier otro error se responde como 500 sin filtrar detalles internos.

- **HURL en vez de curl suelto para la prueba de endpoints**: un solo archivo (`requests.hurl`) declara cada request con su código esperado y asserts sobre el JSON. `hurl --test` lo ejecuta de corrido y devuelve un exit code, lo que hace la prueba repetible y apta para CI.

- **`make test-api` como comando único**: junta tareas previas (sqlc, compilación, limpieza, levantar base y app), la prueba de endpoints y la limpieza posterior. Clonar, `git checkout tp3`, `make test-api`.

- **La app corre local y la base en Docker para `test-api`**: el `test-api` levanta solo el servicio `database` de Compose (con `--wait` sobre el healthcheck) y ejecuta el binario local, para mantener el ciclo de prueba liviano.

- **Las herramientas de desarrollo no van en el Makefile ni en Docker**: sqlc, Atlas y Air solo se necesitan en la máquina del desarrollador, no para correr la app. Meterlas en el Makefile para instalarlas se podría hacer, pero se decantó por dejarlas como instalación manual única; y en Docker directamente no conviene, porque la imagen debería ser mínima (sqlc/atlas/air la agrandarían y contradicen el Dockerfile liviano de dos etapas). HURL sí es necesario para `make test-api` y por eso se documenta su instalación.

- **sqlc en vez de escribir SQL a mano o un ORM**: las queries se escriben en SQL y sqlc genera todo el código Go tipado (mapper). Nos da seguridad (tipo de compilación), evitamos errores de tipeo y no perdemos control sobre el SQL. Se alinea con lo que plantean las filminas como alternativa intermedia entre SQL a mano y un ORM completo.

- **Tests separados por operación (una por función)**: en vez de un solo `main` que recorre todo el CRUD, cada operación (`Create`, `Get`, `List`, `Update`, `Delete`) tiene su propio test. Permite aislar y fallar con claridad cada operación.

- **IDs auto-generados con `SERIAL PRIMARY KEY`**: fiel a las filminas. El `Create` no recibe el ID, lo genera la base y los tests usan el ID que devuelve la DB (nada de IDs hardcodeados).

- **Limpieza con `t.Cleanup`**: los tests crean datos de prueba y los borran solos al terminar, falle o no. Así la base queda limpia entre corridas y los tests son determinísticos (dan igual la primera que la centésima vez). Esto solo corre en `go test`; la producción no se toca.

- **Verificación del borrado con `sql.ErrNoRows`**: en los tests de `Delete`, al volver a pedir el registro eliminado, se espera exactamente `sql.ErrNoRows` (y no "cualquier error"), igual que en las filminas. Así distinguimos "no existe" de un error real de conexión.

- **Tests afuera en `db/testing/`**: los tests no viven en la misma carpeta que el código generado sino en `db/testing/`, como paquete de test externo (`package db_test`) que importa `Los5/db/sqlc` y antepone `sqlc.` a los símbolos del paquete. Así `go test ./...` los encuentra (los descubre solo por el patrón `./...`) y el código generado por sqlc queda intacto y sin tocar.

- **Solo se versionan los tests, no el código generado**: `.gitignore` ignora los `.go` de `db/sqlc/` (los regenera `sqlc generate`), mientras que los tests viven en `db/testing/`, que no está ignorado. Así se versionan y revisan los tests, y el código generado se mantiene fuera del repo.

- **PostgreSQL 18 con montaje en `/var/lib/postgresql`**: las imágenes de postgres 18+ cambiaron dónde guardan los datos, así que el volumen se monta en `/var/lib/postgresql` (en vez de `/var/lib/postgresql/data`) para evitar el error de "unused mount/volume".

- **El esquema se auto-ejecuta al levantar la DB**: `schema.sql` está montado en `/docker-entrypoint-initdb.d/`. Postgres lo corre automáticamente la primera vez (o tras `docker compose down -v`), así no hace falta hacer `docker exec` manual para crear las tablas.

- **Migraciones con Atlas**: el esquema en `db/schema/schema.sql` es la "source de verdad"; Atlas compara ese estado deseado contra la base y genera migraciones versionadas. Después de tocar el esquema, conviene `make migrate` + `make apply`.

- **Air para desarrollo**: con `make run` Air observa los archivos, recompila y reinicia solo. El hot reload del navegador no lo maneja (habría que sumar otra herramienta).

- **Archivo .env.example**: se creo el archivo .env.example que si se sube al github para permitir que se pueda correr en otras maquinas con las credenciales genericas sin tener que subir nuestro .env local ni nuestros secretos. La automatizacion en Make con el target .env verifica si existe o no un .env y si no existe copia los valores de ese archivo plantilla al .env local de la maquina.

## Pendientes

- **Los tests de datos corren contra la misma base de Docker**: alcanza con `make test` que levanta la DB y corre todo. Si algún día se quiere una base separada para tests, sería el paso a dar (ej. `mydb_test`).

- **`make test-api` espera con un `sleep` fijo**: el target levanta la app y espera unos segundos antes de correr HURL. Si la máquina es lenta podría no alcanzar; la mejora natural es esperar activamente (polling) a que la app responda antes de probar.

## Problema solucionado

Cuando cambiamos el esquema a `SERIAL PRIMARY KEY`, a quien hacía pull le rompían los tests con errores raros (`duplicate key violates ..._pkey` y "id debería ser generado"), aunque la base se levantaba bien. Nos llevó un rato darnos cuenta de qué pasaba. `make test` levanta la base y corre los tests, pero **no** regenera el código de sqlc (eso solo lo hace `make build`). Como ese código generado **no se versiona** (está en `.gitignore`), la persona que hizo pull seguía compilando con sus `.go` locales, que eran de la versión anterior del esquema.

No era un problema de la base ni de los tests, sino de que el código de sqlc local estaba desactualizado. Lo resolvimos regenerando a mano con `sqlc generate` y ahí sí, `make test` pasó. Ahora con la incorporacion de herramientas de desarrollo esto no deberia volver a suceder.
