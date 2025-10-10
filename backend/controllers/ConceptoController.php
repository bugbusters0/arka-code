<?php
// backend/controllers/ConceptoController.php
require_once ROOT . '/backend/commons/BaseController.php';
require_once ROOT . '/backend/models/ConceptoModel.php';

class ConceptoController extends BaseController
{
	private $conceptoModel;

	public function __construct()
	{
		parent::__construct();
		$this->requireAuth();
		$this->conceptoModel = new ConceptoModel();
	}

	public function mostrarConceptos($tipo = "gasto")
	{
		if (!isset($_SESSION['miembro_id']) || !isset($_SESSION['user_id'])) {
			header('Location: ' . URLROOT . '/seleccionar-perfil');
			exit;
		}

		$userId = $_SESSION['miembro_id'];
		$conceptos = $this->conceptoModel->getConceptosPorUsuario($userId, $tipo);
		$icons = $this->conceptoModel->getIconos();

		$data = [
			'title' => ($tipo === 'gasto') ? 'Gestión de Gastos - Arka' : 'Gestión de Ingresos - Arka',
			'conceptos' => $conceptos,
			'icons' => $icons,
			'tipo_actual' => $tipo
		];

		$this->viewWithLayout('concepto/gasto', 'main', $data);
	}

	public function guardarConcepto()
	{
		// if ($_SERVER['REQUEST_METHOD'] !== 'POST') {
		// 	header('Location: ' . URLROOT . '/concepto/gasto');
		// 	exit;
		// }

		// if (!isset($_SESSION['miembro_id']) || !isset($_SESSION['user_id'])) {
		// 	$_SESSION['error_message'] = 'Sesión inválida';
		// 	header('Location: ' . URLROOT . '/seleccionar-perfil');
		// 	exit;
		// }

		$idUsuarioCreador = $_SESSION['miembro_id'];
		$idFamilia = $_SESSION['user_id'];
		$nombre = trim($_POST['nombre'] ?? '');
		$tipo = $_POST['tipo'] ?? 'gasto';
		$color = $_POST['color'] ?? '#FF6B6B';
		$idIcono = intval($_POST['id_icono'] ?? 3);

		// Validaciones básicas
		if (empty($nombre)) {
			$_SESSION['error_message'] = 'El nombre del concepto es obligatorio';
			header('Location: ' . URLROOT . '/concepto/gasto');
			exit;
		}

		try {
			// 1. Crear concepto en tabla 'concepto'
			$conceptoCreado = $this->conceptoModel->crearConcepto($nombre, $tipo, $color, $idIcono);

			if (!$conceptoCreado) {
				throw new Exception('Error al crear el concepto');
			}

			$idConcepto = $this->conceptoModel->getLastInsertId();

			// 2. Obtener configuración del usuario creador
			$desembolsoMonto = floatval($_POST['desembolsoMonto'] ?? 0);
			$desembolsoFrecuencia = $_POST['desembolsoFrecuencia'] ?? null;
			$limiteMonto = floatval($_POST['limiteMonto'] ?? 0);
			$limiteFrecuencia = $_POST['limiteFrecuencia'] ?? null;

			// 3. Obtener TODOS los usuarios de la familia
			$usuariosFamilia = $this->conceptoModel->getUsuariosPorFamilia($idFamilia);

			if (empty($usuariosFamilia)) {
				throw new Exception('No se encontraron usuarios en la familia');
			}

			$relacionesCreadas = 0;

			// 4. Crear relaciones para TODOS los usuarios de la familia
			foreach ($usuariosFamilia as $usuario) {
				$idUsuarioActual = $usuario['id'];

				if ($idUsuarioActual == $idUsuarioCreador) {
					$result = $this->conceptoModel->crearRelacionUsuarioConcepto(
						$idUsuarioActual,
						$idConcepto,
						$desembolsoMonto,
						$desembolsoFrecuencia,
						$limiteMonto,
						$limiteFrecuencia
					);
				} else {
					$result = $this->conceptoModel->crearRelacionUsuarioConceptoBasica(
						$idUsuarioActual,
						$idConcepto
					);
				}

				if ($result) {
					$relacionesCreadas++;
				}
			}

			if ($relacionesCreadas > 0) {
				$_SESSION['success_message'] = "Concepto '{$nombre}' creado y compartido con {$relacionesCreadas} usuarios de la familia";
				header('Location: ' . URLROOT . '/concepto/gasto?tipo=' . $tipo);
				exit;
			} else {
				$this->conceptoModel->eliminarConcepto($idConcepto);
				throw new Exception('Error al crear las relaciones con los usuarios');
			}
		} catch (Exception $e) {
			error_log("❌ Error en guardarConcepto: " . $e->getMessage());
			$_SESSION['error_message'] = 'Error al guardar el concepto: ' . $e->getMessage();
			header('Location: ' . URLROOT . '/concepto/gasto');
			exit;
		}
	}

