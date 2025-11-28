package controllers

import (
	"arka-code/entities"
	"arka-code/models"
	"arka-code/utils"
	"log"
	"net/http"
	"time"
)

// CTRL.3 Gestor Balance
// CU-003 Consultar Balance
type BalanceController struct{}

var BalanceControllerInstance = &BalanceController{}

// BalanceUsuario estructura para mostrar balance financiero de cada usuario
// Campos:
// - NombreUsuario: Identificador único del usuario
// - NombrePersonal: Nombre real del usuario para mostrar
// - BalanceMensual: Diferencia entre ingresos y gastos del mes
// - BalanceAnual: Diferencia entre ingresos y gastos del año
// - Movimientos: Lista de movimientos del usuario en el período
// - TotalIngresos: Suma total de ingresos del usuario
// - TotalGastos: Suma total de gastos del usuario
// Uso: Presentación de balances individuales en vistas administrativas
type BalanceUsuario struct {
	NombreUsuario  string
	NombrePersonal string
	BalanceMensual float64
	BalanceAnual   float64
	Movimientos    []entities.Movimiento
	TotalIngresos  float64
	TotalGastos    float64
}

// DatosGrafico estructura para datos visualizables en gráficos
// Campos:
// - Etiqueta: Nombre a mostrar en el gráfico (concepto o usuario)
// - Monto: Valor monetario representado
// - Porcentaje: Porcentaje relativo al total del grupo
// - Color: Color identificador para el elemento del gráfico
// Uso: Generación de datos para gráficos de torta y barras
type DatosGrafico struct {
	Etiqueta   string
	Monto      float64
	Porcentaje float64
	Color      string
}

