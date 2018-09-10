package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gobot.io/x/gobot/platforms/parrot/bebop/client"
)

func main() {
	// Connect
	bebop := client.New()
	log.Println("Connecting...")
	if err := bebop.Connect(); err != nil {
		log.Fatalf("Connection error: %s\n", err)
	}
	log.Println("Connected. [Press ctrl+c to disconnect.]")

	// Video setup
	log.Println("Setting up video stream...")
	if err := bebop.VideoStreamMode(0); err != nil {
		log.Fatalf("Error setting setting lowest latency streaming mode: %s\n", err)
	}
	log.Println("Using lowest latency streaming mode.")
	if err := bebop.VideoEnable(true); err != nil {
		log.Fatalf("Error enabling video stream: %s\n", err)
	}
	log.Println("Live video stream enabled.")

	// Signal handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		fmt.Println()
		log.Printf("%s received; shutting down...\n", sig)
		<-time.After(3 * time.Second)
		cancel()
	}()

	// Wait for the context to be canceled
	<-ctx.Done()
}
