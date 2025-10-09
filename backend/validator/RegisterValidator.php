<?php
// backend/validators/RegisterValidator.php

require_once 'ValidatorBase.php';

class RegisterValidator extends ValidatorBase {
  
  public static function validate($postData): array {
      $validator = new self($postData); // Composición con base
      $errors = [];
      $cleanData = [];

      // ✅ 1. Validar campos individuales de familia
      $familiaFields = [
          'email' => trim($postData['email'] ?? ''),
          'telefono' => trim($postData['telefono'] ?? ''),
          'password' => $postData['password'] ?? '',
          'confirmPassword' => $postData['confirmPassword'] ?? ''
      ];

      $result = $validator->validateFields([
          'email' => $familiaFields['email'],
          'telefono' => $familiaFields['telefono']
      ]);

      $errors = array_merge($errors, $result['errors']);
      $cleanData = array_merge($cleanData, $result['cleanData']);

      // ✅ 2. Validar contraseña con reglas estrictas
      if (!$validator->isValidContrasena($familiaFields['password'])) {
          $errors['password'] = $validator->getErrorMessage('contrasena');
      } else {
          $cleanData['password'] = $familiaFields['password'];
      }

      // ✅ 3. Verificar que las contraseñas coincidan
      if (!$validator->passwordsMatch($familiaFields['password'], $familiaFields['confirmPassword'])) {
          $errors['confirmPassword'] = 'Las contraseñas no coinciden';
      }

      // ✅ 4. Validar datos del administrador
      $adminFields = [
          'adminNombre' => trim($postData['adminNombre'] ?? ''),
          'adminNacimiento' => trim($postData['adminNacimiento'] ?? ''),
          'contraPersonal' => $postData['contraPersonal'] ?? ''
      ];

      if (!$validator->isValidNombre($adminFields['adminNombre'])) {
          $errors['adminNombre'] = $validator->getErrorMessage('nombre');
      } else {
          $cleanData['adminNombre'] = $validator->sanitize($adminFields['adminNombre']);
      }

      if (!$validator->isValidFechaNacimiento($adminFields['adminNacimiento'])) {
          $errors['adminNacimiento'] = $validator->getErrorMessage('fecha_nac');
      } else {
          $cleanData['adminNacimiento'] = $adminFields['adminNacimiento'];
      }

      if (!$validator->isValidContraPersonal($adminFields['contraPersonal'])) {
          $errors['contraPersonal'] = $validator->getErrorMessage('contraPersonal');
      } else {
          $cleanData['contraPersonal'] = $adminFields['contraPersonal'];
      }

      // ✅ 5. Validar miembros opcionales
      if (!empty($postData['miembroNombre']) && is_array($postData['miembroNombre'])) {
        $cleanData['miembros'] = [];
        
        foreach ($postData['miembroNombre'] as $index => $nombre) {
          if (!empty(trim($nombre))) {
            $miembro = [
              'nombre' => trim($nombre),
              'nacimiento' => trim($postData['miembroNacimiento'][$index] ?? ''),
              'contra' => $postData['miembroContra'][$index] ?? '',
              'rol' => trim($postData['miembroRol'][$index] ?? 'miembro')
            ];

            $miembroErrors = [];

            // Validar nombre del miembro
            if (!$validator->isValidNombre($miembro['nombre'])) {
              $miembroErrors[] = "El nombre del miembro #{$index} es inválido";
            }

            // Validar fecha de nacimiento del miembro
            if (!$validator->isValidFechaNacimiento($miembro['nacimiento'])) {
              $miembroErrors[] = "La fecha de nacimiento del miembro '{$miembro['nombre']}' es inválida";
            }

            // Validar contraseña personal del miembro
            if (!$validator->isValidContraPersonal($miembro['contra'])) {
              $miembroErrors[] = "La contraseña del miembro '{$miembro['nombre']}' debe tener 4-6 caracteres";
            }

            // Validar rol del miembro
            if (!$validator->isValidRol($miembro['rol'])) {
              $miembroErrors[] = "El rol del miembro '{$miembro['nombre']}' es inválido";
            }

            if (!empty($miembroErrors)) {
              $errors["miembro_{$index}"] = implode(', ', $miembroErrors);
            } else {
              // Sanitizar y agregar miembro válido
              $cleanData['miembros'][] = [
                'nombre' => $validator->sanitize($miembro['nombre']),
                'nacimiento' => $miembro['nacimiento'],
                'contra' => $miembro['contra'], // No sanitizar contraseñas
                'rol' => strtolower($miembro['rol'])
              ];
            }
          }
        }
      }

      if (!empty($errors)) {
          return [
              'success' => false,
              'errors' => $errors
          ];
      }

      return [
          'success' => true,
          'cleanData' => $cleanData
      ];
  }
}
?>