SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;
SET collation_connection = 'utf8mb4_unicode_ci';

-- ARCHIVO: 03_concepto.sql
-- Proyecto: ARKA
-- Descripción: Procedimientos almacenados para gestión de conceptos (categorías de ingresos/gastos)


DELIMITER $$

-- Procedure: sp_create_concepto
-- Descripción: Crea un nuevo concepto de ingreso o gasto para una familia
-- Tipo de operación: Escritura
-- Tablas involucradas: concepto
-- Dependencias: Debe existir la familia y el usuario

DROP PROCEDURE IF EXISTS sp_create_concepto$$

CREATE PROCEDURE sp_create_concepto(
    IN p_nombreConcepto VARCHAR(255),   -- Nombre del concepto a crear
    IN p_correoFamilia VARCHAR(255),    -- Correo de la familia propietaria
    IN p_tipo VARCHAR(50),              -- Tipo de concepto (0=gasto, 1=ingreso)
    IN p_icono VARCHAR(100),            -- Icono visual del concepto
    IN p_color VARCHAR(50),             -- Color asociado al concepto
    IN p_nombreUsuario VARCHAR(255)     -- Usuario que crea el concepto
)
BEGIN

    -- INSERCIÓN DE NUEVO CONCEPTO

    -- Se crea el concepto con sus características visuales y de clasificación
    INSERT INTO concepto (nombreConcepto, correoFamilia, tipo, icono, color, nombreUsuario) VALUES (p_nombreConcepto, p_correoFamilia, p_tipo, p_icono, p_color, p_nombreUsuario);
    

    -- RETORNO DE RESULTADO

    -- Retorno:
    -- rows_affected: Número de filas insertadas
    -- id_concepto: ID del concepto recién creado
    SELECT ROW_COUNT() AS rows_affected, LAST_INSERT_ID() AS id_concepto;
END$$

DELIMITER ;


DELIMITER $$

-- Procedure: sp_create_concepto_transaccional
-- Descripción: Crea un concepto y sus personalizaciones para todos los usuarios de la familia de forma atómica.
-- Tipo de operación: Escritura (Transaccional)
-- Tablas involucradas: concepto, personalizacionconcepto
-- Dependencias: Llama a la lógica de inserción masiva de personalizaciones.

DROP PROCEDURE IF EXISTS sp_create_concepto_transaccional$$

CREATE PROCEDURE sp_create_concepto_transaccional(
    -- Parámetros del Concepto
    IN p_nombreConcepto VARCHAR(255),
    IN p_correoFamilia VARCHAR(255),
    IN p_tipo TINYINT,
    IN p_icono VARCHAR(100),
    IN p_color VARCHAR(50),
    IN p_nombreUsuario VARCHAR(255),
    
    -- Parámetros de Personalización
    IN p_limiteGasto DECIMAL(10,2),
    IN p_montoPlanificado DECIMAL(10,2),
    IN p_tipoPeriodoPlanificado VARCHAR(15),
    IN p_tipoPeriodoLimite VARCHAR(15),
    IN p_diaDesembolsoPlanejado TINYINT,
    IN p_diaPeriodoLimite TINYINT
)
BEGIN
    -- Declarar variables para almacenar el resultado del SP de personalización
    DECLARE v_usuariosTotal INT DEFAULT 0;
    DECLARE v_usuariosCreados INT DEFAULT 0;
    DECLARE v_usuariosError INT DEFAULT 0;
    
    -- Manejador de errores para forzar ROLLBACK si ocurre cualquier SQLEXCEPTION
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
    BEGIN
        -- Si ocurre un error, revertimos todas las operaciones.
        ROLLBACK;
        -- Re-lanza un error estándar para que Go pueda capturarlo
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Transacción fallida: Error al crear concepto o personalizaciones.';
    END;

    START TRANSACTION;
    
    -- 1. CREAR EL CONCEPTO
    INSERT INTO concepto (nombreConcepto, correoFamilia, tipo, icono, color, nombreUsuario)
    VALUES (p_nombreConcepto, p_correoFamilia, p_tipo, p_icono, p_color, p_nombreUsuario);
    
    -- Si la inserción del concepto falló por alguna razón (aunque el handler lo atraparía), se detiene aquí.

    -- 2. CREAR LAS PERSONALIZACIONES MASIVAS
    -- Se llama al SP existente, pero DENTRO de esta transacción.
    CALL sp_create_personalizaciones_all_usuarios(
        p_nombreConcepto,
        p_correoFamilia,
        p_limiteGasto,
        p_montoPlanificado,
        p_tipoPeriodoPlanificado,
        p_tipoPeriodoLimite,
        p_diaDesembolsoPlanejado,
        p_diaPeriodoLimite
    );
    
    -- Se asume que el SP de personalizaciones devolverá los conteos al finalizar.
    -- (Nota: Depende de cómo esté configurado para devolver los datos al SP llamante,
    -- pero para simplificar, usaremos el mismo retorno en el SP.)
    
    -- Se asume que si el SP de personalizaciones no genera un error fatal, la inserción
    -- del concepto es válida. Si tu SP de personalizaciones retorna errores lógicos (v_usuariosError > 0), 
    -- puedes manejarlo aquí con un IF y forzar un ROLLBACK o COMMIT según tu regla de negocio.
    
    -- Si llegamos aquí sin que se active el EXIT HANDLER, hacemos COMMIT.
    COMMIT;
    
    -- Retorno opcional de éxito
    SELECT 'success' AS status;

