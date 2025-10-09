-- Tabla: familia
CREATE TABLE familia (
    id INT AUTO_INCREMENT PRIMARY KEY,
    correo VARCHAR(255) UNIQUE NOT NULL,
    telefono VARCHAR(20),
    contrasena VARCHAR(255) NOT NULL, 
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    delete_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabla: usuario
CREATE TABLE usuario (
    id INT AUTO_INCREMENT PRIMARY KEY,
    id_familia INT NOT NULL,
    nombre VARCHAR(255) NOT NULL,
    fecha_nac DATE NOT NULL,
    rol ENUM('admin', 'miembro') DEFAULT 'miembro',
    contra_personal VARCHAR(255) NOT NULL, 
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    delete_at TIMESTAMP NULL DEFAULT NULL,
    FOREIGN KEY (id_familia) REFERENCES familia(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- Tabla: iconos
CREATE TABLE iconos (
    id INT AUTO_INCREMENT PRIMARY KEY,
    path VARCHAR(500) NOT NULL,  
    nombre VARCHAR(255) NOT NULL,
    size INT DEFAULT NULL, 
    type VARCHAR(50) NOT NULL, 
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    delete_at TIMESTAMP NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabla: concepto
CREATE TABLE concepto (
    id INT AUTO_INCREMENT PRIMARY KEY,
    id_icono INT NOT NULL,
    nombre VARCHAR(255) NOT NULL,
    color VARCHAR(7) DEFAULT NULL,  -- e.g., '#FF0000' para hex color
    tipo ENUM('ingreso', 'gasto') NOT NULL,
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    delete_at TIMESTAMP NULL DEFAULT NULL, 
    FOREIGN KEY (id_icono) REFERENCES iconos(id) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- Tabla: concepto_usuarios (tabla de unión con campos adicionales)
CREATE TABLE concepto_usuarios (
    id INT AUTO_INCREMENT PRIMARY KEY,
    id_usuario INT NOT NULL,
    id_concepto INT NOT NULL,
    desembolso_planejado DECIMAL(10, 2) NOT NULL,  -- Monto planeado
    periodo_desembolso TIMESTAMP DEFAULT NULL,
    periodo_tipo ENUM('mensual', 'quinicenal', 'semanal', 'diario') DEFAULT NULL,
    periodo_dia ENUM('lunes', 'martes', 'miercoles', 'jueves', 'viernes', 'sabado', 'domingo') DEFAULT NULL,
    limite_monto DECIMAL(10, 2) NOT NULL,
    limite_tipo ENUM('mensual', 'quinicenal', 'semanal', 'diario') DEFAULT NULL,
    limite_fecha TIMESTAMP DEFAULT NULL,
    has_notificacion TINYINT(1) DEFAULT 0,  -- Boolean: 1=true, 0=false
    visible TINYINT(1) DEFAULT 1,  -- Boolean por defecto true
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    delete_at TIMESTAMP NULL DEFAULT NULL, 
    UNIQUE KEY unique_usuario_concepto (id_usuario, id_concepto),  -- Evita duplicados
    FOREIGN KEY (id_usuario) REFERENCES usuario(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (id_concepto) REFERENCES concepto(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Tabla: movimiento
CREATE TABLE movimiento (
    id INT AUTO_INCREMENT PRIMARY KEY,
    id_concepto INT NOT NULL,
    monto DECIMAL(10, 2) NOT NULL,  -- Monto con decimales para finanzas
    fecha DATETIME NOT NULL,  -- Fecha del movimiento
    descripcion TEXT DEFAULT NULL,  -- Opcional, asumido del diagrama
    delete_at TIMESTAMP NULL DEFAULT NULL,  -- Soft delete
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (id_concepto) REFERENCES concepto(id) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Índices adicionales para rendimiento (opcional, pero recomendado)
CREATE INDEX idx_usuario_familia ON usuario(id_familia);
CREATE INDEX idx_concepto_tipo ON concepto(tipo);
CREATE INDEX idx_movimiento_concepto ON movimiento(id_concepto);
CREATE INDEX idx_movimiento_fecha ON movimiento(fecha);

-- Datos de prueba (opcional: inserta para testing)
-- INSERT INTO familia (nombre, correo, telefono, contrasena) VALUES ('Familia Test', 'test@familia.com', '123456789', '$2y$10$hashedpassword');
-- INSERT INTO iconos (path, nombre, size, type) VALUES ('/icons/gasto.svg', 'Gasto General', 24, 'svg');
-- INSERT INTO concepto (id_icono, nombre, color, tipo) VALUES (1, 'Comida', '#FF6B6B', 'gasto');
-- INSERT INTO usuario (id_familia, nombre, fecha_nac, rol, contrasena_personal) VALUES (1, 'Juan Pérez', '1990-01-01', 'admin', '$2y$10$hashedpassword');
-- INSERT INTO concepto_usuarios (id_usuario, id_concepto, desembolso_planejado, periodo_desembolso) VALUES (1, 1, 500, 'mensual');
-- INSERT INTO movimiento (id_concepto, monto, fecha) VALUES (1, -50.00, NOW());