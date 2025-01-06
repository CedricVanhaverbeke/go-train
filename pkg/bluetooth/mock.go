package bluetooth

import (
	"log/slog"
	"math/rand"
	"strconv"
	"time"
)

type mockPowerChar struct{}
type speedPowerChar struct{}

// ContinuousRead extracts instantaneous power (signed 16-bit integer, little-endian)
func (p mockPowerChar) ContinuousRead(c chan int) error {
	go func() {
		for {

			c <- rand.Intn(200) + 100
			time.Sleep(2 * time.Second)
		}
	}()

	return nil
}

func (p mockPowerChar) Write(power int) (int, error) {
	slog.Info("Should write " + strconv.Itoa(power) + " to trainer")
	return 0, nil
}

func (p speedPowerChar) ContinuousRead(c chan int) error {
	go func() {
		for {
			c <- rand.Intn(4000) + 28000
			time.Sleep(2 * time.Second)
		}
	}()

	return nil
}

func (p speedPowerChar) Write(speed int) (int, error) {
	slog.Info("Should write " + strconv.Itoa(speed) + " to trainer")
	return 0, nil
}

func NewMockDevice() Device {
	return NewDevice(
		WithPower(mockPowerChar{}),
		WithSpeed(speedPowerChar{}),
	)
}
