package repositorios

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sqlc "Los5/db/sqlc"
	"Los5/pkg/dominio"
)

type BDRepository struct {
	query *sqlc.Queries
}

func NewBDRepository(queries *sqlc.Queries) *BDRepository {
	return &BDRepository{query: queries}
}

//Departamento funciones Repository

func (r *BDRepository) CreateDepartamento(ctx context.Context, departamento dominio.Departamento) (dominio.Departamento, error) {
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

func (r *BDRepository) GetDepartamento(ctx context.Context, id int32) (dominio.Departamento, error) {
	departamentoDB, err := r.query.GetDepartamento(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { //traduce error de sql a error de dominio
			return dominio.Departamento{}, dominio.NotFoundError()
		}
		return dominio.Departamento{}, err
	}
	return toDominioDepartamento(departamentoDB), nil
}

func (r *BDRepository) ListDepartamentos(ctx context.Context) ([]dominio.Departamento, error) {
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

func (r *BDRepository) UpdateDepartamento(ctx context.Context, departamento dominio.Departamento) error {
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

func (r *BDRepository) DeleteDepartamento(ctx context.Context, id int32) error {
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

// Huesped funciones Repository
func (r *BDRepository) CreateHuesped(ctx context.Context, huesped dominio.Huesped) (dominio.Huesped, error) {
	arg := sqlc.CreateHuespedParams{
		IDReserva:     toNullInt32(huesped.IDReserva),
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

func (r *BDRepository) GetHuesped(ctx context.Context, id int32) (dominio.Huesped, error) {
	huespedDB, err := r.query.GetHuesped(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { //traduce error de sql a error de dominio
			return dominio.Huesped{}, dominio.NotFoundError()
		}
		return dominio.Huesped{}, err
	}
	return toDominioHuesped(huespedDB), nil
}

func (r *BDRepository) ListHuespedes(ctx context.Context) ([]dominio.Huesped, error) {
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

func (r *BDRepository) UpdateHuesped(ctx context.Context, huesped dominio.Huesped) error {
	arg := sqlc.UpdateHuespedParams{
		IDHuesped:     huesped.IDHuesped,
		IDReserva:     toNullInt32(huesped.IDReserva),
		Nombre:        huesped.Nombre,
		Apellido:      huesped.Apellido,
		Telefono:      toNullString(huesped.Telefono),
		Email:         toNullString(huesped.Email),
		Observaciones: toNullString(huesped.Observaciones),
	}
	return r.query.UpdateHuesped(ctx, arg)
}

func (r *BDRepository) DeleteHuesped(ctx context.Context, id int32) error {
	return r.query.DeleteHuesped(ctx, id) //aca no hacemos tartamiento de error not found 404 de si no existe, poruqe si quiere borrar algo que no existe, no es un error, simplemente no hace nada.
}

// conversiones de tipos
func toDominioHuesped(huespedDB sqlc.Huesped) dominio.Huesped {
	return dominio.Huesped{
		IDHuesped:     huespedDB.IDHuesped,
		IDReserva:     fromNullInt32(huespedDB.IDReserva),
		Nombre:        huespedDB.Nombre,
		Apellido:      huespedDB.Apellido,
		Telefono:      fromNullString(huespedDB.Telefono),
		Email:         fromNullString(huespedDB.Email),
		Observaciones: fromNullString(huespedDB.Observaciones),
	}
}

// Reservas funciones Repository
func (r *BDRepository) CreateReserva(ctx context.Context, reserva dominio.Reserva) (dominio.Reserva, error) {
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

func (r *BDRepository) GetReserva(ctx context.Context, id int32) (dominio.Reserva, error) {
	reservaDB, err := r.query.GetReserva(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { //traduce error de sql a error de dominio
			return dominio.Reserva{}, dominio.NotFoundError()
		}
		return dominio.Reserva{}, err
	}
	return toDominioReserva(reservaDB), nil
}

func (r *BDRepository) ListReservas(ctx context.Context) ([]dominio.Reserva, error) {
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

func (r *BDRepository) UpdateReserva(ctx context.Context, reserva dominio.Reserva) error {
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

func (r *BDRepository) DeleteReserva(ctx context.Context, id int32) error {
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

func toNullInt32(n *int32) sql.NullInt32 {
	if n == nil {
		return sql.NullInt32{}
	}
	return sql.NullInt32{Int32: *n, Valid: true}
}

func fromNullInt32(n sql.NullInt32) *int32 {
	if !n.Valid {
		return nil
	}
	return &n.Int32
}

//REPOSITORY COMBINADO CASOS DE USO

func (r *BDRepository) ListReservasByDepto(ctx context.Context, idDepto int32) ([]dominio.Reserva, error) {
	reservasDB, err := r.query.ListReservasByDepto(ctx, idDepto)
	if err != nil {
		return nil, err
	}
	reservas := make([]dominio.Reserva, 0, len(reservasDB))
	for _, reservaDB := range reservasDB {
		reservas = append(reservas, toDominioReserva(reservaDB))
	}
	return reservas, nil
}

func (r *BDRepository) ListHuespedesByDepto(ctx context.Context, idDepto int32) ([]dominio.Huesped, error) {
	huespedesDB, err := r.query.ListHuespedesByDepto(ctx, idDepto)
	if err != nil {
		return nil, err
	}
	huespedes := make([]dominio.Huesped, 0, len(huespedesDB))
	for _, huespedDB := range huespedesDB {
		huespedes = append(huespedes, toDominioHuesped(huespedDB))
	}
	return huespedes, nil
}

func (r *BDRepository) CountReservasSolapadas(ctx context.Context, idDepto int32, desde, hasta time.Time) (int64, error) {
	arg := sqlc.CountReservasSolapadasParams{
		IDDepto: idDepto,
		Hasta:   hasta,
		Desde:   desde,
	}
	return r.query.CountReservasSolapadas(ctx, arg)
}

//FALTA CORRER SQLC DE NUEVO CON NUEVAS QUERIES
