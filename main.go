package main

import (
	"flag"
	"fmt"
	"overlay/game"
	"overlay/internal/route"
	"overlay/internal/training"
	"overlay/pkg/bluetooth"
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

	go func() {
		for p := range fileChan {
			fmt.Println(p)
		}
	}()

	trainer.Listen()
	game.Run(training, trainer, helloWorldRoute)
}
