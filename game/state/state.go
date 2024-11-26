package state

import (
	"overlay/internal/training"
)

type GameState struct {
	Metrics  Metrics
	Progress Progress
	Training training.Training
}
