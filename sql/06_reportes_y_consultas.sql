SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;
SET collation_connection = 'utf8mb4_unicode_ci';

-- ARCHIVO: 06_reportes_y_consultas.sql
-- Proyecto: ARKA
-- Descripción: Procedimientos almacenados para reportes y consultas agregadas de movimientos


DELIMITER $$

-- Procedure: sp_get_totales_by_familia_date
-- Descripción: Calcula totales de gastos, ingresos y balance de una familia en una fecha
-- Tipo de operación: Lectura
-- Tablas involucradas: movimiento, concepto
-- Dependencias: JOIN con concepto para clasificar por tipo

DROP PROCEDURE IF EXISTS sp_get_totales_by_familia_date$$

CREATE PROCEDURE sp_get_totales_by_familia_date(
    IN p_correoFamilia VARCHAR(255),  -- Correo de la familia
    IN p_fecha DATE                   -- Fecha de consulta
)
BEGIN

    -- CÁLCULO DE TOTALES AGREGADOS

    -- Calcula totales usando CASE para separar gastos e ingresos
    -- Balance = Ingresos - Gastos
    

    -- RETORNO DE RESULTADO

    -- Retorno:
    -- totalGastos: Suma de todos los gastos (tipo=0)
    -- totalIngresos: Suma de todos los ingresos (tipo=1)
    -- balance: Diferencia entre ingresos y gastos
    -- total_movimientos: Cantidad total de movimientos
    SELECT COALESCE(SUM(CASE WHEN c.tipo = 0 THEN m.monto ELSE 0 END), 0) AS totalGastos, COALESCE(SUM(CASE WHEN c.tipo = 1 THEN m.monto ELSE 0 END), 0) AS totalIngresos, COALESCE(SUM(CASE WHEN c.tipo = 1 THEN m.monto ELSE 0 END), 0) - COALESCE(SUM(CASE WHEN c.tipo = 0 THEN m.monto ELSE 0 END), 0) AS balance, COUNT(*) AS total_movimientos FROM movimiento m INNER JOIN concepto c ON m.nombreConcepto = c.nombreConcepto AND m.correoFamilia = c.correoFamilia WHERE m.correoFamilia = p_correoFamilia AND DATE(m.fecha) = p_fecha AND m.delete_at IS NULL;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_get_totales_by_familia_and_date
-- Descripción: Calcula totales de gastos e ingresos de una familia en una fecha
-- Tipo de operación: Lectura
-- Tablas involucradas: movimiento, concepto
-- Dependencias: JOIN con concepto para clasificar por tipo

DROP PROCEDURE IF EXISTS sp_get_totales_by_familia_and_date$$

CREATE PROCEDURE sp_get_totales_by_familia_and_date(
    IN p_correo_familia VARCHAR(255),  -- Correo de la familia
    IN p_fecha DATE                    -- Fecha de consulta
)
BEGIN

    -- CÁLCULO DE TOTALES

    -- Calcula suma de gastos e ingresos usando CASE
    

    -- RETORNO DE RESULTADO

    -- Retorno:
    -- totalGastos: Suma de todos los gastos (tipo=0)
    -- totalIngresos: Suma de todos los ingresos (tipo=1)
    SELECT COALESCE(SUM(CASE WHEN c.tipo = 0 THEN m.monto ELSE 0 END), 0) as totalGastos, COALESCE(SUM(CASE WHEN c.tipo = 1 THEN m.monto ELSE 0 END), 0) as totalIngresos FROM movimiento m INNER JOIN concepto c ON m.nombreConcepto = c.nombreConcepto AND m.correoFamilia = c.correoFamilia WHERE m.correoFamilia = p_correo_familia AND DATE(m.fecha) = p_fecha AND m.delete_at IS NULL;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_get_totales_by_usuario_and_period
-- Descripción: Calcula totales de un usuario en un rango de fechas
-- Tipo de operación: Lectura
-- Tablas involucradas: movimiento, concepto
-- Dependencias: JOIN con concepto para clasificar por tipo

DROP PROCEDURE IF EXISTS sp_get_totales_by_usuario_and_period$$

CREATE PROCEDURE sp_get_totales_by_usuario_and_period(
    IN p_nombre_usuario VARCHAR(255),  -- Nombre del usuario
    IN p_correo_familia VARCHAR(255),  -- Correo de la familia
    IN p_fecha_inicio DATETIME,        -- Fecha y hora de inicio del período
    IN p_fecha_fin DATETIME            -- Fecha y hora de fin del período
)
BEGIN

    -- CÁLCULO DE TOTALES POR PERÍODO

    -- Calcula suma de gastos e ingresos del usuario en el rango de fechas
    

    -- RETORNO DE RESULTADO

    -- Retorno:
    -- totalGastos: Suma de gastos en el período
    -- totalIngresos: Suma de ingresos en el período
    SELECT COALESCE(SUM(CASE WHEN c.tipo = 0 THEN m.monto ELSE 0 END), 0) as totalGastos, COALESCE(SUM(CASE WHEN c.tipo = 1 THEN m.monto ELSE 0 END), 0) as totalIngresos FROM movimiento m INNER JOIN concepto c ON m.nombreConcepto = c.nombreConcepto AND m.correoFamilia = c.correoFamilia WHERE m.nombreUsuario = p_nombre_usuario AND m.correoFamilia = p_correo_familia AND m.fecha >= p_fecha_inicio AND m.fecha <= p_fecha_fin AND m.delete_at IS NULL;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_get_total_by_concepto_and_period
