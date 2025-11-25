SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;
SET collation_connection = 'utf8mb4_unicode_ci';

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_create_concepto$$

CREATE PROCEDURE sp_create_concepto(
    IN p_nombreConcepto VARCHAR(255),
    IN p_correoFamilia VARCHAR(255),
    IN p_tipo VARCHAR(50),
    IN p_icono VARCHAR(100),
    IN p_color VARCHAR(50),
    IN p_nombreUsuario VARCHAR(255)
)
BEGIN
    INSERT INTO concepto (nombreConcepto, correoFamilia, tipo, icono, color, nombreUsuario)
    VALUES (p_nombreConcepto, p_correoFamilia, p_tipo, p_icono, p_color, p_nombreUsuario);
    
    SELECT ROW_COUNT() AS rows_affected, LAST_INSERT_ID() AS id_concepto;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_find_conceptos_by_familia$$

CREATE PROCEDURE sp_find_conceptos_by_familia(
    IN p_correoFamilia VARCHAR(255),
    IN p_tipo TINYINT
)
BEGIN
    SELECT 
        c.nombreConcepto, 
        c.correoFamilia, 
        c.tipo, 
        c.icono, 
        c.color, 
        c.nombreUsuario,
        COALESCE(pc.activo, 1) AS activo, 
        c.delete_at
    FROM concepto c
    LEFT JOIN personalizacionconcepto pc 
        ON c.nombreConcepto = pc.nombreConcepto 
        AND c.correoFamilia = pc.correoFamilia 
        AND pc.nombreUsuario = c.nombreUsuario
    WHERE c.correoFamilia = p_correoFamilia 
        AND c.tipo = p_tipo 
        AND c.delete_at IS NULL
    ORDER BY c.nombreConcepto;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_update_concepto_icono_color$$

CREATE PROCEDURE sp_update_concepto_icono_color(
    IN p_icono VARCHAR(100),
    IN p_color VARCHAR(50),
    IN p_nombreConcepto VARCHAR(255),
    IN p_correoFamilia VARCHAR(255)
)
BEGIN
    UPDATE concepto 
    SET icono = p_icono, 
        color = p_color
    WHERE nombreConcepto = p_nombreConcepto 
        AND correoFamilia = p_correoFamilia 
        AND delete_at IS NULL;
    
    SELECT ROW_COUNT() AS rows_affected;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_find_all_conceptos_by_familia$$

CREATE PROCEDURE sp_find_all_conceptos_by_familia(
    IN p_correoFamilia VARCHAR(255)
)
BEGIN
    SELECT 
        nombreConcepto, 
        correoFamilia, 
        tipo, 
        icono, 
        color, 
        nombreUsuario, 
        delete_at
    FROM concepto 
    WHERE correoFamilia = p_correoFamilia 
        AND delete_at IS NULL
    ORDER BY nombreConcepto;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_find_concepto_by_nombre$$

CREATE PROCEDURE sp_find_concepto_by_nombre(
    IN p_nombreConcepto VARCHAR(255),
    IN p_correoFamilia VARCHAR(255)
)
BEGIN
    SELECT 
        nombreConcepto, 
        correoFamilia, 
        tipo, 
        icono, 
        color, 
        nombreUsuario, 
        delete_at
    FROM concepto 
    WHERE nombreConcepto = p_nombreConcepto 
        AND correoFamilia = p_correoFamilia 
        AND delete_at IS NULL;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_find_concepto_by_nombre$$

CREATE PROCEDURE sp_find_concepto_by_nombre(
    IN p_nombreConcepto VARCHAR(255),
    IN p_correoFamilia VARCHAR(255)
)
BEGIN
    SELECT 
        nombreConcepto, 
        correoFamilia, 
        tipo, 
        icono, 
        color, 
        nombreUsuario, 
        delete_at
    FROM concepto 
    WHERE nombreConcepto = p_nombreConcepto 
        AND correoFamilia = p_correoFamilia 
        AND delete_at IS NULL;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_concepto_exists$$

CREATE PROCEDURE sp_concepto_exists(
    IN p_nombreConcepto VARCHAR(255),
    IN p_correoFamilia VARCHAR(255)
)
BEGIN
    SELECT COUNT(*) AS existe
    FROM concepto 
    WHERE nombreConcepto = p_nombreConcepto 
        AND correoFamilia = p_correoFamilia 
        AND delete_at IS NULL;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_create_personalizacion$$

