package main

import (
	"net/http"
	_ "net/http/pprof"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/products/bebop2"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	log.SetLevel(log.InfoLevel)
	_, err := bebop2.NewController()
	if err != nil {
		log.Fatal(err)
	}
	select {}
}