-- Descripción: Calcula el total de movimientos de un concepto en un período con filtro opcional de usuario
-- Tipo de operación: Lectura
-- Tablas involucradas: movimiento
-- Dependencias: Ninguna

DROP PROCEDURE IF EXISTS sp_get_total_by_concepto_and_period$$

CREATE PROCEDURE sp_get_total_by_concepto_and_period(
    IN p_nombre_concepto VARCHAR(255),  -- Nombre del concepto
    IN p_correo_familia VARCHAR(255),   -- Correo de la familia
    IN p_fecha_inicio DATETIME,         -- Fecha y hora de inicio del período
    IN p_fecha_fin DATETIME,            -- Fecha y hora de fin del período
    IN p_usuario_filtro VARCHAR(255)    -- Usuario para filtrar (NULL o vacío para todos)
)
BEGIN

    -- CÁLCULO DE TOTAL POR CONCEPTO

    -- Suma todos los movimientos del concepto en el período
    -- Filtro opcional por usuario usando OR con NULL/vacío
    

    -- RETORNO DE RESULTADO

    -- Retorno:
    -- total: Suma total de movimientos del concepto (0 si no hay)
    SELECT COALESCE(SUM(m.monto), 0) as total FROM movimiento m WHERE m.nombreConcepto = p_nombre_concepto AND m.correoFamilia = p_correo_familia AND m.fecha >= p_fecha_inicio AND m.fecha <= p_fecha_fin AND m.delete_at IS NULL AND (p_usuario_filtro IS NULL OR p_usuario_filtro = '' OR m.nombreUsuario = p_usuario_filtro);
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_get_totales_by_usuario_and_date
-- Descripción: Calcula totales de gastos e ingresos de un usuario en una fecha específica
-- Tipo de operación: Lectura
-- Tablas involucradas: movimiento, concepto
-- Dependencias: JOIN con concepto para clasificar por tipo

DROP PROCEDURE IF EXISTS sp_get_totales_by_usuario_and_date$$

CREATE PROCEDURE sp_get_totales_by_usuario_and_date(
    IN p_nombre_usuario VARCHAR(255),  -- Nombre del usuario
    IN p_correo_familia VARCHAR(255),  -- Correo de la familia
    IN p_fecha DATE                    -- Fecha de consulta
)
BEGIN

    -- CÁLCULO DE TOTALES POR USUARIO Y FECHA

    -- Calcula suma de gastos e ingresos del usuario en la fecha especificada
    

    -- RETORNO DE RESULTADO

    -- Retorno:
    -- totalGastos: Suma de gastos del usuario en la fecha
    -- totalIngresos: Suma de ingresos del usuario en la fecha
    SELECT COALESCE(SUM(CASE WHEN c.tipo = 0 THEN m.monto ELSE 0 END), 0) as totalGastos, COALESCE(SUM(CASE WHEN c.tipo = 1 THEN m.monto ELSE 0 END), 0) as totalIngresos FROM movimiento m INNER JOIN concepto c ON m.nombreConcepto = c.nombreConcepto AND m.correoFamilia = c.correoFamilia WHERE m.nombreUsuario = p_nombre_usuario AND m.correoFamilia = p_correo_familia AND DATE(m.fecha) = p_fecha AND m.delete_at IS NULL;
END$$

DELIMITER ;

DELIMITER $$

-- Procedure: sp_get_consumo_actual
-- Descripción: Calcula el consumo actual de un concepto según el tipo de período
-- Tipo de operación: Lectura
-- Tablas involucradas: movimiento
-- Dependencias: Usa lógica de fechas para calcular períodos (diario, semanal, mensual, anual)

DROP PROCEDURE IF EXISTS sp_get_consumo_actual$$

CREATE PROCEDURE sp_get_consumo_actual(
    IN p_nombre_usuario VARCHAR(255),  -- Nombre del usuario
    IN p_nombre_concepto VARCHAR(255), -- Nombre del concepto
    IN p_periodo VARCHAR(20)           -- Tipo de período (diario, semanal, mensual, anual)
)
BEGIN
    DECLARE v_fecha_inicio DATETIME;
    DECLARE v_año INT;
    DECLARE v_mes INT;
    DECLARE v_dia INT;
    DECLARE v_dia_semana INT;
    

    -- OBTENCIÓN DE VALORES DE FECHA ACTUAL

    SET v_año = YEAR(NOW());
    SET v_mes = MONTH(NOW());
    SET v_dia = DAY(NOW());
    SET v_dia_semana = WEEKDAY(NOW()); -- 0=Lunes, 6=Domingo
    

    -- DETERMINACIÓN DE FECHA DE INICIO

    -- Calcula la fecha de inicio según el tipo de período
    CASE p_periodo
        WHEN 'diario' THEN
            -- Inicio del día actual
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
    

    -- CÁLCULO DE CONSUMO

    -- Suma todos los movimientos del concepto desde la fecha de inicio
    

    -- RETORNO DE RESULTADO

    -- Retorno:
    -- consumoActual: Suma de movimientos del concepto en el período (0 si no hay)
    SELECT COALESCE(SUM(m.monto), 0) as consumoActual FROM movimiento m WHERE m.nombreUsuario = p_nombre_usuario AND m.nombreConcepto = p_nombre_concepto AND m.fecha >= v_fecha_inicio AND m.delete_at IS NULL;
END$$

DELIMITER ;