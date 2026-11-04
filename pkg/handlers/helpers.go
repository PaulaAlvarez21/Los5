package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"Los5/pkg/dominio"
	"Los5/pkg/logica"
	"Los5/views"

	"github.com/a-h/templ"
)

// renderHTML escribe un componente templ en la respuesta HTTP.
func renderHTML(w http.ResponseWriter, r *http.Request, status int, componente templ.Component) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := componente.Render(r.Context(), w); err != nil {
		log.Printf("error renderizando template: %v", err)
	}
}

// responderErrorWeb traduce un error de dominio a una página HTML de error.
func responderErrorWeb(w http.ResponseWriter, r *http.Request, err error) {
	renderHTML(w, r, statusDeError(err), views.ErrorPage(statusDeError(err), mensajeDeError(err)))
}

func statusDeError(err error) int {
	var appErr *dominio.AppError
	if errors.As(err, &appErr) {
		return appErr.Status
	}
	return http.StatusInternalServerError
}

func mensajeDeError(err error) string {
	var appErr *dominio.AppError
	if errors.As(err, &appErr) {
		return appErr.Msg
	}
	return "error interno"
}

func parseID(r *http.Request) (int32, error) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		return 0, dominio.ValidationError("ID inválido")
	}
	return int32(id), nil
}

func parseFecha(valor string) (time.Time, error) {
	fecha, err := time.Parse("2006-01-02", valor)
	if err != nil {
		return time.Time{}, dominio.ValidationError("fecha inválida")
	}
	return fecha, nil
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

func filtroVista(r *http.Request) views.FiltroReservas {
	query := r.URL.Query()
	return views.FiltroReservas{
		DeptoID: query.Get("depto_id"),
		Desde:   query.Get("desde"),
		Hasta:   query.Get("hasta"),
	}
}

func stringOpcional(valor string) *string {
	valor = strings.TrimSpace(valor)
	if valor == "" {
		return nil
	}
	return &valor
}

func parseBoolForm(valor string) bool {
	switch strings.ToLower(strings.TrimSpace(valor)) {
	case "true", "on", "1", "si", "sí":
		return true
	default:
		return false
	}
}
