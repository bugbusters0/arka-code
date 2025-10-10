<?php
// backend/commons/BaseController.php

class BaseController
{

  public function __construct()
  {
    // Inicialización común
    if (session_status() === PHP_SESSION_NONE) {
      session_start();
    }
  }

  /**
   * Renderiza una vista sin layout
   */
  protected function view($vista, $data = [])
  {
    extract($data);

    $vistaPath = ROOT . '/frontend/views/' . str_replace('/', DIRECTORY_SEPARATOR, $vista) . '.php';
    if (file_exists($vistaPath)) {
      require_once $vistaPath;
    } else {
      die("Vista no encontrada: $vista en $vistaPath");
    }
  }

  /**
   * Renderiza una vista con layout
   * Automáticamente busca CSS/JS con el mismo nombre de la vista
   * 
   * Ejemplo: 
   * viewWithLayout('concepto/gasto', 'main', $data)
   * Buscará: concepto-listar.css y concepto-listar.js
   */
  protected function viewWithLayout($vista, $layout = 'main', $data = [])
  {
    // Convertir ruta de vista a nombre de archivo: concepto/gasto -> concepto-listar
    $assetName = str_replace('/', '-', $vista);

    // Agregar rutas de assets automáticamente si no están definidas
    if (!isset($data['pageCSS'])) {
      $data['pageCSS'] = $assetName;
    }
    if (!isset($data['pageJS'])) {
      $data['pageJS'] = $assetName;
    }

    extract($data);

    // Captura el contenido de la vista
    ob_start();
    $vistaPath = ROOT . '/frontend/views/' . str_replace('/', DIRECTORY_SEPARATOR, $vista) . '.php';
    if (file_exists($vistaPath)) {
      require_once $vistaPath;
    } else {
      die("Vista no encontrada: $vista");
    }
    $content = ob_get_clean();

    // Renderiza el layout con el contenido
    $layoutPath = ROOT . '/frontend/views/layouts/' . str_replace('/', DIRECTORY_SEPARATOR, $layout) . '.php';
    if (file_exists($layoutPath)) {
      require_once $layoutPath;
    } else {
      die("Layout no encontrado: $layout");
    }
  }

  protected function addAsset($type, $file)
  {
    $baseUrl = URLROOT . '/frontend/assets';
    if ($type === 'css') {
      echo "<link rel='stylesheet' href='$baseUrl/css/$file.css'>";
    } elseif ($type === 'js') {
      echo "<script src='$baseUrl/js/$file.js'></script>";
    }
  }

  protected function requireAuth()
  {
    if (!isset($_SESSION['user_id'])) {
      header('Location: ' . URLROOT . '/login');
      exit;
    }
  }
}
