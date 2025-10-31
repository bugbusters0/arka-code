FROM golang:1.25-alpine

# Instalar MySQL y configurar inicialización
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

# Script de inicio mejorado
COPY start.sh /start.sh
RUN chmod +x /start.sh

# Inicializar base de datos MySQL
RUN mysql_install_db --user=mysql --datadir=/var/lib/mysql

EXPOSE 8080
CMD ["/start.sh"]