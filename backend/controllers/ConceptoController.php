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

    $data = [
      'title' => 'Gestión de Conceptos - Arka',
      'conceptos' => $conceptos
    ];

    // ✅ Automáticamente cargará:
    // - concepto-listar.css
    // - concepto-listar.js
    $this->viewWithLayout('concepto/listar', 'main', $data);
  }
}
