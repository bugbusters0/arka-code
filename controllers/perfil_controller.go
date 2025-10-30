// controllers/perfil_controller.go
package controllers

import (
	"arka-code/entities"
	"arka-code/models"
	"arka-code/utils"
	"arka-code/validators"
	"log"
	"net/http"
	"strconv"
)

type PerfilController struct{}

var PerfilControllerInstance = &PerfilController{}

func (c *PerfilController) ShowSeleccionarPerfil(w http.ResponseWriter, r *http.Request) {
	// Verificar que la familia esté autenticada
	correoFamilia, ok := utils.GetCorreoFamiliaSession(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Obtener usuarios de la familia
	usuarios, err := models.UserModelInstance.GetAllByFamilia(correoFamilia)
	if err != nil {
		data := map[string]interface{}{
			"Title": "Seleccionar Perfil",
			"Error": "Error al cargar los perfiles",
		}
		utils.RenderTemplate(w, "auth", "auth/seleccionar-perfil", data)
		return
	}

	data := map[string]interface{}{
		"Title":    "Seleccionar Perfil",
		"Usuarios": usuarios,
	}
	utils.RenderTemplate(w, "auth", "auth/seleccionar-perfil", data)
}

func (c *PerfilController) SeleccionarPerfil(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/seleccionar-perfil", http.StatusSeeOther)
		return
	}

	nombreUsuario := r.FormValue("nombreUsuario")
	contraPersonal := r.FormValue("contraPersonal")

	// Obtener el correo de la familia de la sesión para poder cargar los usuarios nuevamente
	correoFamilia, ok := utils.GetCorreoFamiliaSession(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Buscar usuario
	usuario, err := models.UserModelInstance.FindByNombreUsuario(nombreUsuario)
	if err != nil || usuario == nil {
		// Cargar usuarios nuevamente antes de mostrar el error
		usuarios, _ := models.UserModelInstance.GetAllByFamilia(correoFamilia)
		data := map[string]interface{}{
			"Title":    "Seleccionar Perfil",
			"Error":    "Usuario no encontrado",
			"Usuarios": usuarios, // ¡IMPORTANTE: incluir los usuarios!
		}
		utils.RenderTemplate(w, "auth", "auth/seleccionar-perfil", data)
		return
	}

	// Verificar contraseña personal
	passwordMatch, err := models.UserModelInstance.CheckPasswordUser(nombreUsuario, contraPersonal)
	if err != nil || !passwordMatch {
		// Cargar usuarios nuevamente antes de mostrar el error
		usuarios, _ := models.UserModelInstance.GetAllByFamilia(correoFamilia)
		data := map[string]interface{}{
			"Title":    "Seleccionar Perfil",
			"Error":    "Contraseña personal incorrecta",
			"Usuarios": usuarios, // ¡IMPORTANTE: incluir los usuarios!
		}
		utils.RenderTemplate(w, "auth", "auth/seleccionar-perfil", data)
		return
	}

	// Verificar que el usuario no esté eliminado
	if usuario.DeleteAt != nil {
		// Cargar usuarios nuevamente antes de mostrar el error
		usuarios, _ := models.UserModelInstance.GetAllByFamilia(correoFamilia)
		data := map[string]interface{}{
			"Title":    "Seleccionar Perfil",
			"Error":    "Usuario inactivo",
			"Usuarios": usuarios, // ¡IMPORTANTE: incluir los usuarios!
		}
		utils.RenderTemplate(w, "auth", "auth/seleccionar-perfil", data)
		return
	}

	// Crear sesión de usuario
	err = utils.SetUserSession(w, r, usuario)
	if err != nil {
		// Cargar usuarios nuevamente antes de mostrar el error
		usuarios, _ := models.UserModelInstance.GetAllByFamilia(correoFamilia)
		data := map[string]interface{}{
			"Title":    "Seleccionar Perfil",
			"Error":    "Error al crear sesión",
			"Usuarios": usuarios, // ¡IMPORTANTE: incluir los usuarios!
		}
		utils.RenderTemplate(w, "auth", "auth/seleccionar-perfil", data)
		return
	}

	// Redirigir al dashboard
	http.Redirect(w, r, "/concepto-gasto", http.StatusSeeOther)
}

func (c *PerfilController) Index(w http.ResponseWriter, r *http.Request) {
	sessionData, ok := utils.GetSessionData(r)
	if !ok {
		http.Redirect(w, r, "/seleccionar-perfil", http.StatusSeeOther)
		return
	}

	var usuarios []entities.Usuario
	var err error

	// Si es admin, obtener todos los usuarios de la familia
	// Si es miembro, solo obtener su propio perfil
	if sessionData.IsAdmin() {
		usuarios, err = models.UserModelInstance.GetAllByFamilia(sessionData.CorreoFamilia)
		if err != nil {
			log.Printf("❌ Error obteniendo usuarios: %v", err)
			usuarios = []entities.Usuario{}
		}
	} else {
		// Solo obtener el perfil del miembro actual
		usuario, err := models.UserModelInstance.FindByNombreUsuario(sessionData.NombreUsuario)
		if err != nil || usuario == nil {
			log.Printf("❌ Error obteniendo perfil del miembro: %v", err)
			http.Error(w, "Error obteniendo perfil", http.StatusInternalServerError)
			return
		}
		usuarios = []entities.Usuario{*usuario}
	}

	// Obtener límites de cada usuario con información adicional
	usuariosConLimites := []map[string]interface{}{}
	for _, usuario := range usuarios {
		// Obtener personalizaciones del usuario
		personalizaciones, err := models.PersonalizacionModelInstance.GetByUsuario(usuario.NombreUsuario)
		if err != nil {
			log.Printf("❌ Error obteniendo personalizaciones de %s: %v", usuario.NombreUsuario, err)
			personalizaciones = []entities.PersonalizacionConcepto{}
		}

		// Convertir personalizaciones a estructura con consumo
		limites := []map[string]interface{}{}
		for _, p := range personalizaciones {
			// Obtener concepto para icono
			concepto, _ := models.ConceptoModelInstance.FindByNombre(p.NombreConcepto, sessionData.CorreoFamilia)
			icono := "fa-circle"
			color := "#6b7280"
			if concepto != nil {
				icono = *concepto.Icono
				color = *concepto.Color
			}

			// Calcular consumo actual
			periodo := "mensual"
			if p.TipoPeriodoLimite != nil {
				periodo = *p.TipoPeriodoLimite
			}

			consumo, _ := models.PersonalizacionModelInstance.GetConsumoActual(
				usuario.NombreUsuario,
				p.NombreConcepto,
				periodo,
			)

			// Obtener nombre del usuario que asignó
			asignadoPor, _ := models.PersonalizacionModelInstance.GetUsuarioQueAsigno(p.IdPersonalizacion)

			limite := map[string]interface{}{
				"IdPersonalizacion": p.IdPersonalizacion,
				"Categoria":         p.NombreConcepto,
				"Icon":              icono,
				"Color":             color,
				"Periodo":           periodo,
				"Limite":            0.0,
				"Consumo":           consumo,
				"AsignadoPor":       asignadoPor,
			}

			if p.LimiteGasto != nil {
				limite["Limite"] = *p.LimiteGasto
			}

			limites = append(limites, limite)
		}

		usuarioData := map[string]interface{}{
			"NombreUsuario":  usuario.NombreUsuario,
			"NombrePersonal": usuario.NombrePersonal,
			"Rol":            usuario.Rol,
			"CorreoFamilia":  usuario.CorreoFamilia,
			"Activo":         usuario.DeleteAt == nil,
			"Limites":        limites,
		}

		usuariosConLimites = append(usuariosConLimites, usuarioData)
	}

	// Obtener conceptos disponibles para crear límites
	gastosConceptos, _ := models.ConceptoModelInstance.FindByFamilia(sessionData.CorreoFamilia, "gasto")
	ingresosConceptos, _ := models.ConceptoModelInstance.FindByFamilia(sessionData.CorreoFamilia, "ingreso")
	todosConceptos := append(gastosConceptos, ingresosConceptos...)

	data := map[string]interface{}{
		"Title":            "Perfiles",
		"CurrentPage":      "perfiles",
		"SessionData":      sessionData,
		"Usuarios":         usuariosConLimites,
		"Conceptos":        todosConceptos,
		"ValidationErrors": map[string]string{},
		"Success":          r.URL.Query().Get("success"),
		"Error":            r.URL.Query().Get("error"),
	}

	utils.RenderTemplate(w, "dashboard", "perfiles/index", data)
}

// Crear crea un nuevo miembro (solo admin)
func (c *PerfilController) Crear(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/perfiles", http.StatusSeeOther)
		return
	}

	sessionData, ok := utils.GetSessionData(r)
	if !ok || !sessionData.IsAdmin() {
		http.Error(w, "No autorizado", http.StatusForbidden)
		return
	}

	// Validar datos
	validation := validators.PerfilValidatorInstance.ValidateCrearMiembro(r)
	if !validation.Success {
		log.Printf("❌ Validación fallida: %+v", validation.Errors)
		c.redirectWithError(w, validation.Errors)
		return
	}

	// Verificar si el usuario ya existe
	nombreUsuario := validation.CleanData["nombreUsuario"].(string)
	exists, err := models.UserModelInstance.Exists(nombreUsuario)
	if err != nil {
		log.Printf("❌ Error verificando existencia de usuario: %v", err)
		c.redirectWithError(w, map[string]string{"general": "Error verificando usuario"})
		return
	}

	if exists {
		log.Printf("❌ Usuario ya existe: %s", nombreUsuario)
		c.redirectWithError(w, map[string]string{"nombreUsuario": "El nombre de usuario ya existe"})
		return
	}

	// Hash de contraseña (si se proporcionó)
	contrasenaHash := ""
	if contrasena, ok := validation.CleanData["contrasena"].(string); ok {
		hashedPass, err := utils.HashPassword(contrasena)
		if err != nil {
			log.Printf("❌ Error hasheando contraseña: %v", err)
			c.redirectWithError(w, map[string]string{"general": "Error procesando contraseña"})
			return
		}
		contrasenaHash = hashedPass
	} else {
		// Si no se proporcionó contraseña, generar una temporal
		hashedPass, _ := utils.HashPassword("Temporal123!")
		contrasenaHash = hashedPass
	}

	// Crear usuario
	usuario := &entities.Usuario{
		NombreUsuario:      nombreUsuario,
		NombrePersonal:     validation.CleanData["nombrePersonal"].(string),
		Rol:                validation.CleanData["rol"].(int8),
		ContrasenaPersonal: contrasenaHash,
		CorreoFamilia:      sessionData.CorreoFamilia,
	}

	err = models.UserModelInstance.Create(usuario)
	if err != nil {
		log.Printf("❌ Error creando usuario: %v", err)
		c.redirectWithError(w, map[string]string{"general": "Error al crear el usuario"})
		return
	}

	log.Printf("✅ Usuario creado exitosamente: %s", nombreUsuario)
	http.Redirect(w, r, "/perfiles?success=created", http.StatusSeeOther)
}

