package bluetooth

type readwriter interface {
	Write(int) (int, error)
	ContinuousRead(chan int) error
}

type Trainer struct {
	Power   readwriter
	Speed   readwriter
	Cadence readwriter
}

type trainerOpt func(*Trainer)

func WithPower(pow readwriter) trainerOpt {
	return func(t *Trainer) {
		t.Power = pow
	}
}

func WithSpeed(v readwriter) trainerOpt {
	return func(t *Trainer) {
		t.Speed = v
	}
}

func WithCadence(cad readwriter) trainerOpt {
	return func(t *Trainer) {
		t.Cadence = cad
	}
}

func NewTrainer(opts ...trainerOpt) Trainer {
	t := &Trainer{}
	for _, opt := range opts {
		opt(t)
	}

	return *t
}
