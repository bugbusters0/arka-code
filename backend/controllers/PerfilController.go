package controllers

import (
	"arka-code/backend/commons"
	"arka-code/backend/config"
	"arka-code/backend/models"
	"html/template"
	"net/http"
)

type PerfilController struct {
	userModel *models.UserModel
}

func NewPerfilController() *PerfilController {
	return &PerfilController{
		userModel: models.NewUserModel(),
	}
}

func (c *PerfilController) SolicitarPerfiles(w http.ResponseWriter, r *http.Request) {
	familyID := commons.GetFamilyID(r)
	miembros, err := c.userModel.GetMiembrosByFamily(familyID)
	if err != nil {
		commons.SetError(r, "Error al cargar perfiles")
	}

	tmpl := template.Must(template.ParseFiles("templates/auth/seleccionar-perfil.html"))
	tmpl.Execute(w, map[string]interface{}{
		"Miembros": miembros,
	})
}

func (c *PerfilController) ConsultarVerificacion(w http.ResponseWriter, r *http.Request) {
	password := r.FormValue("password")
	miembroID := commons.StringToInt(r.FormValue("idMiembro"))

	// Verificar contraseña del usuario principal
	userID := commons.GetUserID(r)
	if !c.userModel.VerifyPassword(userID, password) {
		commons.SetError(r, "Contraseña incorrecta")
		http.Redirect(w, r, config.URLROOT+"/seleccionar-perfil", http.StatusFound)
		return
	}

	commons.SetMiembroID(w, r, miembroID)
	http.Redirect(w, r, config.URLROOT+"/concepto/gasto", http.StatusFound)
}
