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
      $_SESSION['miembro_id'] = $verificacion->getId();  // ← Getter para ID
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
}
