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
        $idUsuario = $_SESSION['user_id'] ?? 1; // Ejemplo, reemplaza con tu lógica de sesión
        $search = $_GET['search'] ?? '';
        $tipo = $_GET['tipo'] ?? '';

        $conceptos = $this->conceptoModel->getConceptosPorUsuario($idUsuario, $search, $tipo);

        $data = [
            'title' => 'Gestión de Conceptos - Arka',
            'conceptos' => $conceptos
        ];

        if (isset($_GET['format']) && $_GET['format'] === 'json') {
            header('Content-Type: application/json');
            echo json_encode($data);
            exit;
        }

        $this->viewWithLayout('concepto/listar', 'main', $data);
    }

public function guardarConcepto()
    {
        if ($_SERVER['REQUEST_METHOD'] === 'POST') {
            header('Content-Type: application/json');
            
            $idUsuario = $_SESSION['user_id'] ?? 1;
            $nombre = $_POST['name'] ?? '';
            $tipo = $_POST['tipo'] ?? 'gasto';
            $color = $_POST['color'] ?? '#FF6B6B';
            $idIcono = intval($_POST['idIcono'] ?? 3);

            // Validaciones básicas
            if (empty($nombre)) {
                echo json_encode(['success' => false, 'message' => 'El nombre es requerido']);
                exit;
            }

            try {
                // 1. Crear concepto en tabla 'concepto'
                $conceptoResult = $this->conceptoModel->crearConcepto($nombre, $tipo, $color, $idIcono);
                
                if ($conceptoResult) {
                    $idConcepto = $this->conceptoModel->getLastInsertId();
                    error_log("✅ Concepto creado con ID: " . $idConcepto);
                    
                    $desembolsoMonto = floatval($_POST['desembolsoMonto'] ?? 0);
                    $desembolsoFrecuencia = $_POST['desembolsoFrecuencia'] ?? 'diario';
                    $limiteMonto = floatval($_POST['limiteMonto'] ?? 0);
                    $limiteFrecuencia = $_POST['limiteFrecuencia'] ?? 'diario';

                    // 2. Crear relación en 'concepto_usuarios' usando el método del modelo
                    $relacionResult = $this->conceptoModel->crearRelacionUsuarioConcepto(
                        $idUsuario, 
                        $idConcepto, 
                        $desembolsoMonto, 
                        $desembolsoFrecuencia, 
                        $limiteMonto, 
                        $limiteFrecuencia
                    );
                    
                    error_log("✅ Relación creada: " . ($relacionResult ? "SÍ" : "NO"));
                    
                    if ($relacionResult) {
                        echo json_encode([
                            'success' => true, 
                            'message' => 'Concepto guardado correctamente',
                            'concepto_id' => $idConcepto
                        ]);
                    } else {
                        // Si falla la relación, eliminar el concepto creado
                        $this->conceptoModel->eliminarConcepto($idConcepto);
                        echo json_encode(['success' => false, 'message' => 'Error al crear la relación con el usuario']);
                    }
                } else {
                    echo json_encode(['success' => false, 'message' => 'Error al guardar el concepto']);
                }
            } catch (Exception $e) {
                error_log("❌ Error en guardarConcepto: " . $e->getMessage());
                echo json_encode(['success' => false, 'message' => 'Error interno del servidor: ' . $e->getMessage()]);
            }
        }
        exit;
    }

    // En ConceptoController.php - Agregar método editar
public function editar($idConcepto)
{
    header('Content-Type: application/json');
    
    if ($_SERVER['REQUEST_METHOD'] === 'POST') {
        $idUsuario = $_SESSION['user_id'] ?? 1;
        $nombre = $_POST['name'] ?? '';
        $color = $_POST['color'] ?? '#FF6B6B';
        $idIcono = intval($_POST['idIcono'] ?? 3);

        // Validaciones
        if (empty($nombre)) {
            echo json_encode(['success' => false, 'message' => 'El nombre es requerido']);
            exit;
        }

        try {
            // 1. Actualizar concepto en tabla 'concepto'
            $conceptoResult = $this->conceptoModel->actualizarConcepto($idConcepto, $nombre, $color, $idIcono);
            
            if ($conceptoResult) {
                $desembolsoMonto = floatval($_POST['desembolsoMonto'] ?? 0);
                $desembolsoFrecuencia = $_POST['desembolsoFrecuencia'] ?? 'diario';
                $limiteMonto = floatval($_POST['limiteMonto'] ?? 0);
                $limiteFrecuencia = $_POST['limiteFrecuencia'] ?? 'diario';

                // 2. Actualizar relación en 'concepto_usuarios'
                $updateResult = $this->conceptoModel->actualizarRelacionUsuarioConcepto(
                    $idUsuario, 
                    $idConcepto, 
                    $desembolsoMonto, 
                    $desembolsoFrecuencia, 
                    $limiteMonto, 
                    $limiteFrecuencia
                );
                
                if ($updateResult) {
                    echo json_encode([
                        'success' => true, 
                        'message' => 'Concepto actualizado correctamente',
                        'concepto_id' => $idConcepto
                    ]);
                } else {
                    echo json_encode(['success' => false, 'message' => 'Error al actualizar la configuración del concepto']);
                }
            } else {
                echo json_encode(['success' => false, 'message' => 'Error al actualizar el concepto']);
            }
        } catch (Exception $e) {
            error_log("❌ Error en editar concepto: " . $e->getMessage());
            echo json_encode(['success' => false, 'message' => 'Error interno del servidor']);
        }
    } else {
        // GET request - Obtener datos del concepto para editar
        $concepto = $this->conceptoModel->getConceptoPorIdYUsuario($idConcepto, $_SESSION['user_id'] ?? 1);
        
        if ($concepto) {
            echo json_encode([
                'success' => true,
                'concepto' => $concepto
            ]);
        } else {
            echo json_encode([
                'success' => false,
                'message' => 'Concepto no encontrado'
            ]);
        }
    }
    exit;
}
}