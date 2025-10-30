package models

import (
	"arka-code/database"
	"arka-code/entities"
	"arka-code/utils"
	"database/sql"
)

type FamiliaModel struct{}

var FamiliaModelInstance = &FamiliaModel{}

func (m *FamiliaModel) Create(familia *entities.Familia) error {
	query := `INSERT INTO familia (correo, telefono, contraseña)
	          VALUES (?, ?, ?)`

	_, err := database.DB.Exec(query,
		familia.Correo,
		familia.Telefono,
		familia.Contraseña)

	return err
}

func (m *FamiliaModel) FindByEmail(email string) (*entities.Familia, error) {
	query := `SELECT correo, telefono, contraseña, delete_at
	          FROM familia WHERE correo = ? AND delete_at IS NULL`

	familia := &entities.Familia{}
	err := database.DB.QueryRow(query, email).Scan(
		&familia.Correo,
		&familia.Telefono,
		&familia.Contraseña,
		&familia.DeleteAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return familia, err
}

func (m *FamiliaModel) GetAll() ([]entities.Familia, error) {
	query := `SELECT * FROM familia WHERE correo = ? AND delete_at IS NULL`

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	familias := []entities.Familia{}
	for rows.Next() {
		var familia entities.Familia
		err := rows.Scan(
			&familia.Correo,
			&familia.Telefono,
			&familia.Contraseña,
			&familia.DeleteAt,
		)
		if err != nil {
			return nil, err
		}
		familias = append(familias, familia)
	}

	return familias, nil
}

func (m *FamiliaModel) CheckPasswordFamilia(correo, contrasena string) (bool, error) {
	query := `SELECT contraseña FROM familia WHERE correo = ? AND delete_at IS NULL`

	var hashedPassword string
	err := database.DB.QueryRow(query, correo).Scan(&hashedPassword)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	// Aquí necesitas una función para verificar la contraseña
	return utils.CheckPassword(contrasena, hashedPassword), nil
}
