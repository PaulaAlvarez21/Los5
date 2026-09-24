package logica

import (
	"context"

	"Los5/pkg/dominio"
	"Los5/pkg/repositorios"
)

type DepartamentoService struct {
	departamentoRepository *repositorios.DepartamentosRepository
}

func NewDepartamentoService(departamentoRepository *repositorios.DepartamentosRepository) *DepartamentoService {
	return &DepartamentoService{departamentoRepository: departamentoRepository}
}

func (s *DepartamentoService) CrearDepartamento(ctx context.Context, departamento dominio.Departamento) (dominio.Departamento, error) {
	if departamento.Nombre == "" || departamento.Direccion == "" {
		return dominio.Departamento{}, dominio.ValidationError("datos de departamento inválidos")
	}
	return s.departamentoRepository.CreateDepartamento(ctx, departamento)
}

func (s *DepartamentoService) ObtenerDepartamentos(ctx context.Context) ([]dominio.Departamento, error) {
	return s.departamentoRepository.ListDepartamentos(ctx)
}

func (s *DepartamentoService) ObtenerDepartamento(ctx context.Context, id int32) (dominio.Departamento, error) {
	return s.departamentoRepository.GetDepartamento(ctx, id)
}

func (s *DepartamentoService) ActualizarDepartamento(ctx context.Context, departamento dominio.Departamento) error {
	if departamento.IDDepto <= 0 || departamento.Nombre == "" || departamento.Direccion == "" {
		return dominio.ValidationError("datos de departamento inválidos")
	}
	return s.departamentoRepository.UpdateDepartamento(ctx, departamento)
}

func (s *DepartamentoService) EliminarDepartamento(ctx context.Context, id int32) error {
	return s.departamentoRepository.DeleteDepartamento(ctx, id)
}