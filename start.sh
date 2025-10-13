#!/bin/bash
# Iniciar MySQL en modo seguro para contenedores
mkdir -p /var/run/mysqld
chown mysql:mysql /var/run/mysqld

# Iniciar MySQL en background
mysqld --user=mysql --datadir=/var/lib/mysql --socket=/var/run/mysqld/mysqld.sock &

# Esperar a que el socket esté disponible
echo "Esperando socket de MySQL..."
while [ ! -S /var/run/mysqld/mysqld.sock ]; do
    sleep 1
done
echo "✅ Socket de MySQL listo"

# Configurar base de datos
mysql -e "CREATE DATABASE IF NOT EXISTS arka_code;" || echo "BD ya existe"
mysql -e "CREATE USER IF NOT EXISTS 'arka_user'@'localhost' IDENTIFIED BY 'arka_password';" || echo "Usuario ya existe"
mysql -e "GRANT ALL PRIVILEGES ON arka_code.* TO 'arka_user'@'localhost';" || echo "Permisos ya dados"
mysql -e "FLUSH PRIVILEGES;"

# Importar SQL si existe
if [ -f /var/www/html/sql.sql ]; then
    echo "Importando base de datos..."
    mysql arka_code < /var/www/html/sql.sql && echo "✅ SQL importado" || echo "❌ Error importando SQL"
fi

echo "✅ MySQL configurado, iniciando Apache..."
exec apache2-foreground