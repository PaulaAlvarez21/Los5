package repositorios

import (
	"context"
	"database/sql"
	"errors"

	sqlc "Los5/db/sqlc"
	"Los5/pkg/dominio"
)

type DepartamentosRepository struct {
	query *sqlc.Queries
}

func NewDepartamentosRepository(queries *sqlc.Queries) *DepartamentosRepository {
	return &DepartamentosRepository{query: queries}
}

func (r *DepartamentosRepository) CreateDepartamento(ctx context.Context, departamento dominio.Departamento) (dominio.Departamento, error) {
	arg := sqlc.CreateDepartamentoParams{
		Nombre:      departamento.Nombre,
		Direccion:   departamento.Direccion,
		Disponible:  departamento.Disponible,
		Limpio:      departamento.Limpio,
		Descripcion: toNullString(departamento.Descripcion), //hace la conversion de departamentostruct a sql
	}
	departamentoDB, err := r.query.CreateDepartamento(ctx, arg)
	if err != nil {
		return dominio.Departamento{}, err
	}
	return toDominioDepartamento(departamentoDB), nil
}

func (r *DepartamentosRepository) GetDepartamento(ctx context.Context, id int32) (dominio.Departamento, error) {
	departamentoDB, err := r.query.GetDepartamento(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { //traduce error de sql a error de dominio
			return dominio.Departamento{}, dominio.NotFoundError()
		}
		return dominio.Departamento{}, err
	}
	return toDominioDepartamento(departamentoDB), nil
}

func (r *DepartamentosRepository) ListDepartamentos(ctx context.Context) ([]dominio.Departamento, error) {
	departamentosDB, err := r.query.ListDepartamentos(ctx)
	if err != nil {
		return nil, err
	}
	departamentos := make([]dominio.Departamento, 0, len(departamentosDB))
	for _, departamentoDB := range departamentosDB {
		departamentos = append(departamentos, toDominioDepartamento(departamentoDB))
	}
	return departamentos, nil
}

func (r *DepartamentosRepository) UpdateDepartamento(ctx context.Context, departamento dominio.Departamento) error {
	arg := sqlc.UpdateDepartamentoParams{
		IDDepto:     departamento.IDDepto,
		Nombre:      departamento.Nombre,
		Direccion:   departamento.Direccion,
		Disponible:  departamento.Disponible,
		Limpio:      departamento.Limpio,
		Descripcion: toNullString(departamento.Descripcion),
	}
	return r.query.UpdateDepartamento(ctx, arg)
}

func (r *DepartamentosRepository) DeleteDepartamento(ctx context.Context, id int32) error {
	return r.query.DeleteDepartamento(ctx, id) //aca no hacemos tartamiento de error not found 404 de si no existe, poruqe si quiere borrar algo que no existe, no es un error, simplemente no hace nada.
}

// conversiones de tipos
func toDominioDepartamento(departamentoDB sqlc.Departamento) dominio.Departamento {
	return dominio.Departamento{
		IDDepto:     departamentoDB.IDDepto,
		Nombre:      departamentoDB.Nombre,
		Direccion:   departamentoDB.Direccion,
		Disponible:  departamentoDB.Disponible,
		Limpio:      departamentoDB.Limpio,
		Descripcion: fromNullString(departamentoDB.Descripcion),
	}
}

//utiliza metodos generales definidos en reservaRepository.go para convertir de sql.NullString a *string y viceversa, para poder usar punteros en los structs de dominio y que sean nulos en la base de datos.
