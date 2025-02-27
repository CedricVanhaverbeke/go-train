package game_test

import (
	"overlay/game"
	"overlay/internal/route"
	"overlay/internal/training"
	"overlay/pkg/bluetooth"
	"overlay/pkg/gpx"
	"sync"
	"testing"
	"time"
)

func TestGame(t *testing.T) {
	// for loop until the game exits

	trainer := bluetooth.NewMockDevice()
	training := training.NewRandom()
	helloWorldRoute := route.New()
	title := "Hello World Ride"
	gpxFile := gpx.New(title)

	// listen for data of the trainer
	trainer.Listen()

	// use the data to build a gpx file
	go func() {
		gpxFile.Build(&trainer, &helloWorldRoute)
	}()

	tickDuration := time.Second
	ticker := time.NewTicker(tickDuration)
	opts := game.NewOpts(game.WithHeadless(true), game.WithTickDuration(tickDuration))
	g := game.NewGame(training, &trainer, opts)

	quit := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			select {
			case <-ticker.C:
				err := g.Update()
				t.Error(err)
				wg.Done()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	wg.Wait()

}
