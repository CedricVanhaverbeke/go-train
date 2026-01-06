package workout

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"time"
)

type Workout struct {
	Name     string
	FTP      int
	Segments []WorkoutSegment
}

func New() Workout {
	return Workout{
		Segments: []WorkoutSegment{},
		Name:     "",
		FTP:      0,
	}
}

func FromString(workoutString string) (*Workout, error) {
	workout := New()
	w := []WorkoutSegment{}

	info := strings.Split(workoutString, ";")
	// first two steps are reserved for workout name
	// and ftp value

	name := info[0]
	ftpStr := info[1]
	ftp, err := strconv.Atoi(ftpStr)
	if err != nil {
		return nil, err
	}

	workout.FTP = ftp
	workout.Name = name

	workoutSteps := info[2:]
	for _, s := range workoutSteps {
		powerDuration := strings.Split(s, "-")
		startPower := powerDuration[0]
		endPower := powerDuration[1]
		duration := powerDuration[2]

		pStartInt, err := strconv.Atoi(startPower)
		if err != nil {
			return nil, errors.New("could not parse workout")
		}
		pEndInt, err := strconv.Atoi(endPower)
		if err != nil {
			return nil, errors.New("could not parse workout")
		}

		durationInt, err := strconv.Atoi(duration)
		if err != nil {
			return nil, errors.New("could not parse workout")
		}
		w = append(w, NewSegment(time.Second*time.Duration(durationInt), Watts(pStartInt), Watts(pEndInt)))
	}

	workout.Segments = w

	return &workout, nil
}

func NewRandom() *Workout {
	return &Workout{
		Segments: []WorkoutSegment{
			{Duration: 30 * time.Minute, StartPower: 120, EndPower: 200},
			NewSegment(30*time.Second, 120, 120),
			NewSegment(40*time.Minute, 195, 195),
		},
		Name: "random",
		FTP:  200,
	}
}

func MinPower(t Workout) Watts {
	m := Watts(math.MaxInt)
	for _, s := range t.Segments {
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
	for _, s := range t.Segments {
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
	for _, s := range t.Segments {
		d += s.Duration
	}

	return d
}

func TrainingPowerAt(training Workout, t time.Duration) int {
	progr := t
	for _, tr := range training.Segments {
		if progr < tr.Duration {
			rico := (float64(tr.EndPower) - float64(tr.StartPower)) / float64(tr.Duration)
			p := int(rico*float64(progr)) + int(tr.StartPower)
			return p
		}
		progr -= tr.Duration
	}
	return -1
}

func TrainingSegmentAt(training Workout, t time.Duration) (*WorkoutSegment, int) {
	progr := t
	for i, tr := range training.Segments {
		if progr < tr.Duration {
			return &tr, i
		}
		progr -= tr.Duration
	}
	return nil, -1
}
