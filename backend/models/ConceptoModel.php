<?php
// backend/models/ConceptoModel.php
require_once ROOT . '/backend/commons/BaseModel.php';

class ConceptoModel extends BaseModel
{
  public function getConceptosPorUsuario($idUsuario, $search = '', $tipo = '')
{
    $sql = "SELECT 
                c.id, 
                c.nombre, 
                c.tipo, 
                c.color,
                c.create_at,
                c.update_at,
                cu.desembolso_planejado, 
                cu.periodo_tipo AS desembolso_frecuencia, 
                cu.limite_monto, 
                cu.limite_tipo AS limite_frecuencia, 
                cu.visible,
                i.path AS icono_path, 
                i.nombre AS icono_nombre
            FROM concepto c
            INNER JOIN concepto_usuarios cu ON c.id = cu.id_concepto
            LEFT JOIN iconos i ON c.id_icono = i.id
            WHERE cu.id_usuario = ? 
                AND c.delete_at IS NULL 
                AND cu.delete_at IS NULL
                AND cu.visible = 1"; // ← Solo conceptos visibles para el usuario

    $params = [$idUsuario];

    if (!empty($search)) {
        $sql .= " AND c.nombre LIKE ?";
        $params[] = "%{$search}%";
    }

    if (!empty($tipo)) {
        $sql .= " AND c.tipo = ?";
        $params[] = $tipo;
    }

    $sql .= " ORDER BY c.nombre ASC";

    error_log("📋 Consulta conceptos usuario: " . $sql);
    error_log("📋 Parámetros: " . implode(', ', $params));
    
    $result = $this->executeQuery($sql, $params);
    error_log("📋 Conceptos encontrados: " . count($result));
    
    return $result;
}

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

    // public function actualizarConcepto($id, $nombre, $color)
    // {
    //     $sql = "UPDATE concepto 
    //             SET nombre = ?, color = ?, update_at = CURRENT_TIMESTAMP 
    //             WHERE id = ? AND delete_at IS NULL";

    //     return $this->executeNonQuery($sql, [$nombre, $color, $id]);
    // }

    public function eliminarConcepto($id)
    {
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
// public function crearRelacionUsuarioConcepto($idUsuario, $idConcepto, $desembolsoMonto, $desembolsoFrecuencia, $limiteMonto, $limiteFrecuencia)
//     {
//         $sql = "INSERT INTO concepto_usuarios 
//                 (id_usuario, id_concepto, desembolso_planejado, periodo_tipo, limite_monto, limite_tipo, visible) 
//                 VALUES (?, ?, ?, ?, ?, ?, 1)";
        
//         $params = [$idUsuario, $idConcepto, $desembolsoMonto, $desembolsoFrecuencia, $limiteMonto, $limiteFrecuencia];
        
//         return $this->executeNonQuery($sql, $params);
//     }

  public function getLastInsertId()
  {
      return $this->db->lastInsertId();
  }

  
public function actualizarConcepto($id, $nombre, $color, $idIcono)
{
    $sql = "UPDATE concepto 
            SET nombre = ?, color = ?, id_icono = ?, update_at = CURRENT_TIMESTAMP 
            WHERE id = ? AND delete_at IS NULL";

    return $this->executeNonQuery($sql, [$nombre, $color, $idIcono, $id]);
}

public function actualizarRelacionUsuarioConcepto($idUsuario, $idConcepto, $desembolsoMonto, $desembolsoFrecuencia, $limiteMonto, $limiteFrecuencia)
{
    $sql = "UPDATE concepto_usuarios 
            SET desembolso_planejado = ?, periodo_tipo = ?, limite_monto = ?, limite_tipo = ?, update_at = CURRENT_TIMESTAMP 
            WHERE id_usuario = ? AND id_concepto = ? AND delete_at IS NULL";

    $params = [$desembolsoMonto, $desembolsoFrecuencia, $limiteMonto, $limiteFrecuencia, $idUsuario, $idConcepto];
    
    return $this->executeNonQuery($sql, $params);
}

public function getConceptoPorIdYUsuario($idConcepto, $idUsuario)
{
    $sql = "SELECT 
                c.id, 
                c.nombre, 
                c.tipo, 
                c.color,
                c.id_icono,
                cu.desembolso_planejado, 
                cu.periodo_tipo AS desembolso_frecuencia, 
                cu.limite_monto, 
                cu.limite_tipo AS limite_frecuencia,
                i.nombre AS icono_nombre
            FROM concepto c
            INNER JOIN concepto_usuarios cu ON c.id = cu.id_concepto
            LEFT JOIN iconos i ON c.id_icono = i.id
            WHERE c.id = ? 
                AND cu.id_usuario = ?
                AND c.delete_at IS NULL 
                AND cu.delete_at IS NULL";

    $result = $this->executeQuery($sql, [$idConcepto, $idUsuario]);
    return !empty($result) ? $result[0] : null;
}

public function getUsuariosPorFamilia($idFamilia)
{
    $sql = "SELECT id FROM usuario WHERE id_familia = ? AND delete_at IS NULL";
    return $this->executeQuery($sql, [$idFamilia]);
}

public function crearRelacionUsuarioConcepto($idUsuario, $idConcepto, $desembolsoMonto, $desembolsoFrecuencia, $limiteMonto, $limiteFrecuencia)
{
    $sql = "INSERT INTO concepto_usuarios 
            (id_usuario, id_concepto, desembolso_planejado, periodo_tipo, limite_monto, limite_tipo, visible) 
            VALUES (?, ?, ?, ?, ?, ?, 1)";
    
    $params = [$idUsuario, $idConcepto, $desembolsoMonto, $desembolsoFrecuencia, $limiteMonto, $limiteFrecuencia];
    
    return $this->executeNonQuery($sql, $params);
}

public function crearRelacionUsuarioConceptoBasica($idUsuario, $idConcepto)
{
    $sql = "INSERT INTO concepto_usuarios 
            (id_usuario, id_concepto, desembolso_planejado, periodo_tipo, limite_monto, limite_tipo, visible) 
            VALUES (?, ?, 0, NULL, 0, NULL, 1)";
    
    $params = [$idUsuario, $idConcepto];
    
    return $this->executeNonQuery($sql, $params);
}
public function deshabilitarConceptoUsuario($idUsuario, $idConcepto)
{
    $sql = "UPDATE concepto_usuarios 
            SET visible = 0, update_at = CURRENT_TIMESTAMP 
            WHERE id_usuario = ? AND id_concepto = ? AND delete_at IS NULL";
    
    return $this->executeNonQuery($sql, [$idUsuario, $idConcepto]);
}

public function reactivarConceptoUsuario($idUsuario, $idConcepto)
{
    $sql = "UPDATE concepto_usuarios 
            SET visible = 1, update_at = CURRENT_TIMESTAMP 
            WHERE id_usuario = ? AND id_concepto = ? AND delete_at IS NULL";
    
    return $this->executeNonQuery($sql, [$idUsuario, $idConcepto]);
}
}