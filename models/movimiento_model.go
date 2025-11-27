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

/*****************************/
/*       FnBD_Mov_Create    */
/***************************/
// @Title Create
// @Description Inserta un nuevo movimiento en la base de datos a través de un Stored Procedure.
// @Accept  application/json
// @Param   movimiento  body  entities.Movimiento  true  "Objeto Movimiento a crear"
// @Success 200 {object} error  null  "Movimiento creado exitosamente"
// @Failure 500 {object} error  "Error al ejecutar el Stored Procedure o al escanear resultados"
// @Router /movimiento [post]
func (m *MovimientoModel) Create(movimiento *entities.Movimiento) error {
	// Definición del Stored Procedure a llamar para crear un movimiento.
	query := `CALL sp_create_movimiento(?, ?, ?, ?, ?, ?)`

	// Variables para almacenar los resultados devueltos por el SP: el ID del movimiento
	// recién creado y el número de filas afectadas (aunque solo se usa idMovimiento).
	var idMovimiento, rowsAffected int64

	// Ejecuta el Stored Procedure en la base de datos. Se utiliza QueryRow porque
	// el SP devuelve valores (idMovimiento y rowsAffected).
	err := database.DB.QueryRow(query,
		movimiento.Fecha,                                            // Parámetro 1: Fecha del movimiento
		movimiento.Monto,                                            // Parámetro 2: Monto del movimiento
		movimiento.Descripcion,                                      // Parámetro 3: Descripción del movimiento
		movimiento.NombreUsuario,                                    // Parámetro 4: Nombre del usuario
		movimiento.NombreConcepto,                                   // Parámetro 5: Nombre del concepto
		movimiento.CorreoFamilia).Scan(&idMovimiento, &rowsAffected) // Parámetro 6: Correo de la familia y escaneo de resultados

	// Verifica si ocurrió algún error durante la ejecución del QueryRow o el escaneo.
	if err != nil {
		// Registra el error en el log con un mensaje descriptivo.
		log.Printf("❌ Error al ejecutar sp_create_movimiento: %v", err)
		return err // Retorna el error.
	}

	// Asigna el ID recién creado devuelto por el SP al objeto 'movimiento'.
	movimiento.IdMovimiento = int(idMovimiento)

	// Registra en el log la creación exitosa del movimiento.
	log.Printf("✅ Movimiento creado - ID: %d", idMovimiento)

	// Retorna nil indicando que la operación fue exitosa.
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

	query := `CALL sp_find_conceptos_by_familia(?, ?)`

	rows, err := database.DB.Query(query, correoFamilia, tipoInt)
	if err != nil {
		log.Printf("❌ Error al ejecutar sp_find_conceptos_by_familia: %v", err)
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
			&concepto.Activo,
			&concepto.DeleteAt,
		)
		if err != nil {
			log.Printf("❌ Error escaneando concepto: %v", err)
			return nil, err
		}
		conceptos = append(conceptos, concepto)
	}

	log.Printf("📋 Conceptos encontrados: %d (tipo: %s, familia: %s)", len(conceptos), tipo, correoFamilia)
	return conceptos, nil
}

// FindByFamiliaAndDate busca movimientos de una familia en una fecha específica
func (m *MovimientoModel) FindByFamiliaAndDate(nombreUsuario string, fecha time.Time) ([]entities.Movimiento, error) {
	query := `CALL sp_find_movimientos_by_familia_date(?, ?)`

	rows, err := database.DB.Query(query, nombreUsuario, fecha.Format("2006-01-02"))
	if err != nil {
		log.Printf("❌ Error al ejecutar sp_find_movimientos_by_familia_date: %v", err)
		return nil, err
	}
	defer rows.Close()

	movimientos := []entities.Movimiento{}
	for rows.Next() {
		var movimiento entities.Movimiento
		var descripcion sql.NullString
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

		if descripcion.Valid {
			movimiento.Descripcion = &descripcion.String
		}

		movimientos = append(movimientos, movimiento)
	}

	log.Printf("📋 Movimientos encontrados: %d para usuario %s en fecha %s",
		len(movimientos), nombreUsuario, fecha.Format("2006-01-02"))
	return movimientos, nil
}

