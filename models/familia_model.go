package models

import (
	"arka-code/database"
	"arka-code/entities"
	"arka-code/utils"
	"database/sql"
	"log"
)

type FamiliaModel struct{}

var FamiliaModelInstance = &FamiliaModel{}

// Create crea una nueva familia en el sistema
// Parámetros:
// - familia: Estructura con los datos de la familia a crear
// Retorno: Error si la creación falla
// Flujo:
// - Inserta nuevo registro en tabla familia
// - Campos requeridos: correo, contraseña
// - Campo opcional: telefono
// - No maneja transacciones (creación simple)
// Uso: Registro de nuevas familias en el sistema
func (m *FamiliaModel) Create(familia *entities.Familia) error {
	query := `INSERT INTO familia (correo, telefono, contraseña)
	          VALUES (?, ?, ?)`

	_, err := database.DB.Exec(query,
		familia.Correo,
		familia.Telefono,
		familia.Contraseña)

	return err
}

// FindByEmail busca una familia por su correo electrónico
// Parámetros:
// - email: Correo electrónico de la familia a buscar
// Retorno: Puntero a la familia encontrada o nil si no existe
// Flujo:
// - Consulta familia por correo exacto
// - Excluye familias eliminadas (soft delete)
// - Retorna sql.ErrNoRows si no se encuentra
// Uso: Verificación de existencia, login de familias
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

// GetAll obtiene todas las familias activas del sistema
// Parámetros: Ninguno
// Retorno: Lista de todas las familias no eliminadas
// Flujo:
// - Consulta todas las familias con delete_at IS NULL
// - NOTA: La consulta actual tiene un error (WHERE correo = ? sin parámetro)
// - Debería ser: "SELECT * FROM familia WHERE delete_at IS NULL"
// Uso: Administración del sistema, reportes globales
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

// CheckPasswordFamilia verifica si la contraseña proporcionada coincide con la almacenada
// Parámetros:
// - correo: Correo de la familia a verificar
// - contrasena: Contraseña en texto plano a verificar
// Retorno: true si la contraseña es correcta, false si no coincide o no existe
// Flujo:
// 1. Obtiene la contraseña hash almacenada para la familia
// 2. Si la familia no existe, retorna false
// 3. Usa utils.CheckPassword para comparar contraseñas
// 4. Retorna resultado de la verificación
// Uso: Autenticación durante el login familiar
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

// Delete elimina permanentemente una familia del sistema
// Parámetros:
// - correo: Correo de la familia a eliminar
// Retorno: Error si la eliminación falla
// Flujo:
// - Eliminación física (DELETE) no soft delete
// - Registra éxito/error en logs
// - ADVERTENCIA: Esta eliminación es permanente y puede romper integridad referencial
// Uso: Eliminación administrativa de familias (uso con precaución)
func (m *FamiliaModel) Delete(correo string) error {
	query := `DELETE FROM familia WHERE correo = ?`
	_, err := database.DB.Exec(query, correo)
	if err != nil {
		log.Printf("❌ Error eliminando familia: %v", err)
	} else {
		log.Printf("✅ Familia eliminada: %s", correo)
	}
	return err
}