END$$
DELIMITER ;


DELIMITER $$

-- Procedure: sp_find_conceptos_by_familia
-- Descripción: Obtiene conceptos de una familia filtrados por tipo con estado de activación
-- Tipo de operación: Lectura
-- Tablas involucradas: concepto, personalizacionconcepto
-- Dependencias: LEFT JOIN con personalizacionconcepto para obtener estado activo

DROP PROCEDURE IF EXISTS sp_find_conceptos_by_familia$$

CREATE PROCEDURE sp_find_conceptos_by_familia(
    IN p_correoFamilia VARCHAR(255),  -- Correo de la familia
    IN p_tipo TINYINT                 -- Tipo de concepto (0=gasto, 1=ingreso)
)
BEGIN

    -- CONSULTA CON ESTADO DE PERSONALIZACIÓN

    -- Obtiene conceptos y su estado activo desde personalizaciones
    -- Si no existe personalización, se considera activo por defecto (COALESCE)
    

    -- RETORNO DE RESULTADO

    -- Retorno: Lista de conceptos con su estado de activación
    -- activo: Estado de activación (1 por defecto si no hay personalización)
    SELECT c.nombreConcepto, c.correoFamilia, c.tipo, c.icono, c.color, c.nombreUsuario, COALESCE(pc.activo, 1) AS activo, c.delete_at FROM concepto c LEFT JOIN personalizacionconcepto pc ON c.nombreConcepto = pc.nombreConcepto AND c.correoFamilia = pc.correoFamilia AND pc.nombreUsuario = c.nombreUsuario WHERE c.correoFamilia = p_correoFamilia AND c.tipo = p_tipo AND c.delete_at IS NULL ORDER BY c.nombreConcepto;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_update_concepto_icono_color
-- Descripción: Actualiza el icono y color de un concepto existente
-- Tipo de operación: Escritura
-- Tablas involucradas: concepto
-- Dependencias: Ninguna

DROP PROCEDURE IF EXISTS sp_update_concepto_icono_color$$

CREATE PROCEDURE sp_update_concepto_icono_color(
    IN p_icono VARCHAR(100),            -- Nuevo icono
    IN p_color VARCHAR(50),             -- Nuevo color
    IN p_nombreConcepto VARCHAR(255),   -- Nombre del concepto a actualizar
    IN p_correoFamilia VARCHAR(255)     -- Correo de la familia
)
BEGIN

    -- ACTUALIZACIÓN DE PROPIEDADES VISUALES

    -- Modifica únicamente el icono y color del concepto
    UPDATE concepto SET icono = p_icono, color = p_color WHERE nombreConcepto = p_nombreConcepto AND correoFamilia = p_correoFamilia AND delete_at IS NULL;
    

    -- RETORNO DE RESULTADO

    -- Retorno:
    -- rows_affected: Número de filas actualizadas (0 o 1)
    SELECT ROW_COUNT() AS rows_affected;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_find_all_conceptos_by_familia
-- Descripción: Obtiene todos los conceptos de una familia sin filtro de tipo
-- Tipo de operación: Lectura
-- Tablas involucradas: concepto
-- Dependencias: Ninguna

DROP PROCEDURE IF EXISTS sp_find_all_conceptos_by_familia$$

