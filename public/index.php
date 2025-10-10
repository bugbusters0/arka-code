<?php
require_once '../backend/config/init.php';

$uri = parse_url($_SERVER['REQUEST_URI'], PHP_URL_PATH);
$uri = str_replace('/arka-code', '', $uri);
$uri = trim($uri, '/');

error_log("DEBUG URI LIMPIA: " . $uri);

// Cargar middleware
require_once ROOT . '/backend/commons/AuthMiddleware.php';

switch ($uri) {
  case '/':
    header('Location: ' . URLROOT . '/login');
    exit;
    break;
  case 'index':
    header('Location: ' . URLROOT . '/login');
    exit;
    break;

  case 'login':

    AuthMiddleware::guest();

    if ($_SERVER['REQUEST_METHOD'] === 'POST') {
      require_once ROOT . '/backend/validator/LoginValidator.php';
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
      $controller->enviarCredenciales();
    } else {
      $controller->showLogin();
    }
    break;

  case 'seleccionar-perfil':
    AuthMiddleware::familyOnly();
    if (isset($_SESSION['miembro_id'])) {
      header('Location: ' . URLROOT . '/concepto/listar');
      exit;
    }

    require_once ROOT . '/backend/controllers/AuthController.php';
    require_once ROOT . '/backend/controllers/PerfilController.php';

    if ($_SERVER['REQUEST_METHOD'] === 'POST') {
      $controller = new PerfilController();
      $controller->consultarVerificacion($_POST["password"], $_POST["idMiembro"]);
    } elseif (isset($_GET['profile'])) {
      $controller = new AuthController();
      $controller->processSeleccionarPerfil();
    } else {
      $controller = new PerfilController();
      $controller->solicitarPerfiles();
    }
    break;

  case 'registro':
    AuthMiddleware::guest();

    if ($_SERVER['REQUEST_METHOD'] === 'POST') {
      require_once ROOT . '/backend/validator/RegisterValidator.php';
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
      $controller->processRegister();
    } else {
      $controller->showRegister();
    }
    break;


  case 'concepto/listar':
    AuthMiddleware::auth();

    require_once ROOT . '/backend/controllers/ConceptoController.php';
    $controller = new ConceptoController();
    $controller->listar();
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

  default:
    // ✅ PROTECCIÓN: Para rutas dinámicas también
    if (preg_match('#^concepto/editar/(\d+)$#', $uri, $matches)) {
      AuthMiddleware::auth();
      require_once ROOT . '/backend/controllers/ConceptoController.php';
      $controller = new ConceptoController();
      $controller->editar($matches[1]);
      break;
    }

    if (preg_match('#^concepto/eliminar/(\d+)$#', $uri, $matches)) {
      AuthMiddleware::auth();
      require_once ROOT . '/backend/controllers/ConceptoController.php';
      $controller = new ConceptoController();
      $controller->eliminar($matches[1]);
      break;
    }

    // Muestra 404 en vez de redirigir
    http_response_code(404);
    echo '404 - Página no encontrada: ' . htmlspecialchars($uri);
    exit;
}
