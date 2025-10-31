# Arka Code - Sistema de Gestión Financiera Familiar

## 📋 Descripción

Arka Code es una aplicación web desarrollada en **Go (Golang)** para la gestión financiera familiar. Permite administrar conceptos de gastos e ingresos, registrar movimientos diarios, establecer límites de gasto por usuario y gestionar perfiles familiares con roles diferenciados.

## 🚀 Características

- **Autenticación multinivel** con perfiles familiares y usuarios individuales
- **Gestión de conceptos** personalizables para gastos e ingresos
- **Registro de movimientos** diarios con categorización automática
- **Sistema de límites de gasto** por usuario y concepto
- **Roles diferenciados** (Administrador y Miembro)
- **Interfaz responsive** con diseño moderno
- **Arquitectura MVC** limpia y escalable
- **Base de datos MySQL/MariaDB** con soft deletes
- **Validación robusta** de formularios server-side

## 🛠️ Tecnologías

- **Backend**: Go 1.21+
- **Base de datos**: MySQL 8.0 / MariaDB 10.6+
- **Frontend**: HTML5, CSS3, JavaScript Vanilla
- **Template Engine**: Go HTML Templates
- **Sesiones**: Gorilla Sessions
- **Router**: http.ServeMux (Go estándar)
- **Deploy**: Compatible con cualquier servidor que soporte Go

---

## 📁 Estructura del Proyecto

```
arka-code/
├── config/
│   └── config.go              # Configuración de la aplicación
├── controllers/
│   ├── auth_controller.go     # Autenticación y registro
│   ├── perfil_controller.go   # Selección de perfil
│   ├── concepto_controller.go # Gestión de conceptos
│   ├── movimiento_controller.go # Entrada diaria de movimientos
│   └── perfil_controller.go   # Gestión de usuarios/perfiles
├── database/
│   └── connection.go          # Conexión a MySQL
├── entities/
│   ├── familia.go             # Entidad Familia
│   ├── usuario.go             # Entidad Usuario
│   ├── concepto.go            # Entidad Concepto
│   ├── movimiento.go          # Entidad Movimiento
│   ├── personalizacion_concepto.go # Entidad Límites
│   └── session_data.go        # Datos de sesión
├── middleware/
│   ├── auth.go                # Middleware de autenticación
│   └── role.go                # Middleware de roles
├── models/
│   ├── familia_model.go       # CRUD de familias
│   ├── user_model.go          # CRUD de usuarios
│   ├── concepto_model.go      # CRUD de conceptos
│   ├── movimiento_model.go    # CRUD de movimientos
│   └── personalizacion_model.go # CRUD de límites
├── routes/
│   └── routes.go              # Definición de rutas
├── static/
│   ├── css/                   # Estilos CSS
│   ├── js/                    # Scripts JavaScript
│   └── img/                   # Imágenes y recursos
├── utils/
│   ├── hash.go                # Hash de contraseñas (bcrypt)
│   ├── render.go              # Renderizado de templates
│   └── session.go             # Manejo de sesiones
├── validators/
│   ├── validator_base.go      # Validador base
│   ├── register_validator.go  # Validación de registro
│   ├── login_validator.go     # Validación de login
│   ├── concepto_validator.go  # Validación de conceptos
│   ├── movimiento_validator.go # Validación de movimientos
│   └── perfil_validator.go    # Validación de perfiles
├── views/
│   ├── layouts/
│   │   ├── auth.html          # Layout para auth
│   │   └── dashboard.html     # Layout principal
│   ├── auth/
│   │   ├── login.html         # Vista de login
│   │   └── registro.html      # Vista de registro
│   ├── perfiles/
│   │   ├── seleccionar.html   # Selección de perfil
│   │   └── index.html         # Gestión de perfiles
│   ├── conceptos/
│   │   └── index.html         # Gestión de conceptos
│   └── movimientos/
│       └── index.html         # Entrada diaria
├── main.go                    # Punto de entrada
├── go.mod                     # Dependencias Go
└── go.sum                     # Checksums de dependencias
```

