package controllers

import (
	"arka-code/entities"
	"arka-code/models"
	"arka-code/utils"
	"arka-code/validators"
	"log"
	"net/http"
	"strconv"
	"time"
)

type MovimientoController struct{}

var MovimientoControllerInstance = &MovimientoController{}

// Index muestra la página principal de movimientos (Entrada Diaria)
func (c *MovimientoController) Index(w http.ResponseWriter, r *http.Request) {
	sessionData, ok := utils.GetSessionData(r)
	if !ok {
		http.Redirect(w, r, "/seleccionar-perfil", http.StatusSeeOther)
		return
	}

	// Obtener tipo de vista (gastos, ingresos, resumen)
	tipo := r.URL.Query().Get("tipo")
	log.Printf(" El tipo es:  " + tipo)
	if tipo == "" || tipo == "{{.Tipo}}" {
		tipo = "gasto" // Por defecto mostrar gastos
	}

	// Obtener fecha seleccionada (por defecto hoy)
	fechaStr := r.URL.Query().Get("fecha")
	var fecha time.Time
	var err error

	if fechaStr != "" {
		fecha, err = time.Parse("2006-01-02", fechaStr)
		if err != nil {
			fecha = time.Now()
		}
	} else {
		fecha = time.Now()
	}

	// Preparar datos base - INICIALIZAR ValidationErrors como mapa vacío
	data := map[string]interface{}{
		"Title":            "Entrada Diaria", // CORREGIDO: Cambié "Title" por consistencia
		"CurrentPage":      "movimientos",
		"SessionData":      sessionData,
		"Tipo":             tipo,
		"FechaActual":      fecha.Format("2006-01-02"),
		"ValidationErrors": map[string]string{}, // INICIALIZAR COMO MAPA VACÍO
		"FormData":         map[string]interface{}{},
	}

	// Obtener conceptos según el tipo
	var conceptos []entities.Concepto
	if tipo == "resumen" {
		// Para resumen, obtener ambos tipos
		gastosConceptos, _ := models.MovimientoModelInstance.FindByFamilia(sessionData.CorreoFamilia, "gasto")
		ingresosConceptos, _ := models.MovimientoModelInstance.FindByFamilia(sessionData.CorreoFamilia, "ingreso")
		conceptos = append(gastosConceptos, ingresosConceptos...)
	} else {
		conceptos, err = models.MovimientoModelInstance.FindByFamilia(sessionData.CorreoFamilia, tipo)
		if err != nil {
			log.Printf("❌ Error obteniendo conceptos: %v", err)
			conceptos = []entities.Concepto{}
		}
	}

	// Obtener movimientos según el tipo
	var movimientos []entities.Movimiento
	var totalGastos, totalIngresos float64

	if tipo == "resumen" {
		// Obtener todos los movimientos del día
		movimientos, err = models.MovimientoModelInstance.FindByUsuarioAndDate(sessionData.NombreUsuario, sessionData.CorreoFamilia, fecha)

		if err != nil {
			log.Printf("❌ Error obteniendo movimientos: %v", err)
			movimientos = []entities.Movimiento{}
		}

		// Calcular totales
		totalGastos, totalIngresos, _ = models.MovimientoModelInstance.GetTotalesByUsuarioAndDate(
			sessionData.NombreUsuario, sessionData.CorreoFamilia, fecha)
		balance := totalIngresos - totalGastos
		data["Balance"] = balance

		// Separar movimientos por tipo
		var gastosMovimientos []entities.Movimiento
		var ingresosMovimientos []entities.Movimiento

		for _, mov := range movimientos {
			concepto, _ := models.ConceptoModelInstance.FindByNombre(mov.NombreConcepto, sessionData.CorreoFamilia)
			if concepto != nil {
				if concepto.Tipo == 0 {
					gastosMovimientos = append(gastosMovimientos, mov)
				} else {
					ingresosMovimientos = append(ingresosMovimientos, mov)
				}
			}
		}

		data["MovimientosGastos"] = gastosMovimientos
		data["MovimientosIngresos"] = ingresosMovimientos
	} else {
		// Obtener movimientos filtrados por tipo
		tipoInt := int8(0)
		if tipo == "ingreso" {
			tipoInt = 1
		}

		movimientos, err = models.MovimientoModelInstance.FindByUsuarioDateAndTipo(
			sessionData.NombreUsuario, sessionData.CorreoFamilia, fecha, tipoInt)
		if err != nil {
			log.Printf("❌ Error obteniendo movimientos: %v", err)
			movimientos = []entities.Movimiento{}
		}

		// Calcular total del tipo seleccionado
		totalGastos, totalIngresos, _ = models.MovimientoModelInstance.GetTotalesByUsuarioAndDate(
			sessionData.NombreUsuario, sessionData.CorreoFamilia, fecha)
	}

	// Agregar conceptos con sus datos completos (incluyendo icono y color)
	conceptosConDatos := []map[string]interface{}{}
	for _, concepto := range conceptos {
		conceptosConDatos = append(conceptosConDatos, map[string]interface{}{
			"NombreConcepto": concepto.NombreConcepto,
			"Icono":          concepto.Icono,
			"Color":          concepto.Color,
			"Tipo":           concepto.Tipo,
		})
	}

	// Crear mapa de conceptos para acceso rápido
	conceptosMap := make(map[string]map[string]interface{})
	for _, concepto := range conceptosConDatos {
		conceptosMap[concepto["NombreConcepto"].(string)] = concepto
	}

	data["ConceptosMap"] = conceptosMap
	data["Conceptos"] = conceptosConDatos
	data["Movimientos"] = movimientos
	data["TotalGastos"] = totalGastos
	data["TotalIngresos"] = totalIngresos

	if tipo == "gasto" {
		data["TotalDia"] = totalGastos
	} else if tipo == "ingreso" {
		data["TotalDia"] = totalIngresos
	}

	utils.RenderTemplate(w, "dashboard", "movimientos/index", data)
}

