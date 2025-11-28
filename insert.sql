INSERT INTO `familia` (`correo`, `telefono`, `contraseña`) VALUES
('familia.garcia@email.com', '910234567', '$2a$14$7xGgcRdhPhZPKkcRufVWSOjM1DORG9FZtr5ISaKFh6B9DF58DZV6S'),
('los.rodriguez@email.com', '920987654', '$2a$14$7xGgcRdhPhZPKkcRufVWSOjM1DORG9FZtr5ISaKFh6B9DF58DZV6S'),
('perez.home@email.com', '933112233', '$2a$14$7xGgcRdhPhZPKkcRufVWSOjM1DORG9FZtr5ISaKFh6B9DF58DZV6S'),
('fernandez.fam@email.com', '944445566', '$2a$14$7xGgcRdhPhZPKkcRufVWSOjM1DORG9FZtr5ISaKFh6B9DF58DZV6S'),
('martinez.team@email.com', '955778899', '$2a$14$7xGgcRdhPhZPKkcRufVWSOjM1DORG9FZtr5ISaKFh6B9DF58DZV6S'),
('sanchez.clan@email.com', '966001122', '$2a$14$7xGgcRdhPhZPKkcRufVWSOjM1DORG9FZtr5ISaKFh6B9DF58DZV6S'),
('gomez.group@email.com', '977334455', '$2a$14$7xGgcRdhPhZPKkcRufVWSOjM1DORG9FZtr5ISaKFh6B9DF58DZV6S'),
('diaz.family@email.com', '988667788', '$2a$14$7xGgcRdhPhZPKkcRufVWSOjM1DORG9FZtr5ISaKFh6B9DF58DZV6S'),
('alvarez.house@email.com', '999900111', '$2a$14$7xGgcRdhPhZPKkcRufVWSOjM1DORG9FZtr5ISaKFh6B9DF58DZV6S'),
('moreno.co@email.com', '900223344', '$2a$14$7xGgcRdhPhZPKkcRufVWSOjM1DORG9FZtr5ISaKFh6B9DF58DZV6S');

INSERT INTO `usuario` (`nombreUsuario`, `rol`, `contraseñaPersonal`, `nombrePersonal`, `correoFamilia`) VALUES
/* --- Familia: familia.garcia@email.com --- */
('jefe_garcia', 1, '$2a$14$ftf.caI1Gk3FC7lbEvXvM.rdUdkk6me6RVLwHzbQhNukSL/yn.6/u', 'Javier García (Padre)', 'familia.garcia@email.com'),
('eva_g', 0, '$2a$14$ftf.caI1Gk3FC7lbEvXvM.rdUdkk6me6RVLwHzbQhNukSL/yn.6/u', 'Eva López (Madre)', 'familia.garcia@email.com'),
('hugo_g', 0, '$2a$14$ftf.caI1Gk3FC7lbEvXvM.rdUdkk6me6RVLwHzbQhNukSL/yn.6/u', 'Hugo García (Hijo)', 'familia.garcia@email.com'),
('laura_g', 0, '$2a$14$ftf.caI1Gk3FC7lbEvXvM.rdUdkk6me6RVLwHzbQhNukSL/yn.6/u', 'Laura García (Hija)', 'familia.garcia@email.com'),

/* --- Familia: los.rodriguez@email.com --- */
('admin_rodri', 1, '$2a$14$ftf.caI1Gk3FC7lbEvXvM.rdUdkk6me6RVLwHzbQhNukSL/yn.6/u', 'Miguel Rodríguez (Padre)', 'los.rodriguez@email.com'),
('ana_r', 0, '$2a$14$ftf.caI1Gk3FC7lbEvXvM.rdUdkk6me6RVLwHzbQhNukSL/yn.6/u', 'Ana Torres (Madre)', 'los.rodriguez@email.com'),
('pablo_r', 0, '$2a$14$ftf.caI1Gk3FC7lbEvXvM.rdUdkk6me6RVLwHzbQhNukSL/yn.6/u', 'Pablo Rodríguez (Hijo)', 'los.rodriguez@email.com'),
('sofia_r', 0, '$2a$14$ftf.caI1Gk3FC7lbEvXvM.rdUdkk6me6RVLwHzbQhNukSL/yn.6/u', 'Sofía Rodríguez (Hija)', 'los.rodriguez@email.com'),

/* --- Familia: perez.home@email.com --- */
('jefe_perez', 1, '$2a$14$ftf.caI1Gk3FC7lbEvXvM.rdUdkk6me6RVLwHzbQhNukSL/yn.6/u', 'Mario Pérez (Padre)', 'perez.home@email.com'),
('clara_p', 0, '$2a$14$ftf.caI1Gk3FC7lbEvXvM.rdUdkk6me6RVLwHzbQhNukSL/yn.6/u', 'Clara Ruiz (Madre)', 'perez.home@email.com'),
('diego_p', 0, '$2a$14$ftf.caI1Gk3FC7lbEvXvM.rdUdkk6me6RVLwHzbQhNukSL/yn.6/u', 'Diego Pérez (Hijo)', 'perez.home@email.com'),
('irene_p', 0, '$2a$14$ftf.caI1Gk3FC7lbEvXvM.rdUdkk6me6RVLwHzbQhNukSL/yn.6/u', 'Irene Pérez (Hija)', 'perez.home@email.com'),

/* --- Familia: fernandez.fam@email.com --- */
('admin_fer', 1, '$2a$14$ftf.caI1Gk3FC7lbEvXvM.rdUdkk6me6RVLwHzbQhNukSL/yn.6/u', 'Luis Fernández (Padre)', 'fernandez.fam@email.com'),
('elena_f', 0, '$2a$14$ftf.caI1Gk3FC7lbEvXvM.rdUdkk6me6RVLwHzbQhNukSL/yn.6/u', 'Elena Soto (Madre)', 'fernandez.fam@email.com'),
('sergio_f', 0, '$2a$14$ftf.caI1Gk3FC7lbEvXvM.rdUdkk6me6RVLwHzbQhNukSL/yn.6/u', 'Sergio Fernández (Hijo)', 'fernandez.fam@email.com'),
('cris_f', 0, '$2a$14$ftf.caI1Gk3FC7lbEvXvM.rdUdkk6me6RVLwHzbQhNukSL/yn.6/u', 'Cristina Fernández (Hija)', 'fernandez.fam@email.com'),

/* --- Familia: martinez.team@email.com --- */
('jefe_martinez', 1, '$2a$14$ftf.caI1Gk3FC7lbEvXvM.rdUdkk6me6RVLwHzbQhNukSL/yn.6/u', 'Carlos Martínez (Padre)', 'martinez.team@email.com'),
('isabel_m', 0, '$2a$14$ftf.caI1Gk3FC7lbEvXvM.rdUdkk6me6RVLwHzbQhNukSL/yn.6/u', 'Isabel Blanco (Madre)', 'martinez.team@email.com'),
('raul_m', 0, '$2a$14$ftf.caI1Gk3FC7lbEvXvM.rdUdkk6me6RVLwHzbQhNukSL/yn.6/u', 'Raúl Martínez (Hijo)', 'martinez.team@email.com'),
('nuria_m', 0, '$2a$14$ftf.caI1Gk3FC7lbEvXvM.rdUdkk6me6RVLwHzbQhNukSL/yn.6/u', 'Nuria Martínez (Hija)', 'martinez.team@email.com'),

/* --- Familia: sanchez.clan@email.com --- */
('admin_sanchez', 1, '$2a$14$ftf.caI1Gk3FC7lbEvXvM.rdUdkk6me6RVLwHzbQhNukSL/yn.6/u', 'Alberto Sánchez (Padre)', 'sanchez.clan@email.com'),
('patricia_s', 0, '$2a$14$ftf.caI1Gk3FC7lbEvXvM.rdUdkk6me6RVLwHzbQhNukSL/yn.6/u', 'Patricia Vidal (Madre)', 'sanchez.clan@email.com'),
('axel_s', 0, '$2a$14$ftf.caI1Gk3FC7lbEvXvM.rdUdkk6me6RVLwHzbQhNukSL/yn.6/u', 'Axel Sánchez (Hijo)', 'sanchez.clan@email.com'),
('lidia_s', 0, '$2a$14$ftf.caI1Gk3FC7lbEvXvM.rdUdkk6me6RVLwHzbQhNukSL/yn.6/u', 'Lidia Sánchez (Hija)', 'sanchez.clan@email.com'),

/* --- Familia: gomez.group@email.com --- */
('jefe_gomez', 1, '$2a$14$ftf.caI1Gk3FC7lbEvXvM.rdUdkk6me6RVLwHzbQhNukSL/yn.6/u', 'Fernando Gómez (Padre)', 'gomez.group@email.com'),
('rocio_g', 0, '$2a$14$ftf.caI1Gk3FC7lbEvXvM.rdUdkk6me6RVLwHzbQhNukSL/yn.6/u', 'Rocío Castro (Madre)', 'gomez.group@email.com'),
('oscar_g', 0, '$2a$14$ftf.caI1Gk3FC7lbEvXvM.rdUdkk6me6RVLwHzbQhNukSL/yn.6/u', 'Óscar Gómez (Hijo)', 'gomez.group@email.com'),
('alba_g', 0, '$2a$14$ftf.caI1Gk3FC7lbEvXvM.rdUdkk6me6RVLwHzbQhNukSL/yn.6/u', 'Alba Gómez (Hija)', 'gomez.group@email.com'),

/* --- Familia: diaz.family@email.com --- */
('admin_diaz', 1, '$2a$14$ftf.caI1Gk3FC7lbEvXvM.rdUdkk6me6RVLwHzbQhNukSL/yn.6/u', 'Roberto Díaz (Padre)', 'diaz.family@email.com'),
('mercedes_d', 0, '$2a$14$ftf.caI1Gk3FC7lbEvXvM.rdUdkk6me6RVLwHzbQhNukSL/yn.6/u', 'Mercedes Gil (Madre)', 'diaz.family@email.com'),
('victor_d', 0, '$2a$14$ftf.caI1Gk3FC7lbEvXvM.rdUdkk6me6RVLwHzbQhNukSL/yn.6/u', 'Víctor Díaz (Hijo)', 'diaz.family@email.com'),
('andrea_d', 0, '$2a$14$ftf.caI1Gk3FC7lbEvXvM.rdUdkk6me6RVLwHzbQhNukSL/yn.6/u', 'Andrea Díaz (Hija)', 'diaz.family@email.com'),