---

## 🏃‍♂️ Desarrollo Local

### Prerrequisitos

- **Go 1.21 o superior**
- **MySQL 8.0 / MariaDB 10.6+**
- **Git**

### Instalación

1. **Clonar el repositorio**

```bash
git clone https://github.com/bugbusters0/arka-code.git
cd arka-code
```

2. **Instalar dependencias de Go**

```bash
go mod download
```

3. **Configurar base de datos**

```bash
# Conectar a MySQL
mysql -u root -p

# Crear base de datos y usuario
CREATE DATABASE arka;
CREATE USER 'arka_user'@'localhost' IDENTIFIED BY 'arka_password';
GRANT ALL PRIVILEGES ON arka.* TO 'arka_user'@'localhost';
FLUSH PRIVILEGES;
EXIT;

# Importar estructura
mysql -u root -p arka < sql.sql
```

4. **Configurar variables de entorno**

Crear archivo `.env` en la raíz del proyecto:

```env
APP_NAME=Arka
APP_PORT=8080
DB_HOST=localhost
DB_PORT=3306
DB_NAME=arka
DB_USER=arka_user
DB_PASSWORD=arka_password
SESSION_KEY=tu-clave-secreta-aqui-cambiar-en-produccion
```

5. **Ejecutar la aplicación**

```bash
# Modo desarrollo
go run main.go

# O compilar y ejecutar
go build -o arka-code
./arka-code
```

6. **Acceder a la aplicación**

```
http://localhost:8080
```

---

## 🔧 Configuración

### Variables de Entorno - `config/config.go`

```go
type Config struct {
    AppName     string
    Port        string
    DBHost      string
    DBPort      string
    DBName      string
    DBUser      string
    DBPassword  string
    SessionKey  string
}

func LoadConfig() {
    // Carga variables desde .env o valores por defecto
    AppConfig = Config{
        AppName:    getEnv("APP_NAME", "Arka"),
        Port:       getEnv("APP_PORT", "8080"),
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "3306"),
        DBName:     getEnv("DB_NAME", "arka"),
        DBUser:     getEnv("DB_USER", "arka_user"),
        DBPassword: getEnv("DB_PASSWORD", "arka_password"),
        SessionKey: getEnv("SESSION_KEY", "default-key-change-in-prod"),
    }
}
```

---

## 🔄 Flujo de la Aplicación

### 1. Punto de Entrada - `main.go`

```go
package main

import (
    "arka-code/config"
    "arka-code/database"
    "arka-code/routes"
    "log"
    "net/http"
)

func main() {
    // 1. Cargar configuración
    config.LoadConfig()

    // 2. Conectar base de datos
    database.Connect()
    defer database.Close()

    // 3. Configurar rutas
    router := routes.SetupRoutes()

    // 4. Iniciar servidor
    log.Printf("🚀 Servidor iniciado en http://localhost:%s", config.AppConfig.Port)
    log.Fatal(http.ListenAndServe(":"+config.AppConfig.Port, router))
}
```

### 2. Sistema de Rutas - `routes/routes.go`

