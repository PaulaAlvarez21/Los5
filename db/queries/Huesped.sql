-- name: GetHuesped :one
SELECT id_huesped, id_reserva, nombre, apellido, telefono, email, observaciones
FROM HUESPED
WHERE id_huesped = $1;

-- name: ListHuespedes :many
SELECT id_huesped, id_reserva, nombre, apellido, telefono, email, observaciones
FROM HUESPED
ORDER BY apellido, nombre;

-- name: CreateHuesped :one
INSERT INTO HUESPED (
    id_reserva, nombre, apellido, telefono, email, observaciones
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id_huesped, id_reserva, nombre, apellido, telefono, email, observaciones;

-- name: UpdateHuesped :exec
UPDATE HUESPED
SET id_reserva = $2,
    nombre = $3,
    apellido = $4,
    telefono = $5,
    email = $6,
    observaciones = $7
WHERE id_huesped = $1;

-- name: DeleteHuesped :exec
DELETE FROM HUESPED
WHERE id_huesped = $1;