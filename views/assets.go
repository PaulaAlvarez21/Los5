package views

import "embed"

// StaticFS contiene los estilos que se sirven desde /static/.
// Se embeben para que la imagen Docker no dependa de copiar la carpeta views.
//
//go:embed styles.css pico.min.css
var StaticFS embed.FS
