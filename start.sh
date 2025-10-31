#!/bin/bash
# Iniciar MySQL en modo seguro para contenedores
mkdir -p /var/run/mysqld
chown mysql:mysql /var/run/mysqld

# Iniciar MySQL en background y mantenerlo corriendo
mysqld --user=mysql --datadir=/var/lib/mysql --socket=/var/run/mysqld/mysqld.sock &

# Esperar a que MySQL esté completamente listo
echo "Esperando a que MySQL esté listo..."
for i in {1..30}; do
    mysql -e "SELECT 1;" > /dev/null 2>&1
    if [ $? -eq 0 ]; then
        echo "✅ MySQL completamente inicializado"
        break
    fi
    echo "Intento $i/30 - Esperando a MySQL..."
    sleep 2
done

# Configurar base de datos solo si no existe
mysql -e "CREATE DATABASE IF NOT EXISTS arka_code;"
mysql -e "CREATE USER IF NOT EXISTS 'arka_user'@'localhost' IDENTIFIED BY 'arka_password';"
mysql -e "GRANT ALL PRIVILEGES ON arka_code.* TO 'arka_user'@'localhost';"
mysql -e "FLUSH PRIVILEGES;"

# Importar SQL si existe
if [ -f /app/sql.sql ]; then
    echo "Importando base de datos..."
    mysql arka_code < /app/sql.sql && echo "✅ SQL importado" || echo "❌ Error importando SQL"
fi

echo "✅ MySQL configurado e iniciado"
echo "✅ Iniciando aplicación Go..."

# Ejecutar la aplicación Go (MySQL sigue corriendo en background)
exec ./main