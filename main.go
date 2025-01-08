package main

import (
	"flag"
	"fmt"
	"log/slog"
	"overlay/game"
	"overlay/internal/route"
	"overlay/internal/training"
	"overlay/pkg/bluetooth"
	"overlay/pkg/gpx"
	"sync"
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

	filePowerChan := make(chan int)
	hasPower := trainer.Power.AddListener(filePowerChan)

	fileCadenceChar := make(chan int)
	hasCadence := trainer.Cadence.AddListener(fileCadenceChar)

	fileSpeedChar := make(chan int)
	hasSpeed := trainer.Speed.AddListener(fileSpeedChar)

	gpxFile := gpx.New("Hello World Ride")
	go func() {
		for {
			wg := sync.WaitGroup{}
			var pow int
			var cad int
			var speed int

			wg.Add(1)
			go func() {
				if !hasPower {
					wg.Done()
				}

				p := <-filePowerChan
				pow = p
				wg.Done()
			}()

			wg.Add(1)
			go func() {
				if !hasCadence {
					wg.Done()
				}

				c := <-fileCadenceChar
				cad = c
				wg.Done()
			}()

			wg.Add(1)
			go func() {
				if !hasSpeed {
					wg.Done()
				}

				// calculate distance in route
				s := <-fileSpeedChar
				speed = s
				wg.Done()
			}()

			wg.Wait()

			slog.Info(
				fmt.Sprintf(
					"Adding trackpoint with power %d, cadence %d and speed %d",
					pow,
					cad,
					speed,
				),
			)

			gpxFile.AddTrackpoint(gpx.NewTrackpoint(
				"23.3581890",
				"54.9870280",
				gpx.WithPower(pow),
				gpx.WithCadence(cad),
			))
		}
	}()

	trainer.Listen()
	game.Run(training, trainer, helloWorldRoute)
}