// Editar actualiza un miembro existente
func (c *PerfilController) Editar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/perfiles", http.StatusSeeOther)
		return
	}

	sessionData, ok := utils.GetSessionData(r)
	if !ok {
		http.Redirect(w, r, "/seleccionar-perfil", http.StatusSeeOther)
		return
	}

	// Validar datos
	validation := validators.PerfilValidatorInstance.ValidateEditarMiembro(r)
	if !validation.Success {
		log.Printf("❌ Validación fallida: %+v", validation.Errors)
		c.redirectWithError(w, validation.Errors)
		return
	}

	nombreUsuarioActual := validation.CleanData["nombreUsuarioActual"].(string)

	// Verificar permisos
	// Admin puede editar cualquiera, miembro solo su propio perfil
	if !sessionData.IsAdmin() && nombreUsuarioActual != sessionData.NombreUsuario {
		http.Error(w, "No autorizado", http.StatusForbidden)
		return
	}

	// Obtener usuario actual
	usuario, err := models.UserModelInstance.FindByNombreUsuario(nombreUsuarioActual)
	if err != nil || usuario == nil {
		log.Printf("❌ Usuario no encontrado: %s", nombreUsuarioActual)
		c.redirectWithError(w, map[string]string{"general": "Usuario no encontrado"})
		return
	}

	// Actualizar nombre personal (si cambió)
	if nombrePersonal, ok := validation.CleanData["nombrePersonal"].(string); ok {
		if nombrePersonal != usuario.NombrePersonal {
			err = models.UserModelInstance.UpdateNombrePersonal(nombreUsuarioActual, nombrePersonal)
			if err != nil {
				log.Printf("❌ Error actualizando nombre: %v", err)
				c.redirectWithError(w, map[string]string{"general": "Error actualizando nombre"})
				return
			}
		}
	}

	// Actualizar nombre de usuario (si cambió y el admin lo permite)
	if sessionData.IsAdmin() {
		if nombreUsuarioNuevo, ok := validation.CleanData["nombreUsuario"].(string); ok {
			if nombreUsuarioNuevo != nombreUsuarioActual {
				// Verificar que el nuevo nombre no exista
				exists, _ := models.UserModelInstance.Exists(nombreUsuarioNuevo)
				if exists {
					c.redirectWithError(w, map[string]string{"nombreUsuario": "El nombre de usuario ya existe"})
					return
				}

				err = models.UserModelInstance.UpdateNombreUsuario(nombreUsuarioActual, nombreUsuarioNuevo)
				if err != nil {
					log.Printf("❌ Error actualizando nombre de usuario: %v", err)
					c.redirectWithError(w, map[string]string{"general": "Error actualizando usuario"})
					return
				}
				nombreUsuarioActual = nombreUsuarioNuevo // Actualizar para siguientes operaciones
			}
		}
	}

	// Actualizar rol (solo admin)
	if sessionData.IsAdmin() {
		if rol, ok := validation.CleanData["rol"].(int8); ok {
			if rol != usuario.Rol {
				err = models.UserModelInstance.UpdateRol(nombreUsuarioActual, rol)
				if err != nil {
					log.Printf("❌ Error actualizando rol: %v", err)
					c.redirectWithError(w, map[string]string{"general": "Error actualizando rol"})
					return
				}
			}
		}
	}

	// Actualizar contraseña (si se proporcionó)
	if contrasena, ok := validation.CleanData["contrasena"].(string); ok {
		err = models.UserModelInstance.UpdatePassword(nombreUsuarioActual, contrasena)
		if err != nil {
			log.Printf("❌ Error actualizando contraseña: %v", err)
			c.redirectWithError(w, map[string]string{"general": "Error actualizando contraseña"})
			return
		}
	}

	log.Printf("✅ Usuario actualizado exitosamente: %s", nombreUsuarioActual)
	http.Redirect(w, r, "/perfiles?success=updated", http.StatusSeeOther)
}

