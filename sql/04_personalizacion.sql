SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;
SET collation_connection = 'utf8mb4_unicode_ci';

-- ARCHIVO: 04_personalizacion.sql
-- Proyecto: ARKA
-- Descripción: Procedimientos almacenados para gestión de personalizaciones de conceptos por usuario


DELIMITER $$

-- Procedure: sp_create_personalizacion
-- Descripción: Crea una personalización de concepto para un usuario específico
-- Tipo de operación: Escritura
-- Tablas involucradas: personalizacionconcepto
-- Dependencias: Deben existir el usuario y el concepto

DROP PROCEDURE IF EXISTS sp_create_personalizacion$$

CREATE PROCEDURE sp_create_personalizacion(
    IN p_limiteGasto DECIMAL(10,2),              -- Límite de gasto establecido
    IN p_activo BOOLEAN,                         -- Estado activo del concepto para el usuario
    IN p_montoPlanificado DECIMAL(10,2),         -- Monto planificado para el período
    IN p_tipoPeriodoPlanificado VARCHAR(50),     -- Tipo de período de planificación
    IN p_tipoPeriodoLimite VARCHAR(50),          -- Tipo de período para el límite
    IN p_diaPeriodoPlanificado TINYINT,          -- Día específico del período planificado
    IN p_diaPeriodoLimite TINYINT,               -- Día específico del período límite
    IN p_notificacion BOOLEAN,                   -- Activación de notificaciones
    IN p_nombreUsuario VARCHAR(255),             -- Usuario al que se aplica la personalización
    IN p_nombreConcepto VARCHAR(255),            -- Concepto personalizado
    IN p_correoFamilia VARCHAR(255)              -- Correo de la familia
)
BEGIN

    -- INSERCIÓN DE PERSONALIZACIÓN

    -- Se crea la personalización con todos los parámetros de configuración
    INSERT INTO personalizacionconcepto (limiteGasto, activo, montoPlanificado, tipoPeriodoPlanificado, tipoPeriodoLimite, diaPeriodoPlanificado, diaPeriodoLimite, notificacion, nombreUsuario, nombreConcepto, correoFamilia) VALUES (p_limiteGasto, p_activo, p_montoPlanificado, p_tipoPeriodoPlanificado, p_tipoPeriodoLimite, p_diaPeriodoPlanificado, p_diaPeriodoLimite, p_notificacion, p_nombreUsuario, p_nombreConcepto, p_correoFamilia);
    

   -- Retorno:
    -- rows_affected: Número de filas insertadas
    -- id_personalizacion: ID de la personalización creada
    SELECT ROW_COUNT() AS rows_affected, LAST_INSERT_ID() AS id_personalizacion;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_create_personalizaciones_all_usuarios
-- Descripción: Crea personalizaciones del mismo concepto para todos los usuarios de una familia
-- Tipo de operación: Escritura
-- Tablas involucradas: usuario, personalizacionconcepto
-- Dependencias: Lee usuarios y crea múltiples personalizaciones usando cursor

DROP PROCEDURE IF EXISTS sp_create_personalizaciones_all_usuarios$$

