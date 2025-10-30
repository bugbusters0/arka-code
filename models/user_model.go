package models

import (
	"arka-code/database"
	"arka-code/entities"
	"arka-code/utils"
	"database/sql"
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
	          FROM usuario WHERE correoFamilia = ? AND delete_at IS NULL
	          ORDER BY nombrePersonal`

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
	_, err := database.DB.Exec(query, rol, nombreUsuario)
	return err
}

func (m *UserModel) SoftDelete(nombreUsuario string) error {
	query := `UPDATE usuario SET delete_at = ? WHERE nombreUsuario = ?`
	_, err := database.DB.Exec(query, time.Now(), nombreUsuario)
	return err
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

	// Aquí necesitas una función para verificar la contraseña
	return utils.CheckPassword(contrasena, hashedPassword), nil
}