CREATE PROCEDURE sp_find_all_conceptos_by_familia(
    IN p_correoFamilia VARCHAR(255)  -- Correo de la familia
)
BEGIN

    -- LISTADO COMPLETO DE CONCEPTOS

    -- Retorna todos los conceptos (ingresos y gastos) de la familia
    

    -- RETORNO DE RESULTADO

    -- Retorno: Lista completa de conceptos ordenados alfabéticamente
    SELECT nombreConcepto, correoFamilia, tipo, icono, color, nombreUsuario, delete_at FROM concepto WHERE correoFamilia = p_correoFamilia AND delete_at IS NULL ORDER BY nombreConcepto;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_find_concepto_by_nombre
-- Descripción: Busca un concepto específico por nombre y familia
-- Tipo de operación: Lectura
-- Tablas involucradas: concepto
-- Dependencias: Ninguna

DROP PROCEDURE IF EXISTS sp_find_concepto_by_nombre$$

CREATE PROCEDURE sp_find_concepto_by_nombre(
    IN p_nombreConcepto VARCHAR(255),  -- Nombre del concepto a buscar
    IN p_correoFamilia VARCHAR(255)    -- Correo de la familia
)
BEGIN

    -- BÚSQUEDA DE CONCEPTO ESPECÍFICO

    -- Retorna el concepto si existe y está activo
    

    -- RETORNO DE RESULTADO

    -- Retorno: Datos del concepto o conjunto vacío si no existe
    SELECT nombreConcepto, correoFamilia, tipo, icono, color, nombreUsuario, delete_at FROM concepto WHERE nombreConcepto = p_nombreConcepto AND correoFamilia = p_correoFamilia AND delete_at IS NULL;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_concepto_exists
-- Descripción: Verifica si existe un concepto en una familia
-- Tipo de operación: Lectura
-- Tablas involucradas: concepto
-- Dependencias: Ninguna

DROP PROCEDURE IF EXISTS sp_concepto_exists$$

CREATE PROCEDURE sp_concepto_exists(
    IN p_nombreConcepto VARCHAR(255),  -- Nombre del concepto a verificar
    IN p_correoFamilia VARCHAR(255)    -- Correo de la familia
)
BEGIN

    -- VERIFICACIÓN DE EXISTENCIA

    -- Cuenta cuántos conceptos activos tienen ese nombre en la familia
    

    -- RETORNO DE RESULTADO

    -- Retorno:
    -- existe: Número de conceptos con ese nombre (0 o 1)
    SELECT COUNT(*) AS existe FROM concepto WHERE nombreConcepto = p_nombreConcepto AND correoFamilia = p_correoFamilia AND delete_at IS NULL;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_get_count_by_tipo
-- Descripción: Obtiene el conteo de conceptos de gastos e ingresos de una familia
-- Tipo de operación: Lectura
-- Tablas involucradas: concepto
-- Dependencias: Ninguna

DROP PROCEDURE IF EXISTS sp_get_count_by_tipo$$

CREATE PROCEDURE sp_get_count_by_tipo(
    IN p_correoFamilia VARCHAR(255)  -- Correo de la familia
)
BEGIN

    -- CONTEO POR TIPO DE CONCEPTO

    -- Calcula cuántos conceptos de gastos e ingresos tiene la familia
    -- Usa CASE para separar por tipo: 0=gastos, 1=ingresos
    

    -- RETORNO DE RESULTADO

    -- Retorno:
    -- gastos: Número total de conceptos de tipo gasto (tipo=0)
    -- ingresos: Número total de conceptos de tipo ingreso (tipo=1)
    SELECT SUM(CASE WHEN tipo = 0 THEN 1 ELSE 0 END) AS gastos, SUM(CASE WHEN tipo = 1 THEN 1 ELSE 0 END) AS ingresos FROM concepto WHERE correoFamilia = p_correoFamilia AND delete_at IS NULL;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_find_conceptos_by_familia_activos
-- Descripción: Obtiene conceptos filtrados por tipo con estado de activación personalizado por usuario
-- Tipo de operación: Lectura
-- Tablas involucradas: concepto, personalizacionconcepto
-- Dependencias: LEFT JOIN con personalizacionconcepto para estado por usuario

DROP PROCEDURE IF EXISTS sp_find_conceptos_by_familia_activos$$

