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

func CrearReservaHandler(w http.ResponseWriter, r *http.Request) {

}

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

func ActualizarReservaHandler(w http.ResponseWriter, r *http.Request) {
}

func EliminarReservaHandler(w http.ResponseWriter, r *http.Request) {
}
