package sprites

import (
	"fmt"
	"image/color"
	"overlay/game/state"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

type TotalTimer struct {
	font font.Face
	text string
}

func NewTotalTimer(total time.Duration) (*TotalTimer, error) {
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

	return &TotalTimer{
		font: font,
		text: formatTotalDuration(total),
	}, nil
}

func (t *TotalTimer) Draw(screen *ebiten.Image) {
	text.Draw(screen, t.text, t.font, 230, 100, color.White)
}

func (t *TotalTimer) Update(state state.GameState) {
	// Do nothing
}

func formatTotalDuration(total time.Duration) string {
	totalSeconds := int(total / time.Second)
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60
	return fmt.Sprintf("/%02d:%02d:%02d", hours, minutes, seconds)
}
