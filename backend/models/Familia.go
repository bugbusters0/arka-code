package models

import (
	"arka-code/backend/config"
	"database/sql"
	"time"
)

type Familia struct {
	ID        int
	Name      string
	CreatedAt time.Time
}

type FamiliaModel struct {
	db *sql.DB
}

func NewFamiliaModel() *FamiliaModel {
	return &FamiliaModel{db: config.DB}
}

func (m *FamiliaModel) Create(nombre string) (int, error) {
	var id int
	err := m.db.QueryRow(`
		INSERT INTO families (name, created_at) 
		VALUES ($1, NOW()) RETURNING id`,
		nombre).Scan(&id)

	return id, err
}

func (m *FamiliaModel) GetByID(id int) (*Familia, error) {
	var familia Familia
	err := m.db.QueryRow(`
		SELECT id, name, created_at
		FROM families WHERE id = $1`,
		id).Scan(&familia.ID, &familia.Name, &familia.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &familia, nil
}
