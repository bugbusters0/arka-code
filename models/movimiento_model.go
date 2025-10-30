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

// FindByFamiliaAndDate busca movimientos de una familia en una fecha específica
func (m *MovimientoModel) FindByFamiliaAndDate(correoFamilia string, fecha time.Time) ([]entities.Movimiento, error) {
	query := `SELECT m.idMovimiento, m.fecha, m.monto, m.descripcion, m.nombreUsuario, 
	                 m.nombreConcepto, m.correoFamilia, c.tipo
	          FROM movimiento m
	          INNER JOIN concepto c ON m.nombreConcepto = c.nombreConcepto 
	                                AND m.correoFamilia = c.correoFamilia
	          WHERE m.correoFamilia = ? 
	            AND DATE(m.fecha) = DATE(?) 
	            AND m.delete_at IS NULL
	          ORDER BY m.fecha DESC`

	rows, err := database.DB.Query(query, correoFamilia, fecha)
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
