package handlers

import (
	"encoding/json"
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

// POST /reservas - Crear una reserva
func (h *ReservaHandler) CrearReserva(w http.ResponseWriter, r *http.Request) {
	var reserva dominio.Reserva
	err := json.NewDecoder(r.Body).Decode(&reserva)
	if err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	reservaCreada, err := h.reservaService.CrearReserva(reserva)
	if err != nil {
		http.Error(w, "Error al crear la reserva", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(reservaCreada)
}

// GET /reservas/{id} - Obtener una reserva con ID
func (h *ReservaHandler) ObtenerReserva(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid reserva ID", http.StatusBadRequest)
		return
	}

	reserva, err := h.reservaService.ObtenerReserva(int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reserva)
}

// GET /reservas - Obtener todas las reservas
func (h *ReservaHandler) ObtenerReservas(w http.ResponseWriter, r *http.Request) {
	reservas, err := h.reservaService.ObtenerReservas()
	if err != nil {
		http.Error(w, "Error al obtener las reservas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservas)
}

// PUT /reservas/{id} - Actualizar una reserva
func (h *ReservaHandler) ActualizarReserva(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var reserva dominio.Reserva
	err = json.NewDecoder(r.Body).Decode(&reserva)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	reserva.IDReserva = int32(id)

	err = h.reservaService.ActualizarReserva(reserva)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reserva)
}

// DELETE /reservas/{id} - Eliminar una reserva
func (h *ReservaHandler) EliminarReserva(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	err = h.reservaService.EliminarReserva(int32(id))
	if err != nil {
		http.Error(w, "Error al eliminar la reserva", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
