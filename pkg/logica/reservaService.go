package logica

import (
	"Los5/pkg/repositorios"
	"context"
)

type ReservaService struct {
	reservaRepository *repositorios.ReservasRepository
}

func NewReservaService(reservaRepository *repositorios.ReservasRepository) *ReservaService {
	return &ReservaService{reservaRepository: reservaRepository}
}

func (s *ReservaService) CrearReserva(reserva sqlc.CreateReservaParams) (sqlc.Reserva, error) {
	return s.reservaRepository.CreateReserva(context.Background(), reserva)
}

func (s *ReservaService) ObtenerReservas() ([]sqlc.Reserva, error) {
	return s.reservaRepository.ListReservas(context.Background())
}

func (s *ReservaService) ObtenerReserva(id int32) (sqlc.Reserva, error)
