package scheduler

import (
	"arka-code/models"
	"context"
	"log"
	"time"
)

type Scheduler struct {
	ticker *time.Ticker
}

func NewScheduler() *Scheduler {
	return &Scheduler{}
}

func (s *Scheduler) Start(ctx context.Context) {
	log.Println("🕐 Scheduler iniciado - Verificando movimientos automáticos cada hora")

	// Ejecutar inmediatamente al iniciar
	go s.ejecutarTareasAutomaticas()

	// Configurar ticker cada hora (puedes ajustar según necesidades)
	s.ticker = time.NewTicker(2 * time.Second)

	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.ejecutarTareasAutomaticas()
			case <-ctx.Done():
				log.Println("🛑 Scheduler recibió señal de cierre")
				s.Stop()
				return
			}
		}
	}()
}

func (s *Scheduler) ejecutarTareasAutomaticas() {
	log.Println("🔄 Ejecutando tareas automáticas...")

	// Tarea 1: Crear movimientos para períodos diarios
	if err := s.crearMovimientosPeriodicos(); err != nil {
		log.Printf("❌ Error en tarea de movimientos periódicos: %v", err)
	}

	// Tarea 2: Crear movimientos para días específicos
	if err := s.crearMovimientosDiasEspecificos(); err != nil {
		log.Printf("❌ Error en tarea de movimientos por día específico: %v", err)
	}

	log.Println("✅ Tareas automáticas completadas")
}

// crearMovimientosPeriodicos crea movimientos para períodos diarios
func (s *Scheduler) crearMovimientosPeriodicos() error {
	log.Println("📅 Creando movimientos para períodos diarios...")

	// Obtener personalizaciones con período diario
	personalizaciones, err := models.PersonalizacionModelInstance.FindByPeriodoDiario()
	if err != nil {
		return err
	}

	log.Printf("📋 Personalizaciones diarias encontradas: %d", len(personalizaciones))

	movimientosCreados := 0
	movimientosExistentes := 0
	movimientosError := 0

	for _, personalizacion := range personalizaciones {
		// Verificar si ya existe un movimiento para hoy
		existe, err := models.MovimientoModelInstance.ExisteMovimientoHoy(
			personalizacion.NombreUsuario,
			personalizacion.NombreConcepto,
			personalizacion.CorreoFamilia,
		)
		if err != nil {
			log.Printf("⚠️  Error verificando movimiento existente: %v", err)
			movimientosError++
			continue
		}

		if existe {
			log.Printf("⏭️  Ya existe movimiento hoy para %s - %s",
				personalizacion.NombreUsuario, personalizacion.NombreConcepto)
			movimientosExistentes++
			continue
		}

		// Verificar que MontoPlanificado no sea nil
		if personalizacion.MontoPlanificado == nil {
			log.Printf("⚠️  Monto planificado es nil para %s - %s",
				personalizacion.NombreUsuario, personalizacion.NombreConcepto)
			movimientosError++
			continue
		}

		// Crear el movimiento automático
		err = models.MovimientoModelInstance.CreateAutomatico(
			personalizacion.NombreUsuario,
			personalizacion.NombreConcepto,
			personalizacion.CorreoFamilia,
			*personalizacion.MontoPlanificado,
			"Movimiento automático diario",
		)
		if err != nil {
			log.Printf("❌ Error creando movimiento automático para %s: %v",
				personalizacion.NombreUsuario, err)
			movimientosError++
			continue
		}

		log.Printf("✅ Movimiento diario creado: %s - %s - S/ %.2f",
			personalizacion.NombreUsuario, personalizacion.NombreConcepto, *personalizacion.MontoPlanificado)
		movimientosCreados++
	}

	log.Printf("📊 Resumen diarios: Creados=%d, Existentes=%d, Errores=%d",
		movimientosCreados, movimientosExistentes, movimientosError)

	return nil
}

// crearMovimientosDiasEspecificos crea movimientos para días específicos del mes
func (s *Scheduler) crearMovimientosDiasEspecificos() error {
	log.Println("📆 Creando movimientos para días específicos...")

	diaActual := time.Now().Day()
	log.Printf("📅 Día actual: %d", diaActual)

	// Obtener personalizaciones que tienen día específico planificado
	personalizaciones, err := models.PersonalizacionModelInstance.FindByDiaPlanificado(int8(diaActual))
	if err != nil {
		return err
	}

	log.Printf("📋 Personalizaciones para día %d encontradas: %d", diaActual, len(personalizaciones))

	movimientosCreados := 0
	movimientosExistentes := 0
	movimientosError := 0

	for _, personalizacion := range personalizaciones {
		// Verificar si ya existe un movimiento para hoy
		existe, err := models.MovimientoModelInstance.ExisteMovimientoHoy(
			personalizacion.NombreUsuario,
			personalizacion.NombreConcepto,
			personalizacion.CorreoFamilia,
		)
		if err != nil {
			log.Printf("⚠️  Error verificando movimiento existente: %v", err)
			movimientosError++
			continue
		}

		if existe {
			log.Printf("⏭️  Ya existe movimiento hoy para %s - %s",
				personalizacion.NombreUsuario, personalizacion.NombreConcepto)
			movimientosExistentes++
			continue
		}

		// Verificar que MontoPlanificado no sea nil
		if personalizacion.MontoPlanificado == nil {
			log.Printf("⚠️  Monto planificado es nil para %s - %s",
				personalizacion.NombreUsuario, personalizacion.NombreConcepto)
			movimientosError++
			continue
		}

		// Crear el movimiento automático
		err = models.MovimientoModelInstance.CreateAutomatico(
			personalizacion.NombreUsuario,
			personalizacion.NombreConcepto,
			personalizacion.CorreoFamilia,
			*personalizacion.MontoPlanificado,
			"Movimiento automático mensual",
		)
		if err != nil {
			log.Printf("❌ Error creando movimiento automático para %s: %v",
				personalizacion.NombreUsuario, err)
			movimientosError++
			continue
		}

		log.Printf("✅ Movimiento mensual creado: %s - %s - S/ %.2f",
			personalizacion.NombreUsuario, personalizacion.NombreConcepto, *personalizacion.MontoPlanificado)
		movimientosCreados++
	}

	log.Printf("📊 Resumen días específicos: Creados=%d, Existentes=%d, Errores=%d",
		movimientosCreados, movimientosExistentes, movimientosError)

	return nil
}

func (s *Scheduler) Stop() {
	log.Println("🛑 Deteniendo scheduler...")
	if s.ticker != nil {
		s.ticker.Stop()
	}
	log.Println("✅ Scheduler detenido")
}
