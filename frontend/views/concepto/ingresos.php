<div class="tabs">
  <a href="<?php echo URLROOT . "/concepto/gasto" ?>" class="tab active">Gastos</a>
  <a href="<?php echo URLROOT . "/concepto/ingresos" ?>" class="tab">Ingresos</a>
</div>


<div class="content-area">
  <!-- Lista de conceptos -->
  <div class="conceptos" id="listaConceptos2">

    <?php foreach ($conceptos as $conceptoData): ?>
      <?php
      $concepto = $conceptoData['concepto'];
      $conceptoUsuario = $conceptoData['concepto_usuario'];
      ?>

      <div class="concepto" data-id="<?php echo $concepto->getId(); ?>">
        <div class="concepto-icono" style="color: <?php echo $concepto->getColor(); ?>;">
          <i class="<?php echo $conceptoData['icono_path']; ?>"></i>
        </div>
        <div class="concepto-info">
          <div class="concepto-nombre"><?php echo htmlspecialchars($concepto->getNombre()); ?></div>
          <?php if ($conceptoUsuario->getDesembolsoPlanejado() > 0): ?>
            <div class="concepto-desembolso">
              Desembolso: S/ <?php echo number_format($conceptoUsuario->getDesembolsoPlanejado(), 2); ?>
              <?php if ($conceptoUsuario->getPeriodoTipo()): ?>
                / <?php echo ucfirst($conceptoUsuario->getPeriodoTipo()); ?>
              <?php endif; ?>
            </div>
          <?php endif; ?>
          <?php if ($conceptoUsuario->getLimiteMonto() > 0): ?>
            <div class="concepto-limite">
              Límite: S/ <?php echo number_format($conceptoUsuario->getLimiteMonto(), 2); ?>
              <?php if ($conceptoUsuario->getLimiteTipo()): ?>
                / <?php echo ucfirst($conceptoUsuario->getLimiteTipo()); ?>
              <?php endif; ?>
            </div>
          <?php endif; ?>
        </div>
        <div class="concepto-acciones">
          <button class="btn-editar" onclick="editarConcepto(<?php echo $concepto->getId(); ?>)">✏️</button>
          <button class="btn-deshabilitar" onclick="deshabilitarConcepto(<?php echo $concepto->getId(); ?>)">🚫</button>
        </div>
      </div>
    <?php endforeach; ?>


    <!-- Botón "Nuevo" -->
    <div class="concepto nuevo" onclick="mostrarFormulario()">
      <i class="fa-solid fa-plus"></i>
      <span>Nuevo Concepto</span>
    </div>
  </div>

  <!-- Formulario lateral -->
  <form method="POST" action="<?php echo URLROOT; ?>/concepto/crear">
    <div class="sidebar-right" id="formConcepto">
      <h3 id="tituloForm">Crear Concepto</h3>

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

      <?php if (!empty($success)): ?>
        <div class="success-box">
          <p style="color:green"><?php echo htmlspecialchars($success); ?></p>
        </div>
      <?php endif; ?>

      <input type="hidden" name="tipo" value="gasto">
      <input type="hidden" name="configuracion_usuario" value="true">

      <div class="form-group">
        <label>Nombre del Concepto *</label>
        <input type="text" id="nombre" name="nombre" placeholder="Ej: Transporte" required
          value="<?php echo htmlspecialchars($_POST['nombre'] ?? ''); ?>">
      </div>

      <div class="form-group">
        <label>Ícono *</label>
        <div id="iconosGrid" class="iconos-grid">
          <?php foreach ($icons as $icon): ?>
            <input type="radio" name="id_icono" id="icono_<?php echo $icon->getId(); ?>"
              value="<?php echo $icon->getId(); ?>" class="icono-radio" required
              <?php echo (isset($_POST['id_icono']) && $_POST['id_icono'] == $icon->getId()) ? 'checked' : ''; ?>>
            <label for="icono_<?php echo $icon->getId(); ?>" class="icono-label">
              <i class="<?php echo htmlspecialchars($icon->getPath()); ?>"></i>
              <span><?php echo htmlspecialchars($icon->getNombre()); ?></span>
            </label>
          <?php endforeach; ?>
        </div>
      </div>

      <div class="form-group">
        <label>Color *</label>
        <div class="color-picker">
          <input type="radio" name="color" id="color1" value="#FF6B6B"
            <?php echo (!isset($_POST['color']) || $_POST['color'] == '#FF6B6B') ? 'checked' : ''; ?>>
          <label for="color1" style="background-color: #FF6B6B;"></label>

          <input type="radio" name="color" id="color2" value="#FFD84D"
            <?php echo (isset($_POST['color']) && $_POST['color'] == '#FFD84D') ? 'checked' : ''; ?>>
          <label for="color2" style="background-color: #FFD84D;"></label>

          <input type="radio" name="color" id="color3" value="#8CCB5E"
            <?php echo (isset($_POST['color']) && $_POST['color'] == '#8CCB5E') ? 'checked' : ''; ?>>
          <label for="color3" style="background-color: #8CCB5E;"></label>

          <input type="radio" name="color" id="color4" value="#4DA3FF"
            <?php echo (isset($_POST['color']) && $_POST['color'] == '#4DA3FF') ? 'checked' : ''; ?>>
          <label for="color4" style="background-color: #4DA3FF;"></label>

          <input type="radio" name="color" id="color5" value="#A66BFF"
            <?php echo (isset($_POST['color']) && $_POST['color'] == '#A66BFF') ? 'checked' : ''; ?>>
          <label for="color5" style="background-color: #A66BFF;"></label>
        </div>
      </div>

      <div class="section-title">Configuración Personal</div>

      <div class="form-group2">
        <label>Desembolso planificado</label>
        <input type="number" id="desembolso_planejado" name="desembolso_planejado"
          placeholder="S/. 0" step="0.01" min="0"
          value="<?php echo htmlspecialchars($_POST['desembolso_planejado'] ?? '0'); ?>">
        <select id="periodo_tipo" name="periodo_tipo" class="frecuencia">
          <option value="">Sin frecuencia</option>
          <option value="diario" <?php echo ($_POST['periodo_tipo'] ?? '') === 'diario' ? 'selected' : ''; ?>>Diario</option>
          <option value="semanal" <?php echo ($_POST['periodo_tipo'] ?? '') === 'semanal' ? 'selected' : ''; ?>>Semanal</option>
          <option value="mensual" <?php echo ($_POST['periodo_tipo'] ?? '') === 'mensual' ? 'selected' : ''; ?>>Mensual</option>
        </select>
      </div>

      <div class="form-group2">
        <label>Establecer límite</label>
        <input type="number" id="limite_monto" name="limite_monto"
          placeholder="S/. 0" step="0.01" min="0"
          value="<?php echo htmlspecialchars($_POST['limite_monto'] ?? '0'); ?>">
        <select id="limite_tipo" name="limite_tipo" class="frecuencia">
          <option value="">Sin frecuencia</option>
          <option value="diario" <?php echo ($_POST['limite_tipo'] ?? '') === 'diario' ? 'selected' : ''; ?>>Diario</option>
          <option value="semanal" <?php echo ($_POST['limite_tipo'] ?? '') === 'semanal' ? 'selected' : ''; ?>>Semanal</option>
          <option value="mensual" <?php echo ($_POST['limite_tipo'] ?? '') === 'mensual' ? 'selected' : ''; ?>>Mensual</option>
        </select>
      </div>

      <div class="actions">
        <button id="btnGuardar" class="btn btn-yellow" type="submit">Guardar Concepto</button>
        <button id="btnCancelar" class="btn btn-gray" type="button" onclick="ocultarFormulario()">Cancelar</button>
      </div>
    </div>
  </form>
