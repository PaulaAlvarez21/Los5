# Pruebas del flujo web server-rendered (sin JavaScript).
#
# La base se reinicia antes de correr estas pruebas (make test-web hace
# `docker compose down -v` y vuelve a levantar PostgreSQL), por eso los IDs
# arrancan en 1.

# ==========================================
# DASHBOARD Y ESTÁTICOS
# ==========================================

GET http://localhost:8080/
HTTP 200
[Asserts]
header "Content-Type" contains "text/html"
body contains "Los 5 Alquileres"
body contains "href=\"/departamentos\""
body contains "href=\"/reservas\""
body contains "href=\"/huespedes\""

GET http://localhost:8080/static/pico.min.css
HTTP 200
[Asserts]
header "Content-Type" contains "text/css"

GET http://localhost:8080/static/styles.css
HTTP 200
[Asserts]
header "Content-Type" contains "text/css"


# ==========================================
# DEPARTAMENTOS
# ==========================================

GET http://localhost:8080/departamentos
HTTP 200
[Asserts]
body contains "Catálogo de departamentos"
body contains "No hay departamentos registrados"

# Alta por formulario (PRG: 303 hacia el listado)
POST http://localhost:8080/departamentos
[FormParams]
nombre: Departamento Centro
direccion: Av. España 123
descripcion: Departamento de dos ambientes
disponible: true
limpio: true
HTTP 303
[Asserts]
header "Location" == "/departamentos"

GET http://localhost:8080/departamentos
HTTP 200
[Asserts]
body contains "Departamento Centro"
body contains "Av. España 123"
body contains "Limpio"
body contains "/departamentos/1/reservas"

# Alternar higiene
POST http://localhost:8080/departamentos/1/limpieza
HTTP 303
[Asserts]
header "Location" == "/departamentos"

GET http://localhost:8080/departamentos
HTTP 200
[Asserts]
body contains "Mantenimiento"

# Validación del formulario
POST http://localhost:8080/departamentos
[FormParams]
nombre: Sin direccion
direccion:
HTTP 400
[Asserts]
body contains "datos de departamento inválidos"


# ==========================================
# RESERVAS
# ==========================================

GET http://localhost:8080/reservas
HTTP 200
[Asserts]
body contains "Gestión de reservas"
body contains "No hay reservas registradas"

# El alta de reserva exige un huésped existente: se crea primero con el atajo
# "sin reserva" (mismo flujo que el botón "+ Crear huésped" del formulario).
GET http://localhost:8080/huespedes
HTTP 200
[Asserts]
body contains "Padrón de huéspedes"
body contains "No hay huéspedes registrados"

POST http://localhost:8080/huespedes
[FormParams]
sin_reserva: 1
nombre: Maria
apellido: Gonzalez
telefono: 2494123456
email: maria@example.com
observaciones: Huésped de prueba
HTTP 303
[Asserts]
header "Location" == "/reservas"

POST http://localhost:8080/reservas
[FormParams]
id_depto: 1
id_huesped: 1
fecha_inicio: 2026-11-01
fecha_fin: 2026-11-05
cant_noches: 4
precio_base: 15000.00
descuento: 10
observaciones: Reserva de prueba
HTTP 303
[Asserts]
header "Location" == "/reservas"

GET http://localhost:8080/reservas
HTTP 200
[Asserts]
body contains "Departamento Centro"
body contains "2026-11-01"
body contains "2026-11-05"
body contains "4 noc."
body contains "15000.00"

# Filtro por departamento
GET http://localhost:8080/reservas?depto_id=1&desde=2026-11-01&hasta=2026-11-30
HTTP 200
[Asserts]
body contains "Departamento Centro"

# Validaciones de reservas
POST http://localhost:8080/reservas
[FormParams]
id_depto: 1
id_huesped: 1
fecha_inicio: 2026-11-05
fecha_fin: 2026-11-01
cant_noches: 4
precio_base: 15000.00
HTTP 400
[Asserts]
body contains "cantidad de noches inválida"

GET http://localhost:8080/reservas?depto_id=abc
HTTP 400

GET http://localhost:8080/reservas?desde=fecha-invalida
HTTP 400


# ==========================================
# HUÉSPEDES
# ==========================================

# Alta de un segundo huésped, ya asociado a la reserva creada.
POST http://localhost:8080/huespedes
[FormParams]
id_reserva: 1
nombre: Juan
apellido: Perez
telefono: 2494654321
email: juan@example.com
HTTP 303
[Asserts]
header "Location" == "/huespedes"

GET http://localhost:8080/huespedes
HTTP 200
[Asserts]
body contains "Maria"
body contains "Gonzalez"
body contains "maria@example.com"

POST http://localhost:8080/huespedes
[FormParams]
nombre:
apellido: Perez
telefono: 2494123456
HTTP 400
[Asserts]
body contains "datos de huesped inválidos"


# ==========================================
# CASOS DE USO COMBINADOS
# ==========================================

# Reservas de un departamento
GET http://localhost:8080/departamentos/1/reservas
HTTP 200
[Asserts]
body contains "Reservas de Departamento Centro"
body contains "2026-11-01"
body contains "15000.00"

# Huéspedes de un departamento
GET http://localhost:8080/departamentos/1/huespedes
HTTP 200
[Asserts]
body contains "Huéspedes de Departamento Centro"
body contains "Maria"
body contains "Gonzalez"

# Disponibilidad: rango sin reservas
GET http://localhost:8080/departamentos/1/disponibilidad?desde=2026-12-01&hasta=2026-12-05
HTTP 200
[Asserts]
body contains "Disponible"
body contains "0"

# Disponibilidad: rango que se solapa con la reserva existente
GET http://localhost:8080/departamentos/1/disponibilidad?desde=2026-11-02&hasta=2026-11-04
HTTP 200
[Asserts]
body contains "No disponible"

# Formulario de disponibilidad (sin fechas)
GET http://localhost:8080/departamentos/1/disponibilidad
HTTP 200
[Asserts]
body contains "Consultar"

# Total calculado: 15000.00 * 4 noches - 10% = 54000.00
GET http://localhost:8080/reservas/1/total
HTTP 200
[Asserts]
body contains "60000.00"
body contains "10.00%"
body contains "54000.00"

# Errores de los casos de uso combinados
GET http://localhost:8080/departamentos/abc/reservas
HTTP 400

GET http://localhost:8080/departamentos/999/reservas
HTTP 404

GET http://localhost:8080/departamentos/999/huespedes
HTTP 404

GET http://localhost:8080/reservas/999/total
HTTP 404

GET http://localhost:8080/departamentos/1/disponibilidad?desde=2026-11-05&hasta=2026-11-01
HTTP 400

GET http://localhost:8080/departamentos/1/disponibilidad?desde=fecha-invalida&hasta=2026-11-01
HTTP 400
