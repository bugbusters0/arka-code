<?php
// backend/models/ConceptoModel.php
require_once ROOT . '/backend/commons/BaseModel.php';

class ConceptoModel extends BaseModel
{

  public function getConceptos($search = '', $tipo = '')
  {
    $sql = "SELECT c.*, i.path as icono_path, i.nombre as icono_nombre,
                       (SELECT COUNT(*) FROM concepto_usuarios WHERE id_concepto = c.id) as usuarios_count,
                       (SELECT COUNT(*) FROM movimiento WHERE id_concepto = c.id) as movimientos_count
                FROM concepto c
                LEFT JOIN iconos i ON c.id_icono = i.id
                WHERE c.delete_at IS NULL";

    $params = [];

    if (!empty($search)) {
      $sql .= " AND c.nombre LIKE ?";
      $params[] = "%{$search}%";
    }

    if (!empty($tipo)) {
      $sql .= " AND c.tipo = ?";
      $params[] = $tipo;
    }

    $sql .= " ORDER BY c.nombre ASC";

    return $this->executeQuery($sql, $params);
  }

  public function getConceptoPorId($id)
  {
    $sql = "SELECT c.*, i.path as icono_path, i.nombre as icono_nombre
                FROM concepto c
                LEFT JOIN iconos i ON c.id_icono = i.id
                WHERE c.id = ? AND c.delete_at IS NULL";

    $result = $this->executeQuery($sql, [$id]);
    return !empty($result) ? $result[0] : null;
  }

  public function crearConcepto($nombre, $tipo, $color, $idIcono)
  {
    $sql = "INSERT INTO concepto (nombre, tipo, color, id_icono) 
                VALUES (?, ?, ?, ?)";

    return $this->executeNonQuery($sql, [$nombre, $tipo, $color, $idIcono]);
  }

  public function actualizarConcepto($id, $nombre, $color)
  {
    $sql = "UPDATE concepto 
                SET nombre = ?, color = ?, update_at = CURRENT_TIMESTAMP 
                WHERE id = ? AND delete_at IS NULL";

    return $this->executeNonQuery($sql, [$nombre, $color, $id]);
  }

  public function eliminarConcepto($id)
  {
    // Soft delete
    $sql = "UPDATE concepto 
                SET delete_at = CURRENT_TIMESTAMP 
                WHERE id = ?";

    return $this->executeNonQuery($sql, [$id]);
  }

  public function getIconosDisponibles()
  {
    $sql = "SELECT id, nombre, path FROM iconos WHERE delete_at IS NULL ORDER BY nombre";
    return $this->executeQuery($sql);
  }

  public function getConceptosPorUsuario($idUsuario)
  {
    $sql = "SELECT c.*, cu.desembolso_planejado, cu.limite_monto, cu.visible
                FROM concepto c
                INNER JOIN concepto_usuarios cu ON c.id = cu.id_concepto
                WHERE cu.id_usuario = ? AND c.delete_at IS NULL AND cu.visible = 1
                ORDER BY c.tipo, c.nombre";

    return $this->executeQuery($sql, [$idUsuario]);
  }
}
