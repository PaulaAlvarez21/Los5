package main

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	db "Los5.com/ServidorWeb/db/sqlc"
)

func datosTestReserva(t *testing.T) db.Reserva {
	t.Helper()

	// crear departamento dependiente primero
	depto, err := testQueries.CreateDepartamento(context.Background(), db.CreateDepartamentoParams{
		Nombre:      "Depto Reserva",
		Direccion:   "Av. Principal 100",
		Disponible:  true,
		Limpio:      true,
		Descripcion: sql.NullString{String: "Depto para test de reserva", Valid: true},
	})
	if err != nil {
		t.Fatalf("no se pudo crear departamento para reserva: %v", err)
	}
	t.Cleanup(func() { testQueries.DeleteDepartamento(context.Background(), depto.IDDepto) })

	fechaInicio := time.Date(2026, 9, 10, 0, 0, 0, 0, time.UTC)
	fechaFin := time.Date(2026, 9, 15, 0, 0, 0, 0, time.UTC)

	reserva, err := testQueries.CreateReserva(context.Background(), db.CreateReservaParams{
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
	t.Cleanup(func() { testQueries.DeleteReserva(context.Background(), reserva.IDReserva) })
	return reserva
}

func TestCreateReserva(t *testing.T) {
	reserva := datosTestReserva(t)

	if reserva.IDReserva == 0 {
		t.Error("id_reserva deberia ser generado por la base de datos")
	}
	if reserva.CantNoches != 5 {
		t.Errorf("cant_noches incorrecto, esperado 5, got %d", reserva.CantNoches)
	}
	if reserva.PrecioBase != "50000.00" {
		t.Errorf("precio_base incorrecto, esperado '50000.00', got '%s'", reserva.PrecioBase)
	}
}

func TestGetReserva(t *testing.T) {
	reserva := datosTestReserva(t)

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
}

func TestListReservas(t *testing.T) {
	datosTestReserva(t)
	datosTestReserva(t)

	reservas, err := testQueries.ListReservas(context.Background())
	if err != nil {
		t.Fatalf("no se pudo listar reservas: %v", err)
	}
	if len(reservas) < 2 {
		t.Errorf("esperaba al menos 2 reservas, got %d", len(reservas))
	}
}

func TestUpdateReserva(t *testing.T) {
	reserva := datosTestReserva(t)

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
}

func TestDeleteReserva(t *testing.T) {
	reserva := datosTestReserva(t)

	err := testQueries.DeleteReserva(context.Background(), reserva.IDReserva)
	if err != nil {
		t.Fatalf("no se pudo eliminar reserva: %v", err)
	}

	_, err = testQueries.GetReserva(context.Background(), reserva.IDReserva)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("esperaba sql.ErrNoRows al obtener reserva eliminada, got %v", err)
	}
}
