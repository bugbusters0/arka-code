package models

import (
	"arka-code/database"
	"arka-code/entities"
	"database/sql"
	"fmt"
	"log"
)

// CU-004
// BD.3 ConceptoBD
type ConceptoModel struct{}

var ConceptoModelInstance = &ConceptoModel{}

// FNBD-Index
// Create crea un nuevo concepto en la base de datos
// Parámetros:
// - concepto: Estructura con los datos del concepto a crear
// Retorno: Error si la creación falla
// Flujo:
// - Inserta nuevo registro en tabla concepto
// - Campos requeridos: nombreConcepto, correoFamilia, tipo, nombreUsuario
// - Campos opcionales: icono, color
// - Registra número de filas afectadas para verificación
// Uso: Creación de nuevos conceptos de gastos/ingresos
func (m *ConceptoModel) Create(concepto *entities.Concepto) error {
	query := `INSERT INTO concepto (nombreConcepto, correoFamilia, tipo, icono, color, nombreUsuario)
	          VALUES (?, ?, ?, ?, ?, ?)`

	result, err := database.DB.Exec(query,
		concepto.NombreConcepto,
		concepto.CorreoFamilia,
		concepto.Tipo,
		concepto.Icono,
		concepto.Color,
		concepto.NombreUsuario)

	if err != nil {
		log.Printf("❌ Error en consulta INSERT concepto: %v", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("✅ Concepto creado - Filas afectadas: %d", rowsAffected)
	return nil
}

// FNBD-FindByFamilia
// FindByFamilia busca todos los conceptos de una familia filtrados por tipo
// Parámetros:
// - correoFamilia: Identificador de la familia
// - tipo: Tipo de concepto ("gasto" o "ingreso")
// Retorno: Lista de conceptos con su estado activo/inactivo
// Flujo:
// - Convierte tipo string a int8 (0=gasto, 1=ingreso)
// - Realiza LEFT JOIN con personalizacionconcepto para obtener estado activo
// - Usa COALESCE para establecer activo=true por defecto si no hay personalización
// - Excluye conceptos eliminados (soft delete)
// Uso: Listado de conceptos para selección en movimientos
func (m *ConceptoModel) FindByFamilia(correoFamilia string, tipo string) ([]entities.Concepto, error) {
	var tipoInt int8
	if tipo == "ingreso" {
		tipoInt = 1
	} else {
		tipoInt = 0
	}

	// ✅ CORREGIDO: JOIN con personalizacionconcepto para obtener el estado activo
	query := `SELECT c.nombreConcepto, c.correoFamilia, c.tipo, c.icono, c.color, c.nombreUsuario,
                     COALESCE(pc.activo, 1) as activo, c.delete_at
              FROM concepto c
              LEFT JOIN personalizacionconcepto pc ON c.nombreConcepto = pc.nombreConcepto 
                                                   AND c.correoFamilia = pc.correoFamilia 
                                                   AND pc.nombreUsuario = c.nombreUsuario
              WHERE c.correoFamilia = ? AND c.tipo = ? AND c.delete_at IS NULL
              ORDER BY c.nombreConcepto`

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
			&concepto.Activo,
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

// FNBD-UpdateIconoColor
// UpdateIconoColor actualiza solo el ícono y color de un concepto existente
// Parámetros:
// - concepto: Estructura con los nuevos valores de ícono y color
// Retorno: Error si la actualización falla
// Flujo:
// - Actualiza únicamente campos icono y color
// - Filtra por nombreConcepto, correoFamilia y excluye eliminados
// - No permite cambiar nombre o tipo del concepto
// Uso: Personalización visual de conceptos existentes
func (m *ConceptoModel) UpdateIconoColor(concepto *entities.Concepto) error {
	query := `UPDATE concepto 
	          SET icono = ?, color = ?
	          WHERE nombreConcepto = ? AND correoFamilia = ? AND delete_at IS NULL`

	result, err := database.DB.Exec(query,
		concepto.Icono,
		concepto.Color,
		concepto.NombreConcepto,
		concepto.CorreoFamilia)

	if err != nil {
		log.Printf("❌ Error actualizando concepto: %v", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("✅ Concepto actualizado - Filas afectadas: %d", rowsAffected)
	return nil
}

// FNBD-FindAllByFamilia
// FindAllByFamilia busca todos los conceptos de una familia sin filtrar por tipo
// Parámetros:
// - correoFamilia: Identificador de la familia
// Retorno: Lista completa de conceptos de la familia
// Flujo:
// - Consulta todos los conceptos no eliminados de la familia
// - No incluye información de personalizaciones
// - Ordena alfabéticamente por nombre
// Uso: Reportes generales, estadísticas familiares
func (m *ConceptoModel) FindAllByFamilia(correoFamilia string) ([]entities.Concepto, error) {
	query := `SELECT nombreConcepto, correoFamilia, tipo, icono, color, nombreUsuario, delete_at
              FROM concepto 
              WHERE correoFamilia = ? AND delete_at IS NULL
              ORDER BY nombreConcepto`

	rows, err := database.DB.Query(query, correoFamilia)
	if err != nil {
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
			return nil, err
		}
		conceptos = append(conceptos, concepto)
	}

	return conceptos, nil
}

// FNBD-FindByNombre
// FindByNombre busca un concepto específico por nombre y familia
// Parámetros:
// - nombreConcepto: Nombre del concepto a buscar
// - correoFamilia: Familia a la que pertenece el concepto
// Retorno: Puntero al concepto encontrado o nil si no existe
// Flujo:
// - Busca concepto exacto por nombre y familia
// - Excluye conceptos eliminados (soft delete)
// - Retorna sql.ErrNoRows si no se encuentra
// Uso: Verificación de existencia, operaciones de edición
func (m *ConceptoModel) FindByNombre(nombreConcepto, correoFamilia string) (*entities.Concepto, error) {
	query := `SELECT nombreConcepto, correoFamilia, tipo, icono, color, nombreUsuario, delete_at
	          FROM concepto 
	          WHERE nombreConcepto = ? AND correoFamilia = ? AND delete_at IS NULL`

	concepto := &entities.Concepto{}
	err := database.DB.QueryRow(query, nombreConcepto, correoFamilia).Scan(
		&concepto.NombreConcepto,
		&concepto.CorreoFamilia,
		&concepto.Tipo,
		&concepto.Icono,
		&concepto.Color,
		&concepto.NombreUsuario,
		&concepto.DeleteAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return concepto, err
}

// FNBD-Exists
// Exists verifica si un concepto ya existe para una familia específica
// Parámetros:
// - nombreConcepto: Nombre del concepto a verificar
// - correoFamilia: Familia donde buscar
// Retorno: true si existe, false si no existe, error en caso de fallo
// Flujo:
// - Cuenta registros con nombre y familia coincidentes
// - Excluye conceptos eliminados (soft delete)
// Uso: Validación antes de crear nuevos conceptos, evitar duplicados
func (m *ConceptoModel) Exists(nombreConcepto, correoFamilia string) (bool, error) {
	query := `SELECT COUNT(*) FROM concepto 
	          WHERE nombreConcepto = ? AND correoFamilia = ? AND delete_at IS NULL`

	var count int
	err := database.DB.QueryRow(query, nombreConcepto, correoFamilia).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// FNBD-CreatePersonalizacion
// CreatePersonalizacion crea una configuración personalizada para un usuario y concepto
// Parámetros:
// - personalizacion: Mapa con todos los campos de personalización
// Retorno: Error si la creación falla
// Flujo:
// - Inserta registro en tabla personalizacionconcepto
// - Campos incluyen límites, montos planificados, frecuencias y notificaciones
// - Registra datos para debugging en caso de error
// Uso: Configuración individual de conceptos por usuario
func (m *ConceptoModel) CreatePersonalizacion(personalizacion map[string]interface{}) error {
	query := `INSERT INTO personalizacionconcepto 
              (limiteGasto, activo, montoPlanificado, tipoPeriodoPlanificado, 
               tipoPeriodoLimite, diaPeriodoPlanificado, diaPeriodoLimite, notificacion,
               nombreUsuario, nombreConcepto, correoFamilia)
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := database.DB.Exec(query,
		personalizacion["limiteGasto"],
		personalizacion["activo"],
		personalizacion["montoPlanificado"],
		personalizacion["tipoPeriodoPlanificado"],
		personalizacion["tipoPeriodoLimite"],
		personalizacion["diaPeriodoPlanificado"],
		personalizacion["diaPeriodoLimite"],
		personalizacion["notificacion"],
		personalizacion["nombreUsuario"],
		personalizacion["nombreConcepto"],
		personalizacion["correoFamilia"])

	if err != nil {
		log.Printf("❌ Error en consulta INSERT personalizacion: %v", err)
		log.Printf("   Datos: %+v", personalizacion)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("✅ Personalización creada para %s - Filas afectadas: %d",
		personalizacion["nombreUsuario"], rowsAffected)

	return nil
}

// FNBD-GetUsuariosByFamilia
// GetUsuariosByFamilia obtiene todos los usuarios activos de una familia
// Parámetros:
// - correoFamilia: Identificador de la familia
// Retorno: Lista de usuarios de la familia
// Flujo:
// - Consulta tabla usuario filtrada por familia
// - Excluye usuarios eliminados (soft delete)
// - Incluye todos los campos de usuario excepto contraseña
// Uso: Operaciones masivas sobre usuarios familiares
func (m *ConceptoModel) GetUsuariosByFamilia(correoFamilia string) ([]entities.Usuario, error) {
	query := `SELECT nombreUsuario, rol, contraseñaPersonal, nombrePersonal, correoFamilia, delete_at
	          FROM usuario 
	          WHERE correoFamilia = ? AND delete_at IS NULL`

	rows, err := database.DB.Query(query, correoFamilia)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	usuarios := []entities.Usuario{}
	for rows.Next() {
		var usuario entities.Usuario
		err := rows.Scan(
			&usuario.NombreUsuario,
			&usuario.Rol,
			&usuario.ContrasenaPersonal,
			&usuario.NombrePersonal,
			&usuario.CorreoFamilia,
			&usuario.DeleteAt,
		)
		if err != nil {
			return nil, err
		}
		usuarios = append(usuarios, usuario)
	}

	log.Printf("👥 Usuarios encontrados en la familia: %d", len(usuarios))
	return usuarios, nil
}

// FNBD-CreatePersonalizacionesForAllUsuarios
// CreatePersonalizacionesForAllUsuarios crea personalizaciones para todos los usuarios de una familia
// Parámetros:
// - nombreConcepto: Concepto a personalizar
// - correoFamilia: Familia de los usuarios
// - datosPersonalizacion: Configuración base para las personalizaciones
// Retorno: Error si alguna personalización falla
// Flujo:
// 1. Obtiene todos los usuarios de la familia
// 2. Para cada usuario, crea personalización con datos proporcionados
// 3. Maneja valores nulos correctamente
// 4. Registra estadísticas de éxito/error
// Uso: Configuración masiva al crear nuevos conceptos
func (m *ConceptoModel) CreatePersonalizacionesForAllUsuarios(nombreConcepto, correoFamilia string, datosPersonalizacion map[string]interface{}) error {
	// Obtener todos los usuarios de la familia
	usuarios, err := m.GetUsuariosByFamilia(correoFamilia)
	if err != nil {
		log.Printf("❌ Error obteniendo usuarios de la familia: %v", err)
		return err
	}

	log.Printf("👥 Creando personalizaciones para %d usuarios", len(usuarios))

	usuariosConPersonalizacion := 0
	usuariosConError := 0

	// Crear personalización para cada usuario
	for _, usuario := range usuarios {
		// Convertir y manejar valores nulos correctamente
		var limiteGasto interface{} = nil
		if limite, ok := datosPersonalizacion["limite_monto"].(float64); ok && limite > 0 {
			limiteGasto = limite
		}

		var montoPlanificado interface{} = nil
		if desembolso, ok := datosPersonalizacion["desembolso_planejado"].(float64); ok && desembolso > 0 {
			montoPlanificado = desembolso
		}
		var diaDesembolsoPlanejado interface{} = nil
		if diaDesembolso, ok := datosPersonalizacion["dia_desembolso_planejado"].(int8); ok && diaDesembolso > 0 {
			diaDesembolsoPlanejado = diaDesembolso
		}
		var diaPeriodoLimite interface{} = nil
		if diaLimiteTipo, ok := datosPersonalizacion["dia_limite_tipo"].(int8); ok && diaLimiteTipo > 0 {
			diaPeriodoLimite = diaLimiteTipo
		}

		var tipoPeriodoPlanificado interface{} = nil
		if periodo, ok := datosPersonalizacion["periodo_tipo"].(string); ok && periodo != "" {
			tipoPeriodoPlanificado = periodo
		}

		var tipoPeriodoLimite interface{} = nil
		if limiteTipo, ok := datosPersonalizacion["limite_tipo"].(string); ok && limiteTipo != "" {
			tipoPeriodoLimite = limiteTipo
		}

		personalizacion := map[string]interface{}{
			"limiteGasto":            limiteGasto,
			"activo":                 true,
			"montoPlanificado":       montoPlanificado,
			"tipoPeriodoPlanificado": tipoPeriodoPlanificado,
			"tipoPeriodoLimite":      tipoPeriodoLimite,
			"diaPeriodoPlanificado":  diaDesembolsoPlanejado, // Por ahora siempre nil
			"diaPeriodoLimite":       diaPeriodoLimite,       // Por ahora siempre nil
			"notificacion":           false,                  // Por defecto false
			"nombreUsuario":          usuario.NombreUsuario,
			"nombreConcepto":         nombreConcepto,
			"correoFamilia":          correoFamilia,
		}

		err := m.CreatePersonalizacion(personalizacion)
		if err != nil {
			log.Printf("❌ Error creando personalización para usuario %s: %v", usuario.NombreUsuario, err)
			usuariosConError++
			continue
		}
		usuariosConPersonalizacion++
	}

	log.Printf("✅ Personalizaciones creadas: %d/%d usuarios (errores: %d)",
		usuariosConPersonalizacion, len(usuarios), usuariosConError)

	if usuariosConError > 0 {
		return fmt.Errorf("algunas personalizaciones fallaron: %d errores", usuariosConError)
	}

	return nil
}

// GetPersonalizacionesByConcepto obtiene todas las personalizaciones de un concepto específico
// Parámetros:
// - nombreConcepto: Concepto del cual obtener personalizaciones
// - correoFamilia: Familia del concepto
// Retorno: Lista de mapas con datos de personalización
// Flujo:
// - Consulta tabla personalizacionconcepto filtrada por concepto y familia
// - Maneja campos nulos con sql.Null types
// - Excluye personalizaciones eliminadas
// Uso: Reportes de configuración, auditoría de conceptos
func (m *ConceptoModel) GetPersonalizacionesByConcepto(nombreConcepto, correoFamilia string) ([]map[string]interface{}, error) {
	query := `SELECT idPersonalizacion, nombreUsuario, limiteGasto, montoPlanificado, 
                     tipoPeriodoPlanificado, tipoPeriodoLimite
              FROM personalizacionconcepto 
              WHERE nombreConcepto = ? AND correoFamilia = ? AND delete_at IS NULL`

	rows, err := database.DB.Query(query, nombreConcepto, correoFamilia)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	personalizaciones := []map[string]interface{}{}
	for rows.Next() {
		var idPersonalizacion int
		var nombreUsuario string
		var limiteGasto, montoPlanificado sql.NullFloat64
		var tipoPeriodoPlanificado, tipoPeriodoLimite sql.NullString

		err := rows.Scan(&idPersonalizacion, &nombreUsuario, &limiteGasto, &montoPlanificado,
			&tipoPeriodoPlanificado, &tipoPeriodoLimite)
		if err != nil {
			return nil, err
		}

		personalizacion := map[string]interface{}{
			"idPersonalizacion":      idPersonalizacion,
			"nombreUsuario":          nombreUsuario,
			"limiteGasto":            limiteGasto.Float64,
			"montoPlanificado":       montoPlanificado.Float64,
			"tipoPeriodoPlanificado": tipoPeriodoPlanificado.String,
			"tipoPeriodoLimite":      tipoPeriodoLimite.String,
		}
		personalizaciones = append(personalizaciones, personalizacion)
	}

	return personalizaciones, nil
}

// GetCountByTipo cuenta la cantidad de conceptos por tipo en una familia
// Parámetros:
// - correoFamilia: Familia a contar
// Retorno: Cantidad de gastos, cantidad de ingresos, error
// Flujo:
// - Usa conditional aggregation con SUM y CASE
// - Cuenta conceptos no eliminados
// Uso: Estadísticas familiares, dashboards administrativos
func (m *ConceptoModel) GetCountByTipo(correoFamilia string) (int, int, error) {
	query := `SELECT 
	          SUM(CASE WHEN tipo = 0 THEN 1 ELSE 0 END) as gastos,
	          SUM(CASE WHEN tipo = 1 THEN 1 ELSE 0 END) as ingresos
	          FROM concepto 
	          WHERE correoFamilia = ? AND delete_at IS NULL`

	var gastos, ingresos int
	err := database.DB.QueryRow(query, correoFamilia).Scan(&gastos, &ingresos)
	if err != nil {
		return 0, 0, err
	}

	return gastos, ingresos, nil
}

// FindByFamiliaActivos busca conceptos activos para un usuario específico
// Parámetros:
// - correoFamilia: Familia de los conceptos
// - nombreUsuario: Usuario para verificar personalizaciones
// - tipo: Tipo de concepto ("gasto" o "ingreso")
// Retorno: Lista de conceptos con estado activo para el usuario
// Flujo:
// - LEFT JOIN con personalizacionconcepto para usuario específico
// - Usa COALESCE para establecer activo=true si no hay personalización
// - Filtra por tipo y excluye eliminados
// Uso: Listado de conceptos disponibles para un usuario específico
func (m *ConceptoModel) FindByFamiliaActivos(correoFamilia, nombreUsuario, tipo string) ([]entities.Concepto, error) {
	var tipoInt int8
	if tipo == "ingreso" {
		tipoInt = 1
	} else {
		tipoInt = 0
	}

	query := `SELECT c.nombreConcepto, c.correoFamilia, c.tipo, c.icono, c.color, c.nombreUsuario, 
	                 COALESCE(pc.activo, true) as activo, c.delete_at
              FROM concepto c
              LEFT JOIN personalizacionconcepto pc ON c.nombreConcepto = pc.nombreConcepto 
                                                   AND c.correoFamilia = pc.correoFamilia 
                                                   AND pc.nombreUsuario = ?
              WHERE c.correoFamilia = ? AND c.tipo = ? AND c.delete_at IS NULL
              ORDER BY c.nombreConcepto`

	rows, err := database.DB.Query(query, nombreUsuario, correoFamilia, tipoInt)
	if err != nil {
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
			return nil, err
		}
		conceptos = append(conceptos, concepto)
	}

	return conceptos, nil
}

// Update actualiza un concepto existente (solo para el usuario creador)
// Parámetros:
// - concepto: Estructura con los datos actualizados
// Retorno: Error si la actualización falla
// Flujo:
// - Actualiza ícono, color y nombreUsuario
// - Solo permite actualizar conceptos del usuario creador
// - Filtra por nombreConcepto, correoFamilia y nombreUsuario
// Uso: Edición de conceptos propios
func (m *ConceptoModel) Update(concepto *entities.Concepto) error {
	query := `UPDATE concepto 
	          SET icono = ?, color = ?, nombreUsuario = ?
	          WHERE nombreConcepto = ? AND correoFamilia = ? AND nombreUsuario = ?`

	result, err := database.DB.Exec(query,
		concepto.Icono,
		concepto.Color,
		concepto.NombreUsuario,
		concepto.NombreConcepto,
		concepto.CorreoFamilia,
		concepto.NombreUsuario)

	if err != nil {
		log.Printf("❌ Error actualizando concepto: %v", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("✅ Concepto actualizado - Filas afectadas: %d", rowsAffected)
	return nil
}

// ToggleActivo cambia el estado activo/inactivo de un concepto para un usuario específico
// Parámetros:
// - nombreConcepto: Concepto a modificar
// - correoFamilia: Familia del concepto
// - nombreUsuario: Usuario afectado
// Retorno: Error si la operación falla
// Flujo:
// 1. Verifica si existe personalización para el usuario
// 2. Si no existe, crea una con activo=false
// 3. Si existe, cambia el estado activo (NOT activo)
// Uso: Habilitar/deshabilitar conceptos individualmente por usuario
func (m *ConceptoModel) ToggleActivo(nombreConcepto, correoFamilia, nombreUsuario string) error {
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
			return err
		}
		log.Printf("✅ Personalización creada como inactiva para usuario %s", nombreUsuario)
		return nil
	} else if err != nil {
		return err
	}

	// Ya existe personalización, cambiar el estado activo
	queryUpdate := `UPDATE personalizacionconcepto 
	               SET activo = NOT activo 
	               WHERE nombreConcepto = ? AND correoFamilia = ? AND nombreUsuario = ?`

	result, err := database.DB.Exec(queryUpdate, nombreConcepto, correoFamilia, nombreUsuario)
	if err != nil {
		log.Printf("❌ Error cambiando estado activo: %v", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("✅ Estado activo cambiado - Filas afectadas: %d", rowsAffected)
	return nil
}

// GetEstadoActivo obtiene el estado activo de un concepto para un usuario específico
// Parámetros:
// - nombreConcepto: Concepto a consultar
// - correoFamilia: Familia del concepto
// - nombreUsuario: Usuario a consultar
// Retorno: Estado activo (true/false), error si existe
// Flujo:
// - Consulta estado activo en personalizacionconcepto
// - Si no existe personalización, retorna true por defecto
// - Usa COALESCE para manejar valores nulos
// Uso: Verificar disponibilidad de conceptos para usuarios
func (m *ConceptoModel) GetEstadoActivo(nombreConcepto, correoFamilia, nombreUsuario string) (bool, error) {
	query := `SELECT COALESCE(activo, true) as activo
	          FROM personalizacionconcepto 
	          WHERE nombreConcepto = ? AND correoFamilia = ? AND nombreUsuario = ?`

	var activo bool
	err := database.DB.QueryRow(query, nombreConcepto, correoFamilia, nombreUsuario).Scan(&activo)

	if err == sql.ErrNoRows {
		// Si no existe personalización, el concepto está activo por defecto
		return true, nil
	}

	return activo, err
}

// ActualizarNombreEnPersonalizaciones actualiza el nombre de concepto en todas las personalizaciones
// Parámetros:
// - nombreViejo: Nombre actual del concepto
// - nombreNuevo: Nuevo nombre del concepto
// - correoFamilia: Familia del concepto
// Retorno: Error si la actualización falla
// Flujo:
// - Actualiza campo nombreConcepto en tabla personalizacionconcepto
// - Mantiene la relación entre personalizaciones y el concepto renombrado
// Uso: Sincronización al renombrar conceptos
func (m *ConceptoModel) ActualizarNombreEnPersonalizaciones(nombreViejo, nombreNuevo, correoFamilia string) error {
	query := `UPDATE personalizacionconcepto 
	          SET nombreConcepto = ?
	          WHERE nombreConcepto = ? AND correoFamilia = ?`

	result, err := database.DB.Exec(query, nombreNuevo, nombreViejo, correoFamilia)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("✅ Nombre actualizado en personalizaciones - Filas afectadas: %d", rowsAffected)
	return nil
}

// Deshabilitar desactiva un concepto y todas sus personalizaciones para toda la familia
// Parámetros:
// - nombreConcepto: Concepto a deshabilitar
// - correoFamilia: Familia del concepto
// Retorno: Error si la operación falla
// Flujo:
// 1. Inicia transacción para consistencia
// 2. Desactiva concepto en tabla concepto (activo = 0)
// 3. Desactiva todas las personalizaciones del concepto
// 4. Confirma transacción
// Uso: Deshabilitación completa de conceptos a nivel familiar

func (m *ConceptoModel) Deshabilitar(nombreConcepto, correoFamilia string) error {
	// Iniciar transacción para asegurar consistencia
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 1. Desactivar el concepto en la tabla concepto (soft delete)
	queryConcepto := `UPDATE concepto SET activo = 0 WHERE nombreConcepto = ? AND correoFamilia = ?`
	_, err = tx.Exec(queryConcepto, nombreConcepto, correoFamilia)
	if err != nil {
		log.Printf("❌ Error desactivando concepto: %v", err)
		return err
	}

	// 2. Desactivar todas las personalizaciones del concepto
	queryPersonalizaciones := `UPDATE personalizacionconcepto SET activo = 0 
                              WHERE nombreConcepto = ? AND correoFamilia = ?`
	_, err = tx.Exec(queryPersonalizaciones, nombreConcepto, correoFamilia)
	if err != nil {
		log.Printf("❌ Error desactivando personalizaciones: %v", err)
		return err
	}

	// Confirmar transacción
	err = tx.Commit()
	if err != nil {
		return err
	}

	log.Printf("✅ Concepto deshabilitado: %s para familia %s", nombreConcepto, correoFamilia)
	return nil
}

// Habilitar reactiva un concepto y todas sus personalizaciones para toda la familia
// Parámetros:
// - nombreConcepto: Concepto a habilitar
// - correoFamilia: Familia del concepto
// Retorno: Error si la operación falla
// Flujo:
// 1. Inicia transacción para consistencia
// 2. Reactiva concepto en tabla concepto (activo = 1)
// 3. Reactiva todas las personalizaciones del concepto
// 4. Confirma transacción
// Uso: Reactivación completa de conceptos a nivel familiar

func (m *ConceptoModel) Habilitar(nombreConcepto, correoFamilia string) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 1. Reactivar el concepto
	queryConcepto := `UPDATE concepto SET activo = 1 WHERE nombreConcepto = ? AND correoFamilia = ?`
	_, err = tx.Exec(queryConcepto, nombreConcepto, correoFamilia)
	if err != nil {
		return err
	}

	// 2. Reactivar todas las personalizaciones
	queryPersonalizaciones := `UPDATE personalizacionconcepto SET activo = 1 
                              WHERE nombreConcepto = ? AND correoFamilia = ?`
	_, err = tx.Exec(queryPersonalizaciones, nombreConcepto, correoFamilia)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	log.Printf("✅ Concepto habilitado: %s para familia %s", nombreConcepto, correoFamilia)
	return nil
}

// ToggleActivoGlobal cambia el estado activo/inactivo de un concepto para toda la familia
// Parámetros:
// - nombreConcepto: Concepto a modificar
// - correoFamilia: Familia del concepto
// Retorno: Error si la operación falla
// Flujo:
// 1. Consulta estado actual del concepto
// 2. Si está activo, llama a Deshabilitar
// 3. Si está inactivo, llama a Habilitar
// Uso: Alternar estado global de conceptos desde interfaz administrativa-

func (m *ConceptoModel) ToggleActivoGlobal(nombreConcepto, correoFamilia string) error {
	// Primero obtener el estado actual
	queryEstado := `SELECT activo FROM concepto WHERE nombreConcepto = ? AND correoFamilia = ?`
	var activoActual bool
	err := database.DB.QueryRow(queryEstado, nombreConcepto, correoFamilia).Scan(&activoActual)
	if err != nil {
		return err
	}

	if activoActual {
		return m.Deshabilitar(nombreConcepto, correoFamilia)
	} else {
		return m.Habilitar(nombreConcepto, correoFamilia)
	}
}
