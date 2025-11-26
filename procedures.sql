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

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_consultar_movimientos$$

CREATE PROCEDURE sp_consultar_movimientos(
    IN p_correoFamilia VARCHAR(255),
    IN p_fecha DATE,
    IN p_nombreUsuarios TEXT  -- Lista separada por comas: "user1,user2,user3"
)
BEGIN
    -- Crear tabla temporal para los usuarios
    DROP TEMPORARY TABLE IF EXISTS temp_usuarios;
    CREATE TEMPORARY TABLE temp_usuarios (
        nombreUsuario VARCHAR(255)
    );
    
    -- Insertar usuarios desde la lista separada por comas
    -- Esto requiere que p_nombreUsuarios venga como: "user1,user2,user3"
    SET @sql = CONCAT(
        'INSERT INTO temp_usuarios (nombreUsuario) VALUES ',
        REPLACE(REPLACE(p_nombreUsuarios, ',', '"),("'), ',', '("'),
        '")'
    );
    
    PREPARE stmt FROM @sql;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;
    
    -- Consultar movimientos usando la tabla temporal
    SELECT 
        m.idMovimiento, 
        m.fecha, 
        m.monto, 
        m.descripcion, 
        m.nombreUsuario, 
        m.nombreConcepto, 
        m.correoFamilia, 
        c.tipo
    FROM movimiento m
    INNER JOIN concepto c 
        ON m.nombreConcepto = c.nombreConcepto 
        AND m.correoFamilia = c.correoFamilia
    INNER JOIN temp_usuarios tu
        ON m.nombreUsuario = tu.nombreUsuario
    WHERE m.correoFamilia = p_correoFamilia 
        AND DATE(m.fecha) = p_fecha 
        AND m.delete_at IS NULL
    ORDER BY m.nombreUsuario ASC, m.fecha DESC;
    
    -- Limpiar tabla temporal
    DROP TEMPORARY TABLE IF EXISTS temp_usuarios;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_find_movimientos_by_familia_date_tipo$$

CREATE PROCEDURE sp_find_movimientos_by_familia_date_tipo(
    IN p_correoFamilia VARCHAR(255),
    IN p_fecha DATE,
    IN p_tipo TINYINT
)
BEGIN
    SELECT 
        m.idMovimiento, 
        m.fecha, 
        m.monto, 
        m.descripcion, 
        m.nombreUsuario, 
        m.nombreConcepto, 
        m.correoFamilia
    FROM movimiento m
    INNER JOIN concepto c 
        ON m.nombreConcepto = c.nombreConcepto 
        AND m.correoFamilia = c.correoFamilia
    WHERE m.correoFamilia = p_correoFamilia 
        AND DATE(m.fecha) = p_fecha 
        AND c.tipo = p_tipo
        AND m.delete_at IS NULL
    ORDER BY m.fecha DESC;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_update_movimiento$$

CREATE PROCEDURE sp_update_movimiento(
    IN p_idMovimiento INT,
    IN p_fecha DATETIME,
    IN p_monto DECIMAL(10,2),
    IN p_descripcion TEXT,
    IN p_nombreConcepto VARCHAR(255)
)
BEGIN
    UPDATE movimiento 
    SET fecha = p_fecha, 
        monto = p_monto, 
        descripcion = p_descripcion, 
        nombreConcepto = p_nombreConcepto
    WHERE idMovimiento = p_idMovimiento 
        AND delete_at IS NULL;
    
    SELECT ROW_COUNT() AS rows_affected;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_delete_movimiento$$

CREATE PROCEDURE sp_delete_movimiento(
    IN p_idMovimiento INT
)
BEGIN
    UPDATE movimiento 
    SET delete_at = NOW() 
    WHERE idMovimiento = p_idMovimiento 
        AND delete_at IS NULL;
    
    SELECT ROW_COUNT() AS rows_affected, NOW() AS deleted_at;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_get_totales_by_familia_date$$

CREATE PROCEDURE sp_get_totales_by_familia_date(
    IN p_correoFamilia VARCHAR(255),
    IN p_fecha DATE
)
BEGIN
    SELECT 
        COALESCE(SUM(CASE WHEN c.tipo = 0 THEN m.monto ELSE 0 END), 0) AS totalGastos,
        COALESCE(SUM(CASE WHEN c.tipo = 1 THEN m.monto ELSE 0 END), 0) AS totalIngresos,
        COALESCE(SUM(CASE WHEN c.tipo = 1 THEN m.monto ELSE 0 END), 0) - 
        COALESCE(SUM(CASE WHEN c.tipo = 0 THEN m.monto ELSE 0 END), 0) AS balance,
        COUNT(*) AS total_movimientos
    FROM movimiento m
    INNER JOIN concepto c 
        ON m.nombreConcepto = c.nombreConcepto 
        AND m.correoFamilia = c.correoFamilia
    WHERE m.correoFamilia = p_correoFamilia 
        AND DATE(m.fecha) = p_fecha
        AND m.delete_at IS NULL;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_create_movimiento$$