// Deshabilitar desactiva un miembro (solo admin)
func (c *PerfilController) Deshabilitar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/perfiles", http.StatusSeeOther)
		return
	}

	sessionData, ok := utils.GetSessionData(r)
	if !ok || !sessionData.IsAdmin() {
		http.Error(w, "No autorizado", http.StatusForbidden)
		return
	}

	nombreUsuario := r.FormValue("nombreUsuario")
	if nombreUsuario == "" {
		c.redirectWithError(w, map[string]string{"general": "Usuario no especificado"})
		return
	}

	// No permitir que el admin se deshabilite a sí mismo
	if nombreUsuario == sessionData.NombreUsuario {
		c.redirectWithError(w, map[string]string{"general": "No puedes deshabilitarte a ti mismo"})
		return
	}

	err := models.UserModelInstance.SoftDelete(nombreUsuario)
	if err != nil {
		log.Printf("❌ Error deshabilitando usuario: %v", err)
		c.redirectWithError(w, map[string]string{"general": "Error al deshabilitar usuario"})
		return
	}

	log.Printf("✅ Usuario deshabilitado exitosamente: %s", nombreUsuario)
	http.Redirect(w, r, "/perfiles?success=disabled", http.StatusSeeOther)
}

// CrearLimite crea un nuevo límite de gasto para un usuario
func (c *PerfilController) CrearLimite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/perfiles", http.StatusSeeOther)
		return
	}

	sessionData, ok := utils.GetSessionData(r)
	if !ok {
		http.Redirect(w, r, "/seleccionar-perfil", http.StatusSeeOther)
		return
	}

	// Validar datos
	validation := validators.PerfilValidatorInstance.ValidateLimite(r)
	if !validation.Success {
		log.Printf("❌ Validación de límite fallida: %+v", validation.Errors)
		c.redirectWithError(w, validation.Errors)
		return
	}

	nombreUsuarioDestino := validation.CleanData["nombreUsuario"].(string)

	// Verificar permisos: admin puede crear para cualquiera, miembro solo para sí mismo
	if !sessionData.IsAdmin() && nombreUsuarioDestino != sessionData.NombreUsuario {
		http.Error(w, "No autorizado", http.StatusForbidden)
		return
	}

	// Crear personalización
	limite := validation.CleanData["limite"].(float64)
	periodo := validation.CleanData["periodo"].(string)
	nombreConcepto := validation.CleanData["nombreConcepto"].(string)

	personalizacion := &entities.PersonalizacionConcepto{
		LimiteGasto:       &limite,
		Activo:            true,
		TipoPeriodoLimite: &periodo,
		Notificacion:      false,
		NombreUsuario:     nombreUsuarioDestino,
		NombreConcepto:    nombreConcepto,
		CorreoFamilia:     sessionData.CorreoFamilia,
	}

	err := models.PersonalizacionModelInstance.Create(personalizacion)
	if err != nil {
		log.Printf("❌ Error creando límite: %v", err)
		c.redirectWithError(w, map[string]string{"general": "Error al crear el límite"})
		return
	}

	log.Printf("✅ Límite creado exitosamente para %s", nombreUsuarioDestino)
	http.Redirect(w, r, "/perfiles?success=limit_created", http.StatusSeeOther)
}

