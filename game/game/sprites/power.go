package sprites

import (
	"image/color"
	"strconv"

	"overlay/game/state"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

type power struct {
	font font.Face
	text string
}

func NewPower() (*power, error) {
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

	return &power{
		font: font,
		text: "0",
	}, nil
}

func (p *power) Update(state state.GameState) {
	p.text = strconv.Itoa(state.Metrics.Power)
}

func (p *power) Draw(screen *ebiten.Image) {
	dx := len(p.text) * 20
	text.Draw(screen, p.text, p.font, screen.Bounds().Dx()-50-dx, 100, color.White)
}
