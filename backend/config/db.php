<?php
// Configuración para diferentes entornos
if (getenv('FLY_APP_NAME') !== false) {
    // Fly.io - usar socket
    define('DB_HOST', 'localhost');
    define('DB_NAME', 'arka_code');
    define('DB_USER', 'arka_user');
    define('DB_PASS', 'arka_password');
    define('DB_SOCKET', '/var/run/mysqld/mysqld.sock');
} else {
    // Desarrollo local
    define('DB_HOST', 'mysql');
    define('DB_NAME', 'arka_code');
    define('DB_USER', 'root');
    define('DB_PASS', 'password');
    define('DB_SOCKET', null);
}

function conectarBD()
{
    try {
        if (defined('DB_SOCKET') && DB_SOCKET) {
            // Usar socket en Fly.io
            $dsn = "mysql:unix_socket=" . DB_SOCKET . ";dbname=" . DB_NAME . ";charset=utf8mb4";
        } else {
            // Usar TCP en desarrollo local
            $dsn = "mysql:host=" . DB_HOST . ";dbname=" . DB_NAME . ";charset=utf8mb4";
        }

        $pdo = new PDO($dsn, DB_USER, DB_PASS);
        $pdo->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
        return $pdo;
    } catch (PDOException $e) {
        die("Error de conexión: " . $e->getMessage());
    }
}
