package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"Los5/pkg/dominio"
	"Los5/pkg/logica"
	"Los5/views"
)

type HuespedWebHandler struct {
	huespedService *logica.HuespedService
	reservaService *logica.ReservaService
}

func NewHuespedWebHandler(
	huespedService *logica.HuespedService,
	reservaService *logica.ReservaService,
) *HuespedWebHandler {
	return &HuespedWebHandler{
		huespedService: huespedService,
		reservaService: reservaService,
	}
}

func (h *HuespedWebHandler) HuespedesWeb(w http.ResponseWriter, r *http.Request) {
	h.renderHuespedes(w, r, http.StatusOK, "")
}

func (h *HuespedWebHandler) renderHuespedes(w http.ResponseWriter, r *http.Request, status int, mensaje string) {
	huespedes, err := h.huespedService.ObtenerHuespedes(r.Context())
	if err != nil {
		responderErrorWeb(w, r, err)
		return
	}

	reservas, err := h.reservaService.ObtenerReservas(r.Context())
	if err != nil {
		responderErrorWeb(w, r, err)
		return
	}

	renderHTML(w, r, status, views.HuespedesPage(views.HuespedesPageData{
		Huespedes: huespedes,
		Reservas:  reservas,
		Error:     mensaje,
	}))
}

func (h *HuespedWebHandler) CrearHuespedWeb(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderHuespedes(w, r, http.StatusBadRequest, "formulario inválido")
		return
	}

	var idReserva *int32
	if rawIDReserva := strings.TrimSpace(r.FormValue("id_reserva")); rawIDReserva != "" {
		parsed, err := strconv.Atoi(rawIDReserva)
		if err != nil || parsed <= 0 {
			h.renderHuespedes(w, r, http.StatusBadRequest, "la reserva es inválida")
			return
		}
		valor := int32(parsed)
		idReserva = &valor
	}

	huesped := dominio.Huesped{
		IDReserva:     idReserva,
		Nombre:        r.FormValue("nombre"),
		Apellido:      r.FormValue("apellido"),
		Telefono:      stringOpcional(r.FormValue("telefono")),
		Email:         stringOpcional(r.FormValue("email")),
		Observaciones: stringOpcional(r.FormValue("observaciones")),
	}

	if _, err := h.huespedService.CrearHuesped(r.Context(), huesped); err != nil {
		h.renderHuespedes(w, r, statusDeError(err), mensajeDeError(err))
		return
	}

	http.Redirect(w, r, "/huespedes", http.StatusSeeOther)
}

func (h *HuespedWebHandler) EliminarHuespedWeb(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		h.renderHuespedes(w, r, statusDeError(err), mensajeDeError(err))
		return
	}

	if err := h.huespedService.EliminarHuesped(r.Context(), id); err != nil {
		h.renderHuespedes(w, r, statusDeError(err), mensajeDeError(err))
		return
	}

	http.Redirect(w, r, "/huespedes", http.StatusSeeOther)
}
