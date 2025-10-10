// ======= DATOS INICIALES =======
const conceptos = []; // Will be populated from DB

// ======= ELEMENTOS =======
const lista = document.getElementById("listaConceptos");
const tituloForm = document.getElementById("tituloForm");
const btnGuardar = document.getElementById("btnGuardar");
const btnDeshabilitar = document.getElementById("btnDeshabilitar");
const btnCancelar = document.getElementById("btnCancelar");
const popup = document.getElementById("popup");

let modo = "crear";
let conceptoSeleccionado = null;

// ======= FUNCIONES =======
function renderConceptos() {
    lista.innerHTML = "";

    conceptos.forEach((c, i) => {
        const div = document.createElement("div");
        div.className = "concepto";
        div.style.borderColor = conceptoSeleccionado === i ? "#4DA3FF" : "transparent";
        div.innerHTML = `<i class="fa-solid ${c.icono}"></i>${c.nombre}`;
        div.addEventListener("click", () => editarConcepto(i));
        lista.appendChild(div);
    });

    // Botón "Nuevo"
    const nuevo = document.createElement("div");
    nuevo.className = "concepto nuevo";
    nuevo.innerHTML = '<i class="fa-solid fa-plus"></i>Nuevo';
    nuevo.addEventListener("click", nuevoConcepto);
    lista.appendChild(nuevo);
}

// ======= MODAL CALENDARIO =======
const selectsFrecuencia = document.querySelectorAll(".frecuencia");
const modal = document.getElementById("modalCalendario");
const closeBtn = document.querySelector(".modal .close");
const guardarFecha = document.getElementById("guardarFecha");
const calendarInput = document.getElementById("calendarInput");

let currentSelect = null;

function mostrarFechaSeleccionada(select) {
    let fechaLabel = select.parentElement.querySelector(".fecha-seleccionada");
    if (!fechaLabel) {
        fechaLabel = document.createElement("div");
        fechaLabel.className = "fecha-seleccionada";
        select.parentElement.appendChild(fechaLabel);
    }
    if (select.dataset.fecha) {
        const date = new Date(select.dataset.fecha);
        let textoFecha = "";
        if (select.value === "semanal") {
            const dias = ["Domingo", "Lunes", "Martes", "Miércoles", "Jueves", "Viernes", "Sábado"];
            textoFecha = `📅 ${dias[date.getDay()]} de cada semana`;
        } else if (select.value === "mensual") {
            textoFecha = `📅 Día ${date.getDate()} de cada mes`;
        } else {
            textoFecha = `📅 ${select.dataset.fecha}`;
        }
        fechaLabel.innerHTML = `<small>${textoFecha}</small>`;
    } else {
        fechaLabel.innerHTML = "";
    }
}

selectsFrecuencia.forEach(select => {
    select.addEventListener("change", () => {
        if (select.value !== "diario") {
            modal.style.display = "flex";
            currentSelect = select;
        } else {
            delete select.dataset.fecha;
            mostrarFechaSeleccionada(select);
        }
    });

    select.addEventListener("click", () => {
        if (select.value !== "diario") {
            modal.style.display = "flex";
            currentSelect = select;
        }
    });
});

if (closeBtn) {
    closeBtn.addEventListener("click", () => {
        modal.style.display = "none";
    });
}

if (guardarFecha) {
    guardarFecha.addEventListener("click", () => {
        if (currentSelect) {
            currentSelect.dataset.fecha = calendarInput.value;
            mostrarFechaSeleccionada(currentSelect);
        }
        modal.style.display = "none";
    });
}

window.addEventListener("click", (e) => {
    if (e.target === modal) {
        modal.style.display = "none";
    }
});

// ======= RESTO DE FUNCIONES =======
function nuevoConcepto() {
    mostrarFormulario();
    modo = "crear";
    conceptoSeleccionado = null;
    tituloForm.textContent = "Crear Concepto";
    btnGuardar.textContent = "Guardar";
    btnDeshabilitar.style.display = "none";
    btnCancelar.style.display = "none";
    limpiarFormulario();
    renderConceptos();
}

