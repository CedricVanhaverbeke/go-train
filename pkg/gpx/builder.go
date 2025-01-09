package gpx

import (
	"encoding/xml"
	"fmt"
	"log/slog"
	"os"
	"overlay/pkg/bluetooth"
	"strings"
	"sync"
)

type chanAvailability struct {
	c         chan int
	available bool
}

func setupChannels(
	trainer *bluetooth.Device,
) (power chanAvailability, speed chanAvailability, cadence chanAvailability) {
	filePowerChan := make(chan int)
	hasPower := trainer.Power.AddListener(filePowerChan)
	power = chanAvailability{c: filePowerChan, available: hasPower}

	fileCadenceChar := make(chan int)
	hasCadence := trainer.Cadence.AddListener(fileCadenceChar)
	speed = chanAvailability{c: fileCadenceChar, available: hasCadence}

	fileSpeedChar := make(chan int)
	hasSpeed := trainer.Speed.AddListener(fileSpeedChar)
	cadence = chanAvailability{c: fileSpeedChar, available: hasSpeed}

	return power, speed, cadence
}

// Build waits for a trackpoint
// to be ready to be added to
// the gpx struct
func (data *Gpx) Build(trainer *bluetooth.Device) {
	power, speed, cadence := setupChannels(trainer)

	for {
		wg := sync.WaitGroup{}
		var powV int
		var cadV int
		var speedV int

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

		wg.Add(1)
		go func() {
			if !speed.available {
				wg.Done()
			}

			// calculate distance in route
			s := <-speed.c
			speedV = s
			wg.Done()
		}()

		wg.Wait()

		slog.Info(
			fmt.Sprintf(
				"Adding trackpoint with power %d, cadence %d and speed %d",
				powV,
				cadV,
				speedV,
			),
		)

		data.AddTrackpoint(NewTrackpoint(
			"23.3581890",
			"54.9870280",
			WithPower(powV),
			WithCadence(cadV),
		))

		// not ideal, but this is a way to write to a file
		// figure out how we could catch a terminate or interrupt signal
		slog.Info("Should append to file")
		fileTitle := strings.ReplaceAll(data.Trk.Name, " ", "_")
		fileTitle += ".gpx"
		file, err := os.Create(fileTitle)
		if err != nil {
			slog.Error(err.Error())
			return
		}

		gpxBytes, err := xml.Marshal(data)
		if err != nil {
			slog.Error(err.Error())
			return
		}

		n, err := file.Write(gpxBytes)
		if err != nil {
			slog.Error(err.Error())
			return
		}

		if n != len(gpxBytes) {
			slog.Error("Did not write whole xml file")
			return
		}

		slog.Info("Done writing to file")
	}
}
