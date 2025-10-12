<?php
// Detectar el entorno
$isFly = getenv('FLY_APP_NAME') !== false;
$isLocal = !$isFly;

if ($isFly) {
  // Fly.io - usar la URL de Fly
  define('URLROOT', 'https://arka-code-db.fly.dev');
} else {
  // Desarrollo local
  define('URLROOT', 'http://localhost:8080');
}

define('APP_NAME', 'Arka App');
date_default_timezone_set('America/Lima');