CREATE PROCEDURE sp_create_movimiento(
    IN p_fecha DATETIME,
    IN p_monto DECIMAL(10,2),
    IN p_descripcion TEXT,
    IN p_nombreUsuario VARCHAR(255),
    IN p_nombreConcepto VARCHAR(255),
    IN p_correoFamilia VARCHAR(255)
)
BEGIN
    INSERT INTO movimiento 
        (fecha, monto, descripcion, nombreUsuario, nombreConcepto, correoFamilia)
    VALUES 
        (p_fecha, p_monto, p_descripcion, p_nombreUsuario, p_nombreConcepto, p_correoFamilia);
    
    SELECT LAST_INSERT_ID() AS id_movimiento, ROW_COUNT() AS rows_affected;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_find_movimiento_by_id$$

CREATE PROCEDURE sp_find_movimiento_by_id(
    IN p_idMovimiento INT
)
BEGIN
    SELECT 
        idMovimiento, 
        fecha, 
        monto, 
        descripcion, 
        nombreUsuario, 
        nombreConcepto, 
        correoFamilia, 
        delete_at
    FROM movimiento 
    WHERE idMovimiento = p_idMovimiento 
        AND delete_at IS NULL;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_find_movimientos_by_familia_date$$

CREATE PROCEDURE sp_find_movimientos_by_familia_date(
    IN p_nombreUsuario VARCHAR(255),
    IN p_fecha DATE
)
BEGIN
    SELECT 
        m.idMovimiento, 
        m.fecha, 
        m.monto, 
        m.descripcion, 
        m.nombreUsuario, 
        m.nombreConcepto, 
        m.correoFamilia, 
        c.tipo
    FROM movimiento m
    INNER JOIN concepto c 
        ON m.nombreConcepto = c.nombreConcepto 
        AND m.correoFamilia = c.correoFamilia
    INNER JOIN personalizacionconcepto pc 
        ON m.nombreConcepto = pc.nombreConcepto 
        AND m.correoFamilia = pc.correoFamilia 
        AND m.nombreUsuario = pc.nombreUsuario
    WHERE m.nombreUsuario = p_nombreUsuario 
        AND DATE(m.fecha) = p_fecha 
        AND m.delete_at IS NULL
        AND pc.activo = 1
        AND pc.delete_at IS NULL
    ORDER BY m.fecha DESC;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_find_movimientos_familia_date_tipo$$

CREATE PROCEDURE sp_find_movimientos_familia_date_tipo(
    IN p_correoFamilia VARCHAR(255),
    IN p_fecha DATE,
    IN p_tipo TINYINT
)
BEGIN
    SELECT 
        m.idMovimiento, 
        m.fecha, 
        m.monto, 
        m.descripcion, 
        m.nombreUsuario, 
        m.nombreConcepto, 
        m.correoFamilia
    FROM movimiento m
    INNER JOIN concepto c 
        ON m.nombreConcepto = c.nombreConcepto 
        AND m.correoFamilia = c.correoFamilia
    WHERE m.correoFamilia = p_correoFamilia 
        AND DATE(m.fecha) = p_fecha 
        AND c.tipo = p_tipo
        AND m.delete_at IS NULL
    ORDER BY m.fecha DESC;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_get_totales_by_familia_and_date$$

CREATE PROCEDURE sp_get_totales_by_familia_and_date(
    IN p_correo_familia VARCHAR(255),
    IN p_fecha DATE
)
BEGIN
    SELECT 
        COALESCE(SUM(CASE WHEN c.tipo = 0 THEN m.monto ELSE 0 END), 0) as totalGastos,
        COALESCE(SUM(CASE WHEN c.tipo = 1 THEN m.monto ELSE 0 END), 0) as totalIngresos
    FROM movimiento m
    INNER JOIN concepto c ON m.nombreConcepto = c.nombreConcepto 
                          AND m.correoFamilia = c.correoFamilia
    WHERE m.correoFamilia = p_correo_familia
      AND DATE(m.fecha) = p_fecha
      AND m.delete_at IS NULL;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_get_totales_by_usuario_and_period$$

CREATE PROCEDURE sp_get_totales_by_usuario_and_period(
    IN p_nombre_usuario VARCHAR(255),
    IN p_correo_familia VARCHAR(255),
    IN p_fecha_inicio DATETIME,
    IN p_fecha_fin DATETIME
)
BEGIN
    SELECT 
        COALESCE(SUM(CASE WHEN c.tipo = 0 THEN m.monto ELSE 0 END), 0) as totalGastos,
        COALESCE(SUM(CASE WHEN c.tipo = 1 THEN m.monto ELSE 0 END), 0) as totalIngresos
    FROM movimiento m
    INNER JOIN concepto c ON m.nombreConcepto = c.nombreConcepto 
                          AND m.correoFamilia = c.correoFamilia
    WHERE m.nombreUsuario = p_nombre_usuario
      AND m.correoFamilia = p_correo_familia
      AND m.fecha >= p_fecha_inicio
      AND m.fecha <= p_fecha_fin
      AND m.delete_at IS NULL;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_find_movimientos_by_usuario_and_date$$

