<?php
// Detectar el entorno
$isFly = getenv('FLY_APP_NAME') !== false;

if ($isFly) {
  // Fly.io
  define('URLROOT', 'https://arka-code.fly.dev');
} else {
  // Desarrollo local - SIN subcarpeta
  define('URLROOT', 'http://localhost/arka-code');
}

define('APP_NAME', 'Arka App');
date_default_timezone_set('America/Lima');