CREATE PROCEDURE sp_create_personalizaciones_all_usuarios(
    IN p_nombreConcepto VARCHAR(255),            -- Nombre del concepto a personalizar
    IN p_correoFamilia VARCHAR(255),             -- Correo de la familia
    IN p_limiteGasto DECIMAL(10,2),              -- Límite de gasto común
    IN p_montoPlanificado DECIMAL(10,2),         -- Monto planificado común
    IN p_tipoPeriodoPlanificado VARCHAR(50),     -- Tipo de período planificado común
    IN p_tipoPeriodoLimite VARCHAR(50),          -- Tipo de período límite común
    IN p_diaPeriodoPlanificado TINYINT,          -- Día del período planificado
    IN p_diaPeriodoLimite TINYINT                -- Día del período límite
)
BEGIN
    DECLARE v_usuarios_total INT DEFAULT 0;
    DECLARE v_usuarios_creados INT DEFAULT 0;
    DECLARE v_usuarios_error INT DEFAULT 0;
    DECLARE v_finished INTEGER DEFAULT 0;
    DECLARE v_nombreUsuario VARCHAR(255);
       -- DECLARACIÓN DE CURSOR Y HANDLERS

    -- Cursor para iterar sobre todos los usuarios activos de la familia
    DECLARE cur_usuarios CURSOR FOR SELECT nombreUsuario FROM usuario WHERE correoFamilia = p_correoFamilia AND delete_at IS NULL;
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET v_finished = 1;
    DECLARE CONTINUE HANDLER FOR SQLEXCEPTION SET v_usuarios_error = v_usuarios_error + 1;
       -- CONTEO TOTAL DE USUARIOS

    SELECT COUNT(*) INTO v_usuarios_total FROM usuario WHERE correoFamilia = p_correoFamilia AND delete_at IS NULL;
       -- ITERACIÓN Y CREACIÓN DE PERSONALIZACIONES

    OPEN cur_usuarios;
    
    usuario_loop: LOOP
        FETCH cur_usuarios INTO v_nombreUsuario;
        
        IF v_finished = 1 THEN 
            LEAVE usuario_loop;
        END IF;
        
        -- Crear personalización para el usuario actual con valores por defecto (activo=TRUE, notificacion=FALSE)
        INSERT INTO personalizacionconcepto (limiteGasto, activo, montoPlanificado, tipoPeriodoPlanificado, tipoPeriodoLimite, diaPeriodoPlanificado, diaPeriodoLimite, notificacion, nombreUsuario, nombreConcepto, correoFamilia) VALUES (p_limiteGasto, TRUE, p_montoPlanificado, p_tipoPeriodoPlanificado, p_tipoPeriodoLimite, p_diaPeriodoPlanificado, p_diaPeriodoLimite, FALSE, v_nombreUsuario, p_nombreConcepto, p_correoFamilia);
        
        SET v_usuarios_creados = v_usuarios_creados + 1;
        
    END LOOP usuario_loop;
    
    CLOSE cur_usuarios;
    

   -- Retorno:
    -- usuarios_total: Total de usuarios en la familia
    -- usuarios_creados: Personalizaciones creadas exitosamente
    -- usuarios_error: Personalizaciones que fallaron
    SELECT v_usuarios_total AS usuarios_total, v_usuarios_creados AS usuarios_creados, v_usuarios_error AS usuarios_error;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_get_personalizaciones_by_concepto
-- Descripción: Obtiene todas las personalizaciones de un concepto específico
-- Tipo de operación: Lectura
-- Tablas involucradas: personalizacionconcepto
-- Dependencias: Ninguna

DROP PROCEDURE IF EXISTS sp_get_personalizaciones_by_concepto$$

CREATE PROCEDURE sp_get_personalizaciones_by_concepto(
    IN p_nombreConcepto VARCHAR(255),  -- Nombre del concepto
    IN p_correoFamilia VARCHAR(255)    -- Correo de la familia
)
BEGIN

    -- CONSULTA DE PERSONALIZACIONES

    -- Retorna todas las personalizaciones activas del concepto
    

   -- Retorno: Lista de personalizaciones con sus configuraciones principales
    SELECT idPersonalizacion, nombreUsuario, limiteGasto, montoPlanificado, tipoPeriodoPlanificado, tipoPeriodoLimite FROM personalizacionconcepto WHERE nombreConcepto = p_nombreConcepto AND correoFamilia = p_correoFamilia AND delete_at IS NULL;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_get_estado_activo_concepto
-- Descripción: Obtiene el estado activo de un concepto para un usuario específico
-- Tipo de operación: Lectura
-- Tablas involucradas: personalizacionconcepto
-- Dependencias: Ninguna

