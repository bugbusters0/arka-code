package entities

import "time"

type Familia struct {
	Correo     string     `json:"correo"`
	Telefono   *string    `json:"telefono"`
	Contraseña string     `json:"-"`
	DeleteAt   *time.Time `json:"delete_at"`
}