CREATE PROCEDURE sp_create_personalizacion(
    IN p_limiteGasto DECIMAL(10,2),
    IN p_activo BOOLEAN,
    IN p_montoPlanificado DECIMAL(10,2),
    IN p_tipoPeriodoPlanificado VARCHAR(50),
    IN p_tipoPeriodoLimite VARCHAR(50),
    IN p_diaPeriodoPlanificado TINYINT,
    IN p_diaPeriodoLimite TINYINT,
    IN p_notificacion BOOLEAN,
    IN p_nombreUsuario VARCHAR(255),
    IN p_nombreConcepto VARCHAR(255),
    IN p_correoFamilia VARCHAR(255)
)
BEGIN
    INSERT INTO personalizacionconcepto 
        (limiteGasto, activo, montoPlanificado, tipoPeriodoPlanificado, 
         tipoPeriodoLimite, diaPeriodoPlanificado, diaPeriodoLimite, notificacion,
         nombreUsuario, nombreConcepto, correoFamilia)
    VALUES 
        (p_limiteGasto, p_activo, p_montoPlanificado, p_tipoPeriodoPlanificado, 
         p_tipoPeriodoLimite, p_diaPeriodoPlanificado, p_diaPeriodoLimite, p_notificacion,
         p_nombreUsuario, p_nombreConcepto, p_correoFamilia);
    
    SELECT ROW_COUNT() AS rows_affected, LAST_INSERT_ID() AS id_personalizacion;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_get_usuarios_by_familia$$

CREATE PROCEDURE sp_get_usuarios_by_familia(
    IN p_correoFamilia VARCHAR(255)
)
BEGIN
    SELECT 
        nombreUsuario, 
        rol, 
        contrasenaPersonal, 
        nombrePersonal, 
        correoFamilia, 
        delete_at
    FROM usuario 
    WHERE correoFamilia = p_correoFamilia 
        AND delete_at IS NULL;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_create_personalizaciones_all_usuarios$$

CREATE PROCEDURE sp_create_personalizaciones_all_usuarios(
    IN p_nombreConcepto VARCHAR(255),
    IN p_correoFamilia VARCHAR(255),
    IN p_limiteGasto DECIMAL(10,2),
    IN p_montoPlanificado DECIMAL(10,2),
    IN p_tipoPeriodoPlanificado VARCHAR(50),
    IN p_tipoPeriodoLimite VARCHAR(50),
    IN p_diaPeriodoPlanificado TINYINT,
    IN p_diaPeriodoLimite TINYINT
)
BEGIN
    DECLARE v_usuarios_total INT DEFAULT 0;
    DECLARE v_usuarios_creados INT DEFAULT 0;
    DECLARE v_usuarios_error INT DEFAULT 0;
    DECLARE v_finished INTEGER DEFAULT 0;
    DECLARE v_nombreUsuario VARCHAR(255);
    
    DECLARE cur_usuarios CURSOR FOR 
        SELECT nombreUsuario 
        FROM usuario 
        WHERE correoFamilia = p_correoFamilia 
            AND delete_at IS NULL;
    
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET v_finished = 1;
    DECLARE CONTINUE HANDLER FOR SQLEXCEPTION SET v_usuarios_error = v_usuarios_error + 1;
    
    SELECT COUNT(*) INTO v_usuarios_total
    FROM usuario 
    WHERE correoFamilia = p_correoFamilia 
        AND delete_at IS NULL;
    
    OPEN cur_usuarios;
    
    usuario_loop: LOOP
        FETCH cur_usuarios INTO v_nombreUsuario;
        
        IF v_finished = 1 THEN 
            LEAVE usuario_loop;
        END IF;
        
        INSERT INTO personalizacionconcepto 
            (limiteGasto, activo, montoPlanificado, tipoPeriodoPlanificado, 
             tipoPeriodoLimite, diaPeriodoPlanificado, diaPeriodoLimite, notificacion,
             nombreUsuario, nombreConcepto, correoFamilia)
        VALUES 
            (p_limiteGasto, TRUE, p_montoPlanificado, p_tipoPeriodoPlanificado, 
             p_tipoPeriodoLimite, p_diaPeriodoPlanificado, p_diaPeriodoLimite, FALSE,
             v_nombreUsuario, p_nombreConcepto, p_correoFamilia);
        
        SET v_usuarios_creados = v_usuarios_creados + 1;
        
    END LOOP usuario_loop;
    
    CLOSE cur_usuarios;
    
    SELECT 
        v_usuarios_total AS usuarios_total,
        v_usuarios_creados AS usuarios_creados,
        v_usuarios_error AS usuarios_error;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_get_personalizaciones_by_concepto$$