// FNBalance-Index
// Index muestra la página principal de Balance con diferentes vistas y filtros
// Parámetros:
// - w: ResponseWriter para enviar respuesta HTTP
// - r: Request HTTP con parámetros de filtrado
// Flujo:
// 1. Verifica autenticación del usuario
// 2. Obtiene y procesa parámetros de URL (tipo, gráficos, fechas)
// 3. Configura rangos de fecha según período seleccionado
// 4. Carga datos según el tipo de vista (balance o gráficos)
// 5. Renderiza la plantilla con los datos procesados
// Uso: Punto de entrada principal para el módulo de balances
func (c *BalanceController) Index(w http.ResponseWriter, r *http.Request) {
	sessionData, ok := utils.GetSessionData(r)
	if !ok {
		http.Redirect(w, r, "/seleccionar-perfil", http.StatusSeeOther)
		return
	}

	// Obtener parámetros de URL
	tipo := r.URL.Query().Get("tipo")
	if tipo == "" {
		tipo = "balance"
	}

	tipoGrafico := r.URL.Query().Get("tipoGrafico")
	if tipoGrafico == "" {
		tipoGrafico = "barras"
	}

	vistaGrafico := r.URL.Query().Get("vistaGrafico")
	if vistaGrafico == "" {
		vistaGrafico = "concepto"
	}

	usuarioFiltro := r.URL.Query().Get("usuario")

	// ===== NUEVAS LÍNEAS: Capturar parámetros de fecha =====
	periodo := r.URL.Query().Get("periodo")
	if periodo == "" {
		periodo = "mes" // Por defecto mes actual
	}

	var inicioRango, finRango time.Time
	now := time.Now()

	if periodo == "mes" {
		// Filtrar por mes específico
		mesStr := r.URL.Query().Get("mes")
		if mesStr != "" {
			// Parsear formato YYYY-MM
			mesDate, err := time.Parse("2006-01", mesStr)
			if err == nil {
				inicioRango = time.Date(mesDate.Year(), mesDate.Month(), 1, 0, 0, 0, 0, mesDate.Location())
				finRango = inicioRango.AddDate(0, 1, 0).Add(-time.Second)
			} else {
				// Si hay error, usar mes actual
				inicioRango = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
				finRango = inicioRango.AddDate(0, 1, 0).Add(-time.Second)
			}
		} else {
			// Mes actual por defecto
			inicioRango = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
			finRango = inicioRango.AddDate(0, 1, 0).Add(-time.Second)
		}
	} else if periodo == "rango" {
		// Filtrar por rango personalizado
		fechaInicioStr := r.URL.Query().Get("fechaInicio")
		fechaFinStr := r.URL.Query().Get("fechaFin")

		if fechaInicioStr != "" && fechaFinStr != "" {
			fechaInicio, err1 := time.Parse("2006-01-02", fechaInicioStr)
			fechaFin, err2 := time.Parse("2006-01-02", fechaFinStr)

			if err1 == nil && err2 == nil {
				inicioRango = time.Date(fechaInicio.Year(), fechaInicio.Month(), fechaInicio.Day(), 0, 0, 0, 0, fechaInicio.Location())
				finRango = time.Date(fechaFin.Year(), fechaFin.Month(), fechaFin.Day(), 23, 59, 59, 0, fechaFin.Location())
			} else {
				// Si hay error, usar mes actual
				inicioRango = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
				finRango = inicioRango.AddDate(0, 1, 0).Add(-time.Second)
			}
		} else {
			// Por defecto últimos 30 días
			inicioRango = now.AddDate(0, 0, -30)
			finRango = now
		}
	} else {
		// Por defecto mes actual
		inicioRango = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		finRango = inicioRango.AddDate(0, 1, 0).Add(-time.Second)
	}

	log.Printf("📅 Período seleccionado: %s - Desde: %s Hasta: %s",
		periodo, inicioRango.Format("2006-01-02"), finRango.Format("2006-01-02"))
	// ===== FIN NUEVAS LÍNEAS =====

	// Preparar datos base
	data := map[string]interface{}{
		"Title":        "Balance",
		"CurrentPage":  "balance",
		"SessionData":  sessionData,
		"Tipo":         tipo,
		"TipoGrafico":  tipoGrafico,
		"VistaGrafico": vistaGrafico,
		"MesActual":    now.Format("2006-01"), // Para el input type="month"
	}

	if tipo == "balance" {
		c.cargarDatosBalance(w, r, sessionData, usuarioFiltro, inicioRango, finRango, data)
	} else if tipo == "graficos" {
		c.cargarDatosGraficos(w, r, sessionData, usuarioFiltro, tipoGrafico, vistaGrafico, inicioRango, finRango, data)
	}

	utils.RenderTemplate(w, "dashboard", "balance/index", data)
}

