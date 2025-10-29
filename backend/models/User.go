package models

import (
	"arka-code/backend/commons"
)

type User struct {
	commons.BaseModel
	FamilyID int    `json:"family_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Type     string `json:"type"`
	Status   string `json:"status"`
}

type UserModel struct {
	commons.BaseModel
}

func NewUserModel() *UserModel {
	return &UserModel{}
}

func (m *UserModel) GetTableName() string {
	return "users"
}

func (m *UserModel) FindByID(id int) (*User, error) {
	var user User
	err := m.GetDB().QueryRow(`
		SELECT id, family_id, name, email, password, type, status, created_at, updated_at
		FROM users WHERE id = $1`,
		id).Scan(&user.ID, &user.FamilyID, &user.Name, &user.Email, &user.Password, &user.Type, &user.Status, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &user, nil
}
