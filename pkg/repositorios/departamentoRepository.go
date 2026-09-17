package repositorios

import (
	"context"

	sqlc "Los5/db/sqlc"
)

type DepartamentosRepository struct {
	query *sqlc.Queries
}

func NewDepartamentosRepository(queries *sqlc.Queries) *DepartamentosRepository {
	return &DepartamentosRepository{query: queries}
}

func (r *DepartamentosRepository) CreateDepartamento(ctx context.Context, arg sqlc.CreateDepartamentoParams) (sqlc.Departamento, error) {
	return r.query.CreateDepartamento(ctx, arg)
}

func (r *DepartamentosRepository) GetDepartamento(ctx context.Context, id int32) (sqlc.Departamento, error) {
	return r.query.GetDepartamento(ctx, id)
}

func (r *DepartamentosRepository) ListDepartamentos(ctx context.Context) ([]sqlc.Departamento, error) {
	return r.query.ListDepartamentos(ctx)
}

func (r *DepartamentosRepository) UpdateDepartamento(ctx context.Context, arg sqlc.UpdateDepartamentoParams) error {
	return r.query.UpdateDepartamento(ctx, arg)
}

func (r *DepartamentosRepository) DeleteDepartamento(ctx context.Context, id int32) error {
	return r.query.DeleteDepartamento(ctx, id)
}