DROP PROCEDURE IF EXISTS sp_get_estado_activo_concepto$$

CREATE PROCEDURE sp_get_estado_activo_concepto(
    IN p_nombreConcepto VARCHAR(255),   -- Nombre del concepto
    IN p_correoFamilia VARCHAR(255),    -- Correo de la familia
    IN p_nombreUsuario VARCHAR(255)     -- Usuario a consultar
)
BEGIN

    -- CONSULTA DE ESTADO ACTIVO

    -- Retorna el estado activo o TRUE por defecto si no existe personalización
    SELECT COALESCE(activo, TRUE) AS activo FROM personalizacionconcepto WHERE nombreConcepto = p_nombreConcepto AND correoFamilia = p_correoFamilia AND nombreUsuario = p_nombreUsuario LIMIT 1;
       -- RETORNO ALTERNATIVO

    -- Si no hay resultado, retornar TRUE por defecto
    SELECT TRUE AS activo WHERE NOT EXISTS (SELECT 1 FROM personalizacionconcepto WHERE nombreConcepto = p_nombreConcepto AND correoFamilia = p_correoFamilia AND nombreUsuario = p_nombreUsuario);
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_actualizar_nombre_personalizaciones
-- Descripción: Actualiza el nombre de concepto en todas las personalizaciones
-- Tipo de operación: Escritura
-- Tablas involucradas: personalizacionconcepto
-- Dependencias: Se usa cuando se renombra un concepto

DROP PROCEDURE IF EXISTS sp_actualizar_nombre_personalizaciones$$

CREATE PROCEDURE sp_actualizar_nombre_personalizaciones(
    IN p_nombreViejo VARCHAR(255),  -- Nombre actual del concepto
    IN p_nombreNuevo VARCHAR(255),  -- Nuevo nombre del concepto
    IN p_correoFamilia VARCHAR(255) -- Correo de la familia
)
BEGIN

    -- ACTUALIZACIÓN EN CASCADA

    -- Actualiza el nombre del concepto en todas sus personalizaciones
    UPDATE personalizacionconcepto SET nombreConcepto = p_nombreNuevo WHERE nombreConcepto = p_nombreViejo AND correoFamilia = p_correoFamilia;
    

   -- Retorno:
    -- rows_affected: Número de personalizaciones actualizadas
    SELECT ROW_COUNT() AS rows_affected;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_get_personalizaciones_by_usuario
-- Descripción: Obtiene todas las personalizaciones de un usuario específico
-- Tipo de operación: Lectura
-- Tablas involucradas: personalizacionconcepto
-- Dependencias: Ninguna

DROP PROCEDURE IF EXISTS sp_get_personalizaciones_by_usuario$$

CREATE PROCEDURE sp_get_personalizaciones_by_usuario(
    IN p_nombre_usuario VARCHAR(255)  -- Nombre del usuario
)
BEGIN

    -- CONSULTA DE PERSONALIZACIONES POR USUARIO

    -- Retorna todas las personalizaciones activas del usuario ordenadas alfabéticamente
    

   -- Retorno: Lista completa de personalizaciones del usuario
    SELECT p.idPersonalizacion, p.limiteGasto, p.activo, p.montoPlanificado, p.tipoPeriodoPlanificado, p.tipoPeriodoLimite, p.diaPeriodoPlanificado, p.notificacion, p.nombreUsuario, p.nombreConcepto, p.correoFamilia, p.delete_at FROM personalizacionconcepto p WHERE p.nombreUsuario = p_nombre_usuario AND p.delete_at IS NULL ORDER BY p.nombreConcepto;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_get_personalizacion_by_id
-- Descripción: Obtiene una personalización específica por su ID
-- Tipo de operación: Lectura
-- Tablas involucradas: personalizacionconcepto
-- Dependencias: Ninguna

DROP PROCEDURE IF EXISTS sp_get_personalizacion_by_id$$

