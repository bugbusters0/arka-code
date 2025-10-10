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
        $idFamilia = $_SESSION['family_id'] ?? 1; 
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

                // 2. Obtener todos los usuarios de la familia
                $usuariosFamilia = $this->conceptoModel->getUsuariosPorFamilia($idFamilia);
                error_log("👨‍👩‍👧‍👦 Usuarios en familia: " . count($usuariosFamilia));
                
                $relacionesCreadas = 0;
                $totalUsuarios = count($usuariosFamilia);
                
                // 3. Crear relaciones para todos los usuarios de la familia
                foreach ($usuariosFamilia as $usuario) {
                    $idUsuarioActual = $usuario['id'];
                    
                    if ($idUsuarioActual == $idUsuario) {
                        // Para el usuario que crea el concepto - con valores configurados
                        $result = $this->conceptoModel->crearRelacionUsuarioConcepto(
                            $idUsuarioActual, 
                            $idConcepto, 
                            $desembolsoMonto, 
                            $desembolsoFrecuencia, 
                            $limiteMonto, 
                            $limiteFrecuencia
                        );
                    } else {
                        // Para los demás usuarios - con valores NULL
                        $result = $this->conceptoModel->crearRelacionUsuarioConceptoBasica(
                            $idUsuarioActual, 
                            $idConcepto
                        );
                    }
                    
                    if ($result) {
                        $relacionesCreadas++;
                        error_log("✅ Relación creada para usuario: " . $idUsuarioActual);
                    } else {
                        error_log("❌ Error creando relación para usuario: " . $idUsuarioActual);
                    }
                }
                
                if ($relacionesCreadas > 0) {
                    echo json_encode([
                        'success' => true, 
                        'message' => "Concepto guardado correctamente y compartido con " . $relacionesCreadas . " usuarios",
                        'concepto_id' => $idConcepto,
                        'usuarios_compartidos' => $relacionesCreadas
                    ]);
                } else {
                    // Si fallan todas las relaciones, eliminar el concepto creado
                    $this->conceptoModel->eliminarConcepto($idConcepto);
                    echo json_encode(['success' => false, 'message' => 'Error al crear las relaciones con los usuarios']);
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
            // 1. Verificar que el usuario tiene acceso a este concepto
            $conceptoUsuario = $this->conceptoModel->getConceptoPorIdYUsuario($idConcepto, $idUsuario);
            if (!$conceptoUsuario) {
                echo json_encode(['success' => false, 'message' => 'No tienes permisos para editar este concepto']);
                exit;
            }

            // 2. Actualizar concepto en tabla 'concepto' (esto afecta a todos los usuarios)
            $conceptoResult = $this->conceptoModel->actualizarConcepto($idConcepto, $nombre, $color, $idIcono);
            
            if ($conceptoResult) {
                $desembolsoMonto = floatval($_POST['desembolsoMonto'] ?? 0);
                $desembolsoFrecuencia = $_POST['desembolsoFrecuencia'] ?? 'diario';
                $limiteMonto = floatval($_POST['limiteMonto'] ?? 0);
                $limiteFrecuencia = $_POST['limiteFrecuencia'] ?? 'diario';

                // 3. Actualizar SOLO la relación del usuario actual en 'concepto_usuarios'
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
                        'message' => 'Concepto actualizado correctamente (solo para tu usuario)',
                        'concepto_id' => $idConcepto
                    ]);
                } else {
                    echo json_encode(['success' => false, 'message' => 'Error al actualizar tu configuración personal del concepto']);
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
        $idUsuario = $_SESSION['user_id'] ?? 1;
        $concepto = $this->conceptoModel->getConceptoPorIdYUsuario($idConcepto, $idUsuario);
        
        if ($concepto) {
            echo json_encode([
                'success' => true,
                'concepto' => $concepto
            ]);
        } else {
            echo json_encode([
                'success' => false,
                'message' => 'Concepto no encontrado o no tienes permisos'
            ]);
        }
    }
    exit;
}
}