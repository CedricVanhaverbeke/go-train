package main

import (
	"overlay/game"
	"overlay/internal/training"
	"overlay/pkg/bluetooth"
)

func main() {
	trainer, err := bluetooth.Connect()
	if err != nil {
		panic(err)
	}

	training := training.NewRandom()

	game.Run(training, trainer)
}
