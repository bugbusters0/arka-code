SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;
SET collation_connection = 'utf8mb4_unicode_ci';

-- ARCHIVO: 01_familia.sql
-- Proyecto: ARKA
-- Descripción: Procedimientos almacenados para gestión de familias


DELIMITER $$

-- Procedure: sp_familia_create
-- Descripción: Crea un nuevo registro de familia en el sistema
-- Tipo de operación: Escritura
-- Tablas involucradas: familia
-- Dependencias: Ninguna

DROP PROCEDURE IF EXISTS sp_familia_create$$

CREATE PROCEDURE sp_familia_create(
    IN p_correo VARCHAR(255),        -- Correo electrónico de la familia
    IN p_telefono VARCHAR(20),       -- Número de teléfono de contacto
    IN p_contraseña VARCHAR(255),    -- Contraseña hasheada de la familia
    OUT p_id_familia INT             -- ID generado del nuevo registro de familia
)
BEGIN

    -- INSERCIÓN DE NUEVA FAMILIA

    -- Se crea un nuevo registro con los datos básicos de contacto y autenticación
    INSERT INTO familia (correo, telefono, contraseña) VALUES (p_correo, p_telefono, p_contraseña);
    
   -- Retorno: 
    -- p_id_familia: ID de la familia recién creada
    SET p_id_familia = LAST_INSERT_ID();
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_familia_find_by_email
-- Descripción: Busca una familia activa por su correo electrónico
-- Tipo de operación: Lectura
-- Tablas involucradas: familia
-- Dependencias: Ninguna

DROP PROCEDURE IF EXISTS sp_familia_find_by_email$$

CREATE PROCEDURE sp_familia_find_by_email(
    IN p_correo VARCHAR(255)  -- Correo electrónico de la familia a buscar
)
BEGIN

    -- BÚSQUEDA DE FAMILIA POR CORREO

    -- Retorna los datos de la familia si existe y no está eliminada
    
   -- Retorno:
    -- correo: Correo electrónico de la familia
    -- telefono: Número de teléfono de contacto
    -- contraseña: Contraseña hasheada
    -- delete_at: Fecha de eliminación lógica (NULL si está activa)
    SELECT correo, telefono, contraseña, delete_at FROM familia WHERE correo = p_correo AND delete_at IS NULL;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_familia_get_all
-- Descripción: Obtiene todas las familias activas del sistema
-- Tipo de operación: Lectura
-- Tablas involucradas: familia
-- Dependencias: Ninguna

DROP PROCEDURE IF EXISTS sp_familia_get_all$$

CREATE PROCEDURE sp_familia_get_all()
BEGIN

    -- LISTADO DE TODAS LAS FAMILIAS ACTIVAS

    -- Retorna todas las familias no eliminadas ordenadas por ID
    
   -- Retorno:
    -- id_familia: Identificador único de la familia
    -- correo: Correo electrónico de la familia
    -- telefono: Número de teléfono de contacto
    -- contraseña: Contraseña hasheada
    -- delete_at: Fecha de eliminación lógica (NULL para activas)
    SELECT id_familia, correo, telefono, contraseña, delete_at FROM familia WHERE delete_at IS NULL ORDER BY id_familia;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_get_password_familia
-- Descripción: Obtiene la contraseña hasheada de una familia para autenticación
-- Tipo de operación: Lectura
-- Tablas involucradas: familia
-- Dependencias: Ninguna

DROP PROCEDURE IF EXISTS sp_get_password_familia$$

CREATE PROCEDURE sp_get_password_familia(
    IN p_correo VARCHAR(255)  -- Correo electrónico de la familia
)
BEGIN

    -- OBTENCIÓN DE CONTRASEÑA

    -- Retorna la contraseña hasheada para validación de autenticación
    
   -- Retorno:
    -- contraseña: Contraseña hasheada de la familia
    SELECT contraseña FROM familia WHERE correo = p_correo AND delete_at IS NULL LIMIT 1;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_delete_familia
-- Descripción: Elimina permanentemente una familia del sistema
-- Tipo de operación: Escritura
-- Tablas involucradas: familia
-- Dependencias: Puede afectar registros relacionados en usuario, concepto, movimiento, etc.

DROP PROCEDURE IF EXISTS sp_delete_familia$$

CREATE PROCEDURE sp_delete_familia(
    IN p_correo VARCHAR(255)  -- Correo electrónico de la familia a eliminar
)
BEGIN

    -- ELIMINACIÓN PERMANENTE DE FAMILIA

    -- Elimina físicamente el registro de la base de datos
    -- ATENCIÓN: Esta es una eliminación permanente, no lógica
    DELETE FROM familia WHERE correo = p_correo;
    
   -- Retorno:
    -- rows_affected: Número de filas eliminadas (0 o 1)
    SELECT ROW_COUNT() AS rows_affected;
END$$

DELIMITER ;