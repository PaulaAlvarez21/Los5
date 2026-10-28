-- name: ListReservasByDepto :many
SELECT id_reserva, fecha_inicio, id_depto, fecha_fin, precio_base, cant_noches, descuento, observaciones
FROM RESERVA
WHERE id_depto = sqlc.arg(id_depto)
ORDER BY fecha_inicio;

-- name: ListHuespedesByDepto :many
SELECT id_huesped, id_reserva, nombre, apellido, telefono, email, observaciones
FROM HUESPED
WHERE id_reserva IN (
    SELECT id_reserva
    FROM RESERVA
    WHERE id_depto = sqlc.arg(id_depto)
)
ORDER BY apellido, nombre;

-- name: CountReservasSolapadas :one
SELECT COUNT(*)
FROM RESERVA
WHERE id_depto = sqlc.arg(id_depto)
  AND fecha_inicio < sqlc.arg(hasta)
  AND fecha_fin > sqlc.arg(desde);
