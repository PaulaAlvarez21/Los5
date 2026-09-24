package dominio

import (
	"time"
)

type Reserva struct {
	IDReserva     int32     `json:"id_reserva"`
	FechaInicio   time.Time `json:"fecha_inicio"`
	IDDepto       int32     `json:"id_depto"`
	FechaFin      time.Time `json:"fecha_fin"`
	PrecioBase    string    `json:"precio_base"`
	CantNoches    int32     `json:"cant_noches"`
	Descuento     *string   `json:"descuento"`
	Observaciones *string   `json:"observaciones"`
}
