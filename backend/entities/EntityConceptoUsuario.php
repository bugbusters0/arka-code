<?php
// backend/entities/EntityConceptoUsuario.php

class EntityConceptoUsuario
{
  private $id;
  private $id_usuario;
  private $id_concepto;
  private $desembolso_planejado;
  private $limite_monto;
  private $periodo_desembolso;
  private $periodo_tipo;
  private $periodo_dia;
  private $limite_tipo;
  private $limite_fecha;
  private $has_notificacion;
  private $visible;
  private $create_at;
  private $update_at;
  private $delete_at;

  public function __construct(
    $id,
    $id_usuario,
    $id_concepto,
    $desembolso_planejado,
    $limite_monto, // ← Parámetros requeridos primero
    $periodo_desembolso = null,
    $periodo_tipo = null,
    $periodo_dia = null,
    $limite_tipo = null,
    $limite_fecha = null,
    $has_notificacion = false,
    $visible = true,
    $create_at = null,
    $update_at = null,
    $delete_at = null
  ) {
    $this->id = $id;
    $this->id_usuario = $id_usuario;
    $this->id_concepto = $id_concepto;
    $this->desembolso_planejado = $desembolso_planejado;
    $this->limite_monto = $limite_monto;
    $this->periodo_desembolso = $periodo_desembolso;
    $this->periodo_tipo = $periodo_tipo;
    $this->periodo_dia = $periodo_dia;
    $this->limite_tipo = $limite_tipo;
    $this->limite_fecha = $limite_fecha;
    $this->has_notificacion = $has_notificacion;
    $this->visible = $visible;
    $this->create_at = $create_at;
    $this->update_at = $update_at;
    $this->delete_at = $delete_at;
  }

  public function getId()
  {
    return $this->id;
  }

  public function getIdUsuario()
  {
    return $this->id_usuario;
  }

  public function getIdConcepto()
  {
    return $this->id_concepto;
  }

  public function getDesembolsoPlanejado()
  {
    return $this->desembolso_planejado;
  }

  public function getLimiteMonto()
  {
    return $this->limite_monto;
  }

  public function getPeriodoDesembolso()
  {
    return $this->periodo_desembolso;
  }

  public function getPeriodoTipo()
  {
    return $this->periodo_tipo;
  }

  public function getPeriodoDia()
  {
    return $this->periodo_dia;
  }

  public function getLimiteTipo()
  {
    return $this->limite_tipo;
  }

  public function getLimiteFecha()
  {
    return $this->limite_fecha;
  }

  public function hasNotificacion()
  {
    return $this->has_notificacion;
  }

  public function isVisible()
  {
    return $this->visible;
  }

  public function getCreateAt()
  {
    return $this->create_at;
  }

  public function getUpdateAt()
  {
    return $this->update_at;
  }

  public function getDeleteAt()
  {
    return $this->delete_at;
  }

  public function isDeleted()
  {
    return $this->delete_at !== null;
  }

  // Métodos de validación/business logic
  public function isPeriodoMensual()
  {
    return $this->periodo_tipo === 'mensual';
  }

  public function isPeriodoSemanal()
  {
    return $this->periodo_tipo === 'semanal';
  }

  public function isPeriodoDiario()
  {
    return $this->periodo_tipo === 'diario';
  }

  public function isPeriodoQuincenal()
  {
    return $this->periodo_tipo === 'quincenal';
  }

  public function isNotificacionActiva()
  {
    return $this->has_notificacion && $this->visible;
  }
}
