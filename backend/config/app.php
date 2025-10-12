<?php
// Detectar el entorno
$isFly = getenv('FLY_APP_NAME') !== false;
$isLocal = !$isFly;

if ($isFly) {
  // Fly.io - usar la URL de Fly
  define('URLROOT', 'https://arka-code.fly.dev');
} else {
  // Desarrollo local
  define('URLROOT', 'http://localhost/arka-code');
}

define('APP_NAME', 'Arka App');
date_default_timezone_set('America/Lima');
