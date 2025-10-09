<?php
// backend/entities/EntityUser.php

class EntityUser {
    private $id;
    private $email;
    private $password;
    private $nombre;
    private $idFamilia;

    public function __construct($id, $email, $password, $nombre = null, $idFamilia = null) {
        $this->id = $id;
        $this->email = $email;
        $this->password = $password;
        $this->nombre = $nombre;
        $this->idFamilia = $idFamilia;
    }

    public function getId() {
        return $this->id;
    }

    public function getEmail() {
        return $this->email;
    }

    public function getPassword() {
        return $this->password;
    }
    
    public function getNombre() {
        return $this->nombre;
    }
    
    public function getIdFamilia() {
        return $this->idFamilia;
    }

    public function setEmail($email) {
        $this->email = $email;
    }
    
    public function setNombre($nombre) {
        $this->nombre = $nombre;
    }
}
?>