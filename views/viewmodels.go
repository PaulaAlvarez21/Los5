package views

import "Los5/pkg/dominio"

// FiltroReservas conserva los valores ingresados en el formulario de filtros.
type FiltroReservas struct {
	DeptoID string
	Desde   string
	Hasta   string
}

type DepartamentosPageData struct {
	Departamentos []dominio.Departamento
	Error         string
}

type ReservasPageData struct {
	Reservas      []dominio.ReservaConDetalle
	Departamentos []dominio.Departamento
	Huespedes     []dominio.Huesped
	Filtro        FiltroReservas
	Error         string
}

type HuespedesPageData struct {
	Huespedes []dominio.Huesped
	Reservas  []dominio.Reserva
	Error     string
}

type ReservasPorDepartamentoPageData struct {
	Departamento dominio.Departamento
	Reservas     []dominio.Reserva
	Error        string
}

type HuespedesPorDepartamentoPageData struct {
	Departamento dominio.Departamento
	Huespedes    []dominio.Huesped
	Error        string
}
