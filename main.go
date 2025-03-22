package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"overlay/game"
	"overlay/internal/route"
	"overlay/internal/workout"
	"overlay/pkg/bluetooth"
	"overlay/pkg/gpx"
	"path"
	"strings"
	"time"
)

var mock = flag.Bool(
	"mock",
	false,
	"Sets up a mock trainer instead of connecting to a real trainer",
)
var headless = flag.Bool("headless", false, "Sets up the game in headless mode for testing")

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

func writeGpxToFile(fileTitle string, gpxFile gpx.Gpx) error {
	file, err := os.Create(fileTitle)
	if err != nil {
		return err
	}

	err = gpxFile.Write(file)
	if err != nil {
		return err
	}

	return nil

}

func main() {
	flag.Parse()

	trainer, err := newDevice()
	if err != nil {
		panic(err)
	}

	training := workout.NewRandom()
	helloWorldRoute := route.NewExample()

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
		gpxFile.Build(trainer, &helloWorldRoute)
	}()

	fmt.Println("distance of route (in m)", helloWorldRoute.Distance())

	// use the data to run the game
	// the game needs to run in the main thread according
	// to the ebiten spec
	opts := game.NewOpts(game.WithHeadless(*headless), game.WithTickDuration(time.Second))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			slog.Info("Program interrupted, writing to file...")
			err = writeGpxToFile(fileTitle, gpxFile)
			if err != nil {
				slog.Error(err.Error())
			}

			os.Exit(0)
		}
	}()

	game.Run(training, trainer, helloWorldRoute, opts)

	slog.Info("Game ended")
	err = writeGpxToFile(fileTitle, gpxFile)
	if err != nil {
		slog.Error(err.Error())
	}

	cmd := exec.Command("open", "-a", "GpxSee", gpxFile.Path)
	err = cmd.Run()
	if err != nil {
		slog.Error(err.Error())
	}
}
