APP_NAME := Los5
DB_URL := postgres://user:password@localhost:5432/mydb?sslmode=disable

.PHONY: all run run-app reset generate sqlc templ migrate apply status build test test-api clean db borrar

all: build

# levanta todo desde un clon: .env + codigo generado + binario + DB + app
run-app: .env build db
	@set -a; . ./.env; set +a; \
	echo "App disponible en http://localhost:8080 (Ctrl+C para cortar)"; \
	./tmp/$(APP_NAME)

# hot reload con Air (tambien funciona desde un clon)
run: .env db
	@set -a; . ./.env; set +a; air

# vuelve la base a cero: borra el volumen y reaplica schema.sql
reset: borrar db

generate: sqlc templ

sqlc:
	@sqlc generate

templ:
	@templ generate

migrate:
	@test -n "$(name)" || (echo "Uso: make migrate name=nombre" && exit 1)
	atlas migrate diff "$(name)" --dir "file://db/migrations" --to \
	"file://db/schema/schema.sql" --dev-url "docker://postgres/18/dev?search_path=public"

apply:
	atlas migrate apply --dir "file://db/migrations" --url "$(DB_URL)"

status:
	atlas migrate status --dir "file://db/migrations" --url "$(DB_URL)"

db: .env
	@docker compose up -d --wait database

build: generate
	@mkdir -p tmp
	@go build -o tmp/$(APP_NAME) .


clean:
	@rm -rf ./tmp
 
borrar:
	@docker compose down -v


.env:
	@if [ ! -f .env ]; then \
        cp .env.example .env; \
    fi

test: .env build borrar db 
	@go test -v ./...
	@rm -rf ./tmp
	@docker compose down -v

test-api: .env build borrar db
	@set -a; . ./.env; set +a; \
	./tmp/$(APP_NAME) & echo $$! > tmp/server.pid; \
	sleep 3; \
	hurl --test requests.hurl; status=$$?; \
	kill $$(cat tmp/server.pid) 2>/dev/null; \
	rm -rf ./tmp; \
	docker compose down -v; \
	exit $$status