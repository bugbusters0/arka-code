// Debug inicial
console.log("✅ conceptos.js cargado");

function mostrarFormulario() {
    console.log("📝 Mostrando formulario");
    document.getElementById('formConcepto').classList.add('active');
}

function ocultarFormulario() {
    console.log("❌ Ocultando formulario");
    document.getElementById('formConcepto').classList.remove('active');
}

function editarConcepto(nombreConcepto) {
  window.location.href = `/conceptos/editar?nombre=${encodeURIComponent(nombreConcepto)}`;
}

// Deshabilitar concepto
function deshabilitarConcepto(nombreConcepto, tipo) {
    if (confirm('¿Estás seguro de que quieres deshabilitar este concepto?\n\nLos usuarios ya no podrán seleccionarlo para nuevos movimientos.')) {
        const formData = new FormData();
        formData.append('nombreConcepto', nombreConcepto);
        formData.append('tipo', tipo);

        fetch('/conceptos/deshabilitar', {
            method: 'POST',
            body: formData
        })
        .then(response => {
            if (response.redirected) {
                window.location.href = response.url;
            } else if (!response.ok) {
                alert('Error al deshabilitar el concepto');
            }
        })
        .catch(error => {
            console.error('Error:', error);
            alert('Error al deshabilitar el concepto');
        });
    }
}

// Habilitar concepto
function habilitarConcepto(nombreConcepto, tipo) {
    const formData = new FormData();
    formData.append('nombreConcepto', nombreConcepto);
    formData.append('tipo', tipo);

    fetch('/conceptos/habilitar', {
        method: 'POST',
        body: formData
    })
    .then(response => {
        if (response.redirected) {
            window.location.href = response.url;
        } else if (!response.ok) {
            alert('Error al habilitar el concepto');
        }
    })
    .catch(error => {
        console.error('Error:', error);
        alert('Error al habilitar el concepto');
    });
}

// Debug del formulario cuando se carga la página
document.addEventListener('DOMContentLoaded', function() {
    console.log("📄 Página conceptos cargada");
    
    // Auto-seleccionar primer ícono si no hay ninguno seleccionado
    const iconosSeleccionados = document.querySelectorAll('.icono-radio:checked');
    if (iconosSeleccionados.length === 0) {
        const primerIcono = document.querySelector('.icono-radio');
        if (primerIcono) {
            primerIcono.checked = true;
            console.log("✅ Primer ícono auto-seleccionado");
        }
    }
    
    // Auto-seleccionar primer color si no hay ninguno seleccionado
    const coloresSeleccionados = document.querySelectorAll('input[name="color"]:checked');
    if (coloresSeleccionados.length === 0) {
        const primerColor = document.querySelector('input[name="color"]');
        if (primerColor) {
            primerColor.checked = true;
            console.log("🎨 Primer color auto-seleccionado");
        }
    }
    
    // Debug del formulario de creación/edición
    const formularios = document.querySelectorAll('form.form-crear');
    formularios.forEach(formulario => {
        console.log("✅ Formulario encontrado");
        
        formulario.addEventListener('submit', function(e) {
            console.log("🚀 FORMULARIO ENVIÁNDOSE...");
            
            const formData = new FormData(formulario);
            console.log("📝 DATOS DEL FORMULARIO:");
            for (let [key, value] of formData.entries()) {
                console.log(`   ${key}: ${value}`);
            }
        });
    });
    
    // Verificar botón "Nuevo Concepto"
    const btnNuevo = document.querySelector('.concepto.nuevo');
    if (btnNuevo) {
        console.log("✅ Botón 'Nuevo Concepto' encontrado");
    }
    
    // Verificar que el sidebar del formulario existe
    const sidebarForm = document.getElementById('formConcepto');
    if (sidebarForm) {
        console.log("✅ Sidebar del formulario encontrado");
        console.log("   Clases:", sidebarForm.className);
    }

    // ===== CONFIGURAR CAMPOS DEPENDIENTES =====
    configurarCampoDependiente("periodo_tipo", "dia_desembolso_planejado", "mensual");
    configurarCampoDependiente("limite_tipo", "dia_limite_tipo", "mensual");
});

function configurarCampoDependiente(selectId, targetId, valorActivador) {
    const selectElement = document.getElementById(selectId);
    const targetElement = document.getElementById(targetId);

    // ✅ VALIDAR QUE LOS ELEMENTOS EXISTAN
    if (!selectElement) {
        console.log(`⚠️ Select '${selectId}' no encontrado (puede ser modo edición)`);
        return;
    }
    
    if (!targetElement) {
        console.log(`⚠️ Input '${targetId}' no encontrado (puede ser modo edición)`);
        return;
    }

    console.log(`✅ Configurando campo dependiente: ${selectId} -> ${targetId}`);

    const handlePeriodChange = (e) => {
        const valorSeleccionado = e.target.value.toLowerCase();
        const esMensual = valorSeleccionado === valorActivador.toLowerCase();
        const esQuincenal = valorSeleccionado === "quincenal";

        targetElement.disabled = !esMensual && !esQuincenal;
        
        if (esMensual) {
            targetElement.placeholder = "Día (1-31)";
            targetElement.min = "1";
            targetElement.max = "31";
        } else if (esQuincenal) {
            targetElement.placeholder = "Día (1-15)";
            targetElement.min = "1";
            targetElement.max = "15";
        } else {
            targetElement.placeholder = "";
            targetElement.value = "";
        }

        console.log(`📝 ${selectId} cambió a: ${valorSeleccionado}, campo ${esMensual || esQuincenal ? 'habilitado' : 'deshabilitado'}`);
    };

    // Escuchar cambios
    selectElement.addEventListener('change', handlePeriodChange);
    
    // Ejecutar al cargar para aplicar estado inicial
    handlePeriodChange({ target: selectElement });
}