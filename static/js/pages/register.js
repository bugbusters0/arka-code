// En static/js/pages/register.js
function togglePassword(id, btn) {
  const input = document.getElementById(id);
  const isHidden = input.type === "password";
  input.type = isHidden ? "text" : "password";
  btn.textContent = isHidden ? "🙈" : "👁";
}

const addBtn = document.getElementById('addMiembro');
const container = document.getElementById('miembrosContainer');
let memberCount = 0;

addBtn.addEventListener('click', () => {
  memberCount++;
  const miembro = document.createElement('div');
  miembro.classList.add('registro-miembro');
  miembro.innerHTML = `
    <div class="registro-grid">
      <div class="registro-group">
        <label>Nombre completo</label>
        <input type="text" name="miembroNombre[]" pattern="[A-Za-zÁÉÍÓÚáéíóúÑñ ]+" placeholder="Nombre completo" required>
      </div>
      <div class="registro-group">
        <label>Nombre de usuario</label>
        <input type="text" name="miembroUsuario[]" placeholder="nombre.usuario" required>
      </div>
    </div>
    <div class="registro-group">
      <label>Contraseña personal</label>
      <div class="registro-password">
        <input type="password" name="miembroContra[]" placeholder="PIN 4-6 dígitos" required>
        <button type="button" class="registro-toggle" onclick="toggleMemberPassword(this)">👁</button>
      </div>
    </div>
    <div class="registro-group">
      <label>Rol</label>
      <select name="miembroRol[]">
        <option value="miembro">Miembro</option>
        <option value="admin">Administrador</option>
      </select>
    </div>
    <button type="button" class="registro-remove" onclick="this.parentElement.remove(); memberCount--;">🗑 Eliminar</button>
    <hr>
  `;
  container.appendChild(miembro);
});

// Función para toggle password en miembros dinámicos
function toggleMemberPassword(btn) {
  const input = btn.previousElementSibling;
  const isHidden = input.type === "password";
  input.type = isHidden ? "text" : "password";
  btn.textContent = isHidden ? "🙈" : "👁";
}