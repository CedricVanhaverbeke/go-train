package bluetooth

import (
	"log/slog"
	"math/rand"
	"strconv"
	"time"
)

type mockPowerChar struct {
	listeners
}
type speedPowerChar struct {
	listeners
}

// ContinuousRead extracts instantaneous power (signed 16-bit integer, little-endian)
func (p *mockPowerChar) ContinuousRead() error {
	go func() {
		for {
			p.listeners.WriteValue(rand.Intn(200) + 100)
			time.Sleep(2 * time.Second)
		}
	}()

	return nil
}

func (p *mockPowerChar) Write(power int) (int, error) {
	slog.Info("Should write " + strconv.Itoa(power) + " to trainer")
	return 0, nil
}

func (p *speedPowerChar) ContinuousRead() error {
	go func() {
		for {
			p.listeners.WriteValue(rand.Intn(4) + 28)
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
