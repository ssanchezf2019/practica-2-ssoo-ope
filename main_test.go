// EJECUCIÓN: go test -v

package main

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

var resultadosTests []resultadoTest

// Función de ayuda para ejecutar la simulación y calcular los promedios
func ejecutarSimulacion() (float64, float64) {
	torre := &TorreControl{
		pistas:                       make(chan struct{}, NumPistas),
		puertas:                      make(chan struct{}, NumPuertas),
		esperaAterrizaje:             make(chan struct{}, MaxEspera),
		wg:                           &sync.WaitGroup{},
		mutex:                        sync.Mutex{},
		inicio:                       time.Now(),
		tiemposConexiónAterrizaje:    []float64{},
		tiemposAterrizajeDesembarque: []float64{},
	}

	for i := 0; i < NumAviones; i++ {
		torre.wg.Add(1)
		avion := Avion{id: i}
		go torre.aterrizar(avion)
	}

	// Esperar a que todos los aviones terminen
	torre.wg.Wait()

	// Calcular los promedios
	promedioConexionAterrizaje := calcularPromedio(torre.tiemposConexiónAterrizaje)
	promedioAterrizajeDesembarque := calcularPromedio(torre.tiemposAterrizajeDesembarque)

	return promedioConexionAterrizaje, promedioAterrizajeDesembarque
}

// TestMain para ejecutar todos los tests y luego mostrar la tabla comparativa
func TestMain(m *testing.M) {
	// Ejecuta todas las pruebas
	m.Run()

	// Generar tabla comparativa de los resultados
	fmt.Printf("\nTabla Comparativa de Tiempos Promedio\n")
	fmt.Printf("%-50s | %-40s | %-40s\n", "Prueba", "Promedio Conexión a Aterrizaje (s)", "Promedio Aterrizaje a Desembarque (s)")
	fmt.Println(strings.Repeat("-", 110))
	for _, res := range resultadosTests {
		fmt.Printf("%-50s | %-40.2f | %-40.2f\n", res.nombre, res.promedioConexionAterrizaje, res.promedioAterrizajeDesembarque)
	}
}

// Función para almacenar resultados de cada prueba en la tabla comparativa
func registrarResultados(nombre string, promedioConexionAterrizaje, promedioAterrizajeDesembarque float64) {
	resultadosTests = append(resultadosTests, resultadoTest{
		nombre:                        nombre,
		promedioConexionAterrizaje:    promedioConexionAterrizaje,
		promedioAterrizajeDesembarque: promedioAterrizajeDesembarque,
	})
}

// Prueba 0: Ejecución estándar sin cambios en la configuración
func TestEjecucionEstandar(t *testing.T) {
	promedioConexionAterrizaje, promedioAterrizajeDesembarque := ejecutarSimulacion()

	registrarResultados("Ejecución Estándar", promedioConexionAterrizaje, promedioAterrizajeDesembarque)
}

// Prueba 1: Duplicar la cantidad máxima de aviones esperando
func TestDuplicarMaxEspera(t *testing.T) {
	MaxEspera = MaxEspera * 2
	defer func() { MaxEspera = 5 }() // Restaura el valor original después de la prueba

	promedioConexionAterrizaje, promedioAterrizajeDesembarque := ejecutarSimulacion()

	registrarResultados("Duplicar MaxEspera", promedioConexionAterrizaje, promedioAterrizajeDesembarque)
}

// Prueba 2: Variación en el tiempo de uso/espera al 25% sobre el nominal
func TestVariacionTiempo(t *testing.T) {
	VariacionTiempo = VariacionTiempo + 25
	defer func() { VariacionTiempo = 15 }() // Restaura el valor original después de la prueba

	promedioConexionAterrizaje, promedioAterrizajeDesembarque := ejecutarSimulacion()

	registrarResultados("Variación Tiempo +25%", promedioConexionAterrizaje, promedioAterrizajeDesembarque)
}

// Prueba 3: Duplicar la cantidad máxima de aviones esperando y variar el tiempo al 25% sobre el nominal
func TestDuplicarMaxEsperaYVariacionTiempo(t *testing.T) {
	MaxEspera = MaxEspera * 2
	VariacionTiempo = VariacionTiempo + 25
	defer func() {
		MaxEspera = 5 // Restaura el valor original después de la prueba
		VariacionTiempo = 15
	}()

	promedioConexionAterrizaje, promedioAterrizajeDesembarque := ejecutarSimulacion()

	registrarResultados("Duplicar MaxEspera y Variación Tiempo +25%", promedioConexionAterrizaje, promedioAterrizajeDesembarque)
}

// Prueba 4: Multiplicar el número de pistas por 5
func TestMultiplicarPistas(t *testing.T) {
	NumPistas = NumPistas * 5
	defer func() { NumPistas = 3 }() // Restaura el valor original después de la prueba

	promedioConexionAterrizaje, promedioAterrizajeDesembarque := ejecutarSimulacion()

	registrarResultados("Multiplicar Pistas x5", promedioConexionAterrizaje, promedioAterrizajeDesembarque)
}

// Prueba 5: Multiplicar el número de pistas por 5 y aumentar el tiempo de uso de cada pista 5 veces
func TestMultiplicarPistasYTiempo(t *testing.T) {
	NumPistas = NumPistas * 5
	TiempoAterrizaje = TiempoAterrizaje * 5
	defer func() {
		NumPistas = 3 // Restaura el valor original después de la prueba
		TiempoAterrizaje = 1000
	}()

	promedioConexionAterrizaje, promedioAterrizajeDesembarque := ejecutarSimulacion()

	registrarResultados("Multiplicar Pistas x5 y Tiempo Aterrizaje x5", promedioConexionAterrizaje, promedioAterrizajeDesembarque)
}

// Prueba 6: Multiplicar el tiempo de uso de cada puerta por 3
func TestMultiplicarTiempoPuerta(t *testing.T) {
	TiempoPuerta = TiempoPuerta * 3
	defer func() { TiempoPuerta = 1500 }() // Restaura el valor original después de la prueba

	promedioConexionAterrizaje, promedioAterrizajeDesembarque := ejecutarSimulacion()

	registrarResultados("Multiplicar Tiempo Puerta x3", promedioConexionAterrizaje, promedioAterrizajeDesembarque)
}
