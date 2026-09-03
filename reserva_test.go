package main

import (
	"context"
	"database/sql"
	"testing"
	"time"

	db "Los5.com/ServidorWeb/db/sqlc"
)

func createTestReserva(t *testing.T) db.Reserva {
	t.Helper()

	// crear departamento dependiente primero
	depto, err := testQueries.CreateDepartamento(context.Background(), db.CreateDepartamentoParams{
		IDDepto:     100,
		Nombre:      "Depto Reserva",
		Direccion:   "Av. Principal 100",
		Disponible:  true,
		Limpio:      true,
		Descripcion: sql.NullString{String: "Depto para test de reserva", Valid: true},
	})
	if err != nil {
		t.Fatalf("no se pudo crear departamento para reserva: %v", err)
	}

	fechaInicio := time.Date(2026, 9, 10, 0, 0, 0, 0, time.UTC)
	fechaFin := time.Date(2026, 9, 15, 0, 0, 0, 0, time.UTC)

	reserva, err := testQueries.CreateReserva(context.Background(), db.CreateReservaParams{
		IDReserva:     1,
		FechaInicio:   fechaInicio,
		IDDepto:       depto.IDDepto,
		FechaFin:      fechaFin,
		PrecioBase:    "50000.00",
		CantNoches:    5,
		Descuento:     sql.NullString{String: "10.00", Valid: true},
		Observaciones: sql.NullString{String: "Reserva de prueba", Valid: true},
	})
	if err != nil {
		t.Fatalf("no se pudo crear reserva: %v", err)
	}
	return reserva
}

func TestCreateReserva(t *testing.T) {
	reserva := createTestReserva(t)

	if reserva.IDReserva != 1 {
		t.Errorf("id_reserva incorrecto, esperado 1, got %d", reserva.IDReserva)
	}
	if reserva.IDDepto != 100 {
		t.Errorf("id_depto incorrecto, esperado 100, got %d", reserva.IDDepto)
	}
	if reserva.CantNoches != 5 {
		t.Errorf("cant_noches incorrecto, esperado 5, got %d", reserva.CantNoches)
	}
	if reserva.PrecioBase != "50000.00" {
		t.Errorf("precio_base incorrecto, esperado '50000.00', got '%s'", reserva.PrecioBase)
	}

	// limpiar
	testQueries.DeleteReserva(context.Background(), reserva.IDReserva)
	testQueries.DeleteDepartamento(context.Background(), reserva.IDDepto)
}

func TestGetReserva(t *testing.T) {
	reserva := createTestReserva(t)

	got, err := testQueries.GetReserva(context.Background(), reserva.IDReserva)
	if err != nil {
		t.Fatalf("no se pudo obtener reserva: %v", err)
	}
	if got.IDReserva != reserva.IDReserva {
		t.Errorf("id_reserva incorrecto, esperado %d, got %d", reserva.IDReserva, got.IDReserva)
	}
	if got.FechaInicio != reserva.FechaInicio {
		t.Errorf("fecha_inicio incorrecta, esperado %v, got %v", reserva.FechaInicio, got.FechaInicio)
	}

	// limpiar
	testQueries.DeleteReserva(context.Background(), reserva.IDReserva)
	testQueries.DeleteDepartamento(context.Background(), reserva.IDDepto)
}

func TestListReservas(t *testing.T) {
	reserva1 := createTestReserva(t)
	reserva2 := createTestReserva(t)

	reservas, err := testQueries.ListReservas(context.Background())
	if err != nil {
		t.Fatalf("no se pudo listar reservas: %v", err)
	}
	if len(reservas) < 2 {
		t.Errorf("esperaba al menos 2 reservas, got %d", len(reservas))
	}

	// limpiar
	testQueries.DeleteReserva(context.Background(), reserva1.IDReserva)
	testQueries.DeleteReserva(context.Background(), reserva2.IDReserva)
	testQueries.DeleteDepartamento(context.Background(), reserva1.IDDepto)
}

func TestUpdateReserva(t *testing.T) {
	reserva := createTestReserva(t)

	nuevaFecha := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	nuevaFechaFin := time.Date(2026, 10, 5, 0, 0, 0, 0, time.UTC)

	arg := db.UpdateReservaParams{
		IDReserva:     reserva.IDReserva,
		FechaInicio:   nuevaFecha,
		IDDepto:       reserva.IDDepto,
		FechaFin:      nuevaFechaFin,
		PrecioBase:    "75000.00",
		CantNoches:    4,
		Descuento:     sql.NullString{String: "15.00", Valid: true},
		Observaciones: sql.NullString{String: "Reserva actualizada", Valid: true},
	}
	err := testQueries.UpdateReserva(context.Background(), arg)
	if err != nil {
		t.Fatalf("no se pudo actualizar reserva: %v", err)
	}

	got, err := testQueries.GetReserva(context.Background(), reserva.IDReserva)
	if err != nil {
		t.Fatalf("no se pudo obtener reserva actualizada: %v", err)
	}
	if got.PrecioBase != "75000.00" {
		t.Errorf("precio_base no se actualizo, esperado '75000.00', got '%s'", got.PrecioBase)
	}
	if got.CantNoches != 4 {
		t.Errorf("cant_noches no se actualizo, esperado 4, got %d", got.CantNoches)
	}

	// limpiar
	testQueries.DeleteReserva(context.Background(), reserva.IDReserva)
	testQueries.DeleteDepartamento(context.Background(), reserva.IDDepto)
}

func TestDeleteReserva(t *testing.T) {
	reserva := createTestReserva(t)

	err := testQueries.DeleteReserva(context.Background(), reserva.IDReserva)
	if err != nil {
		t.Fatalf("no se pudo eliminar reserva: %v", err)
	}

	_, err = testQueries.GetReserva(context.Background(), reserva.IDReserva)
	if err == nil {
		t.Error("esperaba error al obtener reserva eliminada, pero no hubo error")
	}

	// limpiar
	testQueries.DeleteDepartamento(context.Background(), reserva.IDDepto)
}