/* --- Familia: alvarez.house@email.com --- */
('jefe_alvarez', 1, '$2a$14$ftf.caI1Gk3FC7lbEvXvM.rdUdkk6me6RVLwHzbQhNukSL/yn.6/u', 'José Álvarez (Padre)', 'alvarez.house@email.com'),
('rosa_a', 0, '$2a$14$ftf.caI1Gk3FC7lbEvXvM.rdUdkk6me6RVLwHzbQhNukSL/yn.6/u', 'Rosa Márquez (Madre)', 'alvarez.house@email.com'),
('dani_a', 0, '$2a$14$ftf.caI1Gk3FC7lbEvXvM.rdUdkk6me6RVLwHzbQhNukSL/yn.6/u', 'Daniel Álvarez (Hijo)', 'alvarez.house@email.com'),
('eva_a', 0, '$2a$14$ftf.caI1Gk3FC7lbEvXvM.rdUdkk6me6RVLwHzbQhNukSL/yn.6/u', 'Eva Álvarez (Hija)', 'alvarez.house@email.com'),

/* --- Familia: moreno.co@email.com --- */
('admin_moreno', 1, '$2a$14$ftf.caI1Gk3FC7lbEvXvM.rdUdkk6me6RVLwHzbQhNukSL/yn.6/u', 'Antonio Moreno (Padre)', 'moreno.co@email.com'),
('marta_m', 0, '$2a$14$ftf.caI1Gk3FC7lbEvXvM.rdUdkk6me6RVLwHzbQhNukSL/yn.6/u', 'Marta Torres (Madre)', 'moreno.co@email.com'),
('felix_m', 0, '$2a$14$ftf.caI1Gk3FC7lbEvXvM.rdUdkk6me6RVLwHzbQhNukSL/yn.6/u', 'Félix Moreno (Hijo)', 'moreno.co@email.com'),
('nerea_m', 0, '$2a$14$ftf.caI1Gk3FC7lbEvXvM.rdUdkk6me6RVLwHzbQhNukSL/yn.6/u', 'Nerea Moreno (Hija)', 'moreno.co@email.com');



INSERT INTO `concepto` (`nombreConcepto`, `correoFamilia`, `tipo`, `icono`, `color`, `nombreUsuario`) VALUES
/* --- Familia: familia.garcia@email.com (6 conceptos/usuario) --- */
/* Usuario: jefe_garcia (Admin) */
('Alquiler/Hipoteca_G', 'familia.garcia@email.com', 0, 'fa-solid fa-building', '#FF6B6B', 'jefe_garcia'),
('Supermercado_G', 'familia.garcia@email.com', 0, 'fa-solid fa-utensils', '#FFD84D', 'jefe_garcia'),
('Electricidad_G', 'familia.garcia@email.com', 0, 'fa-solid fa-plug', '#4DA3FF', 'jefe_garcia'),
('Salario_Ppal_I', 'familia.garcia@email.com', 1, 'fa-solid fa-money-bill-wave', '#8CCB5E', 'jefe_garcia'),
('Inversiones_I', 'familia.garcia@email.com', 1, 'fa-solid fa-chart-line', '#A66BFF', 'jefe_garcia'),
('Regalos_Rec_I', 'familia.garcia@email.com', 1, 'fa-solid fa-gift', '#FFD84D', 'jefe_garcia'),
/* Usuario: eva_g */
('Ropa_eva_G', 'familia.garcia@email.com', 0, 'fa-solid fa-shirt', '#A66BFF', 'eva_g'),
('Salud_eva_G', 'familia.garcia@email.com', 0, 'fa-solid fa-heart-pulse', '#FF6B6B', 'eva_g'),
('Taxi_eva_G', 'familia.garcia@email.com', 0, 'fa-solid fa-taxi', '#4DA3FF', 'eva_g'),
('Bonos_eva_I', 'familia.garcia@email.com', 1, 'fa-solid fa-star', '#A66BFF', 'eva_g'),
('Freelance_eva_I', 'familia.garcia@email.com', 1, 'fa-solid fa-laptop-code', '#8CCB5E', 'eva_g'),
('Comisiones_eva_I', 'familia.garcia@email.com', 1, 'fa-solid fa-percent', '#FFD84D', 'eva_g'),
/* Usuario: hugo_g */
('Videojuegos_G', 'familia.garcia@email.com', 0, 'fa-solid fa-laptop', '#FFD84D', 'hugo_g'),
('Gimnasio_hugo_G', 'familia.garcia@email.com', 0, 'fa-solid fa-weight-hanging', '#A66BFF', 'hugo_g'),
('Comida_rapida_G', 'familia.garcia@email.com', 0, 'fa-solid fa-utensils', '#8CCB5E', 'hugo_g'),
('Paga_hugo_I', 'familia.garcia@email.com', 1, 'fa-solid fa-money-bill-wave', '#8CCB5E', 'hugo_g'),
('Trabajo_Verano_I', 'familia.garcia@email.com', 1, 'fa-solid fa-laptop-code', '#FFD84D', 'hugo_g'),
('Regalo_Abuelos_I', 'familia.garcia@email.com', 1, 'fa-solid fa-gift', '#A66BFF', 'hugo_g'),
/* Usuario: laura_g */
('Cursos_laura_G', 'familia.garcia@email.com', 0, 'fa-solid fa-graduation-cap', '#4DA3FF', 'laura_g'),
('Viaje_laura_G', 'familia.garcia@email.com', 0, 'fa-solid fa-plane', '#8CCB5E', 'laura_g'),
('Deportes_laura_G', 'familia.garcia@email.com', 0, 'fa-solid fa-dumbbell', '#FF6B6B', 'laura_g'),
('Beca_laura_I', 'familia.garcia@email.com', 1, 'fa-solid fa-money-bill-wave', '#FFD84D', 'laura_g'),
('Prestamo_Rec_I', 'familia.garcia@email.com', 1, 'fa-solid fa-piggy-bank', '#4DA3FF', 'laura_g'),
('Premio_I', 'familia.garcia@email.com', 1, 'fa-solid fa-star', '#A66BFF', 'laura_g'),

/* --- Familia: los.rodriguez@email.com --- */
/* Usuario: admin_rodri (Admin) */
('Seguro_admin_G', 'los.rodriguez@email.com', 0, 'fa-solid fa-house', '#FF6B6B', 'admin_rodri'),
('Comida_admin_G', 'los.rodriguez@email.com', 0, 'fa-solid fa-utensils', '#FFD84D', 'admin_rodri'),
('Transporte_admin_G', 'los.rodriguez@email.com', 0, 'fa-solid fa-car', '#4DA3FF', 'admin_rodri'),
('Salario_Sec_I', 'los.rodriguez@email.com', 1, 'fa-solid fa-money-bill-wave', '#8CCB5E', 'admin_rodri'),
('Intereses_I', 'los.rodriguez@email.com', 1, 'fa-solid fa-chart-line', '#4DA3FF', 'admin_rodri'),
('Bonificacion_I', 'los.rodriguez@email.com', 1, 'fa-solid fa-star', '#FFD84D', 'admin_rodri'),
/* Usuario: ana_r */
('Gimnasio_ana_G', 'los.rodriguez@email.com', 0, 'fa-solid fa-weight-hanging', '#A66BFF', 'ana_r'),
('Otros_ana_G', 'los.rodriguez@email.com', 0, 'fa-solid fa-circle', '#4DA3FF', 'ana_r'),
('Salud_ana_G', 'los.rodriguez@email.com', 0, 'fa-solid fa-heart-pulse', '#8CCB5E', 'ana_r'),
('Freelance_ana_I', 'los.rodriguez@email.com', 1, 'fa-solid fa-laptop-code', '#A66BFF', 'ana_r'),
('Regalo_Cumple_I', 'los.rodriguez@email.com', 1, 'fa-solid fa-gift', '#8CCB5E', 'ana_r'),
('Comisiones_ana_I', 'los.rodriguez@email.com', 1, 'fa-solid fa-percent', '#FF6B6B', 'ana_r'),
/* Usuario: pablo_r */
('Cine_pablo_G', 'los.rodriguez@email.com', 0, 'fa-solid fa-clapperboard', '#FFD84D', 'pablo_r'),
('Tecnologia_pablo_G', 'los.rodriguez@email.com', 0, 'fa-solid fa-laptop', '#A66BFF', 'pablo_r'),
('Deportes_pablo_G', 'los.rodriguez@email.com', 0, 'fa-solid fa-dumbbell', '#FF6B6B', 'pablo_r'),
('Paga_pablo_I', 'los.rodriguez@email.com', 1, 'fa-solid fa-money-bill-wave', '#8CCB5E', 'pablo_r'),
('Regalo_Navidad_I', 'los.rodriguez@email.com', 1, 'fa-solid fa-gift', '#A66BFF', 'pablo_r'),
('Ahorro_Vacaciones_I', 'los.rodriguez@email.com', 1, 'fa-solid fa-piggy-bank', '#FF6B6B', 'pablo_r'),
/* Usuario: sofia_r */
('Ropa_sofia_G', 'los.rodriguez@email.com', 0, 'fa-solid fa-shirt', '#4DA3FF', 'sofia_r'),
('Cosmeticos_G', 'los.rodriguez@email.com', 0, 'fa-solid fa-pills', '#8CCB5E', 'sofia_r'),
('Cafe_G', 'los.rodriguez@email.com', 0, 'fa-solid fa-utensils', '#A66BFF', 'sofia_r'),
('Beca_sofia_I', 'los.rodriguez@email.com', 1, 'fa-solid fa-money-bill-wave', '#FFD84D', 'sofia_r'),
('Regalo_Tia_I', 'los.rodriguez@email.com', 1, 'fa-solid fa-gift', '#4DA3FF', 'sofia_r'),
('Bonos_Estudio_I', 'los.rodriguez@email.com', 1, 'fa-solid fa-star', '#FF6B6B', 'sofia_r'),

