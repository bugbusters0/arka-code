package controllers

import (
	"arka-code/entities"
	"arka-code/models"
	"arka-code/utils"
	"log"
	"net/http"
	"time"
)

type BalanceController struct{}

var BalanceControllerInstance = &BalanceController{}

// BalanceUsuario estructura para mostrar balance de cada usuario
type BalanceUsuario struct {
	NombreUsuario  string
	NombrePersonal string
	BalanceMensual float64
	BalanceAnual   float64
	Movimientos    []entities.Movimiento
	TotalIngresos  float64
	TotalGastos    float64
}

// DatosGrafico estructura para datos de gráficos
type DatosGrafico struct {
	Etiqueta   string
	Monto      float64
	Porcentaje float64
	Color      string
}

// Index muestra la página principal de Balance
func (c *BalanceController) Index(w http.ResponseWriter, r *http.Request) {
	sessionData, ok := utils.GetSessionData(r)
	if !ok {
		http.Redirect(w, r, "/seleccionar-perfil", http.StatusSeeOther)
		return
	}

	// Obtener parámetros de URL
	tipo := r.URL.Query().Get("tipo")
	if tipo == "" {
		tipo = "balance" // Por defecto mostrar balance
	}

	tipoGrafico := r.URL.Query().Get("tipoGrafico")
	if tipoGrafico == "" {
		tipoGrafico = "barras" // Por defecto gráfico de barras
	}

	vistaGrafico := r.URL.Query().Get("vistaGrafico")
	if vistaGrafico == "" {
		vistaGrafico = "concepto" // Por defecto por concepto
	}

	usuarioFiltro := r.URL.Query().Get("usuario")

	// Preparar datos base
	data := map[string]interface{}{
		"Title":        "Balance",
		"CurrentPage":  "balance",
		"SessionData":  sessionData,
		"Tipo":         tipo,
		"TipoGrafico":  tipoGrafico,
		"VistaGrafico": vistaGrafico,
	}

	if tipo == "balance" {
		c.cargarDatosBalance(w, r, sessionData, usuarioFiltro, data)
	} else if tipo == "graficos" {
		c.cargarDatosGraficos(w, r, sessionData, usuarioFiltro, tipoGrafico, vistaGrafico, data)
	}

	utils.RenderTemplate(w, "dashboard", "balance/index", data)
}

// cargarDatosBalance carga los datos para la vista de balance
func (c *BalanceController) cargarDatosBalance(w http.ResponseWriter, _ *http.Request, sessionData *entities.SessionData, usuarioFiltro string, data map[string]interface{}) {
	// Calcular fechas para mes y año actual
	inicioMes := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Now().Location())
	finMes := inicioMes.AddDate(0, 1, 0).Add(-time.Second)

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

// cargarDatosGraficos carga los datos para la vista de gráficos
func (c *BalanceController) cargarDatosGraficos(_ http.ResponseWriter, _ *http.Request, sessionData *entities.SessionData, usuarioFiltro, _, vistaGrafico string, data map[string]interface{}) {
	inicioMes := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Now().Location())
	finMes := inicioMes.AddDate(0, 1, 0).Add(-time.Second)

	if vistaGrafico == "concepto" {
		// Gráficos por concepto
		datosIngresos, datosGastos := c.obtenerDatosPorConcepto(sessionData, usuarioFiltro, inicioMes, finMes)
		data["IngresosData"] = datosIngresos
		data["GastosData"] = datosGastos
	} else if vistaGrafico == "miembro" && sessionData.Rol == 1 {
		// Gráficos por miembro (solo admin)
		datosIngresos, datosGastos := c.obtenerDatosPorMiembro(sessionData, inicioMes, finMes)
		data["IngresosData"] = datosIngresos
		data["GastosData"] = datosGastos
	}

}

// calcularBalanceUsuarioCompleto calcula el balance completo de un usuario
func (c *BalanceController) calcularBalanceUsuarioCompleto(nombreUsuario, correoFamilia string, inicioMes, finMes, inicioAnio, finAnio time.Time) BalanceUsuario {
	balance := BalanceUsuario{
		NombreUsuario: nombreUsuario,
	}

	// Obtener nombre personal del usuario
	usuario, err := models.UserModelInstance.FindByNombreUsuario(nombreUsuario)
	if err == nil && usuario != nil {
		balance.NombrePersonal = usuario.NombrePersonal
	}

	// Balance mensual
	gastosMes, ingresosMes, _ := models.MovimientoModelInstance.GetTotalesByUsuarioAndPeriod(
		nombreUsuario, correoFamilia, inicioMes, finMes,
	)
	balance.BalanceMensual = ingresosMes - gastosMes
	balance.TotalIngresos = ingresosMes
	balance.TotalGastos = gastosMes

	// Balance anual
	gastosAnio, ingresosAnio, _ := models.MovimientoModelInstance.GetTotalesByUsuarioAndPeriod(
		nombreUsuario, correoFamilia, inicioAnio, finAnio,
	)
	balance.BalanceAnual = ingresosAnio - gastosAnio

	// Movimientos del mes
	movimientos, _ := models.MovimientoModelInstance.FindByUsuarioAndPeriod(
		nombreUsuario, correoFamilia, inicioMes, finMes,
	)
	balance.Movimientos = movimientos

	return balance
}

// filtrarMovimientosPorTipo filtra movimientos por tipo (ingreso/gasto)
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

// obtenerDatosPorConcepto obtiene datos agrupados por concepto
// obtenerDatosPorConcepto obtiene datos agrupados por concepto
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

// obtenerDatosPorMiembro obtiene datos agrupados por miembro
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
