package logica

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"Los5/pkg/dominio"
	"Los5/pkg/repositorios"
)

type CombinadoService struct {
	repository *repositorios.BDRepository
}

func NewCombinadoService(repository *repositorios.BDRepository) *CombinadoService {
	return &CombinadoService{repository: repository}
}

type FiltroReservas struct {
	IDDepto *int32
	Desde   *time.Time
	Hasta   *time.Time
}

func (s *CombinadoService) ObtenerDepartamento(ctx context.Context, idDepto int32) (dominio.Departamento, error) {
	if idDepto <= 0 {
		return dominio.Departamento{}, dominio.ValidationError("id de departamento inválido")
	}
	return s.repository.GetDepartamento(ctx, idDepto)
}

func (s *CombinadoService) ObtenerPaginaReservas(ctx context.Context, filtro FiltroReservas) (dominio.PaginaReservas, error) {
	if filtro.IDDepto != nil && *filtro.IDDepto <= 0 {
		return dominio.PaginaReservas{}, dominio.ValidationError("id de departamento inválido")
	}
	if filtro.Desde != nil && filtro.Hasta != nil && !filtro.Desde.Before(*filtro.Hasta) {
		return dominio.PaginaReservas{}, dominio.ValidationError("rango de fechas inválido")
	}

	var reservas []dominio.Reserva
	var huespedes []dominio.Huesped

	if filtro.IDDepto != nil {
		if _, err := s.repository.GetDepartamento(ctx, *filtro.IDDepto); err != nil {
			return dominio.PaginaReservas{}, err
		}

		reservasFiltradas, err := s.repository.ListReservasByDepto(ctx, *filtro.IDDepto)
		if err != nil {
			return dominio.PaginaReservas{}, err
		}
		reservas = reservasFiltradas

		huespedesFiltrados, err := s.repository.ListHuespedesByDepto(ctx, *filtro.IDDepto)
		if err != nil {
			return dominio.PaginaReservas{}, err
		}
		huespedes = huespedesFiltrados
	} else {
		reservasTodas, err := s.repository.ListReservas(ctx)
		if err != nil {
			return dominio.PaginaReservas{}, err
		}
		reservas = reservasTodas

		huespedesTodos, err := s.repository.ListHuespedes(ctx)
		if err != nil {
			return dominio.PaginaReservas{}, err
		}
		huespedes = huespedesTodos
	}

	departamentos, err := s.repository.ListDepartamentos(ctx)
	if err != nil {
		return dominio.PaginaReservas{}, err
	}

	nombresDeptos := make(map[int32]string, len(departamentos))
	for _, departamento := range departamentos {
		nombresDeptos[departamento.IDDepto] = departamento.Nombre
	}

	nombresHuespedes := make(map[int32][]string)
	for _, huesped := range huespedes {
		nombreCompleto := strings.TrimSpace(huesped.Nombre + " " + huesped.Apellido)
		nombresHuespedes[huesped.IDReserva] = append(nombresHuespedes[huesped.IDReserva], nombreCompleto)
	}

	reservasConDetalle := make([]dominio.ReservaConDetalle, 0, len(reservas))
	for _, reserva := range reservas {
		if !reservaEnRango(reserva, filtro) {
			continue
		}

		precioBase, err := parseDecimal(reserva.PrecioBase)
		if err != nil || math.IsNaN(precioBase) || math.IsInf(precioBase, 0) {
			return dominio.PaginaReservas{}, fmt.Errorf("precio base inválido: %q", reserva.PrecioBase)
		}

		nombreDepto := nombresDeptos[reserva.IDDepto]
		if nombreDepto == "" {
			nombreDepto = fmt.Sprintf("Departamento #%d", reserva.IDDepto)
		}

		nombreHuesped := "Sin huéspedes"
		if nombres := nombresHuespedes[reserva.IDReserva]; len(nombres) > 0 {
			nombreHuesped = strings.Join(nombres, ", ")
		}

		reservasConDetalle = append(reservasConDetalle, dominio.ReservaConDetalle{
			IDReserva:     reserva.IDReserva,
			NombreDepto:   nombreDepto,
			NombreHuesped: nombreHuesped,
			FechaInicio:   reserva.FechaInicio,
			FechaFin:      reserva.FechaFin,
			CantNoches:    reserva.CantNoches,
			PrecioBase:    precioBase,
		})
	}

	return dominio.PaginaReservas{
		Reservas:      reservasConDetalle,
		Departamentos: departamentos,
		Huespedes:     huespedes,
	}, nil
}

