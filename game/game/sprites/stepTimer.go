package sprites

import (
	"fmt"
	"image/color"
	"time"

	"overlay/game/state"
	"overlay/internal/workout"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

type StepTimer struct {
	font font.Face
	text string
}

func NewStepTimer() (*StepTimer, error) {
	tt, err := opentype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}

	font, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    48,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, err
	}

	return &StepTimer{
		font: font,
		text: "00:00",
	}, nil
}

func (t *StepTimer) Update(s state.GameState) {
	segment, index := workout.TrainingSegmentAt(s.Training, s.Progress.Duration())
	if segment == nil {
		t.text = "00:00"
		return
	}

	var timeInPreviousSegments time.Duration
	for i := 0; i < index; i++ {
		timeInPreviousSegments += s.Training.Segments[i].Duration
	}

	timeInCurrentSegment := s.Progress.Duration() - timeInPreviousSegments
	remainingTime := segment.Duration - timeInCurrentSegment
	t.text = formatStepDuration(remainingTime)
}

func (t *StepTimer) Draw(screen *ebiten.Image) {
	text.Draw(screen, t.text, t.font, 20, 150, color.White)
}

func formatStepDuration(d time.Duration) string {
	totalSeconds := int(d / time.Second)
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}
