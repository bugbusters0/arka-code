package controllers

import (
	"arka-code/utils"
	"net/http"
)

type FamiliaController struct{}

var FamiliaControllerInstance = &FamiliaController{}

// ShowRegister muestra formulario de registro de familia
func (c *FamiliaController) ShowRegister(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":  "Registrar Familia",
		"Errors": map[string]string{},
	}
	utils.RenderTemplate(w, "auth", "auth/registro", data)
}

// // Register procesa el registro de familia
// func (c *FamiliaController) Register(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Redirect(w, r, "/familia/registro", http.StatusSeeOther)
// 		return
// 	}

// 	validation := validators.FamiliaValidatorInstance.ValidateRegister(r)

// 	if !validation.Success {
// 		data := map[string]interface{}{
// 			"Title":    "Registrar Familia",
// 			"Errors":   validation.Errors,
// 			"Correo":   r.FormValue("correo"),
// 			"Telefono": r.FormValue("telefono"),
// 		}
// 		utils.RenderTemplate(w, "auth", "familia/registro", data)
// 		return
// 	}

// 	correo := r.FormValue("correo")
// 	telefono := r.FormValue("telefono")
// 	contraseña := r.FormValue("contraseña")

// 	// Verificar si ya existe
// 	exists, err := models.FamiliaModelInstance.ExistsByCorreo(correo)
// 	if err != nil {
// 		data := map[string]interface{}{
// 			"Title":  "Registrar Familia",
// 			"Errors": map[string]string{"general": "Error al verificar correo"},
// 			"Correo": correo,
// 		}
// 		utils.RenderTemplate(w, "auth", "familia/registro", data)
// 		return
// 	}

// 	if exists {
// 		data := map[string]interface{}{
// 			"Title":  "Registrar Familia",
// 			"Errors": map[string]string{"correo": "Este correo ya está registrado"},
// 			"Correo": correo,
// 		}
// 		utils.RenderTemplate(w, "auth", "familia/registro", data)
// 		return
// 	}

// 	// Hash password
// 	hashedPassword, err := utils.HashPassword(contraseña)
// 	if err != nil {
// 		data := map[string]interface{}{
// 			"Title":  "Registrar Familia",
// 			"Errors": map[string]string{"general": "Error al procesar contraseña"},
// 			"Correo": correo,
// 		}
// 		utils.RenderTemplate(w, "auth", "familia/registro", data)
// 		return
// 	}

// 	// Crear familia
// 	var telefonoPtr *string
// 	if telefono != "" {
// 		telefonoPtr = &telefono
// 	}

// 	familia := &entities.Familia{
// 		Correo:     correo,
// 		Telefono:   telefonoPtr,
// 		Contraseña: hashedPassword,
// 	}

// 	err = models.FamiliaModelInstance.Create(familia)
// 	if err != nil {
// 		data := map[string]interface{}{
// 			"Title":  "Registrar Familia",
// 			"Errors": map[string]string{"general": "Error al crear familia"},
// 			"Correo": correo,
// 		}
// 		utils.RenderTemplate(w, "auth", "familia/registro", data)
// 		return
// 	}

// 	// Redirigir al login con mensaje de éxito
// 	http.Redirect(w, r, "/login?success=familia_creada", http.StatusSeeOther)
// }

// // ShowLogin muestra formulario de login de familia
// func (c *FamiliaController) ShowLogin(w http.ResponseWriter, r *http.Request) {
// 	successMsg := ""
// 	if r.URL.Query().Get("success") == "familia_creada" {
// 		successMsg = "Familia registrada exitosamente. Ahora puedes iniciar sesión."
// 	}

// 	data := map[string]interface{}{
// 		"Title":      "Iniciar Sesión",
// 		"Errors":     map[string]string{},
// 		"SuccessMsg": successMsg,
// 	}
// 	utils.RenderTemplate(w, "auth", "familia/login", data)
// }

// // Login procesa el login de familia (admin)
// func (c *FamiliaController) Login(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Redirect(w, r, "/login", http.StatusSeeOther)
// 		return
// 	}

// 	validation := validators.FamiliaValidatorInstance.ValidateLogin(r)

// 	if !validation.Success {
// 		data := map[string]interface{}{
// 			"Title":  "Iniciar Sesión",
// 			"Errors": validation.Errors,
// 			"Correo": r.FormValue("correo"),
// 		}
// 		utils.RenderTemplate(w, "auth", "familia/login", data)
// 		return
// 	}

// 	correo := r.FormValue("correo")
// 	contraseña := r.FormValue("contraseña")

// 	// Buscar familia
// 	familia, err := models.FamiliaModelInstance.FindByCorreo(correo)
// 	if err != nil {
// 		data := map[string]interface{}{
// 			"Title":  "Iniciar Sesión",
// 			"Errors": map[string]string{"general": "Error al buscar familia"},
// 			"Correo": correo,
// 		}
// 		utils.RenderTemplate(w, "auth", "familia/login", data)
// 		return
// 	}

// 	if familia == nil || !utils.CheckPassword(contraseña, familia.Contraseña) {
// 		data := map[string]interface{}{
// 			"Title":  "Iniciar Sesión",
// 			"Errors": map[string]string{"general": "Credenciales incorrectas"},
// 			"Correo": correo,
// 		}
// 		utils.RenderTemplate(w, "auth", "familia/login", data)
// 		return
// 	}

// 	// Familia autenticada, ahora seleccionar usuario
// 	http.Redirect(w, r, "/usuario/select?familia="+correo, http.StatusSeeOther)
// }