function editarConcepto(index) {
    const c = conceptos[index];
    console.log("🔄 Editando concepto:", c);
    
    // Cargar datos del concepto para editar
    fetch(`/arka-code/concepto/editar/${c.id}`)
        .then(response => {
            if (!response.ok) throw new Error('Network response was not ok');
            return response.json();
        })
        .then(data => {
            if (data.success && data.concepto) {
                const concepto = data.concepto;
                mostrarFormulario();
                modo = "editar";
                conceptoSeleccionado = index;
                
                tituloForm.textContent = "Editar Concepto";
                btnGuardar.textContent = "Actualizar";
                btnDeshabilitar.style.display = "inline-block";
                btnCancelar.style.display = "inline-block";

                // Llenar formulario con datos del concepto
                document.getElementById("nombreConcepto").value = concepto.nombre;
                
                // Seleccionar ícono
                document.querySelectorAll(".iconos-grid i").forEach(icon => {
                    const iconClass = 'fa-' + (concepto.icono_nombre || 'question').toLowerCase().replace(/ /g, '-');
                    icon.classList.toggle("selected", icon.classList.contains(iconClass));
                });
                
                // Seleccionar color
                const colorInput = document.querySelector(`input[name='color'][value='${concepto.color}']`);
                if (colorInput) {
                    colorInput.checked = true;
                } else {
                    // Si no encuentra el color, seleccionar el primero
                    document.querySelector('input[name="color"]').checked = true;
                }
                
                // Llenar configuración personal
                document.getElementById("desembolsoMonto").value = concepto.desembolso_planejado || 0;
                document.getElementById("desembolsoFrecuencia").value = concepto.desembolso_frecuencia || 'diario';
                document.getElementById("limiteMonto").value = concepto.limite_monto || 0;
                document.getElementById("limiteFrecuencia").value = concepto.limite_frecuencia || 'diario';

                // Mostrar fechas si existen
                mostrarInfoFechaEnFormulario("desembolsoFrecuencia", concepto.desembolso_frecuencia, concepto.periodo_desembolso);
                mostrarInfoFechaEnFormulario("limiteFrecuencia", concepto.limite_frecuencia, concepto.limite_fecha);

                renderConceptos();
            } else {
                mostrarPopup(data.message || "Error al cargar el concepto", "info");
            }
        })
        .catch(error => {
            console.error('Error cargando concepto:', error);
            mostrarPopup("Error al cargar el concepto", "info");
        });
}

// Función auxiliar para mostrar fechas en formulario
function mostrarInfoFechaEnFormulario(selectId, frecuencia, fechaValor) {
    const select = document.getElementById(selectId);
    let fechaLabel = select.parentElement.querySelector(".fecha-seleccionada");
    if (!fechaLabel) {
        fechaLabel = document.createElement("div");
        fechaLabel.className = "fecha-seleccionada";
        select.parentElement.appendChild(fechaLabel);
    }

    if (fechaValor && (frecuencia === "semanal" || frecuencia === "mensual")) {
        const fecha = new Date(fechaValor);
        const dias = ["Domingo", "Lunes", "Martes", "Miércoles", "Jueves", "Viernes", "Sábado"];
        
        if (frecuencia === "semanal") {
            fechaLabel.innerHTML = `📅 ${dias[fecha.getDay()]} de cada semana`;
        } else {
            fechaLabel.innerHTML = `📅 Día ${fecha.getDate()} de cada mes`;
        }
        select.dataset.fecha = fechaValor.split('T')[0]; // Guardar fecha en formato YYYY-MM-DD
    } else {
        fechaLabel.innerHTML = "";
        delete select.dataset.fecha;
    }
}

function limpiarFormulario() {
    document.getElementById("nombreConcepto").value = "";
    document.querySelectorAll(".iconos-grid i").forEach(icon => icon.classList.remove("selected"));
    document.querySelectorAll("input[name='color']").forEach(r => r.checked = false);
    document.getElementById("desembolsoMonto").value = "";
    document.getElementById("limiteMonto").value = "";
    document.getElementById("desembolsoFrecuencia").value = "diario";
    document.getElementById("limiteFrecuencia").value = "diario";
    
    // Limpiar labels de fecha
    document.querySelectorAll(".fecha-seleccionada").forEach(label => {
        label.innerHTML = "";
    });
}

