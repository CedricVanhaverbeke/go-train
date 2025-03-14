package gpx

import (
	"encoding/xml"
	"fmt"
	"io"
	"log/slog"
	"overlay/pkg/bluetooth"
	"sync"
	"time"
)

type chanAvailability struct {
	c         chan int
	available bool
}

func setupChannels(
	trainer *bluetooth.Device,
) (power chanAvailability, cadence chanAvailability) {
	filePowerChan := make(chan int)
	hasPower := trainer.Power.AddListener(filePowerChan)
	power = chanAvailability{c: filePowerChan, available: hasPower}

	// fileCadenceChar := make(chan int)
	// hasCadence := trainer.Cadence.AddListener(fileCadenceChar)
	// cadence = chanAvailability{c: fileCadenceChar, available: hasCadence}

	// the values are implicitely returned but I don't like that
	// it's just for readability when using this function
	return power, cadence
}

func valuesWait(power chanAvailability, cadence chanAvailability) (int, int) {
	wg := sync.WaitGroup{}
	var powV int
	var cadV int

	wg.Add(1)
	go func() {
		if !power.available {
			wg.Done()
		}

		p := <-power.c
		powV = p
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		if !cadence.available {
			wg.Done()
		}

		c := <-cadence.c
		cadV = c
		wg.Done()
	}()

	wg.Wait()
	return powV, cadV
}

func (data *Gpx) Write(out io.Writer) error {
	gpxBytes, err := xml.Marshal(data)
	if err != nil {
		return err
	}

	n, err := out.Write(gpxBytes)
	if err != nil {
		return err
	}

	if n != len(gpxBytes) {
		return err
	}

	return nil
}

// Build waits for a trackpoint
// to be ready to be added to
// the gpx struct
func (data *Gpx) Build(trainer *bluetooth.Device, route *Gpx) {
	power, cadence := setupChannels(trainer)
	distance := 0.0

	for {
		before := time.Now()

		powV, cadV := valuesWait(power, cadence)

		after := time.Now()
		timeD := (after.Sub(before)).Seconds()

		vrel := route.Speed(distance, powV)

		vrelms := vrel / 3.6

		distance += vrelms * float64(timeD)
		lat, lng, ele, _, _ := route.CoordInfo(distance)

		slog.Info(
			fmt.Sprintf(
				"Adding trackpoint with power %d, cadence %d and speed %f, distance %f lat %f, lng %f",
				powV,
				cadV,
				vrel,
				distance,
				lat,
				lng,
			),
		)

		tp := NewTrackpoint(
			lat,
			lng,
			WithPower(powV),
			WithCadence(cadV),
			WithElevation(ele),
		)

		data.AddTrackpoint(tp)
	}
}