// FNBalance-cargarDatosBalance
// cargarDatosBalance carga y procesa los datos para la vista de balance tabular
// Parámetros:
// - w: ResponseWriter para manejar errores HTTP
// - r: Request HTTP (no utilizado actualmente)
// - sessionData: Datos de sesión del usuario autenticado
// - usuarioFiltro: Usuario específico a filtrar (opcional)
// - inicioRango, finRango: Rango de fechas para el balance
// - data: Mapa de datos que se populate con la información del balance
// Flujo:
// 1. Define rangos para balance mensual y anual
// 2. Según el rol del usuario, carga balances individuales o familiares
// 3. Calcula totales y balances generales
// 4. Agrega datos al mapa para renderizado
// Uso: Preparación de datos para vista tabular de balances
func (c *BalanceController) cargarDatosBalance(w http.ResponseWriter, _ *http.Request, sessionData *entities.SessionData, usuarioFiltro string, inicioRango, finRango time.Time, data map[string]interface{}) {
	// Usar inicioRango y finRango en lugar de fechas fijas
	inicioMes := inicioRango
	finMes := finRango

	// Para el balance anual, usar el año completo
	inicioAnio := time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.Now().Location())
	finAnio := time.Date(time.Now().Year(), 12, 31, 23, 59, 59, 0, time.Now().Location())

	var balances []BalanceUsuario
	var movimientosIngresos []entities.Movimiento
	var movimientosGastos []entities.Movimiento
	var totalIngresos, totalGastos float64

	if sessionData.Rol == 1 {
		// ADMINISTRADOR
		usuarios, err := models.UserModelInstance.FindByFamilia(sessionData.CorreoFamilia)
		if err != nil {
			log.Printf("❌ Error obteniendo usuarios: %v", err)
			http.Error(w, "Error al cargar usuarios", http.StatusInternalServerError)
			return
		}

		data["UsuariosFamilia"] = usuarios

		// Si hay filtro de usuario, solo mostrar ese
		if usuarioFiltro != "" {
			balance := c.calcularBalanceUsuarioCompleto(usuarioFiltro, sessionData.CorreoFamilia, inicioMes, finMes, inicioAnio, finAnio)
			balances = append(balances, balance)
			movimientosIngresos = c.filtrarMovimientosPorTipo(balance.Movimientos, 1, sessionData.CorreoFamilia)
			movimientosGastos = c.filtrarMovimientosPorTipo(balance.Movimientos, 0, sessionData.CorreoFamilia)
			totalIngresos = balance.TotalIngresos
			totalGastos = balance.TotalGastos
			data["UsuarioFiltro"] = usuarioFiltro
		} else {
			// Mostrar todos los usuarios
			for _, usuario := range usuarios {
				balance := c.calcularBalanceUsuarioCompleto(usuario.NombreUsuario, sessionData.CorreoFamilia, inicioMes, finMes, inicioAnio, finAnio)
				balances = append(balances, balance)

				// Acumular movimientos y totales
				ingresosUsuario := c.filtrarMovimientosPorTipo(balance.Movimientos, 1, sessionData.CorreoFamilia)
				gastosUsuario := c.filtrarMovimientosPorTipo(balance.Movimientos, 0, sessionData.CorreoFamilia)
				movimientosIngresos = append(movimientosIngresos, ingresosUsuario...)
				movimientosGastos = append(movimientosGastos, gastosUsuario...)
				totalIngresos += balance.TotalIngresos
				totalGastos += balance.TotalGastos
			}
		}
	} else {
		// MIEMBRO - Solo su propio balance
		balance := c.calcularBalanceUsuarioCompleto(sessionData.NombreUsuario, sessionData.CorreoFamilia, inicioMes, finMes, inicioAnio, finAnio)
		balances = append(balances, balance)
		movimientosIngresos = c.filtrarMovimientosPorTipo(balance.Movimientos, 1, sessionData.CorreoFamilia)
		movimientosGastos = c.filtrarMovimientosPorTipo(balance.Movimientos, 0, sessionData.CorreoFamilia)
		totalIngresos = balance.TotalIngresos
		totalGastos = balance.TotalGastos
	}

	data["Balances"] = balances
	data["MovimientosIngresos"] = movimientosIngresos
	data["MovimientosGastos"] = movimientosGastos
	data["TotalIngresos"] = totalIngresos
	data["TotalGastos"] = totalGastos
	data["BalanceMensual"] = totalIngresos - totalGastos

	// Calcular balance anual total
	var totalIngresosAnual, totalGastosAnual float64
	for _, balance := range balances {
		totalIngresosAnual += balance.TotalIngresos
		totalGastosAnual += balance.TotalGastos
	}
	data["BalanceAnual"] = totalIngresosAnual - totalGastosAnual
}

