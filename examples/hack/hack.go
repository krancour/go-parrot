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
	controller, err := bebop2.NewController()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for range controller.Common().CommonState().BatteryPercentCh() {
			battery, _ := controller.Common().CommonState().BatteryPercent()
			log.Infof("current battery level: %d%%", battery)
		}
	}()
	select {}
}
