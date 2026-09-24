package dominio

type Departamento struct {
	IDDepto     int32   `json:"id_depto"`
	Nombre      string  `json:"nombre"`
	Direccion   string  `json:"direccion"`
	Disponible  bool    `json:"disponible"`
	Limpio      bool    `json:"limpio"`
	Descripcion *string `json:"descripcion"`
}