CREATE PROCEDURE sp_find_conceptos_by_familia_activos(
    IN p_correoFamilia VARCHAR(255),   -- Correo de la familia
    IN p_nombreUsuario VARCHAR(255),   -- Usuario para filtrar personalizaciones
    IN p_tipo TINYINT                  -- Tipo de concepto (0=gasto, 1=ingreso)
)
BEGIN

    -- CONSULTA CON PERSONALIZACIÓN POR USUARIO

    -- Obtiene conceptos con estado activo específico del usuario solicitante
    -- Si no existe personalización del usuario, se considera activo (TRUE) por defecto
    

    -- RETORNO DE RESULTADO

    -- Retorno: Lista de conceptos con estado de activación personalizado
    -- activo: Estado específico para el usuario (TRUE por defecto)
    SELECT c.nombreConcepto, c.correoFamilia, c.tipo, c.icono, c.color, c.nombreUsuario, COALESCE(pc.activo, TRUE) AS activo, c.delete_at FROM concepto c LEFT JOIN personalizacionconcepto pc ON c.nombreConcepto = pc.nombreConcepto AND c.correoFamilia = pc.correoFamilia AND pc.nombreUsuario = p_nombreUsuario WHERE c.correoFamilia = p_correoFamilia AND c.tipo = p_tipo AND c.delete_at IS NULL ORDER BY c.nombreConcepto;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_update_concepto
-- Descripción: Actualiza icono, color y usuario asignado de un concepto
-- Tipo de operación: Escritura
-- Tablas involucradas: concepto
-- Dependencias: Ninguna

DROP PROCEDURE IF EXISTS sp_update_concepto$$

CREATE PROCEDURE sp_update_concepto(
    IN p_icono VARCHAR(100),            -- Nuevo icono
    IN p_color VARCHAR(50),             -- Nuevo color
    IN p_nombreUsuario VARCHAR(255),    -- Usuario asignado al concepto
    IN p_nombreConcepto VARCHAR(255),   -- Nombre del concepto a actualizar
    IN p_correoFamilia VARCHAR(255)     -- Correo de la familia
)
BEGIN

    -- ACTUALIZACIÓN COMPLETA DE CONCEPTO

    -- Modifica propiedades visuales y usuario responsable
    UPDATE concepto SET icono = p_icono, color = p_color, nombreUsuario = p_nombreUsuario WHERE nombreConcepto = p_nombreConcepto AND correoFamilia = p_correoFamilia AND nombreUsuario = p_nombreUsuario;
    

    -- RETORNO DE RESULTADO

    -- Retorno:
    -- rows_affected: Número de filas actualizadas (0 o 1)
    SELECT ROW_COUNT() AS rows_affected;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_deshabilitar_concepto
-- Descripción: Desactiva un concepto y todas sus personalizaciones en transacción
-- Tipo de operación: Escritura
-- Tablas involucradas: concepto, personalizacionconcepto
-- Dependencias: Modifica ambas tablas en una transacción atómica

DROP PROCEDURE IF EXISTS sp_deshabilitar_concepto$$

CREATE PROCEDURE sp_deshabilitar_concepto(
    IN p_nombreConcepto VARCHAR(255),  -- Nombre del concepto a deshabilitar
    IN p_correoFamilia VARCHAR(255)    -- Correo de la familia
)
BEGIN
    DECLARE v_rows_concepto INT DEFAULT 0;
    DECLARE v_rows_personalizaciones INT DEFAULT 0;
    

    -- TRANSACCIÓN DE DESACTIVACIÓN

    START TRANSACTION;
    
    -- Desactivar el concepto en la tabla concepto
    UPDATE concepto SET activo = 0 WHERE nombreConcepto = p_nombreConcepto AND correoFamilia = p_correoFamilia;
    SET v_rows_concepto = ROW_COUNT();
    
    -- Desactivar todas las personalizaciones del concepto
    UPDATE personalizacionconcepto SET activo = 0 WHERE nombreConcepto = p_nombreConcepto AND correoFamilia = p_correoFamilia;
    SET v_rows_personalizaciones = ROW_COUNT();
    
    COMMIT;
    

    -- RETORNO DE RESULTADO

    -- Retorno:
    -- concepto_rows: Filas afectadas en tabla concepto
    -- personalizaciones_rows: Filas afectadas en personalizaciones
    -- estado: Estado resultante ('deshabilitado')
    SELECT v_rows_concepto AS concepto_rows, v_rows_personalizaciones AS personalizaciones_rows, 'deshabilitado' AS estado;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_habilitar_concepto
-- Descripción: Activa un concepto y todas sus personalizaciones en transacción
-- Tipo de operación: Escritura
-- Tablas involucradas: concepto, personalizacionconcepto
-- Dependencias: Modifica ambas tablas en una transacción atómica