```go
func SetupRoutes() *http.ServeMux {
    utils.InitSession()
    mux := http.NewServeMux()

    // Archivos estáticos
    fs := http.FileServer(http.Dir("./static"))
    mux.Handle("/static/", http.StripPrefix("/static/", fs))

    // ========================================
    // RUTAS PÚBLICAS
    // ========================================
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if utils.IsAuthenticated(r) {
            http.Redirect(w, r, "/movimientos", http.StatusSeeOther)
        } else {
            http.Redirect(w, r, "/login", http.StatusSeeOther)
        }
    })

    // ========================================
    // AUTENTICACIÓN DE FAMILIA
    // ========================================
    mux.HandleFunc("/registro", middleware.RequireGuest(
        controllers.AuthControllerInstance.ShowRegister,
    ))

    mux.HandleFunc("/registro/submit", middleware.RequireGuest(
        controllers.AuthControllerInstance.Register,
    ))

    mux.HandleFunc("/login", middleware.RequireGuest(
        controllers.AuthControllerInstance.ShowLogin,
    ))

    mux.HandleFunc("/login/submit", middleware.RequireGuest(
        controllers.AuthControllerInstance.Login,
    ))

    // ========================================
    // SELECCIÓN DE PERFIL
    // ========================================
    mux.HandleFunc("/seleccionar-perfil", middleware.RequireFamiliaAuth(
        controllers.PerfilControllerInstance.ShowSeleccionarPerfil,
    ))

    mux.HandleFunc("/seleccionar-perfil/submit", middleware.RequireFamiliaAuth(
        controllers.PerfilControllerInstance.SeleccionarPerfil,
    ))

    // ========================================
    // ÁREA PROTEGIDA (REQUIERE USUARIO)
    // ========================================

    // Conceptos
    mux.HandleFunc("/conceptos", middleware.RequireUsuarioAuth(
        controllers.ConceptoControllerInstance.Index,
    ))

    mux.HandleFunc("/conceptos/crear", middleware.RequireUsuarioAuth(
        controllers.ConceptoControllerInstance.Crear,
    ))

    // Movimientos (Entrada Diaria)
    mux.HandleFunc("/movimientos", middleware.RequireUsuarioAuth(
        controllers.MovimientoControllerInstance.Index,
    ))

    mux.HandleFunc("/movimientos/crear", middleware.RequireUsuarioAuth(
        controllers.MovimientoControllerInstance.Crear,
    ))

    mux.HandleFunc("/movimientos/editar", middleware.RequireUsuarioAuth(
        controllers.MovimientoControllerInstance.Editar,
    ))

    mux.HandleFunc("/movimientos/eliminar", middleware.RequireUsuarioAuth(
        controllers.MovimientoControllerInstance.Eliminar,
    ))

    // Perfiles (Solo Admin para crear/deshabilitar)
    mux.HandleFunc("/perfiles", middleware.RequireUsuarioAuth(
        controllers.PerfilControllerInstance.Index,
    ))

    mux.HandleFunc("/perfiles/crear", middleware.RequireUsuarioAuth(
        middleware.RequireAdmin(
            controllers.PerfilControllerInstance.Crear,
        ),
    ))

    mux.HandleFunc("/perfiles/editar", middleware.RequireUsuarioAuth(
        controllers.PerfilControllerInstance.Editar,
    ))

    mux.HandleFunc("/perfiles/deshabilitar", middleware.RequireUsuarioAuth(
        middleware.RequireAdmin(
            controllers.PerfilControllerInstance.Deshabilitar,
        ),
    ))

    // Límites de Gasto
    mux.HandleFunc("/perfiles/limites/crear", middleware.RequireUsuarioAuth(
        controllers.PerfilControllerInstance.CrearLimite,
    ))

    mux.HandleFunc("/perfiles/limites/editar", middleware.RequireUsuarioAuth(
        controllers.PerfilControllerInstance.EditarLimite,
    ))

    mux.HandleFunc("/perfiles/limites/eliminar", middleware.RequireUsuarioAuth(
        controllers.PerfilControllerInstance.EliminarLimite,
    ))

    // Logout
    mux.HandleFunc("/logout-familia", middleware.RequireFamiliaAuth(
        controllers.AuthControllerInstance.LogoutFamilia,
    ))

    mux.HandleFunc("/logout-usuario", middleware.RequireUsuarioAuth(
        controllers.AuthControllerInstance.LogoutUsuario,
    ))

    return mux
}
```

### 3. Middleware de Autenticación - `middleware/auth.go`

