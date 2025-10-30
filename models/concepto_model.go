package models

import (
	"arka-code/database"
	"arka-code/entities"
	"database/sql"
	"fmt"
	"log"
)

type ConceptoModel struct{}

var ConceptoModelInstance = &ConceptoModel{}

// Create crea un nuevo concepto
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

// FindByFamilia busca todos los conceptos de una familia por tipo
func (m *ConceptoModel) FindByFamilia(correoFamilia string, tipo string) ([]entities.Concepto, error) {
	// Convertir string a int8 para la base de datos
	var tipoInt int8
	if tipo == "ingreso" {
		tipoInt = 1
	} else {
		tipoInt = 0 // gasto por defecto
	}

	query := `SELECT nombreConcepto, correoFamilia, tipo, icono, color, nombreUsuario, delete_at
              FROM concepto 
              WHERE correoFamilia = ? AND tipo = ? AND delete_at IS NULL
              ORDER BY nombreConcepto`

	rows, err := database.DB.Query(query, correoFamilia, tipoInt)
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

// FindByNombre busca un concepto por su nombre y familia
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

// Exists verifica si un concepto ya existe para una familia
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

// CreatePersonalizacion crea una personalización para un usuario y concepto
func (m *ConceptoModel) CreatePersonalizacion(personalizacion map[string]interface{}) error {
	query := `INSERT INTO personalizacionconcepto 
              (limiteGasto, activo, montoPlanificado, tipoPeriodoPlanificado, 
               tipoPeriodoLimite, diaPeriodoPlanificado, notificacion,
               nombreUsuario, nombreConcepto, correoFamilia)
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := database.DB.Exec(query,
		personalizacion["limiteGasto"],
		personalizacion["activo"],
		personalizacion["montoPlanificado"],
		personalizacion["tipoPeriodoPlanificado"],
		personalizacion["tipoPeriodoLimite"],
		personalizacion["diaPeriodoPlanificado"],
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

// GetUsuariosByFamilia obtiene todos los usuarios de una familia
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

// CreatePersonalizacionesForAllUsuarios crea personalizaciones para todos los usuarios de la familia
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
			"diaPeriodoPlanificado":  nil,   // Por ahora siempre nil
			"notificacion":           false, // Por defecto false
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

// GetPersonalizacionesByConcepto obtiene todas las personalizaciones de un concepto
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

// GetCountByTipo obtiene la cantidad de conceptos por tipo
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
