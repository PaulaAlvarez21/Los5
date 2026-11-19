package handlers

import (
	"net/http"

	"Los5/pkg/dominio"
	"Los5/pkg/logica"
	"Los5/views"
)

type DepartamentoWebHandler struct {
	departamentoService *logica.DepartamentoService
}

func NewDepartamentoWebHandler(departamentoService *logica.DepartamentoService) *DepartamentoWebHandler {
	return &DepartamentoWebHandler{departamentoService: departamentoService}
}

func (h *DepartamentoWebHandler) DepartamentosWeb(w http.ResponseWriter, r *http.Request) {
	h.renderDepartamentos(w, r, http.StatusOK, "")
}

func (h *DepartamentoWebHandler) renderDepartamentos(w http.ResponseWriter, r *http.Request, status int, mensaje string) {
	departamentos, err := h.departamentoService.ObtenerDepartamentos(r.Context())
	if err != nil {
		responderErrorWeb(w, r, err)
		return
	}

	renderHTML(w, r, status, views.DepartamentosPage(views.DepartamentosPageData{
		Departamentos: departamentos,
		Error:         mensaje,
	}))
}

// NuevoDepartamentoWeb devuelve únicamente el modal con el formulario de alta.
func (h *DepartamentoWebHandler) NuevoDepartamentoWeb(w http.ResponseWriter, r *http.Request) {
	renderHTML(w, r, http.StatusOK, views.DepartamentoFormModal(""))
}

func (h *DepartamentoWebHandler) CrearDepartamentoWeb(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.responderErrorAlta(w, r, http.StatusBadRequest, "formulario inválido")
		return
	}

	departamento := dominio.Departamento{
		Nombre:      r.FormValue("nombre"),
		Direccion:   r.FormValue("direccion"),
		Disponible:  parseBoolForm(r.FormValue("disponible")),
		Limpio:      parseBoolForm(r.FormValue("limpio")),
		Descripcion: stringOpcional(r.FormValue("descripcion")),
	}

	creado, err := h.departamentoService.CrearDepartamento(r.Context(), departamento)
	if err != nil {
		h.responderErrorAlta(w, r, statusDeError(err), mensajeDeError(err))
		return
	}

	if esHTMX(r) {
		departamentos, err := h.departamentoService.ObtenerDepartamentos(r.Context())
		if err != nil {
			h.responderErrorAlta(w, r, statusDeError(err), mensajeDeError(err))
			return
		}

		renderHTML(w, r, http.StatusOK, views.DepartamentoCreado(creado, len(departamentos) == 1))
		return
	}

	http.Redirect(w, r, "/departamentos", http.StatusSeeOther)
}

// EditarDepartamentoWeb devuelve el modal con el formulario precargado.
func (h *DepartamentoWebHandler) EditarDepartamentoWeb(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		h.renderDepartamentos(w, r, statusDeError(err), mensajeDeError(err))
		return
	}

	departamento, err := h.departamentoService.ObtenerDepartamento(r.Context(), id)
	if err != nil {
		h.renderDepartamentos(w, r, statusDeError(err), mensajeDeError(err))
		return
	}

	renderHTML(w, r, http.StatusOK, views.DepartamentoEditarModal(departamento))
}

func (h *DepartamentoWebHandler) ActualizarDepartamentoWeb(w http.ResponseWriter, r *http.Request) {
	fallo := func(err error) {
		if esHTMX(r) {
			responderErrorHTMX(w, r, err, "#page-errors")
			return
		}
		h.renderDepartamentos(w, r, statusDeError(err), mensajeDeError(err))
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

	departamento := dominio.Departamento{
		IDDepto:     id,
		Nombre:      r.FormValue("nombre"),
		Direccion:   r.FormValue("direccion"),
		Disponible:  parseBoolForm(r.FormValue("disponible")),
		Limpio:      parseBoolForm(r.FormValue("limpio")),
		Descripcion: stringOpcional(r.FormValue("descripcion")),
	}

	if err := h.departamentoService.ActualizarDepartamento(r.Context(), departamento); err != nil {
		fallo(err)
		return
	}

	if esHTMX(r) {
		renderHTML(w, r, http.StatusOK, views.DepartamentoEditado(departamento))
		return
	}

	http.Redirect(w, r, "/departamentos", http.StatusSeeOther)
}

func (h *DepartamentoWebHandler) CambiarHigieneWeb(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		h.renderDepartamentos(w, r, statusDeError(err), mensajeDeError(err))
		return
	}

	departamento, err := h.departamentoService.ObtenerDepartamento(r.Context(), id)
	if err != nil {
		h.renderDepartamentos(w, r, statusDeError(err), mensajeDeError(err))
		return
	}

	departamento.Limpio = !departamento.Limpio

	if err := h.departamentoService.ActualizarDepartamento(r.Context(), departamento); err != nil {
		h.renderDepartamentos(w, r, statusDeError(err), mensajeDeError(err))
		return
	}

	http.Redirect(w, r, "/departamentos", http.StatusSeeOther)
}

func (h *DepartamentoWebHandler) EliminarDepartamentoWeb(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		h.responderErrorBorrado(w, r, err)
		return
	}

	if err := h.departamentoService.EliminarDepartamento(r.Context(), id); err != nil {
		h.responderErrorBorrado(w, r, err)
		return
	}

	if esHTMX(r) {
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/departamentos", http.StatusSeeOther)
}

// responderErrorAlta conserva el modal abierto y muestra el error adentro.
func (h *DepartamentoWebHandler) responderErrorAlta(w http.ResponseWriter, r *http.Request, status int, mensaje string) {
	if esHTMX(r) {
		retarget(w, "#modal-root")
		renderHTML(w, r, http.StatusOK, views.DepartamentoFormModal(mensaje))
		return
	}

	h.renderDepartamentos(w, r, status, mensaje)
}

func (h *DepartamentoWebHandler) responderErrorBorrado(w http.ResponseWriter, r *http.Request, err error) {
	if esHTMX(r) {
		responderErrorHTMX(w, r, err, "#page-errors")
		return
	}

	h.renderDepartamentos(w, r, statusDeError(err), mensajeDeError(err))
}
