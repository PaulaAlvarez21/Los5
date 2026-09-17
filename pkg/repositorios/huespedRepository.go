package repositorios

import (
	"context"

	sqlc "Los5/db/sqlc"
)

type HuespedRepository struct {
	query *sqlc.Queries
}

func NewHuespedRepository(queries *sqlc.Queries) *HuespedRepository {
	return &HuespedRepository{query: queries}
}

func (r *HuespedRepository) CreateHuesped(ctx context.Context, arg sqlc.CreateHuespedParams) (sqlc.Huesped, error) {
	return r.query.CreateHuesped(ctx, arg)
}

func (r *HuespedRepository) GetHuesped(ctx context.Context, id int32) (sqlc.Huesped, error) {
	return r.query.GetHuesped(ctx, id)
}

func (r *HuespedRepository) ListHuespedes(ctx context.Context) ([]sqlc.Huesped, error) {
	return r.query.ListHuespedes(ctx)
}

func (r *HuespedRepository) UpdateHuesped(ctx context.Context, arg sqlc.UpdateHuespedParams) error {
	return r.query.UpdateHuesped(ctx, arg)
}

func (r *HuespedRepository) DeleteHuesped(ctx context.Context, id int32) error {
	return r.query.DeleteHuesped(ctx, id)
}
