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

type timer struct {
	source *text.GoTextFaceSource
	text   int
	size   int
}

func NewTimer() *timer {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.PressStart2P_ttf))
	if err != nil {
		panic(err)
	}

	return &timer{
		source: s,
		text:   0,
		size:   32,
	}

}

func (t *timer) Update(state state.GameState) {
	t.text += 1
}

func (t *timer) Draw(screen *ebiten.Image) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(t.size)+5, float64(t.size)+10)
	op.ColorScale.ScaleWithColor(color.White)
	op.LineSpacing = float64(t.size)
	op.PrimaryAlign = text.AlignCenter
	text.Draw(screen, strconv.Itoa(t.text), &text.GoTextFace{
		Source: t.source,
		Size:   float64(t.size),
	}, op)
}
