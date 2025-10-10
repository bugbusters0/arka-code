<?php
// backend/controllers/AuthController.php
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
    $this->userModel = new UserModel();
    $this->miembroModel = new MiembroModel();
  }

  public function consultarVerificacion($contra, $idPerfil)
  {
    $verificacion = $this->miembroModel->verificarMiembro($contra, $idPerfil);

    if ($verificacion) {
      $data['title'] = 'Conceptos Perfil - Arka App';
      $this->view('concepto/listar', $data);
    } else {
      $data['validation_errors'] = "Contraseña inválida";
      $this->view('auth/seleccionarperfil', $data);
    }
  }
}
