package models

import (
	"arka-code/database"
	"arka-code/entities"
	"database/sql"
	"log"
)

type PersonalizacionModel struct{}

var PersonalizacionModelInstance = &PersonalizacionModel{}

// GetByUsuario obtiene todas las personalizaciones de un usuario
func (m *PersonalizacionModel) GetByUsuario(nombreUsuario string) ([]entities.PersonalizacionConcepto, error) {
	query := `CALL sp_get_personalizaciones_by_usuario(?)`

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

	log.Printf("⚙️  Personalizaciones encontradas: %d para usuario %s", len(personalizaciones), nombreUsuario)
	return personalizaciones, nil
}

// GetConsumoActual calcula el consumo actual de un concepto para un usuario
func (m *PersonalizacionModel) GetConsumoActual(nombreUsuario, nombreConcepto, periodo string) (float64, error) {
	query := `CALL sp_get_consumo_actual(?, ?, ?)`

	var consumo float64
	err := database.DB.QueryRow(query, nombreUsuario, nombreConcepto, periodo).Scan(&consumo)
	if err != nil {
		log.Printf("❌ Error calculando consumo: %v", err)
		return 0, err
	}

	log.Printf("📊 Consumo actual - Usuario: %s, Concepto: %s, Período: %s, Monto: %.2f",
		nombreUsuario, nombreConcepto, periodo, consumo)
	return consumo, nil
}

// GetByID obtiene una personalización por su ID
func (m *PersonalizacionModel) GetByID(id int) (*entities.PersonalizacionConcepto, error) {
	query := `CALL sp_get_personalizacion_by_id(?)`

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
		log.Printf("⚠️  Personalización ID %d no encontrada", id)
		return nil, nil
	}

	if err != nil {
		log.Printf("❌ Error obteniendo personalización ID %d: %v", id, err)
		return nil, err
	}

	log.Printf("✅ Personalización ID %d obtenida: Usuario=%s, Concepto=%s",
		id, p.NombreUsuario, p.NombreConcepto)
	return p, nil
}

