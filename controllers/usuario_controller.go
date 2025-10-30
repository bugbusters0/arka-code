package controllers

type UsuarioController struct{}

var UsuarioControllerInstance = &UsuarioController{}

// // ShowSelect muestra lista de usuarios de la familia
// func (c *UsuarioController) ShowSelect(w http.ResponseWriter, r *http.Request) {
// 	correoFamilia := r.URL.Query().Get("familia")

// 	if correoFamilia == "" {
// 		http.Redirect(w, r, "/login", http.StatusSeeOther)
// 		return
// 	}

// 	// Obtener usuarios de la familia
// 	usuarios, err := models.UsuarioModelInstance.FindByCorreoFamilia(correoFamilia)
// 	if err != nil {
// 		http.Error(w, "Error al cargar usuarios", http.StatusInternalServerError)
// 		return
// 	}

// 	data := map[string]interface{}{
// 		"Title":         "Seleccionar Usuario",
// 		"CorreoFamilia": correoFamilia,
// 		"Usuarios":      usuarios,
// 		"Errors":        map[string]string{},
// 	}

// 	utils.RenderTemplate(w, "auth", "usuario/select", data)
// }

// // Login procesa el login de usuario
// func (c *UsuarioController) Login(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Redirect(w, r, "/login", http.StatusSeeOther)
// 		return
// 	}

// 	nombreUsuario := r.FormValue("nombre_usuario")
// 	contraseñaPersonal := r.FormValue("contraseña_personal")
// 	correoFamilia := r.FormValue("correo_familia")

// 	if nombreUsuario == "" || contraseñaPersonal == "" {
// 		http.Redirect(w, r, "/usuario/select?familia="+correoFamilia+"&error=campos_vacios", http.StatusSeeOther)
// 		return
// 	}

// 	// Buscar usuario
// 	usuario, err := models.UsuarioModelInstance.FindByNombreUsuario(nombreUsuario)
// 	if err != nil {
// 		http.Redirect(w, r, "/usuario/select?familia="+correoFamilia+"&error=error_servidor", http.StatusSeeOther)
// 		return
// 	}

// 	if usuario == nil || !utils.CheckPassword(contraseñaPersonal, usuario.ContraseñaPersonal) {
// 		http.Redirect(w, r, "/usuario/select?familia="+correoFamilia+"&error=credenciales_incorrectas", http.StatusSeeOther)
// 		return
// 	}

// 	// Verificar que pertenece a la familia
// 	if usuario.CorreoFamilia != correoFamilia {
// 		http.Redirect(w, r, "/usuario/select?familia="+correoFamilia+"&error=usuario_no_pertenece", http.StatusSeeOther)
// 		return
// 	}

// 	// Crear sesión
// 	err = utils.SetUserSession(w, r, usuario)
// 	if err != nil {
// 		http.Redirect(w, r, "/usuario/select?familia="+correoFamilia+"&error=error_sesion", http.StatusSeeOther)
// 		return
// 	}

// 	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
// }

// // ShowCreate muestra formulario para crear usuario
// func (c *UsuarioController) ShowCreate(w http.ResponseWriter, r *http.Request) {
// 	correoFamilia := r.URL.Query().Get("familia")

// 	if correoFamilia == "" {
// 		http.Redirect(w, r, "/login", http.StatusSeeOther)
// 		return
// 	}

// 	data := map[string]interface{}{
// 		"Title":         "Crear Usuario",
// 		"CorreoFamilia": correoFamilia,
// 		"Errors":        map[string]string{},
// 	}

// 	utils.RenderTemplate(w, "auth", "usuario/create", data)
// }

// // Create procesa la creación de usuario
// func (c *UsuarioController) Create(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Redirect(w, r, "/login", http.StatusSeeOther)
// 		return
// 	}

// 	correoFamilia := r.FormValue("correo_familia")
// 	validation := validators.UsuarioValidatorInstance.ValidateCreate(r)

// 	if !validation.Success {
// 		data := map[string]interface{}{
// 			"Title":          "Crear Usuario",
// 			"CorreoFamilia":  correoFamilia,
// 			"Errors":         validation.Errors,
// 			"NombreUsuario":  r.FormValue("nombre_usuario"),
// 			"NombrePersonal": r.FormValue("nombre_personal"),
// 		}
// 		utils.RenderTemplate(w, "auth", "usuario/create", data)
// 		return
// 	}

// 	nombreUsuario := r.FormValue("nombre_usuario")
// 	contraseñaPersonal := r.FormValue("contraseña_personal")
// 	nombrePersonal := r.FormValue("nombre_personal")

// 	// Verificar si el usuario ya existe
// 	exists, err := models.UsuarioModelInstance.ExistsByNombreUsuario(nombreUsuario)
// 	if err != nil {
// 		data := map[string]interface{}{
// 			"Title":          "Crear Usuario",
// 			"CorreoFamilia":  correoFamilia,
// 			"Errors":         map[string]string{"general": "Error al verificar usuario"},
// 			"NombreUsuario":  nombreUsuario,
// 			"NombrePersonal": nombrePersonal,
// 		}
// 		utils.RenderTemplate(w, "auth", "usuario/create", data)
// 		return
// 	}

