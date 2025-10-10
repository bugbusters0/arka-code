<!DOCTYPE html>
<html lang="es">

<head>
  <meta charset="UTF-8">
  <title><?php echo htmlspecialchars($title ?? 'Seleccionar Perfil - Arka'); ?></title>
  <link rel="stylesheet" href="<?php echo URLROOT; ?>/frontend/assets/css/auth.css">
</head>

<body class="perfil-page">
  <h1>¿Quién está usando Arka?</h1>
  <?php if (!empty($error)): ?>
    <p style="color: red;"><?php echo htmlspecialchars($error); ?></p>
  <?php endif; ?>
  <?php if (!empty($_SESSION['error_message'])): ?>
    <p style="color: red;"><?php echo htmlspecialchars($_SESSION['error_message']); ?></p>
  <?php endif; ?>

  <div class="perfil-container">
    <?php if (!empty($profiles)): ?>
      <?php foreach ($profiles as $index => $profile): ?>
        <?php
        $color = ($profile['rol'] === 'admin') ? '#007bff' : '#6c757d';
        $icon = ($profile['rol'] === 'admin') ? URLROOT . '/frontend/assets/img/icono_usuario1.png' : URLROOT . '/frontend/assets/img/icono_usuario2.png';
        $modalId = 'miModal_' . $profile['id'];
        ?>

        <!-- Perfil normal (sin protección) -->
        <?php if (!$profile['protegido']): ?>
          <a class="perfil-card" href="<?php echo URLROOT; ?>/seleccionar-perfil?profile=<?php echo urlencode($profile['id']); ?>" style="background-color: <?php echo $color; ?>;">
            <img class="perfil-icon" src="<?php echo $icon; ?>" alt="Icono de <?php echo htmlspecialchars($profile['nombre']); ?>">
            <div class="perfil-name"><?php echo htmlspecialchars($profile['nombre']); ?></div>
          </a>
        <?php else: ?>
          <!-- Perfil protegido (con modal) -->
          <div class="perfil-card protegido" style="background-color: <?php echo $color; ?>;">
            <span class="candado">🔒</span>
            <button type="button" onclick="abrirModal('<?php echo $modalId; ?>')">
              <img class="perfil-icon" src="<?php echo $icon; ?>" alt="Icono de <?php echo htmlspecialchars($profile['nombre']); ?>">
              <div class="perfil-name"><?php echo htmlspecialchars($profile['nombre']); ?></div>
            </button>

            <dialog id="<?php echo $modalId; ?>">
              <h2>Protección con contraseña</h2>
              <p>Ingresa la contraseña para <strong><?php echo htmlspecialchars($profile['nombre']); ?></strong></p>
              <form method="post" action="<?php echo URLROOT; ?>/seleccionar-perfil">
                <input type="password" name="password" placeholder="Ingresa tu contraseña" required>
                <input type="hidden" name="idMiembro" value="<?php echo $profile['id']; ?>">
                <div class="modal-buttons">
                  <button type="submit">Aceptar</button>
                  <button type="button" onclick="cerrarModal('<?php echo $modalId; ?>')">Cancelar</button>
                </div>
              </form>
            </dialog>
          </div>
        <?php endif; ?>
      <?php endforeach; ?>
    <?php else: ?>
      <p>No hay perfiles disponibles.</p>
    <?php endif; ?>
  </div>

  <script>
    function abrirModal(modalId) {
      const modal = document.getElementById(modalId);
      if (modal) {
        modal.showModal();
      }
    }

    function cerrarModal(modalId) {
      const modal = document.getElementById(modalId);
      if (modal) {
        modal.close();
      }
    }

    document.addEventListener('DOMContentLoaded', function() {
      document.querySelectorAll('dialog').forEach(dialog => {
        dialog.addEventListener('click', function(event) {
          const rect = dialog.getBoundingClientRect();
          if (event.clientX < rect.left || event.clientX > rect.right ||
            event.clientY < rect.top || event.clientY > rect.bottom) {
            dialog.close();
          }
        });
      });
    });
  </script>
</body>

</html>