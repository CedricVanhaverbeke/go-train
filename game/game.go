package game

import (
	"log"
	"overlay/game/sprites"
	"overlay/game/state"
	"overlay/internal/training"
	"overlay/pkg/bluetooth"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kbinani/screenshot"
	"golang.org/x/exp/slog"
)

type game struct {
	width   int
	height  int
	sprites []sprites.Spriter
	start   time.Time
	timer   time.Time

	trainer *bluetooth.Trainer
	State   state.GameState
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
		_, err := g.trainer.WritePower(
			training.TrainingPowerAt(g.State.Training, g.State.Progress.Duration()),
		)

		if err != nil {
			slog.Error("could not write power: ", err)
		}

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

func newGame(training training.Training, trainer *bluetooth.Trainer) *game {
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
		trainer: trainer,
	}

	return game
}

func (g *game) subscribe(tr *bluetooth.Trainer) {
	powerChan := make(chan int)
	err := tr.ReadPower(powerChan)
	if err != nil {
		slog.Error("Could not read power")
	}

	go func() {
		for p := range powerChan {
			g.State.Metrics.Power = p
		}
	}()
}

func Run(training training.Training, trainer *bluetooth.Trainer) {
	game := newGame(training, trainer)

	ebiten.SetWindowDecorated(false)
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowMousePassthrough(true)
	ebiten.SetWindowSize(game.width, game.height)

	op := &ebiten.RunGameOptions{}
	op.ScreenTransparent = true
	op.SkipTaskbar = true

	game.subscribe(trainer)

	if err := ebiten.RunGameWithOptions(game, op); err != nil {
		log.Fatal(err)
	}
}
