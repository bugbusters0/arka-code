<?php
// backend/models/ConceptoModel.php
require_once ROOT . '/backend/commons/BaseModel.php';
require_once ROOT . '/backend/entities/EntityConcepto.php';
require_once ROOT . '/backend/entities/EntityConceptoUsuario.php';
require_once ROOT . '/backend/entities/EntityIcono.php';

class ConceptoModel extends BaseModel
{
	// ============================================
	// MÉTODOS DE CONSULTA
	// ============================================

	/**
	 * Obtiene todos los conceptos de un usuario específico
	 */
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
                AND cu.visible = 1
            ORDER BY c.nombre ASC";

		$result = $this->executeQuery($sql, [$idUsuario, $tipo]);

		// Convertir a Entities
		$conceptos = [];
		foreach ($result as $row) {
			$conceptos[] = [
				'concepto' => new EntityConcepto(
					$row['id'],
					$row['id_icono'],
					$row['nombre'],
					$row['tipo'],
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
					$row['limite_monto'],
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

	/**
	 * Obtiene un concepto específico por ID
	 */
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
				$row['tipo'],
				$row['color'],
				$row['create_at'],
				$row['update_at'],
				null
			);
		}

		return null;
	}

	/**
	 * Obtiene un concepto específico para un usuario
	 */
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
                i.path AS icono_path,
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
					$row['tipo'],
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
					$row['limite_monto'],
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