func reservaEnRango(reserva dominio.Reserva, filtro FiltroReservas) bool {
	if filtro.Desde != nil && !reserva.FechaFin.After(*filtro.Desde) {
		return false
	}
	if filtro.Hasta != nil && !reserva.FechaInicio.Before(*filtro.Hasta) {
		return false
	}
	return true
}

func (s *CombinadoService) ObtenerReservasPorDepartamento(ctx context.Context, idDepto int32) (dominio.ReservasPorDepartamento, error) {
	if idDepto <= 0 {
		return dominio.ReservasPorDepartamento{}, dominio.ValidationError("id de departamento inválido")
	}

	departamento, err := s.repository.GetDepartamento(ctx, idDepto)
	if err != nil {
		return dominio.ReservasPorDepartamento{}, err
	}

	reservas, err := s.repository.ListReservasByDepto(ctx, idDepto)
	if err != nil {
		return dominio.ReservasPorDepartamento{}, err
	}

	return dominio.ReservasPorDepartamento{
		Departamento: departamento,
		Reservas:     reservas,
	}, nil
}

func (s *CombinadoService) ObtenerHuespedesPorDepartamento(ctx context.Context, idDepto int32) (dominio.HuespedesPorDepartamento, error) {
	if idDepto <= 0 {
		return dominio.HuespedesPorDepartamento{}, dominio.ValidationError("id de departamento inválido")
	}

	departamento, err := s.repository.GetDepartamento(ctx, idDepto)
	if err != nil {
		return dominio.HuespedesPorDepartamento{}, err
	}

	huespedes, err := s.repository.ListHuespedesByDepto(ctx, idDepto)
	if err != nil {
		return dominio.HuespedesPorDepartamento{}, err
	}

	return dominio.HuespedesPorDepartamento{
		Departamento: departamento,
		Huespedes:    huespedes,
	}, nil
}

func (s *CombinadoService) ConsultarDisponibilidad(ctx context.Context, idDepto int32, desde, hasta time.Time) (dominio.Disponibilidad, error) {
	if idDepto <= 0 {
		return dominio.Disponibilidad{}, dominio.ValidationError("id de departamento inválido")
	}
	if desde.IsZero() || hasta.IsZero() || !desde.Before(hasta) {
		return dominio.Disponibilidad{}, dominio.ValidationError("rango de fechas inválido")
	}

	departamento, err := s.repository.GetDepartamento(ctx, idDepto)
	if err != nil {
		return dominio.Disponibilidad{}, err
	}

	cantReservasSolapadas, err := s.repository.CountReservasSolapadas(ctx, idDepto, desde, hasta)
	if err != nil {
		return dominio.Disponibilidad{}, err
	}

	return dominio.Disponibilidad{
		Departamento:          departamento,
		Desde:                 desde,
		Hasta:                 hasta,
		Disponible:            departamento.Disponible && cantReservasSolapadas == 0,
		CantReservasSolapadas: cantReservasSolapadas,
	}, nil
}

func (s *CombinadoService) CalcularTotalReserva(ctx context.Context, idReserva int32) (dominio.TotalReserva, error) {
	if idReserva <= 0 {
		return dominio.TotalReserva{}, dominio.ValidationError("id de reserva inválido")
	}

	reserva, err := s.repository.GetReserva(ctx, idReserva)
	if err != nil {
		return dominio.TotalReserva{}, err
	}
	if reserva.CantNoches <= 0 {
		return dominio.TotalReserva{}, dominio.ValidationError("cantidad de noches inválida")
	}

	precioBase, err := parseDecimal(reserva.PrecioBase)
	if err != nil || precioBase < 0 || math.IsNaN(precioBase) || math.IsInf(precioBase, 0) {
		return dominio.TotalReserva{}, fmt.Errorf("precio base inválido: %q", reserva.PrecioBase)
	}

	descuento := 0.0
	if reserva.Descuento != nil && strings.TrimSpace(*reserva.Descuento) != "" {
		descuento, err = parseDecimal(*reserva.Descuento)
		if err != nil || descuento < 0 || descuento > 100 || math.IsNaN(descuento) || math.IsInf(descuento, 0) {
			return dominio.TotalReserva{}, dominio.ValidationError("descuento inválido")
		}
	}

	subtotal := precioBase * float64(reserva.CantNoches)
	total := subtotal * (1 - descuento/100)

	return dominio.TotalReserva{
		Reserva:             reserva,
		Subtotal:            formatMoney(subtotal),
		DescuentoPorcentaje: formatMoney(descuento),
		Total:               formatMoney(total),
	}, nil
}

func parseDecimal(valor string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(valor), 64)
}

func formatMoney(valor float64) string {
	return strconv.FormatFloat(math.Round(valor*100)/100, 'f', 2, 64)
}