// 	if exists {
// 		data := map[string]interface{}{
// 			"Title":          "Crear Usuario",
// 			"CorreoFamilia":  correoFamilia,
// 			"Errors":         map[string]string{"nombre_usuario": "Este nombre de usuario ya está en uso"},
// 			"NombreUsuario":  nombreUsuario,
// 			"NombrePersonal": nombrePersonal,
// 		}
// 		utils.RenderTemplate(w, "auth", "usuario/create", data)
// 		return
// 	}

// 	// Hash password
// 	hashedPassword, err := utils.HashPassword(contraseñaPersonal)
// 	if err != nil {
// 		data := map[string]interface{}{
// 			"Title":          "Crear Usuario",
// 			"CorreoFamilia":  correoFamilia,
// 			"Errors":         map[string]string{"general": "Error al procesar contraseña"},
// 			"NombreUsuario":  nombreUsuario,
// 			"NombrePersonal": nombrePersonal,
// 		}
// 		utils.RenderTemplate(w, "auth", "usuario/create", data)
// 		return
// 	}

// 	// Determinar rol: primer usuario es admin (1), los demás member (0)
// 	countAdmins, _ := models.UsuarioModelInstance.CountAdminsByFamilia(correoFamilia)
// 	rol := int8(0) // member por defecto
// 	if countAdmins == 0 {
// 		rol = 1 // primer usuario es admin
// 	}

// 	usuario := &entities.Usuario{
// 		NombreUsuario:      nombreUsuario,
// 		Rol:                rol,
// 		ContraseñaPersonal: hashedPassword,
// 		NombrePersonal:     nombrePersonal,
// 		CorreoFamilia:      correoFamilia,
// 	}

// 	err = models.UsuarioModelInstance.Create(usuario)
// 	if err != nil {
// 		data := map[string]interface{}{
// 			"Title":          "Crear Usuario",
// 			"CorreoFamilia":  correoFamilia,
// 			"Errors":         map[string]string{"general": "Error al crear usuario"},
// 			"NombreUsuario":  nombreUsuario,
// 			"NombrePersonal": nombrePersonal,
// 		}
// 		utils.RenderTemplate(w, "auth", "usuario/create", data)
// 		return
// 	}

// 	http.Redirect(w, r, "/usuario/select?familia="+correoFamilia+"&success=usuario_creado", http.StatusSeeOther)
// }

// // Logout cierra sesión
// func (c *UsuarioController) Logout(w http.ResponseWriter, r *http.Request) {
// 	utils.ClearSession(w, r)
// 	http.Redirect(w, r, "/login", http.StatusSeeOther)
// }

// // UpdateRol actualiza el rol de un usuario (solo admin)
// func (c *UsuarioController) UpdateRol(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Redirect(w, r, "/admin/usuarios", http.StatusSeeOther)
// 		return
// 	}

// 	sessionData, ok := utils.GetSessionData(r)
// 	if !ok || !sessionData.IsAdmin() {
// 		http.Error(w, "Acceso denegado", http.StatusForbidden)
// 		return
// 	}

// 	nombreUsuario := r.FormValue("nombre_usuario")
// 	nuevoRolStr := r.FormValue("nuevo_rol")

// 	nuevoRol, err := strconv.Atoi(nuevoRolStr)
// 	if err != nil || (nuevoRol != 0 && nuevoRol != 1) {
// 		http.Redirect(w, r, "/admin/usuarios?error=rol_invalido", http.StatusSeeOther)
// 		return
// 	}

// 	// No permitir que el admin se quite a sí mismo el rol de admin
// 	if nombreUsuario == sessionData.NombreUsuario && nuevoRol == 0 {
// 		http.Redirect(w, r, "/admin/usuarios?error=no_puedes_cambiar_tu_rol", http.StatusSeeOther)
// 		return
// 	}

// 	err = models.UsuarioModelInstance.UpdateRol(nombreUsuario, int8(nuevoRol))
// 	if err != nil {
// 		http.Redirect(w, r, "/admin/usuarios?error=error_actualizando", http.StatusSeeOther)
// 		return
// 	}

// 	http.Redirect(w, r, "/admin/usuarios?success=rol_actualizado", http.StatusSeeOther)
// }

// // SoftDelete elimina un usuario (soft delete)
// func (c *UsuarioController) SoftDelete(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Redirect(w, r, "/admin/usuarios", http.StatusSeeOther)
// 		return
// 	}

// 	sessionData, ok := utils.GetSessionData(r)
// 	if !ok || !sessionData.IsAdmin() {
// 		http.Error(w, "Acceso denegado", http.StatusForbidden)
// 		return
// 	}

// 	nombreUsuario := r.FormValue("nombre_usuario")

// 	// No permitir que el admin se elimine a sí mismo
// 	if nombreUsuario == sessionData.NombreUsuario {
// 		http.Redirect(w, r, "/admin/usuarios?error=no_puedes_eliminarte", http.StatusSeeOther)
// 		return
// 	}

// 	err := models.UsuarioModelInstance.SoftDelete(nombreUsuario)
// 	if err != nil {
// 		http.Redirect(w, r, "/admin/usuarios?error=error_eliminando", http.StatusSeeOther)
// 		return
// 	}

// 	http.Redirect(w, r, "/admin/usuarios?success=usuario_eliminado", http.StatusSeeOther)
// }
