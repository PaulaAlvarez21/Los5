package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"Los5/pkg/dominio"
	"Los5/pkg/logica"
)

type ReservaHandler struct {
	reservaService *logica.ReservaService
}

func NewReservaHandler(reservaService *logica.ReservaService) *ReservaHandler {
	return &ReservaHandler{reservaService: reservaService}
}

func responderError(w http.ResponseWriter, err error) {
	var appErr *dominio.AppError
	if errors.As(err, &appErr) {
		http.Error(w, appErr.Msg, appErr.Status)
		return
	}
	http.Error(w, "error interno", http.StatusInternalServerError) //500
}

// Crear una reserva
func (h *ReservaHandler) CrearReserva(w http.ResponseWriter, r *http.Request) {
	var reserva dominio.Reserva
	if err := json.NewDecoder(r.Body).Decode(&reserva); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest) //400
		return
	}

	reservaCreada, err := h.reservaService.CrearReserva(r.Context(), reserva)
	if err != nil {
		responderError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) //201
	json.NewEncoder(w).Encode(reservaCreada)
}

// Obtener una reserva con ID
func (h *ReservaHandler) ObtenerReserva(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest) //400
		return
	}

	reserva, err := h.reservaService.ObtenerReserva(r.Context(), int32(id))
	if err != nil {
		responderError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reserva)
}

// Obtener todas las reservas
func (h *ReservaHandler) ObtenerReservas(w http.ResponseWriter, r *http.Request) {
	reservas, err := h.reservaService.ObtenerReservas(r.Context())
	if err != nil {
		responderError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservas)
}

// Actualizar una reserva
func (h *ReservaHandler) ActualizarReserva(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var reserva dominio.Reserva
	if err := json.NewDecoder(r.Body).Decode(&reserva); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	reserva.IDReserva = int32(id)

	if err := h.reservaService.ActualizarReserva(r.Context(), reserva); err != nil {
		responderError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reserva)
}

// Eliminar una reserva
func (h *ReservaHandler) EliminarReserva(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := h.reservaService.EliminarReserva(r.Context(), int32(id)); err != nil {
		responderError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent) //204
}
