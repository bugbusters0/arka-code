#!/bin/bash
# Iniciar MySQL
service mysql start

# Configurar base de datos
mysql -e "CREATE DATABASE IF NOT EXISTS arka_code;"
mysql -e "CREATE USER IF NOT EXISTS 'arka_user'@'localhost' IDENTIFIED BY 'arka_password';"
mysql -e "GRANT ALL PRIVILEGES ON arka_code.* TO 'arka_user'@'localhost';"
mysql -e "FLUSH PRIVILEGES;"

# Importar tu SQL inicial
if [ -f /var/www/html/sql.sql ]; then
    mysql arka_code < /var/www/html/sql.sql
fi

# Iniciar Apache
apache2-foreground