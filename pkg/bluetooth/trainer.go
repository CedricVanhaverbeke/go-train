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

func (t Trainer) ReadPower(c chan int) error {
	return t.Power.ContinuousRead(c)
}

func (t Trainer) WritePower(power int) (int, error) {
	return t.Power.Write(power)
}

type trainerOpt func(*Trainer)

func WithPower(pow readwriter) trainerOpt {
	return func(t *Trainer) {
		t.Power = pow
	}
}

func NewTrainer(opts ...trainerOpt) Trainer {
	t := &Trainer{}
	for _, opt := range opts {
		opt(t)
	}

	return *t
}
