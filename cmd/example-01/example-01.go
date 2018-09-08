package main

import (
	"log"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/parrot/bebop"
)

func main() {
	bebopAdaptor := bebop.NewAdaptor()
	drone := bebop.NewDriver(bebopAdaptor)

	work := func() {
		if err := drone.On(bebop.Flying, func(data interface{}) {
			gobot.After(10*time.Second, func() {
				drone.Land()
			})
		}); err != nil {
			log.Fatal(err)
		}

		if err := drone.HullProtection(true); err != nil {
			log.Fatal(err)
		}

		drone.TakeOff()
	}

	robot := gobot.NewRobot("drone",
		[]gobot.Connection{bebopAdaptor},
		[]gobot.Device{drone},
		work,
	)

	if err := robot.Start(); err != nil {
		log.Fatal(err)
	}
}