// EditarLimite actualiza un límite existente
func (c *PerfilController) EditarLimite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/perfiles", http.StatusSeeOther)
		return
	}

	sessionData, ok := utils.GetSessionData(r)
	if !ok {
		http.Redirect(w, r, "/seleccionar-perfil", http.StatusSeeOther)
		return
	}

	// Obtener ID de la personalización
	idStr := r.FormValue("idPersonalizacion")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.redirectWithError(w, map[string]string{"general": "ID de límite inválido"})
		return
	}

	// Obtener personalización existente
	personalizacion, err := models.PersonalizacionModelInstance.GetByID(id)
	if err != nil || personalizacion == nil {
		log.Printf("❌ Límite no encontrado: %d", id)
		c.redirectWithError(w, map[string]string{"general": "Límite no encontrado"})
		return
	}

	// Verificar permisos
	if !sessionData.IsAdmin() && personalizacion.NombreUsuario != sessionData.NombreUsuario {
		http.Error(w, "No autorizado", http.StatusForbidden)
		return
	}

	// Actualizar campos
	limiteStr := r.FormValue("limite")
	if limiteStr != "" {
		limite, err := strconv.ParseFloat(limiteStr, 64)
		if err == nil && limite > 0 {
			personalizacion.LimiteGasto = &limite
		}
	}

	periodo := r.FormValue("periodo")
	if periodo != "" {
		personalizacion.TipoPeriodoLimite = &periodo
	}

	err = models.PersonalizacionModelInstance.Update(personalizacion)
	if err != nil {
		log.Printf("❌ Error actualizando límite: %v", err)
		c.redirectWithError(w, map[string]string{"general": "Error al actualizar el límite"})
		return
	}

	log.Printf("✅ Límite actualizado exitosamente: %d", id)
	http.Redirect(w, r, "/perfiles?success=limit_updated", http.StatusSeeOther)
}

