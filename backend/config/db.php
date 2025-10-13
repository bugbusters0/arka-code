<?php
// Configuración para diferentes entornos
if (getenv('FLY_APP_NAME') !== false) {
    // Fly.io - usar socket
    define('DB_HOST', 'localhost');
    define('DB_NAME', 'arka_code');
    define('DB_USER', 'arka_user');
    define('DB_PASS', 'arka_password');
} else {
    // Desarrollo local - MySQL local
    define('DB_HOST', 'localhost');
    define('DB_NAME', 'arka');
    define('DB_USER', 'root');
    define('DB_PASS', 'root'); // o tu password real
}

function conectarBD()
{
    try {
        if (getenv('FLY_APP_NAME') !== false) {
            // Fly.io - usar socket
            $dsn = "mysql:unix_socket=/var/run/mysqld/mysqld.sock;dbname=" . DB_NAME . ";charset=utf8mb4";
        } else {
            // Desarrollo local - TCP normal
            $dsn = "mysql:host=" . DB_HOST . ";dbname=" . DB_NAME . ";charset=utf8mb4";
        }

        $pdo = new PDO($dsn, DB_USER, DB_PASS);
        $pdo->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
        return $pdo;
    } catch (PDOException $e) {
        die("Error de conexión: " . $e->getMessage());
    }
}