CREATE PROCEDURE sp_find_movimientos_by_usuario_and_date(
    IN p_nombre_usuario VARCHAR(255),
    IN p_correo_familia VARCHAR(255),
    IN p_fecha DATE
)
BEGIN
    SELECT 
        m.idMovimiento,
        m.fecha,
        m.monto,
        m.descripcion,
        m.nombreUsuario,
        m.nombreConcepto,
        m.correoFamilia
    FROM movimiento m
    WHERE m.nombreUsuario = p_nombre_usuario
      AND m.correoFamilia = p_correo_familia
      AND DATE(m.fecha) = p_fecha
      AND m.delete_at IS NULL
    ORDER BY m.fecha DESC;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_find_movimientos_by_usuario_and_period$$

CREATE PROCEDURE sp_find_movimientos_by_usuario_and_period(
    IN p_nombre_usuario VARCHAR(255),
    IN p_correo_familia VARCHAR(255),
    IN p_fecha_inicio DATETIME,
    IN p_fecha_fin DATETIME
)
BEGIN
    SELECT 
        m.idMovimiento,
        m.fecha,
        m.monto,
        m.descripcion,
        m.nombreUsuario,
        m.nombreConcepto,
        m.correoFamilia
    FROM movimiento m
    WHERE m.nombreUsuario = p_nombre_usuario
      AND m.correoFamilia = p_correo_familia
      AND m.fecha >= p_fecha_inicio
      AND m.fecha <= p_fecha_fin
      AND m.delete_at IS NULL
    ORDER BY m.fecha DESC;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_get_total_by_concepto_and_period$$

CREATE PROCEDURE sp_get_total_by_concepto_and_period(
    IN p_nombre_concepto VARCHAR(255),
    IN p_correo_familia VARCHAR(255),
    IN p_fecha_inicio DATETIME,
    IN p_fecha_fin DATETIME,
    IN p_usuario_filtro VARCHAR(255)
)
BEGIN
    SELECT COALESCE(SUM(m.monto), 0) as total
    FROM movimiento m
    WHERE m.nombreConcepto = p_nombre_concepto
      AND m.correoFamilia = p_correo_familia
      AND m.fecha >= p_fecha_inicio
      AND m.fecha <= p_fecha_fin
      AND m.delete_at IS NULL
      AND (p_usuario_filtro IS NULL OR p_usuario_filtro = '' OR m.nombreUsuario = p_usuario_filtro);
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_create_movimiento_automatico$$

CREATE PROCEDURE sp_create_movimiento_automatico(
    IN p_nombre_usuario VARCHAR(255),
    IN p_nombre_concepto VARCHAR(255),
    IN p_correo_familia VARCHAR(255),
    IN p_monto DECIMAL(10,2),
    IN p_descripcion TEXT
)
BEGIN
    INSERT INTO movimiento (
        fecha, 
        monto, 
        descripcion, 
        nombreUsuario, 
        nombreConcepto, 
        correoFamilia
    )
    VALUES (
        CURDATE(), 
        p_monto, 
        p_descripcion, 
        p_nombre_usuario, 
        p_nombre_concepto, 
        p_correo_familia
    );
    
    -- Retorna el ID del movimiento creado
    SELECT LAST_INSERT_ID() as idMovimiento;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_existe_movimiento_hoy$$

CREATE PROCEDURE sp_existe_movimiento_hoy(
    IN p_nombre_usuario VARCHAR(255),
    IN p_nombre_concepto VARCHAR(255),
    IN p_correo_familia VARCHAR(255)
)
BEGIN
    SELECT COUNT(*) > 0 as existe
    FROM movimiento 
    WHERE nombreUsuario = p_nombre_usuario
      AND nombreConcepto = p_nombre_concepto
      AND correoFamilia = p_correo_familia
      AND DATE(fecha) = CURDATE()
      AND delete_at IS NULL;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_find_movimientos_by_usuario_date_and_tipo$$

CREATE PROCEDURE sp_find_movimientos_by_usuario_date_and_tipo(
    IN p_nombre_usuario VARCHAR(255),
    IN p_correo_familia VARCHAR(255),
    IN p_fecha DATE,
    IN p_tipo TINYINT
)
BEGIN
    SELECT 
        m.idMovimiento,
        m.fecha,
        m.monto,
        m.descripcion,
        m.nombreUsuario,
        m.nombreConcepto,
        m.correoFamilia
    FROM movimiento m
    INNER JOIN concepto c ON m.nombreConcepto = c.nombreConcepto 
                          AND m.correoFamilia = c.correoFamilia
    WHERE m.nombreUsuario = p_nombre_usuario
      AND m.correoFamilia = p_correo_familia
      AND DATE(m.fecha) = p_fecha
      AND c.tipo = p_tipo
      AND m.delete_at IS NULL
    ORDER BY m.fecha DESC;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_get_totales_by_usuario_and_date$$

