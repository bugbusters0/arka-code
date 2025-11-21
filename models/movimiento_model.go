package models

import (
	"arka-code/database"
	"arka-code/entities"
	"database/sql"
	"log"
	"time"
)

type MovimientoModel struct{}

var MovimientoModelInstance = &MovimientoModel{}

// Create crea un nuevo movimiento
func (m *MovimientoModel) Create(movimiento *entities.Movimiento) error {
	query := `INSERT INTO movimiento (fecha, monto, descripcion, nombreUsuario, nombreConcepto, correoFamilia)
	          VALUES (?, ?, ?, ?, ?, ?)`

	result, err := database.DB.Exec(query,
		movimiento.Fecha,
		movimiento.Monto,
		movimiento.Descripcion,
		movimiento.NombreUsuario,
		movimiento.NombreConcepto,
		movimiento.CorreoFamilia)

	if err != nil {
		log.Printf("❌ Error en consulta INSERT movimiento: %v", err)
		return err
	}

	id, _ := result.LastInsertId()
	movimiento.IdMovimiento = int(id)

	log.Printf("✅ Movimiento creado - ID: %d", id)
	return nil
}

// FindByFamilia busca todos los conceptos de una familia por tipo
func (m *MovimientoModel) FindByFamilia(correoFamilia string, tipo string) ([]entities.Concepto, error) {
	var tipoInt int8
	if tipo == "ingreso" {
		tipoInt = 1
	} else {
		tipoInt = 0
	}

	// ✅ CORREGIDO: JOIN con personalizacionconcepto para obtener el estado activo
	query := `SELECT c.nombreConcepto, c.correoFamilia, c.tipo, c.icono, c.color, c.nombreUsuario, c.delete_at FROM 
    concepto c 
		LEFT JOIN personalizacionconcepto pc 
			ON c.nombreConcepto = pc.nombreConcepto 
			AND c.correoFamilia = pc.correoFamilia 
			AND pc.nombreUsuario = c.nombreUsuario
		WHERE 
				c.correoFamilia = ?
				AND c.tipo = ? 
				AND c.delete_at IS NULL 
				AND pc.activo = 1
		ORDER BY 
				c.nombreConcepto;`

	rows, err := database.DB.Query(query, correoFamilia, tipoInt)
	if err != nil {
		log.Printf("❌ Error en FindByFamilia: %v", err)
		return nil, err
	}
	defer rows.Close()

	conceptos := []entities.Concepto{}
	for rows.Next() {
		var concepto entities.Concepto
		err := rows.Scan(
			&concepto.NombreConcepto,
			&concepto.CorreoFamilia,
			&concepto.Tipo,
			&concepto.Icono,
			&concepto.Color,
			&concepto.NombreUsuario,
			&concepto.DeleteAt,
		)
		if err != nil {
			log.Printf("❌ Error escaneando concepto: %v", err)
			return nil, err
		}
		conceptos = append(conceptos, concepto)
	}

	log.Printf("📋 Conceptos encontrados: %d (tipo: %s)", len(conceptos), tipo)
	return conceptos, nil
}

// FindByFamiliaAndDate busca movimientos de una familia en una fecha específica
func (m *MovimientoModel) FindByFamiliaAndDate(nombreUsuario string, fecha time.Time) ([]entities.Movimiento, error) {
	query := `SELECT m.idMovimiento, m.fecha, m.monto, m.descripcion, m.nombreUsuario, 
       m.nombreConcepto, m.correoFamilia, c.tipo
				FROM movimiento m
				INNER JOIN concepto c ON m.nombreConcepto = c.nombreConcepto 
															AND m.correoFamilia = c.correoFamilia
				INNER JOIN personalizacionconcepto pc ON m.nombreConcepto = pc.nombreConcepto 
																							AND m.correoFamilia = pc.correoFamilia 
																							AND m.nombreUsuario = pc.nombreUsuario
				WHERE m.nombreUsuario = ? 
					AND DATE(m.fecha) = DATE(?) 
					AND m.delete_at IS NULL
					AND pc.activo = 1
					AND pc.delete_at IS NULL
				ORDER BY m.fecha DESC`

	rows, err := database.DB.Query(query, nombreUsuario, fecha)
	if err != nil {
		log.Printf("❌ Error en consulta FindByFamiliaAndDate: %v", err)
		return nil, err
	}
	defer rows.Close()

	movimientos := []entities.Movimiento{}
	for rows.Next() {
		var movimiento entities.Movimiento
		var tipo int8

		err := rows.Scan(
			&movimiento.IdMovimiento,
			&movimiento.Fecha,
			&movimiento.Monto,
			&movimiento.Descripcion,
			&movimiento.NombreUsuario,
			&movimiento.NombreConcepto,
			&movimiento.CorreoFamilia,
			&tipo,
		)
		if err != nil {
			log.Printf("❌ Error escaneando movimiento: %v", err)
			continue
		}
		movimientos = append(movimientos, movimiento)
	}

	log.Printf("📋 Movimientos encontrados: %d para fecha %s", len(movimientos), fecha.Format("2006-01-02"))
	return movimientos, nil
}