/* --- Familia: perez.home@email.com --- */
/* Usuario: jefe_perez (Admin) */
('Gas_jefe_G', 'perez.home@email.com', 0, 'fa-solid fa-fire-flame-simple', '#FF6B6B', 'jefe_perez'),
('Mantenimiento_jefe_G', 'perez.home@email.com', 0, 'fa-solid fa-hammer', '#FFD84D', 'jefe_perez'),
('Impuestos_G', 'perez.home@email.com', 0, 'fa-solid fa-house', '#4DA3FF', 'jefe_perez'),
('Salario_jefe_I', 'perez.home@email.com', 1, 'fa-solid fa-money-bill-wave', '#8CCB5E', 'jefe_perez'),
('Inversion_jefe_I', 'perez.home@email.com', 1, 'fa-solid fa-chart-line', '#4DA3FF', 'jefe_perez'),
('Dividendos_jefe_I', 'perez.home@email.com', 1, 'fa-solid fa-chart-bar', '#FFD84D', 'jefe_perez'),
/* Usuario: clara_p */
('Comida_clara_G', 'perez.home@email.com', 0, 'fa-solid fa-utensils', '#A66BFF', 'clara_p'),
('Salud_clara_G', 'perez.home@email.com', 0, 'fa-solid fa-heart-pulse', '#FF6B6B', 'clara_p'),
('Entretenimiento_clara_G', 'perez.home@email.com', 0, 'fa-solid fa-film', '#8CCB5E', 'clara_p'),
('Freelance_clara_I', 'perez.home@email.com', 1, 'fa-solid fa-laptop-code', '#A66BFF', 'clara_p'),
('Regalo_clara_I', 'perez.home@email.com', 1, 'fa-solid fa-gift', '#8CCB5E', 'clara_p'),
('Bonos_clara_I', 'perez.home@email.com', 1, 'fa-solid fa-star', '#FF6B6B', 'clara_p'),
/* Usuario: diego_p */
('Libros_diego_G', 'perez.home@email.com', 0, 'fa-solid fa-book-open', '#FFD84D', 'diego_p'),
('Gimnasio_diego_G', 'perez.home@email.com', 0, 'fa-solid fa-weight-hanging', '#A66BFF', 'diego_p'),
('Tecnologia_diego_G', 'perez.home@email.com', 0, 'fa-solid fa-laptop', '#FF6B6B', 'diego_p'),
('Paga_diego_I', 'perez.home@email.com', 1, 'fa-solid fa-money-bill-wave', '#8CCB5E', 'diego_p'),
('Regalo_diego_I', 'perez.home@email.com', 1, 'fa-solid fa-gift', '#FFD84D', 'diego_p'),
('Ahorro_diego_I', 'perez.home@email.com', 1, 'fa-solid fa-piggy-bank', '#A66BFF', 'diego_p'),
/* Usuario: irene_p */
('Viajes_irene_G', 'perez.home@email.com', 0, 'fa-solid fa-plane', '#4DA3FF', 'irene_p'),
('Ropa_irene_G', 'perez.home@email.com', 0, 'fa-solid fa-shirt', '#8CCB5E', 'irene_p'),
('Cine_irene_G', 'perez.home@email.com', 0, 'fa-solid fa-clapperboard', '#FF6B6B', 'irene_p'),
('Beca_irene_I', 'perez.home@email.com', 1, 'fa-solid fa-money-bill-wave', '#FFD84D', 'irene_p'),
('Regalo_tio_I', 'perez.home@email.com', 1, 'fa-solid fa-gift', '#4DA3FF', 'irene_p'),
('Bonos_irene_I', 'perez.home@email.com', 1, 'fa-solid fa-star', '#8CCB5E', 'irene_p'),

/* --- Familia: fernandez.fam@email.com --- */
/* Usuario: admin_fer (Admin) */
('Agua_admin_G', 'fernandez.fam@email.com', 0, 'fa-solid fa-droplet', '#FF6B6B', 'admin_fer'),
('Internet_admin_G', 'fernandez.fam@email.com', 0, 'fa-solid fa-wifi', '#FFD84D', 'admin_fer'),
('Farmacia_admin_G', 'fernandez.fam@email.com', 0, 'fa-solid fa-pills', '#4DA3FF', 'admin_fer'),
('Salario_Admin_I', 'fernandez.fam@email.com', 1, 'fa-solid fa-money-bill-wave', '#8CCB5E', 'admin_fer'),
('Inversion_Bolsa_I', 'fernandez.fam@email.com', 1, 'fa-solid fa-chart-line', '#4DA3FF', 'admin_fer'),
('Alquiler_Prop_I', 'fernandez.fam@email.com', 1, 'fa-solid fa-building', '#FFD84D', 'admin_fer'),
/* Usuario: elena_f */
('Cursos_elena_G', 'fernandez.fam@email.com', 0, 'fa-solid fa-graduation-cap', '#A66BFF', 'elena_f'),
('Ropa_elena_G', 'fernandez.fam@email.com', 0, 'fa-solid fa-shirt', '#FF6B6B', 'elena_f'),
('Viaje_finde_G', 'fernandez.fam@email.com', 0, 'fa-solid fa-plane', '#4DA3FF', 'elena_f'),
('Freelance_elena_I', 'fernandez.fam@email.com', 1, 'fa-solid fa-laptop-code', '#A66BFF', 'elena_f'),
('Regalo_mama_I', 'fernandez.fam@email.com', 1, 'fa-solid fa-gift', '#8CCB5E', 'elena_f'),
('Bonos_elena_I', 'fernandez.fam@email.com', 1, 'fa-solid fa-star', '#FF6B6B', 'elena_f'),
/* Usuario: sergio_f */
('Deportes_sergio_G', 'fernandez.fam@email.com', 0, 'fa-solid fa-dumbbell', '#FFD84D', 'sergio_f'),
('Videojuegos_sergio_G', 'fernandez.fam@email.com', 0, 'fa-solid fa-laptop', '#A66BFF', 'sergio_f'),
('Cine_sergio_G', 'fernandez.fam@email.com', 0, 'fa-solid fa-clapperboard', '#FF6B6B', 'sergio_f'),
('Paga_sergio_I', 'fernandez.fam@email.com', 1, 'fa-solid fa-money-bill-wave', '#8CCB5E', 'sergio_f'),
('Regalo_cumple_sergio_I', 'fernandez.fam@email.com', 1, 'fa-solid fa-gift', '#A66BFF', 'sergio_f'),
('Ahorro_coche_I', 'fernandez.fam@email.com', 1, 'fa-solid fa-piggy-bank', '#FF6B6B', 'sergio_f'),
/* Usuario: cris_f */
('Libros_cris_G', 'fernandez.fam@email.com', 0, 'fa-solid fa-book-open', '#4DA3FF', 'cris_f'),
('Ropa_cris_G', 'fernandez.fam@email.com', 0, 'fa-solid fa-shirt', '#8CCB5E', 'cris_f'),
('Taxi_cris_G', 'fernandez.fam@email.com', 0, 'fa-solid fa-taxi', '#FF6B6B', 'cris_f'),
('Beca_cris_I', 'fernandez.fam@email.com', 1, 'fa-solid fa-money-bill-wave', '#FFD84D', 'cris_f'),
('Aportacion_cris_I', 'fernandez.fam@email.com', 1, 'fa-solid fa-piggy-bank', '#4DA3FF', 'cris_f'),
('Regalo_abuela_I', 'fernandez.fam@email.com', 1, 'fa-solid fa-gift', '#8CCB5E', 'cris_f'),

/* --- Familia: martinez.team@email.com --- */
/* Usuario: jefe_martinez (Admin) */
('Casa_jefe_G', 'martinez.team@email.com', 0, 'fa-solid fa-house', '#FF6B6B', 'jefe_martinez'),
('Comida_jefe_G', 'martinez.team@email.com', 0, 'fa-solid fa-utensils', '#FFD84D', 'jefe_martinez'),
('Gasolina_jefe_G', 'martinez.team@email.com', 0, 'fa-solid fa-car', '#4DA3FF', 'jefe_martinez'),
('Salario_jefe_I_M', 'martinez.team@email.com', 1, 'fa-solid fa-money-bill-wave', '#8CCB5E', 'jefe_martinez'),
('Alquiler_jefe_I', 'martinez.team@email.com', 1, 'fa-solid fa-building', '#4DA3FF', 'jefe_martinez'),
('Bonos_emp_I', 'martinez.team@email.com', 1, 'fa-solid fa-star', '#FFD84D', 'jefe_martinez'),
/* Usuario: isabel_m */
('Ropa_isabel_G', 'martinez.team@email.com', 0, 'fa-solid fa-shirt', '#A66BFF', 'isabel_m'),
('Viaje_isabel_G', 'martinez.team@email.com', 0, 'fa-solid fa-plane', '#FF6B6B', 'isabel_m'),
('Gimnasio_isabel_G', 'martinez.team@email.com', 0, 'fa-solid fa-weight-hanging', '#FFD84D', 'isabel_m'),
('Freelance_isabel_I', 'martinez.team@email.com', 1, 'fa-solid fa-laptop-code', '#A66BFF', 'isabel_m'),
('Regalo_jefe_I', 'martinez.team@email.com', 1, 'fa-solid fa-gift', '#8CCB5E', 'isabel_m'),
('Dividendos_isabel_I', 'martinez.team@email.com', 1, 'fa-solid fa-chart-bar', '#FF6B6B', 'isabel_m'),
/* Usuario: raul_m */
('Tecnologia_raul_G', 'martinez.team@email.com', 0, 'fa-solid fa-laptop', '#FFD84D', 'raul_m'),
('Deportes_raul_G', 'martinez.team@email.com', 0, 'fa-solid fa-dumbbell', '#A66BFF', 'raul_m'),
('Cursos_raul_G', 'martinez.team@email.com', 0, 'fa-solid fa-graduation-cap', '#8CCB5E', 'raul_m'),
('Paga_raul_I', 'martinez.team@email.com', 1, 'fa-solid fa-money-bill-wave', '#8CCB5E', 'raul_m'),
('Ahorro_playa_I', 'martinez.team@email.com', 1, 'fa-solid fa-piggy-bank', '#A66BFF', 'raul_m'),
('Regalo_amigos_I', 'martinez.team@email.com', 1, 'fa-solid fa-gift', '#4DA3FF', 'raul_m'),
/* Usuario: nuria_m */
('Ropa_nuria_G', 'martinez.team@email.com', 0, 'fa-solid fa-shirt', '#4DA3FF', 'nuria_m'),
('Cine_nuria_G', 'martinez.team@email.com', 0, 'fa-solid fa-clapperboard', '#8CCB5E', 'nuria_m'),
('Farmacia_nuria_G', 'martinez.team@email.com', 0, 'fa-solid fa-pills', '#FFD84D', 'nuria_m'),
('Beca_nuria_I', 'martinez.team@email.com', 1, 'fa-solid fa-money-bill-wave', '#FFD84D', 'nuria_m'),
('Aportacion_nuria_I', 'martinez.team@email.com', 1, 'fa-solid fa-piggy-bank', '#4DA3FF', 'nuria_m'),
('Regalo_papa_I', 'martinez.team@email.com', 1, 'fa-solid fa-gift', '#A66BFF', 'nuria_m'),

