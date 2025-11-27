SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;
SET collation_connection = 'utf8mb4_unicode_ci';

-- ARCHIVO: 02_usuario.sql
-- Proyecto: ARKA
-- Descripción: Procedimientos almacenados para gestión de usuarios


DELIMITER $$

-- Procedure: sp_create_usuario
-- Descripción: Crea un nuevo usuario asociado a una familia
-- Tipo de operación: Escritura
-- Tablas involucradas: usuario
-- Dependencias: Debe existir la familia (correoFamilia en tabla familia)

DROP PROCEDURE IF EXISTS sp_create_usuario$$

CREATE PROCEDURE sp_create_usuario(
    IN p_nombreUsuario VARCHAR(255),         -- Nombre de usuario único
    IN p_rol VARCHAR(50),                    -- Rol del usuario en la familia
    IN p_contrasenaPersonal VARCHAR(255),    -- Contraseña personal hasheada
    IN p_nombrePersonal VARCHAR(255),        -- Nombre real del usuario
    IN p_correoFamilia VARCHAR(255)          -- Correo de la familia a la que pertenece
)
BEGIN

    -- INSERCIÓN DE NUEVO USUARIO

    -- Se crea el usuario con sus datos de perfil y autenticación
    INSERT INTO usuario (nombreUsuario, rol, contraseñaPersonal, nombrePersonal, correoFamilia) VALUES (p_nombreUsuario, p_rol, p_contrasenaPersonal, p_nombrePersonal, p_correoFamilia);
    


    -- Retorno:
    -- id: ID del usuario recién creado
    SELECT LAST_INSERT_ID() AS id;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_find_usuario_by_email
-- Descripción: Busca todos los usuarios activos de una familia
-- Tipo de operación: Lectura
-- Tablas involucradas: usuario
-- Dependencias: Ninguna

DROP PROCEDURE IF EXISTS sp_find_usuario_by_email$$

CREATE PROCEDURE sp_find_usuario_by_email(
    IN p_correoFamilia VARCHAR(255)  -- Correo de la familia
)
BEGIN

    -- BÚSQUEDA DE USUARIOS POR FAMILIA

    -- Retorna todos los usuarios activos asociados a una familia
    


    -- Retorno: Lista de usuarios con sus datos completos (sin orden específico)
    SELECT nombreUsuario, rol, contraseñaPersonal, nombrePersonal, correoFamilia, delete_at FROM usuario WHERE correoFamilia = p_correoFamilia AND delete_at IS NULL;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_find_usuario_by_nombre
-- Descripción: Busca un usuario específico por su nombre de usuario
-- Tipo de operación: Lectura
-- Tablas involucradas: usuario
-- Dependencias: Ninguna

DROP PROCEDURE IF EXISTS sp_find_usuario_by_nombre$$

CREATE PROCEDURE sp_find_usuario_by_nombre(
    IN p_nombreUsuario VARCHAR(255)  -- Nombre de usuario a buscar
)
BEGIN

    -- BÚSQUEDA DE USUARIO POR NOMBRE

    -- Retorna los datos completos del usuario si existe y está activo
    


    -- Retorno: Datos del usuario o conjunto vacío si no existe
    SELECT nombreUsuario, rol, contraseñaPersonal, nombrePersonal, correoFamilia, delete_at FROM usuario WHERE nombreUsuario = p_nombreUsuario AND delete_at IS NULL;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_find_usuarios_by_familia
-- Descripción: Obtiene usuarios de una familia ordenados por rol y nombre
-- Tipo de operación: Lectura
-- Tablas involucradas: usuario
-- Dependencias: Ninguna

DROP PROCEDURE IF EXISTS sp_find_usuarios_by_familia$$

CREATE PROCEDURE sp_find_usuarios_by_familia(
    IN p_correoFamilia VARCHAR(255)  -- Correo de la familia
)
BEGIN

    -- LISTADO ORDENADO DE USUARIOS

    -- Retorna usuarios ordenados primero por rol (descendente) y luego por nombre personal (ascendente)
    


    -- Retorno: Lista de usuarios con orden jerárquico
    SELECT nombreUsuario, rol, contraseñaPersonal, nombrePersonal, correoFamilia, delete_at FROM usuario WHERE correoFamilia = p_correoFamilia AND delete_at IS NULL ORDER BY rol DESC, nombrePersonal ASC;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_get_all_usuarios_by_familia
-- Descripción: Obtiene todos los usuarios de una familia ordenados por nombre descendente
-- Tipo de operación: Lectura
-- Tablas involucradas: usuario
-- Dependencias: Ninguna

DROP PROCEDURE IF EXISTS sp_get_all_usuarios_by_familia$$

CREATE PROCEDURE sp_get_all_usuarios_by_familia(
    IN p_correoFamilia VARCHAR(255)  -- Correo de la familia
)
BEGIN

    -- LISTADO DE USUARIOS POR NOMBRE

    -- Retorna usuarios ordenados por nombre personal en orden descendente
    


    -- Retorno: Lista de usuarios ordenados por nombre
    SELECT nombreUsuario, rol, contraseñaPersonal, nombrePersonal, correoFamilia, delete_at FROM usuario WHERE correoFamilia = p_correoFamilia AND delete_at IS NULL ORDER BY nombrePersonal DESC;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_update_usuario_rol
