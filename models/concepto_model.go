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

// FNBD_Concepto_Crear
// CreateTransactional crea un concepto y sus personalizaciones dentro de una transacción única en la BD.
// Parámetros:
// - concepto: Estructura con datos del concepto a crear.
// - datosPersonalizacion: Mapa con la configuración base de personalización.
// Retorno: Error si la transacción falla.
// Flujo:
// - Llama al SP transaccional que ejecuta CREATE y CREATE_PERSONALIZACIONES.
// - La BD maneja el COMMIT/ROLLBACK de ambas operaciones.
// Uso: Creación atómica de conceptos.
func (m *ConceptoModel) Create(concepto *entities.Concepto, datosPersonalizacion map[string]interface{}) error {
	// Ajustamos los punteros a valores nulos o el valor base para los SPs.
	var icono, color *string
	if concepto.Icono != nil {
		icono = concepto.Icono
	}
	if concepto.Color != nil {
		color = concepto.Color
	}

	// Preparación de valores nulos y types para el SP (similar a tu lógica en CreatePersonalizacionesForAllUsuarios)
	var limiteGasto, montoPlanificado interface{} = nil, nil
	if limite, ok := datosPersonalizacion["limite_monto"].(float64); ok && limite > 0 {
		limiteGasto = limite
	}
	if desembolso, ok := datosPersonalizacion["desembolso_planejado"].(float64); ok && desembolso > 0 {
		montoPlanificado = desembolso
	}

	var diaDesembolsoPlanejado, diaPeriodoLimite interface{} = nil, nil
	if diaDesembolso, ok := datosPersonalizacion["dia_desembolso_planejado"].(int8); ok && diaDesembolso > 0 {
		diaDesembolsoPlanejado = diaDesembolso
	}
	if diaLimiteTipo, ok := datosPersonalizacion["dia_limite_tipo"].(int8); ok && diaLimiteTipo > 0 {
		diaPeriodoLimite = diaLimiteTipo
	}

	var tipoPeriodoPlanificado, tipoPeriodoLimite interface{} = nil, nil
	if periodo, ok := datosPersonalizacion["periodo_tipo"].(string); ok && periodo != "" {
		tipoPeriodoPlanificado = periodo
	}
	if limiteTipo, ok := datosPersonalizacion["limite_tipo"].(string); ok && limiteTipo != "" {
		tipoPeriodoLimite = limiteTipo
	}

	query := `CALL sp_create_concepto_transaccional(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// Ejecutamos el SP transaccional
	_, err := database.DB.Exec(query,
		// Parámetros de Concepto
		concepto.NombreConcepto,
		concepto.CorreoFamilia,
		concepto.Tipo,
		icono,
		color,
		concepto.NombreUsuario,

		// Parámetros de Personalización
		limiteGasto,
		montoPlanificado,
		tipoPeriodoPlanificado,
		tipoPeriodoLimite,
		diaDesembolsoPlanejado,
		diaPeriodoLimite)

	if err != nil {
		// MySQL devolverá el error '45000' si se activa el ROLLBACK
		log.Printf("❌ Error al ejecutar transacción de creación de concepto: %v", err)
		return err
	}

	log.Printf("✅ Concepto '%s' y personalizaciones creadas de forma atómica.", concepto.NombreConcepto)
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
			log.Printf("❌ Error escaneando concepto 'Funcion FindByFamilia': %v", err)
			return nil, err
		}
		conceptos = append(conceptos, concepto)
	}

	log.Printf("📋 Conceptos encontrados: %d (tipo: %s, correo: %s)", len(conceptos), tipo, correoFamilia)
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
	query := `CALL sp_update_concepto_icono_color(?, ?, ?, ?)`

	var rowsAffected int64

	err := database.DB.QueryRow(query,
		concepto.Icono,
		concepto.Color,
		concepto.NombreConcepto,
		concepto.CorreoFamilia).Scan(&rowsAffected)

	if err != nil {
		log.Printf("❌ Error al ejecutar sp_update_concepto_icono_color: %v", err)
		return err
	}

	if rowsAffected == 0 {
		log.Printf("⚠️ No se encontró el concepto o ya fue eliminado - Nombre: %s, Familia: %s",
			concepto.NombreConcepto, concepto.CorreoFamilia)
	} else {
		log.Printf("✅ Concepto actualizado - Filas afectadas: %d", rowsAffected)
	}

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
	query := `CALL sp_find_all_conceptos_by_familia(?)`

	rows, err := database.DB.Query(query, correoFamilia)
	if err != nil {
		log.Printf("❌ Error al ejecutar sp_find_all_conceptos_by_familia: %v", err)
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
			log.Printf("❌ Error escaneando concepto 'Funcion FindAllBytFamilia': %v", err)
			return nil, err
		}
		conceptos = append(conceptos, concepto)
	}

	log.Printf("📋 Total de conceptos encontrados para familia '%s': %d", correoFamilia, len(conceptos))
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
	query := `CALL sp_find_concepto_by_nombre(?, ?)`

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
		log.Printf("ℹ️ Concepto no encontrado - Nombre: %s, Familia: %s", nombreConcepto, correoFamilia)
		return nil, nil
	}

	if err != nil {
		log.Printf("❌ Error al ejecutar sp_find_concepto_by_nombre: %v", err)
		return nil, err
	}

	log.Printf("✅ Concepto encontrado - Nombre: %s, Familia: %s", nombreConcepto, correoFamilia)
	return concepto, nil
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
	query := `CALL sp_concepto_exists(?, ?)`

	var count int
	err := database.DB.QueryRow(query, nombreConcepto, correoFamilia).Scan(&count)
	if err != nil {
		log.Printf("❌ Error al ejecutar sp_concepto_exists: %v", err)
		return false, err
	}

	exists := count > 0
	log.Printf("🔍 Verificación de existencia - Concepto: %s, Familia: %s, Existe: %v",
		nombreConcepto, correoFamilia, exists)

	return exists, nil
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
	query := `CALL sp_create_personalizacion(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	var rowsAffected int64
	var idPersonalizacion int64

	err := database.DB.QueryRow(query,
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
		personalizacion["correoFamilia"]).Scan(&rowsAffected, &idPersonalizacion)

	if err != nil {
		log.Printf("❌ Error al ejecutar sp_create_personalizacion: %v", err)
		log.Printf("   Datos: %+v", personalizacion)
		return err
	}

	log.Printf("✅ Personalización creada para %s - ID: %d, Filas afectadas: %d",
		personalizacion["nombreUsuario"], idPersonalizacion, rowsAffected)

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
	query := `CALL sp_get_usuarios_by_familia(?)`

	rows, err := database.DB.Query(query, correoFamilia)
	if err != nil {
		log.Printf("❌ Error al ejecutar sp_get_usuarios_by_familia: %v", err)
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
			log.Printf("❌ Error escaneando usuario: %v", err)
			return nil, err
		}
		usuarios = append(usuarios, usuario)
	}

	log.Printf("👥 Usuarios encontrados en la familia '%s': %d", correoFamilia, len(usuarios))
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
	query := `CALL sp_create_personalizaciones_all_usuarios(?, ?, ?, ?, ?, ?, ?, ?)`

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

	var usuariosTotal, usuariosCreados, usuariosError int

	err := database.DB.QueryRow(query,
		nombreConcepto,
		correoFamilia,
		limiteGasto,
		montoPlanificado,
		tipoPeriodoPlanificado,
		tipoPeriodoLimite,
		diaDesembolsoPlanejado,
		diaPeriodoLimite).Scan(&usuariosTotal, &usuariosCreados, &usuariosError)

	if err != nil {
		log.Printf("❌ Error al ejecutar sp_create_personalizaciones_all_usuarios: %v", err)
		return err
	}

	log.Printf("👥 Creando personalizaciones para %d usuarios", usuariosTotal)
	log.Printf("✅ Personalizaciones creadas: %d/%d usuarios (errores: %d)",
		usuariosCreados, usuariosTotal, usuariosError)

	if usuariosError > 0 {
		return fmt.Errorf("algunas personalizaciones fallaron: %d errores", usuariosError)
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
	query := `CALL sp_get_personalizaciones_by_concepto(?, ?)`

	rows, err := database.DB.Query(query, nombreConcepto, correoFamilia)
	if err != nil {
		log.Printf("❌ Error al ejecutar sp_get_personalizaciones_by_concepto: %v", err)
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
			log.Printf("❌ Error escaneando personalización: %v", err)
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

	log.Printf("📋 Personalizaciones encontradas para concepto '%s': %d", nombreConcepto, len(personalizaciones))
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
	query := `CALL sp_get_count_by_tipo(?)`

	var gastos, ingresos sql.NullInt64
	err := database.DB.QueryRow(query, correoFamilia).Scan(&gastos, &ingresos)
	if err != nil {
		log.Printf("❌ Error al ejecutar sp_get_count_by_tipo: %v", err)
		return 0, 0, err
	}

	// Convertir NULL a 0
	gastosVal := int(0)
	if gastos.Valid {
		gastosVal = int(gastos.Int64)
	}

	ingresosVal := int(0)
	if ingresos.Valid {
		ingresosVal = int(ingresos.Int64)
	}

	log.Printf("📊 Estadísticas familia '%s' - Gastos: %d, Ingresos: %d",
		correoFamilia, gastosVal, ingresosVal)

	return gastosVal, ingresosVal, nil
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

	query := `CALL sp_find_conceptos_by_familia_activos(?, ?, ?)`

	rows, err := database.DB.Query(query, correoFamilia, nombreUsuario, tipoInt)
	if err != nil {
		log.Printf("❌ Error al ejecutar sp_find_conceptos_by_familia_activos: %v", err)
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
			log.Printf("❌ Error escaneando concepto 'Funcion FindFamiliaActivos': %v", err)
			return nil, err
		}
		conceptos = append(conceptos, concepto)
	}

	log.Printf("📋 Conceptos activos encontrados para usuario '%s' (tipo: %s): %d",
		nombreUsuario, tipo, len(conceptos))
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
	query := `CALL sp_update_concepto(?, ?, ?, ?, ?)`

	var rowsAffected int64

	err := database.DB.QueryRow(query,
		concepto.Icono,
		concepto.Color,
		concepto.NombreUsuario,
		concepto.NombreConcepto,
		concepto.CorreoFamilia).Scan(&rowsAffected)

	if err != nil {
		log.Printf("❌ Error al ejecutar sp_update_concepto: %v", err)
		return err
	}

	if rowsAffected == 0 {
		log.Printf("⚠️ Concepto no encontrado o no pertenece al usuario - Concepto: %s, Usuario: %s",
			concepto.NombreConcepto, concepto.NombreUsuario)
	} else {
		log.Printf("✅ Concepto actualizado - Concepto: %s, Filas afectadas: %d",
			concepto.NombreConcepto, rowsAffected)
	}

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
	query := `CALL sp_toggle_activo_concepto(?, ?, ?)`

	var accion string
	var nuevoEstado bool
	var rowsAffected int64

	err := database.DB.QueryRow(query, nombreConcepto, correoFamilia, nombreUsuario).
		Scan(&accion, &nuevoEstado, &rowsAffected)

	if err != nil {
		log.Printf("❌ Error al ejecutar sp_toggle_activo_concepto: %v", err)
		return err
	}

	if accion == "created" {
		log.Printf("✅ Personalización creada como inactiva - Usuario: %s, Concepto: %s",
			nombreUsuario, nombreConcepto)
	} else {
		log.Printf("✅ Estado activo cambiado - Usuario: %s, Concepto: %s, Nuevo estado: %v, Filas: %d",
			nombreUsuario, nombreConcepto, nuevoEstado, rowsAffected)
	}

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
	query := `CALL sp_get_estado_activo_concepto(?, ?, ?)`

	var activo bool
	err := database.DB.QueryRow(query, nombreConcepto, correoFamilia, nombreUsuario).Scan(&activo)

	if err == sql.ErrNoRows {
		// Si no existe personalización, el concepto está activo por defecto
		log.Printf("ℹ️ No existe personalización, concepto activo por defecto - Usuario: %s, Concepto: %s",
			nombreUsuario, nombreConcepto)
		return true, nil
	}

	if err != nil {
		log.Printf("❌ Error al ejecutar sp_get_estado_activo_concepto: %v", err)
		return false, err
	}

	log.Printf("🔍 Estado activo consultado - Usuario: %s, Concepto: %s, Activo: %v",
		nombreUsuario, nombreConcepto, activo)

	return activo, nil
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
	query := `CALL sp_actualizar_nombre_personalizaciones(?, ?, ?)`

	var rowsAffected int64

	err := database.DB.QueryRow(query, nombreViejo, nombreNuevo, correoFamilia).Scan(&rowsAffected)
	if err != nil {
		log.Printf("❌ Error al ejecutar sp_actualizar_nombre_personalizaciones: %v", err)
		return err
	}

	log.Printf("✅ Nombre actualizado en personalizaciones - '%s' → '%s', Filas afectadas: %d",
		nombreViejo, nombreNuevo, rowsAffected)

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
	query := `CALL sp_deshabilitar_concepto(?, ?)`

	var conceptoRows, personalizacionesRows int
	var estado string

	err := database.DB.QueryRow(query, nombreConcepto, correoFamilia).
		Scan(&conceptoRows, &personalizacionesRows, &estado)

	if err != nil {
		log.Printf("❌ Error al ejecutar sp_deshabilitar_concepto: %v", err)
		return err
	}

	log.Printf("✅ Concepto deshabilitado: '%s' para familia '%s'", nombreConcepto, correoFamilia)
	log.Printf("   📊 Concepto: %d fila(s), Personalizaciones: %d fila(s)",
		conceptoRows, personalizacionesRows)

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
	query := `CALL sp_habilitar_concepto(?, ?)`

	var conceptoRows, personalizacionesRows int
	var estado string

	err := database.DB.QueryRow(query, nombreConcepto, correoFamilia).
		Scan(&conceptoRows, &personalizacionesRows, &estado)

	if err != nil {
		log.Printf("❌ Error al ejecutar sp_habilitar_concepto: %v", err)
		return err
	}

	log.Printf("✅ Concepto habilitado: '%s' para familia '%s'", nombreConcepto, correoFamilia)
	log.Printf("   📊 Concepto: %d fila(s), Personalizaciones: %d fila(s)",
		conceptoRows, personalizacionesRows)

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
	query := `CALL sp_toggle_activo_global(?, ?)`

	var conceptoRows, personalizacionesRows int
	var accion string
	var nuevoEstado bool

	err := database.DB.QueryRow(query, nombreConcepto, correoFamilia).
		Scan(&conceptoRows, &personalizacionesRows, &accion, &nuevoEstado)

	if err != nil {
		log.Printf("❌ Error al ejecutar sp_toggle_activo_global: %v", err)
		return err
	}

	log.Printf("✅ Concepto '%s' %s para familia '%s' - Nuevo estado: %v",
		nombreConcepto, accion, correoFamilia, nuevoEstado)
	log.Printf("   📊 Concepto: %d fila(s), Personalizaciones: %d fila(s)",
		conceptoRows, personalizacionesRows)

	return nil
}
