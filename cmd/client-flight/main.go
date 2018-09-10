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
const withVideo = false

func main() {
	// Connect
	bebop := client.New()
	if err := bebop.Connect(); err != nil {
		log.Fatalf("Connection error: %s\n", err)
	}
	log.Println("Connected. [Press ctrl+c to disconnect.]")

	// Configure drone
	if err := bebop.Outdoor(false); err != nil {
		log.Fatalf("Error configuring for indoor flight: %s\n", err)
	}
	log.Println("Configured for indoor flight.")
	if err := bebop.HullProtection(false); err != nil {
		log.Fatalf("Error configuring physical drone attributes: %s\n", err)
	}
	log.Println("Configured physical drone attributes.")

	// Take off
	if err := bebop.TakeOff(); err != nil {
		log.Fatalf("Launch error: %s\n", err)
	}
	log.Println("We have liftoff!")
	log.Printf("Hovering for %s...\n", flightTime)

	// Video is captured by default. Turn it off if not desired.
	if !withVideo {
		if err := bebop.StopRecording(); err != nil {
			log.Fatalf("Error disabling video capure: %s", err)
		}
		log.Println("Video capture disabled.")
	}

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
	if err := bebop.Land(); err != nil {
		log.Fatalf("Landing error: %s\n", err)
	}
	log.Println("We have touch down!")

	bebop.StartRecording()
}
