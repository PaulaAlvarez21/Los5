package logica

import (
	"context"

	"Los5/pkg/dominio"
	"Los5/pkg/repositorios"
)

type ReservaService struct {
	reservaRepository *repositorios.ReservasRepository
}

func NewReservaService(reservaRepository *repositorios.ReservasRepository) *ReservaService {
	return &ReservaService{reservaRepository: reservaRepository}
}

func (s *ReservaService) CrearReserva(reserva dominio.Reserva) (dominio.Reserva, error) {
	return s.reservaRepository.CreateReserva(context.Background(), reserva)
}

func (s *ReservaService) ObtenerReservas() ([]dominio.Reserva, error) {
	return s.reservaRepository.ListReservas(context.Background())
}

func (s *ReservaService) ObtenerReserva(id int32) (dominio.Reserva, error) {
	return s.reservaRepository.GetReserva(context.Background(), id)
}

func (s *ReservaService) ActualizarReserva(reserva dominio.Reserva) error {
	return s.reservaRepository.UpdateReserva(context.Background(), reserva)
}

func (s *ReservaService) EliminarReserva(id int32) error {
	return s.reservaRepository.DeleteReserva(context.Background(), id)
}
