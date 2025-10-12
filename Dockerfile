FROM php:8.2-apache

# Instalar MySQL en el mismo contenedor
RUN apt-get update && apt-get install -y default-mysql-server default-mysql-client

# Extensiones PHP
RUN docker-php-ext-install mysqli pdo pdo_mysql

# Habilitar mod_rewrite
RUN a2enmod rewrite

# Copiar configuración Apache
COPY apache-config/000-default.conf /etc/apache2/sites-available/000-default.conf

# Script para iniciar MySQL + Apache
COPY start.sh /start.sh
RUN chmod +x /start.sh

# Copiar aplicación
COPY . /var/www/html/

# Permisos
RUN chown -R www-data:www-data /var/www/html
RUN chmod -R 755 /var/www/html

EXPOSE 80
CMD ["/start.sh"]