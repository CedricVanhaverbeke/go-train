package bluetooth

type readwriter interface {
	Write(int) (int, error)
	ContinuousRead(chan int) error
}

type Trainer struct {
	characteristics []readwriter
}

func (t Trainer) ReadPower(c chan int) error {
	return t.characteristics[0].ContinuousRead(c)
}

func (t Trainer) WritePower(power int) (int, error) {
	return t.characteristics[0].Write(power)
}

func NewTrainer(chars ...readwriter) Trainer {
	return Trainer{
		characteristics: chars,
	}
}
