package game

import (
	"log"
	"overlay/game/sprites"
	"overlay/game/state"
	"overlay/internal/workout"
	"overlay/pkg/bluetooth"
	"overlay/pkg/gpx"
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

	trainer *bluetooth.Device
	State   state.GameState
	opts    Opts
}

type Opts struct {
	// Determines if the game runs in headless mode
	Headless bool

	// Determines how fast the game moves, the default is one second.
	// This tickDuration exists mainly for testing purposes
	TickDuration time.Duration
}

func WithHeadless(headless bool) func(opts *Opts) {
	return func(opts *Opts) {
		opts.Headless = headless
	}
}

func WithTickDuration(tickDuration time.Duration) func(opts *Opts) {
	return func(opts *Opts) {
		opts.TickDuration = tickDuration
	}
}

func NewOpts(optsArgs ...func(opts *Opts)) Opts {
	opts := Opts{
		Headless:     false,
		TickDuration: time.Second,
	}

	for _, arg := range optsArgs {
		arg(&opts)
	}

	return opts
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.width, g.height
}
func (g *game) Update() error {
	if !g.State.Progress.Started {
		return nil
	}

	now := time.Now()
	prevTimer := g.timer
	prevPower := workout.TrainingPowerAt(g.State.Training, g.State.Progress.Duration())

	if now.Sub(prevTimer) > g.opts.TickDuration {
		g.timer = now
		g.State.Progress.Tick()

		newPower := workout.TrainingPowerAt(g.State.Training, g.State.Progress.Duration())

		// check if we changed the power or the game just started
		if prevPower != newPower || g.State.Progress.Duration() == 1*time.Second {
			_, err := g.trainer.Power.Write(
				newPower,
			)

			if err != nil {
				slog.Error("could not write power: ", err)
			}
		}

		for _, s := range g.sprites {
			s.Update(g.State)
		}
	}

	if g.State.Progress.Duration() >= workout.Duration(g.State.Training) {
		// return a termination error
		return ebiten.Termination
	}

	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	if g.opts.Headless {
		return
	}

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

func NewGame(training workout.Workout, trainer *bluetooth.Device, opts Opts) *game {
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
		opts:    opts,
	}

	return game
}

func (g *game) subscribe(tr *bluetooth.Device) {
	g.subscribePwr(tr)
	g.subscribeCadence(tr)
}

func Run(training workout.Workout, trainer *bluetooth.Device, route gpx.Gpx, opts Opts) {
	game := NewGame(training, trainer, opts)

	ebiten.SetWindowDecorated(false)
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowMousePassthrough(true)
	ebiten.SetWindowSize(game.width, game.height)

	op := &ebiten.RunGameOptions{}
	op.ScreenTransparent = true
	op.SkipTaskbar = true
	// subscribe to all trainer metrics
	game.subscribe(trainer)

	if err := ebiten.RunGameWithOptions(game, op); err != nil {
		log.Fatal(err)
	}
}