// FNCtrl_Movimiento_Crear
// Crear maneja la lógica para insertar un nuevo movimiento en el sistema a través de un formulario POST.
// Parámetros:
// - w: http.ResponseWriter para manejar la respuesta HTTP (redirecciones o renderizado de error).
// - r: *http.Request que contiene la solicitud, el método y los datos del formulario POST.
// Retorno: Redirecciona al dashboard con la nueva información o renderiza la vista de entrada con errores.
// Flujo:
// 1. **Verificación de Método:** Asegura que la solicitud sea **POST**. Si no, redirecciona a `/movimientos`.
// 2. **Verificación de Sesión:** Obtiene y valida los datos de sesión. Si falla, redirecciona a `/seleccionar-perfil`.
// 3. **Validación de Datos:**
//   - Llama al `MovimientoValidator` para procesar y validar los campos del formulario.
//   - Si la validación falla, registra el error y usa `renderWithError` para mostrar los mensajes de error al usuario.
//
// 4. **Construcción de Entidad:** Prepara la entidad `Movimiento` con los datos limpios y convertidos, incluyendo los datos de usuario y familia de la sesión.
// 5. **Persistencia:** Llama al método `Create` del `MovimientoModel` para guardar el registro en la base de datos.
//   - Si la creación falla, registra el error y usa `renderWithError`.
//
// 6. **Éxito:** Redirecciona al dashboard (`/movimientos`) con un indicador de éxito.
func (c *MovimientoController) Crear(w http.ResponseWriter, r *http.Request) {
	// Muestra el inicio del proceso en el log.
	log.Println("Inicio de Crear movimiento")

	// 1. Verificación de Método
	if r.Method != http.MethodPost {
		// Si el método no es POST, se registra el error y se redirige al dashboard.
		log.Println("Método no es POST")
		http.Redirect(w, r, "/movimientos", http.StatusSeeOther)
		return
	}

	log.Println("Método POST verificado")

	// 2. Verificación de Sesión
	sessionData, ok := utils.GetSessionData(r)
	if !ok {
		// Si no hay datos de sesión, se registra y se redirige a la selección de perfil.
		log.Println("❌ No hay sesión activa")
		http.Redirect(w, r, "/seleccionar-perfil", http.StatusSeeOther)
		return
	}

	// Se confirma la sesión activa y los datos de usuario/familia.
	log.Printf("Sesión verificada - Usuario: %s, Familia: %s", sessionData.NombreUsuario, sessionData.CorreoFamilia)

	// 3. Validación de Datos
	log.Println("Iniciando validación...")
	// Se invoca el validador para procesar los datos del formulario (r).
	validation := validators.MovimientoValidatorInstance.Validate(r)

	// Se registran los resultados de la validación.
	log.Printf("Resultado validación - Success: %v, Errors: %v", validation.Success, validation.Errors)
	log.Printf("CleanData: %+v", validation.CleanData)

	if !validation.Success {
		// Si la validación falló, se registra y se renderiza la vista con los errores.
		log.Println("Validación falló")
		// c.renderWithError es una función auxiliar para renderizar la vista con errores.
		c.renderWithError(w, r, sessionData, validation.Errors, r.FormValue("tipo"))
		return
	}

	log.Println("Validación exitosa")

	// 4. Construcción de Entidad
	var descripcion *string
	// Se verifica y procesa el campo opcional 'descripcion' de los datos limpios.
	if descVal, ok := validation.CleanData["descripcion"]; ok && descVal != nil {
		descStr := descVal.(string)
		if descStr != "" {
			descripcion = &descStr // Asigna un puntero a la descripción si no está vacía.
			log.Printf("Descripción: %s", descStr)
		} else {
			log.Println("Descripción vacía")
		}
	} else {
		log.Println("Sin descripción")
	}

	// Se construye el objeto Movimiento con los datos limpios y los datos de sesión.
	movimiento := &entities.Movimiento{
		Fecha:          validation.CleanData["fecha"].(time.Time),
		Monto:          validation.CleanData["monto"].(float64),
		Descripcion:    descripcion,
		NombreUsuario:  sessionData.NombreUsuario,
		NombreConcepto: validation.CleanData["nombreConcepto"].(string),
		CorreoFamilia:  sessionData.CorreoFamilia,
	}

	log.Printf("Intentando crear movimiento: %+v", movimiento)

	// 5. Persistencia
	// Se llama al modelo para crear el registro en la base de datos.
	err := models.MovimientoModelInstance.Create(movimiento)
	if err != nil {
		// Si hay un error de base de datos/modelo, se registra.
		log.Printf("Error creando movimiento: %v", err)
		// Se añade un error general y se renderiza la vista con el error.
		validation.Errors["general"] = "Error al crear el movimiento"
		c.renderWithError(w, r, sessionData, validation.Errors, r.FormValue("tipo"))
		return
	}

	log.Printf("Movimiento creado exitosamente - ID: %d", movimiento.IdMovimiento)

	// 6. Éxito
	// Se obtienen los parámetros de filtro del formulario para la redirección.
	tipo := r.FormValue("tipo")
	fecha := r.FormValue("fecha")
	// Redirección exitosa al dashboard, manteniendo los filtros de fecha y tipo, e indicando la creación.
	http.Redirect(w, r, "/movimientos?tipo="+tipo+"&fecha="+fecha+"&success=created", http.StatusSeeOther)
}

