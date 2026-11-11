package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"Los5/pkg/dominio"
	"Los5/pkg/logica"
	"Los5/views"
)

type ReservaWebHandler struct {
	reservaService   *logica.ReservaService
	huespedService   *logica.HuespedService
	combinadoService *logica.CombinadoService
}

func NewReservaWebHandler(
	reservaService *logica.ReservaService,
	huespedService *logica.HuespedService,
	combinadoService *logica.CombinadoService,
) *ReservaWebHandler {
	return &ReservaWebHandler{
		reservaService:   reservaService,
		huespedService:   huespedService,
		combinadoService: combinadoService,
	}
}

func (h *ReservaWebHandler) ReservasWeb(w http.ResponseWriter, r *http.Request) {
	filtro, err := parseFiltroReservas(r)
	if err != nil {
		h.renderReservas(w, r, statusDeError(err), mensajeDeError(err))
		return
	}

	pagina, err := h.combinadoService.ObtenerPaginaReservas(r.Context(), filtro)
	if err != nil {
		h.renderReservas(w, r, statusDeError(err), mensajeDeError(err))
		return
	}

	h.renderPaginaReservas(w, r, http.StatusOK, "", pagina)
}

func (h *ReservaWebHandler) renderReservas(w http.ResponseWriter, r *http.Request, status int, mensaje string) {
	pagina, err := h.combinadoService.ObtenerPaginaReservas(r.Context(), logica.FiltroReservas{})
	if err != nil {
		responderErrorWeb(w, r, err)
		return
	}

	h.renderPaginaReservas(w, r, status, mensaje, pagina)
}

func (h *ReservaWebHandler) renderPaginaReservas(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	mensaje string,
	pagina dominio.PaginaReservas,
) {
	renderHTML(w, r, status, views.ReservasPage(views.ReservasPageData{
		Reservas:      pagina.Reservas,
		Departamentos: pagina.Departamentos,
		Huespedes:     pagina.Huespedes,
		Filtro:        filtroVista(r),
		Error:         mensaje,
	}))
}

func (h *ReservaWebHandler) CrearReservaWeb(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderReservas(w, r, http.StatusBadRequest, "formulario inválido")
		return
	}

	idDepto, err := strconv.Atoi(r.FormValue("id_depto"))
	if err != nil || idDepto <= 0 {
		h.renderReservas(w, r, http.StatusBadRequest, "departamento inválido")
		return
	}

	cantNoches, err := strconv.Atoi(r.FormValue("cant_noches"))
	if err != nil || cantNoches <= 0 {
		h.renderReservas(w, r, http.StatusBadRequest, "cantidad de noches inválida")
		return
	}

	fechaInicio, err := parseFecha(r.FormValue("fecha_inicio"))
	if err != nil {
		h.renderReservas(w, r, statusDeError(err), mensajeDeError(err))
		return
	}

	fechaFin, err := parseFecha(r.FormValue("fecha_fin"))
	if err != nil {
		h.renderReservas(w, r, statusDeError(err), mensajeDeError(err))
		return
	}

	precioBase := strings.TrimSpace(r.FormValue("precio_base"))
	if precioBase == "" {
		h.renderReservas(w, r, http.StatusBadRequest, "precio base inválido")
		return
	}

	reserva := dominio.Reserva{
		FechaInicio:   fechaInicio,
		IDDepto:       int32(idDepto),
		FechaFin:      fechaFin,
		PrecioBase:    precioBase,
		CantNoches:    int32(cantNoches),
		Descuento:     stringOpcional(r.FormValue("descuento")),
		Observaciones: stringOpcional(r.FormValue("observaciones")),
	}

	reservaCreada, err := h.reservaService.CrearReserva(r.Context(), reserva)
	if err != nil {
		h.renderReservas(w, r, statusDeError(err), mensajeDeError(err))
		return
	}

	if rawIDHuesped := strings.TrimSpace(r.FormValue("id_huesped")); rawIDHuesped != "" {
		idHuesped, err := strconv.Atoi(rawIDHuesped)
		if err != nil || idHuesped <= 0 {
			h.renderReservas(w, r, http.StatusBadRequest, "huésped inválido")
			return
		}

		huesped, err := h.huespedService.ObtenerHuesped(r.Context(), int32(idHuesped))
		if err != nil {
			h.renderReservas(w, r, statusDeError(err), mensajeDeError(err))
			return
		}

		huesped.IDReserva = &reservaCreada.IDReserva
		if err := h.huespedService.ActualizarHuesped(r.Context(), huesped); err != nil {
			h.renderReservas(w, r, statusDeError(err), mensajeDeError(err))
			return
		}
	}

	http.Redirect(w, r, "/reservas", http.StatusSeeOther)
}

func (h *ReservaWebHandler) EliminarReservaWeb(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		h.renderReservas(w, r, statusDeError(err), mensajeDeError(err))
		return
	}

	if err := h.reservaService.EliminarReserva(r.Context(), id); err != nil {
		h.renderReservas(w, r, statusDeError(err), mensajeDeError(err))
		return
	}

	http.Redirect(w, r, "/reservas", http.StatusSeeOther)
}
