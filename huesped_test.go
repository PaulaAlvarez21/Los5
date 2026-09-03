package main

import (
	"context"
	"database/sql"
	"testing"
	"time"

	db "Los5.com/ServidorWeb/db/sqlc"
)

func createTestHuesped(t *testing.T) db.Huesped {
	t.Helper()

	// crear departamento dependiente
	depto, err := testQueries.CreateDepartamento(context.Background(), db.CreateDepartamentoParams{
		IDDepto:     200,
		Nombre:      "Depto Huesped",
		Direccion:   "Calle Huesped 200",
		Disponible:  true,
		Limpio:      true,
		Descripcion: sql.NullString{String: "Depto para test de huesped", Valid: true},
	})
	if err != nil {
		t.Fatalf("no se pudo crear departamento para huesped: %v", err)
	}

	// crear reserva dependiente
	fechaInicio := time.Date(2026, 9, 20, 0, 0, 0, 0, time.UTC)
	fechaFin := time.Date(2026, 9, 25, 0, 0, 0, 0, time.UTC)

	reserva, err := testQueries.CreateReserva(context.Background(), db.CreateReservaParams{
		IDReserva:     200,
		FechaInicio:   fechaInicio,
		IDDepto:       depto.IDDepto,
		FechaFin:      fechaFin,
		PrecioBase:    "60000.00",
		CantNoches:    5,
		Descuento:     sql.NullString{String: "5.00", Valid: true},
		Observaciones: sql.NullString{String: "Reserva para huesped", Valid: true},
	})
	if err != nil {
		t.Fatalf("no se pudo crear reserva para huesped: %v", err)
	}

	huesped, err := testQueries.CreateHuesped(context.Background(), db.CreateHuespedParams{
		IDHuesped:     1,
		IDReserva:     reserva.IDReserva,
		Nombre:        "Juan",
		Apellido:      "Perez",
		Telefono:      sql.NullString{String: "123456789", Valid: true},
		Email:         sql.NullString{String: "juan@test.com", Valid: true},
		Observaciones: sql.NullString{String: "Huesped de prueba", Valid: true},
	})
	if err != nil {
		t.Fatalf("no se pudo crear huesped: %v", err)
	}
	return huesped
}

func TestCreateHuesped(t *testing.T) {
	huesped := createTestHuesped(t)

	if huesped.IDHuesped != 1 {
		t.Errorf("id_huesped incorrecto, esperado 1, got %d", huesped.IDHuesped)
	}
	if huesped.Nombre != "Juan" {
		t.Errorf("nombre incorrecto, esperado 'Juan', got '%s'", huesped.Nombre)
	}
	if huesped.Apellido != "Perez" {
		t.Errorf("apellido incorrecto, esperado 'Perez', got '%s'", huesped.Apellido)
	}
	if huesped.Email.String != "juan@test.com" {
		t.Errorf("email incorrecto, esperado 'juan@test.com', got '%s'", huesped.Email.String)
	}

	// limpiar
	testQueries.DeleteHuesped(context.Background(), huesped.IDHuesped)
	testQueries.DeleteReserva(context.Background(), huesped.IDReserva)
	testQueries.DeleteDepartamento(context.Background(), 200)
}

func TestGetHuesped(t *testing.T) {
	huesped := createTestHuesped(t)

	got, err := testQueries.GetHuesped(context.Background(), huesped.IDHuesped)
	if err != nil {
		t.Fatalf("no se pudo obtener huesped: %v", err)
	}
	if got.IDHuesped != huesped.IDHuesped {
		t.Errorf("id_huesped incorrecto, esperado %d, got %d", huesped.IDHuesped, got.IDHuesped)
	}
	if got.Nombre != huesped.Nombre {
		t.Errorf("nombre incorrecto, esperado '%s', got '%s'", huesped.Nombre, got.Nombre)
	}

	// limpiar
	testQueries.DeleteHuesped(context.Background(), huesped.IDHuesped)
	testQueries.DeleteReserva(context.Background(), huesped.IDReserva)
	testQueries.DeleteDepartamento(context.Background(), 200)
}

func TestListHuespedes(t *testing.T) {
	huesped1 := createTestHuesped(t)
	huesped2 := createTestHuesped(t)

	huespedes, err := testQueries.ListHuespedes(context.Background())
	if err != nil {
		t.Fatalf("no se pudo listar huespedes: %v", err)
	}
	if len(huespedes) < 2 {
		t.Errorf("esperaba al menos 2 huespedes, got %d", len(huespedes))
	}

	// limpiar
	testQueries.DeleteHuesped(context.Background(), huesped1.IDHuesped)
	testQueries.DeleteHuesped(context.Background(), huesped2.IDHuesped)
	testQueries.DeleteReserva(context.Background(), huesped1.IDReserva)
	testQueries.DeleteDepartamento(context.Background(), 200)
}

func TestUpdateHuesped(t *testing.T) {
	huesped := createTestHuesped(t)

	arg := db.UpdateHuespedParams{
		IDHuesped:     huesped.IDHuesped,
		IDReserva:     huesped.IDReserva,
		Nombre:        "Maria",
		Apellido:      "Gomez",
		Telefono:      sql.NullString{String: "987654321", Valid: true},
		Email:         sql.NullString{String: "maria@test.com", Valid: true},
		Observaciones: sql.NullString{String: "Huesped actualizado", Valid: true},
	}
	err := testQueries.UpdateHuesped(context.Background(), arg)
	if err != nil {
		t.Fatalf("no se pudo actualizar huesped: %v", err)
	}

	got, err := testQueries.GetHuesped(context.Background(), huesped.IDHuesped)
	if err != nil {
		t.Fatalf("no se pudo obtener huesped actualizado: %v", err)
	}
	if got.Nombre != "Maria" {
		t.Errorf("nombre no se actualizo, esperado 'Maria', got '%s'", got.Nombre)
	}
	if got.Apellido != "Gomez" {
		t.Errorf("apellido no se actualizo, esperado 'Gomez', got '%s'", got.Apellido)
	}
	if got.Email.String != "maria@test.com" {
		t.Errorf("email no se actualizo, esperado 'maria@test.com', got '%s'", got.Email.String)
	}

	// limpiar
	testQueries.DeleteHuesped(context.Background(), huesped.IDHuesped)
	testQueries.DeleteReserva(context.Background(), huesped.IDReserva)
	testQueries.DeleteDepartamento(context.Background(), 200)
}

func TestDeleteHuesped(t *testing.T) {
	huesped := createTestHuesped(t)

	err := testQueries.DeleteHuesped(context.Background(), huesped.IDHuesped)
	if err != nil {
		t.Fatalf("no se pudo eliminar huesped: %v", err)
	}

	_, err = testQueries.GetHuesped(context.Background(), huesped.IDHuesped)
	if err == nil {
		t.Error("esperaba error al obtener huesped eliminado, pero no hubo error")
	}

	// limpiar
	testQueries.DeleteReserva(context.Background(), huesped.IDReserva)
	testQueries.DeleteDepartamento(context.Background(), 200)
}
