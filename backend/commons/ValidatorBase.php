<?php
// backend/validators/ValidatorBase.php

class ValidatorBase
{
    protected $email;
    protected $nombre;
    protected $numero;
    protected $contrasena;
    protected $telefono;
    protected $color;
    protected $rol;
    protected $fecha_nac;
    protected $confirmPassword;

    public function __construct($data = [])
    {
        $this->email = $data['email'] ?? '';
        $this->nombre = $data['nombre'] ?? '';
        $this->numero = (int) ($data['numero'] ?? 0);
        $this->contrasena = $data['contrasena'] ?? '';
        $this->telefono = $data['telefono'] ?? '';
        $this->color = $data['color'] ?? '';
        $this->rol = $data['rol'] ?? '';
        $this->fecha_nac = $data['fecha_nac'] ?? '';
        $this->confirmPassword = $data['confirmPassword'] ?? '';
    }

    public function isValidEmail($email = ''): bool
    {
        $emailToCheck = $email ?: $this->email;
        return filter_var($emailToCheck, FILTER_VALIDATE_EMAIL) !== false;
    }

    public function isValidNombre($nombre = ''): bool
    {
        $nombreToCheck = $nombre ?: $this->nombre;
        return strlen(trim($nombreToCheck)) >= 2
            && strlen(trim($nombreToCheck)) <= 100
            && preg_match('/^[a-zA-ZáéíóúÁÉÍÓÚñÑ\s]+$/', trim($nombreToCheck));
    }

    public function isValidNumero($numero = 0): bool
    {
        $numeroToCheck = $numero ?: $this->numero;
        return is_numeric($numeroToCheck) && $numeroToCheck > 0;
    }

    public function isValidContrasena($contrasena = ''): bool
    {
        $contrasenaToCheck = $contrasena ?: $this->contrasena;
        return strlen($contrasenaToCheck) >= 6
            && preg_match('/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)/', $contrasenaToCheck);
    }

    public function isValidTelefono($telefono = ''): bool
    {
        $telefonoToCheck = $telefono ?: $this->telefono;
        // Para Perú: acepta 9 dígitos sin espacios
        return preg_match('/^[0-9]{9}$/', $telefonoToCheck);
    }

    public function isValidColor($color = ''): bool
    {
        $colorToCheck = $color ?: $this->color;
        return preg_match('/^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$/', $colorToCheck);
    }

    public function isValidRol($rol = ''): bool
    {
        $rolToCheck = $rol ?: $this->rol;
        $rolesValidos = ['admin', 'miembro'];
        return in_array(strtolower(trim($rolToCheck)), $rolesValidos);
    }

    public function isValidFechaNacimiento($fecha = ''): bool
    {
        $fechaToCheck = $fecha ?: $this->fecha_nac;

        $fechaObj = DateTime::createFromFormat('Y-m-d', $fechaToCheck);
        if (!$fechaObj || $fechaObj->format('Y-m-d') !== $fechaToCheck) {
            return false;
        }

        // No puede ser futura
        if ($fechaObj > new DateTime()) {
            return false;
        }

        // Debe tener al menos 5 años (opcional, ajusta según necesites)
        $edad = (new DateTime())->diff($fechaObj)->y;
        return $edad >= 0;
    }

    public function isValidContraPersonal($contra = ''): bool
    {
        $contraToCheck = $contra ?: $this->contrasena;
        // PIN de 4-6 caracteres
        return strlen($contraToCheck) >= 4 && strlen($contraToCheck) <= 6;
    }

    public function passwordsMatch($password, $confirmPassword): bool
    {
        return $password === $confirmPassword;
    }

    public function sanitize($value, $type = 'string')
    {
        switch ($type) {
            case 'email':
                return filter_var($value, FILTER_SANITIZE_EMAIL);
            case 'int':
                return filter_var($value, FILTER_SANITIZE_NUMBER_INT);
            case 'string':
            default:
                return htmlspecialchars(trim($value), ENT_QUOTES, 'UTF-8');
        }
    }

    // Método helper: Valida múltiples y devuelve errores
    public function validateFields($fields = []): array
    {
        $errors = [];
        $cleanData = [];

        foreach ($fields as $field => $value) {
            $method = 'isValid' . ucfirst($field);

            if (method_exists($this, $method)) {
                if (!$this->$method($value)) {
                    $errors[$field] = $this->getErrorMessage($field);
                } else {
                    // Sanitiza solo si es válido
                    $cleanData[$field] = $this->sanitizeField($field, $value);
                }
            } else {
                // Si no hay validador específico, solo sanitiza
                $cleanData[$field] = $this->sanitize($value);
            }
        }

        return [
            'success' => empty($errors),
            'errors' => $errors,
            'cleanData' => $cleanData
        ];
    }

    // ✅ Mensajes de error personalizados
    protected function getErrorMessage($field): string
    {
        $messages = [
            'email' => 'El correo electrónico no es válido',
            'nombre' => 'El nombre debe tener entre 2 y 100 caracteres y solo letras',
            'contrasena' => 'Contraseña inválida',
            'telefono' => 'El teléfono debe tener 9 dígitos',
            'color' => 'El color debe estar en formato hexadecimal (#RRGGBB)',
            'rol' => 'El rol debe ser "admin" o "miembro"',
            'fecha_nac' => 'La fecha de nacimiento no es válida',
            'contraPersonal' => 'La contraseña personal debe tener entre 4 y 6 caracteres',
            'numero' => 'El número debe ser mayor a 0'
        ];

        return $messages[$field] ?? ucfirst($field) . ' inválido';
    }

    // ✅ Sanitización por campo
    protected function sanitizeField($field, $value)
    {
        switch ($field) {
            case 'email':
                return $this->sanitize($value, 'email');
            case 'numero':
            case 'telefono':
                return $this->sanitize($value, 'int');
            case 'contrasena':
            case 'contraPersonal':
            case 'confirmPassword':
                // NO sanitizar contraseñas
                return $value;
            default:
                return $this->sanitize($value, 'string');
        }
    }
}
