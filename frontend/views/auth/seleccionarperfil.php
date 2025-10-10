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

  <div class="perfil-container">
    <?php if (!empty($profiles)): ?>
      <?php foreach ($profiles as $profile): ?>

        <?php
        // Color básico por rol (expande si tienes icons/colores en DB)
        $color = ($profile['rol'] === 'admin') ? '#007bff' : '#6c757d';
        $icon = ($profile['rol'] === 'admin') ? URLROOT . '/frontend/assets/img/icono_usuario1.png' : URLROOT . '/frontend/assets/img/icono_usuario2.png';
        ?>
        <a class="perfil-card" href="<?php echo URLROOT; ?>/seleccionar-perfil?profile=<?php echo urlencode($profile['id']); ?>" style="background-color: <?php echo $color; ?>;">
          <img class="perfil-icon" src="<?php echo $icon; ?>" alt="Icono de <?php echo htmlspecialchars($profile['nombre']); ?>">
          <div class="perfil-name"><?php echo htmlspecialchars($profile['nombre']); ?></div>
        </a>

        <?php if ($profile['protegido']) { ?>
          <div class="perfil-card <?php echo $profile['protegido'] ? 'protegido' : ''; ?>">
            <span class="candado">🔒</span>
            <button type="button" onclick="miModal.showModal()">Contraseña</button>

            <dialog id="miModal">
              <h2>Protección con contraseña</h2>
              <form method="post" action="<?php echo URLROOT; ?>/seleccionar-perfil">
                <input type="password" name="password" placeholder="Ingresa tu contraseña">
                <input type="hidden" name="id" value="<?php echo $profile['id']; ?>">
                <button type="submit">Aceptar</button>
                <form method="dialog">
                  <button type="button" onclick="miModal.close()">Cancelar</button>
                </form>
              </form>
            </dialog>
          </div>
        <?php } ?>
      <?php endforeach; ?>
    <?php else: ?>
      <p>No hay perfiles disponibles.</p>
    <?php endif; ?>
  </div>
</body>

</html>