package bluetooth

import (
	"log/slog"
	"math/rand"
	"strconv"
	"time"
)

type mockPowerChar struct {
	listeners []chan int
}
type speedPowerChar struct {
	listeners []chan int
}

// ContinuousRead extracts instantaneous power (signed 16-bit integer, little-endian)
func (p *mockPowerChar) ContinuousRead() error {
	go func() {
		for {
			for _, listener := range p.listeners {
				listener <- rand.Intn(200) + 100

			}
		}
	}()

	return nil
}

func (p *mockPowerChar) Write(power int) (int, error) {
	slog.Info("Should write " + strconv.Itoa(power) + " to trainer")
	return 0, nil
}

func (p *speedPowerChar) AddListener(c chan int) {
	if p.listeners == nil {
		p.listeners = make([]chan int, 0)
	}

	p.listeners = append(p.listeners, c)
}

func (p *mockPowerChar) AddListener(c chan int) {
	if p.listeners == nil {
		p.listeners = make([]chan int, 0)
	}

	p.listeners = append(p.listeners, c)
}

func (p *speedPowerChar) ContinuousRead() error {
	go func() {
		for {
			for _, listener := range p.listeners {
				listener <- rand.Intn(4) + 28

			}
			time.Sleep(2 * time.Second)
		}
	}()

	return nil
}

func (p *speedPowerChar) Write(speed int) (int, error) {
	slog.Info("Should write " + strconv.Itoa(speed) + " to trainer")
	return 0, nil
}

func NewMockDevice() Device {
	return NewDevice(
		WithPower(&mockPowerChar{}),
		WithSpeed(&speedPowerChar{}),
	)
}
