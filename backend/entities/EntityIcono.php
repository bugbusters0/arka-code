<?php
// backend/entities/EntityIcono.php

class EntityIcono
{
  private $id;
  private $path;
  private $nombre;
  private $size;
  private $type;
  private $create_at;
  private $update_at;
  private $delete_at;

  public function __construct(
    $id,
    $path,
    $nombre,
    $size = null,
    $type = 'image/svg+xml',
    $create_at = null,
    $update_at = null,
    $delete_at = null
  ) {
    $this->id = $id;
    $this->path = $path;
    $this->nombre = $nombre;
    $this->size = $size;
    $this->type = $type;
    $this->create_at = $create_at;
    $this->update_at = $update_at;
    $this->delete_at = $delete_at;
  }

  public function getId()
  {
    return $this->id;
  }

  public function getPath()
  {
    return $this->path;
  }

  public function getNombre()
  {
    return $this->nombre;
  }

  public function getSize()
  {
    return $this->size;
  }

  public function getType()
  {
    return $this->type;
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

  public function isFontAwesome()
  {
    return strpos($this->path, 'fa-') !== false;
  }

  public function isImage()
  {
    return !$this->isFontAwesome();
  }

  public function getFullPath()
  {
    if ($this->isFontAwesome()) {
      return $this->path; // Para FontAwesome, el path es el nombre de la clase
    }

    // Para imágenes, retornar la ruta completa
    return URLROOT . '/' . ltrim($this->path, '/');
  }

  public function getDisplayName()
  {
    return ucfirst(str_replace(['-', '_'], ' ', $this->nombre));
  }

  public function getFileExtension()
  {
    if ($this->isFontAwesome()) {
      return 'font-awesome';
    }

    $pathInfo = pathinfo($this->path);
    return $pathInfo['extension'] ?? 'unknown';
  }

  public function isSvg()
  {
    return $this->type === 'image/svg+xml' ||
      strtolower($this->getFileExtension()) === 'svg';
  }

  public function isPng()
  {
    return $this->type === 'image/png' ||
      strtolower($this->getFileExtension()) === 'png';
  }

  public function isJpg()
  {
    return $this->type === 'image/jpeg' ||
      strtolower($this->getFileExtension()) === 'jpg' ||
      strtolower($this->getFileExtension()) === 'jpeg';
  }

  public function getHumanReadableSize()
  {
    if ($this->size === null) {
      return 'Desconocido';
    }

    $units = ['B', 'KB', 'MB', 'GB'];
    $bytes = $this->size;
    $factor = floor((strlen($bytes) - 1) / 3);

    if ($factor > 0) {
      $bytes = $bytes / pow(1024, $factor);
    }

    return sprintf('%.2f %s', $bytes, $units[$factor]);
  }

  public function toArray()
  {
    return [
      'id' => $this->id,
      'path' => $this->path,
      'nombre' => $this->nombre,
      'size' => $this->size,
      'type' => $this->type,
      'create_at' => $this->create_at,
      'update_at' => $this->update_at,
      'delete_at' => $this->delete_at,
      'is_font_awesome' => $this->isFontAwesome(),
      'full_path' => $this->getFullPath(),
      'display_name' => $this->getDisplayName(),
      'file_extension' => $this->getFileExtension(),
      'human_readable_size' => $this->getHumanReadableSize()
    ];
  }
}
