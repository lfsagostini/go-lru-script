package main

import (
	"fmt"
	"go-lru-script/cache"
	"log"
	"time"
)
var tokenCache *cache.Cache[string, bool]

func handleRequest(shopperId string) {
	fmt.Printf("\n--- Petición recibida para Shopper ID: %s ---\n", shopperId)

	_, tokenExists := tokenCache.Get(shopperId)

	if !tokenExists {
		fmt.Printf("Token para [%s] no encontrado o expirado. ¡Procesando y añadiendo al caché!\n", shopperId)
		tokenCache.Set(shopperId, true)
		fmt.Printf("Token para [%s] añadido al caché.\n", shopperId)
	} else {
		fmt.Printf("Token para [%s] ya existe en el caché. Saltando proceso duplicado.\n", shopperId)
	}
}

func main() {
	var err error
	// Inicializar el caché con un tamaño máximo de 1000 entradas y un TTL de 5 segundos
	tokenCache, err = cache.New[string, bool](1000, 5*time.Second)
	if err != nil {
		log.Fatalf("No se pudo crear el caché: %v", err)
	}

	handleRequest("ana-123")
	handleRequest("ana-123")

	fmt.Println("\n...Esperando 6 segundos para que el TTL invalide el caché...")
	time.Sleep(6 * time.Second)

	handleRequest("ana-123")
	
	
	if _, found := tokenCache.Get("ana-123"); found {
		fmt.Println("\nConfirmado: El token de Ana está de nuevo en el caché después de la re-validación.")
	}
}