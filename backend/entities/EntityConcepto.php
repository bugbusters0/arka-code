<?php
// backend/entities/EntityConcepto.php

class EntityConcepto
{
  private $id;
  private $id_icono;
  private $nombre;
  private $tipo;
  private $color;
  private $create_at;
  private $update_at;
  private $delete_at;

  public function __construct(
    $id,
    $id_icono,
    $nombre,
    $tipo, // ← Parámetro requerido antes de los opcionales
    $color = null,
    $create_at = null,
    $update_at = null,
    $delete_at = null
  ) {
    $this->id = $id;
    $this->id_icono = $id_icono;
    $this->nombre = $nombre;
    $this->tipo = $tipo;
    $this->color = $color;
    $this->create_at = $create_at;
    $this->update_at = $update_at;
    $this->delete_at = $delete_at;
  }

  public function getId()
  {
    return $this->id;
  }

  public function getIdIcono()
  {
    return $this->id_icono;
  }

  public function getNombre()
  {
    return $this->nombre;
  }

  public function getTipo()
  {
    return $this->tipo;
  }

  public function getColor()
  {
    return $this->color;
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

  public function isIngreso()
  {
    return $this->tipo === 'ingreso';
  }

  public function isGasto()
  {
    return $this->tipo === 'gasto';
  }
}
