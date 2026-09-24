package repositorios

import (
	"context"
	"database/sql"
	"errors"

	sqlc "Los5/db/sqlc"
	"Los5/pkg/dominio"
)

type ReservasRepository struct {
	query *sqlc.Queries
}

func NewReservasRepository(queries *sqlc.Queries) *ReservasRepository {
	return &ReservasRepository{query: queries}
}

func (r *ReservasRepository) CreateReserva(ctx context.Context, reserva dominio.Reserva) (dominio.Reserva, error) {
	arg := sqlc.CreateReservaParams{
		FechaInicio:   reserva.FechaInicio,
		IDDepto:       reserva.IDDepto,
		FechaFin:      reserva.FechaFin,
		PrecioBase:    reserva.PrecioBase,
		CantNoches:    reserva.CantNoches,
		Descuento:     toNullString(reserva.Descuento),
		Observaciones: toNullString(reserva.Observaciones), //hace la conversion de reservastruct a sql
	}
	reservaDB, err := r.query.CreateReserva(ctx, arg)
	if err != nil {
		return dominio.Reserva{}, err
	}
	return toDominioReserva(reservaDB), nil
}

func (r *ReservasRepository) GetReserva(ctx context.Context, id int32) (dominio.Reserva, error) {
	reservaDB, err := r.query.GetReserva(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { //traduce error de sql a error de dominio
			return dominio.Reserva{}, dominio.NotFoundError()
		}
		return dominio.Reserva{}, err
	}
	return toDominioReserva(reservaDB), nil
}

func (r *ReservasRepository) ListReservas(ctx context.Context) ([]dominio.Reserva, error) {
	reservasDB, err := r.query.ListReservas(ctx)
	if err != nil {
		return nil, err
	}
	reservas := make([]dominio.Reserva, 0, len(reservasDB))
	for _, reservaDB := range reservasDB {
		reservas = append(reservas, toDominioReserva(reservaDB))
	}
	return reservas, nil
}

func (r *ReservasRepository) UpdateReserva(ctx context.Context, reserva dominio.Reserva) error {
	arg := sqlc.UpdateReservaParams{
		IDReserva:     reserva.IDReserva,
		FechaInicio:   reserva.FechaInicio,
		IDDepto:       reserva.IDDepto,
		FechaFin:      reserva.FechaFin,
		PrecioBase:    reserva.PrecioBase,
		CantNoches:    reserva.CantNoches,
		Descuento:     toNullString(reserva.Descuento),
		Observaciones: toNullString(reserva.Observaciones),
	}
	return r.query.UpdateReserva(ctx, arg) //sqlc devuelve si modifico o no
}

func (r *ReservasRepository) DeleteReserva(ctx context.Context, id int32) error {
	return r.query.DeleteReserva(ctx, id) //aca no hacemos tartamiento de error not found 404 de si no existe, poruqe si quiere borrar algo que no existe, no es un error, simplemente no hace nada.
}

// conversiones de tipos
func toDominioReserva(reservaDB sqlc.Reserva) dominio.Reserva {
	return dominio.Reserva{
		IDReserva:     reservaDB.IDReserva,
		FechaInicio:   reservaDB.FechaInicio,
		IDDepto:       reservaDB.IDDepto,
		FechaFin:      reservaDB.FechaFin,
		PrecioBase:    reservaDB.PrecioBase,
		CantNoches:    reservaDB.CantNoches,
		Descuento:     fromNullString(reservaDB.Descuento),
		Observaciones: fromNullString(reservaDB.Observaciones),
	}
}

// metodos generales usados en reservarepository, huespedesrepository y departamentosrepository, para convertir de sql.NullString a *string y viceversa, para poder usar punteros en los structs de dominio y que sean nulos en la base de datos.
func toNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *s, Valid: true}
}

func fromNullString(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}
