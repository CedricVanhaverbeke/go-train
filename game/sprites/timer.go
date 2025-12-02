package sprites

import (
	"bytes"
	"image/color"
	"overlay/game/state"
	"strconv"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

type timer struct {
	face    *text.GoTextFace
	text    string
	size    int
	font    *sfnt.Font
	fontBuf sfnt.Buffer
}

func NewTimer() *timer {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.PressStart2P_ttf))
	if err != nil {
		panic(err)
	}

	fnt, err := sfnt.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		panic(err)
	}

	return &timer{
		face: &text.GoTextFace{
			Source: s,
			Size:   32,
		},
		text: "0",
		size: 32,
		font: fnt,
	}

}

func (t *timer) Update(state state.GameState) {
	seconds := int(state.Progress.Duration() / time.Second)
	if seconds < 0 {
		seconds = 0
	}

	t.text = formatDuration(seconds)
}

func (t *timer) Draw(screen *ebiten.Image) {
	op := &text.DrawOptions{}
	width := t.measureText()
	op.GeoM.Translate(float64(t.size)+5+width/2, float64(t.size)+10)
	op.ColorScale.ScaleWithColor(color.White)
	op.LineSpacing = float64(t.size)
	op.PrimaryAlign = text.AlignCenter
	text.Draw(screen, t.text, t.face, op)
}

func (t *timer) measureText() float64 {
	if t.font == nil || t.text == "" {
		return 0
	}

	ppem := fixed.I(t.size)
	var width fixed.Int26_6
	var prev sfnt.GlyphIndex
	hasPrev := false

	for _, r := range t.text {
		glyph, err := t.font.GlyphIndex(&t.fontBuf, r)
		if err != nil {
			continue
		}

		if hasPrev {
			if kern, err := t.font.Kern(&t.fontBuf, prev, glyph, ppem, font.HintingNone); err == nil {
				width += kern
			}
		}

		adv, err := t.font.GlyphAdvance(&t.fontBuf, glyph, ppem, font.HintingNone)
		if err != nil {
			continue
		}

		width += adv
		prev = glyph
		hasPrev = true
	}

	return float64(width) / 64
}

func formatDuration(totalSeconds int) string {
	if totalSeconds < 60 {
		return strconv.Itoa(totalSeconds)
	}

	minutes := totalSeconds / 60
	seconds := totalSeconds % 60
	if minutes < 60 {
		return buildMinuteSecond(minutes, seconds)
	}

	hours := minutes / 60
	minutes = minutes % 60

	return buildHourMinuteSecond(hours, minutes, seconds)
}

func buildMinuteSecond(minutes, seconds int) string {
	var b strings.Builder
	// at least m:ss, e.g. 4 characters for "0:59"
	b.Grow(5)
	b.WriteString(strconv.Itoa(minutes))
	b.WriteByte(':')
	if seconds < 10 {
		b.WriteByte('0')
	}
	b.WriteString(strconv.Itoa(seconds))
	return b.String()
}

func buildHourMinuteSecond(hours, minutes, seconds int) string {
	var b strings.Builder
	// at least h:mm:ss
	b.Grow(8)
	b.WriteString(strconv.Itoa(hours))
	b.WriteByte(':')
	appendTwoDigits(&b, minutes)
	b.WriteByte(':')
	appendTwoDigits(&b, seconds)
	return b.String()
}

func appendTwoDigits(b *strings.Builder, value int) {
	if value < 10 {
		b.WriteByte('0')
	}
	b.WriteString(strconv.Itoa(value))
}