CREATE PROCEDURE sp_get_personalizacion_by_id(
    IN p_id_personalizacion INT  -- ID de la personalización a buscar
)
BEGIN

    -- BÚSQUEDA POR ID

    -- Retorna la personalización si existe y está activa
    

   -- Retorno: Datos completos de la personalización o conjunto vacío
    SELECT idPersonalizacion, limiteGasto, activo, montoPlanificado, tipoPeriodoPlanificado, tipoPeriodoLimite, diaPeriodoPlanificado, notificacion, nombreUsuario, nombreConcepto, correoFamilia, delete_at FROM personalizacionconcepto WHERE idPersonalizacion = p_id_personalizacion AND delete_at IS NULL;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_update_personalizacion
-- Descripción: Actualiza los parámetros principales de una personalización existente
-- Tipo de operación: Escritura
-- Tablas involucradas: personalizacionconcepto
-- Dependencias: Ninguna

DROP PROCEDURE IF EXISTS sp_update_personalizacion$$

CREATE PROCEDURE sp_update_personalizacion(
    IN p_id_personalizacion INT,                 -- ID de la personalización a actualizar
    IN p_limite_gasto DECIMAL(10,2),             -- Nuevo límite de gasto
    IN p_monto_planificado DECIMAL(10,2),        -- Nuevo monto planificado
    IN p_tipo_periodo_planificado VARCHAR(50),   -- Nuevo tipo de período planificado
    IN p_tipo_periodo_limite VARCHAR(50),        -- Nuevo tipo de período límite
    IN p_dia_periodo_planificado TINYINT,        -- Nuevo día del período
    IN p_activo BOOLEAN                          -- Nuevo estado activo
)
BEGIN

    -- ACTUALIZACIÓN DE PARÁMETROS

    -- Actualiza múltiples campos de configuración de la personalización
    UPDATE personalizacionconcepto SET limiteGasto = p_limite_gasto, montoPlanificado = p_monto_planificado, tipoPeriodoPlanificado = p_tipo_periodo_planificado, tipoPeriodoLimite = p_tipo_periodo_limite, diaPeriodoPlanificado = p_dia_periodo_planificado, activo = p_activo WHERE idPersonalizacion = p_id_personalizacion AND delete_at IS NULL;
    

   -- Retorno:
    -- filasAfectadas: Número de filas actualizadas (0 o 1)
    SELECT ROW_COUNT() as filasAfectadas;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_delete_personalizacion
-- Descripción: Realiza eliminación lógica de una personalización
-- Tipo de operación: Escritura
-- Tablas involucradas: personalizacionconcepto
-- Dependencias: Ninguna

DROP PROCEDURE IF EXISTS sp_delete_personalizacion$$

CREATE PROCEDURE sp_delete_personalizacion(
    IN p_id_personalizacion INT  -- ID de la personalización a eliminar
)
BEGIN

    -- ELIMINACIÓN LÓGICA

    -- Marca la personalización con fecha de eliminación
    UPDATE personalizacionconcepto SET delete_at = NOW() WHERE idPersonalizacion = p_id_personalizacion AND delete_at IS NULL;
    

   -- Retorno:
    -- filasAfectadas: Número de filas eliminadas lógicamente (0 o 1)
    SELECT ROW_COUNT() as filasAfectadas;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_find_personalizacion_by_usuario_and_concepto
-- Descripción: Busca la personalización de un usuario para un concepto específico
-- Tipo de operación: Lectura
-- Tablas involucradas: personalizacionconcepto
-- Dependencias: Ninguna

DROP PROCEDURE IF EXISTS sp_find_personalizacion_by_usuario_and_concepto$$

