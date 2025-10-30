// controllers/perfil_controller.go
package controllers

import (
	"arka-code/models"
	"arka-code/utils"
	"net/http"
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