/* --- Familia: sanchez.clan@email.com --- */
/* Usuario: admin_sanchez (Admin) */
('Alquiler_admin_G', 'sanchez.clan@email.com', 0, 'fa-solid fa-building', '#FF6B6B', 'admin_sanchez'),
('Luz_Gas_admin_G', 'sanchez.clan@email.com', 0, 'fa-solid fa-plug', '#FFD84D', 'admin_sanchez'),
('Salud_Familiar_G', 'sanchez.clan@email.com', 0, 'fa-solid fa-heart-pulse', '#4DA3FF', 'admin_sanchez'),
('Salario_Admin_S_I', 'sanchez.clan@email.com', 1, 'fa-solid fa-money-bill-wave', '#8CCB5E', 'admin_sanchez'),
('Dividendos_admin_I', 'sanchez.clan@email.com', 1, 'fa-solid fa-chart-bar', '#4DA3FF', 'admin_sanchez'),
('Inversion_admin_I', 'sanchez.clan@email.com', 1, 'fa-solid fa-chart-line', '#A66BFF', 'admin_sanchez'),
/* Usuario: patricia_s */
('Comida_pat_G', 'sanchez.clan@email.com', 0, 'fa-solid fa-utensils', '#A66BFF', 'patricia_s'),
('Ropa_pat_G', 'sanchez.clan@email.com', 0, 'fa-solid fa-shirt', '#FF6B6B', 'patricia_s'),
('Viaje_pat_G', 'sanchez.clan@email.com', 0, 'fa-solid fa-plane', '#8CCB5E', 'patricia_s'),
('Freelance_pat_I', 'sanchez.clan@email.com', 1, 'fa-solid fa-laptop-code', '#A66BFF', 'patricia_s'),
('Bonos_pat_I', 'sanchez.clan@email.com', 1, 'fa-solid fa-star', '#8CCB5E', 'patricia_s'),
('Regalo_admin_I', 'sanchez.clan@email.com', 1, 'fa-solid fa-gift', '#FFD84D', 'admin_sanchez'),
/* Usuario: axel_s */
('Tecnologia_axel_G', 'sanchez.clan@email.com', 0, 'fa-solid fa-laptop', '#FFD84D', 'axel_s'),
('Deportes_axel_G', 'sanchez.clan@email.com', 0, 'fa-solid fa-dumbbell', '#A66BFF', 'axel_s'),
('Cine_axel_G', 'sanchez.clan@email.com', 0, 'fa-solid fa-clapperboard', '#4DA3FF', 'axel_s'),
('Paga_axel_I', 'sanchez.clan@email.com', 1, 'fa-solid fa-money-bill-wave', '#8CCB5E', 'axel_s'),
('Regalo_axel_I', 'sanchez.clan@email.com', 1, 'fa-solid fa-gift', '#FFD84D', 'axel_s'),
('Ahorro_estudio_I', 'sanchez.clan@email.com', 1, 'fa-solid fa-piggy-bank', '#A66BFF', 'axel_s'),
/* Usuario: lidia_s */
('Ropa_lidia_G', 'sanchez.clan@email.com', 0, 'fa-solid fa-shirt', '#4DA3FF', 'lidia_s'),
('Salud_lidia_G', 'sanchez.clan@email.com', 0, 'fa-solid fa-heart-pulse', '#8CCB5E', 'lidia_s'),
('Gimnasio_lidia_G', 'sanchez.clan@email.com', 0, 'fa-solid fa-weight-hanging', '#FF6B6B', 'lidia_s'),
('Beca_lidia_I', 'sanchez.clan@email.com', 1, 'fa-solid fa-money-bill-wave', '#FFD84D', 'lidia_s'),
('Aportacion_lidia_I', 'sanchez.clan@email.com', 1, 'fa-solid fa-piggy-bank', '#4DA3FF', 'lidia_s'),
('Bonos_lidia_I', 'sanchez.clan@email.com', 1, 'fa-solid fa-star', '#A66BFF', 'lidia_s'),

/* --- Familia: gomez.group@email.com --- */
/* Usuario: jefe_gomez (Admin) */
('Alquiler_jefe_G_G', 'gomez.group@email.com', 0, 'fa-solid fa-building', '#FF6B6B', 'jefe_gomez'),
('Agua_Internet_G', 'gomez.group@email.com', 0, 'fa-solid fa-droplet', '#FFD84D', 'jefe_gomez'),
('Mantenimiento_casa_G', 'gomez.group@email.com', 0, 'fa-solid fa-hammer', '#4DA3FF', 'jefe_gomez'),
('Salario_jefe_G_I', 'gomez.group@email.com', 1, 'fa-solid fa-money-bill-wave', '#8CCB5E', 'jefe_gomez'),
('Inversion_jefe_G_I', 'gomez.group@email.com', 1, 'fa-solid fa-chart-line', '#4DA3FF', 'jefe_gomez'),
('Dividendos_jefe_G_I', 'gomez.group@email.com', 1, 'fa-solid fa-chart-bar', '#FFD84D', 'jefe_gomez'),
/* Usuario: rocio_g */
('Ropa_rocio_G', 'gomez.group@email.com', 0, 'fa-solid fa-shirt', '#A66BFF', 'rocio_g'),
('Salud_rocio_G', 'gomez.group@email.com', 0, 'fa-solid fa-heart-pulse', '#FF6B6B', 'rocio_g'),
('Cursos_rocio_G', 'gomez.group@email.com', 0, 'fa-solid fa-graduation-cap', '#4DA3FF', 'rocio_g'),
('Freelance_rocio_I', 'gomez.group@email.com', 1, 'fa-solid fa-laptop-code', '#A66BFF', 'rocio_g'),
('Regalo_rocio_I', 'gomez.group@email.com', 1, 'fa-solid fa-gift', '#8CCB5E', 'rocio_g'),
('Comisiones_rocio_I', 'gomez.group@email.com', 1, 'fa-solid fa-percent', '#FF6B6B', 'rocio_g'),
/* Usuario: oscar_g */
('Tecnologia_oscar_G', 'gomez.group@email.com', 0, 'fa-solid fa-laptop', '#FFD84D', 'oscar_g'),
('Deportes_oscar_G', 'gomez.group@email.com', 0, 'fa-solid fa-dumbbell', '#A66BFF', 'oscar_g'),
('Cine_oscar_G', 'gomez.group@email.com', 0, 'fa-solid fa-clapperboard', '#FF6B6B', 'oscar_g'),
('Paga_oscar_I', 'gomez.group@email.com', 1, 'fa-solid fa-money-bill-wave', '#8CCB5E', 'oscar_g'),
('Regalo_oscar_I', 'gomez.group@email.com', 1, 'fa-solid fa-gift', '#FFD84D', 'oscar_g'),
('Ahorro_oscar_I', 'gomez.group@email.com', 1, 'fa-solid fa-piggy-bank', '#A66BFF', 'oscar_g'),
/* Usuario: alba_g */
('Ropa_alba_G', 'gomez.group@email.com', 0, 'fa-solid fa-shirt', '#4DA3FF', 'alba_g'),
('Libros_alba_G', 'gomez.group@email.com', 0, 'fa-solid fa-book-open', '#8CCB5E', 'alba_g'),
('Farmacia_alba_G', 'gomez.group@email.com', 0, 'fa-solid fa-pills', '#FFD84D', 'alba_g'),
('Beca_alba_I', 'gomez.group@email.com', 1, 'fa-solid fa-money-bill-wave', '#FFD84D', 'alba_g'),
('Aportacion_alba_I', 'gomez.group@email.com', 1, 'fa-solid fa-piggy-bank', '#4DA3FF', 'alba_g'),
('Regalo_mama_alba_I', 'gomez.group@email.com', 1, 'fa-solid fa-gift', '#A66BFF', 'alba_g'),

/* --- Familia: diaz.family@email.com --- */
/* Usuario: admin_diaz (Admin) */
('Casa_admin_D_G', 'diaz.family@email.com', 0, 'fa-solid fa-house', '#FF6B6B', 'admin_diaz'),
('Electricidad_admin_D_G', 'diaz.family@email.com', 0, 'fa-solid fa-plug', '#FFD84D', 'admin_diaz'),
('Gas_admin_D_G', 'diaz.family@email.com', 0, 'fa-solid fa-fire-flame-simple', '#4DA3FF', 'admin_diaz'),
('Salario_admin_D_I', 'diaz.family@email.com', 1, 'fa-solid fa-money-bill-wave', '#8CCB5E', 'admin_diaz'),
('Inversion_admin_D_I', 'diaz.family@email.com', 1, 'fa-solid fa-chart-line', '#4DA3FF', 'admin_diaz'),
('Alquiler_admin_D_I', 'diaz.family@email.com', 1, 'fa-solid fa-building', '#FFD84D', 'admin_diaz'),
/* Usuario: mercedes_d */
('Ropa_mercedes_G', 'diaz.family@email.com', 0, 'fa-solid fa-shirt', '#A66BFF', 'mercedes_d'),
('Comida_mercedes_G', 'diaz.family@email.com', 0, 'fa-solid fa-utensils', '#FF6B6B', 'mercedes_d'),
('Ocio_mercedes_G', 'diaz.family@email.com', 0, 'fa-solid fa-film', '#4DA3FF', 'mercedes_d'),
('Freelance_mercedes_I', 'diaz.family@email.com', 1, 'fa-solid fa-laptop-code', '#A66BFF', 'mercedes_d'),
('Regalo_mercedes_I', 'diaz.family@email.com', 1, 'fa-solid fa-gift', '#8CCB5E', 'mercedes_d'),
('Dividendos_mercedes_I', 'diaz.family@email.com', 1, 'fa-solid fa-chart-bar', '#FF6B6B', 'mercedes_d'),
/* Usuario: victor_d */
('Tecnologia_victor_G', 'diaz.family@email.com', 0, 'fa-solid fa-laptop', '#FFD84D', 'victor_d'),
('Deportes_victor_G', 'diaz.family@email.com', 0, 'fa-solid fa-dumbbell', '#A66BFF', 'victor_d'),
('Transporte_victor_G', 'diaz.family@email.com', 0, 'fa-solid fa-car', '#4DA3FF', 'victor_d'),
('Paga_victor_I', 'diaz.family@email.com', 1, 'fa-solid fa-money-bill-wave', '#8CCB5E', 'victor_d'),
('Regalo_victor_I', 'diaz.family@email.com', 1, 'fa-solid fa-gift', '#FFD84D', 'victor_d'),
('Ahorro_viaje_I', 'diaz.family@email.com', 1, 'fa-solid fa-piggy-bank', '#A66BFF', 'victor_d'),
/* Usuario: andrea_d */
('Viajes_andrea_G', 'diaz.family@email.com', 0, 'fa-solid fa-plane', '#4DA3FF', 'andrea_d'),
('Ropa_andrea_G', 'diaz.family@email.com', 0, 'fa-solid fa-shirt', '#8CCB5E', 'andrea_d'),
('Farmacia_andrea_G', 'diaz.family@email.com', 0, 'fa-solid fa-pills', '#FF6B6B', 'andrea_d'),
('Beca_andrea_I', 'diaz.family@email.com', 1, 'fa-solid fa-money-bill-wave', '#FFD84D', 'andrea_d'),
('Aportacion_andrea_I', 'diaz.family@email.com', 1, 'fa-solid fa-piggy-bank', '#4DA3FF', 'andrea_d'),
('Regalo_navidad_andrea_I', 'diaz.family@email.com', 1, 'fa-solid fa-gift', '#A66BFF', 'andrea_d'),