```go
// RequireGuest - Solo permite acceso a usuarios no autenticados
func RequireGuest(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if utils.IsAuthenticated(r) {
            http.Redirect(w, r, "/seleccionar-perfil", http.StatusSeeOther)
            return
        }
        next(w, r)
    }
}

// RequireFamiliaAuth - Requiere autenticación de familia
func RequireFamiliaAuth(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if !utils.IsFamiliaAuthenticated(r) {
            http.Redirect(w, r, "/login", http.StatusSeeOther)
            return
        }
        next(w, r)
    }
}

// RequireUsuarioAuth - Requiere usuario seleccionado
func RequireUsuarioAuth(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if !utils.IsUsuarioAuthenticated(r) {
            http.Redirect(w, r, "/seleccionar-perfil", http.StatusSeeOther)
            return
        }
        next(w, r)
    }
}
```

### 4. Controlador de Ejemplo - `controllers/movimiento_controller.go`

```go
type MovimientoController struct{}

func (c *MovimientoController) Index(w http.ResponseWriter, r *http.Request) {
    // 1. Obtener datos de sesión
    sessionData, ok := utils.GetSessionData(r)
    if !ok {
        http.Redirect(w, r, "/seleccionar-perfil", http.StatusSeeOther)
        return
    }

    // 2. Obtener parámetros de la URL
    tipo := r.URL.Query().Get("tipo")
    if tipo == "" {
        tipo = "gasto" // Por defecto
    }

    fechaStr := r.URL.Query().Get("fecha")
    fecha, _ := time.Parse("2006-01-02", fechaStr)
    if fechaStr == "" {
        fecha = time.Now()
    }

    // 3. Obtener conceptos y movimientos desde el modelo
    conceptos, _ := models.ConceptoModelInstance.FindByFamilia(
        sessionData.CorreoFamilia,
        tipo,
    )

    movimientos, _ := models.MovimientoModelInstance.FindByFamiliaDateAndTipo(
        sessionData.CorreoFamilia,
        fecha,
        int8(0), // 0 = gasto
    )

    // 4. Preparar datos para la vista
    data := map[string]interface{}{
        "Title":       "Entrada Diaria",
        "CurrentPage": "dashboard",
        "SessionData": sessionData,
        "Tipo":        tipo,
        "Conceptos":   conceptos,
        "Movimientos": movimientos,
        "FechaActual": fecha.Format("2006-01-02"),
    }

    // 5. Renderizar template
    utils.RenderTemplate(w, "dashboard", "movimientos/index", data)
}
```

### 5. Modelo de Ejemplo - `models/movimiento_model.go`

```go
type MovimientoModel struct{}

func (m *MovimientoModel) Create(movimiento *entities.Movimiento) error {
    query := `INSERT INTO movimiento (fecha, monto, descripcion, nombreUsuario,
                                     nombreConcepto, correoFamilia)
              VALUES (?, ?, ?, ?, ?, ?)`

    result, err := database.DB.Exec(query,
        movimiento.Fecha,
        movimiento.Monto,
        movimiento.Descripcion,
        movimiento.NombreUsuario,
        movimiento.NombreConcepto,
        movimiento.CorreoFamilia)

    if err != nil {
        log.Printf("❌ Error creando movimiento: %v", err)
        return err
    }

    id, _ := result.LastInsertId()
    movimiento.IdMovimiento = int(id)

    log.Printf("✅ Movimiento creado - ID: %d", id)
    return nil
}

func (m *MovimientoModel) FindByFamiliaAndDate(correoFamilia string, fecha time.Time) ([]entities.Movimiento, error) {
    query := `SELECT m.idMovimiento, m.fecha, m.monto, m.descripcion,
                     m.nombreUsuario, m.nombreConcepto, m.correoFamilia
              FROM movimiento m
              WHERE m.correoFamilia = ?
                AND DATE(m.fecha) = DATE(?)
                AND m.delete_at IS NULL
              ORDER BY m.fecha DESC`

    rows, err := database.DB.Query(query, correoFamilia, fecha)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    movimientos := []entities.Movimiento{}
    for rows.Next() {
        var m entities.Movimiento
        rows.Scan(&m.IdMovimiento, &m.Fecha, &m.Monto, &m.Descripcion,
                 &m.NombreUsuario, &m.NombreConcepto, &m.CorreoFamilia)
        movimientos = append(movimientos, m)
    }

    return movimientos, nil
}
```

