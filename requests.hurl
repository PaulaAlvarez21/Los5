# DEPARTAMENTOS

# Crear departamento (POST)
POST http://localhost:8080/departamentos
Content-Type: application/json
{
  "nombre": "Departamento Centro",
  "direccion": "Av. España 123",
  "disponible": true,
  "limpio": true,
  "descripcion": "Departamento de dos ambientes"
}
HTTP 201
[Asserts]
jsonpath "$.id_depto" isInteger
jsonpath "$.nombre" == "Departamento Centro"

# Listar departamentos (GET)
GET http://localhost:8080/departamentos
HTTP 200
[Asserts]
jsonpath "$.length()" > 0
jsonpath "$[0].nombre" == "Departamento Centro"

# Obtener departamento por ID (GET)
GET http://localhost:8080/departamentos/1
HTTP 200
[Asserts]
jsonpath "$.nombre" == "Departamento Centro"

# Actualizar departamento (PUT)
PUT http://localhost:8080/departamentos/1
Content-Type: application/json
{
  "nombre": "Departamento Centro Actualizado",
  "direccion": "Av. España 123",
  "disponible": false,
  "limpio": true,
  "descripcion": "Departamento actualizado"
}
HTTP 200
[Asserts]
jsonpath "$.nombre" == "Departamento Centro Actualizado"


# ==========================================

# RESERVAS

# Crear reserva (POST)
POST http://localhost:8080/reservas
Content-Type: application/json
{
  "id_depto": 1,
  "fecha_inicio": "2026-11-01T00:00:00Z",
  "fecha_fin": "2026-11-05T00:00:00Z",
  "precio_base": "15000.00",
  "cant_noches": 4,
  "descuento": "0.00",
  "observaciones": "reserva nueva"
}
HTTP 201
[Asserts]
jsonpath "$.id_reserva" isInteger
jsonpath "$.precio_base" == "15000.00"
jsonpath "$.observaciones" == "reserva nueva"

# Listar reservas (GET)
GET http://localhost:8080/reservas
HTTP 200
[Asserts]
jsonpath "$.length()" > 0
jsonpath "$[0].id_depto" == 1
jsonpath "$[0].precio_base" == "15000.00"
jsonpath "$[0].observaciones" == "reserva nueva"

# Obtener reserva por ID (GET)
GET http://localhost:8080/reservas/1
HTTP 200
[Asserts]
jsonpath "$.id_depto" == 1
jsonpath "$.precio_base" == "15000.00"
jsonpath "$.observaciones" == "reserva nueva"

# Actualizar reserva (PUT)
PUT http://localhost:8080/reservas/1
Content-Type: application/json
{
  "id_depto": 1,
  "fecha_inicio": "2026-11-01T00:00:00Z",
  "fecha_fin": "2026-11-05T00:00:00Z",
  "precio_base": "18000.00",
  "cant_noches": 4,
  "descuento": "1000.00",
  "observaciones": "reserva actualizada"
}
HTTP 200
[Asserts]
jsonpath "$.precio_base" == "18000.00"
jsonpath "$.observaciones" == "reserva actualizada"


# ==========================================

# HUESPEDES

# Crear huésped (POST)
POST http://localhost:8080/huespedes
Content-Type: application/json
{
  "id_reserva": 1,
  "nombre": "Maria",
  "apellido": "Gonzalez",
  "telefono": "2494123456",
  "email": "maria@example.com",
  "observaciones": "Sin observaciones"
}
HTTP 201
[Asserts]
jsonpath "$.id_huesped" isInteger
jsonpath "$.nombre" == "Maria"

# Listar huéspedes (GET)
GET http://localhost:8080/huespedes
HTTP 200
[Asserts]
jsonpath "$.length()" > 0
jsonpath "$[0].nombre" == "Maria"
jsonpath "$[0].apellido" == "Gonzalez"

# Obtener huésped por ID (GET)
GET http://localhost:8080/huespedes/1
HTTP 200
[Asserts]
jsonpath "$.nombre" == "Maria"
jsonpath "$.apellido" == "Gonzalez"

# Actualizar huésped (PUT)
PUT http://localhost:8080/huespedes/1
Content-Type: application/json
{
  "id_reserva": 1,
  "nombre": "Maria",
  "apellido": "Gonzalez Actualizada",
  "telefono": "2494123456",
  "email": "maria@example.com",
  "observaciones": "Datos actualizados"
}
HTTP 200
[Asserts]
jsonpath "$.apellido" == "Gonzalez Actualizada"


# ==========================================

# ELIMINACIONES


# 1. Eliminar huésped (depende de reserva)
DELETE http://localhost:8080/huespedes/1
HTTP 204

# 2. Eliminar reserva (depende de departamento)
DELETE http://localhost:8080/reservas/1
HTTP 204

# 3. Eliminar departamento
DELETE http://localhost:8080/departamentos/1
HTTP 204