package models

import (
	"arka-code/database"
	"arka-code/entities"
	"database/sql"
	"log"
	"time"
)

type PersonalizacionModel struct{}

var PersonalizacionModelInstance = &PersonalizacionModel{}

// GetByUsuario obtiene todas las personalizaciones de un usuario
func (m *PersonalizacionModel) GetByUsuario(nombreUsuario string) ([]entities.PersonalizacionConcepto, error) {
	query := `SELECT p.idPersonalizacion, p.limiteGasto, p.activo, p.montoPlanificado,
	                 p.tipoPeriodoPlanificado, p.tipoPeriodoLimite, p.diaPeriodoPlanificado,
	                 p.notificacion, p.nombreUsuario, p.nombreConcepto, p.correoFamilia, p.delete_at
	          FROM personalizacionconcepto p
	          WHERE p.nombreUsuario = ? AND p.delete_at IS NULL
	          ORDER BY p.nombreConcepto`

	rows, err := database.DB.Query(query, nombreUsuario)
	if err != nil {
		log.Printf("❌ Error obteniendo personalizaciones: %v", err)
		return nil, err
	}
	defer rows.Close()

	personalizaciones := []entities.PersonalizacionConcepto{}
	for rows.Next() {
		var p entities.PersonalizacionConcepto
		err := rows.Scan(
			&p.IdPersonalizacion,
			&p.LimiteGasto,
			&p.Activo,
			&p.MontoPlanificado,
			&p.TipoPeriodoPlanificado,
			&p.TipoPeriodoLimite,
			&p.DiaPeriodoPlanificado,
			&p.Notificacion,
			&p.NombreUsuario,
			&p.NombreConcepto,
			&p.CorreoFamilia,
			&p.DeleteAt,
		)
		if err != nil {
			log.Printf("❌ Error escaneando personalización: %v", err)
			continue
		}
		personalizaciones = append(personalizaciones, p)
	}

	return personalizaciones, nil
}

// GetConsumoActual calcula el consumo actual de un concepto para un usuario
func (m *PersonalizacionModel) GetConsumoActual(nombreUsuario, nombreConcepto, periodo string) (float64, error) {
	// Determinar el rango de fechas según el periodo
	var fechaInicio time.Time
	now := time.Now()

	switch periodo {
	case "diario":
		fechaInicio = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	case "semanal":
		// Inicio de la semana (Lunes)
		diasDesdeInicio := int(now.Weekday())
		if diasDesdeInicio == 0 {
			diasDesdeInicio = 7 // Domingo es el último día
		}
		fechaInicio = now.AddDate(0, 0, -(diasDesdeInicio - 1))
		fechaInicio = time.Date(fechaInicio.Year(), fechaInicio.Month(), fechaInicio.Day(), 0, 0, 0, 0, fechaInicio.Location())
	case "mensual":
		fechaInicio = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	case "anual":
		fechaInicio = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	default:
		fechaInicio = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	}

	query := `SELECT COALESCE(SUM(m.monto), 0)
	          FROM movimiento m
	          WHERE m.nombreUsuario = ? 
	            AND m.nombreConcepto = ?
	            AND m.fecha >= ?
	            AND m.delete_at IS NULL`

	var consumo float64
	err := database.DB.QueryRow(query, nombreUsuario, nombreConcepto, fechaInicio).Scan(&consumo)
	if err != nil {
		log.Printf("❌ Error calculando consumo: %v", err)
		return 0, err
	}

	return consumo, nil
}

// GetByID obtiene una personalización por su ID
func (m *PersonalizacionModel) GetByID(id int) (*entities.PersonalizacionConcepto, error) {
	query := `SELECT idPersonalizacion, limiteGasto, activo, montoPlanificado,
	                 tipoPeriodoPlanificado, tipoPeriodoLimite, diaPeriodoPlanificado,
	                 notificacion, nombreUsuario, nombreConcepto, correoFamilia, delete_at
	          FROM personalizacionconcepto
	          WHERE idPersonalizacion = ? AND delete_at IS NULL`

	p := &entities.PersonalizacionConcepto{}
	err := database.DB.QueryRow(query, id).Scan(
		&p.IdPersonalizacion,
		&p.LimiteGasto,
		&p.Activo,
		&p.MontoPlanificado,
		&p.TipoPeriodoPlanificado,
		&p.TipoPeriodoLimite,
		&p.DiaPeriodoPlanificado,
		&p.Notificacion,
		&p.NombreUsuario,
		&p.NombreConcepto,
		&p.CorreoFamilia,
		&p.DeleteAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Printf("❌ Error obteniendo personalización ID %d: %v", id, err)
		return nil, err
	}

	return p, nil
}