### 6. Validador - `validators/movimiento_validator.go`

```go
type MovimientoValidator struct {
    *ValidatorBase
}

func (v *MovimientoValidator) Validate(r *http.Request) ValidationResult {
    result := NewValidationResult()

    // Validar monto
    montoStr := strings.TrimSpace(r.FormValue("monto"))
    if montoStr == "" {
        result.Errors["monto"] = "El monto es obligatorio"
        result.Success = false
    } else {
        monto, err := strconv.ParseFloat(montoStr, 64)
        if err != nil || monto <= 0 {
            result.Errors["monto"] = "El monto debe ser mayor a 0"
            result.Success = false
        } else {
            result.CleanData["monto"] = monto
        }
    }

    // Validar concepto
    concepto := strings.TrimSpace(r.FormValue("nombreConcepto"))
    if concepto == "" {
        result.Errors["nombreConcepto"] = "Debe seleccionar un concepto"
        result.Success = false
    } else {
        result.CleanData["nombreConcepto"] = concepto
    }

    // Validar fecha
    fechaStr := strings.TrimSpace(r.FormValue("fecha"))
    if fechaStr == "" {
        result.Errors["fecha"] = "La fecha es obligatoria"
        result.Success = false
    } else {
        fecha, err := time.Parse("2006-01-02", fechaStr)
        if err != nil {
            result.Errors["fecha"] = "Formato de fecha inválido"
            result.Success = false
        } else {
            result.CleanData["fecha"] = fecha
        }
    }

    return result
}
```

### 7. Template - `views/movimientos/index.html`

```html
{{define "content"}}
<div class="container">
  <div class="topbar">
    <div class="fecha-display" id="fechaDisplay">
      <i id="fechaTexto">{{.FechaActual}}</i>
      <input type="date" id="fechaInput" value="{{.FechaActual}}" style="display:none;">
    </div>

    <div class="tabs">
      <a href="/movimientos?tipo=gasto&fecha={{.FechaActual}}"
         class="tab {{if eq .Tipo "gasto"}}active{{end}}">Gastos</a>
      <a href="/movimientos?tipo=ingreso&fecha={{.FechaActual}}"
         class="tab {{if eq .Tipo "ingreso"}}active{{end}}">Ingresos</a>
    </div>
  </div>

  <div class="main">
    <form method="POST" action="/movimientos/crear">
      <input type="hidden" name="tipo" value="{{.Tipo}}">
      <input type="hidden" name="fecha" value="{{.FechaActual}}">

      <div class="form-group">
        <label>Concepto *</label>
        <div class="concepts-strip">
          {{range .Conceptos}}
          <div class="concept-item" onclick="selectConcept('{{.NombreConcepto}}')">
            <div class="concept-circle" style="background-color: {{.Color}};">
              <i class="{{.Icono}}"></i>
            </div>
            <div class="concept-label">{{.NombreConcepto}}</div>
          </div>
          {{end}}
        </div>
        <input type="hidden" name="nombreConcepto" id="nombreConcepto">
      </div>

      <div class="form-group">
        <label>Monto (S/) *</label>
        <input type="number" step="0.01" name="monto" required>
      </div>

      <button type="submit" class="btn btn-primary">Guardar</button>
    </form>

    <div class="movements">
      {{range .Movimientos}}
      <div class="movement-card">
        <div class="movement-title">{{.NombreConcepto}}</div>
        <div class="movement-amount">S/ {{printf "%.2f" .Monto}}</div>
      </div>
      {{end}}
    </div>
  </div>
</div>
{{end}}
```

