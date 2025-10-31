FROM golang:1.21-alpine

# Instalar MySQL en el mismo contenedor
RUN apk add --no-cache \
    mariadb \
    mariadb-client \
    bash

# Crear directorio de trabajo
WORKDIR /app

# Copiar go mod y descargar dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fuente
COPY . .

# Compilar la aplicación Go
RUN go build -o main .

# Script de inicio que inicia MySQL y luego la app Go
COPY start.sh /start.sh
RUN chmod +x /start.sh

# Exponer puerto
EXPOSE 8080

# Usar el mismo script de inicio que tenías
CMD ["/start.sh"]