/* --- Familia: alvarez.house@email.com --- */
/* Usuario: jefe_alvarez (Admin) */
('Alquiler_jefe_A_G', 'alvarez.house@email.com', 0, 'fa-solid fa-building', '#FF6B6B', 'jefe_alvarez'),
('Internet_jefe_A_G', 'alvarez.house@email.com', 0, 'fa-solid fa-wifi', '#FFD84D', 'jefe_alvarez'),
('Comida_jefe_A_G', 'alvarez.house@email.com', 0, 'fa-solid fa-utensils', '#FFD84D', 'jefe_alvarez'),
('Salario_jefe_A_I', 'alvarez.house@email.com', 1, 'fa-solid fa-money-bill-wave', '#8CCB5E', 'jefe_alvarez'),
('Inversion_jefe_A_I', 'alvarez.house@email.com', 1, 'fa-solid fa-chart-line', '#4DA3FF', 'jefe_alvarez'),
('Dividendos_jefe_A_I', 'alvarez.house@email.com', 1, 'fa-solid fa-chart-bar', '#FFD84D', 'jefe_alvarez'),
/* Usuario: rosa_a */
('Ropa_rosa_G', 'alvarez.house@email.com', 0, 'fa-solid fa-shirt', '#A66BFF', 'rosa_a'),
('Salud_rosa_G', 'alvarez.house@email.com', 0, 'fa-solid fa-heart-pulse', '#FF6B6B', 'rosa_a'),
('Cine_rosa_G', 'alvarez.house@email.com', 0, 'fa-solid fa-clapperboard', '#4DA3FF', 'rosa_a'),
('Freelance_rosa_I', 'alvarez.house@email.com', 1, 'fa-solid fa-laptop-code', '#A66BFF', 'rosa_a'),
('Bonos_rosa_I', 'alvarez.house@email.com', 1, 'fa-solid fa-star', '#8CCB5E', 'rosa_a'),
('Aportacion_rosa_I', 'alvarez.house@email.com', 1, 'fa-solid fa-piggy-bank', '#FFD84D', 'rosa_a'),
/* Usuario: dani_a */
('Tecnologia_dani_G', 'alvarez.house@email.com', 0, 'fa-solid fa-laptop', '#FFD84D', 'dani_a'),
('Deportes_dani_G', 'alvarez.house@email.com', 0, 'fa-solid fa-dumbbell', '#A66BFF', 'dani_a'),
('Gimnasio_dani_G', 'alvarez.house@email.com', 0, 'fa-solid fa-weight-hanging', '#4DA3FF', 'dani_a'),
('Paga_dani_I', 'alvarez.house@email.com', 1, 'fa-solid fa-money-bill-wave', '#8CCB5E', 'dani_a'),
('Regalo_dani_I', 'alvarez.house@email.com', 1, 'fa-solid fa-gift', '#FFD84D', 'dani_a'),
('Ahorro_dani_I', 'alvarez.house@email.com', 1, 'fa-solid fa-piggy-bank', '#A66BFF', 'dani_a'),
/* Usuario: eva_a */
('Viajes_eva_A_G', 'alvarez.house@email.com', 0, 'fa-solid fa-plane', '#4DA3FF', 'eva_a'),
('Ropa_eva_A_G', 'alvarez.house@email.com', 0, 'fa-solid fa-shirt', '#8CCB5E', 'eva_a'),
('Cursos_eva_A_G', 'alvarez.house@email.com', 0, 'fa-solid fa-graduation-cap', '#A66BFF', 'eva_a'),
('Beca_eva_A_I', 'alvarez.house@email.com', 1, 'fa-solid fa-money-bill-wave', '#FFD84D', 'eva_a'),
('Aportacion_eva_A_I', 'alvarez.house@email.com', 1, 'fa-solid fa-piggy-bank', '#4DA3FF', 'eva_a'),
('Regalo_papa_eva_I', 'alvarez.house@email.com', 1, 'fa-solid fa-gift', '#A66BFF', 'eva_a'),

/* --- Familia: moreno.co@email.com --- */
/* Usuario: admin_moreno (Admin) */
('Casa_admin_M_G', 'moreno.co@email.com', 0, 'fa-solid fa-house', '#FF6B6B', 'admin_moreno'),
('Mantenimiento_admin_M_G', 'moreno.co@email.com', 0, 'fa-solid fa-hammer', '#FFD84D', 'admin_moreno'),
('Seguro_admin_M_G', 'moreno.co@email.com', 0, 'fa-solid fa-heart-pulse', '#4DA3FF', 'admin_moreno'),
('Salario_admin_M_I', 'moreno.co@email.com', 1, 'fa-solid fa-money-bill-wave', '#8CCB5E', 'admin_moreno'),
('Inversion_admin_M_I', 'moreno.co@email.com', 1, 'fa-solid fa-chart-line', '#4DA3FF', 'admin_moreno'),
('Alquiler_admin_M_I', 'moreno.co@email.com', 1, 'fa-solid fa-building', '#FFD84D', 'admin_moreno'),
/* Usuario: marta_m */
('Ropa_marta_G', 'moreno.co@email.com', 0, 'fa-solid fa-shirt', '#A66BFF', 'marta_m'),
('Viajes_marta_G', 'moreno.co@email.com', 0, 'fa-solid fa-plane', '#FF6B6B', 'marta_m'),
('Ocio_marta_G', 'moreno.co@email.com', 0, 'fa-solid fa-film', '#FFD84D', 'marta_m'),
('Freelance_marta_I', 'moreno.co@email.com', 1, 'fa-solid fa-laptop-code', '#A66BFF', 'marta_m'),
('Regalo_marta_I', 'moreno.co@email.com', 1, 'fa-solid fa-gift', '#8CCB5E', 'marta_m'),
('Bonos_marta_I', 'moreno.co@email.com', 1, 'fa-solid fa-star', '#FF6B6B', 'marta_m'),
/* Usuario: felix_m */
('Tecnologia_felix_G', 'moreno.co@email.com', 0, 'fa-solid fa-laptop', '#FFD84D', 'felix_m'),
('Deportes_felix_G', 'moreno.co@email.com', 0, 'fa-solid fa-dumbbell', '#A66BFF', 'felix_m'),
('Cine_felix_G', 'moreno.co@email.com', 0, 'fa-solid fa-clapperboard', '#FF6B6B', 'felix_m'),
('Paga_felix_I', 'moreno.co@email.com', 1, 'fa-solid fa-money-bill-wave', '#8CCB5E', 'felix_m'),
('Regalo_felix_I', 'moreno.co@email.com', 1, 'fa-solid fa-gift', '#FFD84D', 'felix_m'),
('Ahorro_playa_felix_I', 'moreno.co@email.com', 1, 'fa-solid fa-piggy-bank', '#A66BFF', 'felix_m'),
/* Usuario: nerea_m */
('Ropa_nerea_G', 'moreno.co@email.com', 0, 'fa-solid fa-shirt', '#4DA3FF', 'nerea_m'),
('Farmacia_nerea_G', 'moreno.co@email.com', 0, 'fa-solid fa-pills', '#8CCB5E', 'nerea_m'),
('Cursos_nerea_G', 'moreno.co@email.com', 0, 'fa-solid fa-graduation-cap', '#FFD84D', 'nerea_m'),
('Beca_nerea_I', 'moreno.co@email.com', 1, 'fa-solid fa-money-bill-wave', '#FFD84D', 'nerea_m'),
('Aportacion_nerea_I', 'moreno.co@email.com', 1, 'fa-solid fa-piggy-bank', '#4DA3FF', 'nerea_m'),
('Regalo_abuela_nerea_I', 'moreno.co@email.com', 1, 'fa-solid fa-gift', '#A66BFF', 'nerea_m');


INSERT INTO `movimiento` (`fecha`, `monto`, `descripcion`, `nombreUsuario`, `nombreConcepto`, `correoFamilia`) VALUES
/* --- Familia: familia.garcia@email.com --- */
/* jefe_garcia */
('2025-10-01', 1200.00, 'Pago mensual de la hipoteca', 'jefe_garcia', 'Alquiler/Hipoteca_G', 'familia.garcia@email.com'),
('2025-10-05', 3500.50, 'Nómina principal del mes de Octubre', 'jefe_garcia', 'Salario_Ppal_I', 'familia.garcia@email.com'),
/* eva_g */
('2025-10-03', 150.00, 'Consulta médica de rutina', 'eva_g', 'Salud_eva_G', 'familia.garcia@email.com'),
('2025-10-10', 450.00, 'Ingreso extra por proyecto de diseño freelance', 'eva_g', 'Freelance_eva_I', 'familia.garcia@email.com'),
/* hugo_g */
('2025-10-12', 39.99, 'Suscripción a servicio de juegos online', 'hugo_g', 'Videojuegos_G', 'familia.garcia@email.com'),
('2025-10-07', 50.00, 'Paga semanal recibida', 'hugo_g', 'Paga_hugo_I', 'familia.garcia@email.com'),
/* laura_g */
('2025-10-09', 120.00, 'Matrícula de curso de idiomas', 'laura_g', 'Cursos_laura_G', 'familia.garcia@email.com'),
('2025-10-02', 200.00, 'Depósito de beca estudiantil', 'laura_g', 'Beca_laura_I', 'familia.garcia@email.com'),

/* --- Familia: los.rodriguez@email.com --- */
/* admin_rodri */
('2025-10-01', 550.00, 'Pago del seguro del hogar y vehículo', 'admin_rodri', 'Seguro_admin_G', 'los.rodriguez@email.com'),
('2025-10-04', 4100.00, 'Nómina secundaria', 'admin_rodri', 'Salario_Sec_I', 'los.rodriguez@email.com'),
/* ana_r */
('2025-10-06', 40.00, 'Cuota mensual del club deportivo', 'ana_r', 'Gimnasio_ana_G', 'los.rodriguez@email.com'),
('2025-10-13', 280.00, 'Ingreso por diseño de logotipo freelance', 'ana_r', 'Freelance_ana_I', 'los.rodriguez@email.com'),
/* pablo_r */
('2025-10-08', 30.00, 'Entradas de cine y palomitas', 'pablo_r', 'Cine_pablo_G', 'los.rodriguez@email.com'),
('2025-10-16', 80.00, 'Paga por ayudar en la tienda familiar', 'pablo_r', 'Paga_pablo_I', 'los.rodriguez@email.com'),
/* sofia_r */
('2025-10-24', 88.00, 'Compra de vestido y accesorios', 'sofia_r', 'Ropa_sofia_G', 'los.rodriguez@email.com'),
('2025-10-17', 25.00, 'Regalo de cumpleaños de una tía', 'sofia_r', 'Regalo_Tia_I', 'los.rodriguez@email.com'),