CREATE PROCEDURE sp_get_totales_by_usuario_and_date(
    IN p_nombre_usuario VARCHAR(255),
    IN p_correo_familia VARCHAR(255),
    IN p_fecha DATE
)
BEGIN
    SELECT 
        COALESCE(SUM(CASE WHEN c.tipo = 0 THEN m.monto ELSE 0 END), 0) as totalGastos,
        COALESCE(SUM(CASE WHEN c.tipo = 1 THEN m.monto ELSE 0 END), 0) as totalIngresos
    FROM movimiento m
    INNER JOIN concepto c ON m.nombreConcepto = c.nombreConcepto 
                          AND m.correoFamilia = c.correoFamilia
    WHERE m.nombreUsuario = p_nombre_usuario
      AND m.correoFamilia = p_correo_familia
      AND DATE(m.fecha) = p_fecha
      AND m.delete_at IS NULL;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_get_personalizaciones_by_usuario$$

CREATE PROCEDURE sp_get_personalizaciones_by_usuario(
    IN p_nombre_usuario VARCHAR(255)
)
BEGIN
    SELECT 
        p.idPersonalizacion,
        p.limiteGasto,
        p.activo,
        p.montoPlanificado,
        p.tipoPeriodoPlanificado,
        p.tipoPeriodoLimite,
        p.diaPeriodoPlanificado,
        p.notificacion,
        p.nombreUsuario,
        p.nombreConcepto,
        p.correoFamilia,
        p.delete_at
    FROM personalizacionconcepto p
    WHERE p.nombreUsuario = p_nombre_usuario
      AND p.delete_at IS NULL
    ORDER BY p.nombreConcepto;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_get_consumo_actual$$

CREATE PROCEDURE sp_get_consumo_actual(
    IN p_nombre_usuario VARCHAR(255),
    IN p_nombre_concepto VARCHAR(255),
    IN p_periodo VARCHAR(20)
)
BEGIN
    DECLARE v_fecha_inicio DATETIME;
    DECLARE v_año INT;
    DECLARE v_mes INT;
    DECLARE v_dia INT;
    DECLARE v_dia_semana INT;
    
    -- Obtener fecha y tiempo actual
    SET v_año = YEAR(NOW());
    SET v_mes = MONTH(NOW());
    SET v_dia = DAY(NOW());
    SET v_dia_semana = WEEKDAY(NOW()); -- 0=Lunes, 6=Domingo
    
    -- Determinar fecha de inicio según el período
    CASE p_periodo
        WHEN 'diario' THEN
            SET v_fecha_inicio = DATE(NOW());
            
        WHEN 'semanal' THEN
            -- Inicio de semana (Lunes)
            SET v_fecha_inicio = DATE_SUB(DATE(NOW()), INTERVAL v_dia_semana DAY);
            
        WHEN 'mensual' THEN
            -- Primer día del mes actual
            SET v_fecha_inicio = DATE(CONCAT(v_año, '-', v_mes, '-01'));
            
        WHEN 'anual' THEN
            -- Primer día del año actual
            SET v_fecha_inicio = DATE(CONCAT(v_año, '-01-01'));
            
        ELSE
            -- Por defecto: mensual
            SET v_fecha_inicio = DATE(CONCAT(v_año, '-', v_mes, '-01'));
    END CASE;
    
    -- Calcular consumo
    SELECT COALESCE(SUM(m.monto), 0) as consumoActual
    FROM movimiento m
    WHERE m.nombreUsuario = p_nombre_usuario
      AND m.nombreConcepto = p_nombre_concepto
      AND m.fecha >= v_fecha_inicio
      AND m.delete_at IS NULL;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_get_personalizacion_by_id$$

CREATE PROCEDURE sp_get_personalizacion_by_id(
    IN p_id_personalizacion INT
)
BEGIN
    SELECT 
        idPersonalizacion,
        limiteGasto,
        activo,
        montoPlanificado,
        tipoPeriodoPlanificado,
        tipoPeriodoLimite,
        diaPeriodoPlanificado,
        notificacion,
        nombreUsuario,
        nombreConcepto,
        correoFamilia,
        delete_at
    FROM personalizacionconcepto
    WHERE idPersonalizacion = p_id_personalizacion
      AND delete_at IS NULL;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_create_personalizacion$$

