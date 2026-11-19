package logica

import (
	"context"
	"strings"

	"Los5/pkg/dominio"
	"Los5/pkg/repositorios"
)

type HuespedService struct {
	huespedRepository *repositorios.BDRepository
}

func NewHuespedService(huespedRepository *repositorios.BDRepository) *HuespedService {
	return &HuespedService{huespedRepository: huespedRepository}
}

func (s *HuespedService) CrearHuesped(ctx context.Context, huesped dominio.Huesped) (dominio.Huesped, error) {
	if strings.TrimSpace(huesped.Nombre) == "" {
		return dominio.Huesped{}, dominio.ValidationError("datos de huesped inválidos")
	}
	if huesped.Telefono == nil || strings.TrimSpace(*huesped.Telefono) == "" {
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
	if huesped.IDHuesped <= 0 || strings.TrimSpace(huesped.Nombre) == "" {
		return dominio.ValidationError("datos de huesped inválidos")
	}
	if huesped.Telefono == nil || strings.TrimSpace(*huesped.Telefono) == "" {
		return dominio.ValidationError("datos de huesped inválidos")
	}
	return s.huespedRepository.UpdateHuesped(ctx, huesped)
}

func (s *HuespedService) EliminarHuesped(ctx context.Context, id int32) error {
	return s.huespedRepository.DeleteHuesped(ctx, id)
}
