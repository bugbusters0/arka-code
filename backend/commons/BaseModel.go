package commons

import (
	"arka-code/backend/config"
	"database/sql"
	"time"
)

type BaseModel struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type IModel interface {
	GetTableName() string
	FindByID(id int) error
	Save() error
	Delete() error
}

func (bm *BaseModel) GetDB() *sql.DB {
	return config.DB
}

func (bm *BaseModel) Timestamps() {
	if bm.CreatedAt.IsZero() {
		bm.CreatedAt = time.Now()
	}
	bm.UpdatedAt = time.Now()
}
