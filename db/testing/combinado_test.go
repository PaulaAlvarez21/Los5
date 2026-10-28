package db_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	sqlc "Los5/db/sqlc"
)

func crearDepartamentoCombinado(t *testing.T, nombre, direccion string) sqlc.Departamento {
	t.Helper()

	depto, err := testQueries.CreateDepartamento(context.Background(), sqlc.CreateDepartamentoParams{
		Nombre:      nombre,
		Direccion:   direccion,
		Disponible:  true,
		Limpio:      true,
		Descripcion: sql.NullString{String: "Departamento para test combinado", Valid: true},
	})
	if err != nil {
		t.Fatalf("no se pudo crear departamento combinado: %v", err)
	}
	t.Cleanup(func() { testQueries.DeleteDepartamento(context.Background(), depto.IDDepto) })
	return depto
}

func crearReservaCombinada(t *testing.T, idDepto int32, fechaInicio, fechaFin time.Time) sqlc.Reserva {
	t.Helper()

	reserva, err := testQueries.CreateReserva(context.Background(), sqlc.CreateReservaParams{
		FechaInicio:   fechaInicio,
		IDDepto:       idDepto,
		FechaFin:      fechaFin,
		PrecioBase:    "10000.00",
		CantNoches:    4,
		Descuento:     sql.NullString{String: "10.00", Valid: true},
		Observaciones: sql.NullString{String: "Reserva para test combinado", Valid: true},
	})
	if err != nil {
		t.Fatalf("no se pudo crear reserva combinada: %v", err)
	}
	t.Cleanup(func() { testQueries.DeleteReserva(context.Background(), reserva.IDReserva) })
	return reserva
}

func crearHuespedCombinado(t *testing.T, idReserva int32, nombre, apellido string) sqlc.Huesped {
	t.Helper()

	huesped, err := testQueries.CreateHuesped(context.Background(), sqlc.CreateHuespedParams{
		IDReserva:     idReserva,
		Nombre:        nombre,
		Apellido:      apellido,
		Telefono:      sql.NullString{String: "123456789", Valid: true},
		Email:         sql.NullString{String: "combinado@test.com", Valid: true},
		Observaciones: sql.NullString{String: "Huesped para test combinado", Valid: true},
	})
	if err != nil {
		t.Fatalf("no se pudo crear huesped combinado: %v", err)
	}
	t.Cleanup(func() { testQueries.DeleteHuesped(context.Background(), huesped.IDHuesped) })
	return huesped
}

func datosTestCombinado(t *testing.T) (sqlc.Departamento, sqlc.Reserva, sqlc.Huesped) {
	t.Helper()

	depto := crearDepartamentoCombinado(t, "Depto Combinado", "Calle Combinada 100")
	reserva := crearReservaCombinada(
		t,
		depto.IDDepto,
		time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 11, 5, 0, 0, 0, 0, time.UTC),
	)
	huesped := crearHuespedCombinado(t, reserva.IDReserva, "Ana", "Combinada")

	otroDepto := crearDepartamentoCombinado(t, "Otro Depto Combinado", "Calle Combinada 200")
	otraReserva := crearReservaCombinada(
		t,
		otroDepto.IDDepto,
		time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 12, 5, 0, 0, 0, 0, time.UTC),
	)
	crearHuespedCombinado(t, otraReserva.IDReserva, "Beto", "Ajeno")

	return depto, reserva, huesped
}

func TestListReservasByDepto(t *testing.T) {
	depto, reserva, _ := datosTestCombinado(t)

	reservas, err := testQueries.ListReservasByDepto(context.Background(), depto.IDDepto)
	if err != nil {
		t.Fatalf("no se pudo listar reservas por departamento: %v", err)
	}
	if len(reservas) != 1 {
		t.Fatalf("esperaba 1 reserva, got %d", len(reservas))
	}
	if reservas[0].IDReserva != reserva.IDReserva {
		t.Errorf("id_reserva incorrecto, esperado %d, got %d", reserva.IDReserva, reservas[0].IDReserva)
	}
}

func TestListHuespedesByDepto(t *testing.T) {
	depto, _, huesped := datosTestCombinado(t)

	huespedes, err := testQueries.ListHuespedesByDepto(context.Background(), depto.IDDepto)
	if err != nil {
		t.Fatalf("no se pudo listar huespedes por departamento: %v", err)
	}
	if len(huespedes) != 1 {
		t.Fatalf("esperaba 1 huesped, got %d", len(huespedes))
	}
	if huespedes[0].IDHuesped != huesped.IDHuesped {
		t.Errorf("id_huesped incorrecto, esperado %d, got %d", huesped.IDHuesped, huespedes[0].IDHuesped)
	}
}

func TestCountReservasSolapadas(t *testing.T) {
	depto, _, _ := datosTestCombinado(t)
	ctx := context.Background()

	casos := []struct {
		nombre string
		desde  time.Time
		hasta  time.Time
		espera int64
	}{
		{
			nombre: "se solapa con la reserva existente",
			desde:  time.Date(2026, 11, 2, 0, 0, 0, 0, time.UTC),
			hasta:  time.Date(2026, 11, 4, 0, 0, 0, 0, time.UTC),
			espera: 1,
		},
		{
			nombre: "no se solapa con la reserva existente",
			desde:  time.Date(2026, 12, 10, 0, 0, 0, 0, time.UTC),
			hasta:  time.Date(2026, 12, 15, 0, 0, 0, 0, time.UTC),
			espera: 0,
		},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			cant, err := testQueries.CountReservasSolapadas(ctx, sqlc.CountReservasSolapadasParams{
				IDDepto: depto.IDDepto,
				Desde:   caso.desde,
				Hasta:   caso.hasta,
			})
			if err != nil {
				t.Fatalf("no se pudo contar reservas solapadas: %v", err)
			}
			if cant != caso.espera {
				t.Errorf("esperaba %d reservas solapadas, got %d", caso.espera, cant)
			}
		})
	}
}
