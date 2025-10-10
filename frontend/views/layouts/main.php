<!DOCTYPE html>
<html lang="es">

<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title><?php echo htmlspecialchars($title ?? 'Arka App'); ?></title>

  <!-- CSS Global -->
  <link rel="stylesheet" href="<?php echo URLROOT; ?>/frontend/assets/css/main.css">
  <link rel="stylesheet" href="<?php echo URLROOT; ?>/frontend/assets/css/components.css">

  <!-- CSS Específico de la página (si existe) -->
  <?php if (isset($pageCSS)): ?>
    <link rel="stylesheet" href="<?php echo URLROOT; ?>/frontend/assets/css/<?php echo $pageCSS; ?>.css">
  <?php endif; ?>

  <style>
    * {
      margin: 0;
      padding: 0;
      box-sizing: border-box;
    }

    body {
      font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
      background: #f5f7fa;
      min-height: 100vh;
      display: flex;
      flex-direction: column;
    }

    .navbar {
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      color: white;
      padding: 1rem 2rem;
      box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
    }

    .navbar-content {
      max-width: 1200px;
      margin: 0 auto;
      display: flex;
      justify-content: space-between;
      align-items: center;
    }

    .navbar-brand {
      font-size: 1.5rem;
      font-weight: bold;
      text-decoration: none;
      color: white;
    }

    .navbar-menu {
      display: flex;
      gap: 2rem;
      list-style: none;
    }

    .navbar-menu a {
      color: white;
      text-decoration: none;
      padding: 0.5rem 1rem;
      border-radius: 5px;
      transition: background 0.3s;
    }

    .navbar-menu a:hover {
      background: rgba(255, 255, 255, 0.2);
    }

    .navbar-user {
      display: flex;
      align-items: center;
      gap: 1rem;
    }

    .main-content {
      flex: 1;
      max-width: 1200px;
      width: 100%;
      margin: 2rem auto;
      padding: 0 2rem;
    }

    .footer {
      background: #2c3e50;
      color: white;
      text-align: center;
      padding: 2rem;
      margin-top: auto;
    }

    .btn {
      padding: 0.5rem 1.5rem;
      border: none;
      border-radius: 5px;
      cursor: pointer;
      text-decoration: none;
      display: inline-block;
      transition: all 0.3s;
    }

    .btn-primary {
      background: #667eea;
      color: white;
    }

    .btn-primary:hover {
      background: #5568d3;
      transform: translateY(-2px);
    }

    .btn-danger {
      background: #e74c3c;
      color: white;
    }

    .btn-danger:hover {
      background: #c0392b;
    }
  </style>
</head>

<body>
  <!-- Navbar -->
  <nav class="navbar">
    <div class="navbar-content">
      <a href="<?php echo URLROOT; ?>/dashboard" class="navbar-brand">
        <img src="<?php echo URLROOT; ?>/frontend/assets/img/arka_logo.png"
          alt="Arka" style="height: 30px; vertical-align: middle;">
        Arka
      </a>

      <ul class="navbar-menu">
        <li><a href="<?php echo URLROOT; ?>/dashboard">Dashboard</a></li>
        <li><a href="<?php echo URLROOT; ?>/concepto/listar">Conceptos</a></li>
        <li><a href="<?php echo URLROOT; ?>/movimientos">Movimientos</a></li>
        <li><a href="<?php echo URLROOT; ?>/reportes">Reportes</a></li>
      </ul>

      <div class="navbar-user">
        <span>👤 <?php echo htmlspecialchars($_SESSION['user_email'] ?? 'Usuario'); ?></span>
        <a href="<?php echo URLROOT; ?>/logout" class="btn btn-danger">Cerrar Sesión</a>
      </div>
    </div>
  </nav>

  <!-- Contenido Principal -->
  <main class="main-content">
    <!-- Mensajes Flash -->
    <?php if (isset($_SESSION['success_message'])): ?>
      <div style="background: #2ecc71; color: white; padding: 1rem; border-radius: 5px; margin-bottom: 1rem;">
        ✓ <?php echo htmlspecialchars($_SESSION['success_message']); ?>
      </div>
      <?php unset($_SESSION['success_message']); ?>
    <?php endif; ?>

    <?php if (isset($_SESSION['error_message'])): ?>
      <div style="background: #e74c3c; color: white; padding: 1rem; border-radius: 5px; margin-bottom: 1rem;">
        ✕ <?php echo htmlspecialchars($_SESSION['error_message']); ?>
      </div>
      <?php unset($_SESSION['error_message']); ?>
    <?php endif; ?>

    <!-- Aquí se inyecta el contenido específico -->
    <?php echo $content; ?>
  </main>

  <!-- Footer -->
  <footer class="footer">
    <p>&copy; <?php echo date('Y'); ?> Arka App - Desarrollado por BugBusters</p>
  </footer>

  <!-- JS Global -->
  <script src="<?php echo URLROOT; ?>/frontend/assets/js/main.js"></script>

  <!-- JS Específico de la página (si existe) -->
  <?php if (isset($pageJS)): ?>
    <script src="<?php echo URLROOT; ?>/frontend/assets/js/<?php echo $pageJS; ?>.js"></script>
  <?php endif; ?>
</body>

</html>