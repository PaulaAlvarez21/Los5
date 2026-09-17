package repositorios

import (
	"context"

	sqlc "Los5/db/sqlc"
)

type ReservasRepository struct {
	query *sqlc.Queries
}

func NewReservasRepository(queries *sqlc.Queries) *ReservasRepository {
	return &ReservasRepository{query: queries}
}

func (r *ReservasRepository) CreateReserva(ctx context.Context, arg sqlc.CreateReservaParams) (sqlc.Reserva, error) {
	return r.query.CreateReserva(ctx, arg)
}

func (r *ReservasRepository) GetReserva(ctx context.Context, id int32) (sqlc.Reserva, error) {
	return r.query.GetReserva(ctx, id)
}

func (r *ReservasRepository) ListReservas(ctx context.Context) ([]sqlc.Reserva, error) {
	return r.query.ListReservas(ctx)
}

func (r *ReservasRepository) UpdateReserva(ctx context.Context, arg sqlc.UpdateReservaParams) error {
	return r.query.UpdateReserva(ctx, arg)
}

func (r *ReservasRepository) DeleteReserva(ctx context.Context, id int32) error {
	return r.query.DeleteReserva(ctx, id)
}