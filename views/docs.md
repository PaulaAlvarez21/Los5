en huespedes.templ

huespedes := views.HusspedesPage(...)
struct Datos map[string]any{
	"Huespedes": huespedes 
}


y despeues
			{{render .Huespedes}}