// Create crea una nueva personalización
func (m *PersonalizacionModel) Create(p *entities.PersonalizacionConcepto) error {
	query := `INSERT INTO personalizacionconcepto 
	          (limiteGasto, activo, montoPlanificado, tipoPeriodoPlanificado,
	           tipoPeriodoLimite, diaPeriodoPlanificado, notificacion,
	           nombreUsuario, nombreConcepto, correoFamilia)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := database.DB.Exec(query,
		p.LimiteGasto,
		p.Activo,
		p.MontoPlanificado,
		p.TipoPeriodoPlanificado,
		p.TipoPeriodoLimite,
		p.DiaPeriodoPlanificado,
		p.Notificacion,
		p.NombreUsuario,
		p.NombreConcepto,
		p.CorreoFamilia)

	if err != nil {
		log.Printf("❌ Error creando personalización: %v", err)
		return err
	}

	id, _ := result.LastInsertId()
	p.IdPersonalizacion = int(id)

	log.Printf("✅ Personalización creada - ID: %d", id)
	return nil
}

// Update actualiza una personalización existente
func (m *PersonalizacionModel) Update(p *entities.PersonalizacionConcepto) error {
	query := `UPDATE personalizacionconcepto 
	          SET limiteGasto = ?, montoPlanificado = ?,
	              tipoPeriodoPlanificado = ?, tipoPeriodoLimite = ?,
	              diaPeriodoPlanificado = ?, activo = ?
	          WHERE idPersonalizacion = ?`

	result, err := database.DB.Exec(query,
		p.LimiteGasto,
		p.MontoPlanificado,
		p.TipoPeriodoPlanificado,
		p.TipoPeriodoLimite,
		p.DiaPeriodoPlanificado,
		p.Activo,
		p.IdPersonalizacion)

	if err != nil {
		log.Printf("❌ Error actualizando personalización: %v", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("✅ Personalización actualizada - ID: %d, Filas afectadas: %d", p.IdPersonalizacion, rowsAffected)
	return nil
}

// Delete elimina una personalización (soft delete)
func (m *PersonalizacionModel) Delete(id int) error {
	query := `UPDATE personalizacionconcepto SET delete_at = ? WHERE idPersonalizacion = ?`

	result, err := database.DB.Exec(query, time.Now(), id)
	if err != nil {
		log.Printf("❌ Error eliminando personalización: %v", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("✅ Personalización eliminada - ID: %d, Filas afectadas: %d", id, rowsAffected)
	return nil
}

// GetUsuarioQueAsigno obtiene el nombre del usuario que asignó un límite
func (m *PersonalizacionModel) GetUsuarioQueAsigno(idPersonalizacion int) (string, error) {
	// Por ahora, asumimos que el que creó la personalización es el que la asignó
	// Si necesitas un campo adicional en la BD, lo agregamos después
	query := `SELECT u.nombrePersonal 
	          FROM personalizacionconcepto p
	          INNER JOIN usuario u ON p.nombreUsuario = u.nombreUsuario
	          WHERE p.idPersonalizacion = ?`

	var nombrePersonal string
	err := database.DB.QueryRow(query, idPersonalizacion).Scan(&nombrePersonal)
	if err != nil {
		return "Desconocido", nil
	}

	return nombrePersonal, nil
}

// FindByUsuarioAndConcepto obtiene la personalización de un usuario para un concepto
func (m *PersonalizacionModel) FindByUsuarioAndConcepto(nombreUsuario, nombreConcepto, correoFamilia string) (*entities.PersonalizacionConcepto, error) {
	query := `SELECT idPersonalizacion, montoPlanificado, tipoPeriodoPlanificado, 
	                 diaPeriodoPlanificado, limiteGasto, tipoPeriodoLimite, 
	                 diaPeriodoLimite, notificacion, activo,
	                 nombreUsuario, nombreConcepto, correoFamilia
	          FROM personalizacionconcepto
	          WHERE nombreUsuario = ? 
	            AND nombreConcepto = ? 
	            AND correoFamilia = ?
	            AND delete_at IS NULL`

	p := &entities.PersonalizacionConcepto{}
	err := database.DB.QueryRow(query, nombreUsuario, nombreConcepto, correoFamilia).Scan(
		&p.IdPersonalizacion,
		&p.MontoPlanificado,
		&p.TipoPeriodoPlanificado,
		&p.DiaPeriodoPlanificado,
		&p.LimiteGasto,
		&p.TipoPeriodoLimite,
		&p.DiaPeriodoLimite,
		&p.Notificacion,
		&p.Activo,
		&p.NombreUsuario,
		&p.NombreConcepto,
		&p.CorreoFamilia,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		log.Printf("❌ Error buscando personalización: %v", err)
		return nil, err
	}

	return p, nil
}

func (m *PersonalizacionModel) FindByPeriodoDiario() ([]entities.PersonalizacionConcepto, error) {
	query := `SELECT p.idPersonalizacion, p.montoPlanificado, p.tipoPeriodoPlanificado, 
                     p.diaPeriodoPlanificado, p.limiteGasto, p.tipoPeriodoLimite, 
                     p.diaPeriodoLimite, p.notificacion, p.activo,
                     p.nombreUsuario, p.nombreConcepto, p.correoFamilia
              FROM personalizacionconcepto p
              WHERE p.activo = 1 
                AND p.montoPlanificado IS NOT NULL 
                AND p.montoPlanificado > 0
                AND p.tipoPeriodoPlanificado = 'diario'
                AND p.delete_at IS NULL`

	rows, err := database.DB.Query(query)
	if err != nil {
		log.Printf("❌ Error en consulta FindByPeriodoDiario: %v", err)
		return nil, err
	}
	defer rows.Close()

	personalizaciones := []entities.PersonalizacionConcepto{}
	for rows.Next() {
		var p entities.PersonalizacionConcepto

		err := rows.Scan(
			&p.IdPersonalizacion,
			&p.MontoPlanificado,
			&p.TipoPeriodoPlanificado,
			&p.DiaPeriodoPlanificado,
			&p.LimiteGasto,
			&p.TipoPeriodoLimite,
			&p.DiaPeriodoLimite,
			&p.Notificacion,
			&p.Activo,
			&p.NombreUsuario,
			&p.NombreConcepto,
			&p.CorreoFamilia,
		)
		if err != nil {
			log.Printf("❌ Error escaneando personalización diaria: %v", err)
			continue
		}
		personalizaciones = append(personalizaciones, p)
	}

	log.Printf("📋 Personalizaciones diarias encontradas: %d", len(personalizaciones))
	return personalizaciones, nil
}

// FindByDiaPlanificado obtiene personalizaciones para un día específico (modificado)
func (m *PersonalizacionModel) FindByDiaPlanificado(dia int8) ([]entities.PersonalizacionConcepto, error) {
	query := `SELECT p.idPersonalizacion, p.montoPlanificado, p.tipoPeriodoPlanificado, 
                     p.diaPeriodoPlanificado, p.limiteGasto, p.tipoPeriodoLimite, 
                     p.diaPeriodoLimite, p.notificacion, p.activo,
                     p.nombreUsuario, p.nombreConcepto, p.correoFamilia
              FROM personalizacionconcepto p
              WHERE p.activo = 1 
                AND p.montoPlanificado IS NOT NULL 
                AND p.montoPlanificado > 0
                AND p.diaPeriodoPlanificado = ?
                AND p.tipoPeriodoPlanificado IN ('mensual', 'quincenal')
                AND p.delete_at IS NULL`

	rows, err := database.DB.Query(query, dia)
	if err != nil {
		log.Printf("❌ Error en consulta FindByDiaPlanificado: %v", err)
		return nil, err
	}
	defer rows.Close()

	personalizaciones := []entities.PersonalizacionConcepto{}
	for rows.Next() {
		var p entities.PersonalizacionConcepto

		err := rows.Scan(
			&p.IdPersonalizacion,
			&p.MontoPlanificado,
			&p.TipoPeriodoPlanificado,
			&p.DiaPeriodoPlanificado,
			&p.LimiteGasto,
			&p.TipoPeriodoLimite,
			&p.DiaPeriodoLimite,
			&p.Notificacion,
			&p.Activo,
			&p.NombreUsuario,
			&p.NombreConcepto,
			&p.CorreoFamilia,
		)
		if err != nil {
			log.Printf("❌ Error escaneando personalización: %v", err)
			continue
		}
		personalizaciones = append(personalizaciones, p)
	}

	log.Printf("📋 Personalizaciones encontradas para día %d: %d", dia, len(personalizaciones))
	return personalizaciones, nil
}

// DeshabilitarParaUsuario desactiva un concepto para un usuario específico
func (m *PersonalizacionModel) DeshabilitarParaUsuario(nombreConcepto, correoFamilia, nombreUsuario string) error {
	// Primero verificar si existe la personalización
	queryCheck := `SELECT idPersonalizacion FROM personalizacionconcepto 
                   WHERE nombreConcepto = ? AND correoFamilia = ? AND nombreUsuario = ?`

	var idPersonalizacion int
	err := database.DB.QueryRow(queryCheck, nombreConcepto, correoFamilia, nombreUsuario).Scan(&idPersonalizacion)

	if err == sql.ErrNoRows {
		// No existe personalización, crear una con activo = false
		queryInsert := `INSERT INTO personalizacionconcepto 
                       (limiteGasto, activo, montoPlanificado, tipoPeriodoPlanificado, 
                        tipoPeriodoLimite, diaPeriodoPlanificado, notificacion,
                        nombreUsuario, nombreConcepto, correoFamilia)
                       VALUES (NULL, false, NULL, NULL, NULL, NULL, false, ?, ?, ?)`

		_, err = database.DB.Exec(queryInsert, nombreUsuario, nombreConcepto, correoFamilia)
		if err != nil {
			log.Printf("❌ Error creando personalización inactiva: %v", err)
			return err
		}
		log.Printf("✅ Personalización creada como inactiva para usuario %s", nombreUsuario)
		return nil
	} else if err != nil {
		log.Printf("❌ Error verificando personalización: %v", err)
		return err
	}

	// Ya existe personalización, desactivarla
	queryUpdate := `UPDATE personalizacionconcepto 
                   SET activo = false 
                   WHERE nombreConcepto = ? AND correoFamilia = ? AND nombreUsuario = ?`

	result, err := database.DB.Exec(queryUpdate, nombreConcepto, correoFamilia, nombreUsuario)
	if err != nil {
		log.Printf("❌ Error deshabilitando concepto: %v", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("✅ Concepto deshabilitado para usuario %s - Filas afectadas: %d", nombreUsuario, rowsAffected)
	return nil
}

// HabilitarParaUsuario activa un concepto para un usuario específico
func (m *PersonalizacionModel) HabilitarParaUsuario(nombreConcepto, correoFamilia, nombreUsuario string) error {
	// Primero verificar si existe la personalización
	queryCheck := `SELECT idPersonalizacion FROM personalizacionconcepto 
                   WHERE nombreConcepto = ? AND correoFamilia = ? AND nombreUsuario = ?`

	var idPersonalizacion int
	err := database.DB.QueryRow(queryCheck, nombreConcepto, correoFamilia, nombreUsuario).Scan(&idPersonalizacion)

	if err == sql.ErrNoRows {
		// No existe personalización, crear una con activo = true
		queryInsert := `INSERT INTO personalizacionconcepto 
                       (limiteGasto, activo, montoPlanificado, tipoPeriodoPlanificado, 
                        tipoPeriodoLimite, diaPeriodoPlanificado, notificacion,
                        nombreUsuario, nombreConcepto, correoFamilia)
                       VALUES (NULL, true, NULL, NULL, NULL, NULL, false, ?, ?, ?)`

		_, err = database.DB.Exec(queryInsert, nombreUsuario, nombreConcepto, correoFamilia)
		if err != nil {
			log.Printf("❌ Error creando personalización activa: %v", err)
			return err
		}
		log.Printf("✅ Personalización creada como activa para usuario %s", nombreUsuario)
		return nil
	} else if err != nil {
		log.Printf("❌ Error verificando personalización: %v", err)
		return err
	}

	// Ya existe personalización, activarla
	queryUpdate := `UPDATE personalizacionconcepto 
                   SET activo = true 
                   WHERE nombreConcepto = ? AND correoFamilia = ? AND nombreUsuario = ?`

	result, err := database.DB.Exec(queryUpdate, nombreConcepto, correoFamilia, nombreUsuario)
	if err != nil {
		log.Printf("❌ Error habilitando concepto: %v", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("✅ Concepto habilitado para usuario %s - Filas afectadas: %d", nombreUsuario, rowsAffected)
	return nil
}

// ToggleActivoParaUsuario cambia el estado activo de un concepto para un usuario
func (m *PersonalizacionModel) ToggleActivoParaUsuario(nombreConcepto, correoFamilia, nombreUsuario string) error {
	// Primero obtener el estado actual
	queryEstado := `SELECT COALESCE(activo, true) as activo
                   FROM personalizacionconcepto 
                   WHERE nombreConcepto = ? AND correoFamilia = ? AND nombreUsuario = ?`

	var activoActual bool
	err := database.DB.QueryRow(queryEstado, nombreConcepto, correoFamilia, nombreUsuario).Scan(&activoActual)

	if err == sql.ErrNoRows {
		// No existe personalización, crear una con activo = false (toggle de true a false)
		queryInsert := `INSERT INTO personalizacionconcepto 
                       (limiteGasto, activo, montoPlanificado, tipoPeriodoPlanificado, 
                        tipoPeriodoLimite, diaPeriodoPlanificado, notificacion,
                        nombreUsuario, nombreConcepto, correoFamilia)
                       VALUES (NULL, false, NULL, NULL, NULL, NULL, false, ?, ?, ?)`

		_, err = database.DB.Exec(queryInsert, nombreUsuario, nombreConcepto, correoFamilia)
		if err != nil {
			return err
		}
		log.Printf("✅ Concepto deshabilitado (nueva personalización) para usuario %s", nombreUsuario)
		return nil
	} else if err != nil {
		return err
	}

	// Cambiar el estado
	if activoActual {
		return m.DeshabilitarParaUsuario(nombreConcepto, correoFamilia, nombreUsuario)
	} else {
		return m.HabilitarParaUsuario(nombreConcepto, correoFamilia, nombreUsuario)
	}
}