CREATE PROCEDURE sp_find_personalizacion_by_usuario_and_concepto(
    IN p_nombre_usuario VARCHAR(255),   -- Nombre del usuario
    IN p_nombre_concepto VARCHAR(255),  -- Nombre del concepto
    IN p_correo_familia VARCHAR(255)    -- Correo de la familia
)
BEGIN

    -- BÚSQUEDA ESPECÍFICA

    -- Retorna la personalización única del usuario para el concepto
    

   -- Retorno: Datos completos de la personalización o conjunto vacío
    SELECT idPersonalizacion, montoPlanificado, tipoPeriodoPlanificado, diaPeriodoPlanificado, limiteGasto, tipoPeriodoLimite, diaPeriodoLimite, notificacion, activo, nombreUsuario, nombreConcepto, correoFamilia FROM personalizacionconcepto WHERE nombreUsuario = p_nombre_usuario AND nombreConcepto = p_nombre_concepto AND correoFamilia = p_correo_familia AND delete_at IS NULL LIMIT 1;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_find_personalizaciones_by_periodo_diario
-- Descripción: Obtiene personalizaciones activas con período diario planificado
-- Tipo de operación: Lectura
-- Tablas involucradas: personalizacionconcepto
-- Dependencias: Usado para procesamiento automático diario

DROP PROCEDURE IF EXISTS sp_find_personalizaciones_by_periodo_diario$$

CREATE PROCEDURE sp_find_personalizaciones_by_periodo_diario()
BEGIN

    -- CONSULTA DE PERSONALIZACIONES DIARIAS

    -- Filtra personalizaciones activas con monto planificado y tipo diario
    

   -- Retorno: Lista de personalizaciones para procesamiento diario automático
    SELECT p.idPersonalizacion, p.montoPlanificado, p.tipoPeriodoPlanificado, p.diaPeriodoPlanificado, p.limiteGasto, p.tipoPeriodoLimite, p.diaPeriodoLimite, p.notificacion, p.activo, p.nombreUsuario, p.nombreConcepto, p.correoFamilia FROM personalizacionconcepto p WHERE p.activo = TRUE AND p.montoPlanificado IS NOT NULL AND p.montoPlanificado > 0 AND p.tipoPeriodoPlanificado = 'diario' AND p.delete_at IS NULL ORDER BY p.nombreUsuario, p.nombreConcepto;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_find_personalizaciones_by_dia_planificado
-- Descripción: Obtiene personalizaciones que deben ejecutarse en un día específico
-- Tipo de operación: Lectura
-- Tablas involucradas: personalizacionconcepto
-- Dependencias: Usado para procesamiento automático mensual/quincenal

DROP PROCEDURE IF EXISTS sp_find_personalizaciones_by_dia_planificado$$

CREATE PROCEDURE sp_find_personalizaciones_by_dia_planificado(
    IN p_dia TINYINT  -- Día del mes a consultar
)
BEGIN

    -- CONSULTA POR DÍA ESPECÍFICO

    -- Filtra personalizaciones mensuales/quincenales que se ejecutan en el día especificado
    

   -- Retorno: Lista de personalizaciones programadas para el día
    SELECT p.idPersonalizacion, p.montoPlanificado, p.tipoPeriodoPlanificado, p.diaPeriodoPlanificado, p.limiteGasto, p.tipoPeriodoLimite, p.diaPeriodoLimite, p.notificacion, p.activo, p.nombreUsuario, p.nombreConcepto, p.correoFamilia FROM personalizacionconcepto p WHERE p.activo = TRUE AND p.montoPlanificado IS NOT NULL AND p.montoPlanificado > 0 AND p.diaPeriodoPlanificado = p_dia AND p.tipoPeriodoPlanificado IN ('mensual', 'quincenal') AND p.delete_at IS NULL ORDER BY p.nombreUsuario, p.nombreConcepto;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_deshabilitar_personalizacion_para_usuario
-- Descripción: Desactiva o crea una personalización inactiva para un usuario
-- Tipo de operación: Escritura
-- Tablas involucradas: personalizacionconcepto
-- Dependencias: Ninguna