CREATE PROCEDURE sp_get_personalizaciones_by_concepto(
    IN p_nombreConcepto VARCHAR(255),
    IN p_correoFamilia VARCHAR(255)
)
BEGIN
    SELECT 
        idPersonalizacion, 
        nombreUsuario, 
        limiteGasto, 
        montoPlanificado, 
        tipoPeriodoPlanificado, 
        tipoPeriodoLimite
    FROM personalizacionconcepto 
    WHERE nombreConcepto = p_nombreConcepto 
        AND correoFamilia = p_correoFamilia 
        AND delete_at IS NULL;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_get_count_by_tipo$$

CREATE PROCEDURE sp_get_count_by_tipo(
    IN p_correoFamilia VARCHAR(255)
)
BEGIN
    SELECT 
        SUM(CASE WHEN tipo = 0 THEN 1 ELSE 0 END) AS gastos,
        SUM(CASE WHEN tipo = 1 THEN 1 ELSE 0 END) AS ingresos
    FROM concepto 
    WHERE correoFamilia = p_correoFamilia 
        AND delete_at IS NULL;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_find_conceptos_by_familia_activos$$

CREATE PROCEDURE sp_find_conceptos_by_familia_activos(
    IN p_correoFamilia VARCHAR(255),
    IN p_nombreUsuario VARCHAR(255),
    IN p_tipo TINYINT
)
BEGIN
    SELECT 
        c.nombreConcepto, 
        c.correoFamilia, 
        c.tipo, 
        c.icono, 
        c.color, 
        c.nombreUsuario, 
        COALESCE(pc.activo, TRUE) AS activo, 
        c.delete_at
    FROM concepto c
    LEFT JOIN personalizacionconcepto pc 
        ON c.nombreConcepto = pc.nombreConcepto 
        AND c.correoFamilia = pc.correoFamilia 
        AND pc.nombreUsuario = p_nombreUsuario
    WHERE c.correoFamilia = p_correoFamilia 
        AND c.tipo = p_tipo 
        AND c.delete_at IS NULL
    ORDER BY c.nombreConcepto;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_update_concepto$$

CREATE PROCEDURE sp_update_concepto(
    IN p_icono VARCHAR(100),
    IN p_color VARCHAR(50),
    IN p_nombreUsuario VARCHAR(255),
    IN p_nombreConcepto VARCHAR(255),
    IN p_correoFamilia VARCHAR(255)
)
BEGIN
    UPDATE concepto 
    SET icono = p_icono, 
        color = p_color, 
        nombreUsuario = p_nombreUsuario
    WHERE nombreConcepto = p_nombreConcepto 
        AND correoFamilia = p_correoFamilia 
        AND nombreUsuario = p_nombreUsuario;
    
    SELECT ROW_COUNT() AS rows_affected;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_toggle_activo_concepto$$

CREATE PROCEDURE sp_toggle_activo_concepto(
    IN p_nombreConcepto VARCHAR(255),
    IN p_correoFamilia VARCHAR(255),
    IN p_nombreUsuario VARCHAR(255)
)
BEGIN
    DECLARE v_existe INT DEFAULT 0;
    DECLARE v_estado_anterior BOOLEAN;
    
    -- Verificar si existe personalización
    SELECT COUNT(*), COALESCE(MAX(activo), TRUE) 
    INTO v_existe, v_estado_anterior
    FROM personalizacionconcepto 
    WHERE nombreConcepto = p_nombreConcepto 
        AND correoFamilia = p_correoFamilia 
        AND nombreUsuario = p_nombreUsuario;
    
    IF v_existe = 0 THEN
        -- No existe personalización, crear una con activo = false
        INSERT INTO personalizacionconcepto 
            (limiteGasto, activo, montoPlanificado, tipoPeriodoPlanificado, 
             tipoPeriodoLimite, diaPeriodoPlanificado, diaPeriodoLimite, notificacion,
             nombreUsuario, nombreConcepto, correoFamilia)
        VALUES 
            (NULL, FALSE, NULL, NULL, NULL, NULL, NULL, FALSE, 
             p_nombreUsuario, p_nombreConcepto, p_correoFamilia);
        
        SELECT 'created' AS accion, FALSE AS nuevo_estado, 1 AS rows_affected;
    ELSE
        -- Ya existe personalización, cambiar el estado activo
        UPDATE personalizacionconcepto 
        SET activo = NOT activo 
        WHERE nombreConcepto = p_nombreConcepto 
            AND correoFamilia = p_correoFamilia 
            AND nombreUsuario = p_nombreUsuario;
        
        SELECT 'toggled' AS accion, NOT v_estado_anterior AS nuevo_estado, ROW_COUNT() AS rows_affected;
    END IF;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_get_estado_activo_concepto$$

