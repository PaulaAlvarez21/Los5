# Los5

Proyecto de Programación Web (Los 5): servidor web en Go con acceso a una base de datos PostgreSQL usando **sqlc** para generar el código de acceso a datos.

## Stack

- **Go** (backend)
- **PostgreSQL** (base de datos)
- **sqlc** (genera el Go tipado a partir de las queries SQL)
- **Atlas** (migraciones de esquema)
- **Air** (recompilación automática en desarrollo)
- **Docker Compose** (levanta la base de datos)
- **make** (orquesta todo el flujo)

## Estructura

```
.
├── db/
│   ├── queries/        # queries SQL que sqlc convierte en funciones Go
│   ├── schema/         # esquema deseado de la base (source de verdad)
│   └── sqlc/           # código Go generado por sqlc + los tests del paquete
│       ├── *_test.go   # tests de acceso a datos (se versionan)
│       └── *.go        # código generado por sqlc (no se edita ni versiona) 
├── Makefile            # automatiza build, test, migraciones, etc.
├── docker-compose.yaml # base de datos PostgreSQL
├── sqlc.yaml           # config de sqlc
```

## Cómo correr

> Casi todo está automatizado con **make**, así que con un par de comandos alcanza.

### Primera vez

Instalar las herramientas de desarrollo (solo una vez por máquina):

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
curl -sSf https://atlasgo.sh | sh          # o: go install ariga.io/atlas/cmd/atlas@latest
go install github.com/air-verse/air@latest
```

Asegurate de tener `$(go env GOPATH)/bin` en el `PATH`.

### Ejecutar los tests

Levanta la base con Docker y corre los tests:

```bash
make test
```

El target `test` hace `docker compose up -d` y después `go test ./...`.

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
```

## Decisiones que fuimos tomando

- **Las herramientas de desarrollo no van en el Makefile ni en Docker**: sqlc, Atlas y Air solo se necesitan en la máquina del desarrollador, no para correr la app. Meterlas en el Makefile para instalarlas se podría hacer, pero se decantó por dejarlas como instalación manual única; y en Docker directamente no conviene, porque la imagen debería ser mínima (sqlc/atlas/air la agrandarían y contradicen el Dockerfile liviano de dos etapas).

- **sqlc en vez de escribir SQL a mano o un ORM**: las queries se escriben en SQL y sqlc genera todo el código Go tipado (mapper). Nos da seguridad (tipo de compilación), evitamos errores de tipeo y no perdemos control sobre el SQL. Se alinea con lo que plantean las filminas como alternativa intermedia entre SQL a mano y un ORM completo.

- **Tests separados por operación (una por función)**: en vez de un solo `main` que recorre todo el CRUD, cada operación (`Create`, `Get`, `List`, `Update`, `Delete`) tiene su propio test. Permite aislar y fallar con claridad cada operación.

- **IDs auto-generados con `SERIAL PRIMARY KEY`**: fiel a las filminas. El `Create` no recibe el ID, lo genera la base y los tests usan el ID que devuelve la DB (nada de IDs hardcodeados).

- **Limpieza con `t.Cleanup`**: los tests crean datos de prueba y los borran solos al terminar, falle o no. Así la base queda limpia entre corridas y los tests son determinísticos (dan igual la primera que la centésima vez). Esto solo corre en `go test`; la producción no se toca.

- **Verificación del borrado con `sql.ErrNoRows`**: en los tests de `Delete`, al volver a pedir el registro eliminado, se espera exactamente `sql.ErrNoRows` (y no "cualquier error"), igual que en las filminas. Así distinguimos "no existe" de un error real de conexión.

- **Tests adentro de `db/sqlc/`**: los tests viven en la misma carpeta que el código generado, como tests internos del paquete `db` (sin prefijo `db.`), para que `go test ./...` los encuentre dentro del paquete donde están las queries. La raíz del proyecto queda sin tests.

- **Solo se versionan los tests, no el código generado**: `.gitignore` ignora los `.go` de `db/sqlc/` (los regenera `sqlc generate`), pero re-incluye los `_test.go` con una regla de negación (`!`). Así se versionan y revisan los tests, y el código generado se mantiene fuera del repo.

- **PostgreSQL 18 con montaje en `/var/lib/postgresql`**: las imágenes de postgres 18+ cambiaron dónde guardan los datos, así que el volumen se monta en `/var/lib/postgresql` (en vez de `/var/lib/postgresql/data`) para evitar el error de "unused mount/volume".

- **El esquema se auto-ejecuta al levantar la DB**: `schema.sql` está montado en `/docker-entrypoint-initdb.d/`. Postgres lo corre automáticamente la primera vez (o tras `docker compose down -v`), así no hace falta hacer `docker exec` manual para crear las tablas.

- **Migraciones con Atlas**: el esquema en `db/schema/schema.sql` es la "source de verdad"; Atlas compara ese estado deseado contra la base y genera migraciones versionadas. Después de tocar el esquema, conviene `make migrate` + `make apply`.

- **Air para desarrollo**: con `make run` Air observa los archivos, recompila y reinicia solo. El hot reload del navegador no lo maneja (habría que sumar otra herramienta).

## Pendientes

- **Atlas**: los targets del Makefile (`migrate`, `apply`, `status`) están definidos, pero todavía no hay carpeta `db/migrations/` ni migraciones generadas. Falta generar la primera migración con `make migrate`.
- **Air**: falta crear el `.air.toml` (config de archivos a observar y delegar el build en `make build`); hoy el target `run` solo llama a `air` sin configuración.

- **Los tests corren contra la misma base de Docker**: alcanza con `make test` que levanta la DB y corre todo. Si algún día se quiere una base separada para tests, sería el paso a dar (ej. `mydb_test`).

## Consultas
 - ¿Incluimos la instalacion de las herramientas de desarrollo en el makefile? (por ahora no los incorporamos)
 - ¿los tests van en raiz o dentro de db/sqlc/? (por ahora los dejamos dentro de db/sqlc)