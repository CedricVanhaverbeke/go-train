package bluetooth

import (
	"fmt"
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

// mock characteristic
// that is not implemented

type errorCadenceChar struct {
	listeners
}

func (char *errorCadenceChar) AddListener(c chan int) bool {
	return false
}

func (p *errorCadenceChar) ContinuousRead() error {
	return fmt.Errorf("Could not read cadence")
}

func (p *errorCadenceChar) Write(power int) (int, error) {
	return 0, fmt.Errorf("Could not write")
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
		WithCadence(&errorCadenceChar{}),
	)
}
