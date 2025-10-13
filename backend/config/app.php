<?php
// Detectar el entorno
$isFly = getenv('FLY_APP_NAME') !== false;

if ($isFly) {
  // Fly.io - raíz del dominio
  define('URLROOT', 'https://arka-code.fly.dev');
} else {
  // Desarrollo local -agregar el nombre de la carpeta
  define('URLROOT', 'http://localhost');
}

define('APP_NAME', 'Arka App');
date_default_timezone_set('America/Lima');
