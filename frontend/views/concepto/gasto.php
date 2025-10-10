<div class="tabs">
  <a href="<?php echo URLROOT . "/concepto/gasto" ?>" class="tab active">Gastos</a>
  <a href="<?php echo URLROOT . "/concepto/ingresos" ?>" class="tab">Ingresos</a>
</div>

<div class="content-area">
  <!-- Lista de conceptos -->
  <div class="conceptos" id="listaConceptos"></div>

  <!-- Formulario lateral -->
  <form method="post" action="<?php echo URLROOT; ?>/concepto/crear">
    <div class="sidebar-right" id="formConcepto">
      <h3 id="tituloForm">Crear Concepto</h3>
      <input type="hidden" name="type" value="gasto">
      <div class="form-group">
        <label>Nombre del Concepto</label>
        <input type="text" id="nombreConcepto" placeholder="Ej: Transporte">
      </div>

      <div class="form-group">
        <label>Ícono</label>
        <div id="iconosGrid" class="iconos-grid">
          <?php foreach ($icons as $icon): ?>
            <input type="radio" name="id_icono" id="icono_<?php echo $icon->getId(); ?>"
              value="<?php echo $icon->getId(); ?>" class="icono-radio" required>
            <label for="icono_<?php echo $icon->getId(); ?>" class="icono-label">
              <i class="<?php echo htmlspecialchars($icon->getPath()); ?>"></i>
              <span><?php echo htmlspecialchars($icon->getNombre()); ?></span>
            </label>
          <?php endforeach; ?>
        </div>
      </div>

      <div class="form-group">
        <label>Color</label>
        <div class="color-picker">
          <input type="radio" name="color" id="color1" value="#FF6B6B">
          <label for="color1" style="background-color: #FF6B6B;"></label>

          <input type="radio" name="color" id="color2" value="#FFD84D">
          <label for="color2" style="background-color: #FFD84D;"></label>

          <input type="radio" name="color" id="color3" value="#8CCB5E">
          <label for="color3" style="background-color: #8CCB5E;"></label>

          <input type="radio" name="color" id="color4" value="#4DA3FF">
          <label for="color4" style="background-color: #4DA3FF;"></label>

          <input type="radio" name="color" id="color5" value="#A66BFF">
          <label for="color5" style="background-color: #A66BFF;"></label>
        </div>
      </div>

      <div class="section-title">Configuración Personal</div>

      <div class="form-group2">
        <label>Desembolso planificado</label>
        <input type="number" id="desembolsoMonto" placeholder="S/. 0">
        <select id="desembolsoFrecuencia" class="frecuencia">
          <option value="diario">Diario</option>
          <option value="semanal">Semanal</option>
          <option value="mensual">Mensual</option>
        </select>
      </div>

      <div class="form-group2">
        <label>Establecer límite</label>
        <input type="number" id="limiteMonto" placeholder="S/. 0">
        <select id="limiteFrecuencia" class="frecuencia">
          <option value="diario">Diario</option>
          <option value="semanal">Semanal</option>
          <option value="mensual">Mensual</option>
        </select>
      </div>

      <div class="actions">
        <button id="btnDeshabilitar" class="btn btn-red" style="display:none;">Deshabilitar</button>
        <button id="btnGuardar" class="btn btn-yellow">Guardar</button>
        <button id="btnCancelar" class="btn btn-gray" style="display:none;">Cancelar</button>
      </div>
    </div>
  </form>

</div>

<div id="popup" class="popup"></div>

<!-- Modal -->
<div id="modalCalendario" class="modal">
  <div class="modal-content">
    <span class="close">&times;</span>
    <h3>Selecciona la fecha</h3>
    <input type="date" id="calendarInput">
    <button id="guardarFecha" class="btn btn-yellow">Guardar</button>
  </div>
</div>