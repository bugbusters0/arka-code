<?php
class AuthMiddleware
{

  /**
   * Verifica si el usuario NO está autenticado
   * (para páginas como login/registro)
   */
  public static function guest()
  {
    // Si tiene user_id (familia autenticada) pero NO tiene miembro_id
    if (isset($_SESSION['user_id']) && !isset($_SESSION['miembro_id'])) {
      // Ya inició sesión de familia, redirigir a selección de perfil
      header('Location: ' . URLROOT . '/seleccionar-perfil');
      exit;
    }

    // Si tiene ambos (completamente autenticado)
    if (isset($_SESSION['user_id']) && isset($_SESSION['miembro_id'])) {
      // Ya está completamente autenticado, redirigir al dashboard
      header('Location: ' . URLROOT . '/concepto/listar');
      exit;
    }
  }

  /**
   * Verifica si el usuario está autenticado
   * (para páginas protegidas como dashboard)
   */
  public static function auth()
  {
    if (!isset($_SESSION['user_id']) || !isset($_SESSION['miembro_id'])) {
      // Usuario no autenticado, redirigir al login
      header('Location: ' . URLROOT . '/login');
      exit;
    }
  }

  /**
   * Verifica solo la familia (para selección de perfil)
   */
  public static function familyOnly()
  {
    // Si no tiene user_id, redirigir a login
    if (!isset($_SESSION['user_id'])) {
      header('Location: ' . URLROOT . '/login');
      exit;
    }

    // Si ya tiene miembro_id, redirigir al dashboard (ya seleccionó perfil)
    if (isset($_SESSION['miembro_id'])) {
      header('Location: ' . URLROOT . '/concepto/listar');
      exit;
    }
  }
}