// FindByFamiliaDateAndTipo busca movimientos filtrados por tipo (gasto/ingreso)
func (m *MovimientoModel) FindByFamiliaDateAndTipo(correoFamilia string, fecha time.Time, tipo int8) ([]entities.Movimiento, error) {
	query := `CALL sp_find_movimientos_familia_date_tipo(?, ?, ?)`

	rows, err := database.DB.Query(query, correoFamilia, fecha.Format("2006-01-02"), tipo)
	if err != nil {
		log.Printf("❌ Error al ejecutar sp_find_movimientos_familia_date_tipo: %v", err)
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
			&movimiento.Descripcion,
			&movimiento.NombreUsuario,
			&movimiento.NombreConcepto,
			&movimiento.CorreoFamilia,
		)
		if err != nil {
			log.Printf("❌ Error escaneando movimiento: %v", err)
			continue
		}

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

// FindByID busca un movimiento por su ID
func (m *MovimientoModel) FindByID(idMovimiento int) (*entities.Movimiento, error) {
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

	if descripcion.Valid {
		movimiento.Descripcion = &descripcion.String
	}

	return movimiento, nil
}

// Update actualiza un movimiento existente
func (m *MovimientoModel) Update(movimiento *entities.Movimiento) error {
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
		log.Printf("✅ Movimiento actualizado - ID: %d, Filas afectadas: %d",
			movimiento.IdMovimiento, rowsAffected)
	}

	return nil
}

// Delete realiza soft delete de un movimiento
func (m *MovimientoModel) Delete(idMovimiento int) error {
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
		log.Printf("✅ Movimiento eliminado (soft delete) - ID: %d, Fecha: %s",
			idMovimiento, deletedAt.Format("2006-01-02 15:04:05"))
	}

	return nil
}

// GetTotalesByFamiliaAndDate obtiene totales de ingresos y gastos de una fecha
func (m *MovimientoModel) GetTotalesByFamiliaAndDate(correoFamilia string, fecha time.Time) (float64, float64, error) {
	query := `CALL sp_get_totales_by_familia_and_date(?, ?)`

	var totalGastos, totalIngresos float64
	err := database.DB.QueryRow(query, correoFamilia, fecha.Format("2006-01-02")).Scan(&totalGastos, &totalIngresos)
	if err != nil {
		log.Printf("❌ Error obteniendo totales: %v", err)
		return 0, 0, err
	}

	log.Printf("💰 Totales - Gastos: %.2f, Ingresos: %.2f", totalGastos, totalIngresos)
	return totalGastos, totalIngresos, nil
}

// Agregar estas funciones al archivo models/movimiento.go

// GetTotalesByUsuarioAndPeriod obtiene totales de ingresos y gastos de un usuario en un período
func (m *MovimientoModel) GetTotalesByUsuarioAndPeriod(nombreUsuario, correoFamilia string, inicio, fin time.Time) (float64, float64, error) {
	query := `CALL sp_get_totales_by_usuario_and_period(?, ?, ?, ?)`

	var totalGastos, totalIngresos float64

	err := database.DB.QueryRow(query, nombreUsuario, correoFamilia, inicio, fin).Scan(&totalGastos, &totalIngresos)
	if err != nil {
		log.Printf("❌ Error obteniendo totales por período: %v", err)
		log.Printf("   Usuario: %s, Inicio: %s, Fin: %s", nombreUsuario, inicio.Format("2006-01-02"), fin.Format("2006-01-02"))
		return 0, 0, err
	}

	log.Printf("💰 Totales Usuario %s - Período %s a %s - Gastos: %.2f, Ingresos: %.2f",
		nombreUsuario, inicio.Format("2006-01-02"), fin.Format("2006-01-02"), totalGastos, totalIngresos)

	return totalIngresos, totalGastos, nil
}

