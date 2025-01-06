package bluetooth

type readwriter interface {
	Write(int) (int, error)
	ContinuousRead(chan int) error
}

type Device struct {
	Power   readwriter
	Speed   readwriter
	Cadence readwriter
}

type trainerOpt func(*Device)

func WithPower(pow readwriter) trainerOpt {
	return func(t *Device) {
		t.Power = pow
	}
}

func WithSpeed(v readwriter) trainerOpt {
	return func(t *Device) {
		t.Speed = v
	}
}

func WithCadence(cad readwriter) trainerOpt {
	return func(t *Device) {
		t.Cadence = cad
	}
}

func NewDevice(opts ...trainerOpt) Device {
	t := &Device{}
	for _, opt := range opts {
		opt(t)
	}

	return *t
}
