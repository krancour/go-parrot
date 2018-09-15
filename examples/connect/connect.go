package main

import (
	"fmt"
	"log"

	arnetworkal "github.com/krancour/go-parrot/protocols/arnetworkal/wifi"
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