		return null;
	}

	/**
	 * Obtiene todos los conceptos (para administración)
	 */
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

	// ============================================
	// MÉTODOS DE ICONOS
	// ============================================

	/**
	 * Obtiene todos los iconos disponibles
	 */
	public function getIconos()
	{
		$sql = "SELECT id, path, nombre, size, type 
		        FROM iconos 
		        WHERE delete_at IS NULL 
		        ORDER BY nombre ASC";

		$result = $this->executeQuery($sql);
		$icons = [];

		foreach ($result as $row) {
			$icons[] = new EntityIcono(
				$row['id'],
				$row['path'],
				$row['nombre'],
				$row['size'],
				$row['type'],
				null,
				null,
				null
			);
		}

		return $icons;
	}

	/**
	 * Alias de getIconos() para compatibilidad
	 */
	public function getIconosDisponibles()
	{
		return $this->getIconos();
	}

	// ============================================
	// MÉTODOS DE USUARIOS Y FAMILIA
	// ============================================

	/**
	 * Obtiene todos los usuarios de una familia
	 */
	public function getUsuariosPorFamilia($idFamilia)
	{
		$sql = "SELECT id, nombre, rol 
		        FROM usuario 
		        WHERE id_familia = ? AND delete_at IS NULL
		        ORDER BY rol DESC, nombre ASC";

		return $this->executeQuery($sql, [$idFamilia]);
	}

	// ============================================
	// MÉTODOS DE CREACIÓN
	// ============================================

	/**
	 * Crea un nuevo concepto en la tabla concepto
	 */
	public function crearConcepto($nombre, $tipo, $color, $idIcono)
	{
		$sql = "INSERT INTO concepto (id_icono, nombre, tipo, color) 
                VALUES (?, ?, ?, ?)";

		return $this->executeNonQuery($sql, [$idIcono, $nombre, $tipo, $color]);
	}

	/**
	 * Crea una relación usuario-concepto CON configuración
	 * (Para el usuario que crea el concepto)
	 */
	public function crearRelacionUsuarioConcepto($idUsuario, $idConcepto, $desembolsoMonto, $desembolsoFrecuencia, $limiteMonto, $limiteFrecuencia)
	{
		$sql = "INSERT INTO concepto_usuarios 
            (id_usuario, id_concepto, desembolso_planejado, periodo_tipo, limite_monto, limite_tipo, visible) 
            VALUES (?, ?, ?, ?, ?, ?, 1)";

		$params = [
			$idUsuario,
			$idConcepto,
			$desembolsoMonto,
			$desembolsoFrecuencia, // puede ser NULL
			$limiteMonto,
			$limiteFrecuencia       // puede ser NULL
		];

		return $this->executeNonQuery($sql, $params);
	}

	/**
	 * Crea una relación usuario-concepto SIN configuración (NULL)
	 * (Para los demás usuarios de la familia)
	 */
	public function crearRelacionUsuarioConceptoBasica($idUsuario, $idConcepto)
	{
		$sql = "INSERT INTO concepto_usuarios 
            (id_usuario, id_concepto, desembolso_planejado, periodo_tipo, limite_monto, limite_tipo, visible) 
            VALUES (?, ?, NULL, NULL, NULL, NULL, 1)";

		$params = [$idUsuario, $idConcepto];

		return $this->executeNonQuery($sql, $params);
	}

	// ============================================
	// MÉTODOS DE ACTUALIZACIÓN
	// ============================================

	/**
	 * Actualiza un concepto (nombre, color, icono)
	 * AFECTA A TODOS LOS USUARIOS QUE TIENEN ESTE CONCEPTO
	 */
	public function actualizarConcepto($id, $nombre, $color, $idIcono)
	{
		$sql = "UPDATE concepto 
            SET nombre = ?, color = ?, id_icono = ?, update_at = CURRENT_TIMESTAMP 
            WHERE id = ? AND delete_at IS NULL";

		return $this->executeNonQuery($sql, [$nombre, $color, $idIcono, $id]);
	}

	/**
	 * Actualiza la configuración personal de un usuario para un concepto
	 * SOLO AFECTA AL USUARIO ACTUAL
	 */
	public function actualizarRelacionUsuarioConcepto($idUsuario, $idConcepto, $desembolsoMonto, $desembolsoFrecuencia, $limiteMonto, $limiteFrecuencia)
	{
		$sql = "UPDATE concepto_usuarios 
            SET desembolso_planejado = ?, 
                periodo_tipo = ?, 
                limite_monto = ?, 
                limite_tipo = ?, 
                update_at = CURRENT_TIMESTAMP 
            WHERE id_usuario = ? AND id_concepto = ? AND delete_at IS NULL";

		$params = [
			$desembolsoMonto,
			$desembolsoFrecuencia,
			$limiteMonto,
			$limiteFrecuencia,
			$idUsuario,
			$idConcepto
		];

		return $this->executeNonQuery($sql, $params);
	}

	// ============================================
	// MÉTODOS DE DESHABILITACIÓN
	// ============================================

	/**
	 * Deshabilita un concepto para un usuario específico
	 * (marca visible = 0 en concepto_usuarios)
	 */
	public function deshabilitarConceptoUsuario($idUsuario, $idConcepto)
	{
		$sql = "UPDATE concepto_usuarios 
            SET visible = 0, update_at = CURRENT_TIMESTAMP 
            WHERE id_usuario = ? AND id_concepto = ? AND delete_at IS NULL";

		return $this->executeNonQuery($sql, [$idUsuario, $idConcepto]);
	}

	/**
	 * Reactiva un concepto para un usuario específico
	 * (marca visible = 1 en concepto_usuarios)
	 */
	public function reactivarConceptoUsuario($idUsuario, $idConcepto)
	{
		$sql = "UPDATE concepto_usuarios 
            SET visible = 1, update_at = CURRENT_TIMESTAMP 
            WHERE id_usuario = ? AND id_concepto = ? AND delete_at IS NULL";

		return $this->executeNonQuery($sql, [$idUsuario, $idConcepto]);
	}

	// ============================================
	// MÉTODOS DE ELIMINACIÓN
	// ============================================

	/**
	 * Elimina un concepto (soft delete)
	 * AFECTA A TODOS LOS USUARIOS
	 */
	public function eliminarConcepto($id)
	{
		$sql = "UPDATE concepto 
                SET delete_at = CURRENT_TIMESTAMP 
                WHERE id = ?";

		return $this->executeNonQuery($sql, [$id]);
	}

	/**
	 * Elimina la relación usuario-concepto (soft delete)
	 * SOLO AFECTA AL USUARIO ESPECÍFICO
	 */
	public function eliminarRelacionUsuarioConcepto($idUsuario, $idConcepto)
	{
		$sql = "UPDATE concepto_usuarios 
                SET delete_at = CURRENT_TIMESTAMP 
                WHERE id_usuario = ? AND id_concepto = ?";

		return $this->executeNonQuery($sql, [$idUsuario, $idConcepto]);
	}

	// ============================================
	// MÉTODOS DE UTILIDAD
	// ============================================

	/**
	 * Obtiene el último ID insertado
	 */
	public function getLastInsertId()
	{
		return $this->db->lastInsertId();
	}

	/**
	 * Verifica si un concepto existe
	 */
	public function existeConcepto($id)
	{
		$sql = "SELECT COUNT(*) as count FROM concepto WHERE id = ? AND delete_at IS NULL";
		$result = $this->executeQuery($sql, [$id]);
		return !empty($result) && $result[0]['count'] > 0;
	}

	/**
	 * Verifica si un usuario tiene acceso a un concepto
	 */
	public function usuarioTieneAccesoConcepto($idUsuario, $idConcepto)
	{
		$sql = "SELECT COUNT(*) as count 
		        FROM concepto_usuarios 
		        WHERE id_usuario = ? AND id_concepto = ? AND delete_at IS NULL";

		$result = $this->executeQuery($sql, [$idUsuario, $idConcepto]);
		return !empty($result) && $result[0]['count'] > 0;
	}

	/**
	 * Cuenta cuántos conceptos tiene un usuario
	 */
	public function contarConceptosUsuario($idUsuario, $tipo = null)
	{
		$sql = "SELECT COUNT(*) as count 
		        FROM concepto c
		        INNER JOIN concepto_usuarios cu ON c.id = cu.id_concepto
		        WHERE cu.id_usuario = ? 
		        AND c.delete_at IS NULL 
		        AND cu.delete_at IS NULL
		        AND cu.visible = 1";

		$params = [$idUsuario];

		if ($tipo) {
			$sql .= " AND c.tipo = ?";
			$params[] = $tipo;
		}

		$result = $this->executeQuery($sql, $params);
		return !empty($result) ? $result[0]['count'] : 0;
	}
}
