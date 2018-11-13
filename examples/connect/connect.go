package main

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	arnetworkal "github.com/krancour/go-parrot/protocols/arnetworkal/wifi"
)

func main() {
	frameSender, frameReceiver, err := arnetworkal.Connect()
	if err != nil {
		log.Fatal(err)
	}
	<-time.After(5 * time.Second)
	defer func() {
		frameSender.Close()
		frameReceiver.Close()
		fmt.Println("Done")
	}()
}
