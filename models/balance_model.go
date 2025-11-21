package models

import (
	"arka-code/database"
	"arka-code/entities"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

/*
CASO DE USO:CU-003 Consultar Balance
*/

type BDBalance struct{}

var BalanceModelInstance = &BDBalance{}

// ConsultarMovimientos busca y agrupa movimientos por usuario para una familia y fecha específica
// Parámetros:
// - correoFamilia: Identificador de la familia
// - nombreUsuarios: Lista de usuarios a consultar
// - fecha: Fecha específica para filtrar movimientos
// Retorno: Mapa donde la clave es el nombre de usuario y el valor es su lista de movimientos
// Flujo:
// 1. Valida que haya usuarios para consultar
// 2. Construye consulta SQL dinámica con cláusula IN para múltiples usuarios
// 3. Ejecuta consulta con parámetros preparados
// 4. Procesa resultados y agrupa movimientos por usuario
// 5. Asegura que todos los usuarios estén en el mapa (incluso sin movimientos)
// Uso: Balance familiar, reportes agrupados por miembro

func (m *BDBalance) ConsultarMovimientos(correoFamilia string, nombreUsuarios []string, fecha time.Time) (map[string][]entities.Movimiento, error) {
	if len(nombreUsuarios) == 0 {
		return make(map[string][]entities.Movimiento), nil
	}

	placeholders := make([]string, len(nombreUsuarios))
	for i := range nombreUsuarios {
		placeholders[i] = "?"
	}
	inClause := strings.Join(placeholders, ", ")

	// 2. CONSTRUCCIÓN DE LA CONSULTA SQL
	query := fmt.Sprintf(`
        SELECT m.idMovimiento, m.fecha, m.monto, m.descripcion, m.nombreUsuario, 
               m.nombreConcepto, m.correoFamilia, c.tipo
        FROM movimiento m
        INNER JOIN concepto c ON m.nombreConcepto = c.nombreConcepto 
                              AND m.correoFamilia = c.correoFamilia
        WHERE m.correoFamilia = ? 
          AND DATE(m.fecha) = DATE(?) 
          AND m.nombreUsuario IN (%s)
          AND m.delete_at IS NULL
        ORDER BY m.nombreUsuario ASC, m.fecha DESC`,
		inClause)

	// 3. PREPARACIÓN DE ARGUMENTOS
	// Los argumentos deben ser: [correoFamilia, fecha, nombreUsuario1, nombreUsuario2, ...]
	args := []interface{}{correoFamilia, fecha.Format("2006-01-02")}
	for _, user := range nombreUsuarios {
		args = append(args, user)
	}

	// 4. EJECUCIÓN DE LA CONSULTA
	// Nota: Es crucial que database.DB esté inicializado
	rows, err := database.DB.Query(query, args...)
	if err != nil {
		log.Printf("❌ Error en consulta ConsultarMovimientos: %v", err)
		return nil, err
	}
	defer rows.Close()
	// 5. PROCESAMIENTO Y AGRUPACIÓN DE RESULTADOS
	movimientosPorUsuario := make(map[string][]entities.Movimiento)

	for rows.Next() {
		var movimiento entities.Movimiento
		var monto float64
		var descripcion sql.NullString // Usamos sql.NullString para campos TEXT/DEFAULT NULL

		err := rows.Scan(
			&movimiento.IdMovimiento,
			&movimiento.Fecha,
			&monto,
			&descripcion,
			&movimiento.NombreUsuario,
			&movimiento.NombreConcepto,
			&movimiento.CorreoFamilia,
		)
		if err != nil {
			log.Printf("❌ Error escaneando movimiento: %v", err)
			continue
		}

		// Convertir tipos de SQL a Go
		movimiento.Monto = monto
		if descripcion.Valid {
			movimiento.Descripcion = &descripcion.String
		}

		// Agrupar el movimiento en el mapa usando el nombre de usuario como clave
		usuario := movimiento.NombreUsuario
		movimientosPorUsuario[usuario] = append(movimientosPorUsuario[usuario], movimiento)
	}

	// 6. ASEGURAR QUE TODOS LOS USUARIOS ESTÉN EN EL MAPA (incluso si no tienen movimientos)
	// Esto es opcional, pero asegura que el controlador siempre reciba una clave para cada usuario.
	for _, user := range nombreUsuarios {
		if _, exists := movimientosPorUsuario[user]; !exists {
			movimientosPorUsuario[user] = []entities.Movimiento{}
		}
	}

	log.Printf("📋 Consulta exitosa. Movimientos agrupados por %d usuarios.", len(movimientosPorUsuario))
	return movimientosPorUsuario, nil
}

// FindByFamiliaDateAndTipo busca movimientos de una familia filtrados por tipo y fecha
// Parámetros:
// - correoFamilia: Identificador de la familia
// - fecha: Fecha específica para filtrar
// - tipo: Tipo de movimiento (0 = gasto, 1 = ingreso)
// Retorno: Lista de movimientos que cumplen con los criterios
// Flujo:
// - Realiza JOIN con tabla concepto para filtrar por tipo
// - Filtra por fecha exacta y familia
// - Excluye movimientos eliminados (soft delete)
// - Ordena por fecha descendente
// Uso: Vista diaria de gastos/ingresos, reportes por tipo

func (m *BDBalance) FindByFamiliaDateAndTipo(correoFamilia string, fecha time.Time, tipo int8) ([]entities.Movimiento, error) {
	query := `SELECT m.idMovimiento, m.fecha, m.monto, m.descripcion, m.nombreUsuario, 
	                 m.nombreConcepto, m.correoFamilia
	          FROM movimiento m
	          INNER JOIN concepto c ON m.nombreConcepto = c.nombreConcepto 
	                                AND m.correoFamilia = c.correoFamilia
	          WHERE m.correoFamilia = ? 
	            AND DATE(m.fecha) = DATE(?) 
	            AND c.tipo = ?
	            AND m.delete_at IS NULL
	          ORDER BY m.fecha DESC`

	rows, err := database.DB.Query(query, correoFamilia, fecha, tipo)
	if err != nil {
		log.Printf("❌ Error en consulta FindByFamiliaDateAndTipo: %v", err)
		return nil, err
	}
	defer rows.Close()

	movimientos := []entities.Movimiento{}
	for rows.Next() {
		var movimiento entities.Movimiento

		err := rows.Scan(
			&movimiento.IdMovimiento,
			&movimiento.Fecha,
			&movimiento.Monto,
			&movimiento.Descripcion,
			&movimiento.NombreUsuario,
			&movimiento.NombreConcepto,
			&movimiento.CorreoFamilia,
		)
		if err != nil {
			log.Printf("❌ Error escaneando movimiento: %v", err)
			continue
		}
		movimientos = append(movimientos, movimiento)
	}

	tipoStr := "gastos"
	if tipo == 1 {
		tipoStr = "ingresos"
	}
	log.Printf("📋 %s encontrados: %d para fecha %s", tipoStr, len(movimientos), fecha.Format("2006-01-02"))
	return movimientos, nil
}

// FindByID busca un movimiento específico por su ID único
// Parámetros:
// - idMovimiento: Identificador único del movimiento
// Retorno: Puntero al movimiento encontrado o nil si no existe
// Flujo:
// - Consulta movimiento por ID exacto
// - Excluye movimientos eliminados (soft delete)
// - Retorna error si hay problemas de base de datos
// Uso: Edición de movimientos, verificación de existencia, operaciones CRUD

func (m *BDBalance) FindByID(idMovimiento int) (*entities.Movimiento, error) {
	query := `SELECT idMovimiento, fecha, monto, descripcion, nombreUsuario, 
	                 nombreConcepto, correoFamilia, delete_at
	          FROM movimiento 
	          WHERE idMovimiento = ? AND delete_at IS NULL`

	movimiento := &entities.Movimiento{}
	err := database.DB.QueryRow(query, idMovimiento).Scan(
		&movimiento.IdMovimiento,
		&movimiento.Fecha,
		&movimiento.Monto,
		&movimiento.Descripcion,
		&movimiento.NombreUsuario,
		&movimiento.NombreConcepto,
		&movimiento.CorreoFamilia,
		&movimiento.DeleteAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Printf("❌ Error buscando movimiento ID %d: %v", idMovimiento, err)
		return nil, err
	}

	return movimiento, nil
}

// Update actualiza un movimiento existente en la base de datos
// Parámetros:
// - movimiento: Estructura con los datos actualizados del movimiento
// Retorno: Error si la actualización falla
// Flujo:
// - Actualiza campos: fecha, monto, descripción y concepto
// - Solo afecta movimientos no eliminados (soft delete)
// - Registra número de filas afectadas para verificación
// Uso: Modificación de movimientos existentes, corrección de datos

func (m *BDBalance) Update(movimiento *entities.Movimiento) error {
	query := `UPDATE movimiento 
	          SET fecha = ?, monto = ?, descripcion = ?, nombreConcepto = ?
	          WHERE idMovimiento = ? AND delete_at IS NULL`

	result, err := database.DB.Exec(query,
		movimiento.Fecha,
		movimiento.Monto,
		movimiento.Descripcion,
		movimiento.NombreConcepto,
		movimiento.IdMovimiento)

	if err != nil {
		log.Printf("❌ Error actualizando movimiento: %v", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("✅ Movimiento actualizado - ID: %d, Filas afectadas: %d", movimiento.IdMovimiento, rowsAffected)
	return nil
}

// Delete realiza una eliminación lógica (soft delete) de un movimiento
// Parámetros:
// - idMovimiento: Identificador único del movimiento a eliminar
// Retorno: Error si la eliminación falla
// Flujo:
// - Establece delete_at con la fecha/hora actual
// - No elimina físicamente el registro
// - Permite recuperación de datos si es necesario
// Uso: Eliminación segura de movimientos, mantenimiento de historial
func (m *BDBalance) Delete(idMovimiento int) error {
	query := `UPDATE movimiento SET delete_at = NOW() WHERE idMovimiento = ?`

	result, err := database.DB.Exec(query, idMovimiento)
	if err != nil {
		log.Printf("❌ Error eliminando movimiento: %v", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("✅ Movimiento eliminado - ID: %d, Filas afectadas: %d", idMovimiento, rowsAffected)
	return nil
}

// GetTotalesByFamiliaAndDate calcula los totales de ingresos y gastos de una familia en una fecha específica
// Parámetros:
// - correoFamilia: Identificador de la familia
// - fecha: Fecha específica para el cálculo
// Retorno: Total de gastos, total de ingresos y error si existe
// Flujo:
// - Utiliza conditional aggregation con CASE statements
// - Suma montos separados por tipo de concepto (gasto/ingreso)
// - Maneja valores nulos con sql.NullFloat64
// Uso: Resumen diario, cálculo de balance, dashboards
func (m *BDBalance) GetTotalesByFamiliaAndDate(correoFamilia string, fecha time.Time) (float64, float64, error) {
	query := `SELECT 
	          SUM(CASE WHEN c.tipo = 0 THEN m.monto ELSE 0 END) as totalGastos,
	          SUM(CASE WHEN c.tipo = 1 THEN m.monto ELSE 0 END) as totalIngresos
	          FROM movimiento m
	          INNER JOIN concepto c ON m.nombreConcepto = c.nombreConcepto 
	                                AND m.correoFamilia = c.correoFamilia
	          WHERE m.correoFamilia = ? 
	            AND DATE(m.fecha) = DATE(?)
	            AND m.delete_at IS NULL`

	var totalGastos, totalIngresos sql.NullFloat64
	err := database.DB.QueryRow(query, correoFamilia, fecha).Scan(&totalGastos, &totalIngresos)
	if err != nil {
		log.Printf("❌ Error obteniendo totales: %v", err)
		return 0, 0, err
	}

	gastos := 0.0
	if totalGastos.Valid {
		gastos = totalGastos.Float64
	}

	ingresos := 0.0
	if totalIngresos.Valid {
		ingresos = totalIngresos.Float64
	}

	log.Printf("💰 Totales - Gastos: %.2f, Ingresos: %.2f", gastos, ingresos)
	return gastos, ingresos, nil
}
