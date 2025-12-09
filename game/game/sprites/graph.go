package sprites

import (
	"image/color"
	"overlay/game/state"
	"overlay/internal/workout"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type graph struct {
	training  workout.Workout
	Width     int
	Height    int
	x         int
	parent    *ebiten.Image
	gameState state.GameState
}

func NewGraph(x, width int, height int, t workout.Workout) *graph {
	return &graph{
		Width:    width,
		Height:   height,
		training: t,
		x:        x,
	}
}

func (m *graph) Parent() *ebiten.Image {
	return m.parent
}

func (m *graph) Update(state state.GameState) {
	m.gameState = state
}

func (m *graph) Draw(screen *ebiten.Image) {
	screenHeight := screen.Bounds().Dy()

	t := m.training
	totalDuration := workout.Duration(t)

	_, currentSegmentIndex := workout.TrainingSegmentAt(t, m.gameState.Progress.Duration())

	x := m.x
	for i, s := range t {
		c := color.RGBA{85, 165, 34, 50}
		w := scaleWidth(s, totalDuration, m.Width)
		if s.StartPower != s.EndPower {
			rico := float64(s.EndPower-s.StartPower) / float64(w)
			for j := range w {
				p := rico*float64(j) + float64(s.StartPower)
				h := scaleHeight(t, p, screenHeight)
				vector.DrawFilledRect(
					screen,
					float32(x+j),
					float32(screenHeight-h),
					1,
					float32(h),
					c,
					true,
				)
			}
			x += w
			continue
		}

		h := scaleHeightAtIndex(t, i, screenHeight)
		if i == currentSegmentIndex {
			c = color.RGBA{255, 165, 34, 50}
		}
		vector.DrawFilledRect(
			screen,
			float32(x),
			float32(screenHeight-h),
			float32(w),
			float32(h),
			c,
			true,
		)
		x += w
	}
}

// scaleWidth calculates the width of the training block based on the screen size
func scaleWidth(s workout.WorkoutSegment, totalDuration time.Duration, totalWidth int) int {
	totalMinutes := totalDuration.Minutes()
	frac := s.Duration.Minutes() / totalMinutes

	return int(frac * float64(totalWidth))
}

// scaleHeightAtIndex calculates the height of a training block depending on the screen height
func scaleHeightAtIndex(s workout.Workout, index int, screenHeight int) int {
	p := float64(s[index].EndPower)
	return scaleHeight(s, p, screenHeight)
}

func scaleHeight(s workout.Workout, p float64, screenHeight int) int {
	maxHeight := screenHeight / 15
	maxPower := float64(workout.MaxPower(s))

	frac := p / maxPower
	h := frac * float64(maxHeight)
	return int(h)
}
