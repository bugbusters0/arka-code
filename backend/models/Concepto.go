package models

import (
	"arka-code/backend/config"
	"database/sql"
	"time"
)

type Concepto struct {
	ID        int
	MemberID  int
	Name      string
	Type      string
	Status    string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

type ConceptoModel struct {
	db *sql.DB
}

func NewConceptoModel() *ConceptoModel {
	return &ConceptoModel{db: config.DB}
}

func (m *ConceptoModel) GetByMiembroAndTipo(miembroID int, tipo string) ([]Concepto, error) {
	rows, err := m.db.Query(`
		SELECT id, member_id, name, type, status, created_at, updated_at
		FROM concepts 
		WHERE member_id = $1 AND type = $2 AND status = 'active'
		ORDER BY created_at DESC`,
		miembroID, tipo)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conceptos []Concepto
	for rows.Next() {
		var c Concepto
		err := rows.Scan(&c.ID, &c.MemberID, &c.Name, &c.Type, &c.Status, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		conceptos = append(conceptos, c)
	}

	return conceptos, nil
}

func (m *ConceptoModel) Create(miembroID int, nombre, tipo string) error {
	_, err := m.db.Exec(`
		INSERT INTO concepts (member_id, name, type, status, created_at) 
		VALUES ($1, $2, $3, 'active', NOW())`,
		miembroID, nombre, tipo)

	return err
}

func (m *ConceptoModel) GetByID(id int) (*Concepto, error) {
	var c Concepto
	err := m.db.QueryRow(`
		SELECT id, member_id, name, type, status, created_at, updated_at
		FROM concepts WHERE id = $1 AND status = 'active'`,
		id).Scan(&c.ID, &c.MemberID, &c.Name, &c.Type, &c.Status, &c.CreatedAt, &c.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (m *ConceptoModel) Update(id int, nombre, tipo string) error {
	_, err := m.db.Exec(`
		UPDATE concepts SET name = $1, type = $2, updated_at = NOW() 
		WHERE id = $3`,
		nombre, tipo, id)

	return err
}

func (m *ConceptoModel) Delete(id int) error {
	_, err := m.db.Exec(`
		UPDATE concepts SET status = 'inactive', updated_at = NOW() 
		WHERE id = $1`,
		id)

	return err
}

func (m *ConceptoModel) Disable(id int) error {
	_, err := m.db.Exec(`
		UPDATE concepts SET status = 'disabled', updated_at = NOW() 
		WHERE id = $1`,
		id)

	return err
}
