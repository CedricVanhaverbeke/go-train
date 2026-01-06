package main

import (
	"flag"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"strings"
	"time"

	"overlay/game"
	"overlay/internal/workout"
	"overlay/pkg/bluetooth"
	"overlay/pkg/gpx"
	"overlay/pkg/repo"
)

var mock = flag.Bool(
	"mock",
	false,
	"Sets up a mock trainer instead of connecting to a real trainer",
)
var headless = flag.Bool("headless", false, "Sets up the game in headless mode for testing")

var selectedWorkout = flag.String("workout", "", "workout to start")

func newDevice() (*bluetooth.Device, error) {
	if *mock {
		return newMockDevice()
	}

	return bluetooth.Connect()
}

func newMockDevice() (*bluetooth.Device, error) {
	trainer := bluetooth.NewMockDevice()
	return &trainer, nil
}

func newTraining(gpxRepo *repo.GPXRepo) {
	flag.Parse()

	trainer, err := newDevice()
	if err != nil {
		panic(err)
	}

	training := workout.NewRandom()
	if selectedWorkout != nil && *selectedWorkout != "" {
		slog.Info(*selectedWorkout)
		training, err = workout.FromString(*selectedWorkout)
		if err != nil {
			panic(err)
		}
	}

	title := "Hello World Ride"
	gpxFile := gpx.New(title)

	// listen for data of the trainer
	trainer.Listen()

	fileTitle := strings.ReplaceAll(title, " ", "_")
	fileTitle += ".gpx"

	dir, _ := os.Getwd()
	gpxFile.Path = path.Join(dir, fileTitle)

	// use the data to build a gpx file
	go func() {
		gpxFile.Build(trainer)
	}()

	// use the data to run the game
	// the game needs to run in the main thread according
	// to the ebiten spec
	opts := game.NewOpts(game.WithHeadless(*headless), game.WithTickDuration(time.Second))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			slog.Info("Program interrupted, writing to file...")
			_, err = gpxRepo.Create("test", gpxFile)
			if err != nil {
				slog.Error(err.Error())
			}

			os.Exit(0)
		}
	}()

	game.Run(training, trainer, opts)

	slog.Info("Game ended")
	_, err = gpxRepo.Create("test", gpxFile)
	if err != nil {
		slog.Error(err.Error())
	}

	cmd := exec.Command("open", "-a", "GpxSee", gpxFile.Path)
	err = cmd.Run()
	if err != nil {
		slog.Error(err.Error())
	}
}

func main() {
	p, _ := os.Getwd()
	repo, err := repo.NewGPXRepo(path.Join(p, "../", "db"))
	if err != nil {
		panic(err)
	}
	newTraining(repo)
}
