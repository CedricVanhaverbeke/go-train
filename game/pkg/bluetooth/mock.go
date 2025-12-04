package bluetooth

import (
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"time"
)

type mockPowerChar struct {
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
			p.listeners.WriteValue(200)
			time.Sleep(2 * time.Second)
		}
	}()

	return nil
}

func (p *mockPowerChar) AddListener(c chan int) bool {
	p.listeners.AddListener(c)
	return true
}

func (p *mockPowerChar) Write(power int) (int, error) {
	slog.Info("Should write " + strconv.Itoa(power) + " to trainer")
	ep := encode(power)
	// contains three bytes
	// add a new byte in between
	newEp := slices.Insert(ep, 1, byte(0))
	return decode(newEp), nil
}

func NewMockDevice() Device {
	return NewDevice(
		WithPower(&mockPowerChar{}),
		WithCadence(&errorCadenceChar{}),
	)
}
