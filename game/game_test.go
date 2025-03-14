package game_test

import (
	"log/slog"
	"os"
	"overlay/game"
	"overlay/internal/route"
	"overlay/pkg/bluetooth"
	"overlay/pkg/gpx"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestGame(t *testing.T) {
	// for loop until the game exits

	trainer := bluetooth.NewMockDevice()
	tr := Workout.NewRandom()
	helloWorldRoute := route.NewExample()
	title := "Hello World Ride"
	gpxFile := gpx.New(title)

	// listen for data of the trainer
	trainer.Listen()

	fileTitle := strings.ReplaceAll(title, " ", "_")
	fileTitle += ".gpx"

	file, err := os.Create(fileTitle)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	dir, _ := os.Getwd()
	gpxFile.Path = path.Join(dir, fileTitle)

	// use the data to build a gpx file
	go func() {
		gpxFile.Build(&trainer, &helloWorldRoute, file)
	}()

	tickDuration := 1 * time.Millisecond

	opts := game.NewOpts(game.WithHeadless(true), game.WithTickDuration(tickDuration))
	g := game.NewGame(tr, &trainer, opts)

	// wait for game to have started

	seconds := 0
	g.State.Progress.Started = true

	for {
		seconds++
		err := g.Update()
		if err != nil {
			if err == ebiten.Termination {
				break
			}
			t.Error(err.Error())
		}

		time.Sleep(time.Millisecond)
	}

	if float64(seconds) != g.State.Progress.Duration().Seconds() {
		t.Error("Should have the same duration")
	}
}
