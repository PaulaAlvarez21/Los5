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

func (h *DepartamentoWebHandler) CrearDepartamentoWeb(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderDepartamentos(w, r, http.StatusBadRequest, "formulario inválido")
		return
	}

	departamento := dominio.Departamento{
		Nombre:      r.FormValue("nombre"),
		Direccion:   r.FormValue("direccion"),
		Disponible:  parseBoolForm(r.FormValue("disponible")),
		Limpio:      parseBoolForm(r.FormValue("limpio")),
		Descripcion: stringOpcional(r.FormValue("descripcion")),
	}

	if _, err := h.departamentoService.CrearDepartamento(r.Context(), departamento); err != nil {
		h.renderDepartamentos(w, r, statusDeError(err), mensajeDeError(err))
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
		h.renderDepartamentos(w, r, statusDeError(err), mensajeDeError(err))
		return
	}

	if err := h.departamentoService.EliminarDepartamento(r.Context(), id); err != nil {
		h.renderDepartamentos(w, r, statusDeError(err), mensajeDeError(err))
		return
	}

	http.Redirect(w, r, "/departamentos", http.StatusSeeOther)
}
