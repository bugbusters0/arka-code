<?php
// backend/validators/LoginValidator.php

require_once 'ValidatorBase.php';

class LoginValidator extends ValidatorBase {
  
  public static function validate($postData): array {
    $validator = new self($postData); // Composición con base

    $fieldsToValidate = [
        'email' => trim($postData['email'] ?? ''),
        'contrasena' => $postData['contrasena'] ?? '' 
    ];

    // Usar el método heredado de ValidatorBase
    $result = $validator->validateFields($fieldsToValidate);

    if (!$result['success']) {
        return [
            'success' => false, 
            'errors' => $result['errors']
        ];
    }

    return [
        'success' => true,
        'cleanData' => $result['cleanData']
    ];
  }
}
?>