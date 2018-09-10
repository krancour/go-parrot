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

const recordTime = 10 * time.Second

func main() {
	// Connect
	bebop := client.New()
	log.Println("Connecting...")
	if err := bebop.Connect(); err != nil {
		log.Fatalf("Connection error: %s\n", err)
	}
	log.Println("Connected. [Press ctrl+c to disconnect.]")

	// Start recording
	if err := bebop.StartRecording(); err != nil {
		log.Fatalf("Error starting recording: %s\n", err)
	}
	log.Printf("Recording for %s...\n", recordTime)

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

	// Record for specified time, or until context canceled
	select {
	case <-time.After(recordTime):
	case <-ctx.Done():
		log.Printf("Terminating recording...")
	}

	// Stop recording
	if err := bebop.StopRecording(); err != nil {
		log.Fatalf("Error stopping recording: %s\n", err)
	}
	log.Println("Recording complete.")
}
