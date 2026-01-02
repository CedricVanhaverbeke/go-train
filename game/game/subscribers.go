package game

import (
	"log/slog"
	"sync/atomic"

	"overlay/pkg/bluetooth"
)

// TODO: this stuff should not be inside the game.
// the game should only get the channels itself,
// not a bluetooth trainer directly
func (g *game) subscribePwr(tr *bluetooth.Device, pauseCounter *atomic.Int32) {
	if tr.Power == nil {
		return
	}

	powerChan := make(chan int)
	tr.Power.AddListener(powerChan)

	go func() {
		for p := range powerChan {
			// Start when a new reading comes in
			if g.State.Progress.Pause && p > 0 {
				g.State.Progress.Pause = false
				pauseCounter.Store(0)
			}

			if p == 0 {
				pauseCounter.Add(1)
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
	g.trainer.Speed.AddListener(speedChan)
	go func() {
		for p := range speedChan {
			// only start when first power comes in
			if g.State.Progress.Pause && p > 0 {
				g.State.Progress.Pause = true
			}
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
	g.trainer.Cadence.AddListener(cadChan)
	go func() {
		for p := range cadChan {
			// only start when first power comes in
			if p > 0 {
				g.State.Progress.Pause = false
			}
			g.State.Metrics.Speed = p
		}
	}()
}
