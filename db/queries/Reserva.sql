-- name: GetReserva :one
SELECT id_reserva, fecha_inicio, id_depto, fecha_fin, precio_base, cant_noches, descuento, observaciones
FROM RESERVA
WHERE id_reserva = $1;

-- name: ListReservas :many
SELECT id_reserva, fecha_inicio, id_depto, fecha_fin, precio_base, cant_noches, descuento, observaciones
FROM RESERVA
ORDER BY fecha_inicio;

-- name: CreateReserva :one
INSERT INTO RESERVA (
    fecha_inicio, id_depto, fecha_fin,
    precio_base, cant_noches, descuento, observaciones
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id_reserva, fecha_inicio, id_depto, fecha_fin, precio_base, cant_noches, descuento, observaciones;

-- name: UpdateReserva :exec
UPDATE RESERVA
SET fecha_inicio = $2,
    id_depto = $3,
    fecha_fin = $4,
    precio_base = $5,
    cant_noches = $6,
    descuento = $7,
    observaciones = $8
WHERE id_reserva = $1;

-- name: DeleteReserva :exec
DELETE FROM RESERVA
WHERE id_reserva = $1;