// static/js/pages/balance.js

document.addEventListener('DOMContentLoaded', function() {
  animarElementos();
  configurarScrolls();
  dibujarGraficosPie();
});

// ========================================
// NAVEGACIÓN Y FILTROS
// ========================================

function filtrarUsuario(nombreUsuario) {
  const url = new URL(window.location);
  
  if (nombreUsuario === '') {
    url.searchParams.delete('usuario');
  } else {
    url.searchParams.set('usuario', nombreUsuario);
  }
  
  window.location.href = url.toString();
}

function cambiarTipoGrafico(tipo) {
  const url = new URL(window.location);
  url.searchParams.set('tipo', 'graficos');
  url.searchParams.set('tipoGrafico', tipo);
  window.location.href = url.toString();
}

function cambiarVistaGrafico(vista) {
  const url = new URL(window.location);
  url.searchParams.set('tipo', 'graficos');
  url.searchParams.set('vistaGrafico', vista);
  window.location.href = url.toString();
}

// ========================================
// ANIMACIONES
// ========================================

function animarElementos() {
  // Animar columnas de balance
  const columns = document.querySelectorAll('.balance-column');
  columns.forEach((col, index) => {
    col.style.opacity = '0';
    col.style.transform = 'translateY(20px)';
    
    setTimeout(() => {
      col.style.transition = 'all 0.5s ease';
      col.style.opacity = '1';
      col.style.transform = 'translateY(0)';
    }, index * 150);
  });

  // Animar movimientos
  const movements = document.querySelectorAll('.movement-card');
  movements.forEach((mov, index) => {
    if (index < 15) { // Solo animar los primeros 15
      mov.style.opacity = '0';
      mov.style.transform = 'translateX(-20px)';
      
      setTimeout(() => {
        mov.style.transition = 'all 0.3s ease';
        mov.style.opacity = '1';
        mov.style.transform = 'translateX(0)';
      }, 300 + (index * 30));
    }
  });

  // Animar barras de gráfico
  const bars = document.querySelectorAll('.bar');
  bars.forEach((bar, index) => {
    const height = bar.style.height;
    bar.style.height = '0%';
    
    setTimeout(() => {
      bar.style.transition = 'height 0.8s ease';
      bar.style.height = height;
    }, index * 100);
  });

  // Animar gráficos circulares - ahora se dibujan con canvas
  dibujarGraficosPie();

  // Animar tarjetas de usuario
  const userCards = document.querySelectorAll('.balance-user-card');
  userCards.forEach((card, index) => {
    card.style.opacity = '0';
    card.style.transform = 'translateY(10px)';
    
    setTimeout(() => {
      card.style.transition = 'all 0.4s ease';
      card.style.opacity = '1';
      card.style.transform = 'translateY(0)';
    }, index * 100);
  });
}

// ========================================
// CONFIGURAR SCROLLS
// ========================================

function configurarScrolls() {
  const scrollContainers = document.querySelectorAll('.movements-list, .balance-cards, .grafico-leyenda');
  
  scrollContainers.forEach(container => {
    // Estilo de scrollbar personalizado
    container.style.scrollbarWidth = 'thin';
    container.style.scrollbarColor = '#cbd5e1 #f1f5f9';
  });
}

// ========================================
// ACTUALIZAR STATS ANIMADOS
// ========================================

function animarContador(elemento, valorFinal, duracion = 1000) {
  if (!elemento) return;
  
  const textoOriginal = elemento.textContent;
  const esMoneda = textoOriginal.includes('S/');
  const esNegativo = valorFinal < 0;
  const valorAbsoluto = Math.abs(valorFinal);
  
  let valorActual = 0;
  const incremento = valorAbsoluto / (duracion / 16);
  
  const actualizar = () => {
    valorActual += incremento;
    
    if (valorActual >= valorAbsoluto) {
      const valor = esNegativo ? -valorAbsoluto : valorAbsoluto;
      elemento.textContent = esMoneda ? `S/ ${valor.toFixed(2)}` : valor.toFixed(2);
      return;
    }
    
    const valor = esNegativo ? -valorActual : valorActual;
    elemento.textContent = esMoneda ? `S/ ${valor.toFixed(2)}` : valor.toFixed(2);
    requestAnimationFrame(actualizar);
  };
  
  requestAnimationFrame(actualizar);
}

// ========================================
// GRÁFICOS CIRCULARES (PIE CHART)
// ========================================

