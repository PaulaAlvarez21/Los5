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

// NuevoHuespedWeb devuelve únicamente el modal con el formulario de alta.
func (h *HuespedWebHandler) NuevoHuespedWeb(w http.ResponseWriter, r *http.Request) {
	mostrarReserva := r.URL.Query().Get("sin_reserva") != "1"

	var reservas []dominio.Reserva
	if mostrarReserva {
		var err error
		reservas, err = h.reservaService.ObtenerReservas(r.Context())
		if err != nil {
			responderErrorWeb(w, r, err)
			return
		}
	}

	renderHTML(w, r, http.StatusOK, views.HuespedFormModal(reservas, "", mostrarReserva))
}

func (h *HuespedWebHandler) CrearHuespedWeb(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.responderErrorAlta(w, r, http.StatusBadRequest, "formulario inválido")
		return
	}

	var idReserva *int32
	if rawIDReserva := strings.TrimSpace(r.FormValue("id_reserva")); rawIDReserva != "" {
		parsed, err := strconv.Atoi(rawIDReserva)
		if err != nil || parsed <= 0 {
			h.responderErrorAlta(w, r, http.StatusBadRequest, "la reserva es inválida")
			return
		}
		valor := int32(parsed)
		idReserva = &valor
	}

	telefono := strings.TrimSpace(r.FormValue("telefono"))
	if telefono == "" {
		h.responderErrorAlta(w, r, http.StatusBadRequest, "el teléfono es obligatorio")
		return
	}

	huesped := dominio.Huesped{
		IDReserva:     idReserva,
		Nombre:        r.FormValue("nombre"),
		Apellido:      r.FormValue("apellido"),
		Telefono:      &telefono,
		Email:         stringOpcional(r.FormValue("email")),
		Observaciones: stringOpcional(r.FormValue("observaciones")),
	}

	creado, err := h.huespedService.CrearHuesped(r.Context(), huesped)
	if err != nil {
		h.responderErrorAlta(w, r, statusDeError(err), mensajeDeError(err))
		return
	}

	if esHTMX(r) {
		huespedes, err := h.huespedService.ObtenerHuespedes(r.Context())
		if err != nil {
			h.responderErrorAlta(w, r, statusDeError(err), mensajeDeError(err))
			return
		}

		renderHTML(w, r, http.StatusOK, views.HuespedCreado(creado, len(huespedes) == 1))
		return
	}

	http.Redirect(w, r, "/huespedes", http.StatusSeeOther)
}

func (h *HuespedWebHandler) EliminarHuespedWeb(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		h.responderErrorBorrado(w, r, err)
		return
	}

	if err := h.huespedService.EliminarHuesped(r.Context(), id); err != nil {
		h.responderErrorBorrado(w, r, err)
		return
	}

	if esHTMX(r) {
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/huespedes", http.StatusSeeOther)
}

func (h *HuespedWebHandler) responderErrorAlta(w http.ResponseWriter, r *http.Request, status int, mensaje string) {
	if esHTMX(r) {
		h.renderHuespedModalError(w, r, mensaje)
		return
	}

	h.renderHuespedes(w, r, status, mensaje)
}

func (h *HuespedWebHandler) renderHuespedModalError(w http.ResponseWriter, r *http.Request, mensaje string) {
	mostrarReserva := r.URL.Query().Get("sin_reserva") != "1" && r.FormValue("sin_reserva") != "1"

	var reservas []dominio.Reserva
	if mostrarReserva {
		var err error
		reservas, err = h.reservaService.ObtenerReservas(r.Context())
		if err != nil {
			responderErrorHTMX(w, r, err, "#modal-root")
			return
		}
	}

	retarget(w, "#modal-root")
	renderHTML(w, r, http.StatusOK, views.HuespedFormModal(reservas, mensaje, mostrarReserva))
}

func (h *HuespedWebHandler) responderErrorBorrado(w http.ResponseWriter, r *http.Request, err error) {
	if esHTMX(r) {
		responderErrorHTMX(w, r, err, "#page-errors")
		return
	}

	h.renderHuespedes(w, r, statusDeError(err), mensajeDeError(err))
}
