APP_NAME := Los5
DB_URL := postgres://user:password@localhost:5432/mydb?sslmode=disable

.PHONY: all run generate migrate apply status build test clean db

all: build

run:
	@air

generate:
	@sqlc generate

migrate:
	@test -n "$(name)" || (echo "Uso: make migrate name=nombre" && exit 1)
	atlas migrate diff "$(name)" --dir "file://db/migrations" --to \
	"file://db/schema/schema.sql" --dev-url "docker://postgres/18/dev?search_path=public"

apply:
	atlas migrate apply --dir "file://db/migrations" --url "$(DB_URL)"

status:
	atlas migrate status --dir "file://db/migrations" --url "$(DB_URL)"

db:
	@docker compose up -d
	@sleep 3

build: generate
	@mkdir -p tmp
	@go build -o tmp/$(APP_NAME) .

test: db
	@go test ./...

clean:
	@rm -f $(APP_NAME)
	@rm -rf tmp
