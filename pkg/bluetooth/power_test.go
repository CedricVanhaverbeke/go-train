package bluetooth_test

import (
	"overlay/pkg/bluetooth"
	"testing"
)

func TestPowerWrite(t *testing.T) {
	mock := bluetooth.NewMockDevice()

	expected := 300
	actual, err := mock.Power.Write(expected)

	if err != nil {
		t.Error("Failed to write")
	}

	if actual != expected {
		t.Errorf("Should be 300, got %d", actual)
	}
}

func TestPowerRead(t *testing.T) {
	mock := bluetooth.NewMockDevice()
	c := make(chan int)

	mock.Power.AddListener(c)
	_ = mock.Power.ContinuousRead()

	p := <-c
	if p != 200 {
		t.Error("Mock device always reads 200")
	}
}
