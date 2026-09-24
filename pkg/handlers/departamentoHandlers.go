package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"Los5/pkg/dominio"
	"Los5/pkg/logica"
)

type DepartamentoHandler struct {
	departamentoService *logica.DepartamentoService
}

func NewDepartamentoHandler(departamentoService *logica.DepartamentoService) *DepartamentoHandler {
	return &DepartamentoHandler{departamentoService: departamentoService}
}

// Crear un departamento
func (h *DepartamentoHandler) CrearDepartamento(w http.ResponseWriter, r *http.Request) {
	var departamento dominio.Departamento
	if err := json.NewDecoder(r.Body).Decode(&departamento); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest) //400
		return
	}

	departamentoCreado, err := h.departamentoService.CrearDepartamento(r.Context(), departamento)
	if err != nil {
		responderError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) //201
	json.NewEncoder(w).Encode(departamentoCreado)
}

// Obtener un departamento con ID
func (h *DepartamentoHandler) ObtenerDepartamento(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest) //400
		return
	}

	departamento, err := h.departamentoService.ObtenerDepartamento(r.Context(), int32(id))
	if err != nil {
		responderError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(departamento)
}

// Obtener todos los departamentos
func (h *DepartamentoHandler) ObtenerDepartamentos(w http.ResponseWriter, r *http.Request) {
	departamentos, err := h.departamentoService.ObtenerDepartamentos(r.Context())
	if err != nil {
		responderError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(departamentos)
}

// Actualizar un departamento
func (h *DepartamentoHandler) ActualizarDepartamento(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var departamento dominio.Departamento
	if err := json.NewDecoder(r.Body).Decode(&departamento); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	departamento.IDDepto = int32(id)

	if err := h.departamentoService.ActualizarDepartamento(r.Context(), departamento); err != nil {
		responderError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(departamento)
}

// Eliminar un departamento
func (h *DepartamentoHandler) EliminarDepartamento(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.departamentoService.EliminarDepartamento(r.Context(), int32(id)); err != nil {
		responderError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent) //204
}