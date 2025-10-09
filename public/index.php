<?php
require_once '../backend/config/init.php';

$uri = parse_url($_SERVER['REQUEST_URI'], PHP_URL_PATH);
$uri = str_replace('/arka-code', '', $uri); // Quita el prefijo
$uri = trim($uri, '/');

error_log("DEBUG URI LIMPIA: " . $uri);

switch ($uri) {
    case '':
    case 'index':
        header('Location: ' . URLROOT . '/login');
        exit;

    case 'login':
      if ($_SERVER['REQUEST_METHOD'] === 'POST') {
          require_once ROOT . '/backend/validators/LoginValidator.php';
          $validation = LoginValidator::validate($_POST);
          if (!$validation['success']) {
              $_SESSION['validation_errors'] = $validation['errors'];
              header('Location: ' . URLROOT . '/login');
              exit;
          }
          $_POST = array_merge($_POST, $validation['cleanData']);
      }

      require_once ROOT . '/backend/controllers/AuthController.php';
      $controller = new AuthController();
      if ($_SERVER['REQUEST_METHOD'] === 'POST') {
          $controller->processLogin();
      } else {
          $controller->showLogin();
      }
      break;
    case 'registro':
      if ($_SERVER['REQUEST_METHOD'] === 'POST') {
        require_once ROOT . '/backend/validators/RegisterValidator.php';
        $validation = RegisterValidator::validate($_POST);
        if (!$validation['success']) {
            $_SESSION['validation_errors'] = $validation['errors'];
            header('Location: ' . URLROOT . '/registro');
            exit;
        }
        $_POST = array_merge($_POST, $validation['cleanData']);
      }
      
      require_once ROOT . '/backend/controllers/AuthController.php';
      $controller = new AuthController();
      if ($_SERVER['REQUEST_METHOD'] === 'POST') {
          $controller->processRegister(); // ✅ CAMBIA: era processLogin()
      } else {
          $controller->showRegister();
      }
      break;
    case 'home':
        $controllerPath = ROOT . '/backend/controllers/HomeController.php';
        if (file_exists($controllerPath)) {
            require_once $controllerPath;
            $controller = new HomeController();
            $controller->index();
        } else {
            header('Location: ' . URLROOT . '/login');
            exit;
        }
        break;

    

    case 'logout':
        $controllerPath = ROOT . '/backend/controllers/AuthController.php';
        if (file_exists($controllerPath)) {
            require_once $controllerPath;
            $controller = new AuthController();
            $controller->logout();
        } else {
            header('Location: ' . URLROOT . '/login');
            exit;
        }
        break;

    case 'dashboard':
        $controllerPath = ROOT . '/backend/controllers/FinanzasController.php';
        if (file_exists($controllerPath)) {
            require_once $controllerPath;
            $controller = new FinanzasController();
            $controller->dashboard();
        } else {
            header('Location: ' . URLROOT . '/login');
            exit;
        }
        break;

    default:
        // Muestra 404 en vez de redirigir
        http_response_code(404);
        echo '404 - Página no encontrada: ' . htmlspecialchars($uri);
        exit;
}