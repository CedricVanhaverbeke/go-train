package sprites

import (
	"image/color"
	"overlay/game/state"
	"overlay/internal/workout"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type progressLine struct {
	x float64
	y int

	width int
}

func NewProgressLine(x int, y int, width int) *progressLine {
	return &progressLine{
		x:     float64(x),
		y:     y,
		width: width,
	}
}

func (p *progressLine) Update(state state.GameState) {
	step := float64(p.width) / workout.Duration(state.Training).Seconds()
	p.x += step
}

func (p *progressLine) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(
		screen,
		float32(p.x),
		float32(screen.Bounds().Dy()-100),
		float32(1),
		float32(100),
		color.RGBA{255, 0, 0, 50},
		true,
	)
}
