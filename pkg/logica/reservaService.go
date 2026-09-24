package logica

import (
	"context"

	"Los5/pkg/dominio"
	"Los5/pkg/repositorios"
)

//se validan crear reserva y actualizar reserva, porque son los datos ingresados por usuario que pueden no respetar reglas de negocio.

type ReservaService struct {
	reservaRepository *repositorios.ReservasRepository
}

func NewReservaService(reservaRepository *repositorios.ReservasRepository) *ReservaService {
	return &ReservaService{reservaRepository: reservaRepository}
}

func (s *ReservaService) CrearReserva(ctx context.Context, reserva dominio.Reserva) (dominio.Reserva, error) {
	if reserva.IDDepto <= 0 || reserva.CantNoches <= 0 || !reserva.FechaInicio.Before(reserva.FechaFin) {
		return dominio.Reserva{}, dominio.ValidationError("datos de reserva inválidos")
	}
	return s.reservaRepository.CreateReserva(ctx, reserva)
}

func (s *ReservaService) ObtenerReservas(ctx context.Context) ([]dominio.Reserva, error) {
	return s.reservaRepository.ListReservas(ctx)
}

func (s *ReservaService) ObtenerReserva(ctx context.Context, id int32) (dominio.Reserva, error) {
	return s.reservaRepository.GetReserva(ctx, id)
}

func (s *ReservaService) ActualizarReserva(ctx context.Context, reserva dominio.Reserva) error {
	if reserva.IDReserva <= 0 || reserva.IDDepto <= 0 || reserva.CantNoches <= 0 || !reserva.FechaInicio.Before(reserva.FechaFin) {
		return dominio.ValidationError("datos de reserva inválidos")
	}
	return s.reservaRepository.UpdateReserva(ctx, reserva)
}

func (s *ReservaService) EliminarReserva(ctx context.Context, id int32) error {
	return s.reservaRepository.DeleteReserva(ctx, id)
}

//calcular precio base de reserva, con descuento si corresponde. Se puede agregar más adelante la lógica de precios según temporada, etc.

//cancelar reserva, según política de cancelación. Se puede agregar más adelante la lógica de cancelación y reembolso.
