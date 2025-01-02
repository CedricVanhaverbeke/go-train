package bluetooth

import (
	"log/slog"
	"math/rand"
	"strconv"
	"time"
)

type mockPowerChar struct{}

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

func NewRandTrainer() Trainer {
	return NewTrainer(WithPower(mockPowerChar{}))
}
