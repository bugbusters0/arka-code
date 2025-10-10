<?php
// backend/entities/EntityUser.php

class EntityMiembro
{
  private $id;
  private $id_familia;
  private $nombre;
  private $fecha_nac;
  private $rol;
  private $contra_personal;

  public function __construct($id, $id_familia, $nombre, $fecha_nac, $rol, $contra_personal = "")
  {
    $this->id = $id;
    $this->id_familia = $id_familia;
    $this->nombre = $nombre;
    $this->fecha_nac = $fecha_nac;
    $this->rol = $rol;
    $this->contra_personal = $contra_personal;
  }

  public function getId()
  {
    return $this->id;
  }

  public function getIdFamilia()
  {
    return $this->id_familia;
  }

  public function getNombre()
  {
    return $this->nombre;
  }

  public function getContraPersonal()
  {
    return $this->contra_personal;
  }

  public function getFechaNac()
  {
    return $this->fecha_nac;
  }

  public function getRol()
  {
    return $this->rol;
  }
}
