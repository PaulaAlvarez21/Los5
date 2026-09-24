package dominio

import "net/http"

type AppError struct {
	Status int
	Msg    string
}

func (e *AppError) Error() string {
	return e.Msg
}

func NotFoundError() *AppError { //404
	return &AppError{Status: http.StatusNotFound, Msg: "recurso no encontrado"}
}

func ValidationError(msg string) *AppError { //400
	return &AppError{Status: http.StatusBadRequest, Msg: msg}
}
