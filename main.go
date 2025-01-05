package main

import (
	"flag"
	"overlay/game"
	"overlay/internal/route"
	"overlay/internal/training"
	"overlay/pkg/bluetooth"
)

var mock = flag.Bool("m", false, "Sets up a mock trainer instead of connecting to a real trainer")

func newTrainer() (*bluetooth.Trainer, error) {
	if *mock {
		trainer := bluetooth.NewRandTrainer()
		return &trainer, nil
	}

	return bluetooth.Connect()
}

func main() {
	flag.Parse()

	trainer, err := newTrainer()
	if err != nil {
		panic(err)
	}

	training := training.NewRandom()
	helloWorldRoute := route.New()
	game.Run(training, trainer, helloWorldRoute)
}