// FNBalance-cargarDatosGraficos
// cargarDatosGraficos carga los datos para la vista de gráficos estadísticos
// Parámetros:
// - sessionData: Datos de sesión del usuario autenticado
// - usuarioFiltro: Usuario específico a filtrar (opcional)
// - vistaGrafico: Tipo de vista ("concepto" o "miembro")
// - inicioRango, finRango: Rango de fechas para los datos
// - data: Mapa de datos que se populate con información de gráficos
// Flujo:
// - Según la vista seleccionada, carga datos agrupados por concepto o miembro
// - Solo permite vista "miembro" para administradores
// - Agrega datos de ingresos y gastos al mapa
// Uso: Preparación de datos para vistas gráficas del balance
func (c *BalanceController) cargarDatosGraficos(_ http.ResponseWriter, _ *http.Request, sessionData *entities.SessionData, usuarioFiltro, _, vistaGrafico string, inicioRango, finRango time.Time, data map[string]interface{}) {
	// Usar inicioRango y finRango directamente
	if vistaGrafico == "concepto" {
		datosIngresos, datosGastos := c.obtenerDatosPorConcepto(sessionData, usuarioFiltro, inicioRango, finRango)
		data["IngresosData"] = datosIngresos
		data["GastosData"] = datosGastos
	} else if vistaGrafico == "miembro" && sessionData.Rol == 1 {
		datosIngresos, datosGastos := c.obtenerDatosPorMiembro(sessionData, inicioRango, finRango)
		data["IngresosData"] = datosIngresos
		data["GastosData"] = datosGastos
	}
}

// FNBalance-calcularBalanceUsuarioCompleto
// calcularBalanceUsuarioCompleto calcula el balance financiero completo de un usuario
// Parámetros:
// - nombreUsuario: Identificador del usuario
// - correoFamilia: Familia a la que pertenece el usuario
// - inicioMes, finMes: Rango para balance mensual
// - inicioAnio, finAnio: Rango para balance anual
// Retorno: Estructura BalanceUsuario con todos los datos calculados
// Flujo:
// 1. Obtiene nombre personal del usuario
// 2. Calcula ingresos y gastos para el rango mensual
// 3. Calcula ingresos y gastos para el año completo
// 4. Obtiene movimientos detallados del período mensual
// Uso: Cálculo completo de balance individual para reportes
func (c *BalanceController) calcularBalanceUsuarioCompleto(nombreUsuario, correoFamilia string, inicioMes, finMes, inicioAnio, finAnio time.Time) BalanceUsuario {
	balance := BalanceUsuario{
		NombreUsuario: nombreUsuario,
	}

	// Obtener nombre personal del usuario
	usuario, err := models.UserModelInstance.FindByNombreUsuario(nombreUsuario)
	if err == nil && usuario != nil {
		balance.NombrePersonal = usuario.NombrePersonal
	}

	// ✅ USAR EL RANGO SELECCIONADO (inicioMes, finMes)
	ingresosMes, gastosMes, _ := models.MovimientoModelInstance.GetTotalesByUsuarioAndPeriod(
		nombreUsuario, correoFamilia, inicioMes, finMes,
	)
	balance.BalanceMensual = ingresosMes - gastosMes
	balance.TotalIngresos = ingresosMes
	balance.TotalGastos = gastosMes

	// Balance anual (año completo)
	ingresosAnio, gastosAnio, _ := models.MovimientoModelInstance.GetTotalesByUsuarioAndPeriod(
		nombreUsuario, correoFamilia, inicioAnio, finAnio,
	)
	balance.BalanceAnual = ingresosAnio - gastosAnio

	// ✅ MOVIMIENTOS DEL RANGO SELECCIONADO
	movimientos, _ := models.MovimientoModelInstance.FindByUsuarioAndPeriod(
		nombreUsuario, correoFamilia, inicioMes, finMes,
	)
	balance.Movimientos = movimientos

	return balance
}

