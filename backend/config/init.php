<?php

define('ROOT', dirname(__DIR__, 2)); 
define('DEBUG', true);
error_reporting(E_ALL); ini_set('display_errors', 1);

require_once ROOT . '/backend/config/app.php';


if (session_status() === PHP_SESSION_NONE) {
    session_start();
}


require ROOT . '/backend/config/db.php';
require ROOT . '/backend/commons/BaseModel.php';     
require ROOT . '/backend/commons/BaseController.php';

?>