function mostrarPopup(mensaje, tipo = "info") {
    popup.textContent = mensaje;
    popup.style.background = tipo === "success" ? "#2ecc71" : "#3498db";
    popup.classList.add("show");
    setTimeout(() => popup.classList.remove("show"), 2000);
}

// ======= FUNCIÓN PRINCIPAL PARA GUARDAR/ACTUALIZAR =======
function guardarConceptoHandler() {
    if (btnGuardar.disabled) return;
    
    const nombre = document.getElementById("nombreConcepto").value.trim();
    const icono = document.querySelector(".iconos-grid i.selected")?.classList[1];
    const color = document.querySelector("input[name='color']:checked")?.value;
    const desembolso = {
        monto: parseFloat(document.getElementById("desembolsoMonto").value || 0),
        frecuencia: document.getElementById("desembolsoFrecuencia").value,
        fecha: document.getElementById("desembolsoFrecuencia").dataset.fecha
    };
    const limite = {
        monto: parseFloat(document.getElementById("limiteMonto").value || 0),
        frecuencia: document.getElementById("limiteFrecuencia").value,
        fecha: document.getElementById("limiteFrecuencia").dataset.fecha
    };

    if (!nombre || !icono || !color) {
        mostrarPopup("Por favor, complete todos los campos requeridos", "info");
        return;
    }

    const iconMap = {
        'fa-burger': 3,
        'fa-spa': 4,
        'fa-car': 5,
        'fa-heart-pulse': 6,
        'fa-graduation-cap': 7,
        'fa-basket-shopping': 8,
        'fa-utensils': 9,
        'fa-lightbulb': 10,
        'fa-shirt': 11
    };
    const idIcono = iconMap[icono.split('-')[2]] || 3;

    // Determinar URL y método según el modo
    const url = modo === 'crear' 
        ? '/arka-code/concepto/guardarConcepto'
        : `/arka-code/concepto/editar/${conceptos[conceptoSeleccionado].id}`;

    // Deshabilitar botón
    btnGuardar.disabled = true;
    btnGuardar.textContent = modo === 'crear' ? 'Guardando...' : 'Actualizando...';

    console.log("📤 Enviando datos a:", url);

    fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: `name=${encodeURIComponent(nombre)}&tipo=gasto&color=${encodeURIComponent(color)}&idIcono=${idIcono}&desembolsoMonto=${desembolso.monto}&desembolsoFrecuencia=${desembolso.frecuencia}&limiteMonto=${limite.monto}&limiteFrecuencia=${limite.frecuencia}`
    })
    .then(response => {
        if (!response.ok) throw new Error('Network response was not ok');
        return response.json();
    })
    .then(data => {
        console.log("📥 Respuesta recibida:", data);
        if (data.success) {
            mostrarPopup(data.message || (modo === 'crear' ? "✅ Concepto guardado" : "✅ Concepto actualizado"), "success");
            // Recargar conceptos
            setTimeout(() => {
                fetchConceptos().then(() => {
                    if (modo === 'crear') {
                        nuevoConcepto();
                    } else {
                        // En modo editar, mantener el formulario abierto pero resetear selección
                        conceptoSeleccionado = null;
                        renderConceptos();
                    }
                });
            }, 1000);
        } else {
            mostrarPopup(data.message || "Error al procesar el concepto", "info");
        }
    })
    .catch(error => {
        console.error('❌ Error:', error);
        mostrarPopup("Error al procesar el concepto", "info");
    })
    .finally(() => {
        btnGuardar.disabled = false;
        btnGuardar.textContent = modo === 'crear' ? 'Guardar' : 'Actualizar';
    });
}

// ======= EVENTOS =======
// Remover event listeners anteriores para evitar duplicados
btnGuardar.removeEventListener("click", guardarConceptoHandler);
btnGuardar.addEventListener("click", guardarConceptoHandler);

btnCancelar.addEventListener("click", nuevoConcepto);

