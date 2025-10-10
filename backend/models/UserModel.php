<?php
// backend/models/UserModel.php
require_once ROOT . '/backend/commons/BaseModel.php';
require_once ROOT . '/backend/entities/EntityUser.php';

class UserModel extends BaseModel
{

  public function verificarUsuario($email, $password)
  {
    // Primero obtiene la familia
    $sqlFamilia = "SELECT id, contrasena FROM familia WHERE correo = ? AND delete_at IS NULL";
    $familia = $this->executeQuery($sqlFamilia, [$email]);

    if (empty($familia)) {
      return false;
    }

    $familiaData = $familia[0];

    // hash contraseña y verificar
    if (!password_verify($password, $familiaData['contrasena'])) {
      return false;
    }

    // Obtiene el usuario admin de esa familia
    $sqlUsuario = "SELECT id, nombre, rol FROM usuario 
                      WHERE id_familia = ? AND rol = 'admin' AND delete_at IS NULL 
                      LIMIT 1";
    $usuario = $this->executeQuery($sqlUsuario, [$familiaData['id']]);

    if (!empty($usuario)) {
      $userData = $usuario[0];
      return new EntityUser(
        $userData['id'],
        $email,
        $familiaData['contrasena'],
        $userData['nombre'],
        $familiaData['id']
      );
    }

    return false;
  }

  public function crearUsuarioAdmin($idFamilia, $nombre, $fechaNac, $contraPersonalPlain)
  {
    $hashedContra = password_hash($contraPersonalPlain, PASSWORD_DEFAULT);

    $sql = "INSERT INTO usuario (id_familia, nombre, fecha_nac, rol, contra_personal) 
              VALUES (?, ?, ?, 'admin', ?)";

    if ($this->executeNonQuery($sql, [$idFamilia, $nombre, $fechaNac, $hashedContra])) {
      return $this->db->lastInsertId();
    }

    return false;
  }

  public function crearUsuarioMiembro($idFamilia, $nombre, $fechaNac, $contraPersonalPlain, $rol = 'miembro')
  {
    $hashedContra = password_hash($contraPersonalPlain, PASSWORD_DEFAULT);

    $sql = "INSERT INTO usuario (id_familia, nombre, fecha_nac, rol, contra_personal) 
              VALUES (?, ?, ?, ?, ?)";

    if ($this->executeNonQuery($sql, [$idFamilia, $nombre, $fechaNac, $rol, $hashedContra])) {
      return $this->db->lastInsertId();
    }

    return false;
  }

  public function getUsuarioPorId($idMember, $idFamily)
  {
    $sql = "SELECT u.id, u.nombre, u.rol, u.id_familia, f.correo 
              FROM usuario u 
              INNER JOIN familia f ON u.id_familia = f.id 
              WHERE u.id = ? AND f.id = ? AND u.delete_at IS NULL";

    $usuarios = $this->executeQuery($sql, [$idMember, $idFamily]);

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
  public function consultarExisteciaCredenciales($email, $tel)
  {
    $sql = "SELECT id From familia WHERE correo = ?";
    $user = $this->executeQuery($sql, [$email]);

    if (!empty($user)) {
      return $email;
    }

    $sql = "SELECT id From familia WHERE telefono = ?";
    $user = $this->executeQuery($sql, [$tel]);

    if (!empty($user)) {
      return $tel;
    }


    return false;
  }
  public function crearUsuario($correo, $telefono, $passwordPlain)
  {
    $hashedPassword = password_hash($passwordPlain, PASSWORD_DEFAULT);

    $sql = "INSERT INTO familia (correo, telefono, contrasena) VALUES (?, ?, ?)";

    if ($this->executeNonQuery($sql, [$correo, $telefono, $hashedPassword])) {
      return $this->db->lastInsertId();
    }

    return false;
  }

  public function getUsuarioPorCorreo($correo)
  {
    $sql = "SELECT id, correo, telefono FROM familia WHERE correo = ? AND delete_at IS NULL";
    $result = $this->executeQuery($sql, [$correo]);

    return !empty($result) ? $result[0] : null;
  }
}