DROP PROCEDURE IF EXISTS sp_deshabilitar_personalizacion_para_usuario$$

CREATE PROCEDURE sp_deshabilitar_personalizacion_para_usuario(
    IN p_nombre_concepto VARCHAR(255),  -- Nombre del concepto
    IN p_correo_familia VARCHAR(255),   -- Correo de la familia
    IN p_nombre_usuario VARCHAR(255)    -- Usuario al que se aplica
)
BEGIN
    DECLARE v_id_personalizacion INT;
    DECLARE v_existe INT DEFAULT 0;
       -- VERIFICACIÓN DE EXISTENCIA

    -- Verificar si existe la personalización (activa o inactiva, pero no eliminada)
    SELECT COUNT(*), COALESCE(MAX(idPersonalizacion), 0) INTO v_existe, v_id_personalizacion FROM personalizacionconcepto WHERE nombreConcepto = p_nombre_concepto AND correoFamilia = p_correo_familia AND nombreUsuario = p_nombre_usuario AND delete_at IS NULL;
       -- CREACIÓN O ACTUALIZACIÓN

    IF v_existe = 0 THEN
        -- No existe, crear una nueva con activo = FALSE
        INSERT INTO personalizacionconcepto (limiteGasto, activo, montoPlanificado, tipoPeriodoPlanificado, tipoPeriodoLimite, diaPeriodoPlanificado, notificacion, nombreUsuario, nombreConcepto, correoFamilia) VALUES (NULL, FALSE, NULL, NULL, NULL, NULL, FALSE, p_nombre_usuario, p_nombre_concepto, p_correo_familia);
        
       
        -- Retorno:
        -- idPersonalizacion: ID de la personalización creada
        -- operacion: 'created'
        SELECT LAST_INSERT_ID() as idPersonalizacion, 'created' as operacion;
    ELSE
        -- Ya existe, solo desactivarla
        UPDATE personalizacionconcepto SET activo = FALSE WHERE idPersonalizacion = v_id_personalizacion AND delete_at IS NULL;
        
       
        -- Retorno:
        -- idPersonalizacion: ID de la personalización actualizada
        -- operacion: 'updated'
        SELECT v_id_personalizacion as idPersonalizacion, 'updated' as operacion;
    END IF;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_habilitar_personalizacion_para_usuario
-- Descripción: Activa o crea una personalización activa para un usuario
-- Tipo de operación: Escritura
-- Tablas involucradas: personalizacionconcepto
-- Dependencias: Ninguna

DROP PROCEDURE IF EXISTS sp_habilitar_personalizacion_para_usuario$$

CREATE PROCEDURE sp_habilitar_personalizacion_para_usuario(
    IN p_nombre_concepto VARCHAR(255),  -- Nombre del concepto
    IN p_correo_familia VARCHAR(255),   -- Correo de la familia
    IN p_nombre_usuario VARCHAR(255)    -- Usuario al que se aplica
)
BEGIN
    DECLARE v_id_personalizacion INT;
    DECLARE v_existe INT DEFAULT 0;
       -- VERIFICACIÓN DE EXISTENCIA

    -- Verificar si existe la personalización (activa o inactiva, pero no eliminada)
    SELECT COUNT(*), COALESCE(MAX(idPersonalizacion), 0) INTO v_existe, v_id_personalizacion FROM personalizacionconcepto WHERE nombreConcepto = p_nombre_concepto AND correoFamilia = p_correo_familia AND nombreUsuario = p_nombre_usuario AND delete_at IS NULL;
       -- CREACIÓN O ACTUALIZACIÓN

    IF v_existe = 0 THEN
        -- No existe, crear una nueva con activo = TRUE
        INSERT INTO personalizacionconcepto (limiteGasto, activo, montoPlanificado, tipoPeriodoPlanificado, tipoPeriodoLimite, diaPeriodoPlanificado, notificacion, nombreUsuario, nombreConcepto, correoFamilia) VALUES (NULL, TRUE, NULL, NULL, NULL, NULL, FALSE, p_nombre_usuario, p_nombre_concepto, p_correo_familia);
        
       
        -- Retorno:
        -- idPersonalizacion: ID de la personalización creada
        -- operacion: 'created'
        SELECT LAST_INSERT_ID() as idPersonalizacion, 'created' as operacion;
    ELSE
        -- Ya existe, solo activarla
        UPDATE personalizacionconcepto SET activo = TRUE WHERE idPersonalizacion = v_id_personalizacion AND delete_at IS NULL;
        
       
        -- Retorno:
        -- idPersonalizacion: ID de la personalización actualizada
        -- operacion: 'updated'
        SELECT v_id_personalizacion as idPersonalizacion, 'updated' as operacion;
    END IF;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_toggle_activo_para_usuario
