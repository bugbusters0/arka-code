# Arka Code - Sistema de Gestión Financiera

## 📋 Descripción

Arka Code es una aplicación web desarrollada en PHP puro para la gestión de conceptos financieros, gastos e ingresos. La aplicación incluye autenticación de usuarios, gestión de perfiles familiares y un sistema completo de categorización de movimientos financieros.

## 🚀 Características

- **Autenticación de usuarios** con múltiples perfiles familiares
- **Gestión de conceptos** para gastos e ingresos
- **Interfaz responsive** con diseño moderno
- **Base de datos MySQL** para persistencia de datos
- **Arquitectura MVC** personalizada
- **Sistema de validación** de formularios
- **Middleware de autenticación**

## 🛠️ Tecnologías

- **Backend**: PHP 8.2, MySQL
- **Frontend**: HTML5, CSS3, JavaScript vanilla
- **Servidor**: Apache HTTP Server
- **Contenedores**: Docker
- **Deploy**: Fly.io

---

## 📁 Estructura del Proyecto

```
arka-code/
├── backend/
│   ├── config/
│   │   ├── app.php          # Configuración de la aplicación
│   │   └── db.php           # Configuración de base de datos
│   ├── controllers/         # Controladores MVC
│   │   ├── AuthController.php
│   │   ├── ConceptoController.php
│   │   └── PerfilController.php
│   ├── models/              # Modelos de datos
│   ├── validator/           # Validadores de formularios
│   └── commons/             # Utilidades comunes
├── frontend/
│   └── views/               # Vistas y layouts
│   └── assets/
│       ├── css/             # Estilos de la aplicación
│       ├── js/              # Scripts JavaScript
│       └── img/             # Imágenes y recursos
├── public/
│   ├── index.php           # Punto de entrada de la aplicación
│   └── .htaccess           # Configuración de URLs amigables
├── apache-config/
│   └── 000-default.conf    # Configuración de Apache
├── sql.sql                 # Estructura inicial de la base de datos
├── Dockerfile              # Configuración para Docker
├── docker-compose.yml      # Orquestación de contenedores
└── start.sh               # Script de inicio para producción
```

---

## 🏃‍♂️ Desarrollo Local

### Prerrequisitos

- PHP 8.2 o superior
- Apache HTTP Server
- MySQL 8.0 o MariaDB
- Git

### Instalación Local

1. **Clonar el repositorio**

```bash
git clone <tu-repositorio>
cd arka-code
```

2. **Configurar base de datos**

```bash
# Conectar a MySQL
mysql -u root -p

# Crear base de datos y usuario
CREATE DATABASE arka;
CREATE USER 'arka_user'@'localhost' IDENTIFIED BY 'arka_password';
GRANT ALL PRIVILEGES ON arka_code.* TO 'arka_user'@'localhost';
FLUSH PRIVILEGES;
EXIT;

# Importar estructura inicial
mysql -u root -p arka_code < sql.sql
```

3. **Configurar Apache**

```apache
# En /etc/httpd/conf/httpd.conf o virtual host
DocumentRoot "/ruta/a/arka-code/public"

<Directory "/ruta/a/arka-code/public">
    AllowOverride All
    Require all granted
</Directory>
```

4. **Configurar aplicación**
   Editar la ruta de la aplicacion (nombre de la carpeta) `backend/config/app.php`:

```php
define('URLROOT', 'http://localhost/arka-code');
```

Editar `backend/config/db.php`:

```php
define('DB_HOST', 'localhost');
define('DB_NAME', 'arka_code');
define('DB_USER', 'arka_user');
define('DB_PASS', 'arka_password');
```

5. **Iniciar servidores**

```bash
# Reiniciar Apache
sudo systemctl restart httpd

# O usar servidor PHP built-in
php -S localhost:8000 -t public
```

6. **Acceder a la aplicación**

```
http://localhost
```

## ☁️ Despliegue en Producción (Fly.io)

### Prerrequisitos

