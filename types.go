package main

import (
	"sync"
	"time"
)

// Estructura para representar un avión
type Avion struct {
	id int
}

// Torre de control para coordinar el aterrizaje
type TorreControl struct {
	pistas           chan struct{} // Canales para controlar acceso a pistas
	puertas          chan struct{} // Canales para controlar acceso a puertas
	esperaAterrizaje chan struct{} // Canal para limitar la espera activa a MaxEspera
	wg               *sync.WaitGroup
	mutex            sync.Mutex // Mutex para sincronizar todos los mensajes críticos
	inicio           time.Time  // Tiempo de inicio del programa

	// Nuevos campos para almacenar los tiempos específicos solicitados
	tiemposConexiónAterrizaje    []float64
	tiemposAterrizajeDesembarque []float64
}

// Estructura para almacenar los resultados de cada prueba
type resultadoTest struct {
	nombre                        string
	promedioConexionAterrizaje    float64
	promedioAterrizajeDesembarque float64
}
