hacer con templ:

resrevas
departamento
huesped

en el index solo botones que llamen a esos renderizados.


y luego manejo de errores tambien se úede renderizar: si algo invalido que se renderice un cartelito de alerta.

tambien podrias hacer templ calendario:
con un boyton en index.html que lo llame y entonces se renderiza

javascript: evento on-click llama.

todo con template y javascrit (tp5) porque solo javascrit es tp4 y htmx no se va a pedir en la entrega.

conexion html y template:
en index.html ponemos {render.var} y toma ño que guaradamos y renderiza el contenido de var que va a ser la funcion definida en template que va a renderizarse.

Organización tp5:
- unir los repositorios
- agregar handler y logica general. Segun casos de uso

decisiones:
nos quedaron la API con los endpoint crud (solamente) en los handlers por entidad
y los casos de uso tienen solo handlers para web (cambia el formato de entrada y salida cómo se comunica con el cliente (depende de lo que recibe o dá el cliente))