-- Descripción: Actualiza el rol de un usuario específico
-- Tipo de operación: Escritura
-- Tablas involucradas: usuario
-- Dependencias: Ninguna

DROP PROCEDURE IF EXISTS sp_update_usuario_rol$$

CREATE PROCEDURE sp_update_usuario_rol(
    IN p_nombreUsuario VARCHAR(255),  -- Nombre del usuario a actualizar
    IN p_rol TINYINT                  -- Nuevo rol a asignar
)
BEGIN

    -- ACTUALIZACIÓN DE ROL

    -- Modifica el rol del usuario si existe y está activo
    UPDATE usuario SET rol = p_rol WHERE nombreUsuario = p_nombreUsuario AND delete_at IS NULL;
    


    -- Retorno:
    -- rowsAffected: Número de filas actualizadas (0 o 1)
    SELECT ROW_COUNT() AS rowsAffected;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_update_usuario_nombre_personal
-- Descripción: Actualiza el nombre personal (nombre real) de un usuario
-- Tipo de operación: Escritura
-- Tablas involucradas: usuario
-- Dependencias: Ninguna

DROP PROCEDURE IF EXISTS sp_update_usuario_nombre_personal$$

CREATE PROCEDURE sp_update_usuario_nombre_personal(
    IN p_nombreUsuario VARCHAR(255),  -- Nombre del usuario a actualizar
    IN p_nuevoNombre VARCHAR(255)     -- Nuevo nombre personal
)
BEGIN

    -- ACTUALIZACIÓN DE NOMBRE PERSONAL

    -- Modifica el nombre real del usuario
    UPDATE usuario SET nombrePersonal = p_nuevoNombre WHERE nombreUsuario = p_nombreUsuario AND delete_at IS NULL;
    


    -- Retorno:
    -- rowsAffected: Número de filas actualizadas (0 o 1)
    SELECT ROW_COUNT() AS rowsAffected;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_update_usuario_nombre_usuario
-- Descripción: Actualiza el nombre de usuario verificando que no exista duplicado
-- Tipo de operación: Escritura
-- Tablas involucradas: usuario
-- Dependencias: Ninguna

DROP PROCEDURE IF EXISTS sp_update_usuario_nombre_usuario$$

CREATE PROCEDURE sp_update_usuario_nombre_usuario(
    IN p_viejoUsuario VARCHAR(255),  -- Nombre de usuario actual
    IN p_nuevoUsuario VARCHAR(255)   -- Nuevo nombre de usuario deseado
)
BEGIN

    -- VALIDACIÓN DE DISPONIBILIDAD

    -- Verificar si el nuevo nombre de usuario ya existe
    DECLARE usuario_existe INT;
    
    SELECT COUNT(*) INTO usuario_existe FROM usuario WHERE nombreUsuario = p_nuevoUsuario AND delete_at IS NULL;
    
    IF usuario_existe > 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'El nombre de usuario ya existe';
    END IF;
    

    -- ACTUALIZACIÓN DE NOMBRE DE USUARIO

    -- Actualizar el nombre de usuario si no hay conflictos
    UPDATE usuario SET nombreUsuario = p_nuevoUsuario WHERE nombreUsuario = p_viejoUsuario AND delete_at IS NULL;
    


    -- Retorno:
    -- rowsAffected: Número de filas actualizadas (0 o 1)
    SELECT ROW_COUNT() AS rowsAffected;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_update_usuario_password
-- Descripción: Actualiza la contraseña personal de un usuario
-- Tipo de operación: Escritura
-- Tablas involucradas: usuario
-- Dependencias: Ninguna

DROP PROCEDURE IF EXISTS sp_update_usuario_password$$

CREATE PROCEDURE sp_update_usuario_password(
    IN p_nombreUsuario VARCHAR(255),     -- Nombre del usuario
    IN p_hashedPassword VARCHAR(255)     -- Nueva contraseña hasheada
)
BEGIN

    -- ACTUALIZACIÓN DE CONTRASEÑA

    -- Modifica la contraseña personal del usuario
    UPDATE usuario SET contraseñaPersonal = p_hashedPassword WHERE nombreUsuario = p_nombreUsuario AND delete_at IS NULL;
    


    -- Retorno:
    -- rowsAffected: Número de filas actualizadas (0 o 1)
    SELECT ROW_COUNT() AS rowsAffected;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_update_usuario
-- Descripción: Actualiza nombre personal y rol de un usuario simultáneamente
-- Tipo de operación: Escritura
-- Tablas involucradas: usuario
-- Dependencias: Ninguna

DROP PROCEDURE IF EXISTS sp_update_usuario$$

