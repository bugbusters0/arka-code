SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;
SET collation_connection = 'utf8mb4_unicode_ci';

-- ARCHIVO: 05_movimiento.sql
-- Proyecto: ARKA
-- Descripción: Procedimientos almacenados para gestión de movimientos financieros (ingresos y gastos)


DELIMITER $$

-- Procedure: sp_consultar_movimientos
-- Descripción: Consulta movimientos de múltiples usuarios en una fecha específica
-- Tipo de operación: Lectura
-- Tablas involucradas: movimiento, concepto, tabla temporal temp_usuarios
-- Dependencias: Crea tabla temporal para procesar lista de usuarios

DROP PROCEDURE IF EXISTS sp_consultar_movimientos$$

CREATE PROCEDURE sp_consultar_movimientos(
    IN p_correoFamilia VARCHAR(255),  -- Correo de la familia
    IN p_fecha DATE,                  -- Fecha de consulta
    IN p_nombreUsuarios TEXT          -- Lista de usuarios separada por comas
)
BEGIN

    -- CREACIÓN DE TABLA TEMPORAL

    -- Crear tabla temporal para almacenar los usuarios a consultar
    DROP TEMPORARY TABLE IF EXISTS temp_usuarios;
    CREATE TEMPORARY TABLE temp_usuarios (nombreUsuario VARCHAR(255));
    

    -- POBLACIÓN DE TABLA TEMPORAL

    -- Insertar usuarios desde la lista separada por comas usando SQL dinámico
    SET @sql = CONCAT('INSERT INTO temp_usuarios (nombreUsuario) VALUES ', REPLACE(REPLACE(p_nombreUsuarios, ',', '"),("'), ',', '("'), '")');
    PREPARE stmt FROM @sql;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;
    

    -- CONSULTA DE MOVIMIENTOS

    -- Obtener movimientos con JOIN a concepto para obtener el tipo
    -- Filtrados por familia, fecha y usuarios de la tabla temporal
    

    -- RETORNO DE RESULTADO

    -- Retorno: Lista de movimientos ordenados por usuario y fecha
    SELECT m.idMovimiento, m.fecha, m.monto, m.descripcion, m.nombreUsuario, m.nombreConcepto, m.correoFamilia, c.tipo FROM movimiento m INNER JOIN concepto c ON m.nombreConcepto = c.nombreConcepto AND m.correoFamilia = c.correoFamilia INNER JOIN temp_usuarios tu ON m.nombreUsuario = tu.nombreUsuario WHERE m.correoFamilia = p_correoFamilia AND DATE(m.fecha) = p_fecha AND m.delete_at IS NULL ORDER BY m.nombreUsuario ASC, m.fecha DESC;
    

    -- LIMPIEZA

    DROP TEMPORARY TABLE IF EXISTS temp_usuarios;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_find_movimientos_by_familia_date_tipo
-- Descripción: Busca movimientos de una familia en una fecha filtrados por tipo
-- Tipo de operación: Lectura
-- Tablas involucradas: movimiento, concepto
-- Dependencias: JOIN con concepto para filtrar por tipo

DROP PROCEDURE IF EXISTS sp_find_movimientos_by_familia_date_tipo$$

CREATE PROCEDURE sp_find_movimientos_by_familia_date_tipo(
    IN p_correoFamilia VARCHAR(255),  -- Correo de la familia
    IN p_fecha DATE,                  -- Fecha de consulta
    IN p_tipo TINYINT                 -- Tipo de movimiento (0=gasto, 1=ingreso)
)
BEGIN

    -- CONSULTA CON FILTRO DE TIPO

    -- Obtiene movimientos de la familia en la fecha especificada
    -- Filtrados por tipo de concepto (gasto o ingreso)
    

    -- RETORNO DE RESULTADO

    -- Retorno: Lista de movimientos del tipo especificado ordenados por fecha
    SELECT m.idMovimiento, m.fecha, m.monto, m.descripcion, m.nombreUsuario, m.nombreConcepto, m.correoFamilia FROM movimiento m INNER JOIN concepto c ON m.nombreConcepto = c.nombreConcepto AND m.correoFamilia = c.correoFamilia WHERE m.correoFamilia = p_correoFamilia AND DATE(m.fecha) = p_fecha AND c.tipo = p_tipo AND m.delete_at IS NULL ORDER BY m.fecha DESC;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_update_movimiento
