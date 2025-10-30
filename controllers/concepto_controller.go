package controllers

import (
	"arka-code/entities"
	"arka-code/models"
	"arka-code/utils"
	"arka-code/validators"
	"log"
	"net/http"
)

type ConceptoController struct{}

var ConceptoControllerInstance = &ConceptoController{}

// Index muestra los conceptos filtrados por tipo
func (c *ConceptoController) Index(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔍 ConceptoController.Index llamado - Método: %s", r.Method)

	sessionData, ok := utils.GetSessionData(r)
	if !ok {
		log.Printf("❌ No hay sesión, redirigiendo a login")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	log.Printf("👤 Usuario en sesión: %s", sessionData.NombreUsuario)

	// Obtener tipo de concepto (gasto/ingreso)
	tipo := r.URL.Query().Get("tipo")
	if tipo == "" {
		tipo = "gasto" // valor por defecto
	}

	log.Printf("📊 Tipo de concepto: %s", tipo)

	// Obtener conceptos básicos
	conceptos, err := models.ConceptoModelInstance.FindByFamilia(sessionData.CorreoFamilia, tipo)
	if err != nil {
		log.Printf("❌ Error buscando conceptos: %v", err)
	}

	log.Printf("📋 Número de conceptos encontrados: %d", len(conceptos))

	// Obtener lista de íconos disponibles
	iconos := c.getIconosDisponibles()
	log.Printf("🎯 Número de íconos disponibles: %d", len(iconos))

	// Manejar mensajes de éxito/error
	successMsg := ""
	errorMsg := ""

	switch r.URL.Query().Get("success") {
	case "concepto_creado":
		successMsg = "Concepto creado exitosamente"
	}

	switch r.URL.Query().Get("error") {
	case "error_creando":
		errorMsg = "Error al crear el concepto"
	}

	data := map[string]interface{}{
		"Title":            "Conceptos",
		"CurrentPage":      "conceptos",
		"SessionData":      sessionData,
		"Tipo":             tipo,
		"Conceptos":        conceptos,
		"Iconos":           iconos,
		"ValidationErrors": nil,
		"Success":          successMsg,
		"Error":            errorMsg,
		"FormData":         map[string]interface{}{},
	}

	log.Printf("✅ Renderizando template conceptos")
	utils.RenderTemplate(w, "dashboard", "concepto/conceptos", data)
}

// Crear maneja la creación de un nuevo concepto
func (c *ConceptoController) Crear(w http.ResponseWriter, r *http.Request) {
	log.Printf("🚀 ConceptoController.Crear llamado - Método: %s", r.Method)

	if r.Method != http.MethodPost {
		log.Printf("❌ Método no permitido: %s, redirigiendo", r.Method)
		http.Redirect(w, r, "/conceptos", http.StatusSeeOther)
		return
	}

	sessionData, ok := utils.GetSessionData(r)
	if !ok {
		log.Printf("❌ No hay sesión en Crear")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	log.Printf("👤 Usuario creando concepto: %s", sessionData.NombreUsuario)

	// Parsear el formulario
	if err := r.ParseForm(); err != nil {
		log.Printf("❌ Error parseando formulario: %v", err)
		http.Redirect(w, r, "/conceptos?error=error_parseando_formulario", http.StatusSeeOther)
		return
	}

	// Log todos los datos del formulario
	log.Printf("📝 FORMULARIO RECIBIDO:")
	for key, values := range r.Form {
		log.Printf("   %s: %v", key, values)
	}

	// Validar datos del formulario
	validation := validators.ConceptoValidatorInstance.Validate(r)
	log.Printf("🔍 RESULTADO VALIDACIÓN - Éxito: %v, Errores: %v", validation.Success, validation.Errors)

	if !validation.Success {
		log.Printf("❌ Validación fallida, recargando página con errores")
		c.recargarPaginaConErrores(w, r, sessionData, validation.Errors, validation.CleanData)
		return
	}

	log.Printf("✅ Validación exitosa, datos limpios: %v", validation.CleanData)

	// Verificar si el concepto ya existe
	nombreConcepto := validation.CleanData["nombre"].(string)
	correoFamilia := sessionData.CorreoFamilia

	log.Printf("🔎 Verificando si existe concepto: %s para familia: %s", nombreConcepto, correoFamilia)

	existe, err := models.ConceptoModelInstance.Exists(nombreConcepto, correoFamilia)
	if err != nil {
		log.Printf("❌ Error verificando existencia: %v", err)
		c.recargarPaginaConErrores(w, r, sessionData,
			map[string]string{"general": "Error al verificar el concepto"}, validation.CleanData)
		return
	}

	if existe {
		log.Printf("❌ Concepto ya existe: %s", nombreConcepto)
		c.recargarPaginaConErrores(w, r, sessionData,
			map[string]string{"nombre": "Ya existe un concepto con este nombre"}, validation.CleanData)
		return
	}

	log.Printf("✅ Concepto no existe, procediendo a crear")

	// Crear el concepto
	tipoInt := int8(0) // 0 = gasto por defecto
	if validation.CleanData["tipo"] == "ingreso" {
		tipoInt = 1
	}

	// Obtener el ícono seleccionado
	idIcono := validation.CleanData["id_icono"].(int)
	iconoSeleccionado := c.getIconoPorID(idIcono)
	log.Printf("🎨 Ícono seleccionado: %s", iconoSeleccionado)

	// Crear variables para los pointers
	colorStr := validation.CleanData["color"].(string)

	concepto := &entities.Concepto{
		NombreConcepto: nombreConcepto,
		CorreoFamilia:  correoFamilia,
		Tipo:           tipoInt,
		Icono:          &iconoSeleccionado,
		Color:          &colorStr,
		NombreUsuario:  sessionData.NombreUsuario,
	}

	log.Printf("💾 Intentando crear concepto en BD: %+v", concepto)

	// Crear concepto en la base de datos
	err = models.ConceptoModelInstance.Create(concepto)
	if err != nil {
		log.Printf("❌ ERROR creando concepto en BD: %v", err)
		c.recargarPaginaConErrores(w, r, sessionData,
			map[string]string{"general": "Error al crear el concepto: " + err.Error()}, validation.CleanData)
		return
	}

	log.Printf("✅ CONCEPTO CREADO EXITOSAMENTE EN BD")

	// Crear personalizaciones para todos los usuarios de la familia
	log.Printf("🔄 Creando personalizaciones para todos los usuarios")

	// Preparar datos de personalización con valores por defecto
	datosPersonalizacion := map[string]interface{}{
		"limite_monto":            0.0,
		"desembolso_planejado":    0.0,
		"periodo_tipo":            "",
		"limite_tipo":             "",
		"dia_periodo_planificado": nil,
		"notificacion":            false,
	}

	// Sobrescribir con datos del formulario si existen
	if desembolso, ok := validation.CleanData["desembolso_planejado"].(float64); ok {
		datosPersonalizacion["desembolso_planejado"] = desembolso
	}
	if limite, ok := validation.CleanData["limite_monto"].(float64); ok {
		datosPersonalizacion["limite_monto"] = limite
	}
	if periodoTipo, ok := validation.CleanData["periodo_tipo"].(string); ok {
		datosPersonalizacion["periodo_tipo"] = periodoTipo
	}
	if limiteTipo, ok := validation.CleanData["limite_tipo"].(string); ok {
		datosPersonalizacion["limite_tipo"] = limiteTipo
	}

	log.Printf("📊 Datos de personalización: %+v", datosPersonalizacion)

	err = models.ConceptoModelInstance.CreatePersonalizacionesForAllUsuarios(
		nombreConcepto,
		correoFamilia,
		datosPersonalizacion,
	)
	if err != nil {
		log.Printf("⚠️  Error creando algunas personalizaciones: %v", err)
		// No fallamos la creación del concepto aunque falle alguna personalización
	} else {
		log.Printf("✅ Personalizaciones creadas exitosamente")
	}

	// Verificar que se crearon las personalizaciones
	personalizaciones, err := models.ConceptoModelInstance.GetPersonalizacionesByConcepto(nombreConcepto, correoFamilia)
	if err != nil {
		log.Printf("⚠️  Error verificando personalizaciones: %v", err)
	} else {
		log.Printf("🔍 Personalizaciones creadas para concepto '%s': %d", nombreConcepto, len(personalizaciones))
		for _, p := range personalizaciones {
			log.Printf("   👤 %s - Límite: %.2f - Desembolso: %.2f",
				p["nombreUsuario"], p["limiteGasto"], p["montoPlanificado"])
		}
	}

	log.Printf("🎉 REDIRIGIENDO A LISTA DE CONCEPTOS")
	http.Redirect(w, r, "/conceptos?success=concepto_creado&tipo="+validation.CleanData["tipo"].(string), http.StatusSeeOther)
}

// recargarPaginaConErrores recarga la página mostrando errores de validación
func (c *ConceptoController) recargarPaginaConErrores(w http.ResponseWriter, r *http.Request, sessionData *entities.SessionData, errors map[string]string, formData map[string]interface{}) {
	log.Printf("🔄 Recargando página con errores: %v", errors)

	tipo := "gasto"
	if tipoForm, ok := formData["tipo"].(string); ok {
		tipo = tipoForm
	} else {
		tipo = r.FormValue("tipo")
	}

	conceptos, _ := models.ConceptoModelInstance.FindByFamilia(sessionData.CorreoFamilia, tipo)
	iconos := c.getIconosDisponibles()

	data := map[string]interface{}{
		"Title":            "Conceptos",
		"CurrentPage":      "conceptos",
		"SessionData":      sessionData,
		"Tipo":             tipo,
		"Conceptos":        conceptos,
		"Iconos":           iconos,
		"ValidationErrors": errors,
		"FormData":         formData,
	}

	utils.RenderTemplate(w, "dashboard", "concepto/conceptos", data)
}

// getIconosDisponibles retorna la lista de íconos disponibles
func (c *ConceptoController) getIconosDisponibles() []map[string]interface{} {
	return []map[string]interface{}{
		{"ID": 1, "Nombre": "Casa", "Path": "fa-solid fa-house"},
		{"ID": 2, "Nombre": "Transporte", "Path": "fa-solid fa-car"},
		{"ID": 3, "Nombre": "Comida", "Path": "fa-solid fa-utensils"},
		{"ID": 4, "Nombre": "Salud", "Path": "fa-solid fa-heart-pulse"},
		{"ID": 5, "Nombre": "Educación", "Path": "fa-solid fa-graduation-cap"},
		{"ID": 6, "Nombre": "Entretenimiento", "Path": "fa-solid fa-film"},
		{"ID": 7, "Nombre": "Ropa", "Path": "fa-solid fa-shirt"},
		{"ID": 8, "Nombre": "Tecnología", "Path": "fa-solid fa-laptop"},
		{"ID": 9, "Nombre": "Deportes", "Path": "fa-solid fa-dumbbell"},
		{"ID": 10, "Nombre": "Viajes", "Path": "fa-solid fa-plane"},
		{"ID": 11, "Nombre": "Salario", "Path": "fa-solid fa-money-bill-wave"},
		{"ID": 12, "Nombre": "Regalo", "Path": "fa-solid fa-gift"},
		{"ID": 13, "Nombre": "Ahorros", "Path": "fa-solid fa-piggy-bank"},
		{"ID": 14, "Nombre": "Inversión", "Path": "fa-solid fa-chart-line"},
		{"ID": 15, "Nombre": "Otros", "Path": "fa-solid fa-circle"},
	}
}

// getIconoPorID retorna el path del ícono basado en el ID
func (c *ConceptoController) getIconoPorID(id int) string {
	iconos := c.getIconosDisponibles()
	for _, icono := range iconos {
		if icono["ID"] == id {
			return icono["Path"].(string)
		}
	}
	return "fa-solid fa-circle"
}