function dibujarGraficosPie() {
  const canvases = document.querySelectorAll('canvas[data-chart]');
  
  console.log('Canvases encontrados:', canvases.length);
  
  canvases.forEach((canvas, idx) => {
    try {
      const dataStr = canvas.getAttribute('data-chart');
      console.log(`Canvas ${idx} data:`, dataStr);
      
      const data = JSON.parse(dataStr);
      
      if (!data || data.length === 0) {
        console.log(`Canvas ${idx} sin datos`);
        return;
      }
      
      console.log(`Dibujando gráfico ${idx} con ${data.length} segmentos`);
      
      const ctx = canvas.getContext('2d');
      const centerX = canvas.width / 2;
      const centerY = canvas.height / 2;
      const radius = Math.min(centerX, centerY) - 20;
      
      let currentAngle = -Math.PI / 2; // Empezar desde arriba
      
      // Limpiar canvas
      ctx.clearRect(0, 0, canvas.width, canvas.height);
      
      // Dibujar cada segmento
      data.forEach((item, index) => {
        const sliceAngle = (item.value / 100) * 2 * Math.PI;
        
        // Dibujar segmento
        ctx.beginPath();
        ctx.moveTo(centerX, centerY);
        ctx.arc(centerX, centerY, radius, currentAngle, currentAngle + sliceAngle);
        ctx.closePath();
        
        ctx.fillStyle = item.color;
        ctx.fill();
        
        // Borde entre segmentos
        ctx.strokeStyle = '#ffffff';
        ctx.lineWidth = 3;
        ctx.stroke();
        
        currentAngle += sliceAngle;
      });
      
      // Agregar interactividad
      canvas.addEventListener('mousemove', function(e) {
        const rect = canvas.getBoundingClientRect();
        const x = e.clientX - rect.left - centerX;
        const y = e.clientY - rect.top - centerY;
        const distance = Math.sqrt(x * x + y * y);
        
        if (distance <= radius) {
          let angle = Math.atan2(y, x) + Math.PI / 2;
          if (angle < 0) angle += 2 * Math.PI;
          
          let currentAngle = 0;
          for (let i = 0; i < data.length; i++) {
            const sliceAngle = (data[i].value / 100) * 2 * Math.PI;
            if (angle >= currentAngle && angle < currentAngle + sliceAngle) {
              canvas.style.cursor = 'pointer';
              mostrarTooltipCanvas(e, data[i]);
              return;
            }
            currentAngle += sliceAngle;
          }
        } else {
          canvas.style.cursor = 'default';
          ocultarTooltip();
        }
      });
      
      canvas.addEventListener('mouseleave', function() {
        ocultarTooltip();
      });
      
      console.log(`✅ Gráfico ${idx} dibujado correctamente`);
      
    } catch (error) {
      console.error(`Error dibujando canvas ${idx}:`, error);
    }
  });
}

function mostrarTooltipCanvas(evento, data) {
  // Remover tooltip existente
  ocultarTooltip();
  
  const tooltip = document.createElement('div');
  tooltip.className = 'custom-tooltip';
  tooltip.style.cssText = `
    position: fixed;
    background: rgba(0, 0, 0, 0.95);
    color: white;
    padding: 12px 16px;
    border-radius: 8px;
    font-size: 13px;
    font-weight: 600;
    pointer-events: none;
    z-index: 10000;
    white-space: nowrap;
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.4);
    border: 2px solid ${data.color};
  `;
  
  tooltip.innerHTML = `
    <div style="margin-bottom: 6px; font-weight: 700; font-size: 14px;">${data.label}</div>
    <div style="color: #a3e635; font-size: 15px; font-weight: 800;">S/ ${data.amount.toFixed(2)}</div>
    <div style="color: #cbd5e1; font-size: 12px; margin-top: 4px;">${data.value.toFixed(1)}%</div>
  `;
  
  document.body.appendChild(tooltip);
  
  tooltip.style.left = evento.clientX + 15 + 'px';
  tooltip.style.top = evento.clientY + 15 + 'px';
}

// ========================================
// TOOLTIPS PARA GRÁFICOS
// ========================================

function mostrarTooltipPie(evento, elemento) {
  const label = elemento.dataset.label;
  const percentage = elemento.dataset.percentage;
  const amount = elemento.dataset.amount;
  
  const tooltip = document.createElement('div');
  tooltip.className = 'custom-tooltip';
  tooltip.style.cssText = `
    position: fixed;
    background: rgba(0, 0, 0, 0.9);
    color: white;
    padding: 12px 16px;
    border-radius: 8px;
    font-size: 13px;
    font-weight: 600;
    pointer-events: none;
    z-index: 10000;
    white-space: nowrap;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  `;
  
  tooltip.innerHTML = `
    <div style="margin-bottom: 4px; font-weight: 700;">${label}</div>
    <div style="color: #a3e635;">S/ ${amount}</div>
    <div style="color: #94a3b8; font-size: 12px;">${percentage}%</div>
  `;
  
  document.body.appendChild(tooltip);
  
  const rect = evento.target.getBoundingClientRect();
  tooltip.style.left = evento.clientX + 10 + 'px';
  tooltip.style.top = evento.clientY + 10 + 'px';
}

function agregarTooltips() {
  const bars = document.querySelectorAll('.bar');
  
  bars.forEach(elemento => {
    elemento.addEventListener('mouseenter', function(e) {
      const container = this.closest('.bar-container');
      const label = container.querySelector('.bar-label').textContent;
      const value = this.querySelector('.bar-value').textContent;
      
      const tooltip = document.createElement('div');
      tooltip.className = 'custom-tooltip';
      tooltip.style.cssText = `
        position: fixed;
        background: rgba(0, 0, 0, 0.9);
        color: white;
        padding: 8px 12px;
        border-radius: 6px;
        font-size: 13px;
        font-weight: 600;
        pointer-events: none;
        z-index: 10000;
        white-space: nowrap;
      `;
      
      tooltip.textContent = `${label}: ${value}`;
      document.body.appendChild(tooltip);
      
      const rect = this.getBoundingClientRect();
      tooltip.style.left = rect.left + (rect.width / 2) - (tooltip.offsetWidth / 2) + 'px';
      tooltip.style.top = rect.top - tooltip.offsetHeight - 10 + 'px';
    });
    
    elemento.addEventListener('mouseleave', function() {
      ocultarTooltip();
    });
  });
}

function ocultarTooltip() {
  const tooltip = document.querySelector('.custom-tooltip');
  if (tooltip) {
    tooltip.remove();
  }
}

// ========================================
// EXPORT FUNCIONES GLOBALES
// ========================================

window.balanceUtils = {
  filtrarUsuario,
  cambiarTipoGrafico,
  cambiarVistaGrafico,
  animarContador
};

// Log de carga
console.log('%c[Balance] Módulo cargado correctamente', 'color: #22c55e; font-weight: bold');