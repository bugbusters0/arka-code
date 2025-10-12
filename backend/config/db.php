<?php
// Para Fly.io con MySQL en el mismo contenedor
define('DB_HOST', 'localhost');  // ← Ahora es localhost
define('DB_NAME', 'arka_code');
define('DB_USER', 'arka_user');
define('DB_PASS', 'arka_password');

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