// FNCtrl_Movimiento_editar
// Editar maneja la lógica para actualizar un movimiento existente en el sistema a través de un formulario POST.
// Parámetros:
// - w: http.ResponseWriter para manejar la respuesta HTTP (redirecciones, error 403, o renderizado de error).
// - r: *http.Request que contiene la solicitud, el método y los datos del formulario POST (incluyendo el ID).
// Retorno: Redirecciona al dashboard con mensaje de éxito o renderiza la vista de entrada con errores.
// Flujo:
// 1. **Verificación de Método y Sesión:** Asegura que la solicitud sea POST y que la sesión esté activa.
// 2. **Validación de Datos:** Llama al `MovimientoValidator` (usando `ValidateUpdate`) para validar los campos, incluyendo el ID del movimiento. Si falla, renderiza errores.
// 3. **Verificación de Existencia y Permisos:**
//   - Usa `FindByID` para verificar que el movimiento exista.
//   - Comprueba que el `CorreoFamilia` del movimiento coincida con el de la sesión. Si no, retorna error 403 (No Autorizado).
//
// 4. **Construcción de Entidad:** Crea una entidad `Movimiento` con los datos limpios y el ID, manteniendo los campos fijos como `NombreUsuario` y `CorreoFamilia` del movimiento existente.
// 5. **Persistencia:** Llama al método `Update` del `MovimientoModel` para actualizar el registro. Si falla, renderiza errores.
// 6. **Éxito:** Redirecciona al dashboard (`/movimientos`) con los parámetros de filtro y un indicador de éxito.
func (c *MovimientoController) Editar(w http.ResponseWriter, r *http.Request) {
	// Asegura que la solicitud sea POST; si no, redirige al dashboard.
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/movimientos", http.StatusSeeOther)
		return
	}

	// Verifica y obtiene los datos de sesión. Si falla, redirige al selector de perfil.
	sessionData, ok := utils.GetSessionData(r)
	if !ok {
		http.Redirect(w, r, "/seleccionar-perfil", http.StatusSeeOther)
		return
	}

	// Validar
	// 2. Validación de Datos (incluye ID de movimiento)
	validation := validators.MovimientoValidatorInstance.ValidateUpdate(r)
	if !validation.Success {
		// Si la validación falla, renderiza la vista con los errores.
		c.renderWithError(w, r, sessionData, validation.Errors, r.FormValue("tipo"))
		return
	}

	// Obtiene el ID validado del movimiento a actualizar.
	idMovimiento := validation.CleanData["idMovimiento"].(int)

	// 3. Verificación de Existencia y Permisos
	// Busca el movimiento en la DB para verificar su existencia y permisos.
	movimientoExistente, err := models.MovimientoModelInstance.FindByID(idMovimiento)
	if err != nil || movimientoExistente == nil {
		// Si no se encuentra o hay error en DB, se registra y se notifica al usuario.
		log.Printf("❌ Movimiento no encontrado: %d", idMovimiento)
		validation.Errors["general"] = "Movimiento no encontrado"
		c.renderWithError(w, r, sessionData, validation.Errors, r.FormValue("tipo"))
		return
	}

	// Verifica que el movimiento pertenezca a la familia del usuario activo.
	if movimientoExistente.CorreoFamilia != sessionData.CorreoFamilia {
		// Si no tiene permiso, registra y retorna error 403 (Forbidden).
		log.Printf("❌ Intento de editar movimiento de otra familia")
		http.Error(w, "No autorizado", http.StatusForbidden)
		return
	}

	// 4. Construcción de Entidad
	// Obtiene la descripción (ya validada).
	descripcion := validation.CleanData["descripcion"].(string)
	// Construye la entidad Movimiento con los nuevos datos y el ID existente.
	movimiento := &entities.Movimiento{
		IdMovimiento:   idMovimiento,
		Fecha:          validation.CleanData["fecha"].(time.Time),
		Monto:          validation.CleanData["monto"].(float64),
		Descripcion:    &descripcion,
		NombreConcepto: validation.CleanData["nombreConcepto"].(string),
		// Mantiene el usuario y familia originales, ya que estos campos no se modifican en la edición.
		NombreUsuario: movimientoExistente.NombreUsuario,
		CorreoFamilia: movimientoExistente.CorreoFamilia,
	}

	// Maneja el caso de descripción vacía (se debe guardar como NULL en la DB).
	if descripcion == "" {
		movimiento.Descripcion = nil
	}

	// 5. Persistencia
	// Llama al modelo para ejecutar la actualización.
	err = models.MovimientoModelInstance.Update(movimiento)
	if err != nil {
		// Si falla la actualización en DB, registra el error y notifica al usuario.
		log.Printf("❌ Error actualizando movimiento: %v", err)
		validation.Errors["general"] = "Error al actualizar el movimiento"
		c.renderWithError(w, r, sessionData, validation.Errors, r.FormValue("tipo"))
		return
	}

	log.Printf("✅ Movimiento actualizado exitosamente - ID: %d", idMovimiento)

	// 6. Éxito
	// Obtiene los parámetros de filtro del formulario.
	tipo := r.FormValue("tipo")
	fecha := r.FormValue("fecha")
	// Redirige al dashboard, indicando el éxito de la actualización.
	http.Redirect(w, r, "/movimientos?tipo="+tipo+"&fecha="+fecha+"&success=updated", http.StatusSeeOther)
}

