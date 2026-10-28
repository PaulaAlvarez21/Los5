package views

func valorTexto(valor *string) string {
	if valor == nil {
		return ""
	}
	return *valor
}
