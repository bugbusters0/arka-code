-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Servidor: 127.0.0.1
-- Tiempo de generación: 24-10-2025 a las 11:11:10
-- Versión del servidor: 10.4.32-MariaDB
-- Versión de PHP: 8.2.12

CREATE DATABASE IF NOT EXISTS arka_go;
USE arka_go;

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Base de datos: `arka`
--

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `concepto`
--

CREATE TABLE `concepto` (
  `nombreConcepto` varchar(40) NOT NULL,
  `correoFamilia` varchar(50) NOT NULL,
  `tipo` tinyint(1) NOT NULL,
  `icono` varchar(100) DEFAULT NULL,
  `color` char(7) DEFAULT NULL,
  `nombreUsuario` varchar(30) NOT NULL,
  `delete_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `familia`
--

CREATE TABLE `familia` (
  `correo` varchar(50) NOT NULL,
  `telefono` char(9) DEFAULT NULL CHECK (`telefono` regexp '^[0-9]{9}$'),
  `contraseña` varchar(255) NOT NULL,
  `delete_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `movimiento`
--

CREATE TABLE `movimiento` (
  `idMovimiento` int(11) NOT NULL,
  `fecha` date NOT NULL,
  `monto` decimal(10,2) NOT NULL,
  `descripcion` text DEFAULT NULL,
  `nombreUsuario` varchar(30) NOT NULL,
  `nombreConcepto` varchar(40) NOT NULL,
  `correoFamilia` varchar(50) NOT NULL,
  `delete_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `personalizacionconcepto`
--

CREATE TABLE `personalizacionconcepto` (
  `idPersonalizacion` int(11) NOT NULL,
  `montoPlanificado` decimal(10,2) DEFAULT NULL,
  `tipoPeriodoPlanificado` varchar(15) DEFAULT NULL,
  `diaPeriodoPlanificado` tinyint(3) UNSIGNED DEFAULT NULL CHECK (`diaPeriodoPlanificado` between 1 and 31),
  `limiteGasto` decimal(10,2) DEFAULT NULL,
  `tipoPeriodoLimite` varchar(15) DEFAULT NULL,
  `diaPeriodoLimite` tinyint(3) UNSIGNED DEFAULT NULL CHECK (`diaPeriodoLimite` between 1 and 31),
  `notificacion` tinyint(1) DEFAULT 0,
  `activo` tinyint(1) DEFAULT 1,
  `nombreUsuario` varchar(30) NOT NULL,
  `nombreConcepto` varchar(40) NOT NULL,
  `correoFamilia` varchar(50) NOT NULL,
  `delete_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Estructura de tabla para la tabla `usuario`
--

CREATE TABLE `usuario` (
  `nombreUsuario` varchar(30) NOT NULL,
  `rol` tinyint(1) NOT NULL DEFAULT 0,
  `contraseñaPersonal` varchar(255) NOT NULL,
  `nombrePersonal` varchar(100) NOT NULL,
  `correoFamilia` varchar(50) NOT NULL,
  `delete_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Índices para tablas volcadas
--

--
-- Indices de la tabla `concepto`
--
ALTER TABLE `concepto`
  ADD PRIMARY KEY (`nombreConcepto`,`correoFamilia`),
  ADD KEY `correoFamilia` (`correoFamilia`),
  ADD KEY `nombreUsuario` (`nombreUsuario`);

--
-- Indices de la tabla `familia`
--
ALTER TABLE `familia`
  ADD PRIMARY KEY (`correo`);

--
-- Indices de la tabla `movimiento`
--
ALTER TABLE `movimiento`
  ADD PRIMARY KEY (`idMovimiento`),
  ADD KEY `nombreUsuario` (`nombreUsuario`),
  ADD KEY `nombreConcepto` (`nombreConcepto`,`correoFamilia`);

--
-- Indices de la tabla `personalizacionconcepto`
--
ALTER TABLE `personalizacionconcepto`
  ADD PRIMARY KEY (`idPersonalizacion`),
  ADD KEY `nombreUsuario` (`nombreUsuario`),
  ADD KEY `nombreConcepto` (`nombreConcepto`,`correoFamilia`);

--
-- Indices de la tabla `usuario`
--
ALTER TABLE `usuario`
  ADD PRIMARY KEY (`nombreUsuario`),
  ADD KEY `correoFamilia` (`correoFamilia`);

--
-- AUTO_INCREMENT de las tablas volcadas
--

--
-- AUTO_INCREMENT de la tabla `movimiento`
--
ALTER TABLE `movimiento`
  MODIFY `idMovimiento` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT de la tabla `personalizacionconcepto`
--
ALTER TABLE `personalizacionconcepto`
  MODIFY `idPersonalizacion` int(11) NOT NULL AUTO_INCREMENT;

--
-- Restricciones para tablas volcadas
--

--
-- Filtros para la tabla `concepto`
--
ALTER TABLE `concepto`
  ADD CONSTRAINT `concepto_ibfk_1` FOREIGN KEY (`correoFamilia`) REFERENCES `familia` (`correo`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `concepto_ibfk_2` FOREIGN KEY (`nombreUsuario`) REFERENCES `usuario` (`nombreUsuario`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Filtros para la tabla `movimiento`
--
ALTER TABLE `movimiento`
  ADD CONSTRAINT `movimiento_ibfk_1` FOREIGN KEY (`nombreUsuario`) REFERENCES `usuario` (`nombreUsuario`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `movimiento_ibfk_2` FOREIGN KEY (`nombreConcepto`,`correoFamilia`) REFERENCES `concepto` (`nombreConcepto`, `correoFamilia`) ON UPDATE CASCADE;

--
-- Filtros para la tabla `personalizacionconcepto`
--
ALTER TABLE `personalizacionconcepto`
  ADD CONSTRAINT `personalizacionconcepto_ibfk_1` FOREIGN KEY (`nombreUsuario`) REFERENCES `usuario` (`nombreUsuario`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `personalizacionconcepto_ibfk_2` FOREIGN KEY (`nombreConcepto`,`correoFamilia`) REFERENCES `concepto` (`nombreConcepto`, `correoFamilia`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Filtros para la tabla `usuario`
--
ALTER TABLE `usuario`
  ADD CONSTRAINT `usuario_ibfk_1` FOREIGN KEY (`correoFamilia`) REFERENCES `familia` (`correo`) ON DELETE CASCADE ON UPDATE CASCADE;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;


-- Insertar datos en la tabla Familia
INSERT INTO familia (correo, telefono, contraseña) VALUES
('familia.garcia@email.com', '912345678', '$2y$10$rQdS6XG2wWz8pZkLmNqPf.AfQbTgHjKcYvMxNs1qRs2wDtUlVpB0C'),
('familia.rodriguez@email.com', '923456789', '$2y$10$sRfT7YH3xX9qZmMnRrQgQ.BgRcUhIjLdZwNy2oRt3wExEuVmWqC1D'),
('familia.lopez@email.com', '934567890', '$2y$10$tSgU8ZI4yY0rAnNoSsRhR.CgSdVJkMeAxOz3pSu4xEyFvNwXnXr2E'),
('familia.martinez@email.com', '945678901', '$2y$10$uThV9AJ5zZ1sBpOpTtSiS.DhTeKlNfByPa4qRv5yFzGwOxYoYr3fG'),
('familia.gonzalez@email.com', '956789012', '$2y$10$vUiW0BK6A1tCpQqUuTjTt.EiUfLmOgCzQb5sSw6zGyHyPzZpZs4hH'),
('familia.hernandez@email.com', '967890123', '$2y$10$wVjX1CL7B2uDqRrVuUkUu.FjVgMnPhDzRc6tTx7zHzIzQaQqAt5iI'),
('familia.perez@email.com', '978901234', '$2y$10$xWkY2DM8C3vEsSrWvVlVv.GkWgNoQiEySd7uTy8aIAJzRbRrBu6jK'),
('familia.sanchez@email.com', '989012345', '$2y$10$yXlZ3EN9D4wFtTwWwWmWw.HlXhOpRjFzTe8vUz9bJBKScSsCrv7kL'),
('familia.ramirez@email.com', '990123456', '$2y$10$zYm4FO0E5xGuUxXxXnXx.ImYiQqSkGzUf0wVz0cKCKLTdTtDs8mM'),
('familia.torres@email.com', '901234567', '$2y$10$aZn5GP1F6yHvVyYyYoYy.JnZjRrTlH0Ag1xWz1dLDKMUeUuEt9nN'),
('familia.flores@email.com', '912345679', '$2y$10$bAo6HQ2G7zIwWzZzZpZz.KoAkSsUmI1Bh2yXz2eMEKLVfVvFu0oO'),
('familia.rivera@email.com', '923456780', '$2y$10$cBp7IR3H8AjXxAaAaQaQa.LpBlTtVmJ2Ci3yY3fNFKMWgWwGv1pP'),
('familia.gomez@email.com', '934567891', '$2y$10$dCq8JS4I9BkYyBbBbRbRb.MqCmUuWnK3Dj4zZ4gOGLNXhXxHw2qQ'),
('familia.diaz@email.com', '945678902', '$2y$10$eDr9KT5J0ClZzCcCcScSc.NrDnVvXoL4Ek5aA5hPHMOYiYyIx3rR'),
('familia.reyes@email.com', '956789013', '$2y$10$fEs0LU6K1DmAaDdDdTdTd.OsEoWwYpM5Fl6bB6iQINPjZzJy4sS'),
('familia.morales@email.com', '967890124', '$2y$10$gFt1MV7L2EnBbEeEeUeUe.PtFpXxZqN6Gm7cC7jRJQqKaKzKz5tT'),
('familia.ortiz@email.com', '978901235', '$2y$10$hGu2NW8M3FoCcFfFfVfVf.QtGqYyArO7Hn8dD8kSKRrLbLzLz6uU'),
('familia.silva@email.com', '989012346', '$2y$10$iHv3OX9N4GpDdGgGgWgWg.RuHrZzBsP8Io9eE9lTLsMcMcMmA7vV'),
('familia.vargas@email.com', '990123457', '$2y$10$jIw4PY0O5HqEeHhHhXhXh.SvIsAaCtQ9Jp0fF0mUMTtNdNdNnB8w'),
('familia.castro@email.com', '901234568', '$2y$10$kJx5QZ1P6IrFfIiIiYiYi.TwJtBbDuR0Kq1gG1nVNUVuOuOuOoCx'),
('familia.romero@email.com', '912345680', '$2y$10$lKy6Ra2Q7JsGgJjJjZjZj.UxKuCcEvS1Lr2hH2oWOVwvPvPvPpDy'),
('familia.aguilar@email.com', '923456791', '$2y$10$mLz7Sb3R8KtHhKkKkAkAk.VyLvDdFwT2Ms3iI3pXPWxwQwQwQqEz'),
('familia.mendoza@email.com', '934567892', '$2y$10$nMa8Tc4S9LuIiLlLlBlBl.WzMwEeGxU3Nt4jJ4qYQXYxYxYxYrF'),
('familia.herrera@email.com', '945678903', '$2y$10$oNb9Ud5T0MvJjMmMmCmCm.XaNxFfHyV4Ou5kK5rZRYZyZyZyZsG'),
('familia.guzman@email.com', '956789014', '$2y$10$pOc0Ve6U1NwKkKnKnDnDn.YbOoYgIzW5Pv6lL6sASZaZaZaZzAt'),
('familia.jimenez@email.com', '967890125', '$2y$10$qPd1Wf7V2OxLlLoLoEoEo.ZcPpZhJzX6Qw7mM7tBTAbAbAbAaBu'),
('familia.ruiz@email.com', '978901236', '$2y$10$rQe2Xg8W3PyMmMpMpFpFp.AdQqAiKaY7Rx8nN8uCUBcBcBcBbCv'),
('familia.salazar@email.com', '989012347', '$2y$10$sRf3Yh9X4QzNnNqNqGqGq.BeRrBjLbZ8Sy9oO9vDVCdCdCdCcDw'),
('familia.delgado@email.com', '990123458', '$2y$10$tSg4Zi0Y5RaOoOrOrHrHr.CfSsCkMcA9Tz0pP0wEWDeDeDeDdEx'),
('familia.molina@email.com', '901234569', '$2y$10$uTh5Aj1Z6SbPpPsPsIsIs.DgTtDlNdB0Ua1qQ1xFXEfEfEfEeFy'),
('familia.rios@email.com', '912345681', '$2y$10$vUi6Bk2A7TcQqQtQtJtJt.EhUuEmOeC1Vb2rR2yGYFgFgFgGfGz'),
('familia.iglesias@email.com', '923456792', '$2y$10$wVj7Cl3B8UdRrRuRuKuKu.FiVvFnPfD2Wc3sS3zHZGhGhGhHgHa'),
('familia.medina@email.com', '934567893', '$2y$10$xWk8Dm4C9VeSsSvSvLvLv.GjWwGoQgE3Xd4tT4aIaHiHiHiIhIb'),
('familia.nunez@email.com', '945678904', '$2y$10$yXl9En5D0WfTtTwTwMwMw.HkXxHpRhF4Ye5uU5bJbIjIjIjJiJc'),
('familia.leon@email.com', '956789015', '$2y$10$zYm0Fo6E1XgUuUxUxNxNx.IlYyYqSkG5Zf6vV6cKcKkKkKkKjKd'),
('familia.miranda@email.com', '967890126', '$2y$10$aZn1Gp7F2YhVvVyVyOyOy.JmZjZrTlH6Ag7wWz7dLdLlLlLlLkLe'),
('familia.cortes@email.com', '978901237', '$2y$10$bAo2Hq8G3ZiWwWzWzPzPz.KnAkAsUmI7Bh8yXz8eMeMmMmMmMlMf'),
('familia.santos@email.com', '989012348', '$2y$10$cBp3Ir9H4AjXxAxAxQxQx.LoBlBtVmJ8Ci9zY9fNfNnNnNnNkNg'),
('familia.vega@email.com', '990123459', '$2y$10$dCq4Js0I5BkYyByByRyRy.MpCmCvWnK0Dj0aA0gOgOoOoOoOjOh'),
('familia.campos@email.com', '901234570', '$2y$10$eDr5Kt1J6ClZzCzCzSzSz.NqDnDwXoL1Ek1bB1hPhPpPpPpPiPj'),
('familia.marquez@email.com', '912345682', '$2y$10$fEs6Lu2K7DmAaDaDaTaTa.OsEoEwYpM2Fl2cC2iQiQqQqQqQhQk'),
('familia.figueroa@email.com', '923456793', '$2y$10$gFt7Mv3L8EnBbEbEbUbUb.PtFpFxZqN3Gm3dD3jRjRrRrRrRgRl'),
('familia.mejia@email.com', '934567894', '$2y$10$hGu8Nw4M9FoCcFcFcVcVc.QtGqGyArO4Hn4eE4kSkSsSsSsSfSm'),
('familia.carrillo@email.com', '945678905', '$2y$10$iHv9Ox5N0GpDdGdGdWdWd.RuHrRzBsP5Io5fF5lTlTtTtTtTeTn'),
('familia.arias@email.com', '956789016', '$2y$10$jIw0Py1O5HqEeHeHeXeXe.SvIsAsCtQ0Jp0gG0mUmUuUuUuUdUo'),
('familia.espinoza@email.com', '967890127', '$2y$10$kJx1Qz2P6IrFfIfIfYfYf.TwJtBtDuR1Kq1hH1nVnVvVvVvVeVp'),
('familia.contreras@email.com', '978901238', '$2y$10$lKy2Ra3Q7JsGgJgJgZgZg.UxKuCuEvS2Lr2iI2oWoWwWwWwWfWq'),
('familia.valdez@email.com', '989012349', '$2y$10$mLz3Sb4R8KtHhKhKhAhAh.VyLvDvFwT3Ms3jI3pXpXxXxXxXgXr'),
('familia.rosales@email.com', '990123460', '$2y$10$nMa4Tc5S9LuIiLiLiBiBi.WzMwEwGxU4Nt4kK4qYqYyYyYyYhYs'),
('familia.soto@email.com', '901234571', '$2y$10$oNb5Ud6T0MvJjMjMjCjCj.XaNxXfHyV5Ou5lL5rZrZzZzZzZiZt');


-- Insertar datos en la tabla Usuario (3-4 integrantes por familia, 1 admin por familia)
INSERT INTO usuario (nombreUsuario, rol, contraseñaPersonal, nombrePersonal, correoFamilia) VALUES
-- Familia García (4 integrantes)
('juan.garcia', 1, '$2y$10$abc123def456ghi789jkl', 'Juan García Pérez', 'familia.garcia@email.com'),
('maria.garcia', 0, '$2y$10$mno234pqr567stu890vwx', 'María García López', 'familia.garcia@email.com'),
('carlos.garcia', 0, '$2y$10$yza345bcd678efg901hij', 'Carlos García Martínez', 'familia.garcia@email.com'),
('laura.garcia', 0, '$2y$10$klm456nop789qrs012tuv', 'Laura García González', 'familia.garcia@email.com'),

-- Familia Rodríguez (3 integrantes)
('pedro.rodriguez', 1, '$2y$10$uvw567xyz890abc123def', 'Pedro Rodríguez Hernández', 'familia.rodriguez@email.com'),
('ana.rodriguez', 0, '$2y$10$ghi678jkl901mno234pqr', 'Ana Rodríguez Pérez', 'familia.rodriguez@email.com'),
('luis.rodriguez', 0, '$2y$10$stu789vwx012yza345bcd', 'Luis Rodríguez Sánchez', 'familia.rodriguez@email.com'),

-- Familia López (4 integrantes)
('miguel.lopez', 1, '$2y$10$efg890hij123klm456nop', 'Miguel López Ramírez', 'familia.lopez@email.com'),
('elena.lopez', 0, '$2y$10$qrs012tuv234uvw567xyz', 'Elena López Torres', 'familia.lopez@email.com'),
('javier.lopez', 0, '$2y$10$abc234def567ghi890jkl', 'Javier López Flores', 'familia.lopez@email.com'),
('carmen.lopez', 0, '$2y$10$mno345pqr678stu901vwx', 'Carmen López Rivera', 'familia.lopez@email.com'),

-- Familia Martínez (3 integrantes)
('david.martinez', 1, '$2y$10$yza456bcd789efg012hij', 'David Martínez Gómez', 'familia.martinez@email.com'),
('patricia.martinez', 0, '$2y$10$klm567nop890qrs123tuv', 'Patricia Martínez Díaz', 'familia.martinez@email.com'),
('jorge.martinez', 0, '$2y$10$uvw678xyz901abc234def', 'Jorge Martínez Reyes', 'familia.martinez@email.com'),

-- Familia González (4 integrantes)
('sandra.gonzalez', 1, '$2y$10$ghi789jkl012mno345pqr', 'Sandra González Morales', 'familia.gonzalez@email.com'),
('ricardo.gonzalez', 0, '$2y$10$stu890vwx123yza456bcd', 'Ricardo González Ortiz', 'familia.gonzalez@email.com'),
('monica.gonzalez', 0, '$2y$10$efg901hij234klm567nop', 'Mónica González Silva', 'familia.gonzalez@email.com'),
('fernando.gonzalez', 0, '$2y$10$qrs123tuv345uvw678xyz', 'Fernando González Vargas', 'familia.gonzalez@email.com'),

-- Familia Hernández (3 integrantes)
('raquel.hernandez', 1, '$2y$10$abc345def678ghi901jkl', 'Raquel Hernández Castro', 'familia.hernandez@email.com'),
('alejandro.hernandez', 0, '$2y$10$mno456pqr789stu012vwx', 'Alejandro Hernández Romero', 'familia.hernandez@email.com'),
('lucia.hernandez', 0, '$2y$10$yza567bcd890efg123hij', 'Lucía Hernández Aguilar', 'familia.hernandez@email.com'),

-- Familia Pérez (4 integrantes)
('sergio.perez', 1, '$2y$10$klm678nop901qrs234tuv', 'Sergio Pérez Mendoza', 'familia.perez@email.com'),
('elvira.perez', 0, '$2y$10$uvw789xyz012abc345def', 'Elvira Pérez Herrera', 'familia.perez@email.com'),
('oscar.perez', 0, '$2y$10$ghi890jkl123mno456pqr', 'Óscar Pérez Guzmán', 'familia.perez@email.com'),
('adriana.perez', 0, '$2y$10$stu901vwx234yza567bcd', 'Adriana Pérez Jiménez', 'familia.perez@email.com'),

-- Familia Sánchez (3 integrantes)
('humberto.sanchez', 1, '$2y$10$efg012hij345klm678nop', 'Humberto Sánchez Ruiz', 'familia.sanchez@email.com'),
('veronica.sanchez', 0, '$2y$10$qrs234tuv456uvw789xyz', 'Verónica Sánchez Salazar', 'familia.sanchez@email.com'),
('arturo.sanchez', 0, '$2y$10$abc456def789ghi012jkl', 'Arturo Sánchez Delgado', 'familia.sanchez@email.com'),

-- Familia Ramírez (4 integrantes)
('gloria.ramirez', 1, '$2y$10$mno567pqr890stu123vwx', 'Gloria Ramírez Molina', 'familia.ramirez@email.com'),
('manuel.ramirez', 0, '$2y$10$yza678bcd901efg234hij', 'Manuel Ramírez Ríos', 'familia.ramirez@email.com'),
('teresa.ramirez', 0, '$2y$10$klm789nop012qrs345tuv', 'Teresa Ramírez Iglesias', 'familia.ramirez@email.com'),
('guillermo.ramirez', 0, '$2y$10$uvw890xyz123abc456def', 'Guillermo Ramírez Medina', 'familia.ramirez@email.com'),

-- Familia Torres (3 integrantes)
('silvia.torres', 1, '$2y$10$ghi901jkl234mno567pqr', 'Silvia Torres Núñez', 'familia.torres@email.com'),
('ramon.torres', 0, '$2y$10$stu012vwx345yza678bcd', 'Ramón Torres León', 'familia.torres@email.com'),
('beatriz.torres', 0, '$2y$10$efg123hij456klm789nop', 'Beatriz Torres Miranda', 'familia.torres@email.com'),

-- Familia Flores (4 integrantes)
('francisco.flores', 1, '$2y$10$qrs345tuv567uvw890xyz', 'Francisco Flores Cortés', 'familia.flores@email.com'),
('lourdes.flores', 0, '$2y$10$abc567def890ghi123jkl', 'Lourdes Flores Santos', 'familia.flores@email.com'),
('victor.flores', 0, '$2y$10$mno678pqr901stu234vwx', 'Víctor Flores Vega', 'familia.flores@email.com'),
('rosa.flores', 0, '$2y$10$yza789bcd012efg345hij', 'Rosa Flores Campos', 'familia.flores@email.com'),

-- Familia Rivera (3 integrantes)
('eduardo.rivera', 1, '$2y$10$klm890nop123qrs456tuv', 'Eduardo Rivera Márquez', 'familia.rivera@email.com'),
('natalia.rivera', 0, '$2y$10$uvw901xyz234abc567def', 'Natalia Rivera Figueroa', 'familia.rivera@email.com'),
('alfonso.rivera', 0, '$2y$10$ghi012jkl345mno678pqr', 'Alfonso Rivera Mejía', 'familia.rivera@email.com'),

-- Familia Gómez (4 integrantes)
('ines.gomez', 1, '$2y$10$stu123vwx456yza789bcd', 'Inés Gómez Carrillo', 'familia.gomez@email.com'),
('joaquin.gomez', 0, '$2y$10$efg234hij567klm890nop', 'Joaquín Gómez Arias', 'familia.gomez@email.com'),
('margarita.gomez', 0, '$2y$10$qrs456tuv678uvw901xyz', 'Margarita Gómez Espinoza', 'familia.gomez@email.com'),
('agustin.gomez', 0, '$2y$10$abc678def901ghi234jkl', 'Agustín Gómez Contreras', 'familia.gomez@email.com'),

-- Familia Díaz (3 integrantes)
('celia.diaz', 1, '$2y$10$mno789pqr012stu345vwx', 'Celia Díaz Valdez', 'familia.diaz@email.com'),
('sebastian.diaz', 0, '$2y$10$yza890bcd123efg456hij', 'Sebastián Díaz Rosales', 'familia.diaz@email.com'),
('olga.diaz', 0, '$2y$10$klm901nop234qrs567tuv', 'Olga Díaz Soto', 'familia.diaz@email.com'),

-- Familia Reyes (4 integrantes)
('antonio.reyes', 1, '$2y$10$uvw012xyz345abc678def', 'Antonio Reyes García', 'familia.reyes@email.com'),
('isabel.reyes', 0, '$2y$10$ghi123jkl456mno789pqr', 'Isabel Reyes Rodríguez', 'familia.reyes@email.com'),
('roberto.reyes', 0, '$2y$10$stu234vwx567yza890bcd', 'Roberto Reyes López', 'familia.reyes@email.com'),
('elena.reyes', 0, '$2y$10$efg345hij678klm901nop', 'Elena Reyes Martínez', 'familia.reyes@email.com'),

-- Familia Morales (3 integrantes)
('pablo.morales', 1, '$2y$10$qrs456tuv789uvw012xyz', 'Pablo Morales González', 'familia.morales@email.com'),
('claudia.morales', 0, '$2y$10$abc789def012ghi345jkl', 'Claudia Morales Hernández', 'familia.morales@email.com'),
('daniel.morales', 0, '$2y$10$mno012pqr345stu678vwx', 'Daniel Morales Pérez', 'familia.morales@email.com'),

-- Familia Ortiz (4 integrantes)
('lucia.ortiz', 1, '$2y$10$yza123bcd456efg789hij', 'Lucía Ortiz Sánchez', 'familia.ortiz@email.com'),
('mario.ortiz', 0, '$2y$10$klm234nop567qrs890tuv', 'Mario Ortiz Ramírez', 'familia.ortiz@email.com'),
('susana.ortiz', 0, '$2y$10$uvw345xyz678abc901def', 'Susana Ortiz Torres', 'familia.ortiz@email.com'),
('jose.ortiz', 0, '$2y$10$ghi456jkl789mno012pqr', 'José Ortiz Flores', 'familia.ortiz@email.com'),

-- Familia Silva (3 integrantes)
('andrea.silva', 1, '$2y$10$stu567vwx890yza123bcd', 'Andrea Silva Rivera', 'familia.silva@email.com'),
('rodrigo.silva', 0, '$2y$10$efg678hij901klm234nop', 'Rodrigo Silva Gómez', 'familia.silva@email.com'),
('camila.silva', 0, '$2y$10$qrs789tuv012uvw345xyz', 'Camila Silva Díaz', 'familia.silva@email.com'),

-- Familia Vargas (4 integrantes)
('ricardo.vargas', 1, '$2y$10$abc012def345ghi678jkl', 'Ricardo Vargas Reyes', 'familia.vargas@email.com'),
('fernanda.vargas', 0, '$2y$10$mno123pqr456stu789vwx', 'Fernanda Vargas Morales', 'familia.vargas@email.com'),
('hugo.vargas', 0, '$2y$10$yza234bcd567efg890hij', 'Hugo Vargas Ortiz', 'familia.vargas@email.com'),
('valeria.vargas', 0, '$2y$10$klm345nop678qrs901tuv', 'Valeria Vargas Silva', 'familia.vargas@email.com'),

-- Familia Castro (3 integrantes)
('gabriel.castro', 1, '$2y$10$uvw456xyz789abc012def', 'Gabriel Castro Vargas', 'familia.castro@email.com'),
('diana.castro', 0, '$2y$10$ghi567jkl012mno345pqr', 'Diana Castro Castro', 'familia.castro@email.com'),
('tomas.castro', 0, '$2y$10$stu678vwx123yza456bcd', 'Tomás Castro Romero', 'familia.castro@email.com'),

-- Familia Romero (4 integrantes)
('alejandra.romero', 1, '$2y$10$efg789hij234klm567nop', 'Alejandra Romero Aguilar', 'familia.romero@email.com'),
('santiago.romero', 0, '$2y$10$qrs890tuv345uvw678xyz', 'Santiago Romero Mendoza', 'familia.romero@email.com'),
('paula.romero', 0, '$2y$10$abc123def456ghi789jkl', 'Paula Romero Herrera', 'familia.romero@email.com'),
('martin.romero', 0, '$2y$10$mno234pqr567stu890vwx', 'Martín Romero Guzmán', 'familia.romero@email.com'),

-- Familia Aguilar (3 integrantes)
('raul.aguilar', 1, '$2y$10$yza345bcd678efg901hij', 'Raúl Aguilar Jiménez', 'familia.aguilar@email.com'),
('lorena.aguilar', 0, '$2y$10$klm456nop789qrs012tuv', 'Lorena Aguilar Ruiz', 'familia.aguilar@email.com'),
('esteban.aguilar', 0, '$2y$10$uvw567xyz890abc123def', 'Esteban Aguilar Salazar', 'familia.aguilar@email.com'),

-- Familia Mendoza (4 integrantes)
('clara.mendoza', 1, '$2y$10$ghi678jkl901mno234pqr', 'Clara Mendoza Delgado', 'familia.mendoza@email.com'),
('felipe.mendoza', 0, '$2y$10$stu789vwx012yza345bcd', 'Felipe Mendoza Molina', 'familia.mendoza@email.com'),
('renata.mendoza', 0, '$2y$10$efg890hij123klm456nop', 'Renata Mendoza Ríos', 'familia.mendoza@email.com'),
('simon.mendoza', 0, '$2y$10$qrs012tuv234uvw567xyz', 'Simón Mendoza Iglesias', 'familia.mendoza@email.com');

-- Continuación para completar las 50 familias...
-- Familia Herrera (3 integrantes)
('julio.herrera', 1, '$2y$10$abc234def567ghi890jkl', 'Julio Herrera Medina', 'familia.herrera@email.com'),
('vanessa.herrera', 0, '$2y$10$mno345pqr678stu901vwx', 'Vanessa Herrera Núñez', 'familia.herrera@email.com'),
('marcos.herrera', 0, '$2y$10$yza456bcd789efg012hij', 'Marcos Herrera León', 'familia.herrera@email.com'),

-- Familia Guzmán (4 integrantes)
('natalia.guzman', 1, '$2y$10$klm567nop890qrs123tuv', 'Natalia Guzmán Miranda', 'familia.guzman@email.com'),
('leonardo.guzman', 0, '$2y$10$uvw678xyz901abc234def', 'Leonardo Guzmán Cortés', 'familia.guzman@email.com'),
('ximena.guzman', 0, '$2y$10$ghi789jkl012mno345pqr', 'Ximena Guzmán Santos', 'familia.guzman@email.com'),
('octavio.guzman', 0, '$2y$10$stu890vwx123yza456bcd', 'Octavio Guzmán Vega', 'familia.guzman@email.com'),

-- Familia Jiménez (3 integrantes)
('diego.jimenez', 1, '$2y$10$efg901hij234klm567nop', 'Diego Jiménez Campos', 'familia.jimenez@email.com'),
('sofia.jimenez', 0, '$2y$10$qrs123tuv345uvw678xyz', 'Sofía Jiménez Márquez', 'familia.jimenez@email.com'),
('emilio.jimenez', 0, '$2y$10$abc345def678ghi901jkl', 'Emilio Jiménez Figueroa', 'familia.jimenez@email.com'),

-- Familia Ruiz (4 integrantes)
('carolina.ruiz', 1, '$2y$10$mno456pqr789stu012vwx', 'Carolina Ruiz Mejía', 'familia.ruiz@email.com'),
('javier.ruiz', 0, '$2y$10$yza567bcd890efg123hij', 'Javier Ruiz Carrillo', 'familia.ruiz@email.com'),
('adriana.ruiz', 0, '$2y$10$klm678nop901qrs234tuv', 'Adriana Ruiz Arias', 'familia.ruiz@email.com'),
('raul.ruiz', 0, '$2y$10$uvw789xyz012abc345def', 'Raúl Ruiz Espinoza', 'familia.ruiz@email.com'),

-- Familia Salazar (3 integrantes)
('patricio.salazar', 1, '$2y$10$ghi890jkl123mno456pqr', 'Patricio Salazar Contreras', 'familia.salazar@email.com'),
('liliana.salazar', 0, '$2y$10$stu901vwx234yza567bcd', 'Liliana Salazar Valdez', 'familia.salazar@email.com'),
('ernesto.salazar', 0, '$2y$10$efg012hij345klm678nop', 'Ernesto Salazar Rosales', 'familia.salazar@email.com'),

-- Familia Delgado (4 integrantes)
('veronica.delgado', 1, '$2y$10$qrs234tuv456uvw789xyz', 'Verónica Delgado Soto', 'familia.delgado@email.com'),
('cristian.delgado', 0, '$2y$10$abc456def789ghi012jkl', 'Cristian Delgado García', 'familia.delgado@email.com'),
('isabella.delgado', 0, '$2y$10$mno567pqr890stu123vwx', 'Isabella Delgado Rodríguez', 'familia.delgado@email.com'),
('federico.delgado', 0, '$2y$10$yza678bcd901efg234hij', 'Federico Delgado López', 'familia.delgado@email.com'),

-- Familia Molina (3 integrantes)
('oscar.molina', 1, '$2y$10$klm789nop012qrs345tuv', 'Óscar Molina Martínez', 'familia.molina@email.com'),
('elisa.molina', 0, '$2y$10$uvw890xyz123abc456def', 'Elisa Molina González', 'familia.molina@email.com'),
('ricardo.molina', 0, '$2y$10$ghi901jkl234mno567pqr', 'Ricardo Molina Hernández', 'familia.molina@email.com'),

-- Familia Ríos (4 integrantes)
('lucia.rios', 1, '$2y$10$stu012vwx345yza678bcd', 'Lucía Ríos Pérez', 'familia.rios@email.com'),
('miguel.rios', 0, '$2y$10$efg123hij456klm789nop', 'Miguel Ríos Sánchez', 'familia.rios@email.com'),
('catalina.rios', 0, '$2y$10$qrs345tuv567uvw890xyz', 'Catalina Ríos Ramírez', 'familia.rios@email.com'),
('andres.rios', 0, '$2y$10$abc567def890ghi123jkl', 'Andrés Ríos Torres', 'familia.rios@email.com'),

-- Familia Iglesias (3 integrantes)
('jorge.iglesias', 1, '$2y$10$mno678pqr901stu234vwx', 'Jorge Iglesias Flores', 'familia.iglesias@email.com'),
('beatriz.iglesias', 0, '$2y$10$yza789bcd012efg345hij', 'Beatriz Iglesias Rivera', 'familia.iglesias@email.com'),
('rodolfo.iglesias', 0, '$2y$10$klm890nop123qrs456tuv', 'Rodolfo Iglesias Gómez', 'familia.iglesias@email.com'),

-- Familia Medina (4 integrantes)
('silvia.medina', 1, '$2y$10$uvw901xyz234abc567def', 'Silvia Medina Díaz', 'familia.medina@email.com'),
('guillermo.medina', 0, '$2y$10$ghi012jkl345mno678pqr', 'Guillermo Medina Reyes', 'familia.medina@email.com'),
('elena.medina', 0, '$2y$10$stu123vwx456yza789bcd', 'Elena Medina Morales', 'familia.medina@email.com'),
('pablo.medina', 0, '$2y$10$efg234hij567klm890nop', 'Pablo Medina Ortiz', 'familia.medina@email.com'),

-- Familia Núñez (3 integrantes)
('carmen.nunez', 1, '$2y$10$qrs456tuv678uvw901xyz', 'Carmen Núñez Silva', 'familia.nunez@email.com'),
('sergio.nunez', 0, '$2y$10$abc678def901ghi234jkl', 'Sergio Núñez Vargas', 'familia.nunez@email.com'),
('laura.nunez', 0, '$2y$10$mno789pqr012stu345vwx', 'Laura Núñez Castro', 'familia.nunez@email.com'),

-- Familia León (4 integrantes)
('ramiro.leon', 1, '$2y$10$yza890bcd123efg456hij', 'Ramiro León Romero', 'familia.leon@email.com'),
('teresa.leon', 0, '$2y$10$klm901nop234qrs567tuv', 'Teresa León Aguilar', 'familia.leon@email.com'),
('humberto.leon', 0, '$2y$10$uvw012xyz345abc678def', 'Humberto León Mendoza', 'familia.leon@email.com'),
('claudia.leon', 0, '$2y$10$ghi123jkl456mno789pqr', 'Claudia León Herrera', 'familia.leon@email.com'),

-- Familia Miranda (3 integrantes)
('esteban.miranda', 1, '$2y$10$stu234vwx567yza890bcd', 'Esteban Miranda Guzmán', 'familia.miranda@email.com'),
('valentina.miranda', 0, '$2y$10$efg345hij678klm901nop', 'Valentina Miranda Jiménez', 'familia.miranda@email.com'),
('ricardo.miranda', 0, '$2y$10$qrs456tuv789uvw012xyz', 'Ricardo Miranda Ruiz', 'familia.miranda@email.com'),

-- Familia Cortés (4 integrantes)
('lorenzo.cortes', 1, '$2y$10$abc789def012ghi345jkl', 'Lorenzo Cortés Salazar', 'familia.cortes@email.com'),
('mariana.cortes', 0, '$2y$10$mno012pqr345stu678vwx', 'Mariana Cortés Delgado', 'familia.cortes@email.com'),
('fernando.cortes', 0, '$2y$10$yza123bcd456efg789hij', 'Fernando Cortés Molina', 'familia.cortes@email.com'),
('isabel.cortes', 0, '$2y$10$klm234nop567qrs890tuv', 'Isabel Cortés Ríos', 'familia.cortes@email.com'),

-- Familia Santos (3 integrantes)
('raul.santos', 1, '$2y$10$uvw345xyz678abc901def', 'Raúl Santos Iglesias', 'familia.santos@email.com'),
('elena.santos', 0, '$2y$10$ghi456jkl789mno012pqr', 'Elena Santos Medina', 'familia.santos@email.com'),
('javier.santos', 0, '$2y$10$stu567vwx890yza123bcd', 'Javier Santos Núñez', 'familia.santos@email.com'),

-- Familia Vega (4 integrantes)
('cristina.vega', 1, '$2y$10$efg678hij901klm234nop', 'Cristina Vega León', 'familia.vega@email.com'),
('alberto.vega', 0, '$2y$10$qrs789tuv012uvw345xyz', 'Alberto Vega Miranda', 'familia.vega@email.com'),
('lucia.vega', 0, '$2y$10$abc012def345ghi678jkl', 'Lucía Vega Cortés', 'familia.vega@email.com'),
('roberto.vega', 0, '$2y$10$mno123pqr456stu789vwx', 'Roberto Vega Santos', 'familia.vega@email.com'),

-- Familia Campos (3 integrantes)
('diana.campos', 1, '$2y$10$yza234bcd567efg890hij', 'Diana Campos Vega', 'familia.campos@email.com'),
('marcos.campos', 0, '$2y$10$klm345nop678qrs901tuv', 'Marcos Campos Campos', 'familia.campos@email.com'),
('natalia.campos', 0, '$2y$10$uvw456xyz789abc012def', 'Natalia Campos Márquez', 'familia.campos@email.com'),

-- Familia Márquez (4 integrantes)
('julio.marquez', 1, '$2y$10$ghi567jkl012mno345pqr', 'Julio Márquez Figueroa', 'familia.marquez@email.com'),
('carla.marquez', 0, '$2y$10$stu678vwx123yza456bcd', 'Carla Márquez Mejía', 'familia.marquez@email.com'),
('antonio.marquez', 0, '$2y$10$efg789hij234klm567nop', 'Antonio Márquez Carrillo', 'familia.marquez@email.com'),
('patricia.marquez', 0, '$2y$10$qrs890tuv345uvw678xyz', 'Patricia Márquez Arias', 'familia.marquez@email.com'),

-- Familia Figueroa (3 integrantes)
('manuel.figueroa', 1, '$2y$10$abc123def456ghi789jkl', 'Manuel Figueroa Espinoza', 'familia.figueroa@email.com'),
('lidia.figueroa', 0, '$2y$10$mno234pqr567stu890vwx', 'Lidia Figueroa Contreras', 'familia.figueroa@email.com'),
('hugo.figueroa', 0, '$2y$10$yza345bcd678efg901hij', 'Hugo Figueroa Valdez', 'familia.figueroa@email.com'),

-- Familia Mejía (4 integrantes)
('rosa.mejia', 1, '$2y$10$klm456nop789qrs012tuv', 'Rosa Mejía Rosales', 'familia.mejia@email.com'),
('eduardo.mejia', 0, '$2y$10$uvw567xyz890abc123def', 'Eduardo Mejía Soto', 'familia.mejia@email.com'),
('gloria.mejia', 0, '$2y$10$ghi678jkl901mno234pqr', 'Gloria Mejía García', 'familia.mejia@email.com'),
('ricardo.mejia', 0, '$2y$10$stu789vwx012yza345bcd', 'Ricardo Mejía Rodríguez', 'familia.mejia@email.com'),

-- Familia Carrillo (3 integrantes)
('alfonso.carrillo', 1, '$2y$10$efg890hij123klm456nop', 'Alfonso Carrillo López', 'familia.carrillo@email.com'),
('susana.carrillo', 0, '$2y$10$qrs012tuv234uvw567xyz', 'Susana Carrillo Martínez', 'familia.carrillo@email.com'),
('jose.carrillo', 0, '$2y$10$abc234def567ghi890jkl', 'José Carrillo González', 'familia.carrillo@email.com'),

-- Familia Arias (4 integrantes)
('lucia.arias', 1, '$2y$10$mno345pqr678stu901vwx', 'Lucía Arias Hernández', 'familia.arias@email.com'),
('carlos.arias', 0, '$2y$10$yza456bcd789efg012hij', 'Carlos Arias Pérez', 'familia.arias@email.com'),
('elena.arias', 0, '$2y$10$klm567nop890qrs123tuv', 'Elena Arias Sánchez', 'familia.arias@email.com'),
('miguel.arias', 0, '$2y$10$uvw678xyz901abc234def', 'Miguel Arias Ramírez', 'familia.arias@email.com'),

-- Familia Espinoza (3 integrantes)
('jorge.espinoza', 1, '$2y$10$ghi789jkl012mno345pqr', 'Jorge Espinoza Torres', 'familia.espinoza@email.com'),
('claudia.espinoza', 0, '$2y$10$stu890vwx123yza456bcd', 'Claudia Espinoza Flores', 'familia.espinoza@email.com'),
('rodrigo.espinoza', 0, '$2y$10$efg901hij234klm567nop', 'Rodrigo Espinoza Rivera', 'familia.espinoza@email.com'),

-- Familia Contreras (4 integrantes)
('maria.contreras', 1, '$2y$10$qrs123tuv345uvw678xyz', 'María Contreras Gómez', 'familia.contreras@email.com'),
('pablo.contreras', 0, '$2y$10$abc345def678ghi901jkl', 'Pablo Contreras Díaz', 'familia.contreras@email.com'),
('laura.contreras', 0, '$2y$10$mno456pqr789stu012vwx', 'Laura Contreras Reyes', 'familia.contreras@email.com'),
('david.contreras', 0, '$2y$10$yza567bcd890efg123hij', 'David Contreras Morales', 'familia.contreras@email.com'),

-- Familia Valdez (3 integrantes)
('ana.valdez', 1, '$2y$10$klm678nop901qrs234tuv', 'Ana Valdez Ortiz', 'familia.valdez@email.com'),
('juan.valdez', 0, '$2y$10$uvw789xyz012abc345def', 'Juan Valdez Silva', 'familia.valdez@email.com'),
('carmen.valdez', 0, '$2y$10$ghi890jkl123mno456pqr', 'Carmen Valdez Vargas', 'familia.valdez@email.com'),

-- Familia Rosales (4 integrantes)
('roberto.rosales', 1, '$2y$10$stu901vwx234yza567bcd', 'Roberto Rosales Castro', 'familia.rosales@email.com'),
('elena.rosales', 0, '$2y$10$efg012hij345klm678nop', 'Elena Rosales Romero', 'familia.rosales@email.com'),
('mario.rosales', 0, '$2y$10$qrs234tuv456uvw789xyz', 'Mario Rosales Aguilar', 'familia.rosales@email.com'),
('lucia.rosales', 0, '$2y$10$abc456def789ghi012jkl', 'Lucía Rosales Mendoza', 'familia.rosales@email.com'),

-- Familia Soto (3 integrantes)
('diego.soto', 1, '$2y$10$mno567pqr890stu123vwx', 'Diego Soto Herrera', 'familia.soto@email.com'),
('patricia.soto', 0, '$2y$10$yza678bcd901efg234hij', 'Patricia Soto Guzmán', 'familia.soto@email.com'),
('raul.soto', 0, '$2y$10$klm789nop012qrs345tuv', 'Raúl Soto Jiménez', 'familia.soto@email.com');



-- Insertar datos en la tabla Concepto
INSERT INTO concepto (nombreConcepto, correoFamilia, tipo, icono, color, nombreUsuario) VALUES
-- Familia García (4 conceptos - 2 gastos, 2 ingresos)
('Alimentación', 'familia.garcia@email.com', 0, '🛒', '#FF6B6B', 'juan.garcia'),
('Transporte', 'familia.garcia@email.com', 0, '🚗', '#4ECDC4', 'juan.garcia'),
('Salario', 'familia.garcia@email.com', 1, '💰', '#45B7D1', 'juan.garcia'),
('Inversiones', 'familia.garcia@email.com', 1, '📈', '#96CEB4', 'juan.garcia'),

-- Familia Rodríguez (4 conceptos - 2 gastos, 2 ingresos)
('Vivienda', 'familia.rodriguez@email.com', 0, '🏠', '#FFEAA7', 'pedro.rodriguez'),
('Servicios', 'familia.rodriguez@email.com', 0, '💡', '#DDA0DD', 'pedro.rodriguez'),
('Freelance', 'familia.rodriguez@email.com', 1, '💻', '#98D8C8', 'pedro.rodriguez'),
('Bonos', 'familia.rodriguez@email.com', 1, '🎯', '#F7DC6F', 'pedro.rodriguez'),

-- Familia López (4 conceptos - 2 gastos, 2 ingresos)
('Educación', 'familia.lopez@email.com', 0, '📚', '#BB8FCE', 'miguel.lopez'),
('Salud', 'familia.lopez@email.com', 0, '🏥', '#85C1E9', 'miguel.lopez'),
('Comisiones', 'familia.lopez@email.com', 1, '📊', '#F8C471', 'miguel.lopez'),
('Regalos', 'familia.lopez@email.com', 1, '🎁', '#82E0AA', 'miguel.lopez'),

-- Familia Martínez (4 conceptos - 2 gastos, 2 ingresos)
('Entretenimiento', 'familia.martinez@email.com', 0, '🎬', '#F1948A', 'david.martinez'),
('Ropa', 'familia.martinez@email.com', 0, '👕', '#A9CCE3', 'david.martinez'),
('Alquiler Ingresos', 'familia.martinez@email.com', 1, '🏢', '#D7BDE2', 'david.martinez'),
('Dividendos', 'familia.martinez@email.com', 1, '📘', '#A3E4D7', 'david.martinez'),

-- Familia González (4 conceptos - 2 gastos, 2 ingresos)
('Deportes', 'familia.gonzalez@email.com', 0, '⚽', '#E59866', 'sandra.gonzalez'),
('Viajes', 'familia.gonzalez@email.com', 0, '✈', '#C39BD3', 'sandra.gonzalez'),
('Ventas', 'familia.gonzalez@email.com', 1, '🛒', '#76D7C4', 'sandra.gonzalez'),
('Intereses', 'familia.gonzalez@email.com', 1, '📈', '#F7DC6F', 'sandra.gonzalez'),

-- Familia Hernández (4 conceptos - 2 gastos, 2 ingresos)
('Restaurantes', 'familia.hernandez@email.com', 0, '🍽', '#D7BDE2', 'raquel.hernandez'),
('Seguros', 'familia.hernandez@email.com', 0, '🛡', '#A9DFBF', 'raquel.hernandez'),
('Consultoría', 'familia.hernandez@email.com', 1, '💼', '#F5B7B1', 'raquel.hernandez'),
('Herencia', 'familia.hernandez@email.com', 1, '📜', '#AED6F1', 'raquel.hernandez'),

-- Familia Pérez (4 conceptos - 2 gastos, 2 ingresos)
('Impuestos', 'familia.perez@email.com', 0, '📝', '#ABEBC6', 'sergio.perez'),
('Mascotas', 'familia.perez@email.com', 0, '🐕', '#FAD7A0', 'sergio.perez'),
('Préstamo', 'familia.perez@email.com', 1, '🏦', '#D2B4DE', 'sergio.perez'),
('Regalías', 'familia.perez@email.com', 1, '📖', '#AED6F1', 'sergio.perez'),

-- Familia Sánchez (4 conceptos - 2 gastos, 2 ingresos)
('Tecnología', 'familia.sanchez@email.com', 0, '💻', '#F9E79F', 'humberto.sanchez'),
('Hogar', 'familia.sanchez@email.com', 0, '🛋', '#D5DBDB', 'humberto.sanchez'),
('Sueldo', 'familia.sanchez@email.com', 1, '💰', '#F1948A', 'humberto.sanchez'),
('Inversiones', 'familia.sanchez@email.com', 1, '📊', '#82E0AA', 'humberto.sanchez'),

-- Familia Ramírez (4 conceptos - 2 gastos, 2 ingresos)
('Cuidado Personal', 'familia.ramirez@email.com', 0, '💄', '#BB8FCE', 'gloria.ramirez'),
('Supermercado', 'familia.ramirez@email.com', 0, '🛒', '#85C1E9', 'gloria.ramirez'),
('Salario', 'familia.ramirez@email.com', 1, '💵', '#F8C471', 'gloria.ramirez'),
('Bonos', 'familia.ramirez@email.com', 1, '🎁', '#82E0AA', 'gloria.ramirez'),

-- Familia Torres (4 conceptos - 2 gastos, 2 ingresos)
('Cine', 'familia.torres@email.com', 0, '🎬', '#F1948A', 'silvia.torres'),
('Gimnasio', 'familia.torres@email.com', 0, '💪', '#A9CCE3', 'silvia.torres'),
('Freelance', 'familia.torres@email.com', 1, '💻', '#D7BDE2', 'silvia.torres'),
('Comisiones', 'familia.torres@email.com', 1, '📈', '#A3E4D7', 'silvia.torres'),

-- Familia Flores (4 conceptos - 2 gastos, 2 ingresos)
('Farmacia', 'familia.flores@email.com', 0, '💊', '#E59866', 'francisco.flores'),
('Libros', 'familia.flores@email.com', 0, '📖', '#C39BD3', 'francisco.flores'),
('Alquiler', 'familia.flores@email.com', 1, '🏠', '#76D7C4', 'francisco.flores'),
('Dividendos', 'familia.flores@email.com', 1, '📊', '#F7DC6F', 'francisco.flores'),

-- Familia Rivera (4 conceptos - 2 gastos, 2 ingresos)
('Electricidad', 'familia.rivera@email.com', 0, '💡', '#D7BDE2', 'eduardo.rivera'),
('Agua', 'familia.rivera@email.com', 0, '💧', '#A9DFBF', 'eduardo.rivera'),
('Ventas', 'familia.rivera@email.com', 1, '🛒', '#F5B7B1', 'eduardo.rivera'),
('Intereses', 'familia.rivera@email.com', 1, '💰', '#AED6F1', 'eduardo.rivera'),

-- Familia Gómez (4 conceptos - 2 gastos, 2 ingresos)
('Internet', 'familia.gomez@email.com', 0, '🌐', '#ABEBC6', 'ines.gomez'),
('Teléfono', 'familia.gomez@email.com', 0, '📱', '#FAD7A0', 'ines.gomez'),
('Consultoría', 'familia.gomez@email.com', 1, '💼', '#D2B4DE', 'ines.gomez'),
('Herencia', 'familia.gomez@email.com', 1, '📜', '#AED6F1', 'ines.gomez'),

-- Familia Díaz (4 conceptos - 2 gastos, 2 ingresos)
('Gas', 'familia.diaz@email.com', 0, '🔥', '#F9E79F', 'celia.diaz'),
('Transporte Público', 'familia.diaz@email.com', 0, '🚌', '#D5DBDB', 'celia.diaz'),
('Préstamo', 'familia.diaz@email.com', 1, '🏦', '#F1948A', 'celia.diaz'),
('Regalías', 'familia.diaz@email.com', 1, '📖', '#82E0AA', 'celia.diaz'),

-- Familia Reyes (4 conceptos - 2 gastos, 2 ingresos)
('Taxi', 'familia.reyes@email.com', 0, '🚕', '#BB8FCE', 'antonio.reyes'),
('Mantenimiento', 'familia.reyes@email.com', 0, '🔧', '#85C1E9', 'antonio.reyes'),
('Sueldo', 'familia.reyes@email.com', 1, '💵', '#F8C471', 'antonio.reyes'),
('Inversiones', 'familia.reyes@email.com', 1, '📈', '#82E0AA', 'antonio.reyes'),

-- Familia Morales (4 conceptos - 2 gastos, 2 ingresos)
('Seguro Médico', 'familia.morales@email.com', 0, '🏥', '#F1948A', 'pablo.morales'),
('Clases', 'familia.morales@email.com', 0, '🎓', '#A9CCE3', 'pablo.morales'),
('Salario', 'familia.morales@email.com', 1, '💰', '#D7BDE2', 'pablo.morales'),
('Bonos', 'familia.morales@email.com', 1, '🎯', '#A3E4D7', 'pablo.morales'),

-- Familia Ortiz (4 conceptos - 2 gastos, 2 ingresos)
('Gasolina', 'familia.ortiz@email.com', 0, '⛽', '#E59866', 'lucia.ortiz'),
('Alimentación', 'familia.ortiz@email.com', 0, '🛒', '#C39BD3', 'lucia.ortiz'),
('Freelance', 'familia.ortiz@email.com', 1, '💻', '#76D7C4', 'lucia.ortiz'),
('Comisiones', 'familia.ortiz@email.com', 1, '📊', '#F7DC6F', 'lucia.ortiz'),

-- Familia Silva (4 conceptos - 2 gastos, 2 ingresos)
('Vivienda', 'familia.silva@email.com', 0, '🏠', '#D7BDE2', 'andrea.silva'),
('Servicios', 'familia.silva@email.com', 0, '💡', '#A9DFBF', 'andrea.silva'),
('Alquiler Ingresos', 'familia.silva@email.com', 1, '🏢', '#F5B7B1', 'andrea.silva'),
('Dividendos', 'familia.silva@email.com', 1, '📈', '#AED6F1', 'andrea.silva'),

-- Familia Vargas (4 conceptos - 2 gastos, 2 ingresos)
('Educación', 'familia.vargas@email.com', 0, '📚', '#ABEBC6', 'ricardo.vargas'),
('Salud', 'familia.vargas@email.com', 0, '🏥', '#FAD7A0', 'ricardo.vargas'),
('Ventas', 'familia.vargas@email.com', 1, '🛒', '#D2B4DE', 'ricardo.vargas'),
('Intereses', 'familia.vargas@email.com', 1, '💰', '#AED6F1', 'ricardo.vargas'),

-- Familia Castro (4 conceptos - 2 gastos, 2 ingresos)
('Entretenimiento', 'familia.castro@email.com', 0, '🎬', '#F9E79F', 'gabriel.castro'),
('Ropa', 'familia.castro@email.com', 0, '👕', '#D5DBDB', 'gabriel.castro'),
('Consultoría', 'familia.castro@email.com', 1, '💼', '#F1948A', 'gabriel.castro'),
('Herencia', 'familia.castro@email.com', 1, '📜', '#82E0AA', 'gabriel.castro'),

-- Familia Romero (4 conceptos - 2 gastos, 2 ingresos)
('Deportes', 'familia.romero@email.com', 0, '⚽', '#BB8FCE', 'alejandra.romero'),
('Viajes', 'familia.romero@email.com', 0, '✈', '#85C1E9', 'alejandra.romero'),
('Préstamo', 'familia.romero@email.com', 1, '🏦', '#F8C471', 'alejandra.romero'),
('Regalías', 'familia.romero@email.com', 1, '📖', '#82E0AA', 'alejandra.romero'),

-- Familia Aguilar (4 conceptos - 2 gastos, 2 ingresos)
('Restaurantes', 'familia.aguilar@email.com', 0, '🍽', '#F1948A', 'raul.aguilar'),
('Seguros', 'familia.aguilar@email.com', 0, '🛡', '#A9CCE3', 'raul.aguilar'),
('Sueldo', 'familia.aguilar@email.com', 1, '💵', '#D7BDE2', 'raul.aguilar'),
('Inversiones', 'familia.aguilar@email.com', 1, '📈', '#A3E4D7', 'raul.aguilar'),

-- Familia Mendoza (4 conceptos - 2 gastos, 2 ingresos)
('Impuestos', 'familia.mendoza@email.com', 0, '📝', '#E59866', 'clara.mendoza'),
('Mascotas', 'familia.mendoza@email.com', 0, '🐕', '#C39BD3', 'clara.mendoza'),
('Salario', 'familia.mendoza@email.com', 1, '💰', '#76D7C4', 'clara.mendoza'),
('Bonos', 'familia.mendoza@email.com', 1, '🎯', '#F7DC6F', 'clara.mendoza'),

-- Familia Herrera (4 conceptos - 2 gastos, 2 ingresos)
('Tecnología', 'familia.herrera@email.com', 0, '💻', '#D7BDE2', 'julio.herrera'),
('Hogar', 'familia.herrera@email.com', 0, '🛋', '#A9DFBF', 'julio.herrera'),
('Freelance', 'familia.herrera@email.com', 1, '💻', '#F5B7B1', 'julio.herrera'),
('Comisiones', 'familia.herrera@email.com', 1, '📊', '#AED6F1', 'julio.herrera'),

-- Familia Guzmán (4 conceptos - 2 gastos, 2 ingresos)
('Cuidado Personal', 'familia.guzman@email.com', 0, '💄', '#ABEBC6', 'natalia.guzman'),
('Supermercado', 'familia.guzman@email.com', 0, '🛒', '#FAD7A0', 'natalia.guzman'),
('Alquiler Ingresos', 'familia.guzman@email.com', 1, '🏢', '#D2B4DE', 'natalia.guzman'),
('Dividendos', 'familia.guzman@email.com', 1, '📈', '#AED6F1', 'natalia.guzman'),

-- Familia Jiménez (4 conceptos - 2 gastos, 2 ingresos)
('Cine', 'familia.jimenez@email.com', 0, '🎬', '#F9E79F', 'diego.jimenez'),
('Gimnasio', 'familia.jimenez@email.com', 0, '💪', '#D5DBDB', 'diego.jimenez'),
('Ventas', 'familia.jimenez@email.com', 1, '🛒', '#F1948A', 'diego.jimenez'),
('Intereses', 'familia.jimenez@email.com', 1, '💰', '#82E0AA', 'diego.jimenez'),

-- Familia Ruiz (4 conceptos - 2 gastos, 2 ingresos)
('Farmacia', 'familia.ruiz@email.com', 0, '💊', '#BB8FCE', 'carolina.ruiz'),
('Libros', 'familia.ruiz@email.com', 0, '📖', '#85C1E9', 'carolina.ruiz'),
('Consultoría', 'familia.ruiz@email.com', 1, '💼', '#F8C471', 'carolina.ruiz'),
('Herencia', 'familia.ruiz@email.com', 1, '📜', '#82E0AA', 'carolina.ruiz'),

-- Familia Salazar (4 conceptos - 2 gastos, 2 ingresos)
('Electricidad', 'familia.salazar@email.com', 0, '💡', '#F1948A', 'patricio.salazar'),
('Agua', 'familia.salazar@email.com', 0, '💧', '#A9CCE3', 'patricio.salazar'),
('Préstamo', 'familia.salazar@email.com', 1, '🏦', '#D7BDE2', 'patricio.salazar'),
('Regalías', 'familia.salazar@email.com', 1, '📖', '#A3E4D7', 'patricio.salazar'),

-- Familia Delgado (4 conceptos - 2 gastos, 2 ingresos)
('Internet', 'familia.delgado@email.com', 0, '🌐', '#E59866', 'veronica.delgado'),
('Teléfono', 'familia.delgado@email.com', 0, '📱', '#C39BD3', 'veronica.delgado'),
('Sueldo', 'familia.delgado@email.com', 1, '💵', '#76D7C4', 'veronica.delgado'),
('Inversiones', 'familia.delgado@email.com', 1, '📈', '#F7DC6F', 'veronica.delgado'),

-- Familia Molina (4 conceptos - 2 gastos, 2 ingresos)
('Gas', 'familia.molina@email.com', 0, '🔥', '#D7BDE2', 'oscar.molina'),
('Transporte Público', 'familia.molina@email.com', 0, '🚌', '#A9DFBF', 'oscar.molina'),
('Salario', 'familia.molina@email.com', 1, '💰', '#F5B7B1', 'oscar.molina'),
('Bonos', 'familia.molina@email.com', 1, '🎯', '#AED6F1', 'oscar.molina'),

-- Familia Ríos (4 conceptos - 2 gastos, 2 ingresos)
('Taxi', 'familia.rios@email.com', 0, '🚕', '#ABEBC6', 'lucia.rios'),
('Mantenimiento', 'familia.rios@email.com', 0, '🔧', '#FAD7A0', 'lucia.rios'),
('Freelance', 'familia.rios@email.com', 1, '💻', '#D2B4DE', 'lucia.rios'),
('Comisiones', 'familia.rios@email.com', 1, '📊', '#AED6F1', 'lucia.rios'),

-- Familia Iglesias (4 conceptos - 2 gastos, 2 ingresos)
('Seguro Médico', 'familia.iglesias@email.com', 0, '🏥', '#F9E79F', 'jorge.iglesias'),
('Clases', 'familia.iglesias@email.com', 0, '🎓', '#D5DBDB', 'jorge.iglesias'),
('Alquiler Ingresos', 'familia.iglesias@email.com', 1, '🏢', '#F1948A', 'jorge.iglesias'),
('Dividendos', 'familia.iglesias@email.com', 1, '📈', '#82E0AA', 'jorge.iglesias'),

-- Familia Medina (4 conceptos - 2 gastos, 2 ingresos)
('Gasolina', 'familia.medina@email.com', 0, '⛽', '#BB8FCE', 'silvia.medina'),
('Alimentación', 'familia.medina@email.com', 0, '🛒', '#85C1E9', 'silvia.medina'),
('Ventas', 'familia.medina@email.com', 1, '🛒', '#F8C471', 'silvia.medina'),
('Intereses', 'familia.medina@email.com', 1, '💰', '#82E0AA', 'silvia.medina'),

-- Familia Núñez (4 conceptos - 2 gastos, 2 ingresos)
('Vivienda', 'familia.nunez@email.com', 0, '🏠', '#F1948A', 'carmen.nunez'),
('Servicios', 'familia.nunez@email.com', 0, '💡', '#A9CCE3', 'carmen.nunez'),
('Consultoría', 'familia.nunez@email.com', 1, '💼', '#D7BDE2', 'carmen.nunez'),
('Herencia', 'familia.nunez@email.com', 1, '📜', '#A3E4D7', 'carmen.nunez'),

-- Familia León (4 conceptos - 2 gastos, 2 ingresos)
('Educación', 'familia.leon@email.com', 0, '📚', '#E59866', 'ramiro.leon'),
('Salud', 'familia.leon@email.com', 0, '🏥', '#C39BD3', 'ramiro.leon'),
('Préstamo', 'familia.leon@email.com', 1, '🏦', '#76D7C4', 'ramiro.leon'),
('Regalías', 'familia.leon@email.com', 1, '📖', '#F7DC6F', 'ramiro.leon'),

-- Familia Miranda (4 conceptos - 2 gastos, 2 ingresos)
('Entretenimiento', 'familia.miranda@email.com', 0, '🎬', '#D7BDE2', 'esteban.miranda'),
('Ropa', 'familia.miranda@email.com', 0, '👕', '#A9DFBF', 'esteban.miranda'),
('Sueldo', 'familia.miranda@email.com', 1, '💵', '#F5B7B1', 'esteban.miranda'),
('Inversiones', 'familia.miranda@email.com', 1, '📈', '#AED6F1', 'esteban.miranda'),

-- Familia Cortés (4 conceptos - 2 gastos, 2 ingresos)
('Deportes', 'familia.cortes@email.com', 0, '⚽', '#ABEBC6', 'lorenzo.cortes'),
('Viajes', 'familia.cortes@email.com', 0, '✈', '#FAD7A0', 'lorenzo.cortes'),
('Salario', 'familia.cortes@email.com', 1, '💰', '#D2B4DE', 'lorenzo.cortes'),
('Bonos', 'familia.cortes@email.com', 1, '🎯', '#AED6F1', 'lorenzo.cortes'),

-- Familia Santos (4 conceptos - 2 gastos, 2 ingresos)
('Restaurantes', 'familia.santos@email.com', 0, '🍽', '#F9E79F', 'raul.santos'),
('Seguros', 'familia.santos@email.com', 0, '🛡', '#D5DBDB', 'raul.santos'),
('Freelance', 'familia.santos@email.com', 1, '💻', '#F1948A', 'raul.santos'),
('Comisiones', 'familia.santos@email.com', 1, '📊', '#82E0AA', 'raul.santos'),

-- Familia Vega (4 conceptos - 2 gastos, 2 ingresos)
('Impuestos', 'familia.vega@email.com', 0, '📝', '#BB8FCE', 'cristina.vega'),
('Mascotas', 'familia.vega@email.com', 0, '🐕', '#85C1E9', 'cristina.vega'),
('Alquiler Ingresos', 'familia.vega@email.com', 1, '🏢', '#F8C471', 'cristina.vega'),
('Dividendos', 'familia.vega@email.com', 1, '📈', '#82E0AA', 'cristina.vega'),

-- Familia Campos (4 conceptos - 2 gastos, 2 ingresos)
('Tecnología', 'familia.campos@email.com', 0, '💻', '#F1948A', 'diana.campos'),
('Hogar', 'familia.campos@email.com', 0, '🛋', '#A9CCE3', 'diana.campos'),
('Ventas', 'familia.campos@email.com', 1, '🛒', '#D7BDE2', 'diana.campos'),
('Intereses', 'familia.campos@email.com', 1, '💰', '#A3E4D7', 'diana.campos'),

-- Familia Márquez (4 conceptos - 2 gastos, 2 ingresos)
('Cuidado Personal', 'familia.marquez@email.com', 0, '💄', '#E59866', 'julio.marquez'),
('Supermercado', 'familia.marquez@email.com', 0, '🛒', '#C39BD3', 'julio.marquez'),
('Consultoría', 'familia.marquez@email.com', 1, '💼', '#76D7C4', 'julio.marquez'),
('Herencia', 'familia.marquez@email.com', 1, '📜', '#F7DC6F', 'julio.marquez'),

-- Familia Figueroa (4 conceptos - 2 gastos, 2 ingresos)
('Cine', 'familia.figueroa@email.com', 0, '🎬', '#D7BDE2', 'manuel.figueroa'),
('Gimnasio', 'familia.figueroa@email.com', 0, '💪', '#A9DFBF', 'manuel.figueroa'),
('Préstamo', 'familia.figueroa@email.com', 1, '🏦', '#F5B7B1', 'manuel.figueroa'),
('Regalías', 'familia.figueroa@email.com', 1, '📖', '#AED6F1', 'manuel.figueroa'),

-- Familia Mejía (4 conceptos - 2 gastos, 2 ingresos)
('Farmacia', 'familia.mejia@email.com', 0, '💊', '#ABEBC6', 'rosa.mejia'),
('Libros', 'familia.mejia@email.com', 0, '📖', '#FAD7A0', 'rosa.mejia'),
('Sueldo', 'familia.mejia@email.com', 1, '💵', '#D2B4DE', 'rosa.mejia'),
('Inversiones', 'familia.mejia@email.com', 1, '📈', '#AED6F1', 'rosa.mejia'),

-- Familia Carrillo (4 conceptos - 2 gastos, 2 ingresos)
('Electricidad', 'familia.carrillo@email.com', 0, '💡', '#F9E79F', 'alfonso.carrillo'),
('Agua', 'familia.carrillo@email.com', 0, '💧', '#D5DBDB', 'alfonso.carrillo'),
('Salario', 'familia.carrillo@email.com', 1, '💰', '#F1948A', 'alfonso.carrillo'),
('Bonos', 'familia.carrillo@email.com', 1, '🎯', '#82E0AA', 'alfonso.carrillo'),

-- Familia Arias (4 conceptos - 2 gastos, 2 ingresos)
('Internet', 'familia.arias@email.com', 0, '🌐', '#BB8FCE', 'lucia.arias'),
('Teléfono', 'familia.arias@email.com', 0, '📱', '#85C1E9', 'lucia.arias'),
('Freelance', 'familia.arias@email.com', 1, '💻', '#F8C471', 'lucia.arias'),
('Comisiones', 'familia.arias@email.com', 1, '📊', '#82E0AA', 'lucia.arias'),

-- Familia Espinoza (4 conceptos - 2 gastos, 2 ingresos)
('Gas', 'familia.espinoza@email.com', 0, '🔥', '#F1948A', 'jorge.espinoza'),
('Transporte Público', 'familia.espinoza@email.com', 0, '🚌', '#A9CCE3', 'jorge.espinoza'),
('Alquiler Ingresos', 'familia.espinoza@email.com', 1, '🏢', '#D7BDE2', 'jorge.espinoza'),
('Dividendos', 'familia.espinoza@email.com', 1, '📈', '#A3E4D7', 'jorge.espinoza'),

-- Familia Contreras (4 conceptos - 2 gastos, 2 ingresos)
('Taxi', 'familia.contreras@email.com', 0, '🚕', '#E59866', 'maria.contreras'),
('Mantenimiento', 'familia.contreras@email.com', 0, '🔧', '#C39BD3', 'maria.contreras'),
('Ventas', 'familia.contreras@email.com', 1, '🛒', '#76D7C4', 'maria.contreras'),
('Intereses', 'familia.contreras@email.com', 1, '💰', '#F7DC6F', 'maria.contreras'),

-- Familia Valdez (4 conceptos - 2 gastos, 2 ingresos)
('Seguro Médico', 'familia.valdez@email.com', 0, '🏥', '#D7BDE2', 'ana.valdez'),
('Clases', 'familia.valdez@email.com', 0, '🎓', '#A9DFBF', 'ana.valdez'),
('Consultoría', 'familia.valdez@email.com', 1, '💼', '#F5B7B1', 'ana.valdez'),
('Herencia', 'familia.valdez@email.com', 1, '📜', '#AED6F1', 'ana.valdez'),

-- Familia Rosales (4 conceptos - 2 gastos, 2 ingresos)
('Gasolina', 'familia.rosales@email.com', 0, '⛽', '#ABEBC6', 'roberto.rosales'),
('Alimentación', 'familia.rosales@email.com', 0, '🛒', '#FAD7A0', 'roberto.rosales'),
('Préstamo', 'familia.rosales@email.com', 1, '🏦', '#D2B4DE', 'roberto.rosales'),
('Regalías', 'familia.rosales@email.com', 1, '📖', '#AED6F1', 'roberto.rosales'),

-- Familia Soto (4 conceptos - 2 gastos, 2 ingresos)
('Vivienda', 'familia.soto@email.com', 0, '🏠', '#F9E79F', 'diego.soto'),
('Servicios', 'familia.soto@email.com', 0, '💡', '#D5DBDB', 'diego.soto'),
('Sueldo', 'familia.soto@email.com', 1, '💵', '#F1948A', 'diego.soto'),
('Inversiones', 'familia.soto@email.com', 1, '📈', '#82E0AA', 'diego.soto')







-- Insertar datos en la tabla Movimiento
INSERT INTO movimiento (idMovimiento, fecha, monto, descripcion, nombreUsuario, nombreConcepto, correoFamilia) VALUES
-- Movimientos Familia García
(1, '2024-10-01', 150.75, 'Compra semanal en supermercado', 'juan.garcia', 'Alimentación', 'familia.garcia@email.com'),
(2, '2024-10-02', 45.50, 'Gasolina para el coche', 'maria.garcia', 'Transporte', 'familia.garcia@email.com'),
(3, '2024-10-05', 2500.00, 'Salario mensual octubre', 'juan.garcia', 'Salario', 'familia.garcia@email.com'),
(4, '2024-10-08', 300.00, 'Inversión en fondos indexados', 'carlos.garcia', 'Inversiones', 'familia.garcia@email.com'),

-- Movimientos Familia Rodríguez
(5, '2024-10-03', 850.00, 'Pago de alquiler mensual', 'pedro.rodriguez', 'Vivienda', 'familia.rodriguez@email.com'),
(6, '2024-10-04', 120.50, 'Pago de luz y agua', 'ana.rodriguez', 'Servicios', 'familia.rodriguez@email.com'),
(7, '2024-10-10', 800.00, 'Pago por trabajo freelance', 'pedro.rodriguez', 'Freelance', 'familia.rodriguez@email.com'),
(8, '2024-10-15', 150.00, 'Bono por productividad', 'luis.rodriguez', 'Bonos', 'familia.rodriguez@email.com'),

-- Movimientos Familia López
(9, '2024-10-06', 200.00, 'Libros de texto universitarios', 'miguel.lopez', 'Educación', 'familia.lopez@email.com'),
(10, '2024-10-07', 85.00, 'Consulta médica especialista', 'elena.lopez', 'Salud', 'familia.lopez@email.com'),
(11, '2024-10-12', 450.00, 'Comisión por ventas', 'miguel.lopez', 'Comisiones', 'familia.lopez@email.com'),
(12, '2024-10-20', 100.00, 'Regalo de cumpleaños', 'javier.lopez', 'Regalos', 'familia.lopez@email.com'),

-- Movimientos Familia Martínez
(13, '2024-10-09', 75.00, 'Entradas de cine', 'david.martinez', 'Entretenimiento', 'familia.martinez@email.com'),
(14, '2024-10-11', 120.00, 'Compra de ropa de invierno', 'patricia.martinez', 'Ropa', 'familia.martinez@email.com'),
(15, '2024-10-14', 1200.00, 'Ingreso por alquiler propiedad', 'david.martinez', 'Alquiler Ingresos', 'familia.martinez@email.com'),
(16, '2024-10-18', 150.00, 'Dividendos de inversiones', 'jorge.martinez', 'Dividendos', 'familia.martinez@email.com'),

-- Movimientos Familia González
(17, '2024-10-13', 60.00, 'Membresía gimnasio mensual', 'sandra.gonzalez', 'Deportes', 'familia.gonzalez@email.com'),
(18, '2024-10-16', 500.00, 'Reserva para viaje a playa', 'ricardo.gonzalez', 'Viajes', 'familia.gonzalez@email.com'),
(19, '2024-10-19', 300.00, 'Venta de muebles usados', 'sandra.gonzalez', 'Ventas', 'familia.gonzalez@email.com'),
(20, '2024-10-22', 75.00, 'Intereses cuenta de ahorros', 'monica.gonzalez', 'Intereses', 'familia.gonzalez@email.com'),

-- Movimientos Familia Hernández
(21, '2024-10-17', 95.00, 'Cena familiar en restaurante', 'raquel.hernandez', 'Restaurantes', 'familia.hernandez@email.com'),
(22, '2024-10-21', 200.00, 'Pago seguro del auto', 'alejandro.hernandez', 'Seguros', 'familia.hernandez@email.com'),
(23, '2024-10-24', 600.00, 'Pago consultoría proyecto', 'raquel.hernandez', 'Consultoría', 'familia.hernandez@email.com'),
(24, '2024-10-25', 5000.00, 'Herencia recibida', 'lucia.hernandez', 'Herencia', 'familia.hernandez@email.com'),

-- Movimientos Familia Pérez
(25, '2024-10-23', 350.00, 'Pago impuestos trimestrales', 'sergio.perez', 'Impuestos', 'familia.perez@email.com'),
(26, '2024-10-26', 80.00, 'Visita al veterinario', 'elvira.perez', 'Mascotas', 'familia.perez@email.com'),
(27, '2024-10-28', 1000.00, 'Préstamo personal recibido', 'sergio.perez', 'Préstamo', 'familia.perez@email.com'),
(28, '2024-10-30', 250.00, 'Regalías por libro publicado', 'oscar.perez', 'Regalías', 'familia.perez@email.com'),

-- Movimientos Familia Sánchez
(29, '2024-10-27', 899.00, 'Compra de laptop nueva', 'humberto.sanchez', 'Tecnología', 'familia.sanchez@email.com'),
(30, '2024-10-29', 150.00, 'Compra de mueble para sala', 'veronica.sanchez', 'Hogar', 'familia.sanchez@email.com'),
(31, '2024-11-01', 3200.00, 'Sueldo quincenal', 'humberto.sanchez', 'Sueldo', 'familia.sanchez@email.com'),
(32, '2024-11-03', 500.00, 'Inversión en criptomonedas', 'arturo.sanchez', 'Inversiones', 'familia.sanchez@email.com'),

-- Movimientos Familia Ramírez
(33, '2024-11-02', 65.00, 'Corte de cabello y cuidado', 'gloria.ramirez', 'Cuidado Personal', 'familia.ramirez@email.com'),
(34, '2024-11-04', 180.25, 'Compra mensual de víveres', 'manuel.ramirez', 'Supermercado', 'familia.ramirez@email.com'),
(35, '2024-11-05', 2800.00, 'Salario mensual noviembre', 'gloria.ramirez', 'Salario', 'familia.ramirez@email.com'),
(36, '2024-11-08', 300.00, 'Bono por metas cumplidas', 'teresa.ramirez', 'Bonos', 'familia.ramirez@email.com'),

-- Movimientos Familia Torres
(37, '2024-11-06', 45.00, 'Entradas para película', 'silvia.torres', 'Cine', 'familia.torres@email.com'),
(38, '2024-11-07', 70.00, 'Pago mensual gimnasio', 'ramon.torres', 'Gimnasio', 'familia.torres@email.com'),
(39, '2024-11-09', 450.00, 'Pago proyecto freelance', 'silvia.torres', 'Freelance', 'familia.torres@email.com'),
(40, '2024-11-11', 120.00, 'Comisión por referido', 'beatriz.torres', 'Comisiones', 'familia.torres@email.com'),

-- Movimientos Familia Flores
(41, '2024-11-10', 35.50, 'Medicamentos recetados', 'francisco.flores', 'Farmacia', 'familia.flores@email.com'),
(42, '2024-11-12', 89.99, 'Libros de desarrollo personal', 'lourdes.flores', 'Libros', 'familia.flores@email.com'),
(43, '2024-11-14', 950.00, 'Ingreso por alquiler local', 'francisco.flores', 'Alquiler', 'familia.flores@email.com'),
(44, '2024-11-16', 180.00, 'Dividendos acciones', 'victor.flores', 'Dividendos', 'familia.flores@email.com'),

-- Movimientos Familia Rivera
(45, '2024-11-13', 120.00, 'Pago de factura eléctrica', 'eduardo.rivera', 'Electricidad', 'familia.rivera@email.com'),
(46, '2024-11-15', 45.00, 'Pago de servicio de agua', 'natalia.rivera', 'Agua', 'familia.rivera@email.com'),
(47, '2024-11-17', 750.00, 'Venta de arte online', 'eduardo.rivera', 'Ventas', 'familia.rivera@email.com'),
(48, '2024-11-19', 45.25, 'Intereses cuenta corriente', 'alfonso.rivera', 'Intereses', 'familia.rivera@email.com'),

-- Movimientos Familia Gómez
(49, '2024-11-18', 59.99, 'Pago mensual internet', 'ines.gomez', 'Internet', 'familia.gomez@email.com'),
(50, '2024-11-20', 35.00, 'Recarga telefónica', 'joaquin.gomez', 'Teléfono', 'familia.gomez@email.com', NULL);




-- Insertar datos en la tabla PersonalizacionConcepto
INSERT INTO personalizacionconcepto (idPersonalizacion, limiteGasto, activo, montoPlanificado, tipoPeriodoPlanificado, tipoPeriodoLimite, diaPeriodoPlanificado, notificacion, nombreUsuario, nombreConcepto, correoFamilia) VALUES
-- Personalizaciones Familia García
(1, 500.00, 1, 2000.00, 'mensual', 'mensual', 1, 1, 'juan.garcia', 'Alimentación', 'familia.garcia@email.com'),
(2, 300.00, 1, 150.00, 'mensual', 'mensual', 5, 0, 'juan.garcia', 'Transporte', 'familia.garcia@email.com'),

-- Personalizaciones Familia Rodríguez
(3, 1000.00, 1, 850.00, 'mensual', 'mensual', 3, 1, 'pedro.rodriguez', 'Vivienda', 'familia.rodriguez@email.com'),
(4, 200.00, 1, 120.00, 'mensual', 'mensual', 10, 0, 'pedro.rodriguez', 'Servicios', 'familia.rodriguez@email.com'),

-- Personalizaciones Familia López
(5, 300.00, 1, 250.00, 'quincenal', 'mensual', 15, 1, 'miguel.lopez', 'Educación', 'familia.lopez@email.com'),
(6, 150.00, 1, 100.00, 'mensual', 'mensual', 20, 0, 'miguel.lopez', 'Salud', 'familia.lopez@email.com'),

-- Personalizaciones Familia Martínez
(7, 200.00, 1, 100.00, 'mensual', 'mensual', 25, 1, 'david.martinez', 'Entretenimiento', 'familia.martinez@email.com'),
(8, 300.00, 1, 150.00, 'diario', 'mensual', 12, 0, 'david.martinez', 'Ropa', 'familia.martinez@email.com'),

-- Personalizaciones Familia González
(9, 100.00, 1, 60.00, 'mensual', 'mensual', 8, 1, 'sandra.gonzalez', 'Deportes', 'familia.gonzalez@email.com'),
(10, 800.00, 1, 500.00, 'anual', 'mensual', 1, 0, 'sandra.gonzalez', 'Viajes', 'familia.gonzalez@email.com'),

-- Personalizaciones Familia Hernández
(11, 150.00, 1, 100.00, 'mensual', 'mensual', 18, 1, 'raquel.hernandez', 'Restaurantes', 'familia.hernandez@email.com'),
(12, 250.00, 1, 200.00, 'semestral', 'mensual', 5, 0, 'raquel.hernandez', 'Seguros', 'familia.hernandez@email.com'),

-- Personalizaciones Familia Pérez
(13, 400.00, 1, 350.00, 'quincenal', 'mensual', 28, 1, 'sergio.perez', 'Impuestos', 'familia.perez@email.com'),
(14, 100.00, 1, 80.00, 'mensual', 'mensual', 15, 0, 'sergio.perez', 'Mascotas', 'familia.perez@email.com'),

-- Personalizaciones Familia Sánchez
(15, 1000.00, 1, 500.00, 'anual', 'mensual', 10, 1, 'humberto.sanchez', 'Tecnología', 'familia.sanchez@email.com'),
(16, 200.00, 1, 150.00, 'mensual', 'mensual', 22, 0, 'humberto.sanchez', 'Hogar', 'familia.sanchez@email.com'),

-- Personalizaciones Familia Ramírez
(17, 100.00, 1, 65.00, 'mensual', 'mensual', 5, 1, 'gloria.ramirez', 'Cuidado Personal', 'familia.ramirez@email.com'),
(18, 250.00, 1, 200.00, 'mensual', 'mensual', 1, 0, 'gloria.ramirez', 'Supermercado', 'familia.ramirez@email.com'),

-- Personalizaciones Familia Torres
(19, 100.00, 1, 50.00, 'mensual', 'mensual', 12, 1, 'silvia.torres', 'Cine', 'familia.torres@email.com'),
(20, 100.00, 1, 70.00, 'mensual', 'mensual', 7, 0, 'silvia.torres', 'Gimnasio', 'familia.torres@email.com'),

-- Personalizaciones Familia Flores
(21, 50.00, 1, 35.00, 'mensual', 'mensual', 20, 1, 'francisco.flores', 'Farmacia', 'familia.flores@email.com'),
(22, 100.00, 1, 90.00, 'mensual', 'mensual', 25, 0, 'francisco.flores', 'Libros', 'familia.flores@email.com'),

-- Personalizaciones Familia Rivera
(23, 150.00, 1, 120.00, 'mensual', 'mensual', 15, 1, 'eduardo.rivera', 'Electricidad', 'familia.rivera@email.com'),
(24, 60.00, 1, 45.00, 'mensual', 'mensual', 18, 0, 'eduardo.rivera', 'Agua', 'familia.rivera@email.com'),

-- Personalizaciones Familia Gómez
(25, 70.00, 1, 60.00, 'mensual', 'mensual', 5, 1, 'ines.gomez', 'Internet', 'familia.gomez@email.com'),
(26, 50.00, 1, 35.00, 'mensual', 'mensual', 10, 0, 'ines.gomez', 'Teléfono', 'familia.gomez@email.com'),

-- Personalizaciones Familia Díaz
(27, 80.00, 1, 45.00, 'mensual', 'mensual', 12, 1, 'celia.diaz', 'Gas', 'familia.diaz@email.com'),
(28, 120.00, 1, 100.00, 'mensual', 'mensual', 8, 0, 'celia.diaz', 'Transporte Público', 'familia.diaz@email.com'),

-- Personalizaciones Familia Reyes
(29, 150.00, 1, 80.00, 'mensual', 'mensual', 20, 1, 'antonio.reyes', 'Taxi', 'familia.reyes@email.com'),
(30, 200.00, 1, 150.00, 'quincenal', 'mensual', 15, 0, 'antonio.reyes', 'Mantenimiento', 'familia.reyes@email.com'),

-- Personalizaciones Familia Morales
(31, 300.00, 1, 250.00, 'mensual', 'mensual', 1, 1, 'pablo.morales', 'Seguro Médico', 'familia.morales@email.com'),
(32, 200.00, 1, 150.00, 'mensual', 'mensual', 10, 0, 'pablo.morales', 'Clases', 'familia.morales@email.com'),

-- Personalizaciones Familia Ortiz
(33, 200.00, 1, 150.00, 'mensual', 'mensual', 5, 1, 'lucia.ortiz', 'Gasolina', 'familia.ortiz@email.com'),
(34, 600.00, 1, 500.00, 'mensual', 'mensual', 1, 0, 'lucia.ortiz', 'Alimentación', 'familia.ortiz@email.com'),

-- Personalizaciones Familia Silva
(35, 900.00, 1, 850.00, 'mensual', 'mensual', 3, 1, 'andrea.silva', 'Vivienda', 'familia.silva@email.com'),
(36, 250.00, 1, 200.00, 'mensual', 'mensual', 12, 0, 'andrea.silva', 'Servicios', 'familia.silva@email.com'),

-- Personalizaciones Familia Vargas
(37, 400.00, 1, 300.00, 'quincenal', 'mensual', 20, 1, 'ricardo.vargas', 'Educación', 'familia.vargas@email.com'),
(38, 200.00, 1, 150.00, 'mensual', 'mensual', 15, 0, 'ricardo.vargas', 'Salud', 'familia.vargas@email.com'),

-- Personalizaciones Familia Castro
(39, 250.00, 1, 100.00, 'mensual', 'mensual', 25, 1, 'gabriel.castro', 'Entretenimiento', 'familia.castro@email.com'),
(40, 400.00, 1, 200.00, 'diario', 'mensual', 10, 0, 'gabriel.castro', 'Ropa', 'familia.castro@email.com'),

-- Personalizaciones Familia Romero
(41, 150.00, 1, 80.00, 'mensual', 'mensual', 8, 1, 'alejandra.romero', 'Deportes', 'familia.romero@email.com'),
(42, 1200.00, 1, 800.00, 'anual', 'mensual', 1, 0, 'alejandra.romero', 'Viajes', 'familia.romero@email.com'),

-- Personalizaciones Familia Aguilar
(43, 200.00, 1, 120.00, 'mensual', 'mensual', 18, 1, 'raul.aguilar', 'Restaurantes', 'familia.aguilar@email.com'),
(44, 300.00, 1, 250.00, 'semestral', 'mensual', 5, 0, 'raul.aguilar', 'Seguros', 'familia.aguilar@email.com'),

-- Personalizaciones Familia Mendoza
(45, 500.00, 1, 400.00, 'quincenal', 'mensual', 28, 1, 'clara.mendoza', 'Impuestos', 'familia.mendoza@email.com'),
(46, 120.00, 1, 90.00, 'mensual', 'mensual', 15, 0, 'clara.mendoza', 'Mascotas', 'familia.mendoza@email.com'),

-- Personalizaciones Familia Herrera
(47, 800.00, 1, 600.00, 'anual', 'mensual', 10, 1, 'julio.herrera', 'Tecnología', 'familia.herrera@email.com'),
(48, 250.00, 1, 180.00, 'mensual', 'mensual', 22, 0, 'julio.herrera', 'Hogar', 'familia.herrera@email.com'),

-- Personalizaciones Familia Guzmán
(49, 80.00, 1, 50.00, 'mensual', 'mensual', 5, 1, 'natalia.guzman', 'Cuidado Personal', 'familia.guzman@email.com'),
(50, 300.00, 1, 250.00, 'mensual', 'mensual', 1, 0, 'natalia.guzman', 'Supermercado', 'familia.guzman@email.com');


UPDATE usuario 
SET nombreUsuario = REPLACE(nombreUsuario, '.', '')
WHERE nombreUsuario LIKE '%.%';


-- Tabla concepto
UPDATE concepto 
SET nombreUsuario = REPLACE(nombreUsuario, '.', '')
WHERE nombreUsuario LIKE '%.%';

-- Tabla movimiento
UPDATE movimiento 
SET nombreUsuario = REPLACE(nombreUsuario, '.', '')
WHERE nombreUsuario LIKE '%.%';

-- Tabla personalizacionconcepto
UPDATE personalizacionconcepto 
SET nombreUsuario = REPLACE(nombreUsuario, '.', '')
WHERE nombreUsuario LIKE '%.%';