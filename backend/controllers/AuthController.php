<?php
// backend/controllers/AuthController.php
require_once ROOT . '/backend/commons/BaseController.php';
require_once ROOT . '/backend/models/UserModel.php';
require_once ROOT . '/backend/models/MiembroModel.php';

class AuthController extends BaseController
{
  private $userModel;
  private $miembroModel;

  public function __construct()
  {
    parent::__construct();
    $this->userModel = new UserModel();
    $this->miembroModel = new MiembroModel();
  }

  public function showLogin()
  {
    $data['title'] = 'Login - Arka App';
    $data['validation_errors'] = $_SESSION['validation_errors'] ?? [];
    $data['success'] = $_SESSION['success_message'] ?? '';
    unset($_SESSION['validation_errors']);
    unset($_SESSION['success_message']);

    $data['error'] = '';
    $this->view('auth/login', $data);
  }

  // ← Actualizado: Consulta DB y pasa profiles a vista
  public function solicitarPerfiles($error = '')
  {
    if (!isset($_SESSION['user_id'])) {
      header('Location: ' . URLROOT . '/login');
      exit;
    }

    $userId = $_SESSION['user_id'];
    $profiles = $this->miembroModel->getMiembrosPorFamilia($userId);
    foreach ($profiles as &$profile) {
      if (!empty($profile['contra_personal'])) {
        $profile['protegido'] = true;
      } else {
        $profile['protegido'] = false;
      }
      $profile['contra_personal'] = '';
    }
    unset($profile);

    $data['title'] = 'Seleccionar Perfil - Arka App';
    $data['profiles'] = $profiles;  // ← Array de usuarios miembros para loop en vista
    $data['validation_errors'] = $_SESSION['validation_errors'] ?? [];
    $data['success'] = $_SESSION['success_message'] ?? '';
    unset($_SESSION['validation_errors'], $_SESSION['success_message']);
    $data['error'] = $error;
    $this->view('auth/seleccionarperfil', $data);
  }

  public function showRegister()
  {
    $data['title'] = 'Registro - Arka App';
    $data['validation_errors'] = $_SESSION['validation_errors'] ?? [];
    $data['success'] = $_SESSION['success_message'] ?? '';
    unset($_SESSION['validation_errors'], $_SESSION['success_message']);
    $data['error'] = '';
    $this->view('auth/registrar', $data);
  }

  public function enviarCredenciales()
  {
    $email = $_POST['email'] ?? '';
    $contrasena = $_POST['contrasena'] ?? '';

    $userEntity = $this->userModel->verificarUsuario($email, $contrasena);
    if ($userEntity) {
      $_SESSION['user_id'] = $userEntity->getId();
      $_SESSION['user_email'] = $userEntity->getEmail();
      header('Location: ' . URLROOT . '/seleccionar-perfil');
      exit;
    } else {
      $data['error'] = 'Credenciales inválidas';
      $data['title'] = 'Login - Arka App';
      $this->view('auth/login', $data);
    }
  }
  public function processSeleccionarPerfil()
  {
    if (!isset($_GET['profile']) || !is_numeric($_GET['profile'])) {
      $error  = 'Perfil inválido';
      $this->solicitarPerfiles($error);  // Vuelve a vista con error
      return;
    }

    $profileId = (int) $_GET['profile'];
    $_SESSION['member_id'] = $profileId;  // ← Guarda ID seleccionado en sesión
    $userId = $_SESSION['user_id'];
    // Opcional: Verifica que el perfil pertenezca a la familia (seguridad)
    $user = $this->userModel->getUsuarioPorId($profileId, $userId);
    if (!$user) {
      $data['error'] = 'Perfil no válido para esta familia';
      $this->solicitarPerfiles();
      return;
    }

    header('Location: ' . URLROOT . '/dashboard');  // ← Redirige a dashboard
    exit;
  }

  public function processRegister()
  {
    // Los datos ya vienen validados y limpios desde index.php
    try {
      // 1. Crear familia
      $userId = $this->userModel->crearUsuario(
        $_POST['email'],
        $_POST['telefono'],
        $_POST['password']
      );

      if (!$userId) {
        throw new Exception('Error al crear la familia');
      }

      // 2. Crear usuario administrador
      $adminId = $this->userModel->crearUsuarioAdmin(
        $userId,
        $_POST['adminNombre'],
        $_POST['adminNacimiento'],
        $_POST['contraPersonal']
      );

      if (!$adminId) {
        throw new Exception('Error al crear el administrador');
      }

      // 3. Crear miembros adicionales (si existen)
      if (!empty($_POST['miembros'])) {
        foreach ($_POST['miembros'] as $miembro) {
          $this->userModel->crearUsuarioMiembro(
            $userId,
            $miembro['nombre'],
            $miembro['nacimiento'],
            $miembro['contra'],
            $miembro['rol']
          );
        }
      }

      // 4. Éxito: redirigir a login
      $_SESSION['success_message'] = 'Cuenta creada exitosamente. ¡Inicia sesión!';
      header('Location: ' . URLROOT . '/login');
      exit;
    } catch (Exception $e) {
      // Error: mostrar en formulario
      error_log("Error en registro: " . $e->getMessage());
      $data['error'] = 'Ocurrió un error al crear la cuenta. Intenta nuevamente.';
      $data['title'] = 'Registro - Arka App';
      $this->view('auth/registrar', $data);
    }
  }

  public function logout()
  {
    session_destroy();
    header('Location: ' . URLROOT . '/login');
    exit;
  }
}