// FindByUsuarioAndDate busca movimientos de un usuario en una fecha específica
func (m *MovimientoModel) FindByUsuarioAndDate(nombreUsuario, correoFamilia string, fecha time.Time) ([]entities.Movimiento, error) {
	query := `CALL sp_find_movimientos_by_usuario_and_date(?, ?, ?)`

	rows, err := database.DB.Query(query, nombreUsuario, correoFamilia, fecha.Format("2006-01-02"))
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
func (m *MovimientoModel) FindByUsuarioAndPeriod(nombreUsuario, correoFamilia string, inicio, fin time.Time) ([]entities.Movimiento, error) {
	query := `CALL sp_find_movimientos_by_usuario_and_period(?, ?, ?, ?)`

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
func (m *MovimientoModel) GetTotalByConceptoAndPeriod(nombreConcepto, correoFamilia string, inicio, fin time.Time, usuarioFiltro string) (float64, error) {
	query := `CALL sp_get_total_by_concepto_and_period(?, ?, ?, ?, ?)`

	var total float64
	err := database.DB.QueryRow(query, nombreConcepto, correoFamilia, inicio, fin, usuarioFiltro).Scan(&total)
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
	query := `CALL sp_create_movimiento_automatico(?, ?, ?, ?, ?)`

	var idMovimiento int64
	err := database.DB.QueryRow(query, nombreUsuario, nombreConcepto, correoFamilia, monto, descripcion).Scan(&idMovimiento)
	if err != nil {
		log.Printf("❌ Error creando movimiento automático: %v", err)
		return err
	}

	log.Printf("✅ Movimiento automático creado [ID: %d]: %s - %s - S/ %.2f",
		idMovimiento, nombreUsuario, nombreConcepto, monto)
	return nil
}

// ExisteMovimientoHoy verifica si ya existe un movimiento para hoy
func (m *MovimientoModel) ExisteMovimientoHoy(nombreUsuario, nombreConcepto, correoFamilia string) (bool, error) {
	query := `CALL sp_existe_movimiento_hoy(?, ?, ?)`

	var existe bool
	err := database.DB.QueryRow(query, nombreUsuario, nombreConcepto, correoFamilia).Scan(&existe)
	if err != nil {
		log.Printf("❌ Error verificando existencia de movimiento: %v", err)
		return false, err
	}

	log.Printf("🔍 Movimiento hoy - Usuario: %s, Concepto: %s, Existe: %v",
		nombreUsuario, nombreConcepto, existe)
	return existe, nil
}

// FindByUsuarioDateAndTipo busca movimientos de un usuario específico en una fecha y tipo
func (m *MovimientoModel) FindByUsuarioDateAndTipo(nombreUsuario, correoFamilia string, fecha time.Time, tipo int8) ([]entities.Movimiento, error) {
	query := `CALL sp_find_movimientos_by_usuario_date_and_tipo(?, ?, ?, ?)`

	rows, err := database.DB.Query(query, nombreUsuario, correoFamilia, fecha.Format("2006-01-02"), tipo)
	if err != nil {
		log.Printf("❌ Error en consulta FindByUsuarioDateAndTipo: %v", err)
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
	log.Printf("📋 %s encontrados: %d para usuario %s en fecha %s",
		tipoStr, len(movimientos), nombreUsuario, fecha.Format("2006-01-02"))
	return movimientos, nil
}

// GetTotalesByUsuarioAndDate obtiene totales de un usuario específico en una fecha
func (m *MovimientoModel) GetTotalesByUsuarioAndDate(nombreUsuario, correoFamilia string, fecha time.Time) (float64, float64, error) {
	query := `CALL sp_get_totales_by_usuario_and_date(?, ?, ?)`

	var totalGastos, totalIngresos float64
	err := database.DB.QueryRow(query, nombreUsuario, correoFamilia, fecha.Format("2006-01-02")).Scan(&totalGastos, &totalIngresos)
	if err != nil {
		log.Printf("❌ Error obteniendo totales por usuario: %v", err)
		return 0, 0, err
	}

	log.Printf("💰 Totales usuario %s - Gastos: %.2f, Ingresos: %.2f", nombreUsuario, totalGastos, totalIngresos)
	return totalGastos, totalIngresos, nil
}
