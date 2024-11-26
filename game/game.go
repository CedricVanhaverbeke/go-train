package game

import (
	"log"
	"overlay/game/sprites"
	"overlay/game/state"
	"overlay/internal/training"
	"overlay/pkg/consumer"
	"overlay/pkg/producer"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kbinani/screenshot"
)

type game struct {
	width   int
	height  int
	sprites []sprites.Spriter
	start   time.Time
	timer   time.Time

	cM    consumer.Metrics
	State state.GameState
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.width, g.height
}
func (g *game) Update() error {
	now := time.Now()
	if now.Sub(g.timer) > time.Second {
		g.timer = now
		g.State.Progress.Tick()

		// push the power
		g.cM.Power <- training.TrainingPowerAt(g.State.Training, g.State.Progress.Duration())

		for _, s := range g.sprites {
			s.Update(g.State)
		}
	}
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	for _, s := range g.sprites {
		s.Draw(screen)
	}
}

func getCurrentMonitorSize() (int, int) {
	bounds := screenshot.GetDisplayBounds(0) // 0 is the primary display
	width := bounds.Dx()
	height := bounds.Dy()

	return width, height
}

func newGame(training training.Training, cM consumer.Metrics) *game {
	w, h := getCurrentMonitorSize()
	now := time.Now()

	game := &game{
		width:  w,
		height: h,
		start:  now,
		timer:  now,
		sprites: []sprites.Spriter{
			sprites.NewTrainingGraph(w, h, 500, 200, training),
			sprites.NewTimer(),
			sprites.NewPower(),
		},
		State: state.GameState{
			Progress: state.NewProgress(),
			Training: training,
		},
		cM: cM,
	}

	return game
}

func (g *game) subscribe(metrics producer.Metrics) {
	go func() {
		for p := range metrics.Power {
			g.State.Metrics.Power = p
		}
	}()

	go func() {
		for hr := range metrics.Hr {
			g.State.Metrics.Hr = hr
		}
	}()
}

func Run(training training.Training, metrics producer.Metrics, cM consumer.Metrics) {
	game := newGame(training, cM)

	ebiten.SetWindowDecorated(false)
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowMousePassthrough(true)
	ebiten.SetWindowSize(game.width, game.height)

	op := &ebiten.RunGameOptions{}
	op.ScreenTransparent = true
	op.SkipTaskbar = true

	game.subscribe(metrics)

	if err := ebiten.RunGameWithOptions(game, op); err != nil {
		log.Fatal(err)
	}
}
