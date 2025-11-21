package controllers

import (
	"arka-code/entities"
	"arka-code/models"
	"arka-code/utils"
	"arka-code/validators"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type AuthController struct{}

var AuthControllerInstance = &AuthController{}

func (c *AuthController) ShowLogin(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":  "Iniciar Sesión",
		"Errors": map[string]string{},
	}
	utils.RenderTemplate(w, "auth", "auth/login", data)
}

func (c *AuthController) ShowRegister(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":  "Registrarse",
		"Errors": map[string]string{},
	}
	utils.RenderTemplate(w, "auth", "auth/registro", data)
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Buscar familia
	email := r.FormValue("email")
	password := r.FormValue("contrasena") // ← Cambiado a "contrasena" que es el name del input

	log.Printf("🔐 Intentando login para: %s", email)

	familia, err := models.FamiliaModelInstance.FindByEmail(email)
	if err != nil {
		log.Printf("❌ Error buscando familia: %v", err)
		data := map[string]interface{}{
			"Title":            "Iniciar Sesión",
			"ValidationErrors": map[string]string{"general": "Error al buscar familia"},
			"FormData": map[string]string{
				"email": email,
			},
		}
		utils.RenderTemplate(w, "auth", "auth/login", data)
		return
	}

	if familia == nil {
		log.Printf("❌ Familia no encontrada: %s", email)
		data := map[string]interface{}{
			"Title":            "Iniciar Sesión",
			"ValidationErrors": map[string]string{"general": "Credenciales incorrectas"},
			"FormData": map[string]string{
				"email": email,
			},
		}
		utils.RenderTemplate(w, "auth", "auth/login", data)
		return
	}

	passwordMatch, err := models.FamiliaModelInstance.CheckPasswordFamilia(email, password)
	if err != nil || !passwordMatch {
		log.Printf("❌ Contraseña incorrecta para: %s", email)
		data := map[string]interface{}{
			"Title":            "Iniciar Sesión",
			"ValidationErrors": map[string]string{"general": "Credenciales incorrectas"},
			"FormData": map[string]string{
				"email": email,
			},
		}
		utils.RenderTemplate(w, "auth", "auth/login", data)
		return
	}

	if familia.DeleteAt != nil {
		log.Printf("❌ Familia eliminada: %s", email)
		data := map[string]interface{}{
			"Title":            "Iniciar Sesión",
			"ValidationErrors": map[string]string{"general": "Cuenta inactiva"},
			"FormData": map[string]string{
				"email": email,
			},
		}
		utils.RenderTemplate(w, "auth", "auth/login", data)
		return
	}

	// Crear sesión de familia
	err = utils.SetFamiliaSession(w, r, familia)
	if err != nil {
		log.Printf("❌ Error creando sesión: %v", err)
		data := map[string]interface{}{
			"Title":            "Iniciar Sesión",
			"ValidationErrors": map[string]string{"general": "Error al crear sesión"},
			"FormData": map[string]string{
				"email": email,
			},
		}
		utils.RenderTemplate(w, "auth", "auth/login", data)
		return
	}

	log.Printf("✅ Login exitoso, redirigiendo a selección de perfil: %s", email)
	http.Redirect(w, r, "/seleccionar-perfil", http.StatusSeeOther)
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/registro", http.StatusSeeOther)
		return
	}

	// Validar
	validation := validators.RegisterValidatorInstance.Validate(r)
	if !validation.Success {
		// Preparar datos para mostrar errores de miembros en el template
		memberErrors := c.extractMemberErrors(validation.Errors)

		data := map[string]interface{}{
			"Title":            "Registrarse",
			"ValidationErrors": validation.Errors,
			"FieldErrors":      validation.Errors,
			"MemberErrors":     memberErrors, // ← Nuevo: errores específicos de miembros
			"FormData": map[string]string{
				"email":         r.FormValue("email"),
				"telefono":      r.FormValue("telefono"),
				"adminNombre":   r.FormValue("adminNombre"),
				"nombreUsuario": r.FormValue("nombreUsuario"),
			},
			"MemberData": c.extractMemberData(r), // ← Datos de miembros para repoblar
		}
		utils.RenderTemplate(w, "auth", "auth/registro", data)
		return
	}

	email := r.FormValue("email")
	telefono := r.FormValue("telefono")
	contrasena := r.FormValue("password")
	adminNombre := r.FormValue("adminNombre")
	contraPersonal := r.FormValue("contraPersonal")
	nombreUsuario := r.FormValue("nombreUsuario")

	// Verificar si el correo ya existe
	existing, err := models.FamiliaModelInstance.FindByEmail(email)
	if err != nil {
		data := map[string]interface{}{
			"Title":            "Registrarse",
			"ValidationErrors": map[string]string{"general": "Error al verificar email"},
			"FormData": map[string]string{
				"email":         email,
				"adminNombre":   adminNombre,
				"telefono":      telefono,
				"nombreUsuario": nombreUsuario,
			},
		}
		utils.RenderTemplate(w, "auth", "auth/registro", data)
		return
	}

	if existing != nil {
		data := map[string]interface{}{
			"Title":            "Registrarse",
			"ValidationErrors": map[string]string{"email": "El email ya está registrado"},
			"FormData": map[string]string{
				"email":         email,
				"adminNombre":   adminNombre,
				"telefono":      telefono,
				"nombreUsuario": nombreUsuario,
			},
		}
		utils.RenderTemplate(w, "auth", "auth/registro", data)
		return
	}
	existingUser, err := models.UserModelInstance.FindByNombreUsuario(nombreUsuario)
	if err != nil {
		data := map[string]interface{}{
			"Title":            "Registrarse",
			"ValidationErrors": map[string]string{"general": "Error al verificar nomrbe de usuario"},
			"FormData": map[string]string{
				"email":         email,
				"adminNombre":   adminNombre,
				"telefono":      telefono,
				"nombreUsuario": nombreUsuario,
			},
		}
		utils.RenderTemplate(w, "auth", "auth/registro", data)
		return
	}
	if existingUser != nil {
		data := map[string]interface{}{
			"Title":            "Registrarse",
			"ValidationErrors": map[string]string{"general": "El nombre de usuario ya está en uso"},
			"FormData": map[string]string{
				"email":         email,
				"adminNombre":   adminNombre,
				"telefono":      telefono,
				"nombreUsuario": nombreUsuario,
			},
		}
		utils.RenderTemplate(w, "auth", "auth/registro", data)
		return
	}
	memberErrors := c.verificarMiembrosAntesDeCrear(r)
	if len(memberErrors) > 0 {
		data := map[string]interface{}{
			"Title":            "Registrarse",
			"ValidationErrors": map[string]string{"general": "Error el crear miembros"},
			"FormData": map[string]string{
				"email":         email,
				"adminNombre":   adminNombre,
				"telefono":      telefono,
				"nombreUsuario": nombreUsuario,
			},
		}
		utils.RenderTemplate(w, "auth", "auth/registro", data)
		return
	}

	// Hash de la password de la familia
	hashedPassword, err := utils.HashPassword(contrasena)
	if err != nil {
		data := map[string]interface{}{
			"Title":            "Registrarse",
			"ValidationErrors": map[string]string{"general": "Error al procesar contraseña"},
			"FormData": map[string]string{
				"email":       email,
				"adminNombre": adminNombre,
				"telefono":    telefono,
			},
		}
		utils.RenderTemplate(w, "auth", "auth/registro", data)
		return
	}

	// Crear familia
	familia := &entities.Familia{
		Correo:     email,
		Telefono:   &telefono,
		Contraseña: hashedPassword,
	}

	err = models.FamiliaModelInstance.Create(familia)
	if err != nil {
		data := map[string]interface{}{
			"Title":            "Registrarse",
			"ValidationErrors": map[string]string{"general": "Error al crear la familia: " + err.Error()},
			"FormData": map[string]string{
				"email":       email,
				"adminNombre": adminNombre,
				"telefono":    telefono,
			},
		}
		utils.RenderTemplate(w, "auth", "auth/registro", data)
		return
	}

	// Hash de la contraseña personal del admin
	hashedContraPersonal, err := utils.HashPassword(contraPersonal)
	if err != nil {
		models.FamiliaModelInstance.Delete(email)
		data := map[string]interface{}{
			"Title":            "Registrarse",
			"ValidationErrors": map[string]string{"general": "Error al procesar contraseña personal"},
			"FormData": map[string]string{
				"email":       email,
				"adminNombre": adminNombre,
				"telefono":    telefono,
			},
		}
		utils.RenderTemplate(w, "auth", "auth/registro", data)
		return
	}

	// Crear usuario administrador
	adminUser := &entities.Usuario{
		CorreoFamilia:      email,
		NombreUsuario:      nombreUsuario,
		Rol:                1, // 1 = admin
		ContrasenaPersonal: hashedContraPersonal,
		NombrePersonal:     adminNombre,
	}

	err = models.UserModelInstance.Create(adminUser)
	if err != nil {

		models.FamiliaModelInstance.Delete(email)
		data := map[string]interface{}{
			"Title":            "Registrarse",
			"ValidationErrors": map[string]string{"general": "Error al crear el administrador: " + err.Error()},
			"FormData": map[string]string{
				"email":       email,
				"adminNombre": adminNombre,
				"telefono":    telefono,
			},
		}
		utils.RenderTemplate(w, "auth", "auth/registro", data)
		return
	}

	// NUEVO: Procesar miembros adicionales
	membersCreated, memberErrors := c.processAdditionalMembers(r, email)
	if len(memberErrors) > 0 {
		models.UserModelInstance.Delete(nombreUsuario) // Eliminar admin
		models.FamiliaModelInstance.Delete(email)      // Eliminar familia
		// Si hay errores al crear miembros, mostrar mensaje pero no fallar el registro completo
		data := map[string]interface{}{
			"Title":            "Registrarse",
			"ValidationErrors": map[string]string{"general": "Error al crear el miembros, ya exsite ese nombre de usuario: "},
			"FormData": map[string]string{
				"email":       email,
				"adminNombre": adminNombre,
				"telefono":    telefono,
			},
		}
		utils.RenderTemplate(w, "auth", "auth/registro", data)
		return
	}

	log.Printf("Registro completado - Familia: %s, Admin: %s, Miembros creados: %d",
		email, nombreUsuario, membersCreated)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (c *AuthController) verificarMiembrosAntesDeCrear(r *http.Request) []string {
	var errors []string
	memberUsernames := r.Form["memberUsername[]"]

	for i, username := range memberUsernames {
		username = strings.TrimSpace(username)
		if username == "" {
			continue // Saltar vacíos
		}

		// Verificar si el usuario ya existe
		exists, err := models.UserModelInstance.Exists(username)
		if err != nil {
			errors = append(errors, "Error verificando usuario: "+username)
			continue
		}
		if exists {
			errors = append(errors, "El usuario ya existe: "+username)
			continue
		}

		// Verificar duplicados en el mismo formulario
		for j, otherUsername := range memberUsernames {
			if i != j && username == strings.TrimSpace(otherUsername) {
				errors = append(errors, "Usuario duplicado: "+username)
				break
			}
		}
	}

	return errors
}