DROP PROCEDURE IF EXISTS sp_habilitar_concepto$$

CREATE PROCEDURE sp_habilitar_concepto(
    IN p_nombreConcepto VARCHAR(255),  -- Nombre del concepto a habilitar
    IN p_correoFamilia VARCHAR(255)    -- Correo de la familia
)
BEGIN
    DECLARE v_rows_concepto INT DEFAULT 0;
    DECLARE v_rows_personalizaciones INT DEFAULT 0;
    

    -- TRANSACCIÓN DE ACTIVACIÓN

    START TRANSACTION;
    
    -- Reactivar el concepto
    UPDATE concepto SET activo = 1 WHERE nombreConcepto = p_nombreConcepto AND correoFamilia = p_correoFamilia;
    SET v_rows_concepto = ROW_COUNT();
    
    -- Reactivar todas las personalizaciones
    UPDATE personalizacionconcepto SET activo = 1 WHERE nombreConcepto = p_nombreConcepto AND correoFamilia = p_correoFamilia;
    SET v_rows_personalizaciones = ROW_COUNT();
    
    COMMIT;
    

    -- RETORNO DE RESULTADO

    -- Retorno:
    -- concepto_rows: Filas afectadas en tabla concepto
    -- personalizaciones_rows: Filas afectadas en personalizaciones
    -- estado: Estado resultante ('habilitado')
    SELECT v_rows_concepto AS concepto_rows, v_rows_personalizaciones AS personalizaciones_rows, 'habilitado' AS estado;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_toggle_activo_global
-- Descripción: Alterna el estado activo/inactivo de un concepto y sus personalizaciones
-- Tipo de operación: Escritura
-- Tablas involucradas: concepto, personalizacionconcepto
-- Dependencias: Modifica ambas tablas en una transacción atómica

DROP PROCEDURE IF EXISTS sp_toggle_activo_global$$

CREATE PROCEDURE sp_toggle_activo_global(
    IN p_nombreConcepto VARCHAR(255),  -- Nombre del concepto
    IN p_correoFamilia VARCHAR(255)    -- Correo de la familia
)
BEGIN
    DECLARE v_activo_actual BOOLEAN;
    DECLARE v_rows_concepto INT DEFAULT 0;
    DECLARE v_rows_personalizaciones INT DEFAULT 0;
    

    -- OBTENCIÓN DE ESTADO ACTUAL

    -- Obtener el estado actual del concepto
    SELECT activo INTO v_activo_actual FROM concepto WHERE nombreConcepto = p_nombreConcepto AND correoFamilia = p_correoFamilia LIMIT 1;
    

    -- TRANSACCIÓN DE ALTERNANCIA

    START TRANSACTION;
    
    IF v_activo_actual = TRUE THEN
        -- Deshabilitar concepto y personalizaciones
        UPDATE concepto SET activo = 0 WHERE nombreConcepto = p_nombreConcepto AND correoFamilia = p_correoFamilia;
        SET v_rows_concepto = ROW_COUNT();
        
        UPDATE personalizacionconcepto SET activo = 0 WHERE nombreConcepto = p_nombreConcepto AND correoFamilia = p_correoFamilia;
        SET v_rows_personalizaciones = ROW_COUNT();
        
    
        -- RETORNO - DESHABILITADO
    
        -- Retorno para estado deshabilitado
        SELECT v_rows_concepto AS concepto_rows, v_rows_personalizaciones AS personalizaciones_rows, 'deshabilitado' AS accion, FALSE AS nuevo_estado;
    ELSE
        -- Habilitar concepto y personalizaciones
        UPDATE concepto SET activo = 1 WHERE nombreConcepto = p_nombreConcepto AND correoFamilia = p_correoFamilia;
        SET v_rows_concepto = ROW_COUNT();
        
        UPDATE personalizacionconcepto SET activo = 1 WHERE nombreConcepto = p_nombreConcepto AND correoFamilia = p_correoFamilia;
        SET v_rows_personalizaciones = ROW_COUNT();
        
    
        -- RETORNO - HABILITADO
    
        -- Retorno para estado habilitado
        SELECT v_rows_concepto AS concepto_rows, v_rows_personalizaciones AS personalizaciones_rows, 'habilitado' AS accion, TRUE AS nuevo_estado;
    END IF;
    
    COMMIT;
END$$

DELIMITER ;