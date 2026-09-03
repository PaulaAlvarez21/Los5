package db

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"
)

func datosTestHuesped(t *testing.T) Huesped {
	t.Helper()

	// crear departamento dependiente
	depto, err := testQueries.CreateDepartamento(context.Background(), CreateDepartamentoParams{
		Nombre:      "Depto Huesped",
		Direccion:   "Calle Huesped 200",
		Disponible:  true,
		Limpio:      true,
		Descripcion: sql.NullString{String: "Depto para test de huesped", Valid: true},
	})
	if err != nil {
		t.Fatalf("no se pudo crear departamento para huesped: %v", err)
	}
	t.Cleanup(func() { testQueries.DeleteDepartamento(context.Background(), depto.IDDepto) })

	// crear reserva dependiente
	fechaInicio := time.Date(2026, 9, 20, 0, 0, 0, 0, time.UTC)
	fechaFin := time.Date(2026, 9, 25, 0, 0, 0, 0, time.UTC)

	reserva, err := testQueries.CreateReserva(context.Background(), CreateReservaParams{
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
	t.Cleanup(func() { testQueries.DeleteReserva(context.Background(), reserva.IDReserva) })

	huesped, err := testQueries.CreateHuesped(context.Background(), CreateHuespedParams{
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
	t.Cleanup(func() { testQueries.DeleteHuesped(context.Background(), huesped.IDHuesped) })
	return huesped
}

func TestCreateHuesped(t *testing.T) {
	huesped := datosTestHuesped(t)

	if huesped.IDHuesped == 0 {
		t.Error("id_huesped deberia ser generado por la base de datos")
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
}

func TestGetHuesped(t *testing.T) {
	huesped := datosTestHuesped(t)

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
}

func TestListHuespedes(t *testing.T) {
	datosTestHuesped(t)
	datosTestHuesped(t)

	huespedes, err := testQueries.ListHuespedes(context.Background())
	if err != nil {
		t.Fatalf("no se pudo listar huespedes: %v", err)
	}
	if len(huespedes) < 2 {
		t.Errorf("esperaba al menos 2 huespedes, got %d", len(huespedes))
	}
}

func TestUpdateHuesped(t *testing.T) {
	huesped := datosTestHuesped(t)

	arg := UpdateHuespedParams{
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
}

func TestDeleteHuesped(t *testing.T) {
	huesped := datosTestHuesped(t)

	err := testQueries.DeleteHuesped(context.Background(), huesped.IDHuesped)
	if err != nil {
		t.Fatalf("no se pudo eliminar huesped: %v", err)
	}

	_, err = testQueries.GetHuesped(context.Background(), huesped.IDHuesped)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("esperaba sql.ErrNoRows al obtener huesped eliminado, got %v", err)
	}
}
