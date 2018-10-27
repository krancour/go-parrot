package main

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	arnetworkal "github.com/krancour/go-parrot/protocols/arnetworkal/wifi"
)

func main() {
	conn, err := arnetworkal.NewConnection()
	if err != nil {
		log.Fatal(err)
	}
	<-time.After(5 * time.Second)
	defer func() {
		conn.Close()
		fmt.Println("Done")
	}()
}