/* --- Familia: perez.home@email.com --- */
/* jefe_perez */
('2025-10-01', 3800.00, 'Ingreso de salario mensual', 'jefe_perez', 'Salario_jefe_I', 'perez.home@email.com'),
('2025-10-05', 95.70, 'Factura de gas natural', 'jefe_perez', 'Gas_jefe_G', 'perez.home@email.com'),
/* clara_p */
('2025-10-11', 200.00, 'Gasto en alimentación especializada', 'clara_p', 'Comida_clara_G', 'perez.home@email.com'),
('2025-10-03', 150.00, 'Venta de artículo de segunda mano', 'clara_p', 'Bonos_clara_I', 'perez.home@email.com'),
/* diego_p */
('2025-10-20', 15.00, 'Compra de libros de texto', 'diego_p', 'Libros_diego_G', 'perez.home@email.com'),
('2025-10-28', 40.00, 'Paga de la semana', 'diego_p', 'Paga_diego_I', 'perez.home@email.com'),
/* irene_p */
('2025-10-16', 78.00, 'Gastos en ropa para evento', 'irene_p', 'Ropa_irene_G', 'perez.home@email.com'),
('2025-10-09', 50.00, 'Ingreso por comisiones', 'irene_p', 'Bonos_irene_I', 'perez.home@email.com'),

/* --- Familia: fernandez.fam@email.com --- */
/* admin_fer */
('2025-10-02', 3600.00, 'Salario mensual del administrador', 'admin_fer', 'Salario_Admin_I', 'fernandez.fam@email.com'),
('2025-10-14', 45.00, 'Pago del servicio de agua', 'admin_fer', 'Agua_admin_G', 'fernandez.fam@email.com'),
/* elena_f */
('2025-10-04', 180.00, 'Pago por proyecto web freelance', 'elena_f', 'Freelance_elena_I', 'fernandez.fam@email.com'),
('2025-10-12', 60.00, 'Clase de cocina avanzada', 'elena_f', 'Cursos_elena_G', 'fernandez.fam@email.com'),
/* sergio_f */
('2025-10-15', 55.00, 'Compra de un videojuego nuevo', 'sergio_f', 'Videojuegos_sergio_G', 'fernandez.fam@email.com'),
('2025-10-07', 30.00, 'Paga semanal', 'sergio_f', 'Paga_sergio_I', 'fernandez.fam@email.com'),
/* cris_f */
('2025-10-18', 40.00, 'Compra de ropa de invierno', 'cris_f', 'Ropa_cris_G', 'fernandez.fam@email.com'),
('2025-10-11', 100.00, 'Aportación a cuenta de ahorros', 'cris_f', 'Aportacion_cris_I', 'fernandez.fam@email.com'),

/* --- Familia: martinez.team@email.com --- */
/* jefe_martinez */
('2025-10-01', 3200.00, 'Salario del mes', 'jefe_martinez', 'Salario_jefe_I_M', 'martinez.team@email.com'),
('2025-10-06', 800.00, 'Pago del alquiler de la vivienda', 'jefe_martinez', 'Casa_jefe_G', 'martinez.team@email.com'),
/* isabel_m */
('2025-10-12', 45.00, 'Cuota de gimnasio', 'isabel_m', 'Gimnasio_isabel_G', 'martinez.team@email.com'),
('2025-10-05', 90.00, 'Ingreso por trabajo temporal freelance', 'isabel_m', 'Freelance_isabel_I', 'martinez.team@email.com'),
/* raul_m */
('2025-10-16', 75.00, 'Mantenimiento y reparación de tablet', 'raul_m', 'Tecnologia_raul_G', 'martinez.team@email.com'),
('2025-10-08', 50.00, 'Regalo de cumpleaños recibido', 'raul_m', 'Regalo_amigos_I', 'martinez.team@email.com'),
/* nuria_m */
('2025-10-04', 15.00, 'Entrada de cine', 'nuria_m', 'Cine_nuria_G', 'martinez.team@email.com'),
('2025-10-26', 150.00, 'Depósito de beca', 'nuria_m', 'Beca_nuria_I', 'martinez.team@email.com'),

/* --- Familia: sanchez.clan@email.com --- */
/* admin_sanchez */
('2025-10-01', 950.00, 'Pago de alquiler', 'admin_sanchez', 'Alquiler_admin_G', 'sanchez.clan@email.com'),
('2025-10-05', 3100.00, 'Ingreso de salario principal', 'admin_sanchez', 'Salario_Admin_S_I', 'sanchez.clan@email.com'),
/* patricia_s */
('2025-10-17', 220.00, 'Compra mensual de supermercado', 'patricia_s', 'Comida_pat_G', 'sanchez.clan@email.com'),
('2025-10-09', 150.00, 'Ingreso por venta de curso online', 'patricia_s', 'Freelance_pat_I', 'sanchez.clan@email.com'),
/* axel_s */
('2025-10-22', 12.00, 'Entrada de cine', 'axel_s', 'Cine_axel_G', 'sanchez.clan@email.com'),
('2025-10-14', 100.00, 'Paga mensual', 'axel_s', 'Paga_axel_I', 'sanchez.clan@email.com'),
/* lidia_s */
('2025-10-04', 35.00, 'Gasto en vitaminas', 'lidia_s', 'Salud_lidia_G', 'sanchez.clan@email.com'),
('2025-10-15', 50.00, 'Aportación a su cuenta de ahorros', 'lidia_s', 'Aportacion_lidia_I', 'sanchez.clan@email.com'),

/* --- Familia: gomez.group@email.com --- */
/* jefe_gomez */
('2025-10-01', 3400.00, 'Salario del mes', 'jefe_gomez', 'Salario_jefe_G_I', 'gomez.group@email.com'),
('2025-10-06', 750.00, 'Pago de alquiler', 'jefe_gomez', 'Alquiler_jefe_G_G', 'gomez.group@email.com'),
/* rocio_g */
('2025-10-03', 120.00, 'Gastos en vestimenta', 'rocio_g', 'Ropa_rocio_G', 'gomez.group@email.com'),
('2025-10-11', 200.00, 'Ingreso por trabajo de fin de semana', 'rocio_g', 'Freelance_rocio_I', 'gomez.group@email.com'),
/* oscar_g */
('2025-10-15', 120.00, 'Compra de teclado y ratón gaming', 'oscar_g', 'Tecnologia_oscar_G', 'gomez.group@email.com'),
('2025-10-08', 40.00, 'Paga semanal', 'oscar_g', 'Paga_oscar_I', 'gomez.group@email.com'),
/* alba_g */
('2025-10-20', 45.00, 'Gastos en farmacia', 'alba_g', 'Farmacia_alba_G', 'gomez.group@email.com'),
('2025-10-05', 150.00, 'Ingreso de beca deportiva', 'alba_g', 'Beca_alba_I', 'gomez.group@email.com'),

/* --- Familia: diaz.family@email.com --- */
/* admin_diaz */
('2025-10-01', 3700.00, 'Salario mensual', 'admin_diaz', 'Salario_admin_D_I', 'diaz.family@email.com'),
('2025-10-07', 850.00, 'Pago de alquiler de la casa', 'admin_diaz', 'Casa_admin_D_G', 'diaz.family@email.com'),
/* mercedes_d */
('2025-10-11', 250.00, 'Ingreso por proyecto de consultoría', 'mercedes_d', 'Freelance_mercedes_I', 'diaz.family@email.com'),
('2025-10-04', 130.00, 'Compra de ropa', 'mercedes_d', 'Ropa_mercedes_G', 'diaz.family@email.com'),
/* victor_d */
('2025-10-23', 85.00, 'Compra de zapatillas deportivas', 'victor_d', 'Deportes_victor_G', 'diaz.family@email.com'),
('2025-10-08', 50.00, 'Paga recibida', 'victor_d', 'Paga_victor_I', 'diaz.family@email.com'),
/* andrea_d */
('2025-10-17', 200.00, 'Gasto en viaje de fin de semana', 'andrea_d', 'Viajes_andrea_G', 'diaz.family@email.com'),
('2025-10-10', 100.00, 'Regalo de cumpleaños', 'andrea_d', 'Regalo_navidad_andrea_I', 'diaz.family@email.com'),

/* --- Familia: alvarez.house@email.com --- */
/* jefe_alvarez */
('2025-10-01', 4000.00, 'Salario del mes', 'jefe_alvarez', 'Salario_jefe_A_I', 'alvarez.house@email.com'),
('2025-10-06', 700.00, 'Pago mensual de alquiler', 'jefe_alvarez', 'Alquiler_jefe_A_G', 'alvarez.house@email.com'),
/* rosa_a */
('2025-10-18', 88.00, 'Compra de ropa para niños', 'rosa_a', 'Ropa_rosa_G', 'alvarez.house@email.com'),
('2025-10-04', 150.00, 'Aportación a su cuenta de ahorro', 'rosa_a', 'Aportacion_rosa_I', 'alvarez.house@email.com'),
/* dani_a */
('2025-10-15', 90.00, 'Compra de videojuegos y accesorios', 'dani_a', 'Tecnologia_dani_G', 'alvarez.house@email.com'),
('2025-10-08', 35.00, 'Ingreso por pequeños trabajos', 'dani_a', 'Paga_dani_I', 'alvarez.house@email.com'),
/* eva_a */
('2025-10-17', 150.00, 'Gasto en clases extraescolares', 'eva_a', 'Cursos_eva_A_G', 'alvarez.house@email.com'),
('2025-10-10', 120.00, 'Beca escolar recibida', 'eva_a', 'Beca_eva_A_I', 'alvarez.house@email.com'),

/* --- Familia: moreno.co@email.com --- */
/* admin_moreno */
('2025-10-01', 3300.00, 'Salario del mes', 'admin_moreno', 'Salario_admin_M_I', 'moreno.co@email.com'),
('2025-10-06', 780.00, 'Pago del alquiler', 'admin_moreno', 'Casa_admin_M_G', 'moreno.co@email.com'),
/* marta_m */
('2025-10-04', 120.00, 'Compra de vuelo para vacaciones', 'marta_m', 'Viajes_marta_G', 'moreno.co@email.com'),
('2025-10-11', 280.00, 'Ingreso por trabajo de ilustración freelance', 'marta_m', 'Freelance_marta_I', 'moreno.co@email.com'),
/* felix_m */
('2025-10-15', 75.00, 'Compra de nuevos auriculares', 'felix_m', 'Tecnologia_felix_G', 'moreno.co@email.com'),
('2025-10-08', 60.00, 'Regalo de sus padres', 'felix_m', 'Regalo_felix_I', 'moreno.co@email.com'),
/* nerea_m */
('2025-10-17', 45.00, 'Gastos en farmacia y maquillaje', 'nerea_m', 'Farmacia_nerea_G', 'moreno.co@email.com'),
('2025-10-10', 80.00, 'Aportación a su ahorro personal', 'nerea_m', 'Aportacion_nerea_I', 'moreno.co@email.com');


