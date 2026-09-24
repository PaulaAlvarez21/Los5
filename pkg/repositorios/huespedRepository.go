package repositorios

import (
	"context"
	"database/sql"
	"errors"

	sqlc "Los5/db/sqlc"
	"Los5/pkg/dominio"
)

type HuespedRepository struct {
	query *sqlc.Queries
}

func NewHuespedRepository(queries *sqlc.Queries) *HuespedRepository {
	return &HuespedRepository{query: queries}
}

func (r *HuespedRepository) CreateHuesped(ctx context.Context, huesped dominio.Huesped) (dominio.Huesped, error) {
	arg := sqlc.CreateHuespedParams{
		IDReserva:     huesped.IDReserva,
		Nombre:        huesped.Nombre,
		Apellido:      huesped.Apellido,
		Telefono:      toNullString(huesped.Telefono),
		Email:         toNullString(huesped.Email),
		Observaciones: toNullString(huesped.Observaciones), //hace la conversion de huespedstruct a sql
	}
	huespedDB, err := r.query.CreateHuesped(ctx, arg)
	if err != nil {
		return dominio.Huesped{}, err
	}
	return toDominioHuesped(huespedDB), nil
}

func (r *HuespedRepository) GetHuesped(ctx context.Context, id int32) (dominio.Huesped, error) {
	huespedDB, err := r.query.GetHuesped(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { //traduce error de sql a error de dominio
			return dominio.Huesped{}, dominio.NotFoundError()
		}
		return dominio.Huesped{}, err
	}
	return toDominioHuesped(huespedDB), nil
}

func (r *HuespedRepository) ListHuespedes(ctx context.Context) ([]dominio.Huesped, error) {
	huespedesDB, err := r.query.ListHuespedes(ctx)
	if err != nil {
		return nil, err
	}
	huespedes := make([]dominio.Huesped, 0, len(huespedesDB))
	for _, huespedDB := range huespedesDB {
		huespedes = append(huespedes, toDominioHuesped(huespedDB))
	}
	return huespedes, nil
}

func (r *HuespedRepository) UpdateHuesped(ctx context.Context, huesped dominio.Huesped) error {
	arg := sqlc.UpdateHuespedParams{
		IDHuesped:     huesped.IDHuesped,
		IDReserva:     huesped.IDReserva,
		Nombre:        huesped.Nombre,
		Apellido:      huesped.Apellido,
		Telefono:      toNullString(huesped.Telefono),
		Email:         toNullString(huesped.Email),
		Observaciones: toNullString(huesped.Observaciones),
	}
	return r.query.UpdateHuesped(ctx, arg)
}

func (r *HuespedRepository) DeleteHuesped(ctx context.Context, id int32) error {
	return r.query.DeleteHuesped(ctx, id) //aca no hacemos tartamiento de error not found 404 de si no existe, poruqe si quiere borrar algo que no existe, no es un error, simplemente no hace nada.
}

// conversiones de tipos
func toDominioHuesped(huespedDB sqlc.Huesped) dominio.Huesped {
	return dominio.Huesped{
		IDHuesped:     huespedDB.IDHuesped,
		IDReserva:     huespedDB.IDReserva,
		Nombre:        huespedDB.Nombre,
		Apellido:      huespedDB.Apellido,
		Telefono:      fromNullString(huespedDB.Telefono),
		Email:         fromNullString(huespedDB.Email),
		Observaciones: fromNullString(huespedDB.Observaciones),
	}
}

//utiliza metodos generales definidos en reservaRepository.go para convertir de sql.NullString a *string y viceversa, para poder usar punteros en los structs de dominio y que sean nulos en la base de datos.
