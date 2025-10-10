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
      const dias = ["Lunes", "Martes", "Miércoles", "Jueves", "Viernes", "Sábado", "Domingo"];
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

function editarConcepto(i) {
  mostrarFormulario();
  modo = "editar";
  conceptoSeleccionado = i;
  const c = conceptos[i];

  tituloForm.textContent = "Editar Concepto";
  btnGuardar.textContent = "Actualizar";
  btnDeshabilitar.style.display = "inline-block";
  btnCancelar.style.display = "inline-block";

  document.getElementById("nombreConcepto").value = c.nombre;
  document.querySelectorAll(".iconos-grid i").forEach(icon => {
    icon.classList.toggle("selected", icon.classList.contains(c.icono));
  });
  document.querySelector(`input[name='color'][value='${c.color}']`).checked = true;
  document.getElementById("desembolsoMonto").value = c.desembolso.monto;
  document.getElementById("desembolsoFrecuencia").value = c.desembolso.frecuencia;
  document.getElementById("limiteMonto").value = c.limite.monto;
  document.getElementById("limiteFrecuencia").value = c.limite.frecuencia;

  const dias = ["Lunes", "Martes", "Miércoles", "Jueves", "Viernes", "Sábado", "Domingo"];

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
      if (frecuencia === "semanal") {
        fechaLabel.innerHTML = `📅 ${dias[fecha.getDay()]} de cada semana`;
      } else {
        fechaLabel.innerHTML = `📅 Día ${fecha.getDate()} de cada mes`;
      }
    } else {
      fechaLabel.innerHTML = "";
    }
  }

  mostrarInfoFechaEnFormulario("desembolsoFrecuencia", c.desembolso.frecuencia, c.desembolso.fecha);
  mostrarInfoFechaEnFormulario("limiteFrecuencia", c.limite.frecuencia, c.limite.fecha);

  renderConceptos();
}

function limpiarFormulario() {
  document.getElementById("nombreConcepto").value = "";
  document.querySelectorAll(".iconos-grid i").forEach(icon => icon.classList.remove("selected"));
  document.querySelectorAll("input[name='color']").forEach(r => r.checked = false);
  document.getElementById("desembolsoMonto").value = "";
  document.getElementById("limiteMonto").value = "";
  document.getElementById("desembolsoFrecuencia").value = "diario";
  document.getElementById("limiteFrecuencia").value = "diario";
}

function mostrarPopup(mensaje, tipo = "info") {
  popup.textContent = mensaje;
  popup.style.background = tipo === "success" ? "#2ecc71" : "#3498db";
  popup.classList.add("show");
  setTimeout(() => popup.classList.remove("show"), 2000); // Ensure timeout is sufficient
}

// ======= EVENTOS =======
btnGuardar.addEventListener("click", () => {
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

  fetch('/arka-code/concepto/guardarConcepto', {
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
    if (data.success) {
      mostrarPopup(data.message || "✅ Concepto guardado correctamente", "success");
      fetchConceptos().then(() => {
        nuevoConcepto();
      });
    } else {
      mostrarPopup(data.message || "Error al guardar el concepto", "info");
    }
  })
  .catch(error => {
    console.error('Error:', error);
    mostrarPopup("Error al guardar el concepto", "info");
  });
});

btnCancelar.addEventListener("click", nuevoConcepto);

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
// Función mejorada para cargar conceptos
function fetchConceptos() {
    return fetch('/arka-code/concepto/listar?format=json')
        .then(response => {
            if (!response.ok) throw new Error('Network response was not ok');
            return response.json();
        })
        .then(data => {
            console.log('Conceptos cargados:', data.conceptos); // Debug
            conceptos.length = 0; // Clear existing concepts
            
            data.conceptos.forEach(concepto => {
                conceptos.push({
                    id: concepto.id, // Asegúrate de que esto viene del backend
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

// Mejorar la función de guardar para recargar correctamente
btnGuardar.addEventListener("click", () => {
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

    // Mostrar loading
    btnGuardar.disabled = true;
    btnGuardar.textContent = 'Guardando...';

    fetch('/arka-code/concepto/guardarConcepto', {
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
        if (data.success) {
            mostrarPopup(data.message || "✅ Concepto guardado correctamente", "success");
            // Recargar conceptos después de guardar
            return fetchConceptos().then(() => {
                nuevoConcepto();
            });
        } else {
            mostrarPopup(data.message || "Error al guardar el concepto", "info");
        }
    })
    .catch(error => {
        console.error('Error:', error);
        mostrarPopup("Error al guardar el concepto", "info");
    })
    .finally(() => {
        // Restaurar botón
        btnGuardar.disabled = false;
        btnGuardar.textContent = modo === 'crear' ? 'Guardar' : 'Actualizar';
    });
});

// ======= INICIALIZACIÓN =======
fetchConceptos();