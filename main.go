package main

import (
	"flag"
	"overlay/game"
	"overlay/internal/route"
	"overlay/internal/training"
	"overlay/pkg/bluetooth"
	"overlay/pkg/gpx"
)

var mock = flag.Bool("m", false, "Sets up a mock trainer instead of connecting to a real trainer")

func newDevice() (*bluetooth.Device, error) {
	if *mock {
		trainer := bluetooth.NewMockDevice()
		return &trainer, nil
	}

	return bluetooth.Connect()
}

func main() {
	flag.Parse()

	trainer, err := newDevice()
	if err != nil {
		panic(err)
	}

	training := training.NewRandom()
	helloWorldRoute := route.New()

	fileChan := make(chan int)
	trainer.Power.AddListener(fileChan)

	gpxFile := gpx.New("Hello World Ride")
	go func() {
		for p := range fileChan {
			tp := gpx.NewTrackpoint(
				"23.3581890",
				"54.9870280",
				gpx.WithPower(p))
			gpxFile.AddTrackpoint(tp)
		}
	}()

	trainer.Listen()
	game.Run(training, trainer, helloWorldRoute)
}