-- Descripción: Alterna el estado activo/inactivo de una personalización de usuario
-- Tipo de operación: Escritura
-- Tablas involucradas: personalizacionconcepto
-- Dependencias: Ninguna

DROP PROCEDURE IF EXISTS sp_toggle_activo_para_usuario$$

CREATE PROCEDURE sp_toggle_activo_para_usuario(
    IN p_nombre_concepto VARCHAR(255),  -- Nombre del concepto
    IN p_correo_familia VARCHAR(255),   -- Correo de la familia
    IN p_nombre_usuario VARCHAR(255)    -- Usuario al que se aplica
)
BEGIN
    DECLARE v_id_personalizacion INT;
    DECLARE v_existe INT DEFAULT 0;
    DECLARE v_activo_actual BOOLEAN;
    DECLARE v_nuevo_estado BOOLEAN;
       -- VERIFICACIÓN Y OBTENCIÓN DE ESTADO

    -- Verificar si existe y obtener estado actual
    SELECT COUNT(*), COALESCE(MAX(idPersonalizacion), 0), COALESCE(MAX(activo), TRUE) INTO v_existe, v_id_personalizacion, v_activo_actual FROM personalizacionconcepto WHERE nombreConcepto = p_nombre_concepto AND correoFamilia = p_correo_familia AND nombreUsuario = p_nombre_usuario AND delete_at IS NULL;
       -- ALTERNANCIA DE ESTADO

    IF v_existe = 0 THEN
        -- No existe, crear con activo = FALSE (toggle desde "por defecto TRUE")
        INSERT INTO personalizacionconcepto (limiteGasto, activo, montoPlanificado, tipoPeriodoPlanificado, tipoPeriodoLimite, diaPeriodoPlanificado, notificacion, nombreUsuario, nombreConcepto, correoFamilia) VALUES (NULL, FALSE, NULL, NULL, NULL, NULL, FALSE, p_nombre_usuario, p_nombre_concepto, p_correo_familia);
        
       
        -- Retorno:
        -- idPersonalizacion: ID de la personalización creada
        -- operacion: 'created'
        -- nuevoEstado: Estado resultante (FALSE)
        -- estadoAnterior: Estado previo (TRUE por defecto)
        SELECT LAST_INSERT_ID() as idPersonalizacion, 'created' as operacion, FALSE as nuevoEstado, TRUE as estadoAnterior;
    ELSE
        -- Ya existe, invertir estado actual
        SET v_nuevo_estado = NOT v_activo_actual;
        
        UPDATE personalizacionconcepto SET activo = v_nuevo_estado WHERE idPersonalizacion = v_id_personalizacion AND delete_at IS NULL;
        
    
        -- RETORNO - ALTERNANCIA
    
        -- Retorno:
        -- idPersonalizacion: ID de la personalización actualizada
        -- operacion: 'toggled'
        -- nuevoEstado: Estado resultante
        -- estadoAnterior: Estado previo
        SELECT v_id_personalizacion as idPersonalizacion, 'toggled' as operacion, v_nuevo_estado as nuevoEstado, v_activo_actual as estadoAnterior;
    END IF;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_get_usuario_que_asigno