CREATE PROCEDURE sp_update_usuario(
    IN p_nombreUsuario VARCHAR(255),   -- Nombre del usuario a actualizar
    IN p_nombrePersonal VARCHAR(255),  -- Nuevo nombre personal
    IN p_rol TINYINT                   -- Nuevo rol
)
BEGIN

    -- ACTUALIZACIÓN MÚLTIPLE DE DATOS

    -- Actualiza nombre personal y rol en una sola operación
    UPDATE usuario SET nombrePersonal = p_nombrePersonal, rol = p_rol WHERE nombreUsuario = p_nombreUsuario AND delete_at IS NULL;
    


    -- Retorno:
    -- rowsAffected: Número de filas actualizadas (0 o 1)
    SELECT ROW_COUNT() AS rowsAffected;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_soft_delete_usuario
-- Descripción: Realiza eliminación lógica de un usuario marcando fecha de eliminación
-- Tipo de operación: Escritura
-- Tablas involucradas: usuario
-- Dependencias: Ninguna

DROP PROCEDURE IF EXISTS sp_soft_delete_usuario$$

CREATE PROCEDURE sp_soft_delete_usuario(
    IN p_nombreUsuario VARCHAR(255)  -- Nombre del usuario a eliminar lógicamente
)
BEGIN

    -- ELIMINACIÓN LÓGICA

    -- Marca el registro con fecha de eliminación sin borrarlo físicamente
    UPDATE usuario SET delete_at = NOW() WHERE nombreUsuario = p_nombreUsuario AND delete_at IS NULL;
    


    -- Retorno:
    -- rowsAffected: Número de filas marcadas como eliminadas (0 o 1)
    SELECT ROW_COUNT() AS rowsAffected;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_get_usuario_password
-- Descripción: Obtiene la contraseña personal de un usuario para autenticación
-- Tipo de operación: Lectura
-- Tablas involucradas: usuario
-- Dependencias: Ninguna

DROP PROCEDURE IF EXISTS sp_get_usuario_password$$

CREATE PROCEDURE sp_get_usuario_password(
    IN p_nombreUsuario VARCHAR(255)  -- Nombre del usuario
)
BEGIN

    -- OBTENCIÓN DE CONTRASEÑA

    -- Retorna la contraseña hasheada para validación
    


    -- Retorno:
    -- contraseñaPersonal: Contraseña hasheada del usuario
    SELECT contraseñaPersonal FROM usuario WHERE nombreUsuario = p_nombreUsuario AND delete_at IS NULL;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_usuario_exists
-- Descripción: Verifica si existe un usuario activo con el nombre especificado
-- Tipo de operación: Lectura
-- Tablas involucradas: usuario
-- Dependencias: Ninguna

DROP PROCEDURE IF EXISTS sp_usuario_exists$$

CREATE PROCEDURE sp_usuario_exists(
    IN p_nombreUsuario VARCHAR(255)  -- Nombre del usuario a verificar
)
BEGIN

    -- VERIFICACIÓN DE EXISTENCIA

    -- Cuenta cuántos usuarios activos tienen ese nombre
    


    -- Retorno:
    -- count: Número de usuarios con ese nombre (0 o 1)
    SELECT COUNT(*) AS count FROM usuario WHERE nombreUsuario = p_nombreUsuario AND delete_at IS NULL;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_delete_usuario_permanently
-- Descripción: Elimina permanentemente un usuario del sistema
-- Tipo de operación: Escritura
-- Tablas involucradas: usuario
-- Dependencias: Puede afectar registros relacionados en movimiento, personalizacionconcepto

DROP PROCEDURE IF EXISTS sp_delete_usuario_permanently$$

CREATE PROCEDURE sp_delete_usuario_permanently(
    IN p_nombreUsuario VARCHAR(255)  -- Nombre del usuario a eliminar
)
BEGIN

    -- ELIMINACIÓN PERMANENTE

    -- Elimina físicamente el registro de la base de datos
    -- ATENCIÓN: Esta es una eliminación permanente, no reversible
    DELETE FROM usuario WHERE nombreUsuario = p_nombreUsuario;
    


    -- Retorno:
    -- rowsAffected: Número de filas eliminadas (0 o 1)
    SELECT ROW_COUNT() AS rowsAffected;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_get_usuarios_by_familia
-- Descripción: Obtiene todos los usuarios activos de una familia
-- Tipo de operación: Lectura
-- Tablas involucradas: usuario
-- Dependencias: Ninguna

DROP PROCEDURE IF EXISTS sp_get_usuarios_by_familia$$

CREATE PROCEDURE sp_get_usuarios_by_familia(
    IN p_correoFamilia VARCHAR(255)  -- Correo de la familia
)
BEGIN

    -- LISTADO DE USUARIOS POR FAMILIA

    -- Retorna todos los usuarios activos asociados a la familia
    

    -- Retorno: Lista de usuarios con sus datos completos
    SELECT nombreUsuario, rol, contrasenaPersonal, nombrePersonal, correoFamilia, delete_at FROM usuario WHERE correoFamilia = p_correoFamilia AND delete_at IS NULL;
END$$

DELIMITER ;