FROM golang:1.21.3-alpine

# Establecer el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiar los archivos necesarios al contenedor
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copiar el código fuente al contenedor
COPY . .

# Compilar la aplicación
RUN go build -o main cmd/api/main.go

# Exponer el puerto por el que la aplicación será accesible
EXPOSE 3000

# Ejecutar la aplicación
CMD ["./main"]