CREATE PROCEDURE sp_create_personalizacion(
    IN p_limite_gasto DECIMAL(10,2),
    IN p_activo BOOLEAN,
    IN p_monto_planificado DECIMAL(10,2),
    IN p_tipo_periodo_planificado VARCHAR(50),
    IN p_tipo_periodo_limite VARCHAR(50),
    IN p_dia_periodo_planificado TINYINT,
    IN p_notificacion BOOLEAN,
    IN p_nombre_usuario VARCHAR(255),
    IN p_nombre_concepto VARCHAR(255),
    IN p_correo_familia VARCHAR(255)
)
BEGIN
    INSERT INTO personalizacionconcepto (
        limiteGasto,
        activo,
        montoPlanificado,
        tipoPeriodoPlanificado,
        tipoPeriodoLimite,
        diaPeriodoPlanificado,
        notificacion,
        nombreUsuario,
        nombreConcepto,
        correoFamilia
    )
    VALUES (
        p_limite_gasto,
        p_activo,
        p_monto_planificado,
        p_tipo_periodo_planificado,
        p_tipo_periodo_limite,
        p_dia_periodo_planificado,
        p_notificacion,
        p_nombre_usuario,
        p_nombre_concepto,
        p_correo_familia
    );
    
    -- Retornar el ID generado
    SELECT LAST_INSERT_ID() as idPersonalizacion;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_update_personalizacion$$

CREATE PROCEDURE sp_update_personalizacion(
    IN p_id_personalizacion INT,
    IN p_limite_gasto DECIMAL(10,2),
    IN p_monto_planificado DECIMAL(10,2),
    IN p_tipo_periodo_planificado VARCHAR(50),
    IN p_tipo_periodo_limite VARCHAR(50),
    IN p_dia_periodo_planificado TINYINT,
    IN p_activo BOOLEAN
)
BEGIN
    UPDATE personalizacionconcepto 
    SET limiteGasto = p_limite_gasto,
        montoPlanificado = p_monto_planificado,
        tipoPeriodoPlanificado = p_tipo_periodo_planificado,
        tipoPeriodoLimite = p_tipo_periodo_limite,
        diaPeriodoPlanificado = p_dia_periodo_planificado,
        activo = p_activo
    WHERE idPersonalizacion = p_id_personalizacion
      AND delete_at IS NULL;
    
    -- Retornar filas afectadas
    SELECT ROW_COUNT() as filasAfectadas;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_delete_personalizacion$$

CREATE PROCEDURE sp_delete_personalizacion(
    IN p_id_personalizacion INT
)
BEGIN
    UPDATE personalizacionconcepto 
    SET delete_at = NOW()
    WHERE idPersonalizacion = p_id_personalizacion
      AND delete_at IS NULL;
    
    -- Retornar filas afectadas
    SELECT ROW_COUNT() as filasAfectadas;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_get_usuario_que_asigno$$

CREATE PROCEDURE sp_get_usuario_que_asigno(
    IN p_id_personalizacion INT
)
BEGIN
    SELECT COALESCE(u.nombrePersonal, 'Desconocido') as nombrePersonal
    FROM personalizacionconcepto p
    INNER JOIN usuario u ON p.nombreUsuario = u.nombreUsuario
                         AND p.correoFamilia = u.correoFamilia
    WHERE p.idPersonalizacion = p_id_personalizacion
      AND p.delete_at IS NULL
    LIMIT 1;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_find_personalizacion_by_usuario_and_concepto$$

CREATE PROCEDURE sp_find_personalizacion_by_usuario_and_concepto(
    IN p_nombre_usuario VARCHAR(255),
    IN p_nombre_concepto VARCHAR(255),
    IN p_correo_familia VARCHAR(255)
)
BEGIN
    SELECT 
        idPersonalizacion,
        montoPlanificado,
        tipoPeriodoPlanificado,
        diaPeriodoPlanificado,
        limiteGasto,
        tipoPeriodoLimite,
        diaPeriodoLimite,
        notificacion,
        activo,
        nombreUsuario,
        nombreConcepto,
        correoFamilia
    FROM personalizacionconcepto
    WHERE nombreUsuario = p_nombre_usuario
      AND nombreConcepto = p_nombre_concepto
      AND correoFamilia = p_correo_familia
      AND delete_at IS NULL
    LIMIT 1;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_find_personalizaciones_by_periodo_diario$$

CREATE PROCEDURE sp_find_personalizaciones_by_periodo_diario()
BEGIN
    SELECT 
        p.idPersonalizacion,
        p.montoPlanificado,
        p.tipoPeriodoPlanificado,
        p.diaPeriodoPlanificado,
        p.limiteGasto,
        p.tipoPeriodoLimite,
        p.diaPeriodoLimite,
        p.notificacion,
        p.activo,
        p.nombreUsuario,
        p.nombreConcepto,
        p.correoFamilia
    FROM personalizacionconcepto p
    WHERE p.activo = TRUE
      AND p.montoPlanificado IS NOT NULL
      AND p.montoPlanificado > 0
      AND p.tipoPeriodoPlanificado = 'diario'
      AND p.delete_at IS NULL
    ORDER BY p.nombreUsuario, p.nombreConcepto;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_find_personalizaciones_by_dia_planificado$$

