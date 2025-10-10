<?php
// backend/validators/ConceptoValidator.php

require_once ROOT . '/backend/commons/ValidatorBase.php';

class ConceptoValidator extends ValidatorBase
{

  // Añadir propiedades específicas de concepto
  protected $tipo;
  protected $id_icono;
  protected $desembolso_planejado;
  protected $periodo_tipo;
  protected $limite_monto;
  protected $limite_tipo;

  public function __construct($data = [])
  {
    parent::__construct($data);
    $this->tipo = $data['tipo'] ?? '';
    $this->id_icono = $data['id_icono'] ?? 0;
    $this->desembolso_planejado = $data['desembolso_planejado'] ?? 0;
    $this->periodo_tipo = $data['periodo_tipo'] ?? '';
    $this->limite_monto = $data['limite_monto'] ?? 0;
    $this->limite_tipo = $data['limite_tipo'] ?? '';
  }

  public static function validate($postData): array
  {
    $validator = new self($postData);
    $errors = [];
    $cleanData = [];

    $conceptoFields = [
      'nombre' => trim($postData['nombre'] ?? ''),
      'tipo' => trim($postData['tipo'] ?? ''),
      'color' => trim($postData['color'] ?? ''),
      'id_icono' => $postData['id_icono'] ?? 0
    ];

    // Validar nombre
    if (!$validator->isValidNombreConcepto($conceptoFields['nombre'])) {
      $errors['nombre'] = $validator->getErrorMessage('nombre_concepto');
    } else {
      $cleanData['nombre'] = $validator->sanitize($conceptoFields['nombre']);
    }

    // Validar tipo
    if (!$validator->isValidTipoConcepto($conceptoFields['tipo'])) {
      $errors['tipo'] = $validator->getErrorMessage('tipo');
    } else {
      $cleanData['tipo'] = strtolower($conceptoFields['tipo']);
    }

    // Validar color (opcional)
    if (!empty($conceptoFields['color']) && !$validator->isValidColor($conceptoFields['color'])) {
      $errors['color'] = $validator->getErrorMessage('color');
    } else {
      $cleanData['color'] = $conceptoFields['color'] ?: null;
    }

    // Validar icono
    if (!$validator->isValidIcono($conceptoFields['id_icono'])) {
      $errors['id_icono'] = $validator->getErrorMessage('icono');
    } else {
      $cleanData['id_icono'] = (int) $conceptoFields['id_icono'];
    }

    if (isset($postData['configuracion_usuario']) && $postData['configuracion_usuario'] === 'true') {
      $configFields = [
        'desembolso_planejado' => $postData['desembolso_planejado'] ?? 0,
        'periodo_tipo' => $postData['periodo_tipo'] ?? '',
        'limite_monto' => $postData['limite_monto'] ?? 0,
        'limite_tipo' => $postData['limite_tipo'] ?? ''
      ];

      $configResult = $validator->validateConfiguracionUsuario($configFields);
      $errors = array_merge($errors, $configResult['errors']);
      $cleanData = array_merge($cleanData, $configResult['cleanData']);
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

  // ✅ Métodos de validación específicos para conceptos
  public function isValidNombreConcepto($nombre = ''): bool
  {
    $nombreToCheck = $nombre ?: $this->nombre;
    return strlen(trim($nombreToCheck)) >= 2
      && strlen(trim($nombreToCheck)) <= 100
      && preg_match('/^[a-zA-ZáéíóúÁÉÍÓÚñÑ0-9\s\-]+$/', trim($nombreToCheck));
  }

  public function isValidTipoConcepto($tipo = ''): bool
  {
    $tipoToCheck = $tipo ?: $this->tipo;
    $tiposValidos = ['ingreso', 'gasto'];
    return in_array(strtolower(trim($tipoToCheck)), $tiposValidos);
  }

  public function isValidIcono($id_icono = 0): bool
  {
    $idToCheck = $id_icono ?: $this->id_icono;
    return is_numeric($idToCheck) && $idToCheck > 0;
  }

  // ✅ Validar configuración de usuario para el concepto
  public function validateConfiguracionUsuario($configFields): array
  {
    $errors = [];
    $cleanData = [];

    // Validar desembolso planeado
    if (!$this->isValidMonto($configFields['desembolso_planejado'])) {
      $errors['desembolso_planejado'] = $this->getErrorMessage('monto');
    } else {
      $cleanData['desembolso_planejado'] = (float) $configFields['desembolso_planejado'];
    }

    // Validar período tipo (opcional)
    if (!empty($configFields['periodo_tipo']) && !$this->isValidPeriodoTipo($configFields['periodo_tipo'])) {
      $errors['periodo_tipo'] = $this->getErrorMessage('periodo_tipo');
    } else {
      $cleanData['periodo_tipo'] = !empty($configFields['periodo_tipo']) ? $configFields['periodo_tipo'] : null;
    }

    // Validar límite monto
    if (!$this->isValidMonto($configFields['limite_monto'])) {
      $errors['limite_monto'] = $this->getErrorMessage('monto');
    } else {
      $cleanData['limite_monto'] = (float) $configFields['limite_monto'];
    }

    // Validar límite tipo (opcional)
    if (!empty($configFields['limite_tipo']) && !$this->isValidPeriodoTipo($configFields['limite_tipo'])) {
      $errors['limite_tipo'] = $this->getErrorMessage('periodo_tipo');
    } else {
      $cleanData['limite_tipo'] = !empty($configFields['limite_tipo']) ? $configFields['limite_tipo'] : null;
    }

    return [
      'errors' => $errors,
      'cleanData' => $cleanData
    ];
  }

  public function isValidMonto($monto): bool
  {
    return is_numeric($monto) && $monto >= 0;
  }

  public function isValidPeriodoTipo($periodo_tipo): bool
  {
    $tiposValidos = ['mensual', 'quincenal', 'semanal', 'diario'];
    return in_array(strtolower(trim($periodo_tipo)), $tiposValidos);
  }

  protected function getErrorMessage($field): string
  {
    $messages = [
      'nombre_concepto' => 'El nombre del concepto debe tener entre 2 y 100 caracteres y solo puede contener letras, números y espacios',
      'tipo' => 'El tipo debe ser "ingreso" o "gasto"',
      'color' => 'El color debe estar en formato hexadecimal (#RRGGBB)',
      'icono' => 'El icono seleccionado no es válido',
      'monto' => 'El monto debe ser un número mayor o igual a 0',
      'periodo_tipo' => 'El tipo de período debe ser: mensual, quincenal, semanal o diario'
    ];

    return $messages[$field] ?? parent::getErrorMessage($field);
  }

  protected function sanitizeField($field, $value)
  {
    switch ($field) {
      case 'desembolso_planejado':
      case 'limite_monto':
        return (float) $value;
      case 'id_icono':
        return (int) $value;
      case 'periodo_tipo':
      case 'limite_tipo':
        return !empty($value) ? strtolower(trim($value)) : null;
      default:
        return parent::sanitizeField($field, $value);
    }
  }
}