-- Descripción: Obtiene el nombre personal del usuario que creó una personalización
-- Tipo de operación: Lectura
-- Tablas involucradas: personalizacionconcepto, usuario
-- Dependencias: JOIN con tabla usuario

DROP PROCEDURE IF EXISTS sp_get_usuario_que_asigno$$

CREATE PROCEDURE sp_get_usuario_que_asigno(
    IN p_id_personalizacion INT  -- ID de la personalización
)
BEGIN

    -- CONSULTA CON JOIN

    -- Obtiene el nombre personal del usuario desde la tabla usuario
    -- Si no se encuentra, retorna 'Desconocido'
    

   -- Retorno:
    -- nombrePersonal: Nombre real del usuario que asignó la personalización
    SELECT COALESCE(u.nombrePersonal, 'Desconocido') as nombrePersonal FROM personalizacionconcepto p INNER JOIN usuario u ON p.nombreUsuario = u.nombreUsuario AND p.correoFamilia = u.correoFamilia WHERE p.idPersonalizacion = p_id_personalizacion AND p.delete_at IS NULL LIMIT 1;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_toggle_activo_concepto
-- Descripción: Alterna el estado activo de un concepto para un usuario específico
-- Tipo de operación: Escritura
-- Tablas involucradas: personalizacionconcepto
-- Dependencias: Ninguna

DROP PROCEDURE IF EXISTS sp_toggle_activo_concepto$$

CREATE PROCEDURE sp_toggle_activo_concepto(
    IN p_nombreConcepto VARCHAR(255),  -- Nombre del concepto
    IN p_correoFamilia VARCHAR(255),   -- Correo de la familia
    IN p_nombreUsuario VARCHAR(255)    -- Usuario al que se aplica
)
BEGIN
    DECLARE v_existe INT DEFAULT 0;
    DECLARE v_estado_anterior BOOLEAN;
       -- VERIFICACIÓN DE EXISTENCIA Y ESTADO

    -- Verificar si existe personalización y obtener estado actual
    SELECT COUNT(*), COALESCE(MAX(activo), TRUE) INTO v_existe, v_estado_anterior FROM personalizacionconcepto WHERE nombreConcepto = p_nombreConcepto AND correoFamilia = p_correoFamilia AND nombreUsuario = p_nombreUsuario;
       -- CREACIÓN O ALTERNANCIA

    IF v_existe = 0 THEN
        -- No existe personalización, crear una con activo = FALSE
        INSERT INTO personalizacionconcepto (limiteGasto, activo, montoPlanificado, tipoPeriodoPlanificado, tipoPeriodoLimite, diaPeriodoPlanificado, diaPeriodoLimite, notificacion, nombreUsuario, nombreConcepto, correoFamilia) VALUES (NULL, FALSE, NULL, NULL, NULL, NULL, NULL, FALSE, p_nombreUsuario, p_nombreConcepto, p_correoFamilia);
        
       
        -- Retorno:
        -- accion: 'created'
        -- nuevo_estado: FALSE
        -- rows_affected: 1
        SELECT 'created' AS accion, FALSE AS nuevo_estado, 1 AS rows_affected;
    ELSE
        -- Ya existe personalización, cambiar el estado activo
        UPDATE personalizacionconcepto SET activo = NOT activo WHERE nombreConcepto = p_nombreConcepto AND correoFamilia = p_correoFamilia AND nombreUsuario = p_nombreUsuario;
        
    
        -- RETORNO - ALTERNANCIA
    
        -- Retorno:
        -- accion: 'toggled'
        -- nuevo_estado: Estado invertido
        -- rows_affected: Número de filas afectadas
        SELECT 'toggled' AS accion, NOT v_estado_anterior AS nuevo_estado, ROW_COUNT() AS rows_affected;
    END IF;
END$$

DELIMITER ;