package dominio

import "time"

type ReservasPorDepartamento struct {
	Departamento Departamento
	Reservas     []Reserva
}

type HuespedesPorDepartamento struct {
	Departamento Departamento
	Huespedes    []Huesped
}

type Disponibilidad struct {
	Departamento          Departamento
	Desde                 time.Time
	Hasta                 time.Time
	Disponible            bool
	CantReservasSolapadas int64
}

type TotalReserva struct {
	Reserva             Reserva
	Subtotal            string
	DescuentoPorcentaje string
	Total               string
}

type ReservaConDetalle struct {
	IDReserva       int32
	NombreDepto     string
	NombreHuesped   string
	ApellidoHuesped string
	FechaInicio     time.Time
	FechaFin        time.Time
	CantNoches      int32
	PrecioBase      float64
}

type PaginaReservas struct {
	Reservas      []ReservaConDetalle
	Departamentos []Departamento
	Huespedes     []Huesped
}