	public function editar($idConcepto)
	{
		if (!isset($_SESSION['miembro_id'])) {
			$_SESSION['error_message'] = 'Sesión inválida';
			header('Location: ' . URLROOT . '/seleccionar-perfil');
			exit;
		}

		$idUsuario = $_SESSION['miembro_id'];

		if ($_SERVER['REQUEST_METHOD'] === 'POST') {
			$nombre = trim($_POST['name'] ?? '');
			$color = $_POST['color'] ?? '#FF6B6B';
			$idIcono = intval($_POST['id_icono'] ?? 3);

			if (empty($nombre)) {
				$_SESSION['error_message'] = 'El nombre es requerido';
				header('Location: ' . URLROOT . '/concepto/gasto');
				exit;
			}

			try {
				// Verificar acceso
				$conceptoUsuario = $this->conceptoModel->getConceptoPorIdYUsuario($idConcepto, $idUsuario);
				if (!$conceptoUsuario) {
					$_SESSION['error_message'] = 'No tienes permisos para editar este concepto';
					header('Location: ' . URLROOT . '/concepto/gasto');
					exit;
				}

				// Actualizar concepto (afecta a todos)
				$conceptoResult = $this->conceptoModel->actualizarConcepto($idConcepto, $nombre, $color, $idIcono);

				if (!$conceptoResult) {
					throw new Exception('Error al actualizar el concepto');
				}

				// Actualizar configuración personal del usuario
				$desembolsoMonto = floatval($_POST['desembolsoMonto'] ?? 0);
				$desembolsoFrecuencia = $_POST['desembolsoFrecuencia'] ?? null;
				$limiteMonto = floatval($_POST['limiteMonto'] ?? 0);
				$limiteFrecuencia = $_POST['limiteFrecuencia'] ?? null;

				$updateResult = $this->conceptoModel->actualizarRelacionUsuarioConcepto(
					$idUsuario,
					$idConcepto,
					$desembolsoMonto,
					$desembolsoFrecuencia,
					$limiteMonto,
					$limiteFrecuencia
				);

				if ($updateResult) {
					$_SESSION['success_message'] = "Concepto '{$nombre}' actualizado correctamente";
				} else {
					$_SESSION['error_message'] = 'Error al actualizar tu configuración personal';
				}

				header('Location: ' . URLROOT . '/concepto/gasto');
				exit;
			} catch (Exception $e) {
				error_log("❌ Error en editar: " . $e->getMessage());
				$_SESSION['error_message'] = 'Error al actualizar el concepto';
				header('Location: ' . URLROOT . '/concepto/gasto');
				exit;
			}
		} else {
			// GET: Mostrar formulario de edición
			$concepto = $this->conceptoModel->getConceptoPorIdYUsuario($idConcepto, $idUsuario);

			if (!$concepto) {
				$_SESSION['error_message'] = 'Concepto no encontrado';
				header('Location: ' . URLROOT . '/concepto/gasto');
				exit;
			}

			$icons = $this->conceptoModel->getIconos();

			$data = [
				'title' => 'Editar Concepto - Arka',
				'concepto' => $concepto,
				'icons' => $icons
			];

			$this->viewWithLayout('concepto/editar', 'main', $data);
		}
	}

	public function deshabilitar($idConcepto)
	{
		if ($_SERVER['REQUEST_METHOD'] !== 'POST') {
			header('Location: ' . URLROOT . '/concepto/gasto');
			exit;
		}

		if (!isset($_SESSION['miembro_id'])) {
			$_SESSION['error_message'] = 'Sesión inválida';
			header('Location: ' . URLROOT . '/seleccionar-perfil');
			exit;
		}

		$idUsuario = $_SESSION['miembro_id'];

		try {
			// Verificar acceso
			$conceptoUsuario = $this->conceptoModel->getConceptoPorIdYUsuario($idConcepto, $idUsuario);
			if (!$conceptoUsuario) {
				$_SESSION['error_message'] = 'No tienes permisos para deshabilitar este concepto';
				header('Location: ' . URLROOT . '/concepto/gasto');
				exit;
			}

			// Deshabilitar solo para este usuario
			$result = $this->conceptoModel->deshabilitarConceptoUsuario($idUsuario, $idConcepto);

			if ($result) {
				$_SESSION['success_message'] = 'Concepto deshabilitado correctamente';
			} else {
				$_SESSION['error_message'] = 'Error al deshabilitar el concepto';
			}
		} catch (Exception $e) {
			error_log("❌ Error en deshabilitar: " . $e->getMessage());
			$_SESSION['error_message'] = 'Error al deshabilitar el concepto';
		}

		header('Location: ' . URLROOT . '/concepto/gasto');
		exit;
	}
}
