package bluetooth

import "fmt"

type readwriter interface {
	Write([]byte) (int, error)
	ContinuousRead(chan int) error
}

type Trainer struct {
	characteristics []readwriter
}

func (t Trainer) ReadPower(c chan int) error {
	return t.characteristics[0].ContinuousRead(c)
}

func (t Trainer) WritePower(power int) (int, error) {
	fmt.Println("Should write ", power)
	return t.characteristics[0].Write([]byte("ello"))
}

func NewTrainer(chars ...readwriter) Trainer {
	return Trainer{
		characteristics: chars,
	}
}
