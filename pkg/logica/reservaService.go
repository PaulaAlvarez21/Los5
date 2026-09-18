package logica

import (
	"context"

	"Los5/pkg/repositorios"
)

type ReservaService struct {
	reservaRepository *repositorios.ReservasRepository
}

// constructor
func NewReservaService(reservaRepository *repositorios.ReservasRepository) *ReservaService {
	return &ReservaService{reservaRepository: reservaRepository}
}

func (s *ReservaService) CrearReserva(reserva repositorios.CreateReserva) (repositorios.Reserva, error) {
	return s.reservaRepository.CreateReserva(context.Background(), reserva)
}

func (s *ReservaService) ObtenerReservas() ([]repositorios.Reserva, error) {
	return s.reservaRepository.ListReservas(context.Background())
}

func (s *ReservaService) ObtenerReserva(id int32) (repositorios.Reserva, error) {
	return s.reservaRepository.GetReserva(context.Background(), id)
}
