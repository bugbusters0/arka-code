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
