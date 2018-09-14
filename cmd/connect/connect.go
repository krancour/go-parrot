package main

import (
	"fmt"
	"log"

	arnetworkal "github.com/krancour/drone-examples/pkg/parrot/protocols/arnetworkal/wifi"
)

func main() {
	conn, err := arnetworkal.NewConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		conn.Close()
		fmt.Println("Done")
	}()
}
