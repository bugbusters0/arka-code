#!/bin/bash
# Inicializar MySQL si no existe
if [ ! -d "/var/lib/mysql/mysql" ]; then
    echo "🔧 Inicializando base de datos MySQL..."
    mysql_install_db --user=mysql --datadir=/var/lib/mysql
fi

# Iniciar MySQL
mkdir -p /var/run/mysqld
chown mysql:mysql /var/run/mysqld

echo "🚀 Iniciando MySQL..."
mysqld_safe --datadir=/var/lib/mysql --socket=/var/run/mysqld/mysqld.sock &

# Esperar a que MySQL esté listo
echo "⏳ Esperando a que MySQL esté listo..."
for i in {1..30}; do
    mysqladmin ping --silent
    if [ $? -eq 0 ]; then
        echo "✅ MySQL listo y respondiendo"
        break
    fi
    echo "Intento $i/30 - Esperando a MySQL..."
    sleep 2
done

# Configurar base de datos y usuario
echo "🔧 Configurando base de datos..."
mysql -e "CREATE DATABASE IF NOT EXISTS arka_go;" || echo "BD ya existe"
mysql -e "CREATE USER IF NOT EXISTS 'root'@'localhost' IDENTIFIED BY 'root';" || echo "Usuario ya existe"
mysql -e "GRANT ALL PRIVILEGES ON arka_go.* TO 'root'@'localhost';" || echo "Permisos ya dados"
mysql -e "FLUSH PRIVILEGES;" || echo "Privilegios ya dados"

# Importar SQL si existe
if [ -f /app/sql.sql ]; then
    echo "📥 Importando estructura de base de datos..."
    mysql arka_go < /app/sql.sql && echo "✅ SQL importado" || echo "❌ Error importando SQL"
fi

echo "🎉 Configuración completada, iniciando aplicación Go..."
exec ./main