- Cuenta en [Fly.io](https://fly.io)
- Fly CLI instalado

### Comandos de despliegue

```bash
# Login en Fly.io
fly auth login

# Inicializar aplicación (primera vez)
fly launch

# Desplegar aplicación
fly deploy

# Ver logs
fly logs

# Conectar a la base de datos
fly ssh console
mysql -u arka_user -p arka_code
```

### URLs de producción

- **Aplicación**: https://arka-code.fly.dev
- **Dashboard**: https://fly.io/apps/arka-code

---

## 🔧 Configuración

### Variables de Entorno

La aplicación detecta automáticamente el entorno:

**Desarrollo Local:**

```php
URLROOT = 'http://localhost'
DB_HOST = 'localhost'
```

**Producción (Fly.io):**

```php
URLROOT = 'https://arka-code.fly.dev'
DB_HOST = 'localhost' (con socket)
```

### Archivos de Configuración

#### `backend/config/app.php`

```php
// Configuración general de la aplicación
define('URLROOT', 'http://localhost');
define('APP_NAME', 'Arka App');
date_default_timezone_set('America/Lima');
```

#### `backend/config/db.php`

```php
// Configuración de conexión a base de datos
define('DB_HOST', 'localhost');
define('DB_NAME', 'arka_code');
define('DB_USER', 'arka_user');
define('DB_PASS', 'arka_password');
```

---

## 🏗️ Arquitectura del Sistema

### Estructura del Proyecto

```
arka-code/
├── backend/
│   ├── config/           # Configuración de app y BD
│   ├── controllers/      # Lógica de negocio (MVC)
│   ├── models/          # Acceso a datos
│   ├── entities/        # Objetos de dominio
│   ├── validators/      # Validación de formularios
│   └── commons/         # Utilidades compartidas
├── frontend/
│   ├── assets/          # CSS, JS, imágenes
│   └── views/           # Vistas y layouts
├── public/              # Punto de entrada público
└── sql.sql             # Estructura de base de datos
```

## 🔄 Flujo de la Aplicación

### 1. Punto de Entrada - `public/index.php`

```php
<?php
require_once '../backend/config/init.php';

// Procesamiento de rutas
$uri = parse_url($_SERVER['REQUEST_URI'], PHP_URL_PATH);
$uri = str_replace('/arka-code', '', $uri);
$uri = trim($uri, '/');

// Sistema de enrutamiento manual
switch ($uri) {
    case 'login':
        // 1. Verifica si el usuario NO está autenticado
        AuthMiddleware::guest();

        // 2. Si es POST, valida y procesa login
        if ($_SERVER['REQUEST_METHOD'] === 'POST') {
            $validation = LoginValidator::validate($_POST);
            if (!$validation['success']) {
                // Guarda errores en sesión y redirige
                $_SESSION['validation_errors'] = $validation['errors'];
                header('Location: ' . URLROOT . '/login');
                exit;
            }

            // Combina datos limpios
            $_POST = array_merge($_POST, $validation['cleanData']);
        }

        // 3. Crea controlador y ejecuta acción
        require_once ROOT . '/backend/controllers/AuthController.php';
        $controller = new AuthController();
        if ($_SERVER['REQUEST_METHOD'] === 'POST') {
            $controller->enviarCredenciales(); // Procesa login
        } else {
            $controller->showLogin(); // Muestra formulario
        }
        break;

    // ... más rutas
}
?>
```

**Explicación del flujo:**

1. **Inicialización**: Carga configuración y sesiones
2. **Parseo de URI**: Limpia y normaliza la ruta solicitada
3. **Middleware**: Verifica permisos (guest, auth, familyOnly)
4. **Validación**: Para POST, valida datos del formulario
5. **Controlador**: Ejecuta la lógica correspondiente
6. **Vista**: Renderiza la interfaz al usuario

### 2. Controladores - `backend/controllers/AuthController.php`

```php
class AuthController extends BaseController
{
    public function enviarCredenciales()
    {
        $email = $_POST['email'] ?? '';
        $contrasena = $_POST['contrasena'] ?? '';

        // Verifica credenciales con el modelo
        $userEntity = $this->userModel->verificarUsuario($email, $contrasena);

        if ($userEntity) {
            // Establece sesión
            $_SESSION['user_id'] = $userEntity->getIdFamilia();
            $_SESSION['user_email'] = $userEntity->getEmail();

            // Redirige a selección de perfil
            header('Location: ' . URLROOT . '/seleccionar-perfil');
            exit;
        } else {
            // Muestra error en la vista
            $data['error'] = 'Credenciales inválidas';
            $this->view('auth/login', $data);
        }
    }

    public function processRegister()
    {
        try {
            // 1. Verifica si el email/teléfono ya existen
            $existsUser = $this->userModel->consultarExisteciaCredenciales(
                $_POST["email"],
                $_POST["telefono"]
            );

            if ($existsUser) {
                throw new Exception("{$existsUser} ya existe");
            }

            // 2. Crea la familia (usuario principal)
            $userId = $this->userModel->crearUsuario(
                $_POST['email'],
                $_POST['telefono'],
                $_POST['password']
            );

            // 3. Crea el administrador de la familia
            $adminId = $this->userModel->crearUsuarioAdmin(
                $userId,
                $_POST['adminNombre'],
                $_POST['adminNacimiento'],
                $_POST['contraPersonal']
            );

            // 4. Crea miembros adicionales si existen
            if (!empty($_POST['miembroNombre'])) {
                foreach ($_POST['miembroNombre'] as $index => $nombre) {
                    if (!empty(trim($nombre))) {
                        $this->userModel->crearUsuarioMiembro(
                            $userId,
                            $nombre,
                            $_POST['miembroNacimiento'][$index] ?? null,
                            $_POST['miembroContra'][$index] ?? '',
                            $_POST['miembroRol'][$index] ?? 'miembro'
                        );
                    }
                }
            }

            // Éxito - redirige al login
            $_SESSION['success_message'] = 'Cuenta creada exitosamente';
            header('Location: ' . URLROOT . '/login');
            exit;

        } catch (Exception $e) {
            // Error - muestra en formulario
            $data['error'] = $e->getMessage();
            $this->view('auth/registrar', $data);
        }
    }
}
```

### 3. Modelos - `backend/models/UserModel.php`

```php
class UserModel extends BaseModel
{
    public function verificarUsuario($email, $password)
    {
        // 1. Busca la familia por email
        $sqlFamilia = "SELECT id, contrasena FROM familia WHERE correo = ?";
        $familia = $this->executeQuery($sqlFamilia, [$email]);

        if (empty($familia)) {
            return false;
        }

        $familiaData = $familia[0];

        // 2. Verifica contraseña hasheada
        if (!password_verify($password, $familiaData['contrasena'])) {
            return false;
        }

        // 3. Obtiene el usuario administrador
        $sqlUsuario = "SELECT id, nombre, rol FROM usuario
                      WHERE id_familia = ? AND rol = 'admin'";
        $usuario = $this->executeQuery($sqlUsuario, [$familiaData['id']]);

        if (!empty($usuario)) {
            $userData = $usuario[0];
            return new EntityUser(
                $userData['id'],
                $email,
                $familiaData['contrasena'],
                $userData['nombre'],
                $familiaData['id']
            );
        }

        return false;
    }
}
```

### 4. Entidades - `backend/entities/EntityUser.php`

```php
class EntityUser {
    private $id;
    private $email;
    private $password;
    private $nombre;
    private $idFamilia;

    // Representa un usuario del sistema con sus datos básicos
    // Se usa para transferir datos entre capas de forma tipada
}
```

### 5. Validadores - `backend/validators/RegisterValidator.php`

```php
class RegisterValidator extends ValidatorBase {
    public static function validate($postData): array {
        $validator = new self($postData);
        $errors = [];
        $cleanData = [];

        // Validación de campos de familia
        $result = $validator->validateFields([
            'email' => $postData['email'] ?? '',
            'telefono' => $postData['telefono'] ?? ''
        ]);

        $errors = array_merge($errors, $result['errors']);
        $cleanData = array_merge($cleanData, $result['cleanData']);

        // Validación de contraseñas
        if (!$validator->isValidContrasena($postData['password'] ?? '')) {
            $errors['password'] = $validator->getErrorMessage('contrasena');
        }

        if (!$validator->passwordsMatch(
            $postData['password'] ?? '',
            $postData['confirmPassword'] ?? ''
        )) {
            $errors['confirmPassword'] = 'Las contraseñas no coinciden';
        }

        // Validación de administrador
        if (!$validator->isValidNombre($postData['adminNombre'] ?? '')) {
            $errors['adminNombre'] = $validator->getErrorMessage('nombre');
        }

        // ... más validaciones

        return [
            'success' => empty($errors),
            'errors' => $errors,
            'cleanData' => $cleanData
        ];
    }
}
```

### 6. Vistas - `frontend/views/auth/login.php`

```php
<!DOCTYPE html>
<html lang="es">
<head>
  <meta charset="UTF-8">
  <title><?php echo htmlspecialchars($title ?? 'Login - Arka'); ?></title>
  <link rel="stylesheet" href="<?php echo URLROOT; ?>/frontend/assets/css/auth.css">
</head>
<body>
  <form method="POST" action="<?php echo URLROOT; ?>/login">
      <img src="<?php echo URLROOT; ?>/frontend/assets/img/arka_logo.png" alt="Logo Arka">
      <h2>Iniciar sesión</h2>

      <!-- Muestra errores de validación -->
      <?php if (!empty($validation_errors)): ?>
          <ul class="errors">
              <?php foreach ($validation_errors as $field => $msg): ?>
                  <li style="color: red;"><?php echo htmlspecialchars($msg); ?></li>
              <?php endforeach; ?>
          </ul>
      <?php endif; ?>

      <!-- Muestra errores generales -->
      <?php if (!empty($error)): ?>
          <p style="color:red"><?php echo htmlspecialchars($error); ?></p>
      <?php endif; ?>

      <input type="email" name="email" placeholder="Email" required>
      <input type="password" name="contrasena" placeholder="Contraseña" required>

      <button type="submit" class="btn btn-oscuro">Entrar</button>

      <p>¿No tienes una cuenta?
         <a href="<?php echo URLROOT; ?>/registro">Crea una nueva aquí!</a>
      </p>
  </form>
</body>
</html>
```

## 🗃️ Estructura de Base de Datos

### Tablas Principales:

- **familia**: Familias/grupos (email, teléfono, contraseña)
- **usuario**: Miembros de cada familia (nombre, rol, contraseña personal)
- **concepto**: Categorías de gastos/ingresos
- **concepto_usuarios**: Configuración personal por usuario
- **movimiento**: Registros financieros
- **iconos**: Íconos para categorías

## 🔐 Sistema de Autenticación

### Niveles de Acceso:

1. **Guest**: Solo puede acceder a login/registro
2. **Authenticated**: Usuario logueado pero sin perfil seleccionado
3. **Family Member**: Usuario con perfil familiar seleccionado

### Middleware:

```php
class AuthMiddleware {
    public static function guest() {
        if (isset($_SESSION['user_id'])) {
            header('Location: ' . URLROOT . '/dashboard');
            exit;
        }
    }

    public static function auth() {
        if (!isset($_SESSION['user_id'])) {
            header('Location: ' . URLROOT . '/login');
            exit;
        }
    }

    public static function familyOnly() {
        if (!isset($_SESSION['user_id']) || !isset($_SESSION['miembro_id'])) {
            header('Location: ' . URLROOT . '/login');
            exit;
        }
    }
}
```

---

## 🎯 Funcionalidades

### Autenticación

- Registro de nuevos usuarios
- Login con email y contraseña
- Selección de perfiles familiares
- Middleware de protección de rutas

### Gestión de Conceptos

- Crear los conceptos
- Asignar tipos (gasto/ingreso)
- Asociar íconos y colores
- Gestión de estado (activo/inactivo)

---

## 🔒 Seguridad

- Validación de formularios en backend
- Protección contra SQL Injection con PDO
- Sanitización de datos de entrada
- Middleware de autenticación
- Manejo seguro de sesiones

---

## 🐛 Troubleshooting

### Problemas Comunes

**Error de conexión a MySQL:**

- Verificar credenciales en `db.php`
- Asegurar que MySQL esté corriendo
- Verificar permisos de usuario

**Assets no cargan:**

- Verificar rutas en templates
- Revisar configuración de Apache
- Verificar permisos de archivos

**Error 404 en rutas:**

- Verificar mod_rewrite habilitado
- Revisar configuración de .htaccess
- Confirmar DocumentRoot correcto

---

## 📄 Licencia

Este proyecto está bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para detalles.

---

## 👥 Autores

- **Tu Nombre** - _Desarrollo inicial_ - [TuUsuario](https://github.com/tuusuario)

## 🙏 Agradecimientos

- Equipo de desarrollo
- Comunidad de PHP
- Documentación de Fly.io

---