// Eliminar elimina un movimiento (soft delete)
func (c *MovimientoController) Eliminar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/movimientos", http.StatusSeeOther)
		return
	}

	sessionData, ok := utils.GetSessionData(r)
	if !ok {
		http.Redirect(w, r, "/seleccionar-perfil", http.StatusSeeOther)
		return
	}

	// Obtener ID del movimiento
	idStr := r.FormValue("idMovimiento")
	idMovimiento, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("❌ ID de movimiento inválido: %s", idStr)
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Verificar que el movimiento existe y pertenece a la familia
	movimiento, err := models.MovimientoModelInstance.FindByID(idMovimiento)
	if err != nil || movimiento == nil {
		log.Printf("❌ Movimiento no encontrado: %d", idMovimiento)
		http.Error(w, "Movimiento no encontrado", http.StatusNotFound)
		return
	}

	if movimiento.CorreoFamilia != sessionData.CorreoFamilia {
		log.Printf("❌ Intento de eliminar movimiento de otra familia")
		http.Error(w, "No autorizado", http.StatusForbidden)
		return
	}

	// Eliminar movimiento
	err = models.MovimientoModelInstance.Delete(idMovimiento)
	if err != nil {
		log.Printf("❌ Error eliminando movimiento: %v", err)
		http.Error(w, "Error al eliminar", http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Movimiento eliminado exitosamente - ID: %d", idMovimiento)

	// Redirigir con mensaje de éxito
	tipo := r.FormValue("tipo")
	fecha := r.FormValue("fecha")
	http.Redirect(w, r, "/movimientos?tipo="+tipo+"&fecha="+fecha+"&success=deleted", http.StatusSeeOther)
}

// renderWithError helper para renderizar con errores
func (c *MovimientoController) renderWithError(w http.ResponseWriter, r *http.Request, sessionData *entities.SessionData, errors map[string]string, tipo string) {
	if tipo == "" {
		tipo = "gasto"
	}

	fechaStr := r.FormValue("fecha")
	fecha, err := time.Parse("2006-01-02", fechaStr)
	if err != nil {
		fecha = time.Now()
	}

	// Obtener conceptos
	conceptos, _ := models.ConceptoModelInstance.FindByFamilia(sessionData.CorreoFamilia, tipo)

	// Convertir conceptos a map con datos completos
	conceptosConDatos := []map[string]interface{}{}
	for _, concepto := range conceptos {
		conceptosConDatos = append(conceptosConDatos, map[string]interface{}{
			"NombreConcepto": concepto.NombreConcepto,
			"Icono":          concepto.Icono,
			"Color":          concepto.Color,
			"Tipo":           concepto.Tipo,
		})
	}

	// Obtener movimientos del USUARIO ACTUAL
	tipoInt := int8(0)
	if tipo == "ingreso" {
		tipoInt = 1
	}
	movimientos, _ := models.MovimientoModelInstance.FindByUsuarioDateAndTipo(
		sessionData.NombreUsuario, sessionData.CorreoFamilia, fecha, tipoInt)

	// Obtener totales del USUARIO ACTUAL
	totalGastos, totalIngresos, _ := models.MovimientoModelInstance.GetTotalesByUsuarioAndDate(
		sessionData.NombreUsuario, sessionData.CorreoFamilia, fecha)

	data := map[string]interface{}{
		"Title":            "Entrada Diaria",
		"CurrentPage":      "dashboard",
		"SessionData":      sessionData,
		"Tipo":             tipo,
		"FechaActual":      fecha.Format("2006-01-02"),
		"ValidationErrors": errors,
		"FormData": map[string]interface{}{
			"monto":          r.FormValue("monto"),
			"descripcion":    r.FormValue("descripcion"),
			"nombreConcepto": r.FormValue("nombreConcepto"),
		},
		"Conceptos":   conceptosConDatos,
		"Movimientos": movimientos,
	}

	// Agregar total según el tipo
	if tipo == "gasto" {
		data["TotalDia"] = totalGastos
	} else if tipo == "ingreso" {
		data["TotalDia"] = totalIngresos
	}

	utils.RenderTemplate(w, "dashboard", "movimientos/index", data)
}
