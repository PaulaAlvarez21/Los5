package logica

import (
	"context"

	"Los5/pkg/repositorios"
)

type ReservaService struct {
	reservaRepository *repositorios.ReservasRepository
}

//hay que crear el struct de request para crear reserva y actualizar reserva, porque no se puede usar el struct de sqlc porque tiene el id que no se necesita al crear y al actualizar se necesita el id pero no se necesita el id_depto ni la fecha_inicio ni la fecha_fin
/*
type ReservaActualizar struct {
    IDReserva     int32
    FechaInicio   time.Time
    IDDepto       int32
    FechaFin      time.Time
    PrecioBase    string
    CantNoches    int32
    Descuento     sql.NullString
    Observaciones sql.NullString
}*/
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
func (s *ReservaService) ActualizarReserva(reserva repositorios.UpdateReserva) error { //(reserva ReservaActualizar) error {
	return s.reservaRepository.UpdateReserva(context.Background(), reserva)
}

func (s *ReservaService) EliminarReserva(id int32) error {
	return s.reservaRepository.DeleteReserva(context.Background(), id)
}
