package workout

import (
	"math"
	"time"
)

type Workout = []WorkoutSegment

func New() Workout {
	return []WorkoutSegment{}
}

func NewRandom() Workout {
	return []WorkoutSegment{
		NewSegment(5*time.Minute, 120),
		NewSegment(30*time.Second, 360),
		NewSegment(30*time.Second, 120),
		NewSegment(40*time.Minute, 195),
	}
}

func MinPower(t Workout) Watts {
	m := Watts(math.MaxInt)
	for _, s := range t {
		if s.StartPower < m {
			m = s.StartPower
		}

		if s.EndPower < m {
			m = s.EndPower
		}
	}
	return m
}

func MaxPower(t Workout) Watts {
	m := Watts(-1)
	for _, s := range t {
		if s.StartPower > m {
			m = s.StartPower
		}

		if s.EndPower > m {
			m = s.EndPower
		}
	}
	return m
}

func Duration(t Workout) time.Duration {
	d := 0 * time.Second
	for _, s := range t {
		d += s.Duration
	}

	return d
}

func TrainingPowerAt(training Workout, t time.Duration) int {
	progr := t
	for _, tr := range training {
		if progr < tr.Duration {
			return int(tr.StartPower)
		}
		progr -= tr.Duration
	}
	return -1
}