CREATE PROCEDURE sp_get_estado_activo_concepto(
    IN p_nombreConcepto VARCHAR(255),
    IN p_correoFamilia VARCHAR(255),
    IN p_nombreUsuario VARCHAR(255)
)
BEGIN
    SELECT COALESCE(activo, TRUE) AS activo
    FROM personalizacionconcepto 
    WHERE nombreConcepto = p_nombreConcepto 
        AND correoFamilia = p_correoFamilia 
        AND nombreUsuario = p_nombreUsuario
    LIMIT 1;
    
    -- Si no hay resultado, retornar TRUE por defecto
    SELECT TRUE AS activo 
    WHERE NOT EXISTS (
        SELECT 1 FROM personalizacionconcepto 
        WHERE nombreConcepto = p_nombreConcepto 
            AND correoFamilia = p_correoFamilia 
            AND nombreUsuario = p_nombreUsuario
    );
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_actualizar_nombre_personalizaciones$$

CREATE PROCEDURE sp_actualizar_nombre_personalizaciones(
    IN p_nombreViejo VARCHAR(255),
    IN p_nombreNuevo VARCHAR(255),
    IN p_correoFamilia VARCHAR(255)
)
BEGIN
    UPDATE personalizacionconcepto 
    SET nombreConcepto = p_nombreNuevo
    WHERE nombreConcepto = p_nombreViejo 
        AND correoFamilia = p_correoFamilia;
    
    SELECT ROW_COUNT() AS rows_affected;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_deshabilitar_concepto$$

CREATE PROCEDURE sp_deshabilitar_concepto(
    IN p_nombreConcepto VARCHAR(255),
    IN p_correoFamilia VARCHAR(255)
)
BEGIN
    DECLARE v_rows_concepto INT DEFAULT 0;
    DECLARE v_rows_personalizaciones INT DEFAULT 0;
    
    -- Iniciar transacción
    START TRANSACTION;
    
    -- 1. Desactivar el concepto en la tabla concepto
    UPDATE concepto 
    SET activo = 0 
    WHERE nombreConcepto = p_nombreConcepto 
        AND correoFamilia = p_correoFamilia;
    
    SET v_rows_concepto = ROW_COUNT();
    
    -- 2. Desactivar todas las personalizaciones del concepto
    UPDATE personalizacionconcepto 
    SET activo = 0 
    WHERE nombreConcepto = p_nombreConcepto 
        AND correoFamilia = p_correoFamilia;
    
    SET v_rows_personalizaciones = ROW_COUNT();
    
    -- Confirmar transacción
    COMMIT;
    
    -- Retornar estadísticas
    SELECT 
        v_rows_concepto AS concepto_rows,
        v_rows_personalizaciones AS personalizaciones_rows,
        'deshabilitado' AS estado;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_habilitar_concepto$$

CREATE PROCEDURE sp_habilitar_concepto(
    IN p_nombreConcepto VARCHAR(255),
    IN p_correoFamilia VARCHAR(255)
)
BEGIN
    DECLARE v_rows_concepto INT DEFAULT 0;
    DECLARE v_rows_personalizaciones INT DEFAULT 0;
    
    -- Iniciar transacción
    START TRANSACTION;
    
    -- 1. Reactivar el concepto
    UPDATE concepto 
    SET activo = 1 
    WHERE nombreConcepto = p_nombreConcepto 
        AND correoFamilia = p_correoFamilia;
    
    SET v_rows_concepto = ROW_COUNT();
    
    -- 2. Reactivar todas las personalizaciones
    UPDATE personalizacionconcepto 
    SET activo = 1 
    WHERE nombreConcepto = p_nombreConcepto 
        AND correoFamilia = p_correoFamilia;
    
    SET v_rows_personalizaciones = ROW_COUNT();
    
    -- Confirmar transacción
    COMMIT;
    
    -- Retornar estadísticas
    SELECT 
        v_rows_concepto AS concepto_rows,
        v_rows_personalizaciones AS personalizaciones_rows,
        'habilitado' AS estado;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_toggle_activo_global$$

