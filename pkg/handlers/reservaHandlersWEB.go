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

// NuevaReservaWeb devuelve únicamente el modal con el formulario de alta.
func (h *ReservaWebHandler) NuevaReservaWeb(w http.ResponseWriter, r *http.Request) {
	pagina, err := h.combinadoService.ObtenerPaginaReservas(r.Context(), logica.FiltroReservas{})
	if err != nil {
		responderErrorWeb(w, r, err)
		return
	}

	renderHTML(w, r, http.StatusOK, views.ReservaFormModal(
		pagina.Departamentos,
		pagina.Huespedes,
		"",
	))
}

// CalcularNochesWeb devuelve el campo de noches calculado a partir de las fechas.
func (h *ReservaWebHandler) CalcularNochesWeb(w http.ResponseWriter, r *http.Request) {
	valor := ""
	fechaInicio, errInicio := parseFecha(r.URL.Query().Get("fecha_inicio"))
	fechaFin, errFin := parseFecha(r.URL.Query().Get("fecha_fin"))
	if errInicio == nil && errFin == nil {
		if noches := int(fechaFin.Sub(fechaInicio).Hours() / 24); noches > 0 {
			valor = strconv.Itoa(noches)
		}
	}

	renderHTML(w, r, http.StatusOK, views.CampoCantNoches(valor))
}

func (h *ReservaWebHandler) CrearReservaWeb(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.responderErrorAlta(w, r, http.StatusBadRequest, "formulario inválido")
		return
	}

	idDepto, err := strconv.Atoi(r.FormValue("id_depto"))
	if err != nil || idDepto <= 0 {
		h.responderErrorAlta(w, r, http.StatusBadRequest, "departamento inválido")
		return
	}

	rawIDHuesped := strings.TrimSpace(r.FormValue("id_huesped"))
	idHuesped, err := strconv.Atoi(rawIDHuesped)
	if err != nil || idHuesped <= 0 {
		h.responderErrorAlta(w, r, http.StatusBadRequest, "huésped inválido")
		return
	}

	huesped, err := h.huespedService.ObtenerHuesped(r.Context(), int32(idHuesped))
	if err != nil {
		h.responderErrorAlta(w, r, statusDeError(err), mensajeDeError(err))
		return
	}

	fechaInicio, err := parseFecha(r.FormValue("fecha_inicio"))
	if err != nil {
		h.responderErrorAlta(w, r, statusDeError(err), mensajeDeError(err))
		return
	}

	fechaFin, err := parseFecha(r.FormValue("fecha_fin"))
	if err != nil {
		h.responderErrorAlta(w, r, statusDeError(err), mensajeDeError(err))
		return
	}

	cantNoches := int(fechaFin.Sub(fechaInicio).Hours() / 24)
	if cantNoches <= 0 {
		h.responderErrorAlta(w, r, http.StatusBadRequest, "cantidad de noches inválida")
		return
	}

	precioBase := strings.TrimSpace(r.FormValue("precio_base"))
	if precioBase == "" {
		h.responderErrorAlta(w, r, http.StatusBadRequest, "precio base inválido")
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
		h.responderErrorAlta(w, r, statusDeError(err), mensajeDeError(err))
		return
	}

	huesped.IDReserva = &reservaCreada.IDReserva
	if err := h.huespedService.ActualizarHuesped(r.Context(), huesped); err != nil {
		h.responderErrorAlta(w, r, statusDeError(err), mensajeDeError(err))
		return
	}

	if esHTMX(r) {
		pagina, err := h.combinadoService.ObtenerPaginaReservas(r.Context(), logica.FiltroReservas{})
		if err != nil {
			h.renderReservaModalError(w, r, mensajeDeError(err))
			return
		}

		detalle, ok := buscarReservaDetalle(pagina.Reservas, reservaCreada.IDReserva)
		if !ok {
			h.renderReservaModalError(w, r, "no se pudo recuperar la reserva creada")
			return
		}

		renderHTML(w, r, http.StatusOK, views.ReservaCreada(detalle, len(pagina.Reservas) == 1))
		return
	}

	http.Redirect(w, r, "/reservas", http.StatusSeeOther)
}

// primerHuespedDeReserva devuelve el huésped asociado con menor ID (el primero
// cargado). ListHuespedes ordena por apellido, así que no alcanza con tomar el
// primero del slice.
func primerHuespedDeReserva(huespedes []dominio.Huesped, idReserva int32) *dominio.Huesped {
	var primero *dominio.Huesped
	for i := range huespedes {
		huesped := &huespedes[i]
		if huesped.IDReserva == nil || *huesped.IDReserva != idReserva {
			continue
		}
		if primero == nil || huesped.IDHuesped < primero.IDHuesped {
			primero = huesped
		}
	}
	return primero
}

// EditarReservaWeb devuelve el modal con el formulario precargado. El selector
// de huésped muestra el primero asociado a la reserva.
func (h *ReservaWebHandler) EditarReservaWeb(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		h.renderReservas(w, r, statusDeError(err), mensajeDeError(err))
		return
	}

	reserva, err := h.reservaService.ObtenerReserva(r.Context(), id)
	if err != nil {
		h.renderReservas(w, r, statusDeError(err), mensajeDeError(err))
		return
	}

	pagina, err := h.combinadoService.ObtenerPaginaReservas(r.Context(), logica.FiltroReservas{})
	if err != nil {
		h.renderReservas(w, r, statusDeError(err), mensajeDeError(err))
		return
	}

	renderHTML(w, r, http.StatusOK, views.ReservaEditarModal(
		reserva,
		primerHuespedDeReserva(pagina.Huespedes, id),
		pagina.Departamentos,
		pagina.Huespedes,
	))
}

