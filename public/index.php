<?php
require_once '../backend/config/init.php';

$uri = parse_url($_SERVER['REQUEST_URI'], PHP_URL_PATH);
$uri = str_replace('/arka-code', '', $uri); // Quita el prefijo
$uri = trim($uri, '/');

error_log("DEBUG URI LIMPIA: " . $uri);

//localhost/arka-code/login

switch ($uri) {
  case '':
  case 'index':
    header('Location: ' . URLROOT . '/login');
    exit;

  case 'login':
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
    if ($_SERVER['REQUEST_METHOD'] === 'POST') {
      require_once ROOT . '/backend/controllers/PerfilController.php';
      $controller = new PerfilController();
      $controller->consultarVerificacion($_POST["password"], $_POST["id"]);
    }

    require_once ROOT . '/backend/controllers/AuthController.php';
    $controller = new AuthController();
    if (isset($_GET['profile'])) {
      $controller->processSeleccionarPerfil();
    } else {
      $controller->solicitarPerfiles();
    }
    break;
  case 'registro':
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
      $controller->processRegister(); // ✅ CAMBIA: era processLogin()
    } else {
      $controller->showRegister();
    }
    break;
  case 'concepto/listar':
    require_once ROOT . '/backend/controllers/ConceptoController.php';
    $controller = new ConceptoController();
    $controller->listar();
    break;

  case 'concepto/crear':
    require_once ROOT . '/backend/controllers/ConceptoController.php';
    $controller = new ConceptoController();
    $controller->crear();
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
    if (preg_match('#^concepto/editar/(\d+)$#', $uri, $matches)) {
      require_once ROOT . '/backend/controllers/ConceptoController.php';
      $controller = new ConceptoController();
      $controller->editar($matches[1]);
      break;
    }

    if (preg_match('#^concepto/eliminar/(\d+)$#', $uri, $matches)) {
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
