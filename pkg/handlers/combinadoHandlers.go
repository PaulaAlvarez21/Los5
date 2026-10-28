package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"Los5/pkg/dominio"
	"Los5/pkg/logica"
	"Los5/views"

	"github.com/a-h/templ"
)

type CombinadoHandler struct {
	combinadoService *logica.CombinadoService
}

func NewCombinadoHandler(combinadoService *logica.CombinadoService) *CombinadoHandler {
	return &CombinadoHandler{combinadoService: combinadoService}
}

func (h *CombinadoHandler) ReservasWeb(w http.ResponseWriter, r *http.Request) {
	filtro, err := parseFiltroReservas(r)
	if err != nil {
		responderErrorWeb(w, r, err)
		return
	}

	pagina, err := h.combinadoService.ObtenerPaginaReservas(r.Context(), filtro)
	if err != nil {
		responderErrorWeb(w, r, err)
		return
	}

	renderHTML(w, r, http.StatusOK, views.ReservasPage(
		pagina.Reservas,
		pagina.Departamentos,
		pagina.Huespedes,
		"",
	))
}

func (h *CombinadoHandler) HuespedesPorDepartamentoWeb(w http.ResponseWriter, r *http.Request) {
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

	renderHTML(w, r, http.StatusOK, views.HuespedesPage(pagina.Huespedes))
}

func (h *CombinadoHandler) DisponibilidadWeb(w http.ResponseWriter, r *http.Request) {
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

func (h *CombinadoHandler) TotalReservaWeb(w http.ResponseWriter, r *http.Request) {
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

func renderHTML(w http.ResponseWriter, r *http.Request, status int, componente templ.Component) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := componente.Render(r.Context(), w); err != nil {
		log.Printf("error renderizando template: %v", err)
	}
}

func responderErrorWeb(w http.ResponseWriter, r *http.Request, err error) {
	status := http.StatusInternalServerError
	mensaje := "error interno"

	var appErr *dominio.AppError
	if errors.As(err, &appErr) {
		status = appErr.Status
		mensaje = appErr.Msg
	}

	renderHTML(w, r, status, views.ErrorPage(status, mensaje))
}

func parseID(r *http.Request) (int32, error) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		return 0, dominio.ValidationError("ID inválido")
	}
	return int32(id), nil
}

func parseFiltroReservas(r *http.Request) (logica.FiltroReservas, error) {
	var filtro logica.FiltroReservas
	query := r.URL.Query()

	if raw := query.Get("depto_id"); raw != "" {
		id, err := strconv.Atoi(raw)
		if err != nil || id <= 0 {
			return filtro, dominio.ValidationError("id de departamento inválido")
		}
		idDepto := int32(id)
		filtro.IDDepto = &idDepto
	}

	if raw := query.Get("desde"); raw != "" {
		desde, err := parseFecha(raw)
		if err != nil {
			return filtro, err
		}
		filtro.Desde = &desde
	}

	if raw := query.Get("hasta"); raw != "" {
		hasta, err := parseFecha(raw)
		if err != nil {
			return filtro, err
		}
		filtro.Hasta = &hasta
	}

	if filtro.Desde != nil && filtro.Hasta != nil && !filtro.Desde.Before(*filtro.Hasta) {
		return filtro, dominio.ValidationError("rango de fechas inválido")
	}

	return filtro, nil
}

func parseRango(desdeStr, hastaStr string) (time.Time, time.Time, error) {
	if desdeStr == "" || hastaStr == "" {
		return time.Time{}, time.Time{}, dominio.ValidationError("las fechas desde y hasta son obligatorias")
	}

	desde, err := parseFecha(desdeStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	hasta, err := parseFecha(hastaStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	if !desde.Before(hasta) {
		return time.Time{}, time.Time{}, dominio.ValidationError("rango de fechas inválido")
	}

	return desde, hasta, nil
}

func parseFecha(valor string) (time.Time, error) {
	fecha, err := time.Parse("2006-01-02", valor)
	if err != nil {
		return time.Time{}, dominio.ValidationError("fecha inválida")
	}
	return fecha, nil
}
