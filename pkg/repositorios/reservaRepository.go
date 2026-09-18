package repositorios

import (
	"context"

	sqlc "Los5/db/sqlc"
)

type ReservasRepository struct {
	query *sqlc.Queries
}

// no se si esto esta bien, es para que esto funcione en el service sin poner sqlc.Queries en el service
type CreateReserva = sqlc.CreateReservaParams
type Reserva = sqlc.Reserva
type UpdateReserva = sqlc.UpdateReservaParams

func NewReservasRepository(queries *sqlc.Queries) *ReservasRepository {
	return &ReservasRepository{query: queries}
}

func (r *ReservasRepository) CreateReserva(ctx context.Context, arg CreateReserva) (Reserva, error) {
	return r.query.CreateReserva(ctx, arg)
}

func (r *ReservasRepository) GetReserva(ctx context.Context, id int32) (Reserva, error) {
	return r.query.GetReserva(ctx, id)
}

func (r *ReservasRepository) ListReservas(ctx context.Context) ([]Reserva, error) {
	return r.query.ListReservas(ctx)
}

func (r *ReservasRepository) UpdateReserva(ctx context.Context, arg UpdateReserva) error {
	return r.query.UpdateReserva(ctx, arg)
}

func (r *ReservasRepository) DeleteReserva(ctx context.Context, id int32) error {
	return r.query.DeleteReserva(ctx, id)
}
