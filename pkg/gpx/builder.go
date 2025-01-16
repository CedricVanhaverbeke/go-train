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
) (power chanAvailability, cadence chanAvailability) {
	filePowerChan := make(chan int)
	hasPower := trainer.Power.AddListener(filePowerChan)
	power = chanAvailability{c: filePowerChan, available: hasPower}

	fileCadenceChar := make(chan int)
	hasCadence := trainer.Cadence.AddListener(fileCadenceChar)
	cadence = chanAvailability{c: fileCadenceChar, available: hasCadence}

	// the values are implicitely returned but I don't like that
	// it's just for readability when using this function
	return power, cadence
}

// Build waits for a trackpoint
// to be ready to be added to
// the gpx struct
func (data *Gpx) Build(trainer *bluetooth.Device, route *Gpx) {
	power, cadence := setupChannels(trainer)
	distance := 0.0

	for {
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

		vrel := route.Speed(distance, powV)
		vrelms := vrel / 3.6

		// every two seconds this updates
		// TODO: chek this out better
		distance += vrelms * float64(2)

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