// FindByFamiliaDateAndTipo busca movimientos filtrados por tipo (gasto/ingreso)
func (m *MovimientoModel) FindByFamiliaDateAndTipo(correoFamilia string, fecha time.Time, tipo int8) ([]entities.Movimiento, error) {
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

// FindByID busca un movimiento por su ID
func (m *MovimientoModel) FindByID(idMovimiento int) (*entities.Movimiento, error) {
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

// Update actualiza un movimiento existente
func (m *MovimientoModel) Update(movimiento *entities.Movimiento) error {
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

// Delete realiza soft delete de un movimiento
func (m *MovimientoModel) Delete(idMovimiento int) error {
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

// GetTotalesByFamiliaAndDate obtiene totales de ingresos y gastos de una fecha
func (m *MovimientoModel) GetTotalesByFamiliaAndDate(correoFamilia string, fecha time.Time) (float64, float64, error) {
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

// Agregar estas funciones al archivo models/movimiento.go

// GetTotalesByUsuarioAndPeriod obtiene totales de ingresos y gastos de un usuario en un período
func (m *MovimientoModel) GetTotalesByUsuarioAndPeriod(nombreUsuario, correoFamilia string, inicio, fin time.Time) (float64, float64, error) {
	query := `SELECT 
	          SUM(CASE WHEN c.tipo = 0 THEN m.monto ELSE 0 END) as totalGastos,
	          SUM(CASE WHEN c.tipo = 1 THEN m.monto ELSE 0 END) as totalIngresos
	          FROM movimiento m
	          INNER JOIN concepto c ON m.nombreConcepto = c.nombreConcepto 
	                                AND m.correoFamilia = c.correoFamilia
	          WHERE m.nombreUsuario = ? 
	            AND m.correoFamilia = ?
	            AND m.fecha >= ?
	            AND m.fecha <= ?
	            AND m.delete_at IS NULL`

	var totalGastos, totalIngresos sql.NullFloat64

	// ✅ IMPORTANTE: Pasar inicio y fin correctamente
	err := database.DB.QueryRow(query, nombreUsuario, correoFamilia, inicio, fin).Scan(&totalGastos, &totalIngresos)
	if err != nil {
		log.Printf("❌ Error obteniendo totales por período: %v", err)
		log.Printf("   Usuario: %s, Inicio: %s, Fin: %s", nombreUsuario, inicio.Format("2006-01-02"), fin.Format("2006-01-02"))
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

	log.Printf("💰 Totales Usuario %s - Período %s a %s - Gastos: %.2f, Ingresos: %.2f",
		nombreUsuario, inicio.Format("2006-01-02"), fin.Format("2006-01-02"), gastos, ingresos)
	return gastos, ingresos, nil
}

// FindByUsuarioAndDate busca movimientos de un usuario en una fecha específica
func (m *MovimientoModel) FindByUsuarioAndDate(nombreUsuario, correoFamilia string, fecha time.Time) ([]entities.Movimiento, error) {
	query := `SELECT m.idMovimiento, m.fecha, m.monto, m.descripcion, m.nombreUsuario, 
	                 m.nombreConcepto, m.correoFamilia
	          FROM movimiento m
	          WHERE m.nombreUsuario = ? 
	            AND m.correoFamilia = ?
	            AND DATE(m.fecha) = DATE(?) 
	            AND m.delete_at IS NULL
	          ORDER BY m.fecha DESC`

	rows, err := database.DB.Query(query, nombreUsuario, correoFamilia, fecha)
	if err != nil {
		log.Printf("❌ Error en consulta FindByUsuarioAndDate: %v", err)
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

	log.Printf("📋 Movimientos encontrados: %d para usuario %s en fecha %s",
		len(movimientos), nombreUsuario, fecha.Format("2006-01-02"))
	return movimientos, nil
}

// FindByUsuarioAndPeriod busca movimientos de un usuario en un período
// FindByUsuarioAndPeriod busca movimientos de un usuario en un período
func (m *MovimientoModel) FindByUsuarioAndPeriod(nombreUsuario, correoFamilia string, inicio, fin time.Time) ([]entities.Movimiento, error) {
	query := `SELECT m.idMovimiento, m.fecha, m.monto, m.descripcion, m.nombreUsuario, 
	                 m.nombreConcepto, m.correoFamilia
	          FROM movimiento m
	          WHERE m.nombreUsuario = ? 
	            AND m.correoFamilia = ?
	            AND m.fecha >= ?
	            AND m.fecha <= ?
	            AND m.delete_at IS NULL
	          ORDER BY m.fecha DESC`

	rows, err := database.DB.Query(query, nombreUsuario, correoFamilia, inicio, fin)
	if err != nil {
		log.Printf("❌ Error en consulta FindByUsuarioAndPeriod: %v", err)
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

	log.Printf("📋 Movimientos encontrados: %d para usuario %s en período %s a %s",
		len(movimientos), nombreUsuario, inicio.Format("2006-01-02"), fin.Format("2006-01-02"))
	return movimientos, nil
}

// GetTotalByConceptoAndPeriod obtiene el total de un concepto en un período
// GetTotalByConceptoAndPeriod obtiene el total de un concepto en un período
func (m *MovimientoModel) GetTotalByConceptoAndPeriod(nombreConcepto, correoFamilia string, inicio, fin time.Time, usuarioFiltro string) (float64, error) {
	query := `SELECT COALESCE(SUM(m.monto), 0) as total
	          FROM movimiento m
	          WHERE m.nombreConcepto = ? 
	            AND m.correoFamilia = ?
	            AND m.fecha >= ?
	            AND m.fecha <= ?
	            AND m.delete_at IS NULL`

	args := []interface{}{nombreConcepto, correoFamilia, inicio, fin}

	if usuarioFiltro != "" {
		query += " AND m.nombreUsuario = ?"
		args = append(args, usuarioFiltro)
	}

	var total float64
	err := database.DB.QueryRow(query, args...).Scan(&total)
	if err != nil {
		log.Printf("❌ Error obteniendo total por concepto: %v", err)
		log.Printf("   Concepto: %s, Inicio: %s, Fin: %s", nombreConcepto, inicio.Format("2006-01-02"), fin.Format("2006-01-02"))
		return 0, err
	}

	log.Printf("📊 Total concepto %s: %.2f (período: %s a %s)",
		nombreConcepto, total, inicio.Format("2006-01-02"), fin.Format("2006-01-02"))
	return total, nil
}

// FUNCIONES PARA CRONJOB
// CreateAutomatico crea un movimiento automático
func (m *MovimientoModel) CreateAutomatico(nombreUsuario, nombreConcepto, correoFamilia string, monto float64, descripcion string) error {
	query := `INSERT INTO movimiento (fecha, monto, descripcion, nombreUsuario, nombreConcepto, correoFamilia)
              VALUES (CURDATE(), ?, ?, ?, ?, ?)`

	_, err := database.DB.Exec(query, monto, descripcion, nombreUsuario, nombreConcepto, correoFamilia)
	if err != nil {
		log.Printf("❌ Error creando movimiento automático: %v", err)
		return err
	}

	log.Printf("✅ Movimiento automático creado: %s - %s - S/ %.2f",
		nombreUsuario, nombreConcepto, monto)
	return nil
}

// ExisteMovimientoHoy verifica si ya existe un movimiento para hoy
func (m *MovimientoModel) ExisteMovimientoHoy(nombreUsuario, nombreConcepto, correoFamilia string) (bool, error) {
	query := `SELECT COUNT(*) FROM movimiento 
              WHERE nombreUsuario = ? 
                AND nombreConcepto = ? 
                AND correoFamilia = ?
                AND fecha = CURDATE()
                AND delete_at IS NULL`

	var count int
	err := database.DB.QueryRow(query, nombreUsuario, nombreConcepto, correoFamilia).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
