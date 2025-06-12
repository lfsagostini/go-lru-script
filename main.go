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
    tokenCache  *cache.Cache[string, bool]
    wg          sync.WaitGroup
    cacheHits   atomic.Int64
    cacheMisses atomic.Int64
    setsDone    atomic.Int64
)

// Prints memory, CPU, and goroutine stats every second
func monitorSystem() {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    var m runtime.MemStats
    for range ticker.C {
        runtime.ReadMemStats(&m)
        fmt.Printf("🖥️ Memory: %.2f MB | CPUs: %d | Goroutines: %d\n",
            float64(m.Alloc)/1024/1024,
            runtime.NumCPU(),
            runtime.NumGoroutine(),
        )
    }
}

// Simulates a user making cache requests
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
        // Random pause between 10 and 100 ms
        time.Sleep(time.Duration(10+rand.Intn(90)) * time.Millisecond)
    }
}

// Prints the number of users in cache every second
func monitorCache() {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    for range ticker.C {
        fmt.Printf("⏳ Users in cache: %d\n", tokenCache.Len())
    }
}

func main() {
    var err error
    tokenCache, err = cache.New[string, bool](cacheCapacity, cacheTTL)
    if err != nil {
        log.Fatalf("Could not create cache: %v", err)
    }
    go monitorCache()
    go monitorSystem()

    fmt.Println("🚀 Starting concurrent stress test...")
    fmt.Printf("   - Concurrent users: %d\n", numUsers)
    fmt.Printf("   - Cache capacity: %d\n", cacheCapacity)
    fmt.Printf("   - TTL: %v\n", cacheTTL)
    fmt.Println("-------------------------------------------------")

    startTime := time.Now()

    for i := 0; i < numUsers; i++ {
        wg.Add(1)
        go userWorker(i)
    }

    fmt.Println("... All workers started, waiting for completion ...")
    wg.Wait()

    duration := time.Since(startTime)

    fmt.Println("\n✅ Test completed!")
    fmt.Println("-------------------------------------------------")
    fmt.Printf("Total duration: %v\n", duration)
    fmt.Println("\nCache stats:")

    totalHits := cacheHits.Load()
    totalMisses := cacheMisses.Load()
    totalSets := setsDone.Load()
    totalOps := totalHits + totalMisses

    hitRate := float64(totalHits) / float64(totalOps) * 100

    fmt.Printf("  - Hits:         %d\n", totalHits)
    fmt.Printf("  - Misses:       %d\n", totalMisses)
    fmt.Printf("  - Set ops:      %d\n", totalSets)
    fmt.Printf("  - Hit rate:     %.2f%%\n", hitRate)
    fmt.Printf("  - Final items:  %d\n", tokenCache.Len())
    fmt.Println("-------------------------------------------------")

    if totalMisses != totalSets {
        fmt.Println("⚠️ Warning: Misses and sets count do not match. Possible issue.")
    }
}