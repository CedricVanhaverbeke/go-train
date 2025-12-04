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

type timer struct {
	font font.Face
	text string
}

func NewTimer() (*timer, error) {
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

	return &timer{
		font: font,
		text: "00:00:00",
	}, nil
}

func (t *timer) Update(state state.GameState) {
	seconds := int(state.Progress.Duration() / time.Second)
	if seconds < 0 {
		seconds = 0
	}

	t.text = formatDuration(seconds)
}

func (t *timer) Draw(screen *ebiten.Image) {
	text.Draw(screen, t.text, t.font, 20, 100, color.White)
}

func formatDuration(totalSeconds int) string {
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

