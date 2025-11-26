package models

import (
	"arka-code/database"
	"arka-code/entities"
	"database/sql"
	/*"fmt"*/
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

	// Convertir array de usuarios a string separado por comas
	usuariosStr := strings.Join(nombreUsuarios, ",")
	
	query := `CALL sp_consultar_movimientos(?, ?, ?)`

	rows, err := database.DB.Query(query, correoFamilia, fecha.Format("2006-01-02"), usuariosStr)
	if err != nil {
		log.Printf("❌ Error al ejecutar sp_consultar_movimientos: %v", err)
		return nil, err
	}
	defer rows.Close()

	movimientosPorUsuario := make(map[string][]entities.Movimiento)

	for rows.Next() {
		var movimiento entities.Movimiento
		var monto float64
		var descripcion sql.NullString
		var tipoConcepto int8

		err := rows.Scan(
			&movimiento.IdMovimiento,
			&movimiento.Fecha,
			&monto,
			&descripcion,
			&movimiento.NombreUsuario,
			&movimiento.NombreConcepto,
			&movimiento.CorreoFamilia,
			&tipoConcepto,
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

	// Asegurar que todos los usuarios estén en el mapa (incluso si no tienen movimientos)
	for _, user := range nombreUsuarios {
		if _, exists := movimientosPorUsuario[user]; !exists {
			movimientosPorUsuario[user] = []entities.Movimiento{}
		}
	}

	log.Printf("📋 Consulta exitosa - Fecha: %s, Usuarios: %d, Total movimientos: %d", 
		fecha.Format("2006-01-02"), len(nombreUsuarios), contarMovimientos(movimientosPorUsuario))
	
	return movimientosPorUsuario, nil
}

// Helper function para contar total de movimientos
func contarMovimientos(movimientos map[string][]entities.Movimiento) int {
	total := 0
	for _, movs := range movimientos {
		total += len(movs)
	}
	return total
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
	query := `CALL sp_find_movimientos_by_familia_date_tipo(?, ?, ?)`

	rows, err := database.DB.Query(query, correoFamilia, fecha.Format("2006-01-02"), tipo)
	if err != nil {
		log.Printf("❌ Error al ejecutar sp_find_movimientos_by_familia_date_tipo: %v", err)
		return nil, err
	}
	defer rows.Close()

	movimientos := []entities.Movimiento{}
	for rows.Next() {
		var movimiento entities.Movimiento
		var descripcion sql.NullString

		err := rows.Scan(
			&movimiento.IdMovimiento,
			&movimiento.Fecha,
			&movimiento.Monto,
			&descripcion,
			&movimiento.NombreUsuario,
			&movimiento.NombreConcepto,
			&movimiento.CorreoFamilia,
		)
		if err != nil {
			log.Printf("❌ Error escaneando movimiento: %v", err)
			continue
		}

		// Manejar descripción NULL
		if descripcion.Valid {
			movimiento.Descripcion = &descripcion.String
		}

		movimientos = append(movimientos, movimiento)
	}

	tipoStr := "gastos"
	if tipo == 1 {
		tipoStr = "ingresos"
	}
	log.Printf("📋 %s encontrados: %d para fecha %s (familia: %s)", 
		tipoStr, len(movimientos), fecha.Format("2006-01-02"), correoFamilia)
	
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
	query := `CALL sp_find_movimiento_by_id(?)`

	movimiento := &entities.Movimiento{}
	var descripcion sql.NullString
	
	err := database.DB.QueryRow(query, idMovimiento).Scan(
		&movimiento.IdMovimiento,
		&movimiento.Fecha,
		&movimiento.Monto,
		&descripcion,
		&movimiento.NombreUsuario,
		&movimiento.NombreConcepto,
		&movimiento.CorreoFamilia,
		&movimiento.DeleteAt,
	)

	if err == sql.ErrNoRows {
		log.Printf("ℹ️ Movimiento no encontrado - ID: %d", idMovimiento)
		return nil, nil
	}

	if err != nil {
		log.Printf("❌ Error al ejecutar sp_find_movimiento_by_id - ID: %d, Error: %v", idMovimiento, err)
		return nil, err
	}

	// Manejar descripción NULL
	if descripcion.Valid {
		movimiento.Descripcion = &descripcion.String
	}

	log.Printf("✅ Movimiento encontrado - ID: %d, Usuario: %s, Concepto: %s", 
		idMovimiento, movimiento.NombreUsuario, movimiento.NombreConcepto)

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
	query := `CALL sp_update_movimiento(?, ?, ?, ?, ?)`

	var rowsAffected int64

	err := database.DB.QueryRow(query,
		movimiento.IdMovimiento,
		movimiento.Fecha,
		movimiento.Monto,
		movimiento.Descripcion,
		movimiento.NombreConcepto).Scan(&rowsAffected)

	if err != nil {
		log.Printf("❌ Error al ejecutar sp_update_movimiento - ID: %d, Error: %v", 
			movimiento.IdMovimiento, err)
		return err
	}

	if rowsAffected == 0 {
		log.Printf("⚠️ Movimiento no encontrado o ya eliminado - ID: %d", movimiento.IdMovimiento)
	} else {
		log.Printf("✅ Movimiento actualizado - ID: %d, Concepto: %s, Monto: %.2f, Filas: %d", 
			movimiento.IdMovimiento, movimiento.NombreConcepto, movimiento.Monto, rowsAffected)
	}

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
	query := `CALL sp_delete_movimiento(?)`

	var rowsAffected int64
	var deletedAt time.Time

	err := database.DB.QueryRow(query, idMovimiento).Scan(&rowsAffected, &deletedAt)
	if err != nil {
		log.Printf("❌ Error al ejecutar sp_delete_movimiento - ID: %d, Error: %v", idMovimiento, err)
		return err
	}

	if rowsAffected == 0 {
		log.Printf("⚠️ Movimiento no encontrado o ya eliminado - ID: %d", idMovimiento)
	} else {
		log.Printf("✅ Movimiento eliminado (soft delete) - ID: %d, Fecha eliminación: %s", 
			idMovimiento, deletedAt.Format("2006-01-02 15:04:05"))
	}

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
	query := `CALL sp_get_totales_by_familia_date(?, ?)`

	var totalGastos, totalIngresos, balance float64
	var totalMovimientos int

	err := database.DB.QueryRow(query, correoFamilia, fecha.Format("2006-01-02")).
		Scan(&totalGastos, &totalIngresos, &balance, &totalMovimientos)
	
	if err != nil {
		log.Printf("❌ Error al ejecutar sp_get_totales_by_familia_date - Familia: %s, Fecha: %s, Error: %v", 
			correoFamilia, fecha.Format("2006-01-02"), err)
		return 0, 0, err
	}

	log.Printf("💰 Totales del día %s - Gastos: %.2f, Ingresos: %.2f, Balance: %.2f, Movimientos: %d", 
		fecha.Format("2006-01-02"), totalGastos, totalIngresos, balance, totalMovimientos)
	
	return totalGastos, totalIngresos, nil
}
