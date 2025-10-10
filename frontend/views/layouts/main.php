<!DOCTYPE html>
<html lang="es">

<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title><?php echo htmlspecialchars($title ?? 'Arka App'); ?></title>

  <!-- ✅ CSS Global (siempre se carga) -->
  <link rel="stylesheet" href="<?php echo URLROOT; ?>/frontend/assets/css/global.css">

  <!-- ✅ CSS Específico de la página (si existe el archivo) -->
  <?php if (isset($pageCSS)): ?>
    <?php
    $cssPath = ROOT . "/frontend/assets/css/{$pageCSS}.css";
    if (file_exists($cssPath)):
    ?>
      <link rel="stylesheet" href="<?php echo URLROOT; ?>/frontend/assets/css/<?php echo $pageCSS; ?>.css">
    <?php endif; ?>
  <?php endif; ?>

  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.0/css/all.min.css">


</head>

<body>
  <!-- Sidebar -->
  <?php require_once ROOT . '/frontend/views/partials/sidebar.php'; ?>

  <!-- Contenido Principal -->
  <main class="main">
    <?php require_once ROOT . '/frontend/views/partials/user-card.php'; ?>
    <?php echo $content; ?>
  </main>




  <!-- ✅ JS Global (siempre se carga) -->
  <script src="<?php echo URLROOT; ?>/frontend/assets/js/main.js"></script>

  <!-- ✅ JS Específico de la página (si existe el archivo) -->
  <?php if (isset($pageJS)): ?>
    <?php
    $jsPath = ROOT . "/frontend/assets/js/{$pageJS}.js";
    if (file_exists($jsPath)):
    ?>
      <script src="<?php echo URLROOT; ?>/frontend/assets/js/<?php echo $pageJS; ?>.js"></script>
    <?php endif; ?>
  <?php endif; ?>
</body>

</html>