<?php
// Configuración para Fly.io
if (getenv('FLY_APP_NAME') !== false) {
    // En Fly.io, MySQL está en otra app
    define('DB_HOST', 'fdaa:0:1234:a7b:abc:1234:5678:2'); // ← Usar IP interna
    define('DB_NAME', getenv('DB_NAME') ?: 'arka_code');
    define('DB_USER', getenv('DB_USER') ?: 'root');
    define('DB_PASS', getenv('DB_PASS') ?: 'password');
} else {
    // Desarrollo local con Docker
    define('DB_HOST', getenv('DB_HOST') ?: 'mysql');
    define('DB_NAME', getenv('DB_NAME') ?: 'arka_code');
    define('DB_USER', getenv('DB_USER') ?: 'root');
    define('DB_PASS', getenv('DB_PASS') ?: 'password');
}

function conectarBD()
{
    try {
        $pdo = new PDO("mysql:host=" . DB_HOST . ";dbname=" . DB_NAME, DB_USER, DB_PASS);
        $pdo->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
        return $pdo;
    } catch (PDOException $e) {
        die("Error de conexión: " . $e->getMessage());
    }
}
