<!DOCTYPE html>
<html lang="es">
<head>
  <meta charset="UTF-8">
  <title><?php echo htmlspecialchars($title ?? 'Registro - Arka'); ?></title>
  <link rel="stylesheet" href="<?php echo URLROOT; ?>/frontend/assets/css/auth.css">
</head>
<body class="registro-page">
  <div class="registro-container">
    <div class="registro-logo-column">
      <img src="<?php echo URLROOT; ?>/frontend/assets/img/arka_logo.png" alt="Logo superior Arka" class="registro-logo1">
      <img src="<?php echo URLROOT; ?>/frontend/assets/img/bugbusters.png" alt="Logo inferior Arka" class="registro-logo2">
    </div>
    
    <!-- ✅ CAMBIA: id="registroForm" → method="POST" action -->
    <form method="POST" action="<?php echo URLROOT; ?>/registro" class="registro-form">
      <h2>Crear una nueva cuenta</h2>

      <!-- ✅ Muestra errores de validación -->
      <?php if (!empty($validation_errors)): ?>
          <div class="error-box">
              <ul>
                  <?php foreach ($validation_errors as $field => $msg): ?>
                      <li style="color: red;"><?php echo htmlspecialchars($msg); ?></li>
                  <?php endforeach; ?>
              </ul>
          </div>
      <?php endif; ?>

      <?php if (!empty($error)): ?>
          <p style="color:red"><?php echo htmlspecialchars($error); ?></p>
      <?php endif; ?>

      <!-- DATOS GENERALES -->
      <div class="registro-grid">
        <div class="registro-group">
          <label for="email">Correo electrónico</label>
          <input type="email" id="email" name="email" placeholder="ejemplo@correo.com" required>
        </div>
        <div class="registro-group">
          <label for="telefono">Teléfono</label>
          <div class="registro-tel">
            <span class="registro-prefijo">+51</span>
            <input type="tel" id="telefono" name="telefono" placeholder="987654321" required>
          </div>
        </div>
      </div>

      <!-- CREDENCIALES -->
      <div class="registro-grid">
        <div class="registro-group">
          <label for="password">Contraseña</label>
          <div class="registro-password">
            <input type="password" id="password" name="password" placeholder="********" required>
            <button type="button" class="registro-toggle" onclick="togglePassword('password', this)">👁</button>
          </div>
        </div>
        <div class="registro-group">
          <label for="confirmPassword">Confirmar contraseña</label>
          <div class="registro-password">
            <input type="password" id="confirmPassword" name="confirmPassword" placeholder="********" required>
            <button type="button" class="registro-toggle" onclick="togglePassword('confirmPassword', this)">👁</button>
          </div>
        </div>
      </div>

      <!-- ADMINISTRADOR -->
      <div class="registro-grid">
        <div class="registro-group">
          <label for="adminNombre">Nombre del administrador</label>
          <input type="text" id="adminNombre" name="adminNombre" pattern="[A-Za-zÁÉÍÓÚáéíóúÑñ ]+" placeholder="Nombre completo" required>
        </div>
        <div class="registro-group">
          <label for="adminNacimiento">Fecha de nacimiento</label>
          <input type="date" id="adminNacimiento" name="adminNacimiento" required>
        </div>
      </div>

      <!-- ✅ CONTRASEÑA PERSONAL DEL ADMIN (requerida por BD) -->
      <div class="registro-group">
        <label for="contraPersonal">Contraseña personal del administrador</label>
        <div class="registro-password">
          <input type="password" id="contraPersonal" name="contraPersonal" placeholder="PIN de 4-6 dígitos" required>
          <button type="button" class="registro-toggle" onclick="togglePassword('contraPersonal', this)">👁</button>
        </div>
      </div>

      <!-- MIEMBROS -->
      <div class="registro-miembros">
        <h3>Miembros de la familia (opcional)</h3>
        <div id="miembrosContainer" class="registro-miembros-container"></div>
        <button type="button" id="addMiembro" class="registro-add-member">+ Añadir miembro</button>
      </div>

      <div class="registro-submit">
        <button type="submit" class="registro-btn">Crear Cuenta</button>
      </div>
      
      <p>¿Ya tienes cuenta? <a href="<?php echo URLROOT; ?>/login">Inicia sesión aquí</a></p>
    </form>
  </div>

  <script>
    function togglePassword(id, btn) {
      const input = document.getElementById(id);
      const isHidden = input.type === "password";
      input.type = isHidden ? "text" : "password";
      btn.textContent = isHidden ? "🙈" : "👁";
    }

    const addBtn = document.getElementById('addMiembro');
    const container = document.getElementById('miembrosContainer');
    addBtn.addEventListener('click', () => {
      const miembro = document.createElement('div');
      miembro.classList.add('registro-miembro');
      miembro.innerHTML = `
        <div class="registro-grid">
          <div class="registro-group">
            <label>Nombre</label>
            <input type="text" name="miembroNombre[]" pattern="[A-Za-zÁÉÍÓÚáéíóúÑñ ]+" placeholder="Nombre completo">
          </div>
          <div class="registro-group">
            <label>Fecha de nacimiento</label>
            <input type="date" name="miembroNacimiento[]">
          </div>
        </div>
        <div class="registro-group">
          <label>Contraseña personal</label>
          <input type="password" name="miembroContra[]" placeholder="PIN 4-6 dígitos">
        </div>
        <div class="registro-group">
          <label>Rol</label>
          <select name="miembroRol[]">
            <option value="admin">Administrador</option>
            <option value="miembro">Miembro</option>
          </select>
        </div>
        <button type="button" class="registro-remove" onclick="this.parentElement.remove()">🗑 Eliminar</button>
      `;
      container.appendChild(miembro);
    });
  </script>
</body>
</html>