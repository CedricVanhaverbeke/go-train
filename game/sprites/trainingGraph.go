package sprites

import (
	"overlay/game/state"
	"overlay/internal/training"

	"github.com/hajimehoshi/ebiten/v2"
)

type TrainingGraph struct {
	training     training.Training
	Width        int
	Height       int
	graphSprites []Spriter
	x            int
}

func NewTrainingGraph(
	screenWidth int,
	screenHeight int,
	width int,
	height int,
	t training.Training,
) TrainingGraph {
	startX, _ := CoordCenterRectStart(width, screenWidth)

	return TrainingGraph{
		Width:    width,
		Height:   height,
		training: t,
		x:        startX,
		graphSprites: []Spriter{
			NewGraph(startX, width, height, t),
			NewProgressLine(startX, 0, width),
		},
	}
}

func (m TrainingGraph) Update(state state.GameState) {
	for _, s := range m.graphSprites {
		s.Update(state)
	}
}

func (m TrainingGraph) Draw(screen *ebiten.Image) {
	for _, s := range m.graphSprites {
		s.Draw(screen)
	}
}
