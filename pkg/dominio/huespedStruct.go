package dominio

type Huesped struct {
	IDHuesped     int32   `json:"id_huesped"`
	IDReserva     int32   `json:"id_reserva"`
	Nombre        string  `json:"nombre"`
	Apellido      string  `json:"apellido"`
	Telefono      *string `json:"telefono"`
	Email         *string `json:"email"`
	Observaciones *string `json:"observaciones"`
}