// Create crea una nueva personalización
func (m *PersonalizacionModel) Create(p *entities.PersonalizacionConcepto) error {
	query := `CALL sp_create_personalizacion(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	var idPersonalizacion int64
	err := database.DB.QueryRow(query,
		p.LimiteGasto,
		p.Activo,
		p.MontoPlanificado,
		p.TipoPeriodoPlanificado,
		p.TipoPeriodoLimite,
		p.DiaPeriodoPlanificado,
		p.Notificacion,
		p.NombreUsuario,
		p.NombreConcepto,
		p.CorreoFamilia,
	).Scan(&idPersonalizacion)

	if err != nil {
		log.Printf("❌ Error creando personalización: %v", err)
		log.Printf("   Usuario: %s, Concepto: %s", p.NombreUsuario, p.NombreConcepto)
		return err
	}

	p.IdPersonalizacion = int(idPersonalizacion)

	// log.Printf("Personalización creada - ID: %d, Usuario: %s, Concepto: %s, Límite: %.2f",
	// idPersonalizacion, p.NombreUsuario, p.NombreConcepto, p.LimiteGasto)
	return nil
}

// Update actualiza una personalización existente
func (m *PersonalizacionModel) Update(p *entities.PersonalizacionConcepto) error {
	query := `CALL sp_update_personalizacion(?, ?, ?, ?, ?, ?, ?)`

	var filasAfectadas int64
	err := database.DB.QueryRow(query,
		p.IdPersonalizacion,
		p.LimiteGasto,
		p.MontoPlanificado,
		p.TipoPeriodoPlanificado,
		p.TipoPeriodoLimite,
		p.DiaPeriodoPlanificado,
		p.Activo,
	).Scan(&filasAfectadas)

	if err != nil {
		log.Printf("❌ Error actualizando personalización: %v", err)
		log.Printf("   ID: %d", p.IdPersonalizacion)
		return err
	}

	if filasAfectadas == 0 {
		log.Printf("⚠️  Personalización ID %d no encontrada o ya eliminada", p.IdPersonalizacion)
	}

	log.Printf("✅ Personalización actualizada - ID: %d, Usuario: %s, Concepto: %s",
		p.IdPersonalizacion, p.NombreUsuario, p.NombreConcepto)
	return nil
}

// Delete elimina una personalización (soft delete)
func (m *PersonalizacionModel) Delete(id int) error {
	query := `CALL sp_delete_personalizacion(?)`

	var filasAfectadas int64
	err := database.DB.QueryRow(query, id).Scan(&filasAfectadas)
	if err != nil {
		log.Printf("❌ Error eliminando personalización: %v", err)
		log.Printf("   ID: %d", id)
		return err
	}

	if filasAfectadas == 0 {
		log.Printf("⚠️  Personalización ID %d no encontrada o ya eliminada", id)
	}

	log.Printf("✅ Personalización eliminada (soft delete) - ID: %d", id)
	return nil
}

// GetUsuarioQueAsigno obtiene el nombre del usuario que asignó un límite
func (m *PersonalizacionModel) GetUsuarioQueAsigno(idPersonalizacion int) (string, error) {
	query := `CALL sp_get_usuario_que_asigno(?)`

	var nombrePersonal string
	err := database.DB.QueryRow(query, idPersonalizacion).Scan(&nombrePersonal)
	if err != nil {
		log.Printf("⚠️  No se pudo obtener usuario que asignó personalización ID %d: %v", idPersonalizacion, err)
		return "Desconocido", nil
	}

	log.Printf("👤 Usuario que asignó personalización ID %d: %s", idPersonalizacion, nombrePersonal)
	return nombrePersonal, nil
}

// FindByUsuarioAndConcepto obtiene la personalización de un usuario para un concepto
func (m *PersonalizacionModel) FindByUsuarioAndConcepto(nombreUsuario, nombreConcepto, correoFamilia string) (*entities.PersonalizacionConcepto, error) {
	query := `CALL sp_find_personalizacion_by_usuario_and_concepto(?, ?, ?)`

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
		log.Printf("ℹ️  Personalización no encontrada - Usuario: %s, Concepto: %s", nombreUsuario, nombreConcepto)
		return nil, nil
	}

	if err != nil {
		log.Printf("❌ Error buscando personalización: %v", err)
		log.Printf("   Usuario: %s, Concepto: %s", nombreUsuario, nombreConcepto)
		return nil, err
	}

	log.Printf("✅ Personalización encontrada - ID: %d, Usuario: %s, Concepto: %s, Activo: %v",
		p.IdPersonalizacion, p.NombreUsuario, p.NombreConcepto, p.Activo)
	return p, nil
}

func (m *PersonalizacionModel) FindByPeriodoDiario() ([]entities.PersonalizacionConcepto, error) {
	query := `CALL sp_find_personalizaciones_by_periodo_diario()`

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
	query := `CALL sp_find_personalizaciones_by_dia_planificado(?)`

	rows, err := database.DB.Query(query, dia)
	if err != nil {
		log.Printf("❌ Error en consulta FindByDiaPlanificado: %v", err)
		log.Printf("   Día: %d", dia)
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
	query := `CALL sp_deshabilitar_personalizacion_para_usuario(?, ?, ?)`

	var idPersonalizacion int64
	var operacion string

	err := database.DB.QueryRow(query, nombreConcepto, correoFamilia, nombreUsuario).Scan(&idPersonalizacion, &operacion)
	if err != nil {
		log.Printf("❌ Error deshabilitando personalización: %v", err)
		log.Printf("   Usuario: %s, Concepto: %s", nombreUsuario, nombreConcepto)
		return err
	}

	if operacion == "created" {
		log.Printf("✅ Personalización creada como inactiva - ID: %d, Usuario: %s, Concepto: %s",
			idPersonalizacion, nombreUsuario, nombreConcepto)
	} else {
		log.Printf("✅ Personalización deshabilitada - ID: %d, Usuario: %s, Concepto: %s",
			idPersonalizacion, nombreUsuario, nombreConcepto)
	}

	return nil
}

// HabilitarParaUsuario activa un concepto para un usuario específico
func (m *PersonalizacionModel) HabilitarParaUsuario(nombreConcepto, correoFamilia, nombreUsuario string) error {
	query := `CALL sp_habilitar_personalizacion_para_usuario(?, ?, ?)`

	var idPersonalizacion int64
	var operacion string

	err := database.DB.QueryRow(query, nombreConcepto, correoFamilia, nombreUsuario).Scan(&idPersonalizacion, &operacion)
	if err != nil {
		log.Printf("❌ Error habilitando personalización: %v", err)
		log.Printf("   Usuario: %s, Concepto: %s", nombreUsuario, nombreConcepto)
		return err
	}

	if operacion == "created" {
		log.Printf("✅ Personalización creada como activa - ID: %d, Usuario: %s, Concepto: %s",
			idPersonalizacion, nombreUsuario, nombreConcepto)
	} else {
		log.Printf("✅ Personalización habilitada - ID: %d, Usuario: %s, Concepto: %s",
			idPersonalizacion, nombreUsuario, nombreConcepto)
	}

	return nil
}

// ToggleActivoParaUsuario cambia el estado activo de un concepto para un usuario
func (m *PersonalizacionModel) ToggleActivoParaUsuario(nombreConcepto, correoFamilia, nombreUsuario string) error {
	query := `CALL sp_toggle_activo_para_usuario(?, ?, ?)`

	var idPersonalizacion int64
	var operacion string
	var nuevoEstado bool
	var estadoAnterior bool

	err := database.DB.QueryRow(query, nombreConcepto, correoFamilia, nombreUsuario).Scan(
		&idPersonalizacion,
		&operacion,
		&nuevoEstado,
		&estadoAnterior,
	)
	if err != nil {
		log.Printf("❌ Error haciendo toggle de personalización: %v", err)
		log.Printf("   Usuario: %s, Concepto: %s", nombreUsuario, nombreConcepto)
		return err
	}

	estadoTexto := "habilitado"
	emoji := "✅"
	if !nuevoEstado {
		estadoTexto = "deshabilitado"
		emoji = "❌"
	}

	if operacion == "created" {
		log.Printf("%s Concepto %s (nueva personalización) - ID: %d, Usuario: %s, Concepto: %s",
			emoji, estadoTexto, idPersonalizacion, nombreUsuario, nombreConcepto)
	} else {
		log.Printf("%s Concepto cambiado a %s - ID: %d, Usuario: %s, Concepto: %s (antes: %v)",
			emoji, estadoTexto, idPersonalizacion, nombreUsuario, nombreConcepto, estadoAnterior)
	}

	return nil
}
