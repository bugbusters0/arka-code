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
    $this->requireAuth(); // Protege todas las rutas
    $this->conceptoModel = new ConceptoModel();
  }

  public function listar()
  {
    // Obtener conceptos con filtros opcionales
    $search = $_GET['search'] ?? '';
    $tipo = $_GET['tipo'] ?? '';

    $conceptos = $this->conceptoModel->getConceptos($search, $tipo);

    // Preparar datos para la vista
    $data = [
      'title' => 'Gestión de Conceptos - Arka',
      'conceptos' => $conceptos,
      'pageCSS' => 'conceptos', // Carga CSS específico si existe
      'pageJS' => 'conceptos'   // Carga JS específico si existe
    ];

    // ✅ AQUÍ USAMOS EL LAYOUT
    // Captura el contenido de la vista
    ob_start();
    require_once ROOT . '/frontend/views/concepto/listar.php';
    $content = ob_get_clean();

    // Renderiza con layout
    $this->layout('main', $content, $data);
  }

  public function crear()
  {
    if ($_SERVER['REQUEST_METHOD'] === 'POST') {
      // Validar y crear concepto
      $nombre = $_POST['nombre'] ?? '';
      $tipo = $_POST['tipo'] ?? '';
      $color = $_POST['color'] ?? '#667eea';
      $idIcono = $_POST['id_icono'] ?? 1;

      if ($this->conceptoModel->crearConcepto($nombre, $tipo, $color, $idIcono)) {
        $_SESSION['success_message'] = 'Concepto creado exitosamente';
        header('Location: ' . URLROOT . '/concepto/listar');
        exit;
      } else {
        $_SESSION['error_message'] = 'Error al crear el concepto';
      }
    }

    $data = [
      'title' => 'Crear Concepto - Arka',
      'iconos' => $this->conceptoModel->getIconosDisponibles()
    ];

    ob_start();
    require_once ROOT . '/frontend/views/concepto/crear.php';
    $content = ob_get_clean();

    $this->layout('main', $content, $data);
  }

  public function editar($id)
  {
    $concepto = $this->conceptoModel->getConceptoPorId($id);

    if (!$concepto) {
      $_SESSION['error_message'] = 'Concepto no encontrado';
      header('Location: ' . URLROOT . '/concepto/listar');
      exit;
    }

    if ($_SERVER['REQUEST_METHOD'] === 'POST') {
      // Actualizar concepto
      $nombre = $_POST['nombre'] ?? '';
      $color = $_POST['color'] ?? '';

      if ($this->conceptoModel->actualizarConcepto($id, $nombre, $color)) {
        $_SESSION['success_message'] = 'Concepto actualizado exitosamente';
        header('Location: ' . URLROOT . '/concepto/listar');
        exit;
      }
    }

    $data = [
      'title' => 'Editar Concepto - Arka',
      'concepto' => $concepto
    ];

    ob_start();
    require_once ROOT . '/frontend/views/concepto/editar.php';
    $content = ob_get_clean();

    $this->layout('main', $content, $data);
  }

  public function eliminar($id)
  {
    if ($this->conceptoModel->eliminarConcepto($id)) {
      $_SESSION['success_message'] = 'Concepto eliminado exitosamente';
    } else {
      $_SESSION['error_message'] = 'Error al eliminar el concepto';
    }

    header('Location: ' . URLROOT . '/concepto/listar');
    exit;
  }
}
