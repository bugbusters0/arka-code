<?php
// Detectar si estamos en Railway
$isRailway = getenv('RAILWAY_STATIC_URL') !== false;

if ($isRailway) {
  // Railway proporciona una URL dinámica
  $railwayUrl = getenv('RAILWAY_STATIC_URL') ?: 'https://tu-app.up.railway.app';
  define('URLROOT', $railwayUrl);
} else {
  // Desarrollo local
  define('URLROOT', 'http://localhost:8080');
}

define('APP_NAME', 'Arka App');
date_default_timezone_set('America/Lima');
