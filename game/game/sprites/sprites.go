package sprites

import (
	"overlay/game/state"

	"github.com/hajimehoshi/ebiten/v2"
)

type Spriter interface {
	Draw(screen *ebiten.Image)
	Update(state state.GameState)
}

type CombinedSpriter interface {
	Spriter
	Parent() *ebiten.Image
}

// returns the coordinates for a centered
// rectangle based on it's width. It always
// returns it at the bottom
func CoordCenterRectStart(rectangleWidth int, screenWidth int) (int, int) {
	x := (screenWidth / 2) - (rectangleWidth / 2)
	y := 0

	return x, y
}
