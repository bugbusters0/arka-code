-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Servidor: 127.0.0.1
-- Tiempo de generación: 24-10-2025 a las 11:11:10
-- Versión del servidor: 10.4.32-MariaDB
-- Versión de PHP: 8.2.12

CREATE DATABASE IF NOT EXISTS arka;
USE arka;

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