-- Descripción: Actualiza los datos de un movimiento existente
-- Tipo de operación: Escritura
-- Tablas involucradas: movimiento
-- Dependencias: Ninguna

DROP PROCEDURE IF EXISTS sp_update_movimiento$$

CREATE PROCEDURE sp_update_movimiento(
    IN p_idMovimiento INT,              -- ID del movimiento a actualizar
    IN p_fecha DATETIME,                -- Nueva fecha del movimiento
    IN p_monto DECIMAL(10,2),           -- Nuevo monto
    IN p_descripcion TEXT,              -- Nueva descripción
    IN p_nombreConcepto VARCHAR(255)    -- Nuevo concepto asociado
)
BEGIN

    -- ACTUALIZACIÓN DE MOVIMIENTO

    -- Modifica los campos principales del movimiento
    UPDATE movimiento SET fecha = p_fecha, monto = p_monto, descripcion = p_descripcion, nombreConcepto = p_nombreConcepto WHERE idMovimiento = p_idMovimiento AND delete_at IS NULL;
    

    -- RETORNO DE RESULTADO

    -- Retorno:
    -- rows_affected: Número de filas actualizadas (0 o 1)
    SELECT ROW_COUNT() AS rows_affected;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_delete_movimiento
-- Descripción: Realiza eliminación lógica de un movimiento
-- Tipo de operación: Escritura
-- Tablas involucradas: movimiento
-- Dependencias: Ninguna

DROP PROCEDURE IF EXISTS sp_delete_movimiento$$

CREATE PROCEDURE sp_delete_movimiento(
    IN p_idMovimiento INT  -- ID del movimiento a eliminar
)
BEGIN

    -- ELIMINACIÓN LÓGICA

    -- Marca el movimiento con fecha de eliminación
    UPDATE movimiento SET delete_at = NOW() WHERE idMovimiento = p_idMovimiento AND delete_at IS NULL;
    

    -- RETORNO DE RESULTADO

    -- Retorno:
    -- rows_affected: Número de filas eliminadas (0 o 1)
    -- deleted_at: Fecha y hora de la eliminación
    SELECT ROW_COUNT() AS rows_affected, NOW() AS deleted_at;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_create_movimiento
-- Descripción: Crea un nuevo movimiento financiero
-- Tipo de operación: Escritura
-- Tablas involucradas: movimiento
-- Dependencias: Deben existir el usuario, concepto y familia

DROP PROCEDURE IF EXISTS sp_create_movimiento$$

CREATE PROCEDURE sp_create_movimiento(
    IN p_fecha DATETIME,                -- Fecha y hora del movimiento
    IN p_monto DECIMAL(10,2),           -- Monto del movimiento
    IN p_descripcion TEXT,              -- Descripción del movimiento
    IN p_nombreUsuario VARCHAR(255),    -- Usuario que registra el movimiento
    IN p_nombreConcepto VARCHAR(255),   -- Concepto asociado
    IN p_correoFamilia VARCHAR(255)     -- Correo de la familia
)
BEGIN

    -- INSERCIÓN DE MOVIMIENTO

    -- Crea un nuevo registro de movimiento con todos sus datos
    INSERT INTO movimiento (fecha, monto, descripcion, nombreUsuario, nombreConcepto, correoFamilia) VALUES (p_fecha, p_monto, p_descripcion, p_nombreUsuario, p_nombreConcepto, p_correoFamilia);
    

    -- RETORNO DE RESULTADO

    -- Retorno:
    -- id_movimiento: ID del movimiento recién creado
    -- rows_affected: Número de filas insertadas (1)
    SELECT LAST_INSERT_ID() AS id_movimiento, ROW_COUNT() AS rows_affected;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_find_movimiento_by_id
