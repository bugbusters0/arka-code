<?php
// backend/controllers/AuthController.php
require_once ROOT . '/backend/commons/BaseController.php';
require_once ROOT . '/backend/models/UserModel.php';
require_once ROOT . '/backend/models/FamiliaModel.php';

class AuthController extends BaseController {
    private $userModel;
    private $familiaModel;
    
    public function __construct() {
        parent::__construct();
        $this->userModel = new UserModel();
        $this->familiaModel = new FamiliaModel();
    }
    
    public function showLogin() {
        $data['title'] = 'Login - Arka App';
        $data['validation_errors'] = $_SESSION['validation_errors'] ?? [];
        $data['success'] = $_SESSION['success_message'] ?? '';
        unset($_SESSION['validation_errors']);
        $data['error'] = '';
        $this->view('auth/login', $data);
    }
    
    public function showRegister() {
        $data['title'] = 'Registro - Arka App';
        $data['validation_errors'] = $_SESSION['validation_errors'] ?? [];
        $data['success'] = $_SESSION['success_message'] ?? '';
        unset($_SESSION['validation_errors'], $_SESSION['success_message']);
        $data['error'] = '';
        $this->view('auth/registrar', $data);
    }
    
    public function processLogin() {
        $email = $_POST['email'] ?? '';
        $contrasena = $_POST['contrasena'] ?? '';
        
        $userEntity = $this->userModel->verificarUsuario($email, $contrasena);
        if ($userEntity) {
            $_SESSION['user_id'] = $userEntity->getId();
            $_SESSION['user_email'] = $userEntity->getEmail();
            header('Location: ' . URLROOT . '/dashboard');
            exit;
        } else {
            $data['error'] = 'Credenciales inválidas';
            $data['title'] = 'Login - Arka App';
            $this->view('auth/login', $data);
        }
    }
    
    public function processRegister() {
        // Los datos ya vienen validados y limpios desde index.php
        try {
            // 1. Crear familia
            $familiaId = $this->familiaModel->crearFamilia(
                $_POST['email'],
                $_POST['telefono'],
                $_POST['password']
            );
            
            if (!$familiaId) {
                throw new Exception('Error al crear la familia');
            }
            
            // 2. Crear usuario administrador
            $adminId = $this->userModel->crearUsuarioAdmin(
                $familiaId,
                $_POST['adminNombre'],
                $_POST['adminNacimiento'],
                $_POST['contraPersonal']
            );
            
            if (!$adminId) {
                throw new Exception('Error al crear el administrador');
            }
            
            // 3. Crear miembros adicionales (si existen)
            if (!empty($_POST['miembros'])) {
                foreach ($_POST['miembros'] as $miembro) {
                    $this->userModel->crearUsuarioMiembro(
                        $familiaId,
                        $miembro['nombre'],
                        $miembro['nacimiento'],
                        $miembro['contra'],
                        $miembro['rol']
                    );
                }
            }
            
            // 4. Éxito: redirigir a login
            $_SESSION['success_message'] = 'Cuenta creada exitosamente. ¡Inicia sesión!';
            header('Location: ' . URLROOT . '/login');
            exit;
            
        } catch (Exception $e) {
            // Error: mostrar en formulario
            error_log("Error en registro: " . $e->getMessage());
            $data['error'] = 'Ocurrió un error al crear la cuenta. Intenta nuevamente.';
            $data['title'] = 'Registro - Arka App';
            $this->view('auth/registrar', $data);
        }
    }
    
    public function logout() {
        session_destroy();
        header('Location: ' . URLROOT . '/login');
        exit;
    }
}
?>