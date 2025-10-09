<?php
// frontend/views/layouts/default.php - Layout reutilizable (inyecta $content, $title, $data)
?>
<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title><?php echo htmlspecialchars($title ?? 'Finanzas App'); ?></title>
    <?php if (method_exists($this, 'addAsset')): ?>
        <?php $this->addAsset('css', 'styles'); ?> <!-- Llama helper de BaseController para CSS -->
    <?php endif; ?>
</head>
<body>
    <header>
        <nav style="background: #f8f9fa; padding: 10px;">
            <a href="<?php echo URLROOT; ?>/dashboard">Dashboard</a> |
            Bienvenido, <?php echo htmlspecialchars($_SESSION['user_email'] ?? 'Usuario'); ?> |
            <a href="<?php echo URLROOT; ?>/logout">Cerrar Sesión</a>
        </nav>
    </header>
    <main style="padding: 20px;">
        <?php echo $content; ?> <!-- Inserta la vista principal (e.g., dashboard/index.php) -->
    </main>
    <footer style="background: #f8f9fa; padding: 10px; text-align: center;">
        &copy; 2025 Finanzas App
    </footer>
    <?php if (method_exists($this, 'addAsset')): ?>
        <?php $this->addAsset('js', 'app'); ?> <!-- JS para interacciones -->
    <?php endif; ?>
</body>
</html>