</div>

<style>
  .conceptos {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  .concepto {
    display: flex;
    align-items: center;
    padding: 15px;
    background: white;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    border-left: 4px solid transparent;
  }

  .concepto-icono {
    font-size: 24px;
    margin-right: 15px;
    width: 40px;
    text-align: center;
  }

  .concepto-info {
    flex: 1;
  }

  .concepto-nombre {
    font-weight: bold;
    margin-bottom: 5px;
  }

  .concepto-desembolso,
  .concepto-limite {
    font-size: 12px;
    color: #666;
  }

  .concepto-acciones {
    display: flex;
    gap: 5px;
  }

  .btn-editar,
  .btn-deshabilitar {
    background: none;
    border: none;
    cursor: pointer;
    padding: 5px;
    border-radius: 4px;
    transition: background 0.3s;
  }

  .btn-editar:hover {
    background: #e3f2fd;
  }

  .btn-deshabilitar:hover {
    background: #ffebee;
  }

  .concepto.nuevo {
    justify-content: center;
    text-align: center;
    cursor: pointer;
    background: #f8f9fa;
    border: 2px dashed #dee2e6;
    color: #6c757d;
  }

  .concepto.nuevo:hover {
    background: #e9ecef;
    border-color: #4DA3FF;
    color: #4DA3FF;
  }

  .concepto-vacio {
    text-align: center;
    padding: 40px 20px;
    color: #6c757d;
  }

  .concepto-vacio i {
    font-size: 48px;
    margin-bottom: 15px;
    opacity: 0.5;
  }

  .error-box {
    background: #ffe6e6;
    border: 1px solid #ffcccc;
    border-radius: 5px;
    padding: 10px;
    margin-bottom: 15px;
  }

  .success-box {
    background: #e6ffe6;
    border: 1px solid #ccffcc;
    border-radius: 5px;
    padding: 10px;
    margin-bottom: 15px;
  }
</style>

<script>
  function mostrarFormulario() {
    document.getElementById('formConcepto').classList.add('active');
  }

  function ocultarFormulario() {
    document.getElementById('formConcepto').classList.remove('active');
  }

  function editarConcepto(id) {
    // Aquí implementarás la edición
    alert('Editar concepto ID: ' + id);
  }

  function deshabilitarConcepto(id) {
    if (confirm('¿Estás seguro de que quieres deshabilitar este concepto?')) {
      // Aquí implementarás la deshabilitación
      window.location.href = '<?php echo URLROOT; ?>/concepto/deshabilitar/' + id;
    }
  }

  // Auto-seleccionar primer ícono si no hay ninguno seleccionado
  document.addEventListener('DOMContentLoaded', function() {
    const iconosSeleccionados = document.querySelectorAll('.icono-radio:checked');
    if (iconosSeleccionados.length === 0) {
      const primerIcono = document.querySelector('.icono-radio');
      if (primerIcono) primerIcono.checked = true;
    }
  });
</script>