package sprites

import (
	"image/color"
	"overlay/game/state"
	"overlay/internal/training"

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
	step := float64(p.width) / training.Duration(state.Training).Seconds()
	p.x += step
}

func (p *progressLine) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(
		screen,
		float32(p.x),
		float32(screen.Bounds().Dy()-500),
		float32(1),
		float32(500),
		color.RGBA{85, 165, 34, 50},
		true,
	)
}