// EliminarLimite elimina un límite de gasto
func (c *PerfilController) EliminarLimite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/perfiles", http.StatusSeeOther)
		return
	}

	sessionData, ok := utils.GetSessionData(r)
	if !ok {
		http.Redirect(w, r, "/seleccionar-perfil", http.StatusSeeOther)
		return
	}

	// Obtener ID de la personalización
	idStr := r.FormValue("idPersonalizacion")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.redirectWithError(w, map[string]string{"general": "ID de límite inválido"})
		return
	}

	// Obtener personalización para verificar permisos
	personalizacion, err := models.PersonalizacionModelInstance.GetByID(id)
	if err != nil || personalizacion == nil {
		log.Printf("❌ Límite no encontrado: %d", id)
		c.redirectWithError(w, map[string]string{"general": "Límite no encontrado"})
		return
	}

	// Verificar permisos
	if !sessionData.IsAdmin() && personalizacion.NombreUsuario != sessionData.NombreUsuario {
		http.Error(w, "No autorizado", http.StatusForbidden)
		return
	}

	err = models.PersonalizacionModelInstance.Delete(id)
	if err != nil {
		log.Printf("❌ Error eliminando límite: %v", err)
		c.redirectWithError(w, map[string]string{"general": "Error al eliminar el límite"})
		return
	}

	log.Printf("✅ Límite eliminado exitosamente: %d", id)
	http.Redirect(w, r, "/perfiles?success=limit_deleted", http.StatusSeeOther)
}

// Helper para redirigir con errores
func (c *PerfilController) redirectWithError(w http.ResponseWriter, errors map[string]string) {
	// Por simplicidad, solo redirigimos con el primer error
	for _, msg := range errors {
		http.Redirect(w, nil, "/perfiles?error="+msg, http.StatusSeeOther)
		return
	}
	http.Redirect(w, nil, "/perfiles?error=Error desconocido", http.StatusSeeOther)
}
