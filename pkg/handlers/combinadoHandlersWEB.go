package handlers

import (
	"net/http"

	"Los5/pkg/dominio"
	"Los5/pkg/logica"
	"Los5/views"
)

type CombinadoWebHandler struct {
	combinadoService *logica.CombinadoService
}

func NewCombinadoWebHandler(combinadoService *logica.CombinadoService) *CombinadoWebHandler {
	return &CombinadoWebHandler{combinadoService: combinadoService}
}

func (h *CombinadoWebHandler) ReservasPorDepartamentoWeb(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		responderErrorWeb(w, r, err)
		return
	}

	pagina, err := h.combinadoService.ObtenerReservasPorDepartamento(r.Context(), id)
	if err != nil {
		responderErrorWeb(w, r, err)
		return
	}

	renderHTML(w, r, http.StatusOK, views.ReservasPorDepartamentoPage(views.ReservasPorDepartamentoPageData{
		Departamento: pagina.Departamento,
		Reservas:     pagina.Reservas,
	}))
}

func (h *CombinadoWebHandler) HuespedesPorDepartamentoWeb(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		responderErrorWeb(w, r, err)
		return
	}

	pagina, err := h.combinadoService.ObtenerHuespedesPorDepartamento(r.Context(), id)
	if err != nil {
		responderErrorWeb(w, r, err)
		return
	}

	renderHTML(w, r, http.StatusOK, views.HuespedesPorDepartamentoPage(views.HuespedesPorDepartamentoPageData{
		Departamento: pagina.Departamento,
		Huespedes:    pagina.Huespedes,
	}))
}

func (h *CombinadoWebHandler) DisponibilidadWeb(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		responderErrorWeb(w, r, err)
		return
	}

	desdeStr := r.URL.Query().Get("desde")
	hastaStr := r.URL.Query().Get("hasta")

	if desdeStr == "" && hastaStr == "" {
		departamento, err := h.combinadoService.ObtenerDepartamento(r.Context(), id)
		if err != nil {
			responderErrorWeb(w, r, err)
			return
		}

		renderHTML(w, r, http.StatusOK, views.DisponibilidadPage(dominio.Disponibilidad{
			Departamento: departamento,
		}))
		return
	}

	desde, hasta, err := parseRango(desdeStr, hastaStr)
	if err != nil {
		responderErrorWeb(w, r, err)
		return
	}

	disponibilidad, err := h.combinadoService.ConsultarDisponibilidad(r.Context(), id, desde, hasta)
	if err != nil {
		responderErrorWeb(w, r, err)
		return
	}

	renderHTML(w, r, http.StatusOK, views.DisponibilidadPage(disponibilidad))
}

func (h *CombinadoWebHandler) TotalReservaWeb(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		responderErrorWeb(w, r, err)
		return
	}

	total, err := h.combinadoService.CalcularTotalReserva(r.Context(), id)
	if err != nil {
		responderErrorWeb(w, r, err)
		return
	}

	renderHTML(w, r, http.StatusOK, views.TotalReservaPage(total))
}
