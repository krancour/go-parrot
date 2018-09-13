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

const flightTime = 3 * time.Second
const withVideo = true

func main() {
	// Connect
	bebop := client.New()
	if err := bebop.Connect(); err != nil {
		log.Fatalf("Connection error: %s\n", err)
	}
	log.Println("Connected. [Press ctrl+c to disconnect.]")

	// Configure drone
	log.Println("Configuring drone...")
	log.Println("Configuring drone for indoor flight...")
	if err := bebop.Outdoor(false); err != nil {
		log.Fatalf("Error configuring for indoor flight: %s\n", err)
	}
	log.Println("Configured for indoor flight.")
	if err := bebop.HullProtection(false); err != nil {
		log.Fatalf("Error configuring physical drone attributes: %s\n", err)
	}
	log.Println("Configured physical drone attributes.")
	if !withVideo {
		log.Println("Disabling video capture...")
		if err := bebop.VideoAutoRecordSelection(false); err != nil {
			log.Fatalf("Error disabling video capture: %s\n", err)
		}
		log.Println("Video capture disabled.")
	}
	log.Println("Dron configured.")

	// Take off
	log.Println("Launching...")
	if err := bebop.TakeOff(); err != nil {
		log.Fatalf("Launch error: %s\n", err)
	}
	log.Println("We have liftoff!")
	log.Printf("Hovering for %s...\n", flightTime)

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

	// Hover for specified time, or until context canceled
	select {
	case <-time.After(flightTime):
	case <-ctx.Done():
		log.Printf("Terminating flight early...")
	}

	// Land
	log.Println("Landing...")
	if err := bebop.Land(); err != nil {
		log.Fatalf("Landing error: %s\n", err)
	}
	log.Println("We have touch down!")

	bebop.StartRecording()
}
