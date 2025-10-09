<?php

require_once ROOT . '/backend/config/db.php'; 

abstract class BaseModel {
    protected $db; 

    public function __construct() {
        $this->db = conectarBD();
    }

    // Método genérico para queries preparadas (SELECT, INSERT, etc.)
    protected function executeQuery($sql, $params = [], $fetchMode = PDO::FETCH_ASSOC) {
        try {
            $stmt = $this->db->prepare($sql);
            $stmt->execute($params);
            return $stmt->fetchAll($fetchMode); 
        } catch (PDOException $e) {
            if (defined('DEBUG') && DEBUG) {
                die("Error DB: " . $e->getMessage()); 
            }
            return false; 
        }
    }

    // Método para INSERT/UPDATE (devuelve bool éxito)
    protected function executeNonQuery($sql, $params = []) {
        try {
            $stmt = $this->db->prepare($sql);
            return $stmt->execute($params);
        } catch (PDOException $e) {
            if (defined('DEBUG') && DEBUG) {
                die("Error DB: " . $e->getMessage());
            }
            return false;
        }
    }

    // Opcional: Método genérico para find by ID
    protected function findById($table, $id) {
        $sql = "SELECT * FROM $table WHERE id = ?";
        return $this->executeQuery($sql, [$id])[0] ?? null;
    }
}
?>