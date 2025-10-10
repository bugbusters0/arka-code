<div class="page-header">
  <h1>Gestión de Conceptos</h1>
  <a href="<?php echo URLROOT; ?>/concepto/crear" class="btn btn-primary">+ Nuevo Concepto</a>
</div>

<div class="filters-section">
  <form method="GET" action="<?php echo URLROOT; ?>/concepto/listar">
    <div class="filters-grid">
      <input type="text" name="search" placeholder="Buscar concepto..."
        value="<?php echo htmlspecialchars($_GET['search'] ?? ''); ?>">

      <select name="tipo">
        <option value="">Todos los tipos</option>
        <option value="ingreso" <?php echo ($_GET['tipo'] ?? '') === 'ingreso' ? 'selected' : ''; ?>>
          Ingresos
        </option>
        <option value="gasto" <?php echo ($_GET['tipo'] ?? '') === 'gasto' ? 'selected' : ''; ?>>
          Gastos
        </option>
      </select>

      <button type="submit" class="btn btn-primary">Filtrar</button>
    </div>
  </form>
</div>

<div class="conceptos-grid">
  <?php if (!empty($conceptos)): ?>
    <?php foreach ($conceptos as $concepto): ?>
      <div class="concepto-card" style="border-left: 4px solid <?php echo htmlspecialchars($concepto['color'] ?? '#667eea'); ?>">
        <div class="concepto-header">
          <?php if (!empty($concepto['icono_path'])): ?>
            <img src="<?php echo URLROOT . '/' . htmlspecialchars($concepto['icono_path']); ?>"
              alt="<?php echo htmlspecialchars($concepto['nombre']); ?>"
              class="concepto-icon">
          <?php else: ?>
            <span class="concepto-icon-placeholder">💰</span>
          <?php endif; ?>

          <div class="concepto-info">
            <h3><?php echo htmlspecialchars($concepto['nombre']); ?></h3>
            <span class="badge badge-<?php echo $concepto['tipo']; ?>">
              <?php echo ucfirst($concepto['tipo']); ?>
            </span>
          </div>
        </div>

        <div class="concepto-stats">
          <div class="stat">
            <span class="stat-label">Usuarios</span>
            <span class="stat-value"><?php echo $concepto['usuarios_count'] ?? 0; ?></span>
          </div>
          <div class="stat">
            <span class="stat-label">Movimientos</span>
            <span class="stat-value"><?php echo $concepto['movimientos_count'] ?? 0; ?></span>
          </div>
        </div>

        <div class="concepto-actions">
          <a href="<?php echo URLROOT; ?>/concepto/editar/<?php echo $concepto['id']; ?>"
            class="btn-action btn-edit">✏️ Editar</a>
          <a href="<?php echo URLROOT; ?>/concepto/eliminar/<?php echo $concepto['id']; ?>"
            class="btn-action btn-delete"
            onclick="return confirm('¿Eliminar este concepto?')">🗑️ Eliminar</a>
        </div>
      </div>
    <?php endforeach; ?>
  <?php else: ?>
    <div class="empty-state">
      <p style="font-size: 3rem;">📋</p>
      <h3>No hay conceptos registrados</h3>
      <p>Crea tu primer concepto para comenzar a gestionar tus finanzas</p>
      <a href="<?php echo URLROOT; ?>/concepto/crear" class="btn btn-primary">Crear Concepto</a>
    </div>
  <?php endif; ?>
</div>

<style>
  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
  }

  .page-header h1 {
    font-size: 2rem;
    color: #2c3e50;
  }

  .filters-section {
    background: white;
    padding: 1.5rem;
    border-radius: 10px;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
    margin-bottom: 2rem;
  }

  .filters-grid {
    display: grid;
    grid-template-columns: 1fr 1fr auto;
    gap: 1rem;
  }

  .filters-grid input,
  .filters-grid select {
    padding: 0.75rem;
    border: 1px solid #ddd;
    border-radius: 5px;
    font-size: 1rem;
  }

  .conceptos-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 1.5rem;
  }

  .concepto-card {
    background: white;
    border-radius: 10px;
    padding: 1.5rem;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
    transition: transform 0.3s, box-shadow 0.3s;
  }

  .concepto-card:hover {
    transform: translateY(-5px);
    box-shadow: 0 5px 20px rgba(0, 0, 0, 0.1);
  }

  .concepto-header {
    display: flex;
    gap: 1rem;
    margin-bottom: 1rem;
  }

  .concepto-icon {
    width: 50px;
    height: 50px;
    border-radius: 50%;
    object-fit: cover;
  }

  .concepto-icon-placeholder {
    width: 50px;
    height: 50px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #f0f0f0;
    border-radius: 50%;
    font-size: 1.5rem;
  }

  .concepto-info h3 {
    font-size: 1.2rem;
    color: #2c3e50;
    margin-bottom: 0.5rem;
  }

  .badge {
    padding: 0.25rem 0.75rem;
    border-radius: 20px;
    font-size: 0.85rem;
    font-weight: bold;
  }

  .badge-ingreso {
    background: #d4edda;
    color: #155724;
  }

  .badge-gasto {
    background: #f8d7da;
    color: #721c24;
  }

  .concepto-stats {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
    margin: 1rem 0;
    padding: 1rem 0;
    border-top: 1px solid #f0f0f0;
    border-bottom: 1px solid #f0f0f0;
  }

  .stat {
    text-align: center;
  }

  .stat-label {
    display: block;
    font-size: 0.85rem;
    color: #7f8c8d;
    margin-bottom: 0.25rem;
  }

  .stat-value {
    display: block;
    font-size: 1.5rem;
    font-weight: bold;
    color: #2c3e50;
  }

  .concepto-actions {
    display: flex;
    gap: 0.5rem;
  }

  .btn-action {
    flex: 1;
    padding: 0.5rem;
    border-radius: 5px;
    text-align: center;
    text-decoration: none;
    font-size: 0.9rem;
    transition: all 0.3s;
  }

  .btn-edit {
    background: #3498db;
    color: white;
  }

  .btn-edit:hover {
    background: #2980b9;
  }

  .btn-delete {
    background: #e74c3c;
    color: white;
  }

  .btn-delete:hover {
    background: #c0392b;
  }

  .empty-state {
    grid-column: 1 / -1;
    text-align: center;
    padding: 4rem 2rem;
    background: white;
    border-radius: 10px;
  }

  .empty-state h3 {
    color: #7f8c8d;
    margin: 1rem 0;
  }

  .empty-state p {
    color: #95a5a6;
    margin-bottom: 2rem;
  }
</style>