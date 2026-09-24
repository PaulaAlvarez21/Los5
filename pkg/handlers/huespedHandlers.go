package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"Los5/pkg/dominio"
	"Los5/pkg/logica"
)

type HuespedHandler struct {
	huespedService *logica.HuespedService
}

func NewHuespedHandler(huespedService *logica.HuespedService) *HuespedHandler {
	return &HuespedHandler{huespedService: huespedService}
}

// Crear un huesped
func (h *HuespedHandler) CrearHuesped(w http.ResponseWriter, r *http.Request) {
	var huesped dominio.Huesped
	if err := json.NewDecoder(r.Body).Decode(&huesped); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest) //400
		return
	}

	huespedCreado, err := h.huespedService.CrearHuesped(r.Context(), huesped)
	if err != nil {
		responderError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) //201
	json.NewEncoder(w).Encode(huespedCreado)
}

// Obtener un huesped con ID
func (h *HuespedHandler) ObtenerHuesped(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest) //400
		return
	}

	huesped, err := h.huespedService.ObtenerHuesped(r.Context(), int32(id))
	if err != nil {
		responderError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(huesped)
}

// Obtener todos los huespedes
func (h *HuespedHandler) ObtenerHuespedes(w http.ResponseWriter, r *http.Request) {
	huespedes, err := h.huespedService.ObtenerHuespedes(r.Context())
	if err != nil {
		responderError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(huespedes)
}

// Actualizar un huesped
func (h *HuespedHandler) ActualizarHuesped(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var huesped dominio.Huesped
	if err := json.NewDecoder(r.Body).Decode(&huesped); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	huesped.IDHuesped = int32(id)

	if err := h.huespedService.ActualizarHuesped(r.Context(), huesped); err != nil {
		responderError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(huesped)
}

// Eliminar un huesped
func (h *HuespedHandler) EliminarHuesped(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.huespedService.EliminarHuesped(r.Context(), int32(id)); err != nil {
		responderError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent) //204
}