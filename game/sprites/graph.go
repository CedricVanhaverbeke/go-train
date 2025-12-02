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
	training workout.Workout
	Width    int
	Height   int
	x        int
	parent   *ebiten.Image
}

func NewGraph(x, width int, height int, t workout.Workout) graph {
	return graph{
		Width:    width,
		Height:   height,
		training: t,
		x:        x,
	}
}

func (m graph) Parent() *ebiten.Image {
	return m.parent
}

func (m graph) Update(state state.GameState) {}

func (m graph) Draw(screen *ebiten.Image) {
	screenHeight := screen.Bounds().Dy()

	t := m.training
	totalDuration := workout.Duration(t)

	x := m.x
	for i, s := range t {
		w := scaleWidth(s, totalDuration, m.Width)
		h := scaleHeight(t, i, screenHeight)
		vector.DrawFilledRect(
			screen,
			float32(x),
			float32(screenHeight-h),
			float32(w),
			float32(h),
			color.RGBA{85, 165, 34, 50},
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

// scaleHeight calculates the height of a training block depending on the screen height
func scaleHeight(s workout.Workout, index int, totalHeight int) int {
	// a training segment can take a maximum of 1/15 * screen height
	maxHeight := totalHeight / 15

	// for now, the start and endpower is the same, so just draw that
	frac := float32(s[index].EndPower) / float32(workout.MaxPower(s))
	h := frac * float32(workout.Watts(maxHeight))
	return int(h)
}