// FNBalance-filtrarMovimientosPorTipo
// filtrarMovimientosPorTipo filtra una lista de movimientos por tipo (ingreso/gasto)
// Parámetros:
// - movimientos: Lista de movimientos a filtrar
// - tipo: Tipo de movimiento (0 = gasto, 1 = ingreso)
// - correoFamilia: Familia para verificar el tipo de concepto
// Retorno: Lista filtrada de movimientos del tipo especificado
// Flujo:
// - Itera sobre movimientos y verifica el tipo de concepto asociado
// - Usa FindByNombre para obtener información completa del concepto
// - Solo incluye movimientos cuyo concepto coincida con el tipo
// Uso: Separación de movimientos para visualización por tipo
func (c *BalanceController) filtrarMovimientosPorTipo(movimientos []entities.Movimiento, tipo int8, correoFamilia string) []entities.Movimiento {
	var resultado []entities.Movimiento
	for _, mov := range movimientos {
		concepto, err := models.ConceptoModelInstance.FindByNombre(mov.NombreConcepto, correoFamilia)
		if err == nil && concepto != nil && concepto.Tipo == tipo {
			resultado = append(resultado, mov)
		}
	}
	return resultado
}

// FNBalance-obtenerDatosPorConcepto
// obtenerDatosPorConcepto genera datos para gráficos agrupados por concepto
// Parámetros:
// - sessionData: Datos de sesión del usuario
// - usuarioFiltro: Usuario específico para filtrar (opcional)
// - inicio, fin: Rango de fechas para los datos
// Retorno: Dos listas de DatosGrafico para ingresos y gastos
// Flujo:
// 1. Obtiene todos los conceptos de la familia
// 2. Calcula totales por concepto en el período
// 3. Calcula porcentajes relativos
// 4. Asigna colores y prepara datos para gráficos
// Uso: Generación de datos para gráficos de torta por categoría
func (c *BalanceController) obtenerDatosPorConcepto(sessionData *entities.SessionData, usuarioFiltro string, inicio, fin time.Time) ([]DatosGrafico, []DatosGrafico) {
	var datosIngresos []DatosGrafico
	var datosGastos []DatosGrafico

	conceptos, err := models.ConceptoModelInstance.FindAllByFamilia(sessionData.CorreoFamilia)
	if err != nil {
		log.Printf("❌ Error obteniendo conceptos: %v", err)
		return datosIngresos, datosGastos
	}

	log.Printf("🔍 DEBUG - Total conceptos encontrados: %d", len(conceptos))

	var totalIngresos, totalGastos float64

	// Primero calcular totales y debug de conceptos
	for i, concepto := range conceptos {
		monto, err := models.MovimientoModelInstance.GetTotalByConceptoAndPeriod(
			concepto.NombreConcepto, sessionData.CorreoFamilia, inicio, fin, usuarioFiltro,
		)
		if err != nil {
			log.Printf("❌ Error obteniendo total por concepto: %v", err)
			continue
		}

		log.Printf("🔍 DEBUG - Concepto %d: %s, Tipo: %d, Monto: %.2f",
			i, concepto.NombreConcepto, concepto.Tipo, monto)

		if monto > 0 {
			if concepto.Tipo == 1 {
				totalIngresos += monto
				log.Printf("  → Sumando a INGRESOS: %.2f", monto)
			} else if concepto.Tipo == 0 {
				totalGastos += monto
				log.Printf("  → Sumando a GASTOS: %.2f", monto)
			} else {
				log.Printf("  ⚠️ Tipo desconocido: %d", concepto.Tipo)
			}
		}
	}

	log.Printf("🔍 DEBUG - Totales calculados - Ingresos: %.2f, Gastos: %.2f", totalIngresos, totalGastos)

	// Si no hay datos, retornar arrays vacíos
	if totalIngresos == 0 && totalGastos == 0 {
		log.Printf("⚠️ No hay datos para gráficos")
		return datosIngresos, datosGastos
	}

	// Calcular porcentajes y crear datos
	colores := []string{"#4f46e5", "#22c55e", "#ef4444", "#f59e0b", "#8b5cf6", "#06b6d4", "#84cc16", "#f97316"}
	colorIndex := 0

	for _, concepto := range conceptos {
		monto, err := models.MovimientoModelInstance.GetTotalByConceptoAndPeriod(
			concepto.NombreConcepto, sessionData.CorreoFamilia, inicio, fin, usuarioFiltro,
		)
		if err != nil {
			continue
		}

		if monto > 0 {
			color := colores[colorIndex%len(colores)]
			if concepto.Color != nil && *concepto.Color != "" {
				color = *concepto.Color
			}

			dato := DatosGrafico{
				Etiqueta: concepto.NombreConcepto,
				Monto:    monto,
				Color:    color,
			}

			if concepto.Tipo == 1 && totalIngresos > 0 {
				dato.Porcentaje = (monto / totalIngresos) * 100
				datosIngresos = append(datosIngresos, dato)
				log.Printf("📊 Añadiendo a INGRESOS: %s - %.2f (%.1f%%)",
					concepto.NombreConcepto, monto, dato.Porcentaje)
			} else if concepto.Tipo == 0 && totalGastos > 0 {
				dato.Porcentaje = (monto / totalGastos) * 100
				datosGastos = append(datosGastos, dato)
				log.Printf("📊 Añadiendo a GASTOS: %s - %.2f (%.1f%%)",
					concepto.NombreConcepto, monto, dato.Porcentaje)
			}

			colorIndex++
		}
	}

	log.Printf("🔍 DEBUG - Resultado final - DatosIngresos: %d, DatosGastos: %d",
		len(datosIngresos), len(datosGastos))

	return datosIngresos, datosGastos
}

