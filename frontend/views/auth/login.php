<?php if (!empty($success)): ?>
    <p style="color: green; font-weight: bold;"><?php echo htmlspecialchars($success); ?></p>
<?php endif; ?>
<!DOCTYPE html>
<html lang="es">
<head>
  <meta charset="UTF-8">
  <title><?php echo htmlspecialchars($title ?? 'Login - Arka'); ?></title>
  <link rel="stylesheet" href="<?php echo URLROOT; ?>/frontend/assets/css/auth.css">  <!-- Corrige path con URLROOT -->
</head>
<body>
  <form method="POST" action="<?php echo URLROOT; ?>/login">  <!-- Action explícito para POST seguro -->
      <img src="<?php echo URLROOT; ?>/frontend/assets/img/arka_logo.png" alt="Logo Arka" class="logo-top">  <!-- Corrige path -->
      <h2>Iniciar sesión</h2>
      <?php if (!empty($error)) echo "<p style='color:red'>" . htmlspecialchars($error) . "</p>"; ?>
      
      <!-- Errores de validador -->
      <?php if (!empty($validation_errors)): ?>
          <ul class="errors">
              <?php foreach ($validation_errors as $field => $msg): ?>
                  <li style="color: red;"><?php echo htmlspecialchars($msg); ?></li>
              <?php endforeach; ?>
          </ul>
      <?php endif; ?>
      
      <input type="email" name="email" placeholder="Email" required>  <!-- Cambiado a 'email' -->
      <input type="password" name="contrasena" placeholder="Contraseña" required>  <!-- Cambiado a 'contrasena' -->
      
      <a href="<?php echo URLROOT; ?>/recuperar_cuenta.php" class="link-reset">¿Olvidaste tu contraseña?</a>  <!-- URLROOT para link -->
      <button type="submit" class="btn btn-oscuro">Entrar</button>
      <p>¿No tienes una cuenta? <a href="<?php echo URLROOT; ?>/registro" class="btn-create">Crea una nueva aquí!</a></p>  <!-- URLROOT -->
      <div class="logo-bottom-container">
          <p>Desarrollado por:</p>
          <img src="<?php echo URLROOT; ?>/frontend/assets/img/bugbusters.png" alt="Logo Empresa" class="logo-bottom">  <!-- Corrige path -->
      </div>
  </form>
</body>
</html>