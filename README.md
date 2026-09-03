Ejecutar con comando: docker compose up
(ó docker compose up -d para segundo plano)


1.instalar sqlc (por si no lo tenemos local (las dos))
2. hacer sqlc generate 
3. agregar la carpeta a gitignore
4.hacer los test (reserva_test.go, ...) que llamen a los generados por sglc generate
5. hacer makefile que automatice este proceso make test( instale sqlc, sqlc generate, air, dockercompose up)
6.opcional agregar sqlc generate a docker para que quede dockerizado
8. hacer nuevo main vacio
9. Documentar todo en README