INSERT INTO `movimiento` (`fecha`, `monto`, `descripcion`, `nombreUsuario`, `nombreConcepto`, `correoFamilia`) VALUES
/* --- Movimientos del 18-11-2025 --- */
('2025-11-18', 85.50, 'Pago pendiente de factura de luz (Gasto)', 'jefe_garcia', 'Electricidad_G', 'familia.garcia@email.com'),
('2025-11-18', 49.99, 'Compra de juego nuevo (Gasto)', 'hugo_g', 'Videojuegos_G', 'familia.garcia@email.com'),

/* --- Movimientos del 19-11-2025 --- */
('2025-11-19', 150.00, 'Ingreso por comisión de venta (Ingreso)', 'eva_g', 'Comisiones_eva_I', 'familia.garcia@email.com'),
('2025-11-19', 120.00, 'Taller intensivo de idiomas (Gasto)', 'laura_g', 'Cursos_laura_G', 'familia.garcia@email.com'),

/* --- Movimientos del 20-11-2025 --- */
('2025-11-20', 1500.00, 'Anticipo de nómina (Ingreso)', 'jefe_garcia', 'Salario_Ppal_I', 'familia.garcia@email.com'),
('2025-11-20', 65.90, 'Compra de ropa de invierno (Gasto)', 'eva_g', 'Ropa_eva_G', 'familia.garcia@email.com'),
('2025-11-20', 20.00, 'Paga semanal recibida (Ingreso)', 'hugo_g', 'Paga_hugo_I', 'familia.garcia@email.com'),

/* --- Movimientos del 21-11-2025 --- */
('2025-11-21', 80.00, 'Ingreso de beca parcial (Ingreso)', 'laura_g', 'Beca_laura_I', 'familia.garcia@email.com'),
('2025-11-21', 12.30, 'Almuerzo rápido (Gasto)', 'hugo_g', 'Comida_rapida_G', 'familia.garcia@email.com'),
('2025-11-21', 50.00, 'Regalo por proyecto completado (Ingreso)', 'jefe_garcia', 'Regalos_Rec_I', 'familia.garcia@email.com');


INSERT INTO `personalizacionconcepto` (
  `montoPlanificado`,
  `tipoPeriodoPlanificado`,
  `diaPeriodoPlanificado`,
  `limiteGasto`,
  `tipoPeriodoLimite`,
  `diaPeriodoLimite`,
  `notificacion`,
  `activo`,
  `nombreUsuario`,
  `nombreConcepto`,
  `correoFamilia`
) VALUES
/* --- Familia: familia.garcia@email.com (jefe_garcia - Admin) --- */
('1200.00', 'mensual', 1, '1250.00', 'mensual', 5, 0, 1, 'jefe_garcia', 'Alquiler/Hipoteca_G', 'familia.garcia@email.com'),
('3500.00', 'mensual', 5, '4000.00', 'mensual', 10, 0, 1, 'jefe_garcia', 'Salario_Ppal_I', 'familia.garcia@email.com'),
('150.00', 'mensual', 15, '100.00', 'mensual', 20, 0, 1, 'jefe_garcia', 'Electricidad_G', 'familia.garcia@email.com'),
/* --- familia.garcia@email.com (eva_g) --- */
(NULL, NULL, NULL, '200.00', 'mensual', 10, 0, 1, 'eva_g', 'Salud_eva_G', 'familia.garcia@email.com'),
('500.00', 'mensual', 10, NULL, NULL, NULL, 0, 1, 'eva_g', 'Freelance_eva_I', 'familia.garcia@email.com'),
(NULL, NULL, NULL, '100.00', 'mensual', 15, 0, 1, 'eva_g', 'Ropa_eva_G', 'familia.garcia@email.com'),
/* --- familia.garcia@email.com (hugo_g) --- */
('40.00', 'diario', NULL, '50.00', 'diario', NULL, 0, 1, 'hugo_g', 'Videojuegos_G', 'familia.garcia@email.com'),
('200.00', 'mensual', 7, NULL, NULL, NULL, 0, 1, 'hugo_g', 'Paga_hugo_I', 'familia.garcia@email.com'),
/* --- familia.garcia@email.com (laura_g) --- */
('100.00', 'mensual', 1, '150.00', 'mensual', 7, 0, 1, 'laura_g', 'Cursos_laura_G', 'familia.garcia@email.com'),
(NULL, NULL, NULL, '50.00', 'quincenal', 15, 0, 1, 'laura_g', 'Deportes_laura_G', 'familia.garcia@email.com'),

/* --- Familia: los.rodriguez@email.com (admin_rodri) --- */
('500.00', 'mensual', 1, '600.00', 'mensual', 10, 0, 1, 'admin_rodri', 'Seguro_admin_G', 'los.rodriguez@email.com'),
('4000.00', 'mensual', 4, NULL, NULL, NULL, 0, 1, 'admin_rodri', 'Salario_Sec_I', 'los.rodriguez@email.com'),
('280.00', 'semanal', NULL, '300.00', 'semanal', NULL, 0, 1, 'admin_rodri', 'Comida_admin_G', 'los.rodriguez@email.com'),
/* --- los.rodriguez@email.com (ana_r) --- */
('45.00', 'mensual', 5, '45.00', 'mensual', 5, 0, 1, 'ana_r', 'Gimnasio_ana_G', 'los.rodriguez@email.com'),
('300.00', 'mensual', 13, NULL, NULL, NULL, 0, 1, 'ana_r', 'Freelance_ana_I', 'los.rodriguez@email.com'),
/* --- los.rodriguez@email.com (pablo_r) --- */
(NULL, NULL, NULL, '50.00', 'quincenal', 15, 0, 1, 'pablo_r', 'Cine_pablo_G', 'los.rodriguez@email.com'),
('100.00', 'mensual', 16, NULL, NULL, NULL, 0, 1, 'pablo_r', 'Paga_pablo_I', 'los.rodriguez@email.com'),
(NULL, NULL, NULL, '150.00', 'mensual', 25, 0, 1, 'pablo_r', 'Tecnologia_pablo_G', 'los.rodriguez@email.com'),
/* --- los.rodriguez@email.com (sofia_r) --- */
('10.00', 'diario', NULL, '12.00', 'diario', NULL, 0, 1, 'sofia_r', 'Cafe_G', 'los.rodriguez@email.com'),
('200.00', 'mensual', 1, NULL, NULL, NULL, 0, 1, 'sofia_r', 'Beca_sofia_I', 'los.rodriguez@email.com'),

/* --- Familia: perez.home@email.com (jefe_perez - Admin) --- */
('4000.00', 'mensual', 1, '4500.00', 'mensual', 1, 0, 1, 'jefe_perez', 'Salario_jefe_I', 'perez.home@email.com'),
(NULL, NULL, NULL, '100.00', 'mensual', 10, 0, 1, 'jefe_perez', 'Gas_jefe_G', 'perez.home@email.com'),
/* --- perez.home@email.com (clara_p) --- */
('250.00', 'quincenal', 15, '300.00', 'quincenal', 1, 0, 1, 'clara_p', 'Comida_clara_G', 'perez.home@email.com'),
('300.00', 'mensual', 15, NULL, NULL, NULL, 0, 1, 'clara_p', 'Freelance_clara_I', 'perez.home@email.com'),
/* --- perez.home@email.com (diego_p) --- */
(NULL, NULL, NULL, '100.00', 'mensual', 20, 0, 1, 'diego_p', 'Tecnologia_diego_G', 'perez.home@email.com'),
('50.00', 'semanal', NULL, NULL, NULL, NULL, 0, 1, 'diego_p', 'Paga_diego_I', 'perez.home@email.com'),
/* --- perez.home@email.com (irene_p) --- */
('100.00', 'mensual', 1, '150.00', 'mensual', 15, 0, 1, 'irene_p', 'Ropa_irene_G', 'perez.home@email.com'),
(NULL, NULL, NULL, '50.00', 'mensual', 25, 0, 1, 'irene_p', 'Cine_irene_G', 'perez.home@email.com'),

/* --- Familia: fernandez.fam@email.com (admin_fer - Admin) --- */
('3500.00', 'mensual', 2, NULL, NULL, NULL, 0, 1, 'admin_fer', 'Salario_Admin_I', 'fernandez.fam@email.com'),
(NULL, NULL, NULL, '50.00', 'mensual', 15, 0, 1, 'admin_fer', 'Agua_admin_G', 'fernandez.fam@email.com'),
('60.00', 'mensual', 21, '65.00', 'mensual', 25, 0, 1, 'admin_fer', 'Internet_admin_G', 'fernandez.fam@email.com'),
/* --- fernandez.fam@email.com (elena_f) --- */
('200.00', 'mensual', 4, NULL, NULL, NULL, 0, 1, 'elena_f', 'Freelance_elena_I', 'fernandez.fam@email.com'),
('150.00', 'quincenal', 15, '150.00', 'quincenal', 15, 0, 1, 'elena_f', 'Viaje_finde_G', 'fernandez.fam@email.com'),
/* --- fernandez.fam@email.com (sergio_f) --- */
('50.00', 'semanal', NULL, '75.00', 'semanal', NULL, 0, 1, 'sergio_f', 'Videojuegos_sergio_G', 'fernandez.fam@email.com'),
('100.00', 'mensual', 7, NULL, NULL, NULL, 0, 1, 'sergio_f', 'Paga_sergio_I', 'fernandez.fam@email.com'),
('100.00', 'mensual', 29, NULL, NULL, NULL, 0, 1, 'sergio_f', 'Ahorro_coche_I', 'fernandez.fam@email.com'),
/* --- fernandez.fam@email.com (cris_f) --- */
(NULL, NULL, NULL, '40.00', 'mensual', 20, 0, 1, 'cris_f', 'Ropa_cris_G', 'fernandez.fam@email.com'),
('50.00', 'mensual', 11, NULL, NULL, NULL, 0, 1, 'cris_f', 'Aportacion_cris_I', 'fernandez.fam@email.com'),

