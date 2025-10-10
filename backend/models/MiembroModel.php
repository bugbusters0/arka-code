<?php
// backend/models/FamiliaModel.php
require_once ROOT . '/backend/commons/BaseModel.php';
require_once ROOT . '/backend/entities/EntityMiembro.php';

class MiembroModel extends BaseModel
{

  public function crearFamilia($correo, $telefono, $passwordPlain)
  {
    $hashedPassword = password_hash($passwordPlain, PASSWORD_DEFAULT);

    $sql = "INSERT INTO familia (correo, telefono, contrasena) VALUES (?, ?, ?)";

    if ($this->executeNonQuery($sql, [$correo, $telefono, $hashedPassword])) {
      return $this->db->lastInsertId();
    }

    return false;
  }

  public function getFamiliaPorCorreo($correo)
  {
    $sql = "SELECT id, correo, telefono FROM familia WHERE correo = ? AND delete_at IS NULL";
    $result = $this->executeQuery($sql, [$correo]);

    return !empty($result) ? $result[0] : null;
  }

  public function verificarFamilia($correo, $password)
  {
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

  public function getMiembrosPorFamilia($idFamilia)
  {
    $sql = "SELECT id, nombre, fecha_nac, rol, contra_personal FROM usuario 
              WHERE id_familia = ? AND delete_at IS NULL";
    $result = $this->executeQuery($sql, [$idFamilia]);
    return $result;
  }

  public function getMiembroPorId($id)
  {
    $sql = "SELECT u.id, u.nombre, u.rol, u.id_familia, f.correo 
              FROM usuario u 
              INNER JOIN familia f ON u.id_familia = f.id 
              WHERE u.id = ? AND u.delete_at IS NULL";

    $usuarios = $this->executeQuery($sql, [$id]);

    if (!empty($usuarios)) {
      $row = $usuarios[0];
      return new EntityUser(
        $row['id'],
        $row['correo'],
        null,
        $row['nombre'],
        $row['id_familia']
      );
    }

    return null;
  }
  public function verificarMiembro($pass, $id)
  {
    $sql = "SELECT id, id_familia, nombre, fecha_nac, rol, contra_personal FROM usuario WHERE id = ? AND delete_at IS NULL";
    $miembro = $this->executeQuery($sql, [$id]);
    if (empty($miembro)) {
      return false;
    }

    $data = $miembro[0];

    // hash contraseña y verificar
    if (!password_verify($pass, $data['contra_personal'])) {
      return false;
    }


    return new EntityMiembro(
      $data['id'],
      $data['id_familia'],
      $data['nombre'],
      $data['fecha_nac'],
      $data['rol'],
      "",
    );
  }
}
