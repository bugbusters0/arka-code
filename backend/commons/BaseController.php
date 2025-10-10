<?php

class BaseController {
    public function __construct() {
        // Inicializaciones comunes si las necesitas después
    }
    protected function view($vista, $data = []) {
        extract($data); 
        
        $vistaPath = ROOT . '/frontend/views/' . str_replace('/', DIRECTORY_SEPARATOR, $vista) . '.php';
        if (file_exists($vistaPath)) {
          require_once $vistaPath;
        } else {
          die("Vista no encontrada: $vista");
        }
    }

   
    protected function layout($layout, $content, $data = []) {
        $data['content'] = $content;
        extract($data);
        
        $layoutPath = ROOT . '/frontend/views/layouts/' . str_replace('/', DIRECTORY_SEPARATOR, $layout) . '.php';
        if (file_exists($layoutPath)) {
            require_once $layoutPath;
        } else {
            die("Layout no encontrado: $layout");
        }
    }

   
    protected function addAsset($type, $file) {
        $baseUrl = URLROOT . '/frontend'; 
        if ($type === 'css') {
            echo "<link rel='stylesheet' href='$baseUrl/css/$file.css'>";
        } elseif ($type === 'js') {
            echo "<script src='$baseUrl/js/$file.js'></script>";
        }
    }

    protected function requireAuth() {
        if (!isset($_SESSION['user_id'])) {
            header('Location: /login');
            exit;
        }
    }
}
?>