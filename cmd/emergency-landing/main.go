package main

import (
	"log"

	"gobot.io/x/gobot/platforms/parrot/bebop/client"
)

func main() {
	// Connect
	bebop := client.New()
	log.Println("Connecting...")
	if err := bebop.Connect(); err != nil {
		log.Fatalf("Connection error: %s\n", err)
	}
	log.Println("Connected.")

	// Land
	log.Println("Executing emergency landing...")
	if err := bebop.Land(); err != nil {
		log.Fatalf("Landing error: %s\n", err)
	}
	log.Println("We have touch down!")
}