func (h *ReservaWebHandler) ActualizarReservaWeb(w http.ResponseWriter, r *http.Request) {
	fallo := func(err error) {
		if esHTMX(r) {
			responderErrorHTMX(w, r, err, "#page-errors")
			return
		}
		h.renderReservas(w, r, statusDeError(err), mensajeDeError(err))
	}

	id, err := parseID(r)
	if err != nil {
		fallo(err)
		return
	}

	if err := r.ParseForm(); err != nil {
		fallo(dominio.ValidationError("formulario inválido"))
		return
	}

	idDepto, err := strconv.Atoi(r.FormValue("id_depto"))
	if err != nil || idDepto <= 0 {
		fallo(dominio.ValidationError("departamento inválido"))
		return
	}

	rawIDHuesped := strings.TrimSpace(r.FormValue("id_huesped"))
	idHuesped, err := strconv.Atoi(rawIDHuesped)
	if err != nil || idHuesped <= 0 {
		fallo(dominio.ValidationError("huésped inválido"))
		return
	}

	huesped, err := h.huespedService.ObtenerHuesped(r.Context(), int32(idHuesped))
	if err != nil {
		fallo(err)
		return
	}

	fechaInicio, err := parseFecha(r.FormValue("fecha_inicio"))
	if err != nil {
		fallo(err)
		return
	}

	fechaFin, err := parseFecha(r.FormValue("fecha_fin"))
	if err != nil {
		fallo(err)
		return
	}

	cantNoches := int(fechaFin.Sub(fechaInicio).Hours() / 24)
	if cantNoches <= 0 {
		fallo(dominio.ValidationError("cantidad de noches inválida"))
		return
	}

	precioBase := strings.TrimSpace(r.FormValue("precio_base"))
	if precioBase == "" {
		fallo(dominio.ValidationError("precio base inválido"))
		return
	}

	reserva := dominio.Reserva{
		IDReserva:     id,
		FechaInicio:   fechaInicio,
		IDDepto:       int32(idDepto),
		FechaFin:      fechaFin,
		PrecioBase:    precioBase,
		CantNoches:    int32(cantNoches),
		Descuento:     stringOpcional(r.FormValue("descuento")),
		Observaciones: stringOpcional(r.FormValue("observaciones")),
	}

	if err := h.reservaService.ActualizarReserva(r.Context(), reserva); err != nil {
		fallo(err)
		return
	}

	// Solo se desasigna el huésped que estaba preseleccionado (el de menor ID),
	// nunca los demás asociados a esta reserva.
	if rawActual := strings.TrimSpace(r.FormValue("id_huesped_actual")); rawActual != "" {
		idActual, err := strconv.Atoi(rawActual)
		if err != nil || idActual <= 0 {
			fallo(dominio.ValidationError("huésped anterior inválido"))
			return
		}
		if idActual != idHuesped {
			anterior, err := h.huespedService.ObtenerHuesped(r.Context(), int32(idActual))
			if err != nil {
				fallo(err)
				return
			}
			if anterior.IDReserva != nil && *anterior.IDReserva == id {
				anterior.IDReserva = nil
				if err := h.huespedService.ActualizarHuesped(r.Context(), anterior); err != nil {
					fallo(err)
					return
				}
			}
		}
	}

	huesped.IDReserva = &reserva.IDReserva
	if err := h.huespedService.ActualizarHuesped(r.Context(), huesped); err != nil {
		fallo(err)
		return
	}

	if esHTMX(r) {
		pagina, err := h.combinadoService.ObtenerPaginaReservas(r.Context(), logica.FiltroReservas{})
		if err != nil {
			fallo(err)
			return
		}
		detalle, ok := buscarReservaDetalle(pagina.Reservas, id)
		if !ok {
			fallo(dominio.ValidationError("no se pudo recuperar la reserva editada"))
			return
		}
		renderHTML(w, r, http.StatusOK, views.ReservaEditado(detalle))
		return
	}

	http.Redirect(w, r, "/reservas", http.StatusSeeOther)
}

func (h *ReservaWebHandler) EliminarReservaWeb(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		h.responderErrorBorrado(w, r, err)
		return
	}

	if err := h.reservaService.EliminarReserva(r.Context(), id); err != nil {
		h.responderErrorBorrado(w, r, err)
		return
	}

	if esHTMX(r) {
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/reservas", http.StatusSeeOther)
}

func (h *ReservaWebHandler) responderErrorAlta(w http.ResponseWriter, r *http.Request, status int, mensaje string) {
	if esHTMX(r) {
		h.renderReservaModalError(w, r, mensaje)
		return
	}

	h.renderReservas(w, r, status, mensaje)
}

func (h *ReservaWebHandler) renderReservaModalError(w http.ResponseWriter, r *http.Request, mensaje string) {
	pagina, err := h.combinadoService.ObtenerPaginaReservas(r.Context(), logica.FiltroReservas{})
	if err != nil {
		responderErrorHTMX(w, r, err, "#modal-root")
		return
	}

	retarget(w, "#modal-root")
	renderHTML(w, r, http.StatusOK, views.ReservaFormModal(
		pagina.Departamentos,
		pagina.Huespedes,
		mensaje,
	))
}

func (h *ReservaWebHandler) responderErrorBorrado(w http.ResponseWriter, r *http.Request, err error) {
	if esHTMX(r) {
		responderErrorHTMX(w, r, err, "#page-errors")
		return
	}

	h.renderReservas(w, r, statusDeError(err), mensajeDeError(err))
}

func buscarReservaDetalle(reservas []dominio.ReservaConDetalle, id int32) (dominio.ReservaConDetalle, bool) {
	for _, reserva := range reservas {
		if reserva.IDReserva == id {
			return reserva, true
		}
	}
	return dominio.ReservaConDetalle{}, false
}
