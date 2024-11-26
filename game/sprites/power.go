package sprites

import (
	"bytes"
	"image/color"
	"overlay/game/state"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type power struct {
	source *text.GoTextFaceSource
	text   string
	size   int
}

func NewPower() *power {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.PressStart2P_ttf))
	if err != nil {
		panic(err)
	}

	return &power{
		source: s,
		size:   32,
	}
}

func (t *power) Update(state state.GameState) {
	t.text = strconv.Itoa(state.Metrics.Power)
}

func (p *power) Draw(screen *ebiten.Image) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(screen.Bounds().Dx()-50), float64(p.size)+10)
	op.ColorScale.ScaleWithColor(color.White)
	op.LineSpacing = float64(p.size)
	op.PrimaryAlign = text.AlignCenter
	text.Draw(screen, p.text, &text.GoTextFace{
		Source: p.source,
		Size:   float64(p.size),
	}, op)
}
