# Etapa 1: compilacion de main.go con go
FROM golang:1.26-alpine AS builder 
# usamos la version 1.26 de go sobre alpine, 
#porque alpine es lo mas ligero posible y usamos esa version de go en el go.mod (la 1.26.5)

WORKDIR /app
#habla con docker y "reserva una parte" para la imagen para trabajar ahi

#copia y descarga dependencias de go
COPY go.mod ./
RUN go mod download

#copia nuestro proyecto dentro de docker
COPY . .

#lo compila y crea un ejecutable que llamaremos Los5
RUN CGO_ENABLED=0 go build -o Los5 .

# Etapa 2: imagen final (mucho mas pequeña)
#uso alpine que es como lo minimo requerido en la version que nos mostro la catedra
FROM alpine:3.20 

#crea un usuario sin permisos por seguridad llamado app
RUN adduser -D -u 1000 app
USER app

#especifico de nuevo para la etapa dos porque en main.go pusimos ./static en lugar de /static, el workdir es para asegurarnos que lo encuentra
WORKDIR /app

#copia desde la etapa anterior
COPY --from=builder /app/Los5 .
COPY --from=builder /app/static ./static

#indica que puerto, por defecto dejamos el 8080
EXPOSE 8080

#automatiza el arranque del contenedor y la ejecucion de main.go
CMD ["./Los5"]