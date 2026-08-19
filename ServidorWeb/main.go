    //package main: Le dice a Go que este codigo es un programa ejecutable (no una libreria reutilizable).
    package main

    import (
      "fmt"
      "net/http"
    )
    /*
    fmt: Para imprimir mensajes en la consola.
    net/http: La libreria nativa de Go para manejar todo el protocolo HTTP (servidores, peticiones, respuestas).*/

    func main() {
      // 1. Define el directorio que contiene los archivos estaticos.
      staticDir := "./static"

      // 2. Crea un manejador (handler) de servidor de archivos.
      fileServer := http.FileServer(http.Dir(staticDir))

      /*
      http.Dir(staticDir): Convierte la cadena "./static" en un sistema de archivos accesible mediante HTTP
      http.FileServer(...): Crea un manejador (handler) que toma la peticion HTTP entrante, busca el archivo 
      correspondiente en esa carpeta y lo devuelve automaticamente
      En Go, http.FileServer inspecciona automaticamente la extension del archivo que esta enviando. Cuando 
      detecta un .html, asigna de forma automatica la cabecera Content-Type: text/html; charset=utf-8 
      en la respuesta HTTP (y lo mismo hace para CSS, JS, imagenes, etc.)*/

      // 3. Registra el manejador para que atienda todas las peticiones ("/").
      http.Handle("/", fileServer)
      /*Usamos http.Handle porque fileServer es un http.Handler.
      Le dice al servidor: "Cualquier peticion que llegue a la raiz "/" (o sus subrutas), 
      entregasela a fileServer para que la responda".*/

      // 4. Define el puerto y muestra un mensaje.
      port := ":8080"
      fmt.Printf("Servidor ESTATICO escuchando en http://localhost%s\n", port)
      fmt.Printf("Sirviendo archivos desde: %s\n", staticDir)

      // 5. Inicia el servidor.
      err := http.ListenAndServe(port, nil)
      /*
      http.ListenAndServe(port, nil): Arranca el servidor web en el puerto :8080. Es una funcion bloqueante 
      (se queda corriendo indefinidamente esperando peticiones). El parametro nil le indica que use el enrutador por defecto donde registramos la ruta en el paso 4*/
      
      // 6. Manejo de errores
      if err != nil {
         fmt.Printf("Error al iniciar el servidor: %s\n", err)
      }
}
