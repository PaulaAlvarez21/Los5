-- name: GetDepartamento :one
SELECT id_depto, nombre, direccion, disponible, limpio, descripcion
FROM DEPARTAMENTO
WHERE id_depto = $1;

-- name: ListDepartamentos :many
SELECT id_depto, nombre, direccion, disponible, limpio, descripcion
FROM DEPARTAMENTO
ORDER BY nombre;

-- name: CreateDepartamento :one
INSERT INTO DEPARTAMENTO (
    id_depto, nombre, direccion, disponible, limpio, descripcion
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id_depto, nombre, direccion, disponible, limpio, descripcion;

-- name: UpdateDepartamento :exec
UPDATE DEPARTAMENTO
SET nombre = $2,
    direccion = $3,
    disponible = $4,
    limpio = $5,
    descripcion = $6
WHERE id_depto = $1;

-- name: DeleteDepartamento :exec
DELETE FROM DEPARTAMENTO
WHERE id_depto = $1;