// btnDeshabilitar.addEventListener("click", () => {
//     if (modo === "editar" && conceptoSeleccionado !== null) {
//         const conceptoId = conceptos[conceptoSeleccionado].id;
//         const conceptoNombre = conceptos[conceptoSeleccionado].nombre;
        
//         if (confirm(`¿Estás seguro de que quieres deshabilitar el concepto "${conceptoNombre}"?`)) {
//             // Aquí puedes implementar la lógica para deshabilitar/eliminar
//             mostrarPopup("Función de deshabilitar en desarrollo", "info");
//         }
//     }
// });
// ======= EVENTO PARA EL BOTÓN DESHABILITAR =======
btnDeshabilitar.addEventListener("click", () => {
    if (modo === "editar" && conceptoSeleccionado !== null) {
        const conceptoId = conceptos[conceptoSeleccionado].id;
        const conceptoNombre = conceptos[conceptoSeleccionado].nombre;
        
        if (confirm(`¿Estás seguro de que quieres deshabilitar el concepto "${conceptoNombre}"?\n\nEsto solo lo ocultará para ti, otros usuarios seguirán viéndolo.`)) {
            // Deshabilitar botón
            btnDeshabilitar.disabled = true;
            btnDeshabilitar.textContent = 'Deshabilitando...';
            
            fetch(`/arka-code/concepto/deshabilitar/${conceptoId}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded',
                }
            })
            .then(response => {
                if (!response.ok) throw new Error('Network response was not ok');
                return response.json();
            })
            .then(data => {
                if (data.success) {
                    mostrarPopup(data.message || "✅ Concepto deshabilitado correctamente", "success");
                    // Recargar conceptos y volver al modo crear
                    setTimeout(() => {
                        fetchConceptos().then(() => {
                            nuevoConcepto();
                        });
                    }, 1000);
                } else {
                    mostrarPopup(data.message || "Error al deshabilitar el concepto", "info");
                }
            })
            .catch(error => {
                console.error('❌ Error:', error);
                mostrarPopup("Error al deshabilitar el concepto", "info");
            })
            .finally(() => {
                btnDeshabilitar.disabled = false;
                btnDeshabilitar.textContent = 'Deshabilitar';
            });
        }
    }
});
document.querySelectorAll(".iconos-grid i").forEach(icon => {
    icon.addEventListener("click", () => {
        document.querySelectorAll(".iconos-grid i").forEach(i => i.classList.remove("selected"));
        icon.classList.add("selected");
    });
});

function mostrarFormulario() {
    document.querySelector(".sidebar-right").classList.add("active");
}

function ocultarFormulario() {
    document.querySelector(".sidebar-right").classList.remove("active");
}

// Function to fetch concepts from the server
function fetchConceptos() {
    return fetch('/arka-code/concepto/listar?format=json')
        .then(response => {
            if (!response.ok) throw new Error('Network response was not ok');
            return response.json();
        })
        .then(data => {
            console.log('Conceptos cargados:', data.conceptos);
            conceptos.length = 0;
            
            data.conceptos.forEach(concepto => {
                conceptos.push({
                    id: concepto.id,
                    nombre: concepto.nombre,
                    icono: 'fa-' + (concepto.icono_nombre || 'question').toLowerCase().replace(/ /g, '-'),
                    color: concepto.color || '#FF6B6B',
                    desembolso: {
                        monto: concepto.desembolso_planejado || 0,
                        frecuencia: concepto.desembolso_frecuencia || 'diario',
                        fecha: concepto.periodo_desembolso || null
                    },
                    limite: {
                        monto: concepto.limite_monto || 0,
                        frecuencia: concepto.limite_frecuencia || 'diario',
                        fecha: concepto.limite_fecha || null
                    }
                });
            });
            renderConceptos();
        })
        .catch(error => {
            console.error('Error fetching concepts:', error);
            mostrarPopup("Error al cargar los conceptos", "info");
        });
}

// ======= INICIALIZACIÓN =======
document.addEventListener('DOMContentLoaded', function() {
    console.log("🔄 Script concepto-listar.js cargado e inicializado");
    fetchConceptos();
});