// NUEVA FUNCIÓN: Procesar miembros adicionales
func (c *AuthController) processAdditionalMembers(r *http.Request, correoFamilia string) (int, []string) {
	memberErrors := []string{}
	membersCreated := 0

	// Obtener arrays de miembros del formulario
	miembroNombres := r.Form["miembroNombre[]"]
	miembroUsuarios := r.Form["miembroUsuario[]"]
	miembroContras := r.Form["miembroContra[]"]
	miembroRoles := r.Form["miembroRol[]"]

	// Verificar que tenemos miembros para procesar
	if len(miembroNombres) == 0 {
		return 0, memberErrors
	}

	log.Printf("👥 Procesando %d miembros adicionales", len(miembroNombres))

	for i := 0; i < len(miembroNombres); i++ {
		// Saltar campos vacíos
		if miembroNombres[i] == "" || miembroUsuarios[i] == "" || miembroContras[i] == "" {
			continue
		}

		nombre := miembroNombres[i]
		usuario := miembroUsuarios[i]
		contra := miembroContras[i]
		rolStr := "miembro" // valor por defecto
		if i < len(miembroRoles) {
			rolStr = miembroRoles[i]
		}

		// Convertir rol a int8
		var rol int8 = 0 // 0 = miembro
		if rolStr == "admin" {
			rol = 1 // 1 = admin
		}

		// Verificar si el nombre de usuario ya existe
		existingUser, err := models.UserModelInstance.FindByNombreUsuario(usuario)
		if err != nil {
			memberErrors = append(memberErrors, fmt.Sprintf("Error verificando usuario %s: %v", usuario, err))
			continue
		}

		if existingUser != nil {
			memberErrors = append(memberErrors, fmt.Sprintf("El usuario miembro%s ya existe", usuario))
			continue
		}

		// Hash de la contraseña personal del miembro
		hashedContra, err := utils.HashPassword(contra)
		if err != nil {
			memberErrors = append(memberErrors, fmt.Sprintf("Error procesando contraseña para %s: %v", usuario, err))
			continue
		}

		// Crear usuario miembro
		memberUser := &entities.Usuario{
			CorreoFamilia:      correoFamilia,
			NombreUsuario:      usuario,
			Rol:                rol,
			ContrasenaPersonal: hashedContra,
			NombrePersonal:     nombre,
		}

		err = models.UserModelInstance.Create(memberUser)
		if err != nil {
			memberErrors = append(memberErrors, fmt.Sprintf("Error creando usuario %s: %v", usuario, err))
			continue
		}

		membersCreated++
		log.Printf("Miembro creado: %s (%s)", usuario, rolStr)
	}

	return membersCreated, memberErrors
}
func (c *AuthController) extractMemberErrors(allErrors map[string]string) map[int]map[string]string {
	memberErrors := make(map[int]map[string]string)

	for field, errorMsg := range allErrors {
		if strings.HasPrefix(field, "miembroNombre_") ||
			strings.HasPrefix(field, "miembroUsuario_") ||
			strings.HasPrefix(field, "miembroContra_") {

			// Extraer índice del miembro: "miembroNombre_0" → 0
			parts := strings.Split(field, "_")
			if len(parts) == 2 {
				index := parts[1]
				if idx, err := strconv.Atoi(index); err == nil {
					if memberErrors[idx] == nil {
						memberErrors[idx] = make(map[string]string)
					}

					fieldType := strings.TrimSuffix(parts[0], "miembro")
					memberErrors[idx][fieldType] = errorMsg
				}
			}
		}
	}

	return memberErrors
}

