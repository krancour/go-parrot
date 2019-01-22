package main

import (
	"net/http"
	_ "net/http/pprof"
	"time"

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
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			<-ticker.C
			battery, ok := controller.Common().CommonState().BatteryPercent()
			if ok {
				log.Infof("current battery level: %d%%", battery)
			}
		}
	}()
	select {}
}
