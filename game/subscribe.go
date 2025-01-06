package game

import (
	"fmt"
	"log/slog"
	"overlay/pkg/bluetooth"
)

// TODO: this stuff should not be inside the game.
// the game should only get the channels itself,
// not a bluetooth trainer directly
func (g *game) subscribePwr(tr *bluetooth.Device) {
	if tr.Power == nil {
		return
	}

	// this code can live inside another package
	// let's say the metrics package
	// the metrics package exposes the channels from where
	// a game can be read
	powerChan := make(chan int)
	err := tr.Power.ContinuousRead(powerChan)
	if err != nil {
		slog.Error("Could not read power")
	}

	go func() {
		for p := range powerChan {
			// only start when first power comes in
			if p > 0 {
				g.State.Progress.Started = true
			}
			g.State.Metrics.Power = p
		}
	}()
}

func (g *game) subscribeSpeed(tr *bluetooth.Device) {
	if tr.Speed == nil {
		slog.Info("Speed characteristic is nil on device")
		return
	}

	speedChan := make(chan int)
	err := tr.Speed.ContinuousRead(speedChan)
	if err != nil {
		slog.Error("Could not read speed")
	}

	go func() {
		for p := range speedChan {
			// only start when first power comes in
			if p > 0 {
				g.State.Progress.Started = true
			}
			slog.Info(fmt.Sprintf("Got speed: %d", p))
			g.State.Metrics.Speed = p
		}
	}()
}

func (g *game) subscribeCadence(tr *bluetooth.Device) {
	if tr.Cadence == nil {
		slog.Info("Cadence characteristic is nil on device")
		return
	}

	cadChan := make(chan int)
	err := tr.Power.ContinuousRead(cadChan)
	if err != nil {
		slog.Error("Could not read cadence")
	}

	go func() {
		for p := range cadChan {
			// only start when first power comes in
			if p > 0 {
				g.State.Progress.Started = true
			}
			g.State.Metrics.Speed = p
		}
	}()
}