CREATE PROCEDURE sp_toggle_activo_global(
    IN p_nombreConcepto VARCHAR(255),
    IN p_correoFamilia VARCHAR(255)
)
BEGIN
    DECLARE v_activo_actual BOOLEAN;
    DECLARE v_rows_concepto INT DEFAULT 0;
    DECLARE v_rows_personalizaciones INT DEFAULT 0;
    
    -- Obtener el estado actual
    SELECT activo INTO v_activo_actual
    FROM concepto 
    WHERE nombreConcepto = p_nombreConcepto 
        AND correoFamilia = p_correoFamilia
    LIMIT 1;
    
    -- Iniciar transacción
    START TRANSACTION;
    
    IF v_activo_actual = TRUE THEN
        -- Deshabilitar
        UPDATE concepto 
        SET activo = 0 
        WHERE nombreConcepto = p_nombreConcepto 
            AND correoFamilia = p_correoFamilia;
        
        SET v_rows_concepto = ROW_COUNT();
        
        UPDATE personalizacionconcepto 
        SET activo = 0 
        WHERE nombreConcepto = p_nombreConcepto 
            AND correoFamilia = p_correoFamilia;
        
        SET v_rows_personalizaciones = ROW_COUNT();
        
        SELECT 
            v_rows_concepto AS concepto_rows,
            v_rows_personalizaciones AS personalizaciones_rows,
            'deshabilitado' AS accion,
            FALSE AS nuevo_estado;
    ELSE
        -- Habilitar
        UPDATE concepto 
        SET activo = 1 
        WHERE nombreConcepto = p_nombreConcepto 
            AND correoFamilia = p_correoFamilia;
        
        SET v_rows_concepto = ROW_COUNT();
        
        UPDATE personalizacionconcepto 
        SET activo = 1 
        WHERE nombreConcepto = p_nombreConcepto 
            AND correoFamilia = p_correoFamilia;
        
        SET v_rows_personalizaciones = ROW_COUNT();
        
        SELECT 
            v_rows_concepto AS concepto_rows,
            v_rows_personalizaciones AS personalizaciones_rows,
            'habilitado' AS accion,
            TRUE AS nuevo_estado;
    END IF;
    
    -- Confirmar transacción
    COMMIT;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_familia_create$$

CREATE PROCEDURE sp_familia_create(
    IN p_correo VARCHAR(255),
    IN p_telefono VARCHAR(20),
    IN p_contraseña VARCHAR(255),
    OUT p_id_familia INT
)
BEGIN
    INSERT INTO familia (correo, telefono, contraseña)
    VALUES (p_correo, p_telefono, p_contraseña);
    
    SET p_id_familia = LAST_INSERT_ID();
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_familia_find_by_email$$

CREATE PROCEDURE sp_familia_find_by_email(
    IN p_correo VARCHAR(255)
)
BEGIN
    SELECT correo, telefono, contraseña, delete_at
    FROM familia 
    WHERE correo = p_correo 
      AND delete_at IS NULL;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_familia_get_all$$

CREATE PROCEDURE sp_familia_get_all()
BEGIN
    SELECT id_familia, correo, telefono, contraseña, delete_at
    FROM familia 
    WHERE delete_at IS NULL
    ORDER BY id_familia;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_get_password_familia$$

CREATE PROCEDURE sp_get_password_familia(
    IN p_correo VARCHAR(255)
)
BEGIN
    SELECT contraseña
    FROM familia 
    WHERE correo = p_correo 
        AND delete_at IS NULL
    LIMIT 1;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_delete_familia$$

CREATE PROCEDURE sp_delete_familia(
    IN p_correo VARCHAR(255)
)
BEGIN
    DELETE FROM familia 
    WHERE correo = p_correo;
    
    SELECT ROW_COUNT() AS rows_affected;
END$$

DELIMITER ;