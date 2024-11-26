package training

import "time"

type Watts int

type TrainingSegment struct {
	Duration   time.Duration
	StartPower Watts
	EndPower   Watts
}

func NewSegment(d time.Duration, w Watts) TrainingSegment {
	return TrainingSegment{
		Duration:   d,
		StartPower: w,
		EndPower:   w,
	}
}

func NewBuildupSegment(d time.Duration, beginW Watts, endW Watts) TrainingSegment {
	return TrainingSegment{
		Duration:   d,
		StartPower: beginW,
		EndPower:   endW,
	}
}