// NUEVA FUNCIÓN: Extraer datos de miembros para repoblar el formulario
func (c *AuthController) extractMemberData(r *http.Request) []map[string]string {
	memberData := []map[string]string{}

	miembroNombres := r.Form["miembroNombre[]"]
	miembroUsuarios := r.Form["miembroUsuario[]"]
	miembroContras := r.Form["miembroContra[]"]
	miembroRoles := r.Form["miembroRol[]"]

	for i := 0; i < len(miembroNombres); i++ {
		// Solo agregar si al menos un campo tiene datos
		if miembroNombres[i] != "" || miembroUsuarios[i] != "" || miembroContras[i] != "" {
			member := map[string]string{
				"nombre":  miembroNombres[i],
				"usuario": miembroUsuarios[i],
				"contra":  miembroContras[i],
				"rol":     "miembro", // valor por defecto
			}

			if i < len(miembroRoles) && miembroRoles[i] != "" {
				member["rol"] = miembroRoles[i]
			}

			memberData = append(memberData, member)
		}
	}

	return memberData
}

func (c *AuthController) LogoutFamilia(w http.ResponseWriter, r *http.Request) {
	utils.ClearSessionFamilia(w, r)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
func (c *AuthController) LogoutUsuario(w http.ResponseWriter, r *http.Request) {
	utils.ClearSessionUsuario(w, r)
	http.Redirect(w, r, "/seleccionar-perfil", http.StatusSeeOther)
}
