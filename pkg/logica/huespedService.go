package logica

import (
	"context"

	"Los5/pkg/dominio"
	"Los5/pkg/repositorios"
)

type HuespedService struct {
	huespedRepository *repositorios.HuespedRepository
}

func NewHuespedService(huespedRepository *repositorios.HuespedRepository) *HuespedService {
	return &HuespedService{huespedRepository: huespedRepository}
}

func (s *HuespedService) CrearHuesped(ctx context.Context, huesped dominio.Huesped) (dominio.Huesped, error) {
	if huesped.IDReserva <= 0 || huesped.Nombre == "" || huesped.Apellido == "" {
		return dominio.Huesped{}, dominio.ValidationError("datos de huesped inválidos")
	}
	return s.huespedRepository.CreateHuesped(ctx, huesped)
}

func (s *HuespedService) ObtenerHuespedes(ctx context.Context) ([]dominio.Huesped, error) {
	return s.huespedRepository.ListHuespedes(ctx)
}

func (s *HuespedService) ObtenerHuesped(ctx context.Context, id int32) (dominio.Huesped, error) {
	return s.huespedRepository.GetHuesped(ctx, id)
}

func (s *HuespedService) ActualizarHuesped(ctx context.Context, huesped dominio.Huesped) error {
	if huesped.IDHuesped <= 0 || huesped.IDReserva <= 0 || huesped.Nombre == "" || huesped.Apellido == "" {
		return dominio.ValidationError("datos de huesped inválidos")
	}
	return s.huespedRepository.UpdateHuesped(ctx, huesped)
}

func (s *HuespedService) EliminarHuesped(ctx context.Context, id int32) error {
	return s.huespedRepository.DeleteHuesped(ctx, id)
}