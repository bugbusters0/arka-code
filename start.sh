#!/bin/bash
# Iniciar MySQL
mysqld_safe &

# Esperar a que MySQL esté listo
echo "Esperando a que MySQL inicie..."
for i in {1..30}; do
    mysql -e "status" > /dev/null 2>&1
    if [ $? -eq 0 ]; then
        break
    fi
    echo "Intento $i/30 - MySQL no está listo..."
    sleep 2
done

# Configurar base de datos
echo "Configurando base de datos..."
mysql -e "CREATE DATABASE IF NOT EXISTS arka_code;" || echo "Error creando BD"
mysql -e "CREATE USER IF NOT EXISTS 'arka_user'@'localhost' IDENTIFIED BY 'arka_password';" || echo "Error creando usuario"
mysql -e "GRANT ALL PRIVILEGES ON arka_code.* TO 'arka_user'@'localhost';" || echo "Error dando permisos"
mysql -e "FLUSH PRIVILEGES;" || echo "Error flush privileges"

# Importar SQL
if [ -f /var/www/html/sql.sql ]; then
    echo "Importando estructura de base de datos..."
    mysql arka_code < /var/www/html/sql.sql && echo "SQL importado correctamente" || echo "Error importando SQL"
else
    echo "AVISO: sql.sql no encontrado"
fi

echo "Iniciando Apache..."
exec apache2-foreground