-- Descripción: Busca un movimiento específico por su ID
-- Tipo de operación: Lectura
-- Tablas involucradas: movimiento
-- Dependencias: Ninguna

DROP PROCEDURE IF EXISTS sp_find_movimiento_by_id$$

CREATE PROCEDURE sp_find_movimiento_by_id(
    IN p_idMovimiento INT  -- ID del movimiento a buscar
)
BEGIN

    -- BÚSQUEDA POR ID

    -- Retorna los datos completos del movimiento si existe y está activo
    

    -- RETORNO DE RESULTADO

    -- Retorno: Datos del movimiento o conjunto vacío si no existe
    SELECT idMovimiento, fecha, monto, descripcion, nombreUsuario, nombreConcepto, correoFamilia, delete_at FROM movimiento WHERE idMovimiento = p_idMovimiento AND delete_at IS NULL;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_find_movimientos_by_familia_date
-- Descripción: Busca movimientos de un usuario en una fecha con conceptos activos
-- Tipo de operación: Lectura
-- Tablas involucradas: movimiento, concepto, personalizacionconcepto
-- Dependencias: JOIN con concepto y personalización para filtrar conceptos activos

DROP PROCEDURE IF EXISTS sp_find_movimientos_by_familia_date$$

CREATE PROCEDURE sp_find_movimientos_by_familia_date(
    IN p_nombreUsuario VARCHAR(255),  -- Nombre del usuario
    IN p_fecha DATE                   -- Fecha de consulta
)
BEGIN

    -- CONSULTA CON FILTRO DE CONCEPTOS ACTIVOS

    -- Obtiene movimientos del usuario filtrando solo conceptos activos en su personalización
    -- Triple JOIN: movimiento -> concepto -> personalizacionconcepto
    

    -- RETORNO DE RESULTADO

    -- Retorno: Lista de movimientos con conceptos activos ordenados por fecha
    SELECT m.idMovimiento, m.fecha, m.monto, m.descripcion, m.nombreUsuario, m.nombreConcepto, m.correoFamilia, c.tipo FROM movimiento m INNER JOIN concepto c ON m.nombreConcepto = c.nombreConcepto AND m.correoFamilia = c.correoFamilia INNER JOIN personalizacionconcepto pc ON m.nombreConcepto = pc.nombreConcepto AND m.correoFamilia = pc.correoFamilia AND m.nombreUsuario = pc.nombreUsuario WHERE m.nombreUsuario = p_nombreUsuario AND DATE(m.fecha) = p_fecha AND m.delete_at IS NULL AND pc.activo = 1 AND pc.delete_at IS NULL ORDER BY m.fecha DESC;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_find_movimientos_familia_date_tipo
-- Descripción: Busca movimientos de una familia en una fecha filtrados por tipo (duplicado)
-- Tipo de operación: Lectura
-- Tablas involucradas: movimiento, concepto
-- Dependencias: JOIN con concepto para filtrar por tipo

DROP PROCEDURE IF EXISTS sp_find_movimientos_familia_date_tipo$$

