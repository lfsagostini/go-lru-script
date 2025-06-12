package main

import (
	"fmt"
	"log"
	"time"

	"go-lru-script/cache"
)

func main() {
	
	c, err := cache.New[string, bool](2, 5*time.Second)
	if err != nil {
		log.Fatalf("Error creating cache: %v", err)
	}

	fmt.Println("Starting LRU Cache Test Script...")

	
	fmt.Println("[1 Configuring Cache]")
	c.Set("user:ana@test.com", true)
	c.Set("user:carlos@test.com", true)
	


	fmt.Println("\n[Testting LRU Expulsion]")
	c.Get("user:ana@test.com")
	fmt.Println("Añadiendo 'user:diana@test.com', debería expulsar a Carlos (el menos usado).")
	c.Set("user:diana@test.com", true)

	if _, found := c.Get("user:carlos@test.com"); !found {
		fmt.Println("-> ÉXITO: El token de Carlos fue expulsado por ser LRU.")
	} else {
		fmt.Println("-> ERROR: El token de Carlos no fue expulsado.")
	}
	fmt.Printf("Elementos en el caché: %d\n", c.Len())



	fmt.Println("\n[Paso 3: Probando TTL]")
	fmt.Println("Esperando 6 segundos para que los tokens restantes expiren...")
	time.Sleep(6 * time.Second)

	if _, found := c.Get("user:ana@test.com"); !found {
		fmt.Println("-> ÉXITO: El token de Ana expiró por TTL.")
	} else {
		fmt.Println("-> ERROR: El token de Ana no expiró.")
	}
	
	if _, found := c.Get("user:diana@test.com"); !found {
		fmt.Println("-> ÉXITO: El token de Diana expiró por TTL.")
	} else {
		fmt.Println("-> ERROR: El token de Diana no expiró.")
	}
	fmt.Printf("Elementos finales en el caché: %d\n", c.Len())
}