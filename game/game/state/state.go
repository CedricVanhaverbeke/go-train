package state

import "overlay/internal/workout"

type GameState struct {
	Metrics  Metrics
	Progress Progress
	Training workout.Workout
}