CREATE PROCEDURE sp_find_personalizaciones_by_dia_planificado(
    IN p_dia TINYINT
)
BEGIN
    SELECT 
        p.idPersonalizacion,
        p.montoPlanificado,
        p.tipoPeriodoPlanificado,
        p.diaPeriodoPlanificado,
        p.limiteGasto,
        p.tipoPeriodoLimite,
        p.diaPeriodoLimite,
        p.notificacion,
        p.activo,
        p.nombreUsuario,
        p.nombreConcepto,
        p.correoFamilia
    FROM personalizacionconcepto p
    WHERE p.activo = TRUE
      AND p.montoPlanificado IS NOT NULL
      AND p.montoPlanificado > 0
      AND p.diaPeriodoPlanificado = p_dia
      AND p.tipoPeriodoPlanificado IN ('mensual', 'quincenal')
      AND p.delete_at IS NULL
    ORDER BY p.nombreUsuario, p.nombreConcepto;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_deshabilitar_personalizacion_para_usuario$$

CREATE PROCEDURE sp_deshabilitar_personalizacion_para_usuario(
    IN p_nombre_concepto VARCHAR(255),
    IN p_correo_familia VARCHAR(255),
    IN p_nombre_usuario VARCHAR(255)
)
BEGIN
    DECLARE v_id_personalizacion INT;
    DECLARE v_existe INT DEFAULT 0;
    
    -- Verificar si existe la personalización (activa o inactiva, pero no eliminada)
    SELECT COUNT(*), COALESCE(MAX(idPersonalizacion), 0)
    INTO v_existe, v_id_personalizacion
    FROM personalizacionconcepto
    WHERE nombreConcepto = p_nombre_concepto
      AND correoFamilia = p_correo_familia
      AND nombreUsuario = p_nombre_usuario
      AND delete_at IS NULL;
    
    IF v_existe = 0 THEN
        -- No existe, crear una nueva con activo = FALSE
        INSERT INTO personalizacionconcepto (
            limiteGasto,
            activo,
            montoPlanificado,
            tipoPeriodoPlanificado,
            tipoPeriodoLimite,
            diaPeriodoPlanificado,
            notificacion,
            nombreUsuario,
            nombreConcepto,
            correoFamilia
        )
        VALUES (
            NULL,
            FALSE,
            NULL,
            NULL,
            NULL,
            NULL,
            FALSE,
            p_nombre_usuario,
            p_nombre_concepto,
            p_correo_familia
        );
        
        -- Retornar el ID creado e indicador de operación
        SELECT LAST_INSERT_ID() as idPersonalizacion, 'created' as operacion;
    ELSE
        -- Ya existe, solo desactivarla
        UPDATE personalizacionconcepto
        SET activo = FALSE
        WHERE idPersonalizacion = v_id_personalizacion
          AND delete_at IS NULL;
        
        -- Retornar el ID actualizado e indicador de operación
        SELECT v_id_personalizacion as idPersonalizacion, 'updated' as operacion;
    END IF;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_habilitar_personalizacion_para_usuario$$

CREATE PROCEDURE sp_habilitar_personalizacion_para_usuario(
    IN p_nombre_concepto VARCHAR(255),
    IN p_correo_familia VARCHAR(255),
    IN p_nombre_usuario VARCHAR(255)
)
BEGIN
    DECLARE v_id_personalizacion INT;
    DECLARE v_existe INT DEFAULT 0;
    
    -- Verificar si existe la personalización (activa o inactiva, pero no eliminada)
    SELECT COUNT(*), COALESCE(MAX(idPersonalizacion), 0)
    INTO v_existe, v_id_personalizacion
    FROM personalizacionconcepto
    WHERE nombreConcepto = p_nombre_concepto
      AND correoFamilia = p_correo_familia
      AND nombreUsuario = p_nombre_usuario
      AND delete_at IS NULL;
    
    IF v_existe = 0 THEN
        -- No existe, crear una nueva con activo = TRUE
        INSERT INTO personalizacionconcepto (
            limiteGasto,
            activo,
            montoPlanificado,
            tipoPeriodoPlanificado,
            tipoPeriodoLimite,
            diaPeriodoPlanificado,
            notificacion,
            nombreUsuario,
            nombreConcepto,
            correoFamilia
        )
        VALUES (
            NULL,
            TRUE,
            NULL,
            NULL,
            NULL,
            NULL,
            FALSE,
            p_nombre_usuario,
            p_nombre_concepto,
            p_correo_familia
        );
        
        -- Retornar el ID creado e indicador de operación
        SELECT LAST_INSERT_ID() as idPersonalizacion, 'created' as operacion;
    ELSE
        -- Ya existe, solo activarla
        UPDATE personalizacionconcepto
        SET activo = TRUE
        WHERE idPersonalizacion = v_id_personalizacion
          AND delete_at IS NULL;
        
        -- Retornar el ID actualizado e indicador de operación
        SELECT v_id_personalizacion as idPersonalizacion, 'updated' as operacion;
    END IF;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_toggle_activo_para_usuario$$

