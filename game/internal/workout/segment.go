package workout

import "time"

type Watts int

type WorkoutSegment struct {
	Duration   time.Duration
	StartPower Watts
	EndPower   Watts
}

func NewSegment(d time.Duration, w Watts) WorkoutSegment {
	return WorkoutSegment{
		Duration:   d,
		StartPower: w,
		EndPower:   w,
	}
}

func NewBuildupSegment(d time.Duration, beginW Watts, endW Watts) WorkoutSegment {
	return WorkoutSegment{
		Duration:   d,
		StartPower: beginW,
		EndPower:   endW,
	}
}
