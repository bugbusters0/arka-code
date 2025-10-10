<?php
// backend/models/ConceptoModel.php
require_once ROOT . '/backend/commons/BaseModel.php';
require_once ROOT . '/backend/entities/EntityConcepto.php';
require_once ROOT . '/backend/entities/EntityConceptoUsuario.php';
require_once ROOT . '/backend/entities/EntityIcono.php';

class ConceptoModel extends BaseModel
{
	public function getConceptosPorUsuario($idUsuario, $tipo = 'gasto')
	{
		$sql = "SELECT 
                c.id, 
                c.nombre, 
                c.tipo, 
                c.color,
                c.id_icono,
                c.create_at,
                c.update_at,
                cu.id as concepto_usuario_id,
                cu.desembolso_planejado, 
                cu.periodo_tipo, 
                cu.periodo_dia,
                cu.periodo_desembolso,
                cu.limite_monto, 
                cu.limite_tipo,
                cu.limite_fecha,
                cu.has_notificacion,
                cu.visible,
                i.path AS icono_path, 
                i.nombre AS icono_nombre
            FROM concepto c
            INNER JOIN concepto_usuarios cu ON c.id = cu.id_concepto
            LEFT JOIN iconos i ON c.id_icono = i.id
            WHERE cu.id_usuario = ? AND c.tipo = ?
                AND c.delete_at IS NULL 
                AND cu.delete_at IS NULL
                AND cu.visible = 1";

		$result = $this->executeQuery($sql, [$idUsuario, $tipo]);
		// Convertir a Entities
		$conceptos = [];
		foreach ($result as $row) {
			$conceptos[] = [
				'concepto' => new EntityConcepto(
					$row['id'],
					$row['id_icono'],
					$row['nombre'],
					$row['tipo'], // ← Ahora tipo va antes de color
					$row['color'],
					$row['create_at'],
					$row['update_at'],
					null
				),
				'concepto_usuario' => new EntityConceptoUsuario(
					$row['concepto_usuario_id'],
					$idUsuario,
					$row['id'],
					$row['desembolso_planejado'],
					$row['limite_monto'], // ← Ahora limite_monto va antes de los parámetros de periodo
					$row['periodo_desembolso'],
					$row['periodo_tipo'],
					$row['periodo_dia'],
					$row['limite_tipo'],
					$row['limite_fecha'],
					$row['has_notificacion'],
					$row['visible'],
					$row['create_at'],
					$row['update_at'],
					null
				),
				'icono_path' => $row['icono_path'],
				'icono_nombre' => $row['icono_nombre']
			];
		}

		return $conceptos;
	}

	public function getIconos()
	{
		$sql = "SELECT id, path, nombre from iconos";
		$result = $this->executeQuery($sql);
		$icons = [];
		foreach ($result as $row) {
			$icons[] = new EntityIcono(
				$row['id'],
				$row['path'],
				$row['nombre'],
				null,
				null,
				null,
				null,
				null
			);
		}

		return $icons;
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

		$result = $this->executeQuery($sql, $params);

		// Convertir a Entities
		$conceptos = [];
		foreach ($result as $row) {
			$conceptos[] = new EntityConcepto(
				$row['id'],
				$row['id_icono'],
				$row['nombre'],
				$row['color'],
				$row['tipo'],
				$row['create_at'],
				$row['update_at'],
				null
			);
		}

		return $conceptos;
	}

	public function getConceptoPorId($id)
	{
		$sql = "SELECT c.*, i.path as icono_path, i.nombre as icono_nombre
                FROM concepto c
                LEFT JOIN iconos i ON c.id_icono = i.id
                WHERE c.id = ? AND c.delete_at IS NULL";

		$result = $this->executeQuery($sql, [$id]);

		if (!empty($result)) {
			$row = $result[0];
			return new EntityConcepto(
				$row['id'],
				$row['id_icono'],
				$row['nombre'],
				$row['color'],
				$row['tipo'],
				$row['create_at'],
				$row['update_at'],
				null
			);
		}

		return null;
	}

	public function getConceptoPorIdYUsuario($idConcepto, $idUsuario)
	{
		$sql = "SELECT 
                c.id, 
                c.nombre, 
                c.tipo, 
                c.color,
                c.id_icono,
                c.create_at,
                c.update_at,
                cu.id as concepto_usuario_id,
                cu.desembolso_planejado, 
                cu.periodo_tipo, 
                cu.periodo_dia,
                cu.periodo_desembolso,
                cu.limite_monto, 
                cu.limite_tipo,
                cu.limite_fecha,
                cu.has_notificacion,
                cu.visible,
                i.nombre AS icono_nombre
            FROM concepto c
            INNER JOIN concepto_usuarios cu ON c.id = cu.id_concepto
            LEFT JOIN iconos i ON c.id_icono = i.id
            WHERE c.id = ? 
                AND cu.id_usuario = ?
                AND c.delete_at IS NULL 
                AND cu.delete_at IS NULL";

		$result = $this->executeQuery($sql, [$idConcepto, $idUsuario]);

		if (!empty($result)) {
			$row = $result[0];
			return [
				'concepto' => new EntityConcepto(
					$row['id'],
					$row['id_icono'],
					$row['nombre'],
					$row['color'],
					$row['tipo'],
					$row['create_at'],
					$row['update_at'],
					null
				),
				'concepto_usuario' => new EntityConceptoUsuario(
					$row['concepto_usuario_id'],
					$idUsuario,
					$row['id'],
					$row['desembolso_planejado'],
					$row['periodo_desembolso'],
					$row['periodo_tipo'],
					$row['periodo_dia'],
					$row['limite_monto'],
					$row['limite_tipo'],
					$row['limite_fecha'],
					$row['has_notificacion'],
					$row['visible'],
					$row['create_at'],
					$row['update_at'],
					null
				)
			];
		}

		return null;
	}

	// Los demás métodos se mantienen igual...
	public function crearConcepto($nombre, $tipo, $color, $idIcono)
	{
		$sql = "INSERT INTO concepto (nombre, tipo, color, id_icono) 
                VALUES (?, ?, ?, ?)";

		return $this->executeNonQuery($sql, [$nombre, $tipo, $color, $idIcono]);
	}

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
