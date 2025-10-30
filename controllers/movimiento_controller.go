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
	if tipo == "" {
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

	// Preparar datos base
	data := map[string]interface{}{
		"Title":       "Entrada Diaria",
		"CurrentPage": "dashboard",
		"SessionData": sessionData,
		"Tipo":        tipo,
		"FechaActual": fecha.Format("2006-01-02"),
		"Errors":      map[string]string{},
		"FormData":    map[string]interface{}{},
	}

	// Obtener conceptos según el tipo
	var conceptos []entities.Concepto
	if tipo == "resumen" {
		// Para resumen, obtener ambos tipos
		gastosConceptos, _ := models.ConceptoModelInstance.FindByFamilia(sessionData.CorreoFamilia, "gasto")
		ingresosConceptos, _ := models.ConceptoModelInstance.FindByFamilia(sessionData.CorreoFamilia, "ingreso")
		conceptos = append(gastosConceptos, ingresosConceptos...)
	} else {
		conceptos, err = models.ConceptoModelInstance.FindByFamilia(sessionData.CorreoFamilia, tipo)
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
		movimientos, err = models.MovimientoModelInstance.FindByFamiliaAndDate(sessionData.CorreoFamilia, fecha)
		if err != nil {
			log.Printf("❌ Error obteniendo movimientos: %v", err)
			movimientos = []entities.Movimiento{}
		}

		// Calcular totales
		totalGastos, totalIngresos, _ = models.MovimientoModelInstance.GetTotalesByFamiliaAndDate(sessionData.CorreoFamilia, fecha)
		// CALCULAR BALANCE AQUÍ
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

		movimientos, err = models.MovimientoModelInstance.FindByFamiliaDateAndTipo(sessionData.CorreoFamilia, fecha, tipoInt)
		if err != nil {
			log.Printf("❌ Error obteniendo movimientos: %v", err)
			movimientos = []entities.Movimiento{}
		}

		// Calcular total del tipo seleccionado
		totalGastos, totalIngresos, _ = models.MovimientoModelInstance.GetTotalesByFamiliaAndDate(sessionData.CorreoFamilia, fecha)
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

// Crear crea un nuevo movimiento
func (c *MovimientoController) Crear(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/movimientos", http.StatusSeeOther)
		return
	}

	sessionData, ok := utils.GetSessionData(r)
	if !ok {
		http.Redirect(w, r, "/seleccionar-perfil", http.StatusSeeOther)
		return
	}

	// Validar
	validation := validators.MovimientoValidatorInstance.Validate(r)
	if !validation.Success {
		c.renderWithError(w, r, sessionData, validation.Errors, r.FormValue("tipo"))
		return
	}

	// Crear movimiento
	var descripcion *string
	if descVal, ok := validation.CleanData["descripcion"]; ok && descVal != nil {
		descStr := descVal.(string)
		if descStr != "" {
			descripcion = &descStr
		}
	}

	movimiento := &entities.Movimiento{
		Fecha:          validation.CleanData["fecha"].(time.Time),
		Monto:          validation.CleanData["monto"].(float64),
		Descripcion:    descripcion,
		NombreUsuario:  sessionData.NombreUsuario,
		NombreConcepto: validation.CleanData["nombreConcepto"].(string),
		CorreoFamilia:  sessionData.CorreoFamilia,
	}

	if descripcion == nil {
		movimiento.Descripcion = nil
	}

	err := models.MovimientoModelInstance.Create(movimiento)
	if err != nil {
		log.Printf("❌ Error creando movimiento: %v", err)
		validation.Errors["general"] = "Error al crear el movimiento"
		c.renderWithError(w, r, sessionData, validation.Errors, r.FormValue("tipo"))
		return
	}

	log.Printf("✅ Movimiento creado exitosamente - ID: %d", movimiento.IdMovimiento)

	// Redirigir con mensaje de éxito
	tipo := r.FormValue("tipo")
	fecha := r.FormValue("fecha")
	http.Redirect(w, r, "/movimientos?tipo="+tipo+"&fecha="+fecha+"&success=created", http.StatusSeeOther)
}

// Editar actualiza un movimiento existente
func (c *MovimientoController) Editar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/movimientos", http.StatusSeeOther)
		return
	}

	sessionData, ok := utils.GetSessionData(r)
	if !ok {
		http.Redirect(w, r, "/seleccionar-perfil", http.StatusSeeOther)
		return
	}

	// Validar
	validation := validators.MovimientoValidatorInstance.ValidateUpdate(r)
	if !validation.Success {
		c.renderWithError(w, r, sessionData, validation.Errors, r.FormValue("tipo"))
		return
	}

	idMovimiento := validation.CleanData["idMovimiento"].(int)

	// Verificar que el movimiento existe y pertenece a la familia
	movimientoExistente, err := models.MovimientoModelInstance.FindByID(idMovimiento)
	if err != nil || movimientoExistente == nil {
		log.Printf("❌ Movimiento no encontrado: %d", idMovimiento)
		validation.Errors["general"] = "Movimiento no encontrado"
		c.renderWithError(w, r, sessionData, validation.Errors, r.FormValue("tipo"))
		return
	}

	if movimientoExistente.CorreoFamilia != sessionData.CorreoFamilia {
		log.Printf("❌ Intento de editar movimiento de otra familia")
		http.Error(w, "No autorizado", http.StatusForbidden)
		return
	}

	// Actualizar movimiento
	descripcion := validation.CleanData["descripcion"].(string)
	movimiento := &entities.Movimiento{
		IdMovimiento:   idMovimiento,
		Fecha:          validation.CleanData["fecha"].(time.Time),
		Monto:          validation.CleanData["monto"].(float64),
		Descripcion:    &descripcion,
		NombreConcepto: validation.CleanData["nombreConcepto"].(string),
		NombreUsuario:  movimientoExistente.NombreUsuario, // Mantener usuario original
		CorreoFamilia:  movimientoExistente.CorreoFamilia,
	}

	if descripcion == "" {
		movimiento.Descripcion = nil
	}

	err = models.MovimientoModelInstance.Update(movimiento)
	if err != nil {
		log.Printf("❌ Error actualizando movimiento: %v", err)
		validation.Errors["general"] = "Error al actualizar el movimiento"
		c.renderWithError(w, r, sessionData, validation.Errors, r.FormValue("tipo"))
		return
	}

	log.Printf("✅ Movimiento actualizado exitosamente - ID: %d", idMovimiento)

	// Redirigir con mensaje de éxito
	tipo := r.FormValue("tipo")
	fecha := r.FormValue("fecha")
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

	// Obtener movimientos
	tipoInt := int8(0)
	if tipo == "ingreso" {
		tipoInt = 1
	}
	movimientos, _ := models.MovimientoModelInstance.FindByFamiliaDateAndTipo(sessionData.CorreoFamilia, fecha, tipoInt)

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

	utils.RenderTemplate(w, "dashboard", "movimientos/index", data)
}
