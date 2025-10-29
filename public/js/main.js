


// === Mostrar panel lateral al hacer click en un concepto ===
//const conceptos = document.querySelectorAll('.concepto');
const sidebarRight = document.querySelector('.sidebar-right');



// === Notificaciones ===
const userRole = "admin";
const notifications = [
  { type: "own", icon: "fa-shirt", message: "Has llegado al 90% del límite de gasto en Ropa", time: "hace 2h" },
  { type: "other", user: "Sofía", message: "ha alcanzado el límite de gasto en Entretenimiento", time: "hace 1h" },
  { type: "own", icon: "fa-basket-shopping", message: "Has superado el límite mensual en Mercado", time: "hace 3h" }
];

const bellIcon = document.querySelector('.fa-bell');
const notificationsPanel = document.querySelector('.notifications-panel');
const notificationCount = document.querySelector('.notification-count');

if (bellIcon) {
  bellIcon.style.cursor = "pointer";
}

const filteredNotifications = userRole === "admin"
  ? notifications
  : notifications.filter(n => n.type === "own");

if (notificationCount) {
  notificationCount.textContent = filteredNotifications.length;
}

if (notificationsPanel) {
  notificationsPanel.innerHTML = filteredNotifications.map(n => {
    if (n.type === "own") {
      return `
        <div class="notification-item">
          <div class="notification-icon" style="background:#8CCB5E">
            <i class="fa-solid ${n.icon}"></i>
          </div>
          <div class="notification-info">
            <span>${n.message}</span>
            <span class="notification-time">${n.time}</span>
          </div>
        </div>
      `;
    } else {
      const initial = n.user.charAt(0).toUpperCase();
      return `
        <div class="notification-item">
          <div class="notification-avatar">${initial}</div>
          <div class="notification-info">
            <span><strong>${n.user}</strong> ${n.message}</span>
            <span class="notification-time">${n.time}</span>
          </div>
        </div>
      `;
    }
  }).join('');
}

if (bellIcon && notificationsPanel) {
  bellIcon.addEventListener('click', (e) => {
    e.stopPropagation();
    notificationsPanel.classList.toggle('active');
    bellIcon.classList.toggle('active');
    if (bellIcon.classList.contains('active')) {
      bellIcon.style.color = "#FFD84D";
    } else {
      bellIcon.style.color = "";
    }
  });
}

window.addEventListener('click', (e) => {
  if (notificationsPanel && bellIcon) {
    if (!notificationsPanel.contains(e.target) && e.target !== bellIcon) {
      notificationsPanel.classList.remove('active');
      bellIcon.classList.remove('active');
      bellIcon.style.color = "";
    }
  }
});