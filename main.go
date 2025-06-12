package main

import (
	"fmt"
	"go-lru-script/cache"
	"log"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)


const (
    numUsers      = 3000
    opsPerUser    = 100
    cacheCapacity = 3000
    cacheTTL      = 5 * time.Second 
)


var (
	tokenCache *cache.Cache[string, bool]
	wg         sync.WaitGroup 
	cacheHits   atomic.Int64
	cacheMisses atomic.Int64
	setsDone    atomic.Int64
)

func monitorSystem() {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    var m runtime.MemStats
    for range ticker.C {
        runtime.ReadMemStats(&m)
        fmt.Printf("🖥️ Memoria: %.2f MB | CPUs: %d | Goroutines: %d\n",
            float64(m.Alloc)/1024/1024,
            runtime.NumCPU(),
            runtime.NumGoroutine(),
        )
    }
}

func userWorker(id int) {
    defer wg.Done()
    shopperID := fmt.Sprintf("shopper-%d", id)

    for i := 0; i < opsPerUser; i++ {
        _, found := tokenCache.Get(shopperID)
        if found {
            cacheHits.Add(1)
        } else {
            cacheMisses.Add(1)
            tokenCache.Set(shopperID, true)
            setsDone.Add(1)
        }
        // Pausa aleatoria entre 10 y 100 ms
        time.Sleep(time.Duration(10+rand.Intn(90)) * time.Millisecond)
    }
}

func monitorCache() {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    for range ticker.C {
        fmt.Printf("⏳ Usuarios en caché: %d\n", tokenCache.Len())
    }
}

func main() {
	
	var err error
	tokenCache, err = cache.New[string, bool](cacheCapacity, cacheTTL)
	if err != nil {
		log.Fatalf("No se pudo crear el caché: %v", err)
	}
	go monitorCache()
	go monitorSystem()

	fmt.Println("🚀 Iniciando prueba de estrés concurrente...")
	fmt.Printf("   - Usuarios concurrentes: %d\n", numUsers)
	fmt.Printf("   - Capacidad del caché: %d\n", cacheCapacity)
	fmt.Printf("   - TTL: %v\n", cacheTTL)
	fmt.Println("-------------------------------------------------")

	startTime := time.Now()


	for i := 0; i < numUsers; i++ {
		wg.Add(1) 
		go userWorker(i) 
	}


	fmt.Println("... Todos los workers iniciados, esperando que terminen ...")
	wg.Wait()

	duration := time.Since(startTime)

	fmt.Println("\n✅ ¡Prueba completada!")
	fmt.Println("-------------------------------------------------")
	fmt.Printf("Duración total: %v\n", duration)
	fmt.Println("\nEstadísticas del Caché:")

	totalHits := cacheHits.Load()
	totalMisses := cacheMisses.Load()
	totalSets := setsDone.Load()
	totalOps := totalHits + totalMisses

	hitRate := float64(totalHits) / float64(totalOps) * 100

	fmt.Printf("  - Aciertos (Hits):    %d\n", totalHits)
	fmt.Printf("  - Fallos (Misses):    %d\n", totalMisses)
	fmt.Printf("  - Operaciones Set:    %d\n", totalSets)
	fmt.Printf("  - Tasa de Aciertos:   %.2f%%\n", hitRate)
	fmt.Printf("  - Elementos finales:  %d\n", tokenCache.Len())
	fmt.Println("-------------------------------------------------")

	if totalMisses != totalSets {
		fmt.Println("⚠️ Advertencia: El número de misses no coincide con el de sets. Puede indicar un problema.")
	}
}