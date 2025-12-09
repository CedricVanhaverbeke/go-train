package workout

import "time"

type Watts int

type WorkoutSegment struct {
	Duration   time.Duration
	StartPower Watts
	EndPower   Watts
}

func NewSegment(d time.Duration, startW Watts, endW Watts) WorkoutSegment {
	return WorkoutSegment{
		Duration:   d,
		StartPower: startW,
		EndPower:   endW,
	}
}

func NewBuildupSegment(d time.Duration, beginW Watts, endW Watts) WorkoutSegment {
	return WorkoutSegment{
		Duration:   d,
		StartPower: beginW,
		EndPower:   endW,
	}
}
