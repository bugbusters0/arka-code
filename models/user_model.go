package models

import (
	"arka-code/database"
	"arka-code/entities"
	"arka-code/utils"
	"database/sql"
	"log"
	"time"
)

type UserModel struct{}

var UserModelInstance = &UserModel{}

func (m *UserModel) Create(user *entities.Usuario) error {
	query := `INSERT INTO usuario (nombreUsuario, rol, contraseñaPersonal, nombrePersonal, correoFamilia)
	          VALUES (?, ?, ?, ?, ?)`

	_, err := database.DB.Exec(query,
		user.NombreUsuario,
		user.Rol,
		user.ContrasenaPersonal,
		user.NombrePersonal,
		user.CorreoFamilia)

	return err
}

func (m *UserModel) FindByEmail(email string) (*entities.Usuario, error) {
	query := `SELECT nombreUsuario, rol, contraseñaPersonal, nombrePersonal, correoFamilia, delete_at
	          FROM usuario WHERE correoFamilia = ? AND delete_at IS NULL`

	user := &entities.Usuario{}
	err := database.DB.QueryRow(query, email).Scan(
		&user.NombreUsuario,
		&user.Rol,
		&user.ContrasenaPersonal,
		&user.NombrePersonal,
		&user.CorreoFamilia,
		&user.DeleteAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return user, err
}

func (m *UserModel) FindByNombreUsuario(nombreUsuario string) (*entities.Usuario, error) {
	query := `SELECT nombreUsuario, rol, contraseñaPersonal, nombrePersonal, correoFamilia, delete_at
	          FROM usuario WHERE nombreUsuario = ? AND delete_at IS NULL`

	user := &entities.Usuario{}
	err := database.DB.QueryRow(query, nombreUsuario).Scan(
		&user.NombreUsuario,
		&user.Rol,
		&user.ContrasenaPersonal,
		&user.NombrePersonal,
		&user.CorreoFamilia,
		&user.DeleteAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return user, err
}

func (m *UserModel) GetAllByFamilia(correoFamilia string) ([]entities.Usuario, error) {
	query := `SELECT nombreUsuario, rol, contraseñaPersonal, nombrePersonal, correoFamilia, delete_at
	          FROM usuario WHERE correoFamilia = ?
	          ORDER BY delete_at IS NULL DESC, nombrePersonal`

	rows, err := database.DB.Query(query, correoFamilia)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	usuarios := []entities.Usuario{}
	for rows.Next() {
		var user entities.Usuario
		err := rows.Scan(
			&user.NombreUsuario,
			&user.Rol,
			&user.ContrasenaPersonal,
			&user.NombrePersonal,
			&user.CorreoFamilia,
			&user.DeleteAt,
		)
		if err != nil {
			return nil, err
		}
		usuarios = append(usuarios, user)
	}

	return usuarios, nil
}

func (m *UserModel) UpdateRol(nombreUsuario string, rol int8) error {
	query := `UPDATE usuario SET rol = ? WHERE nombreUsuario = ?`
	result, err := database.DB.Exec(query, rol, nombreUsuario)
	if err != nil {
		log.Printf("❌ Error actualizando rol: %v", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("✅ Rol actualizado - Usuario: %s, Filas afectadas: %d", nombreUsuario, rowsAffected)
	return nil
}

func (m *UserModel) UpdateNombrePersonal(nombreUsuario string, nuevoNombre string) error {
	query := `UPDATE usuario SET nombrePersonal = ? WHERE nombreUsuario = ?`
	result, err := database.DB.Exec(query, nuevoNombre, nombreUsuario)
	if err != nil {
		log.Printf("❌ Error actualizando nombre: %v", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("✅ Nombre actualizado - Usuario: %s, Filas afectadas: %d", nombreUsuario, rowsAffected)
	return nil
}

func (m *UserModel) UpdateNombreUsuario(viejoUsuario string, nuevoUsuario string) error {
	query := `UPDATE usuario SET nombreUsuario = ? WHERE nombreUsuario = ?`
	result, err := database.DB.Exec(query, nuevoUsuario, viejoUsuario)
	if err != nil {
		log.Printf("❌ Error actualizando nombre de usuario: %v", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("✅ Nombre de usuario actualizado - De: %s a %s, Filas afectadas: %d", viejoUsuario, nuevoUsuario, rowsAffected)
	return nil
}

func (m *UserModel) UpdatePassword(nombreUsuario string, nuevaPassword string) error {
	hashedPassword, err := utils.HashPassword(nuevaPassword)
	if err != nil {
		log.Printf("❌ Error hasheando contraseña: %v", err)
		return err
	}

	query := `UPDATE usuario SET contraseñaPersonal = ? WHERE nombreUsuario = ?`
	result, err := database.DB.Exec(query, hashedPassword, nombreUsuario)
	if err != nil {
		log.Printf("❌ Error actualizando contraseña: %v", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("✅ Contraseña actualizada - Usuario: %s, Filas afectadas: %d", nombreUsuario, rowsAffected)
	return nil
}

// Update actualiza múltiples campos del usuario
func (m *UserModel) Update(usuario *entities.Usuario) error {
	query := `UPDATE usuario 
	          SET nombrePersonal = ?, rol = ? 
	          WHERE nombreUsuario = ?`

	result, err := database.DB.Exec(query,
		usuario.NombrePersonal,
		usuario.Rol,
		usuario.NombreUsuario)

	if err != nil {
		log.Printf("❌ Error actualizando usuario: %v", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("✅ Usuario actualizado - Usuario: %s, Filas afectadas: %d", usuario.NombreUsuario, rowsAffected)
	return nil
}

func (m *UserModel) SoftDelete(nombreUsuario string) error {
	query := `UPDATE usuario SET delete_at = ? WHERE nombreUsuario = ?`
	result, err := database.DB.Exec(query, time.Now(), nombreUsuario)
	if err != nil {
		log.Printf("❌ Error deshabilitando usuario: %v", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("✅ Usuario deshabilitado - Usuario: %s, Filas afectadas: %d", nombreUsuario, rowsAffected)
	return nil
}

func (m *UserModel) CheckPasswordUser(nombreUsuario, contrasena string) (bool, error) {
	query := `SELECT contraseñaPersonal FROM usuario WHERE nombreUsuario = ? AND delete_at IS NULL`

	var hashedPassword string
	err := database.DB.QueryRow(query, nombreUsuario).Scan(&hashedPassword)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return utils.CheckPassword(contrasena, hashedPassword), nil
}

// Exists verifica si un usuario existe
func (m *UserModel) Exists(nombreUsuario string) (bool, error) {
	query := `SELECT COUNT(*) FROM usuario WHERE nombreUsuario = ? AND delete_at IS NULL`

	var count int
	err := database.DB.QueryRow(query, nombreUsuario).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