---

## 🗃️ Estructura de Base de Datos

### Tablas Principales:

```sql
-- Familia (grupo familiar)
CREATE TABLE `familia` (
  `correo` VARCHAR(50) PRIMARY KEY,
  `telefono` CHAR(9),
  `contraseña` VARCHAR(255) NOT NULL,
  `delete_at` TIMESTAMP NULL
);

-- Usuarios (miembros de la familia)
CREATE TABLE `usuario` (
  `nombreUsuario` VARCHAR(30) PRIMARY KEY,
  `rol` TINYINT(1) NOT NULL DEFAULT 0,  -- 0=miembro, 1=admin
  `contraseñaPersonal` VARCHAR(255) NOT NULL,
  `nombrePersonal` VARCHAR(100) NOT NULL,
  `correoFamilia` VARCHAR(50) NOT NULL,
  `delete_at` TIMESTAMP NULL,
  FOREIGN KEY (`correoFamilia`) REFERENCES `familia`(`correo`)
);

-- Conceptos (categorías de gastos/ingresos)
CREATE TABLE `concepto` (
  `nombreConcepto` VARCHAR(40),
  `correoFamilia` VARCHAR(50),
  `tipo` TINYINT(1) NOT NULL,  -- 0=gasto, 1=ingreso
  `icono` VARCHAR(100),
  `color` CHAR(7),
  `nombreUsuario` VARCHAR(30) NOT NULL,
  `delete_at` TIMESTAMP NULL,
  PRIMARY KEY (`nombreConcepto`, `correoFamilia`)
);

-- Movimientos (registros financieros)
CREATE TABLE `movimiento` (
  `idMovimiento` INT AUTO_INCREMENT PRIMARY KEY,
  `fecha` DATE NOT NULL,
  `monto` DECIMAL(10,2) NOT NULL,
  `descripcion` TEXT,
  `nombreUsuario` VARCHAR(30) NOT NULL,
  `nombreConcepto` VARCHAR(40) NOT NULL,
  `correoFamilia` VARCHAR(50) NOT NULL,
  `delete_at` TIMESTAMP NULL
);

-- Límites de gasto personalizados
CREATE TABLE `personalizacionconcepto` (
  `idPersonalizacion` INT AUTO_INCREMENT PRIMARY KEY,
  `limiteGasto` DECIMAL(10,2),
  `activo` TINYINT(1) DEFAULT 1,
  `montoPlanificado` DECIMAL(10,2),
  `tipoPeriodoPlanificado` VARCHAR(15),
  `tipoPeriodoLimite` VARCHAR(15),
  `nombreUsuario` VARCHAR(30) NOT NULL,
  `nombreConcepto` VARCHAR(40) NOT NULL,
  `correoFamilia` VARCHAR(50) NOT NULL,
  `delete_at` TIMESTAMP NULL
);
```

---

## 🔐 Sistema de Autenticación

### Niveles de Acceso:

1. **Guest** (No autenticado)

   - Puede acceder a `/login` y `/registro`
   - Middleware: `RequireGuest`

2. **Familia Autenticada** (Familia logueada)

   - Puede acceder a `/seleccionar-perfil`
   - Middleware: `RequireFamiliaAuth`

3. **Usuario Autenticado** (Perfil seleccionado)

   - Puede acceder a todas las rutas protegidas
   - Middleware: `RequireUsuarioAuth`

4. **Administrador** (Usuario con rol=1)
   - Puede gestionar todos los perfiles
   - Middleware: `RequireAdmin`

### Flujo de Autenticación:

```
1. Usuario → /login
2. Validar credenciales → Crear sesión familia
3. Redirigir → /seleccionar-perfil
4. Seleccionar perfil → Crear sesión usuario
5. Redirigir → /movimientos (dashboard)
```

---

## 🎯 Funcionalidades por Módulo

### 📝 Autenticación (`auth_controller.go`)

- Registro de familias con múltiples usuarios
- Login con email y contraseña (bcrypt)
- Selección de perfil individual
- Logout de familia y usuario

### 👥 Perfiles (`perfil_controller.go`)

- Crear nuevo miembro (admin)
- Editar perfil propio o de otros (admin)
- Deshabilitar miembro (admin, soft delete)
- Gestión de límites de gasto por usuario

### 📊 Conceptos (`concepto_controller.go`)

- Crear conceptos de gasto/ingreso
- Asignar íconos Font Awesome y colores
- Configuración personal por usuario
- Activar/desactivar conceptos

### 💰 Movimientos (`movimiento_controller.go`)

- Registro diario de gastos/ingresos
- Vista por tipo (Gastos/Ingresos/Resumen)
- Filtrado por fecha
- Editar/eliminar movimientos
- Cálculo automático de totales

---

## 🔒 Seguridad

- **Contraseñas**: Hash bcrypt (factor 10)
- **SQL Injection**: Prepared statements con `database/sql`
- **XSS**: Auto-escape en Go templates
- **CSRF**: Gorilla Sessions con cookies seguras
- **Validación**: Server-side en todos los formularios
- **Soft Delete**: Registros nunca se eliminan físicamente

---

## 🐛 Troubleshooting

### Error de conexión a base de datos

```bash
# Verificar que MySQL esté corriendo
sudo systemctl status mysql

# Verificar credenciales en .env
mysql -u arka_user -p
```

### Templates no cargan

```bash
# Verificar estructura de carpetas
ls -R views/

# Verificar rutas en render.go
```

### Sesiones no persisten

```bash
# Verificar SESSION_KEY en .env
# Verificar que gorilla/sessions esté instalado
go get github.com/gorilla/sessions
```

---

## 📦 Dependencias

```go
require (
    github.com/go-sql-driver/mysql v1.7.1
    github.com/gorilla/sessions v1.2.2
    golang.org/x/crypto v0.17.0
)
```

Instalar con:

```bash
go mod tidy
```

---

## 🚀 Deploy en Producción

### Build

```bash
# Compilar para Linux
GOOS=linux GOARCH=amd64 go build -o arka-code

# Compilar para Windows
GOOS=windows GOARCH=amd64 go build -o arka-code.exe
```

### Variables de Entorno en Producción

```bash
export APP_PORT=8080
export DB_HOST=tu-servidor-mysql.com
export DB_NAME=arka
export DB_USER=usuario_prod
export DB_PASSWORD=contraseña_segura
export SESSION_KEY=clave-aleatoria-muy-segura-64-caracteres-minimo
```

---

## 📄 Licencia

Este proyecto está bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE.md) para detalles.

---

## 👥 Autores

- **Ana Concha Castro** - _Desarrolladora Backend_
- **Jose Cornejo Castro** - _Desarrollador Backend_ - [GebUCSP](https://github.com/GebUCSP)
- **Pamela Villar Ticona** - _Desarrolladora Frontend_ - [PamelaVillar](https://github.com/pamelavillar)
- **Rodrigo Silva Murillo** - _Desarrollador Backend_ - [GbTechh](https://github.com/gbTechh)
- **Jose Valdivia Castillo** - _Desarrollador Frontend_ - [Pochano](https://github.com/Pochano)

---

## 🙏 Agradecimientos

- Equipo de desarrollo Bugbusters
- Comunidad de Go (Golang)
- Gorilla Toolkit por las librerías de sesiones

---

## 📞 Soporte

Para reportar bugs o solicitar features, abrir un issue en:
https://github.com/bugbusters0/arka-code/issues
