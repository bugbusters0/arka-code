// ======= DATOS INICIALES =======
const conceptos = [
    { nombre: "Mercado", icono: "fa-basket-shopping", color: "#FF6B6B", desembolso: { monto: 150, frecuencia: "mensual",fecha:"2025-09-04" }, limite: { monto: 200, frecuencia: "mensual" } },
    { nombre: "Ropa", icono: "fa-shirt", color: "#FFD84D", desembolso: { monto: 100, frecuencia: "mensual" }, limite: { monto: 150, frecuencia: "mensual" } },
    { nombre: "Restaurante", icono: "fa-utensils", color: "#8CCB5E", desembolso: { monto: 80, frecuencia: "semanal" }, limite: { monto: 100, frecuencia: "semanal" } },
  ];
  
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

    // Función para mostrar la fecha seleccionada debajo del select
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
            const dias = ["Lunes","Martes","Miércoles","Jueves","Viernes","Sábado","Domingo"];
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
    // Abre modal al cambiar a semanal o mensual
    select.addEventListener("change", () => {
        if (select.value !== "diario") {
        modal.style.display = "flex";
        currentSelect = select;
        } else {
        // Si vuelve a diario, se elimina la fecha guardada
        delete select.dataset.fecha;
        mostrarFechaSeleccionada(select);
        }
    });

    // También abre modal si se hace clic en la opción actual (por ejemplo, re-editar)
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

    // Mostrar fecha seleccionada si existe
    const dias = ["Lunes","Martes","Miércoles","Jueves","Viernes","Sábado","Domingo"];

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
    setTimeout(() => popup.classList.remove("show"), 20);
  }
  
  // ======= EVENTOS =======
  btnGuardar.addEventListener("click", () => {
    const nombre = document.getElementById("nombreConcepto").value.trim();
    const icono = document.querySelector(".iconos-grid i.selected")?.classList[1];
    const color = document.querySelector("input[name='color']:checked")?.value;
    const desembolso = {
      monto: parseFloat(document.getElementById("desembolsoMonto").value || 0),
      frecuencia: document.getElementById("desembolsoFrecuencia").value
    };
    const limite = {
      monto: parseFloat(document.getElementById("limiteMonto").value || 0),
      frecuencia: document.getElementById("limiteFrecuencia").value
    };
  
    if (!nombre || !icono || !color) return;
  
    if (modo === "crear") {
      conceptos.unshift({ nombre, icono, color, desembolso, limite });
      mostrarPopup("✅ Concepto guardado", "success");
    } else {
      conceptos[conceptoSeleccionado] = { nombre, icono, color, desembolso, limite };
      mostrarPopup("✏️ Concepto actualizado", "info");
    }
  
    nuevoConcepto();
    renderConceptos();
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
  
  // ======= INICIALIZACIÓN =======
  renderConceptos();