CREATE PROCEDURE sp_find_movimientos_familia_date_tipo(
    IN p_correoFamilia VARCHAR(255),  -- Correo de la familia
    IN p_fecha DATE,                  -- Fecha de consulta
    IN p_tipo TINYINT                 -- Tipo de movimiento (0=gasto, 1=ingreso)
)
BEGIN

    -- CONSULTA CON FILTRO DE TIPO

    -- Obtiene movimientos de la familia en la fecha especificada
    -- Filtrados por tipo de concepto (gasto o ingreso)
    

    -- RETORNO DE RESULTADO

    -- Retorno: Lista de movimientos del tipo especificado ordenados por fecha
    SELECT m.idMovimiento, m.fecha, m.monto, m.descripcion, m.nombreUsuario, m.nombreConcepto, m.correoFamilia 
    FROM movimiento m INNER JOIN concepto c ON m.nombreConcepto = c.nombreConcepto AND m.correoFamilia = c.correoFamilia WHERE m.correoFamilia = p_correoFamilia AND DATE(m.fecha) = p_fecha AND m.delete_at IS NULL ORDER BY m.fecha DESC;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_find_movimientos_by_usuario_and_date
-- Descripción: Busca movimientos de un usuario específico en una fecha
-- Tipo de operación: Lectura
-- Tablas involucradas: movimiento
-- Dependencias: Ninguna

DROP PROCEDURE IF EXISTS sp_find_movimientos_by_usuario_and_date$$

CREATE PROCEDURE sp_find_movimientos_by_usuario_and_date(
    IN p_nombre_usuario VARCHAR(255),  -- Nombre del usuario
    IN p_correo_familia VARCHAR(255),  -- Correo de la familia
    IN p_fecha DATE                    -- Fecha de consulta
)
BEGIN

    -- CONSULTA POR USUARIO Y FECHA

    -- Retorna movimientos del usuario en la fecha especificada
    

    -- RETORNO DE RESULTADO

    -- Retorno: Lista de movimientos ordenados por fecha descendente
    SELECT m.idMovimiento, m.fecha, m.monto, m.descripcion, m.nombreUsuario, m.nombreConcepto, m.correoFamilia 
    FROM movimiento m 
    WHERE m.nombreUsuario = p_nombre_usuario AND m.correoFamilia = p_correo_familia AND DATE(m.fecha) = p_fecha AND m.delete_at IS NULL 
    ORDER BY m.fecha DESC;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_find_movimientos_by_usuario_and_period
-- Descripción: Busca movimientos de un usuario en un rango de fechas
-- Tipo de operación: Lectura
-- Tablas involucradas: movimiento
-- Dependencias: Ninguna

DROP PROCEDURE IF EXISTS sp_find_movimientos_by_usuario_and_period$$

CREATE PROCEDURE sp_find_movimientos_by_usuario_and_period(
    IN p_nombre_usuario VARCHAR(255),  -- Nombre del usuario
    IN p_correo_familia VARCHAR(255),  -- Correo de la familia
    IN p_fecha_inicio DATETIME,        -- Fecha y hora de inicio del período
    IN p_fecha_fin DATETIME            -- Fecha y hora de fin del período
)
BEGIN

    -- CONSULTA POR PERÍODO

    -- Retorna movimientos del usuario en el rango de fechas especificado
    

    -- RETORNO DE RESULTADO

    -- Retorno: Lista de movimientos en el período ordenados por fecha descendente
    SELECT m.idMovimiento, m.fecha, m.monto, m.descripcion, m.nombreUsuario, m.nombreConcepto, m.correoFamilia FROM movimiento m WHERE m.nombreUsuario = p_nombre_usuario AND m.correoFamilia = p_correo_familia AND m.fecha >= p_fecha_inicio AND m.fecha <= p_fecha_fin AND m.delete_at IS NULL ORDER BY m.fecha DESC;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_create_movimiento_automatico
-- Descripción: Crea un movimiento automático con fecha del día actual
-- Tipo de operación: Escritura
-- Tablas involucradas: movimiento
-- Dependencias: Usado para creación automática desde personalizaciones

DROP PROCEDURE IF EXISTS sp_create_movimiento_automatico$$

