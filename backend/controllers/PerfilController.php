<?php
// backend/controllers/PerfilController.php
require_once ROOT . '/backend/commons/BaseController.php';
require_once ROOT . '/backend/models/UserModel.php';
require_once ROOT . '/backend/models/MiembroModel.php';

class PerfilController extends BaseController
{
  private $userModel;
  private $miembroModel;

  public function __construct()
  {
    parent::__construct();
    $this->requireAuth();
    $this->userModel = new UserModel();
    $this->miembroModel = new MiembroModel();
  }

  public function consultarVerificacion($contra, $idPerfil)
  {
    $verificacion = $this->miembroModel->verificarMiembro($contra, $idPerfil);
    if ($verificacion) {
      $_SESSION['miembro_id'] = $verificacion->getId();
      $_SESSION['miembro_idFamilia'] = $verificacion->getIdFamilia();
      $_SESSION['miembro_rol'] = $verificacion->getRol();
      $_SESSION['miembro_nombre'] = $verificacion->getNombre();
      header('Location: ' . URLROOT . '/concepto/listar');
      exit;
    } else {

      $_SESSION['error_message'] = 'Contraseña incorrecta';
      header('Location: ' . URLROOT . '/seleccionar-perfil');
      exit;
    }
  }

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
    $data['profiles'] = $profiles;
    $data['validation_errors'] = $_SESSION['validation_errors'] ?? [];
    $data['success'] = $_SESSION['success_message'] ?? '';
    unset($_SESSION['validation_errors'], $_SESSION['success_message']);
    $data['error'] = $error;
    $this->view('auth/seleccionarperfil', $data);
  }
}