// FNBalance-obtenerDatosPorMiembro
// obtenerDatosPorMiembro genera datos para gráficos agrupados por miembro familiar
// Parámetros:
// - sessionData: Datos de sesión del usuario (debe ser administrador)
// - inicio, fin: Rango de fechas para los datos
// Retorno: Dos listas de DatosGrafico para ingresos y gastos por miembro
// Flujo:
// 1. Obtiene todos los usuarios de la familia
// 2. Calcula totales por usuario en el período
// 3. Calcula porcentajes relativos
// 4. Asigna colores únicos por usuario
// Uso: Generación de datos para gráficos de distribución familiar
func (c *BalanceController) obtenerDatosPorMiembro(sessionData *entities.SessionData, inicio, fin time.Time) ([]DatosGrafico, []DatosGrafico) {
	var datosIngresos []DatosGrafico
	var datosGastos []DatosGrafico

	usuarios, _ := models.UserModelInstance.FindByFamilia(sessionData.CorreoFamilia)

	var totalIngresos, totalGastos float64

	colores := []string{"#4f46e5", "#22c55e", "#ef4444", "#f59e0b", "#8b5cf6", "#06b6d4"}

	for _, usuario := range usuarios {
		gastos, ingresos, _ := models.MovimientoModelInstance.GetTotalesByUsuarioAndPeriod(
			usuario.NombreUsuario, sessionData.CorreoFamilia, inicio, fin,
		)
		totalIngresos += ingresos
		totalGastos += gastos
	}

	for i, usuario := range usuarios {
		gastos, ingresos, _ := models.MovimientoModelInstance.GetTotalesByUsuarioAndPeriod(
			usuario.NombreUsuario, sessionData.CorreoFamilia, inicio, fin,
		)

		color := colores[i%len(colores)]

		if ingresos > 0 {
			datosIngresos = append(datosIngresos, DatosGrafico{
				Etiqueta:   usuario.NombrePersonal,
				Monto:      ingresos,
				Porcentaje: (ingresos / totalIngresos) * 100,
				Color:      color,
			})
		}

		if gastos > 0 {
			datosGastos = append(datosGastos, DatosGrafico{
				Etiqueta:   usuario.NombrePersonal,
				Monto:      gastos,
				Porcentaje: (gastos / totalGastos) * 100,
				Color:      color,
			})
		}
	}

	return datosIngresos, datosGastos
}
