const fechaInput = document.getElementById('fechaInput');
  const fechaTexto = document.getElementById('fechaTexto');
  const fechaFormTexto = document.getElementById('fechaFormTexto');
  const fechaHidden = document.getElementById('fechaHidden');
  const tipoActual = window.tipoActual;

  // Formatear fecha para mostrar
  function formatearFecha(fechaISO) {
    const fecha = new Date(fechaISO + 'T00:00:00');
    const opciones = { year: 'numeric', month: 'long', day: 'numeric' };
    return fecha.toLocaleDateString('es-PE', opciones);
  }

  // Inicializar fecha
  fechaTexto.textContent = formatearFecha(fechaInput.value);
  if (fechaFormTexto) {
    fechaFormTexto.textContent = formatearFecha(fechaInput.value);
  }

  // Selector de fecha
  document.getElementById('fechaDisplay').addEventListener('click', () => {
    if (fechaInput.showPicker) {
      fechaInput.showPicker();
    } else {
      fechaInput.click();
    }
  });

  fechaInput.addEventListener('change', () => {
    const nuevaFecha = fechaInput.value;
    window.location.href = `/movimientos?tipo=${tipoActual}&fecha=${nuevaFecha}`;
  });

  // Selección de concepto
  let selectedConcepto = '';

  function selectConcept(nombre) {
    selectedConcepto = nombre;
    document.getElementById('nombreConcepto').value = nombre;
    
    // Resaltar concepto seleccionado
    document.querySelectorAll('.concept-item').forEach(el => {
      el.classList.remove('selected');
    });
    
    const selectedItem = document.querySelector(`.concept-item[data-nombre="${nombre}"]`);
    if (selectedItem) {
      selectedItem.classList.add('selected');
    }
  }

  // Editar movimiento
  function editarMovimiento(id, concepto, monto, descripcion) {
    document.getElementById('idMovimiento').value = id;
    document.getElementById('nombreConcepto').value = concepto;
    document.getElementById('monto').value = monto;
    document.getElementById('descripcion').value = descripcion || '';
    
    selectConcept(concepto);
    
    document.getElementById('tituloForm').textContent = 
      tipoActual === 'gasto' ? 'Editar Gasto' : 'Editar Ingreso';
    document.getElementById('btnGuardar').textContent = 'Actualizar';
    document.getElementById('movimientoForm').action = '/movimientos/editar';
    
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }

  // Editar desde resumen
  function editarEnResumen(id, tipo) {
    const fecha = fechaInput.value;
    window.location.href = `/movimientos?tipo=${tipo}&fecha=${fecha}&edit=${id}`;
  }

  // Eliminar movimiento
  function eliminarMovimiento(id) {
    if (!confirm('¿Estás seguro de que quieres eliminar este movimiento?')) {
      return;
    }
    
    const form = document.createElement('form');
    form.method = 'POST';
    form.action = '/movimientos/eliminar';
    
    const inputId = document.createElement('input');
    inputId.type = 'hidden';
    inputId.name = 'idMovimiento';
    inputId.value = id;
    
    const inputTipo = document.createElement('input');
    inputTipo.type = 'hidden';
    inputTipo.name = 'tipo';
    inputTipo.value = tipoActual;
    
    const inputFecha = document.createElement('input');
    inputFecha.type = 'hidden';
    inputFecha.name = 'fecha';
    inputFecha.value = fechaInput.value;
    
    form.appendChild(inputId);
    form.appendChild(inputTipo);
    form.appendChild(inputFecha);
    document.body.appendChild(form);
    form.submit();
  }

  // Reset form
  function resetForm() {
    document.getElementById('movimientoForm').reset();
    document.getElementById('idMovimiento').value = '';
    document.getElementById('movimientoForm').action = '/movimientos/crear';
    document.getElementById('tituloForm').textContent = 
      tipoActual === 'gasto' ? 'Crear Gasto' : 'Crear Ingreso';
    document.getElementById('btnGuardar').textContent = 
      tipoActual === 'gasto' ? 'Guardar Gasto' : 'Guardar Ingreso';
    
    selectedConcepto = '';
    document.querySelectorAll('.concept-item').forEach(el => {
      el.classList.remove('selected');
    });
  }

  // Auto-ocultar mensajes después de 3.5 segundos
  setTimeout(() => {
    const messages = document.querySelectorAll('.success-message, .error-message');
    messages.forEach(msg => msg.style.display = 'none');
  }, 3500);