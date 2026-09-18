package handlers

import (
	"Los5/pkg/logica"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type ReservaHandler struct {
	reservaService *logica.ReservaService
}

func NewReservaHandler(reservaService *logica.ReservaService) *ReservaHandler {
	return &ReservaHandler{reservaService: reservaService}
}

// Crear una reserva
func (h *ReservaHandler) CrearReserva(w http.ResponseWriter, r *http.Request) {

	var reserva logica.CrearReservaRequest

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

// obtener una reserva con ID
func (h *ReservaHandler) ObtenerReservasHandler(w http.ResponseWriter, r *http.Request) {

	parts := strings.Split(r.URL.Path, "/")

	id, err := strconv.Atoi(parts[2])
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

// Obtener todas las reservas
func (h *ReservaHandler) ObtenerReservas(w http.ResponseWriter, r *http.Request) {

	reservas, err := h.reservaService.ObtenerReservas()

	if err != nil {
		http.Error(w, "Error al obtener las reservas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservas)
}

// Actualizar una reserva

func (h *ReservaHandler) ActualizarReserva(w http.ResponseWriter, r *http.Request, id int) {

	var reserva logica.ReservaActualizar //en services hay que crear un struct para actualizar reserva, porque no se puede usar el mismo que para crear reserva, ya que no tiene el ID

	err := json.NewDecoder(r.Body).Decode(&reserva)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	reserva.IDReserva = int32(id)

	reservaActualizada, err := h.reservaService.ActualizarReserva(reserva) //en services hay que crear un struct para actualizar reserva, porque no se puede usar el mismo que para crear reserva, ya que no tiene el ID
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservaActualizada)
}

// Eliminar una reserva
func (h *ReservaHandler) EliminarReserva(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	err = h.reservaService.EliminarReserva(id)

	if err != nil {
		http.Error(w, "Error al eliminar la reserva", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
