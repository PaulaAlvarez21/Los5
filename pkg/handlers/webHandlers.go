package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"Los5/pkg/dominio"
	"Los5/pkg/logica"
	"Los5/views"
)

type WebHandler struct {
	departamentoService *logica.DepartamentoService
	reservaService      *logica.ReservaService
	huespedService      *logica.HuespedService
}

func NewWebHandler(
	departamentoService *logica.DepartamentoService,
	reservaService *logica.ReservaService,
	huespedService *logica.HuespedService,
) *WebHandler {
	return &WebHandler{
		departamentoService: departamentoService,
		reservaService:      reservaService,
		huespedService:      huespedService,
	}
}

func (h *WebHandler) InicioWeb(w http.ResponseWriter, r *http.Request) {
	renderHTML(w, r, http.StatusOK, views.Layout("LOS 5 ALQUILERES", ""))
}

func (h *WebHandler) DepartamentosWeb(w http.ResponseWriter, r *http.Request) {
	departamentos, err := h.departamentoService.ObtenerDepartamentos(r.Context())
	if err != nil {
		responderErrorWeb(w, r, err)
		return
	}

	renderHTML(w, r, http.StatusOK, views.DepartamentosPage(departamentos))
}

func (h *WebHandler) HuespedesWeb(w http.ResponseWriter, r *http.Request) {
	huespedes, err := h.huespedService.ObtenerHuespedes(r.Context())
	if err != nil {
		responderErrorWeb(w, r, err)
		return
	}

	renderHTML(w, r, http.StatusOK, views.HuespedesPage(huespedes))
}

func (h *WebHandler) NuevoDepartamentoWeb(w http.ResponseWriter, r *http.Request) {
	renderHTML(w, r, http.StatusOK, views.DepartamentoFormPage())
}

func (h *WebHandler) NuevaReservaWeb(w http.ResponseWriter, r *http.Request) {
	departamentos, err := h.departamentoService.ObtenerDepartamentos(r.Context())
	if err != nil {
		responderErrorWeb(w, r, err)
		return
	}

	huespedes, err := h.huespedService.ObtenerHuespedes(r.Context())
	if err != nil {
		responderErrorWeb(w, r, err)
		return
	}

	renderHTML(w, r, http.StatusOK, views.ReservaFormPage(departamentos, huespedes))
}

func (h *WebHandler) NuevoHuespedWeb(w http.ResponseWriter, r *http.Request) {
	reservas, err := h.reservaService.ObtenerReservas(r.Context())
	if err != nil {
		responderErrorWeb(w, r, err)
		return
	}

	renderHTML(w, r, http.StatusOK, views.HuespedFormPage(reservas))
}

func (h *WebHandler) CrearDepartamentoWeb(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		responderErrorWeb(w, r, dominio.ValidationError("formulario inválido"))
		return
	}

	departamento := dominio.Departamento{
		Nombre:      r.FormValue("nombre"),
		Direccion:   r.FormValue("direccion"),
		Disponible:  true,
		Limpio:      r.FormValue("limpio") != "false",
		Descripcion: stringOpcional(r.FormValue("descripcion")),
	}

	if _, err := h.departamentoService.CrearDepartamento(r.Context(), departamento); err != nil {
		responderErrorWeb(w, r, err)
		return
	}

	http.Redirect(w, r, "/web/departamentos", http.StatusSeeOther)
}

func (h *WebHandler) CrearReservaWeb(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		responderErrorWeb(w, r, dominio.ValidationError("formulario inválido"))
		return
	}

	idDepto, err := strconv.Atoi(r.FormValue("id_depto"))
	if err != nil || idDepto <= 0 {
		responderErrorWeb(w, r, dominio.ValidationError("departamento inválido"))
		return
	}

	cantNoches, err := strconv.Atoi(r.FormValue("cant_noches"))
	if err != nil || cantNoches <= 0 {
		responderErrorWeb(w, r, dominio.ValidationError("cantidad de noches inválida"))
		return
	}

	fechaInicio, err := parseFecha(r.FormValue("fecha_inicio"))
	if err != nil {
		responderErrorWeb(w, r, err)
		return
	}

	fechaFin, err := parseFecha(r.FormValue("fecha_fin"))
	if err != nil {
		responderErrorWeb(w, r, err)
		return
	}

	precioBase := strings.TrimSpace(r.FormValue("precio_base"))
	if precioBase == "" {
		responderErrorWeb(w, r, dominio.ValidationError("precio base inválido"))
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
		responderErrorWeb(w, r, err)
		return
	}

	if rawIDHuesped := r.FormValue("id_huesped"); rawIDHuesped != "" {
		idHuesped, err := strconv.Atoi(rawIDHuesped)
		if err != nil || idHuesped <= 0 {
			responderErrorWeb(w, r, dominio.ValidationError("huésped inválido"))
			return
		}

		huesped, err := h.huespedService.ObtenerHuesped(r.Context(), int32(idHuesped))
		if err != nil {
			responderErrorWeb(w, r, err)
			return
		}

		huesped.IDReserva = &reservaCreada.IDReserva
		if err := h.huespedService.ActualizarHuesped(r.Context(), huesped); err != nil {
			responderErrorWeb(w, r, err)
			return
		}
	}

	http.Redirect(w, r, "/web/reservas", http.StatusSeeOther)
}

func (h *WebHandler) CrearHuespedWeb(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		responderErrorWeb(w, r, dominio.ValidationError("formulario inválido"))
		return
	}

	var idReserva *int32
	if rawIDReserva := strings.TrimSpace(r.FormValue("id_reserva")); rawIDReserva != "" {
		parsed, err := strconv.Atoi(rawIDReserva)
		if err != nil || parsed <= 0 {
			responderErrorWeb(w, r, dominio.ValidationError("la reserva es inválida"))
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
		responderErrorWeb(w, r, err)
		return
	}

	http.Redirect(w, r, "/web/huespedes", http.StatusSeeOther)
}

func (h *WebHandler) CambiarHigieneWeb(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		responderErrorWeb(w, r, err)
		return
	}

	departamento, err := h.departamentoService.ObtenerDepartamento(r.Context(), id)
	if err != nil {
		responderErrorWeb(w, r, err)
		return
	}

	departamento.Limpio = !departamento.Limpio

	if err := h.departamentoService.ActualizarDepartamento(r.Context(), departamento); err != nil {
		responderErrorWeb(w, r, err)
		return
	}

	http.Redirect(w, r, "/web/departamentos", http.StatusSeeOther)
}

func stringOpcional(valor string) *string {
	valor = strings.TrimSpace(valor)
	if valor == "" {
		return nil
	}
	return &valor
}
