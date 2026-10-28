package handlers

import (
	"Los5/pkg/logica"
)

type CombinadoHandler struct {
	departamentoService *logica.DepartamentoService
}

func NewCombinadoHandler(departamentoService *logica.DepartamentoService) *CombinadoHandler {
	return &CombinadoHandler{departamentoService: departamentoService}
}
