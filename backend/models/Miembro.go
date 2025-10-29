package models

import (
	"arka-code/backend/config"
	"database/sql"
)

type Miembro struct {
	ID       int
	FamilyID int
	Name     string
	Email    string
	Type     string
	Status   string
}

type MiembroModel struct {
	db *sql.DB
}

func NewMiembroModel() *MiembroModel {
	return &MiembroModel{db: config.DB}
}

func (m *MiembroModel) GetByFamilyID(familyID int) ([]Miembro, error) {
	rows, err := m.db.Query(`
		SELECT id, family_id, name, email, type, status
		FROM users 
		WHERE family_id = $1 AND status = 'active' AND type != 'family'`,
		familyID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var miembros []Miembro
	for rows.Next() {
		var miembro Miembro
		err := rows.Scan(&miembro.ID, &miembro.FamilyID, &miembro.Name, &miembro.Email, &miembro.Type, &miembro.Status)
		if err != nil {
			return nil, err
		}
		miembros = append(miembros, miembro)
	}

	return miembros, nil
}

func (m *MiembroModel) GetByID(id int) (*Miembro, error) {
	var miembro Miembro
	err := m.db.QueryRow(`
		SELECT id, family_id, name, email, type, status
		FROM users WHERE id = $1 AND status = 'active'`,
		id).Scan(&miembro.ID, &miembro.FamilyID, &miembro.Name, &miembro.Email, &miembro.Type, &miembro.Status)

	if err != nil {
		return nil, err
	}

	return &miembro, nil
}
