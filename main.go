package main

import (
	"overlay/game"
	"overlay/internal/training"
	"overlay/pkg/consumer"
	"overlay/pkg/producer"
	"overlay/pkg/trainer"
)

func main() {
	// 1. Create a training
	t := training.NewRandom()

	// game consumes data from trainers
	powerChan := make(chan int)
	heartRateChan := make(chan int)

	pM := producer.NewMetrics(powerChan, heartRateChan)
	trainer := trainer.NewRandomTrainer(pM)

	// listen for a channel here
	// the channel sends messages about the metrics of the game
	// do not block
	go func() {
		// the connected trainer produces metrics
		trainer.Produce()
	}()

	consumerPowerChan := make(chan int)
	cM := consumer.NewMetrics(consumerPowerChan)
	go func() {
		// the connected trainer must also consume metrics of the game depending on the state
		//	eg. it must consume the amount of watts it should push
		trainer.Consume(cM)
	}()

	// 3. Play the training.
	game.Run(t, pM, cM)
}
