package controllers

import (
	"arka-code/backend/commons"
	"arka-code/backend/config"
	"arka-code/backend/models"
	"arka-code/backend/validator"
	"html/template"
	"net/http"
)

type AuthController struct {
	// Puedes agregar dependencias aquí
}

func NewAuthController() *AuthController {
	return &AuthController{}
}

func (c *AuthController) ShowLogin(w http.ResponseWriter, r *http.Request) {
	// Verificar si ya está autenticado
	if commons.IsAuthenticated(r) {
		if commons.GetUserType(r) == "family" {
			http.Redirect(w, r, config.URLROOT+"/seleccionar-perfil", http.StatusFound)
		} else {
			http.Redirect(w, r, config.URLROOT+"/concepto/gasto", http.StatusFound)
		}
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/auth/login.html"))
	tmpl.Execute(w, nil)
}

func (c *AuthController) EnviarCredenciales(w http.ResponseWriter, r *http.Request) {
	// Validar datos
	validation := validator.LoginValidator.Validate(r)
	if !validation.Success {
		commons.SetValidationErrors(r, validation.Errors)
		http.Redirect(w, r, config.URLROOT+"/login", http.StatusFound)
		return
	}

	// Autenticar usuario
	email := r.FormValue("email")
	password := r.FormValue("password")

	userModel := models.NewUserModel()
	user, err := userModel.Authenticate(email, password)
	if err != nil {
		commons.SetError(r, "Credenciales inválidas")
		http.Redirect(w, r, config.URLROOT+"/login", http.StatusFound)
		return
	}

	// Iniciar sesión
	commons.Login(w, r, user)

	// Redirigir según tipo
	if user.Type == "family" {
		http.Redirect(w, r, config.URLROOT+"/seleccionar-perfil", http.StatusFound)
	} else {
		commons.SetMiembroID(w, r, user.ID)
		http.Redirect(w, r, config.URLROOT+"/concepto/gasto", http.StatusFound)
	}
}

func (c *AuthController) ShowRegister(w http.ResponseWriter, r *http.Request) {
	if commons.IsAuthenticated(r) {
		http.Redirect(w, r, config.URLROOT+"/seleccionar-perfil", http.StatusFound)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/auth/registro.html"))
	tmpl.Execute(w, nil)
}

func (c *AuthController) ProcessRegister(w http.ResponseWriter, r *http.Request) {
	validation := validator.RegisterValidator.Validate(r)
	if !validation.Success {
		commons.SetValidationErrors(r, validation.Errors)
		http.Redirect(w, r, config.URLROOT+"/registro", http.StatusFound)
		return
	}

	userModel := models.NewUserModel()
	err := userModel.CreateUser(
		r.FormValue("nombre"),
		r.FormValue("email"),
		r.FormValue("password"),
		r.FormValue("family_name"),
	)

	if err != nil {
		commons.SetError(r, "Error al crear usuario: "+err.Error())
		http.Redirect(w, r, config.URLROOT+"/registro", http.StatusFound)
		return
	}

	http.Redirect(w, r, config.URLROOT+"/login?registro=exitoso", http.StatusFound)
}

func (c *AuthController) ProcessSeleccionarPerfil(w http.ResponseWriter, r *http.Request) {
	// Implementar según tu lógica
	profileID := r.URL.Query().Get("profile")
	commons.SetMiembroID(w, r, commons.StringToInt(profileID))
	http.Redirect(w, r, config.URLROOT+"/concepto/gasto", http.StatusFound)
}

func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	commons.Logout(w, r)
	http.Redirect(w, r, config.URLROOT+"/login", http.StatusFound)
}