CREATE PROCEDURE sp_toggle_activo_para_usuario(
    IN p_nombre_concepto VARCHAR(255),
    IN p_correo_familia VARCHAR(255),
    IN p_nombre_usuario VARCHAR(255)
)
BEGIN
    DECLARE v_id_personalizacion INT;
    DECLARE v_existe INT DEFAULT 0;
    DECLARE v_activo_actual BOOLEAN;
    DECLARE v_nuevo_estado BOOLEAN;
    
    -- Verificar si existe y obtener estado actual
    SELECT COUNT(*), COALESCE(MAX(idPersonalizacion), 0), COALESCE(MAX(activo), TRUE)
    INTO v_existe, v_id_personalizacion, v_activo_actual
    FROM personalizacionconcepto
    WHERE nombreConcepto = p_nombre_concepto
      AND correoFamilia = p_correo_familia
      AND nombreUsuario = p_nombre_usuario
      AND delete_at IS NULL;
    
    IF v_existe = 0 THEN
        -- No existe, crear con activo = FALSE (toggle desde "por defecto TRUE")
        INSERT INTO personalizacionconcepto (
            limiteGasto,
            activo,
            montoPlanificado,
            tipoPeriodoPlanificado,
            tipoPeriodoLimite,
            diaPeriodoPlanificado,
            notificacion,
            nombreUsuario,
            nombreConcepto,
            correoFamilia
        )
        VALUES (
            NULL,
            FALSE,
            NULL,
            NULL,
            NULL,
            NULL,
            FALSE,
            p_nombre_usuario,
            p_nombre_concepto,
            p_correo_familia
        );
        
        -- Retornar resultado del toggle
        SELECT 
            LAST_INSERT_ID() as idPersonalizacion,
            'created' as operacion,
            FALSE as nuevoEstado,
            TRUE as estadoAnterior;
    ELSE
        -- Ya existe, invertir estado actual
        SET v_nuevo_estado = NOT v_activo_actual;
        
        UPDATE personalizacionconcepto
        SET activo = v_nuevo_estado
        WHERE idPersonalizacion = v_id_personalizacion
          AND delete_at IS NULL;
        
        -- Retornar resultado del toggle
        SELECT 
            v_id_personalizacion as idPersonalizacion,
            'toggled' as operacion,
            v_nuevo_estado as nuevoEstado,
            v_activo_actual as estadoAnterior;
    END IF;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_create_usuario$$

CREATE PROCEDURE sp_create_usuario(
    IN p_nombreUsuario VARCHAR(255),
    IN p_rol VARCHAR(50),
    IN p_contrasenaPersonal VARCHAR(255),
    IN p_nombrePersonal VARCHAR(255),
    IN p_correoFamilia VARCHAR(255)
)
BEGIN
    INSERT INTO usuario (
        nombreUsuario, 
        rol, 
        contraseñaPersonal, 
        nombrePersonal, 
        correoFamilia
    )
    VALUES (
        p_nombreUsuario,
        p_rol,
        p_contrasenaPersonal,
        p_nombrePersonal,
        p_correoFamilia
    );
    
    -- Retornar el ID del usuario creado
    SELECT LAST_INSERT_ID() AS id;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_find_usuario_by_email$$

CREATE PROCEDURE sp_find_usuario_by_email(
    IN p_correoFamilia VARCHAR(255)
)
BEGIN
    SELECT 
        nombreUsuario, 
        rol, 
        contraseñaPersonal, 
        nombrePersonal, 
        correoFamilia, 
        delete_at
    FROM usuario 
    WHERE correoFamilia = p_correoFamilia 
      AND delete_at IS NULL;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_find_usuario_by_nombre$$

CREATE PROCEDURE sp_find_usuario_by_nombre(
    IN p_nombreUsuario VARCHAR(255)
)
BEGIN
    SELECT 
        nombreUsuario, 
        rol, 
        contraseñaPersonal, 
        nombrePersonal, 
        correoFamilia, 
        delete_at
    FROM usuario 
    WHERE nombreUsuario = p_nombreUsuario 
      AND delete_at IS NULL;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_find_usuarios_by_familia$$

CREATE PROCEDURE sp_find_usuarios_by_familia(
    IN p_correoFamilia VARCHAR(255)
)
BEGIN
    SELECT 
        nombreUsuario, 
        rol, 
        contraseñaPersonal, 
        nombrePersonal, 
        correoFamilia, 
        delete_at
    FROM usuario 
    WHERE correoFamilia = p_correoFamilia 
      AND delete_at IS NULL
    ORDER BY rol DESC, nombrePersonal ASC;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_get_all_usuarios_by_familia$$

