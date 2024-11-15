package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Función para calcular el tiempo transcurrido y formatearlo como [t]
func (t *TorreControl) tiempoTranscurrido() string {
	duracion := time.Since(t.inicio)
	return fmt.Sprintf("[%.2fs]", duracion.Seconds())
}

// Función para simular la variación en tiempo
func tiempoConVariacion(base int) time.Duration {
	variacion := base * VariacionTiempo / 100
	return time.Duration(base+rand.Intn(2*variacion)-variacion) * time.Millisecond
}

// Simulación de aterrizaje del avión en una goroutine
func (t *TorreControl) aterrizar(avion Avion) {
	// Registrar el tiempo desde que el avión se conecta a la torre de control
	startTimeConexion := time.Now()

	// Mensaje de conexión inicial
	t.mutex.Lock()
	fmt.Printf(msgConectado, t.tiempoTranscurrido(), avion.id)
	t.mutex.Unlock()

	// Intenta entrar en la cola de espera para aterrizar
	select {
	case t.esperaAterrizaje <- struct{}{}: // Entra en espera activa si hay espacio
		// Intenta acceder a una pista
		select {
		case t.pistas <- struct{}{}: // Si hay una pista disponible, la ocupa inmediatamente
			t.mutex.Lock()
			fmt.Printf(msgAterrizando, t.tiempoTranscurrido(), avion.id)
			t.mutex.Unlock()
		default:
			// Si no hay pista disponible, espera hasta que pueda acceder a una
			t.mutex.Lock()
			fmt.Printf(msgEsperandoAterrizar, t.tiempoTranscurrido(), avion.id)
			t.mutex.Unlock()
			t.pistas <- struct{}{} // Espera hasta que una pista esté libre
			t.mutex.Lock()
			fmt.Printf(msgAterrizando, t.tiempoTranscurrido(), avion.id)
			t.mutex.Unlock()
		}
		<-t.esperaAterrizaje // Sale de la espera activa
	default:
		// Si hay 5 aviones en espera, el avión permanece conectado pero fuera de espera activa
		t.mutex.Lock()
		fmt.Printf(msgEsperandoEntrada, t.tiempoTranscurrido(), avion.id)
		t.mutex.Unlock()

		// Bloquea hasta que pueda acceder a la cola de espera y luego intenta aterrizar
		t.esperaAterrizaje <- struct{}{}
		// Intenta acceder a una pista
		select {
		case t.pistas <- struct{}{}: // Si hay una pista disponible, la ocupa inmediatamente
			t.mutex.Lock()
			fmt.Printf(msgAterrizando, t.tiempoTranscurrido(), avion.id)
			t.mutex.Unlock()
		default:
			// Si no hay pista disponible, espera hasta que pueda acceder a una
			t.mutex.Lock()
			fmt.Printf(msgEsperandoAterrizar, t.tiempoTranscurrido(), avion.id)
			t.mutex.Unlock()
			t.pistas <- struct{}{} // Espera hasta que una pista esté libre
			t.mutex.Lock()
			fmt.Printf(msgAterrizando, t.tiempoTranscurrido(), avion.id)
			t.mutex.Unlock()
		}
		<-t.esperaAterrizaje
	}

	// Simula el tiempo de aterrizaje
	time.Sleep(tiempoConVariacion(TiempoAterrizaje))

	// Calcular tiempo desde conexión hasta aterrizaje
	tiempoConexionAterrizaje := time.Since(startTimeConexion).Seconds()
	t.mutex.Lock()
	fmt.Printf(msgAterrizado, t.tiempoTranscurrido(), avion.id, tiempoConexionAterrizaje)
	t.tiemposConexiónAterrizaje = append(t.tiemposConexiónAterrizaje, tiempoConexionAterrizaje)
	t.mutex.Unlock()

	// Intenta acceder a una puerta de desembarque
	select {
	case t.puertas <- struct{}{}: // Si hay una puerta disponible, entra directamente
		t.mutex.Lock()
		fmt.Printf(msgProcedePuerta, t.tiempoTranscurrido(), avion.id)
		t.mutex.Unlock()
	default:
		// Si todas las puertas están ocupadas, espera y muestra el mensaje
		t.mutex.Lock()
		fmt.Printf(msgEsperandoDesembarque, t.tiempoTranscurrido(), avion.id)
		t.mutex.Unlock()
		t.puertas <- struct{}{} // Espera hasta que una puerta esté libre

		// Mensaje de proceder a la puerta de desembarque
		t.mutex.Lock()
		fmt.Printf(msgProcedePuerta, t.tiempoTranscurrido(), avion.id)
		t.mutex.Unlock()
	}

	// Libera la pista y comienza el desembarque
	<-t.pistas
	go t.desembarcar(avion)
}

// Simulación de desembarque del avión en una goroutine
func (t *TorreControl) desembarcar(avion Avion) {
	// Registrar el tiempo desde el aterrizaje
	startTimeAterrizaje := time.Now()

	t.mutex.Lock()
	fmt.Printf(msgEnPuerta, t.tiempoTranscurrido(), avion.id)
	t.mutex.Unlock()

	// Simula el tiempo en la puerta de desembarque
	time.Sleep(tiempoConVariacion(TiempoPuerta))
	<-t.puertas // Libera la puerta al finalizar el desembarque

	// Calcular tiempo desde aterrizaje hasta finalización del desembarque
	tiempoAterrizajeDesembarque := time.Since(startTimeAterrizaje).Seconds()
	t.mutex.Lock()
	fmt.Printf(msgTerminoDesembarque, t.tiempoTranscurrido(), avion.id, tiempoAterrizajeDesembarque)
	t.tiemposAterrizajeDesembarque = append(t.tiemposAterrizajeDesembarque, tiempoAterrizajeDesembarque)
	t.mutex.Unlock()

	t.wg.Done() // Marca el avión como terminado
}

// Función para calcular el tiempo promedio de una lista de tiempos
func calcularPromedio(tiempos []float64) float64 {
	var suma float64
	for _, tiempo := range tiempos {
		suma += tiempo
	}
	return suma / float64(len(tiempos))
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Inicialización de la torre de control y recursos
	torre := TorreControl{
		pistas:           make(chan struct{}, NumPistas),
		puertas:          make(chan struct{}, NumPuertas),
		esperaAterrizaje: make(chan struct{}, MaxEspera),
		wg:               &sync.WaitGroup{},
		mutex:            sync.Mutex{},
		inicio:           time.Now(),
	}

	// Lanzamiento de goroutines para cada avión
	for i := 0; i < NumAviones; i++ {
		torre.wg.Add(1)
		avion := Avion{id: i}
		go torre.aterrizar(avion)
	}

	// Espera a que todos los aviones terminen su proceso
	torre.wg.Wait()

	// Calcular y mostrar los promedios
	fmt.Printf("\nSimulación completada\n")
	fmt.Printf("Promedio de tiempo desde conexión hasta aterrizaje: %.2f segundos\n", calcularPromedio(torre.tiemposConexiónAterrizaje))
	fmt.Printf("Promedio de tiempo desde aterrizaje hasta desembarque: %.2f segundos\n", calcularPromedio(torre.tiemposAterrizajeDesembarque))
}
