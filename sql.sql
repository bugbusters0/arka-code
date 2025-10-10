-- ===============================================
-- SCRIPT SQL CORREGIDO Y COMPATIBLE
-- ===============================================

-- Eliminar tablas si existen (para reiniciar sin errores)
DROP TABLE IF EXISTS movimiento;
DROP TABLE IF EXISTS concepto_usuarios;
DROP TABLE IF EXISTS concepto;
DROP TABLE IF EXISTS iconos;
DROP TABLE IF EXISTS usuario;
DROP TABLE IF EXISTS familia;

-- ===============================================
-- Tabla: familia
-- ===============================================
CREATE TABLE familia (
    id INT AUTO_INCREMENT PRIMARY KEY,
    correo VARCHAR(255) UNIQUE NOT NULL,
    telefono VARCHAR(20),
    contrasena VARCHAR(255) NOT NULL, 
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    delete_at TIMESTAMP NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- ===============================================
-- Tabla: usuario
-- ===============================================
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
    FOREIGN KEY (id_familia) REFERENCES familia(id) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- ===============================================
-- Tabla: iconos
-- ===============================================
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


-- ===============================================
-- Tabla: concepto
-- ===============================================
CREATE TABLE concepto (
    id INT AUTO_INCREMENT PRIMARY KEY,
    id_icono INT NOT NULL,
    nombre VARCHAR(255) NOT NULL,
    color VARCHAR(7) DEFAULT NULL,  -- Ejemplo: '#FF0000'
    tipo ENUM('ingreso', 'gasto') NOT NULL,
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    delete_at TIMESTAMP NULL DEFAULT NULL, 
    FOREIGN KEY (id_icono) REFERENCES iconos(id) 
        ON DELETE RESTRICT 
        ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- ===============================================
-- Tabla: concepto_usuarios
-- ===============================================
CREATE TABLE concepto_usuarios (
    id INT AUTO_INCREMENT PRIMARY KEY,
    id_usuario INT NOT NULL,
    id_concepto INT NOT NULL,
    desembolso_planejado DECIMAL(10, 2) NOT NULL,
    periodo_desembolso DATETIME DEFAULT NULL,
    periodo_tipo ENUM('mensual', 'quincenal', 'semanal', 'diario') DEFAULT NULL,
    periodo_dia ENUM('lunes', 'martes', 'miercoles', 'jueves', 'viernes', 'sabado', 'domingo') DEFAULT NULL,
    limite_monto DECIMAL(10, 2) NOT NULL,
    limite_tipo ENUM('mensual', 'quincenal', 'semanal', 'diario') DEFAULT NULL,
    limite_fecha DATETIME DEFAULT NULL,
    has_notificacion TINYINT(1) DEFAULT 0,
    visible TINYINT(1) DEFAULT 1,
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    delete_at TIMESTAMP NULL DEFAULT NULL, 
    UNIQUE KEY unique_usuario_concepto (id_usuario, id_concepto),
    FOREIGN KEY (id_usuario) REFERENCES usuario(id) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE,
    FOREIGN KEY (id_concepto) REFERENCES concepto(id) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- ===============================================
-- Tabla: movimiento
-- ===============================================
CREATE TABLE movimiento (
    id INT AUTO_INCREMENT PRIMARY KEY,
    id_concepto INT NOT NULL,
    monto DECIMAL(10, 2) NOT NULL,
    fecha DATETIME NOT NULL,
    descripcion TEXT DEFAULT NULL,
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    delete_at TIMESTAMP NULL DEFAULT NULL,
    FOREIGN KEY (id_concepto) REFERENCES concepto(id) 
        ON DELETE RESTRICT 
        ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- ===============================================
-- Índices recomendados
-- ===============================================
CREATE INDEX idx_usuario_familia ON usuario(id_familia);
CREATE INDEX idx_concepto_tipo ON concepto(tipo);
CREATE INDEX idx_movimiento_concepto ON movimiento(id_concepto);
CREATE INDEX idx_movimiento_fecha ON movimiento(fecha);



INSERT INTO `iconos` (`id`, `path`, `nombre`, `size`, `type`, `create_at`, `update_at`, `delete_at`) VALUES
(3, 'fa-solid fa-burger', 'Hamburguesa', NULL, 'font-awesome', '2025-10-10 03:31:09', '2025-10-10 03:31:09', NULL),
(4, 'fa-solid fa-spa', 'Spa', NULL, 'font-awesome', '2025-10-10 03:31:09', '2025-10-10 03:31:09', NULL),
(5, 'fa-solid fa-car', 'Auto', NULL, 'font-awesome', '2025-10-10 03:31:09', '2025-10-10 03:31:09', NULL),
(6, 'fa-solid fa-heart-pulse', 'Salud', NULL, 'font-awesome', '2025-10-10 03:31:09', '2025-10-10 03:31:09', NULL),
(7, 'fa-solid fa-graduation-cap', 'Educación', NULL, 'font-awesome', '2025-10-10 03:31:09', '2025-10-10 03:31:09', NULL),
(8, 'fa-solid fa-basket-shopping', 'Canasta de Compras', NULL, 'font-awesome', '2025-10-10 03:31:09', '2025-10-10 03:31:09', NULL),
(9, 'fa-solid fa-utensils', 'Utensilios', NULL, 'font-awesome', '2025-10-10 03:31:09', '2025-10-10 03:31:09', NULL),
(10, 'fa-solid fa-lightbulb', 'Bombilla', NULL, 'font-awesome', '2025-10-10 03:31:09', '2025-10-10 03:31:09', NULL),
(11, 'fa-solid fa-shirt', 'Camisa', NULL, 'font-awesome', '2025-10-10 03:31:09', '2025-10-10 03:31:09', NULL);

-- ===============================================
-- Ejemplo de inserción de datos de prueba (opcional)
-- ===============================================
-- INSERT INTO familia (correo, telefono, contrasena) 
-- VALUES ('test@familia.com', '123456789', '$2y$10$hashedpassword');

-- INSERT INTO iconos (path, nombre, size, type) 
-- VALUES ('/icons/gasto.svg', 'Gasto General', 24, 'svg');

-- INSERT INTO concepto (id_icono, nombre, color, tipo) 
-- VALUES (1, 'Comida', '#FF6B6B', 'gasto');

-- INSERT INTO usuario (id_familia, nombre, fecha_nac, rol, contra_personal) 
-- VALUES (1, 'Juan Pérez', '1990-01-01', 'admin', '$2y$10$hashedpassword');

-- INSERT INTO concepto_usuarios (id_usuario, id_concepto, desembolso_planejado, limite_monto) 
-- VALUES (1, 1, 500.00, 600.00);

-- INSERT INTO movimiento (id_concepto, monto, fecha, descripcion) 
-- VALUES (1, -50.00, NOW(), 'Compra de almuerzo');
