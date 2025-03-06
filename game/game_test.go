package game_test

import (
	"fmt"
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
	tr := training.NewRandom()
	helloWorldRoute := route.NewExample()
	title := "Hello World Ride"
	gpxFile := gpx.New(title)

	// listen for data of the trainer
	trainer.Listen()

	// use the data to build a gpx file
	go func() {
		gpxFile.Build(&trainer, &helloWorldRoute)
	}()

	tickDuration := time.Millisecond
	ticker := time.NewTicker(tickDuration)
	duration := training.Duration(tr)
	timer := time.NewTimer(time.Duration(duration.Seconds() * float64(time.Millisecond)))

	opts := game.NewOpts(game.WithHeadless(true), game.WithTickDuration(tickDuration))
	g := game.NewGame(tr, &trainer, opts)

	// wait for game to have started

	seconds := 0
	err := g.Update()
	if err != nil {
		t.Error(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			select {
			case <-ticker.C:
				seconds++
				err := g.Update()
				if err != nil {
					t.Error(err)
				}
			case <-timer.C:
				fmt.Println("timer fired")
				wg.Done()
			}
		}
	}()

	wg.Wait()
	fmt.Println(g.State.Progress.Duration().Seconds())
	fmt.Println(seconds)
	if float64(seconds) != g.State.Progress.Duration().Seconds() {
		t.Error("Should have the same duration")
	}
}
