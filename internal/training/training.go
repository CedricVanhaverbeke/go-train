package training

import (
	"math"
	"time"
)

type Training = []TrainingSegment

func New() Training {
	return []TrainingSegment{}
}

func NewRandom() Training {
	return []TrainingSegment{
		NewSegment(10*time.Second, 150),
		NewSegment(2*time.Second, 200),
		NewSegment(2*time.Second, 180),
		NewSegment(2*time.Second, 200),
		NewSegment(2*time.Second, 180),
		NewSegment(2*time.Second, 200),
		NewSegment(2*time.Second, 180),
		NewSegment(2*time.Second, 180),
		NewSegment(2*time.Second, 180),
		NewSegment(2*time.Second, 180),
		NewSegment(20*time.Second, 120),
	}
}

func MinPower(t Training) Watts {
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

func MaxPower(t Training) Watts {
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

func Duration(t Training) time.Duration {
	d := 0 * time.Second
	for _, s := range t {
		d += s.Duration
	}

	return d
}

func TrainingPowerAt(training Training, t time.Duration) int {
	progr := t
	for _, tr := range training {
		if progr < tr.Duration {
			return int(tr.StartPower)
		}
		progr -= tr.Duration
	}
	return -1
}
