package controllers

import (
	"arka-code/entities"
	"arka-code/models"
	"arka-code/utils"
	"arka-code/validators"
	"log"
	"net/http"
)

/*
CU-004 Gestionar Conceptos
CTRL.5 Gestor Conceptos
*/

type ConceptoController struct{}

var ConceptoControllerInstance = &ConceptoController{}

// FNConcepto-Index
// Index muestra la página principal de conceptos filtrados por tipo (gasto/ingreso)
// Parámetros:
// - w: ResponseWriter para enviar respuesta HTTP
// - r: Request HTTP con parámetros de filtrado
// Flujo:
// 1. Verifica autenticación del usuario
// 2. Obtiene tipo de concepto desde URL (gasto/ingreso)
// 3. Carga conceptos de la familia según tipo
// 4. Prepara lista de íconos disponibles
// 5. Maneja mensajes de éxito/error desde URL
// 6. Renderiza plantilla con todos los datos
// Uso: Vista principal de gestión de conceptos
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

	// Obtener conceptos básicos
	conceptos, err := models.ConceptoModelInstance.FindByFamilia(sessionData.CorreoFamilia, tipo)
	log.Printf("Conceptosssssssssssssssssssssss: %+v", conceptos)
	if err != nil {
		log.Printf("❌ Error buscando conceptos: %v", err)
	}

	// Obtener lista de íconos disponibles
	iconos := c.getIconosDisponibles()

	// Manejar mensajes de éxito/error
	successMsg := ""
	errorMsg := ""

	switch r.URL.Query().Get("success") {
	case "concepto_creado":
		successMsg = "Concepto creado exitosamente"
	case "concepto_deshabilitado":
		successMsg = "Concepto deshabilitado exitosamente"
	case "concepto_habilitado":
		successMsg = "Concepto habilitado exitosamente"
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

// FNConcepto-Crear
// Crear maneja la creación de un nuevo concepto en el sistema
// Parámetros:
// - w: ResponseWriter para enviar respuesta HTTP
// - r: Request HTTP con datos del formulario
// Flujo:
// 1. Verifica método POST y autenticación
// 2. Parsea y valida datos del formulario
// 3. Verifica que el concepto no exista
// 4. Crea el concepto en la base de datos
// 5. Crea personalizaciones para todos los usuarios de la familia
// 6. Redirige con mensaje de éxito o error
// Uso: Creación de nuevos conceptos desde formulario
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
	if diaDesembolsoPlanejado, ok := validation.CleanData["dia_desembolso_planejado"].(int8); ok {
		datosPersonalizacion["dia_desembolso_planejado"] = diaDesembolsoPlanejado
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
	if diaLimiteTipo, ok := validation.CleanData["dia_limite_tipo"].(int8); ok {
		datosPersonalizacion["dia_limite_tipo"] = diaLimiteTipo
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

// FNConcepto-Editar
// Editar maneja la edición de conceptos existentes (solo ícono y color)
// Parámetros:
// - w: ResponseWriter para enviar respuesta HTTP
// - r: Request HTTP con datos de edición
// Flujo GET:
// 1. Muestra formulario de edición con datos actuales del concepto
// 2. Verifica permisos del usuario sobre el concepto
// Flujo POST:
// 1. Valida datos del formulario de edición
// 2. Verifica que el concepto exista y el usuario tenga permisos
// 3. Actualiza ícono y color del concepto
// 4. Redirige con mensaje de éxito
// Uso: Modificación de aspectos visuales de conceptos existentes
func (c *ConceptoController) Editar(w http.ResponseWriter, r *http.Request) {
	log.Printf("✏️ ConceptoController.Editar llamado - Método: %s", r.Method)

	sessionData, ok := utils.GetSessionData(r)
	if !ok {
		log.Printf("❌ No hay sesión en Editar")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		// Mostrar formulario de edición
		nombreConcepto := r.URL.Query().Get("nombre")
		if nombreConcepto == "" {
			log.Printf("❌ Nombre de concepto vacío")
			http.Error(w, "Nombre de concepto requerido", http.StatusBadRequest)
			return
		}

		// Obtener el concepto
		concepto, err := models.ConceptoModelInstance.FindByNombre(nombreConcepto, sessionData.CorreoFamilia)
		if err != nil || concepto == nil {
			log.Printf("❌ Concepto no encontrado: %s", nombreConcepto)
			http.Error(w, "Concepto no encontrado", http.StatusNotFound)
			return
		}

		// Obtener todos los conceptos del mismo tipo para mostrar la lista
		tipo := concepto.TipoNombre()
		conceptos, _ := models.ConceptoModelInstance.FindByFamilia(sessionData.CorreoFamilia, tipo)
		iconos := c.getIconosDisponibles()

		// Encontrar el ID del ícono actual
		idIconoActual := 15 // Por defecto "Otros"
		for _, icono := range iconos {
			if icono["Path"].(string) == *concepto.Icono {
				idIconoActual = icono["ID"].(int)
				break
			}
		}

		data := map[string]interface{}{
			"Title":            "Editar Concepto",
			"CurrentPage":      "conceptos",
			"SessionData":      sessionData,
			"Tipo":             tipo,
			"Conceptos":        conceptos,
			"Iconos":           iconos,
			"ValidationErrors": nil,
			"ModoEdicion":      true,
			"ConceptoEditar":   concepto,
			"FormData": map[string]interface{}{
				"Nombre":  concepto.NombreConcepto,
				"IDIcono": idIconoActual,
				"Color":   *concepto.Color,
			},
		}

		utils.RenderTemplate(w, "dashboard", "concepto/conceptos", data)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Procesar actualización
	if err := r.ParseForm(); err != nil {
		log.Printf("❌ Error parseando formulario: %v", err)
		http.Error(w, "Error procesando solicitud", http.StatusBadRequest)
		return
	}

	log.Printf("📝 FORMULARIO EDICIÓN RECIBIDO:")
	for key, values := range r.Form {
		log.Printf("   %s: %v", key, values)
	}

	// Validar datos
	validation := validators.ConceptoValidatorInstance.ValidateUpdate(r)
	if !validation.Success {
		log.Printf("❌ Validación fallida: %v", validation.Errors)
		c.recargarEdicionConErrores(w, r, sessionData, validation.Errors, validation.CleanData)
		return
	}

	nombreConcepto := validation.CleanData["nombre"].(string)

	// Verificar que el concepto existe
	concepto, err := models.ConceptoModelInstance.FindByNombre(nombreConcepto, sessionData.CorreoFamilia)
	if err != nil || concepto == nil {
		log.Printf("❌ Concepto no encontrado: %s", nombreConcepto)
		http.Error(w, "Concepto no encontrado", http.StatusNotFound)
		return
	}

	// Solo el creador puede editar
	if concepto.NombreUsuario != sessionData.NombreUsuario && sessionData.Rol != 1 {
		log.Printf("❌ Usuario %s no tiene permisos para editar concepto de %s",
			sessionData.NombreUsuario, concepto.NombreUsuario)
		http.Error(w, "No tienes permisos para editar este concepto", http.StatusForbidden)
		return
	}

	// Actualizar ícono y color
	idIcono := validation.CleanData["id_icono"].(int)
	iconoNuevo := c.getIconoPorID(idIcono)
	colorNuevo := validation.CleanData["color"].(string)

	concepto.Icono = &iconoNuevo
	concepto.Color = &colorNuevo

	err = models.ConceptoModelInstance.UpdateIconoColor(concepto)
	if err != nil {
		log.Printf("❌ Error actualizando concepto: %v", err)
		http.Error(w, "Error al actualizar el concepto", http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Concepto actualizado exitosamente: %s", nombreConcepto)
	http.Redirect(w, r, "/conceptos?tipo="+concepto.TipoNombre()+"&success=concepto_editado", http.StatusSeeOther)
}

// FNConcepto-recargarEdicionConErrores
// recargarEdicionConErrores recarga la página de edición mostrando errores de validación
// Parámetros:
// - w: ResponseWriter para enviar respuesta HTTP
// - r: Request HTTP original
// - sessionData: Datos de sesión del usuario
// - errors: Mapa de errores de validación
// - formData: Datos del formulario para repoblar campos
// Uso: Reutilización de vista de edición cuando hay errores de validación
func (c *ConceptoController) recargarEdicionConErrores(w http.ResponseWriter, r *http.Request, sessionData *entities.SessionData, errors map[string]string, formData map[string]interface{}) {
	nombreConcepto := formData["nombre"].(string)

	concepto, _ := models.ConceptoModelInstance.FindByNombre(nombreConcepto, sessionData.CorreoFamilia)
	if concepto == nil {
		http.Error(w, "Concepto no encontrado", http.StatusNotFound)
		return
	}

	tipo := concepto.TipoNombre()
	conceptos, _ := models.ConceptoModelInstance.FindByFamilia(sessionData.CorreoFamilia, tipo)
	iconos := c.getIconosDisponibles()

	data := map[string]interface{}{
		"Title":            "Editar Concepto",
		"CurrentPage":      "conceptos",
		"SessionData":      sessionData,
		"Tipo":             tipo,
		"Conceptos":        conceptos,
		"Iconos":           iconos,
		"ValidationErrors": errors,
		"ModoEdicion":      true,
		"ConceptoEditar":   concepto,
		"FormData":         formData,
	}

	utils.RenderTemplate(w, "dashboard", "concepto/conceptos", data)
}

// FNConcepto-recargarPaginaConErrores
// recargarPaginaConErrores recarga la página principal de conceptos mostrando errores
// Parámetros:
// - w: ResponseWriter para enviar respuesta HTTP
// - r: Request HTTP original
// - sessionData: Datos de sesión del usuario
// - errors: Mapa de errores de validación
// - formData: Datos del formulario para repoblar campos
// Uso: Reutilización de vista principal cuando hay errores en creación
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

// FNConcepto-getIconosDisponibles
// getIconosDisponibles retorna la lista completa de íconos disponibles en el sistema
// Retorno: Slice de mapas con ID, Nombre y Path (clase FontAwesome) de cada ícono
// Uso: Poblar selector de íconos en formularios de concepto
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
		{"ID": 16, "Nombre": "Electricidad", "Path": "fa-solid fa-plug"},
		{"ID": 17, "Nombre": "Agua", "Path": "fa-solid fa-droplet"},
		{"ID": 18, "Nombre": "Internet", "Path": "fa-solid fa-wifi"},
		{"ID": 19, "Nombre": "Teléfono", "Path": "fa-solid fa-mobile-alt"},
		{"ID": 20, "Nombre": "Gas", "Path": "fa-solid fa-fire-flame-simple"},
		{"ID": 21, "Nombre": "Taxi", "Path": "fa-solid fa-taxi"},
		{"ID": 22, "Nombre": "Mantenimiento", "Path": "fa-solid fa-hammer"},
		{"ID": 23, "Nombre": "Gasolina", "Path": "fa-solid fa-oil-can"},
		{"ID": 24, "Nombre": "Cine", "Path": "fa-solid fa-clapperboard"},
		{"ID": 25, "Nombre": "Gimnasio", "Path": "fa-solid fa-weight-hanging"},
		{"ID": 26, "Nombre": "Farmacia", "Path": "fa-solid fa-pills"},
		{"ID": 27, "Nombre": "Libros", "Path": "fa-solid fa-book-open"},
		{"ID": 28, "Nombre": "Alquiler", "Path": "fa-solid fa-building"},
		{"ID": 29, "Nombre": "Dividendos", "Path": "fa-solid fa-chart-bar"},
		{"ID": 30, "Nombre": "Freelance", "Path": "fa-solid fa-laptop-code"},
		{"ID": 31, "Nombre": "Bonos", "Path": "fa-solid fa-star"},
		{"ID": 32, "Nombre": "Comisiones", "Path": "fa-solid fa-percent"},
	}
}

// FNConcepto-getIconoPorID
// getIconoPorID busca y retorna el path de un ícono basado en su ID
// Parámetros:
// - id: Identificador numérico del ícono
// Retorno: String con la clase FontAwesome del ícono, o "fa-solid fa-circle" por defecto
// Uso: Conversión de ID seleccionado a clase CSS para renderizado
func (c *ConceptoController) getIconoPorID(id int) string {
	iconos := c.getIconosDisponibles()
	for _, icono := range iconos {
		if icono["ID"] == id {
			return icono["Path"].(string)
		}
	}
	return "fa-solid fa-circle"
}

// FNConcepto-Deshabilitar
// Deshabilitar maneja la deshabilitación de un concepto para el usuario actual
// Parámetros:
// - w: ResponseWriter para enviar respuesta HTTP
// - r: Request HTTP con datos del concepto a deshabilitar
// Flujo:
// 1. Verifica método POST y autenticación
// 2. Obtiene parámetros del formulario
// 3. Verifica existencia del concepto
// 4. Deshabilita el concepto para el usuario actual
// 5. Redirige con mensaje de éxito
// Uso: Desactivación individual de conceptos por usuario
func (c *ConceptoController) Deshabilitar(w http.ResponseWriter, r *http.Request) {
	log.Printf("🚫 ConceptoController.Deshabilitar llamado - Método: %s", r.Method)

	if r.Method != http.MethodPost {
		log.Printf("❌ Método no permitido: %s", r.Method)
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	sessionData, ok := utils.GetSessionData(r)
	if !ok {
		log.Printf("❌ No hay sesión en Deshabilitar")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Obtener parámetros
	if err := r.ParseForm(); err != nil {
		log.Printf("❌ Error parseando formulario: %v", err)
		http.Error(w, "Error procesando solicitud", http.StatusBadRequest)
		return
	}

	nombreConcepto := r.FormValue("nombreConcepto")
	tipo := r.FormValue("tipo")

	if nombreConcepto == "" {
		log.Printf("❌ Nombre de concepto vacío")
		http.Error(w, "Nombre de concepto requerido", http.StatusBadRequest)
		return
	}

	log.Printf("👤 Usuario %s deshabilitando concepto: %s", sessionData.NombreUsuario, nombreConcepto)

	// Verificar que el concepto existe
	concepto, err := models.ConceptoModelInstance.FindByNombre(nombreConcepto, sessionData.CorreoFamilia)
	if err != nil || concepto == nil {
		log.Printf("❌ Concepto no encontrado: %s", nombreConcepto)
		http.Error(w, "Concepto no encontrado", http.StatusNotFound)
		return
	}

	// Deshabilitar el concepto para el usuario actual
	err = models.PersonalizacionModelInstance.DeshabilitarParaUsuario(
		nombreConcepto,
		sessionData.CorreoFamilia,
		sessionData.NombreUsuario,
	)
	if err != nil {
		log.Printf("❌ Error deshabilitando concepto: %v", err)
		http.Error(w, "Error al deshabilitar el concepto", http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Concepto deshabilitado exitosamente: %s", nombreConcepto)

	// Redirigir con mensaje de éxito
	redirectURL := "/conceptos?tipo=" + tipo + "&success=concepto_deshabilitado"
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

// FNConcepto-Habilitar
// Habilitar maneja la habilitación de un concepto previamente deshabilitado
// Parámetros:
// - w: ResponseWriter para enviar respuesta HTTP
// - r: Request HTTP con datos del concepto a habilitar
// Flujo:
// 1. Verifica método POST y autenticación
// 2. Obtiene parámetros del formulario
// 3. Verifica existencia del concepto
// 4. Habilita el concepto para el usuario actual
// 5. Redirige con mensaje de éxito
// Uso: Reactivación individual de conceptos por usuario
func (c *ConceptoController) Habilitar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	sessionData, ok := utils.GetSessionData(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error procesando solicitud", http.StatusBadRequest)
		return
	}

	nombreConcepto := r.FormValue("nombreConcepto")
	tipo := r.FormValue("tipo")

	// Verificar que el concepto existe
	concepto, err := models.ConceptoModelInstance.FindByNombre(nombreConcepto, sessionData.CorreoFamilia)
	if err != nil || concepto == nil {
		log.Printf("❌ Concepto no encontrado: %s", nombreConcepto)
		http.Error(w, "Concepto no encontrado", http.StatusNotFound)
		return
	}

	// Habilitar el concepto para el usuario actual
	err = models.PersonalizacionModelInstance.HabilitarParaUsuario(
		nombreConcepto,
		sessionData.CorreoFamilia,
		sessionData.NombreUsuario,
	)
	if err != nil {
		log.Printf("❌ Error habilitando concepto: %v", err)
		http.Error(w, "Error al habilitar el concepto", http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Concepto habilitado exitosamente: %s", nombreConcepto)
	redirectURL := "/conceptos?tipo=" + tipo + "&success=concepto_habilitado"
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
