package models

import (
	"arka-code/database"
	"arka-code/entities"
	"arka-code/utils"
	"database/sql"
	"log"
)

type UserModel struct{}

var UserModelInstance = &UserModel{}

func (m *UserModel) Create(user *entities.Usuario) error {
    query := `CALL sp_create_usuario(?, ?, ?, ?, ?)`

    var lastInsertID int64
    err := database.DB.QueryRow(query,
        user.NombreUsuario,
        user.Rol,
        user.ContrasenaPersonal,
        user.NombrePersonal,
        user.CorreoFamilia).Scan(&lastInsertID)

    if err != nil {
        return err
    }

    return nil
}

func (m *UserModel) FindByEmail(email string) (*entities.Usuario, error) {
    query := `CALL sp_find_usuario_by_email(?)`

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

    if err != nil {
        return nil, err
    }

    return user, nil
}

func (m *UserModel) FindByNombreUsuario(nombreUsuario string) (*entities.Usuario, error) {
    query := `CALL sp_find_usuario_by_nombre(?)`

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

    if err != nil {
        return nil, err
    }

    return user, nil
}

func (m *UserModel) FindByFamilia(correoFamilia string) ([]entities.Usuario, error) {
    query := `CALL sp_find_usuarios_by_familia(?)`

    rows, err := database.DB.Query(query, correoFamilia)
    if err != nil {
        log.Printf("❌ Error en consulta FindByFamilia: %v", err)
        return nil, err
    }
    defer rows.Close()

    usuarios := []entities.Usuario{}
    for rows.Next() {
        var usuario entities.Usuario

        err := rows.Scan(
            &usuario.NombreUsuario,
            &usuario.Rol,
            &usuario.ContrasenaPersonal,
            &usuario.NombrePersonal,
            &usuario.CorreoFamilia,
            &usuario.DeleteAt,
        )
        if err != nil {
            log.Printf("❌ Error escaneando usuario: %v", err)
            continue
        }
        usuarios = append(usuarios, usuario)
    }

    log.Printf("👥 Usuarios encontrados: %d para familia %s", len(usuarios), correoFamilia)
    return usuarios, nil
}

func (m *UserModel) GetAllByFamilia(correoFamilia string) ([]entities.Usuario, error) {
    query := `CALL sp_get_all_usuarios_by_familia(?)`

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
    query := `CALL sp_update_usuario_rol(?, ?)`
    
    var rowsAffected int64
    err := database.DB.QueryRow(query, nombreUsuario, rol).Scan(&rowsAffected)
    if err != nil {
        log.Printf("❌ Error actualizando rol: %v", err)
        return err
    }

    log.Printf("✅ Rol actualizado - Usuario: %s, Filas afectadas: %d", nombreUsuario, rowsAffected)
    return nil
}

func (m *UserModel) UpdateNombrePersonal(nombreUsuario string, nuevoNombre string) error {
    query := `CALL sp_update_usuario_nombre_personal(?, ?)`
    
    var rowsAffected int64
    err := database.DB.QueryRow(query, nombreUsuario, nuevoNombre).Scan(&rowsAffected)
    if err != nil {
        log.Printf("❌ Error actualizando nombre: %v", err)
        return err
    }

    log.Printf("✅ Nombre actualizado - Usuario: %s, Filas afectadas: %d", nombreUsuario, rowsAffected)
    return nil
}

func (m *UserModel) UpdateNombreUsuario(viejoUsuario string, nuevoUsuario string) error {
    query := `CALL sp_update_usuario_nombre_usuario(?, ?)`
    
    var rowsAffected int64
    err := database.DB.QueryRow(query, viejoUsuario, nuevoUsuario).Scan(&rowsAffected)
    if err != nil {
        log.Printf("❌ Error actualizando nombre de usuario: %v", err)
        return err
    }

    log.Printf("✅ Nombre de usuario actualizado - De: %s a %s, Filas afectadas: %d", viejoUsuario, nuevoUsuario, rowsAffected)
    return nil
}

func (m *UserModel) UpdatePassword(nombreUsuario string, nuevaPassword string) error {
    hashedPassword, err := utils.HashPassword(nuevaPassword)
    if err != nil {
        log.Printf("❌ Error hasheando contraseña: %v", err)
        return err
    }

    query := `CALL sp_update_usuario_password(?, ?)`
    
    var rowsAffected int64
    err = database.DB.QueryRow(query, nombreUsuario, hashedPassword).Scan(&rowsAffected)
    if err != nil {
        log.Printf("❌ Error actualizando contraseña: %v", err)
        return err
    }

    log.Printf("✅ Contraseña actualizada - Usuario: %s, Filas afectadas: %d", nombreUsuario, rowsAffected)
    return nil
}

// Update actualiza múltiples campos del usuario
func (m *UserModel) Update(usuario *entities.Usuario) error {
    query := `CALL sp_update_usuario(?, ?, ?)`

    var rowsAffected int64
    err := database.DB.QueryRow(query,
        usuario.NombreUsuario,
        usuario.NombrePersonal,
        usuario.Rol).Scan(&rowsAffected)

    if err != nil {
        log.Printf("❌ Error actualizando usuario: %v", err)
        return err
    }

    log.Printf("✅ Usuario actualizado - Usuario: %s, Filas afectadas: %d", usuario.NombreUsuario, rowsAffected)
    return nil
}

func (m *UserModel) SoftDelete(nombreUsuario string) error {
    query := `CALL sp_soft_delete_usuario(?)`
    
    var rowsAffected int64
    err := database.DB.QueryRow(query, nombreUsuario).Scan(&rowsAffected)
    if err != nil {
        log.Printf("❌ Error deshabilitando usuario: %v", err)
        return err
    }

    if rowsAffected == 0 {
        log.Printf("⚠️ Usuario %s no encontrado o ya estaba deshabilitado", nombreUsuario)
    }

    log.Printf("✅ Usuario deshabilitado - Usuario: %s, Filas afectadas: %d", nombreUsuario, rowsAffected)
    return nil
}

func (m *UserModel) CheckPasswordUser(nombreUsuario, contrasena string) (bool, error) {
    query := `CALL sp_get_usuario_password(?)`

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
    query := `CALL sp_usuario_exists(?)`

    var count int
    err := database.DB.QueryRow(query, nombreUsuario).Scan(&count)
    if err != nil {
        return false, err
    }

    return count > 0, nil
}
func (m *UserModel) Delete(nombreUsuario string) error {
    query := `CALL sp_delete_usuario_permanently(?)`
    
    var rowsAffected int64
    err := database.DB.QueryRow(query, nombreUsuario).Scan(&rowsAffected)
    if err != nil {
        log.Printf("❌ Error eliminando usuario: %v", err)
        return err
    }

    if rowsAffected == 0 {
        log.Printf("⚠️ Usuario %s no encontrado", nombreUsuario)
    }

    log.Printf("✅ Usuario eliminado permanentemente: %s", nombreUsuario)
    return nil
}