CREATE PROCEDURE sp_create_movimiento_automatico(
    IN p_nombre_usuario VARCHAR(255),    -- Usuario al que se asigna el movimiento
    IN p_nombre_concepto VARCHAR(255),   -- Concepto asociado
    IN p_correo_familia VARCHAR(255),    -- Correo de la familia
    IN p_monto DECIMAL(10,2),            -- Monto del movimiento
    IN p_descripcion TEXT                -- Descripción del movimiento
)
BEGIN

    -- INSERCIÓN DE MOVIMIENTO AUTOMÁTICO

    -- Crea un movimiento con la fecha actual (CURDATE) de forma automática
    INSERT INTO movimiento (fecha, monto, descripcion, nombreUsuario, nombreConcepto, correoFamilia) VALUES (CURDATE(), p_monto, p_descripcion, p_nombre_usuario, p_nombre_concepto, p_correo_familia);
    

    -- RETORNO DE RESULTADO

    -- Retorno:
    -- idMovimiento: ID del movimiento recién creado
    SELECT LAST_INSERT_ID() as idMovimiento;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_existe_movimiento_hoy
-- Descripción: Verifica si ya existe un movimiento del concepto para el usuario hoy
-- Tipo de operación: Lectura
-- Tablas involucradas: movimiento
-- Dependencias: Usado para evitar duplicados en procesos automáticos

DROP PROCEDURE IF EXISTS sp_existe_movimiento_hoy$$

CREATE PROCEDURE sp_existe_movimiento_hoy(
    IN p_nombre_usuario VARCHAR(255),    -- Nombre del usuario
    IN p_nombre_concepto VARCHAR(255),   -- Nombre del concepto
    IN p_correo_familia VARCHAR(255)     -- Correo de la familia
)
BEGIN

    -- VERIFICACIÓN DE EXISTENCIA HOY

    -- Verifica si ya existe un movimiento del concepto en la fecha actual
    

    -- RETORNO DE RESULTADO

    -- Retorno:
    -- existe: TRUE si existe movimiento hoy, FALSE en caso contrario
    SELECT COUNT(*) > 0 as existe FROM movimiento WHERE nombreUsuario = p_nombre_usuario AND nombreConcepto = p_nombre_concepto AND correoFamilia = p_correo_familia AND DATE(fecha) = CURDATE() AND delete_at IS NULL;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_find_movimientos_by_usuario_date_and_tipo
-- Descripción: Busca movimientos de un usuario en una fecha filtrados por tipo
-- Tipo de operación: Lectura
-- Tablas involucradas: movimiento, concepto
-- Dependencias: JOIN con concepto para filtrar por tipo

DROP PROCEDURE IF EXISTS sp_find_movimientos_by_usuario_date_and_tipo$$

CREATE PROCEDURE sp_find_movimientos_by_usuario_date_and_tipo(
    IN p_nombre_usuario VARCHAR(255),  -- Nombre del usuario
    IN p_correo_familia VARCHAR(255),  -- Correo de la familia
    IN p_fecha DATE,                   -- Fecha de consulta
    IN p_tipo TINYINT                  -- Tipo de movimiento (0=gasto, 1=ingreso)
)
BEGIN

    -- CONSULTA CON FILTRO DE USUARIO Y TIPO

    -- Retorna movimientos del usuario en la fecha filtrados por tipo de concepto
    

    -- RETORNO DE RESULTADO

    -- Retorno: Lista de movimientos del tipo especificado ordenados por fecha
    SELECT m.idMovimiento, m.fecha, m.monto, m.descripcion, m.nombreUsuario, m.nombreConcepto, m.correoFamilia FROM movimiento m INNER JOIN concepto c ON m.nombreConcepto = c.nombreConcepto AND m.correoFamilia = c.correoFamilia WHERE m.nombreUsuario = p_nombre_usuario AND m.correoFamilia = p_correo_familia AND DATE(m.fecha) = p_fecha AND c.tipo = p_tipo AND m.delete_at IS NULL ORDER BY m.fecha DESC;
END$$

DELIMITER ;