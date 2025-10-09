<?php
// backend/models/FamiliaModel.php
require_once ROOT . '/backend/commons/BaseModel.php';

class FamiliaModel extends BaseModel {
  
  public function crearFamilia($correo, $telefono, $passwordPlain) {
    $hashedPassword = password_hash($passwordPlain, PASSWORD_DEFAULT);
    
    $sql = "INSERT INTO familia (correo, telefono, contrasena) VALUES (?, ?, ?)";
    
    if ($this->executeNonQuery($sql, [$correo, $telefono, $hashedPassword])) {
      return $this->db->lastInsertId();
    }
    
    return false;
  }
  
  public function getFamiliaPorCorreo($correo) {
    $sql = "SELECT id, correo, telefono FROM familia WHERE correo = ? AND delete_at IS NULL";
    $result = $this->executeQuery($sql, [$correo]);
    
    return !empty($result) ? $result[0] : null;
  }
  
  public function verificarFamilia($correo, $password) {
    $sql = "SELECT id, correo, contrasena FROM familia WHERE correo = ? AND delete_at IS NULL";
    $result = $this->executeQuery($sql, [$correo]);
    
    if (!empty($result)) {
      $familia = $result[0];
      if (password_verify($password, $familia['contrasena'])) {
        return $familia;
      }
    }
    
    return false;
  }
}
?>