CREATE PROCEDURE sp_get_all_usuarios_by_familia(
    IN p_correoFamilia VARCHAR(255)
)
BEGIN
    SELECT 
        nombreUsuario, 
        rol, 
        contraseñaPersonal, 
        nombrePersonal, 
        correoFamilia, 
        delete_at
    FROM usuario 
    WHERE correoFamilia = p_correoFamilia 
      AND delete_at IS NULL
    ORDER BY nombrePersonal DESC;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_update_usuario_rol$$

CREATE PROCEDURE sp_update_usuario_rol(
    IN p_nombreUsuario VARCHAR(255),
    IN p_rol TINYINT
)
BEGIN
    UPDATE usuario 
    SET rol = p_rol 
    WHERE nombreUsuario = p_nombreUsuario
      AND delete_at IS NULL;
    
    -- Retornar el número de filas afectadas
    SELECT ROW_COUNT() AS rowsAffected;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_update_usuario_nombre_personal$$

CREATE PROCEDURE sp_update_usuario_nombre_personal(
    IN p_nombreUsuario VARCHAR(255),
    IN p_nuevoNombre VARCHAR(255)
)
BEGIN
    UPDATE usuario 
    SET nombrePersonal = p_nuevoNombre 
    WHERE nombreUsuario = p_nombreUsuario
      AND delete_at IS NULL;
    
    -- Retornar el número de filas afectadas
    SELECT ROW_COUNT() AS rowsAffected;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_update_usuario_nombre_usuario$$

CREATE PROCEDURE sp_update_usuario_nombre_usuario(
    IN p_viejoUsuario VARCHAR(255),
    IN p_nuevoUsuario VARCHAR(255)
)
BEGIN
    -- Verificar si el nuevo nombre de usuario ya existe
    DECLARE usuario_existe INT;
    
    SELECT COUNT(*) INTO usuario_existe
    FROM usuario
    WHERE nombreUsuario = p_nuevoUsuario
      AND delete_at IS NULL;
    
    IF usuario_existe > 0 THEN
        SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'El nombre de usuario ya existe';
    END IF;
    
    -- Actualizar el nombre de usuario
    UPDATE usuario 
    SET nombreUsuario = p_nuevoUsuario 
    WHERE nombreUsuario = p_viejoUsuario
      AND delete_at IS NULL;
    
    -- Retornar el número de filas afectadas
    SELECT ROW_COUNT() AS rowsAffected;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_update_usuario_password$$

CREATE PROCEDURE sp_update_usuario_password(
    IN p_nombreUsuario VARCHAR(255),
    IN p_hashedPassword VARCHAR(255)
)
BEGIN
    UPDATE usuario 
    SET contraseñaPersonal = p_hashedPassword 
    WHERE nombreUsuario = p_nombreUsuario
      AND delete_at IS NULL;
    
    -- Retornar el número de filas afectadas
    SELECT ROW_COUNT() AS rowsAffected;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_update_usuario$$

CREATE PROCEDURE sp_update_usuario(
    IN p_nombreUsuario VARCHAR(255),
    IN p_nombrePersonal VARCHAR(255),
    IN p_rol TINYINT
)
BEGIN
    UPDATE usuario 
    SET nombrePersonal = p_nombrePersonal,
        rol = p_rol
    WHERE nombreUsuario = p_nombreUsuario
      AND delete_at IS NULL;
    
    -- Retornar el número de filas afectadas
    SELECT ROW_COUNT() AS rowsAffected;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_soft_delete_usuario$$

CREATE PROCEDURE sp_soft_delete_usuario(
    IN p_nombreUsuario VARCHAR(255)
)
BEGIN
    UPDATE usuario 
    SET delete_at = NOW()
    WHERE nombreUsuario = p_nombreUsuario
      AND delete_at IS NULL;  -- Solo marcar si no está ya eliminado
    
    -- Retornar el número de filas afectadas
    SELECT ROW_COUNT() AS rowsAffected;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_get_usuario_password$$

CREATE PROCEDURE sp_get_usuario_password(
    IN p_nombreUsuario VARCHAR(255)
)
BEGIN
    SELECT contraseñaPersonal
    FROM usuario 
    WHERE nombreUsuario = p_nombreUsuario 
      AND delete_at IS NULL;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_usuario_exists$$

CREATE PROCEDURE sp_usuario_exists(
    IN p_nombreUsuario VARCHAR(255)
)
BEGIN
    SELECT COUNT(*) AS count
    FROM usuario 
    WHERE nombreUsuario = p_nombreUsuario 
      AND delete_at IS NULL;
END$$

DELIMITER ;

DELIMITER $$

DROP PROCEDURE IF EXISTS sp_delete_usuario_permanently$$

CREATE PROCEDURE sp_delete_usuario_permanently(
    IN p_nombreUsuario VARCHAR(255)
)
BEGIN
    DELETE FROM usuario 
    WHERE nombreUsuario = p_nombreUsuario;
    
    -- Retornar el número de filas afectadas
    SELECT ROW_COUNT() AS rowsAffected;
END$$

DELIMITER ;