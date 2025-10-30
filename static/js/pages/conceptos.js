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

function editarConcepto(id) {
    console.log("✏️ Editando concepto:", id);
    // Aquí implementarás la edición
    alert('Editar concepto: ' + id);
}

function deshabilitarConcepto(id) {
    console.log("🚫 Deshabilitando concepto:", id);
    if (confirm('¿Estás seguro de que quieres deshabilitar este concepto?')) {
        // Aquí implementarás la deshabilitación
        window.location.href = '/conceptos/deshabilitar/' + id;
    }
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
    
    // Debug del formulario de creación
    const formulario = document.querySelector('form[action="/conceptos/crear"]');
    if (formulario) {
        console.log("✅ Formulario encontrado");
        
        formulario.addEventListener('submit', function(e) {
            console.log("🚀 FORMULARIO ENVIÁNDOSE...");
            
            // Mostrar todos los datos del formulario
            const formData = new FormData(formulario);
            console.log("📝 DATOS DEL FORMULARIO:");
            for (let [key, value] of formData.entries()) {
                console.log(`   ${key}: ${value}`);
            }
            
            // Verificar campos requeridos
            const nombre = formData.get('nombre');
            const idIcono = formData.get('id_icono');
            const color = formData.get('color');
            
            console.log("🔍 CAMPOS REQUERIDOS:");
            console.log(`   nombre: ${nombre}`);
            console.log(`   id_icono: ${idIcono}`);
            console.log(`   color: ${color}`);
            
            if (!nombre || !idIcono || !color) {
                console.log("❌ FALTAN CAMPOS REQUERIDOS");
                alert('Por favor completa todos los campos requeridos');
                e.preventDefault();
                return;
            }
            
            console.log("✅ FORMULARIO VÁLIDO, ENVIANDO...");
        });
    } else {
        console.log("❌ NO se encontró el formulario");
    }
    
    // Verificar botón "Nuevo Concepto"
    const btnNuevo = document.querySelector('.concepto.nuevo');
    if (btnNuevo) {
        console.log("✅ Botón 'Nuevo Concepto' encontrado");
        btnNuevo.addEventListener('click', function() {
            console.log("🖱️ Botón 'Nuevo Concepto' clickeado");
        });
    }
    
    // Verificar que el sidebar del formulario existe
    const sidebarForm = document.getElementById('formConcepto');
    if (sidebarForm) {
        console.log("✅ Sidebar del formulario encontrado");
        console.log("   Clases:", sidebarForm.className);
    } else {
        console.log("❌ NO se encontró el sidebar del formulario");
    }
});