/* --- Familia: martinez.team@email.com (jefe_martinez - Admin) --- */
('3000.00', 'mensual', 1, NULL, NULL, NULL, 0, 1, 'jefe_martinez', 'Salario_jefe_I_M', 'martinez.team@email.com'),
('800.00', 'mensual', 5, '850.00', 'mensual', 10, 0, 1, 'jefe_martinez', 'Casa_jefe_G', 'martinez.team@email.com'),
('200.00', 'mensual', 13, NULL, NULL, NULL, 0, 1, 'jefe_martinez', 'Bonos_emp_I', 'martinez.team@email.com'),
('200.00', 'mensual', 25, '250.00', 'mensual', 28, 0, 1, 'jefe_martinez', 'Comida_jefe_G', 'martinez.team@email.com'),
/* --- martinez.team@email.com (isabel_m) --- */
(NULL, NULL, NULL, '50.00', 'mensual', 1, 0, 1, 'isabel_m', 'Gimnasio_isabel_G', 'martinez.team@email.com'),
('100.00', 'quincenal', 15, NULL, NULL, NULL, 0, 1, 'isabel_m', 'Freelance_isabel_I', 'martinez.team@email.com'),
/* --- martinez.team@email.com (raul_m) --- */
('80.00', 'mensual', 20, '100.00', 'mensual', 25, 0, 1, 'raul_m', 'Tecnologia_raul_G', 'martinez.team@email.com'),
(NULL, NULL, NULL, '100.00', 'mensual', 1, 0, 1, 'raul_m', 'Cursos_raul_G', 'martinez.team@email.com'),
/* --- martinez.team@email.com (nuria_m) --- */
(NULL, NULL, NULL, '25.00', 'diario', NULL, 0, 1, 'nuria_m', 'Cine_nuria_G', 'martinez.team@email.com'),
('150.00', 'mensual', 26, NULL, NULL, NULL, 0, 1, 'nuria_m', 'Beca_nuria_I', 'martinez.team@email.com'),

/* --- Familia: sanchez.clan@email.com (admin_sanchez - Admin) --- */
('950.00', 'mensual', 1, '1000.00', 'mensual', 5, 0, 1, 'admin_sanchez', 'Alquiler_admin_G', 'sanchez.clan@email.com'),
('3000.00', 'mensual', 5, NULL, NULL, NULL, 0, 1, 'admin_sanchez', 'Salario_Admin_S_I', 'sanchez.clan@email.com'),
('150.00', 'mensual', 10, '180.00', 'mensual', 15, 0, 1, 'admin_sanchez', 'Luz_Gas_admin_G', 'sanchez.clan@email.com'),
/* --- sanchez.clan@email.com (patricia_s) --- */
('250.00', 'quincenal', 15, '300.00', 'quincenal', 15, 0, 1, 'patricia_s', 'Comida_pat_G', 'sanchez.clan@email.com'),
('200.00', 'mensual', 9, NULL, NULL, NULL, 0, 1, 'patricia_s', 'Freelance_pat_I', 'sanchez.clan@email.com'),
/* --- sanchez.clan@email.com (axel_s) --- */
(NULL, NULL, NULL, '50.00', 'mensual', 20, 0, 1, 'axel_s', 'Tecnologia_axel_G', 'sanchez.clan@email.com'),
('100.00', 'mensual', 14, NULL, NULL, NULL, 0, 1, 'axel_s', 'Paga_axel_I', 'sanchez.clan@email.com'),
/* --- sanchez.clan@email.com (lidia_s) --- */
(NULL, NULL, NULL, '40.00', 'mensual', 5, 0, 1, 'lidia_s', 'Salud_lidia_G', 'sanchez.clan@email.com'),
('50.00', 'mensual', 15, NULL, NULL, NULL, 0, 1, 'lidia_s', 'Aportacion_lidia_I', 'sanchez.clan@email.com'),

/* --- Familia: gomez.group@email.com (jefe_gomez - Admin) --- */
('3400.00', 'mensual', 1, NULL, NULL, NULL, 0, 1, 'jefe_gomez', 'Salario_jefe_G_I', 'gomez.group@email.com'),
('750.00', 'mensual', 6, '800.00', 'mensual', 10, 0, 1, 'jefe_gomez', 'Alquiler_jefe_G_G', 'gomez.group@email.com'),
('150.00', 'mensual', 13, '130.00', 'mensual', 15, 0, 1, 'jefe_gomez', 'Agua_Internet_G', 'gomez.group@email.com'),
/* --- gomez.group@email.com (rocio_g) --- */
(NULL, NULL, NULL, '150.00', 'mensual', 5, 0, 1, 'rocio_g', 'Ropa_rocio_G', 'gomez.group@email.com'),
('250.00', 'quincenal', 10, NULL, NULL, NULL, 0, 1, 'rocio_g', 'Freelance_rocio_I', 'gomez.group@email.com'),
/* --- gomez.group@email.com (oscar_g) --- */
('100.00', 'mensual', 20, '150.00', 'mensual', 25, 0, 1, 'oscar_g', 'Tecnologia_oscar_G', 'gomez.group@email.com'),
('50.00', 'semanal', NULL, NULL, NULL, NULL, 0, 1, 'oscar_g', 'Paga_oscar_I', 'gomez.group@email.com'),
/* --- gomez.group@email.com (alba_g) --- */
(NULL, NULL, NULL, '50.00', 'mensual', 25, 0, 1, 'alba_g', 'Farmacia_alba_G', 'gomez.group@email.com'),
('200.00', 'mensual', 5, NULL, NULL, NULL, 0, 1, 'alba_g', 'Beca_alba_I', 'gomez.group@email.com'),

/* --- Familia: diaz.family@email.com (admin_diaz - Admin) --- */
('3700.00', 'mensual', 1, NULL, NULL, NULL, 0, 1, 'admin_diaz', 'Salario_admin_D_I', 'diaz.family@email.com'),
('850.00', 'mensual', 7, '900.00', 'mensual', 10, 0, 1, 'admin_diaz', 'Casa_admin_D_G', 'diaz.family@email.com'),
('100.00', 'mensual', 14, '120.00', 'mensual', 14, 0, 1, 'admin_diaz', 'Electricidad_admin_D_G', 'diaz.family@email.com'),
/* --- diaz.family@email.com (mercedes_d) --- */
('150.00', 'mensual', 4, '180.00', 'mensual', 10, 0, 1, 'mercedes_d', 'Ropa_mercedes_G', 'diaz.family@email.com'),
('300.00', 'mensual', 11, NULL, NULL, NULL, 0, 1, 'mercedes_d', 'Freelance_mercedes_I', 'diaz.family@email.com'),
('50.00', 'quincenal', 15, '60.00', 'quincenal', 15, 0, 1, 'mercedes_d', 'Comida_mercedes_G', 'diaz.family@email.com'),
/* --- diaz.family@email.com (victor_d) --- */
(NULL, NULL, NULL, '80.00', 'mensual', 25, 0, 1, 'victor_d', 'Deportes_victor_G', 'diaz.family@email.com'),
('50.00', 'semanal', NULL, NULL, NULL, NULL, 0, 1, 'victor_d', 'Paga_victor_I', 'diaz.family@email.com'),
/* --- diaz.family@email.com (andrea_d) --- */
('100.00', 'mensual', 1, '150.00', 'mensual', 15, 0, 1, 'andrea_d', 'Viajes_andrea_G', 'diaz.family@email.com'),
('200.00', 'mensual', 10, NULL, NULL, NULL, 0, 1, 'andrea_d', 'Regalo_navidad_andrea_I', 'diaz.family@email.com'),

/* --- Familia: alvarez.house@email.com (jefe_alvarez - Admin) --- */
('4000.00', 'mensual', 1, NULL, NULL, NULL, 0, 1, 'jefe_alvarez', 'Salario_jefe_A_I', 'alvarez.house@email.com'),
('700.00', 'mensual', 5, '750.00', 'mensual', 10, 0, 1, 'jefe_alvarez', 'Alquiler_jefe_A_G', 'alvarez.house@email.com'),
/* --- alvarez.house@email.com (rosa_a) --- */
(NULL, NULL, NULL, '100.00', 'mensual', 20, 0, 1, 'rosa_a', 'Ropa_rosa_G', 'alvarez.house@email.com'),
('100.00', 'mensual', 4, NULL, NULL, NULL, 0, 1, 'rosa_a', 'Aportacion_rosa_I', 'alvarez.house@email.com'),
/* --- alvarez.house@email.com (dani_a) --- */
(NULL, NULL, NULL, '100.00', 'mensual', 20, 0, 1, 'dani_a', 'Tecnologia_dani_G', 'alvarez.house@email.com'),
('50.00', 'semanal', NULL, NULL, NULL, NULL, 0, 1, 'dani_a', 'Paga_dani_I', 'alvarez.house@email.com'),
/* --- alvarez.house@email.com (eva_a) --- */
('150.00', 'quincenal', 15, '200.00', 'quincenal', 15, 0, 1, 'eva_a', 'Cursos_eva_A_G', 'alvarez.house@email.com'),
('120.00', 'mensual', 10, NULL, NULL, NULL, 0, 1, 'eva_a', 'Beca_eva_A_I', 'alvarez.house@email.com'),

/* --- Familia: moreno.co@email.com (admin_moreno - Admin) --- */
('3300.00', 'mensual', 1, NULL, NULL, NULL, 0, 1, 'admin_moreno', 'Salario_admin_M_I', 'moreno.co@email.com'),
('780.00', 'mensual', 5, '800.00', 'mensual', 10, 0, 1, 'admin_moreno', 'Casa_admin_M_G', 'moreno.co@email.com'),
/* --- moreno.co@email.com (marta_m) --- */
('100.00', 'mensual', 4, '150.00', 'mensual', 10, 0, 1, 'marta_m', 'Viajes_marta_G', 'moreno.co@email.com'),
('300.00', 'mensual', 11, NULL, NULL, NULL, 0, 1, 'marta_m', 'Freelance_marta_I', 'moreno.co@email.com'),
/* --- moreno.co@email.com (felix_m) --- */
(NULL, NULL, NULL, '100.00', 'mensual', 20, 0, 1, 'felix_m', 'Tecnologia_felix_G', 'moreno.co@email.com'),
('60.00', 'mensual', 8, NULL, NULL, NULL, 0, 1, 'felix_m', 'Regalo_felix_I', 'moreno.co@email.com'),
/* --- moreno.co@email.com (nerea_m) --- */
(NULL, NULL, NULL, '50.00', 'mensual', 15, 0, 1, 'nerea_m', 'Farmacia_nerea_G', 'moreno.co@email.com'),
('100.00', 'mensual', 10, NULL, NULL, NULL, 0, 1, 'nerea_m', 'Aportacion_nerea_I', 'moreno.co@email.com');