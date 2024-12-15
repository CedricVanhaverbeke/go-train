package bluetooth

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter
var zwiftHubUUID = "c96fb5f7-b4d5-262e-7baf-0a479225f3ab"

var ftmsUUID = "00001826-0000-1000-8000-00805f9b34fb"
var fitnessMachineControlPointUUID = "00002ad9-0000-1000-8000-00805f9b34fb"
var cyclingPower = "00001818-0000-1000-8000-00805f9b34fb"

// notice that the only thing that's different from the ftmsUUID is the first segment.
// this is the case wit hall uuids
// file:///Users/cedricvanhaverbeke/Downloads/GATT_Specification_Supplement_v5.pdf
// https://gist.github.com/sam016/4abe921b5a9ee27f67b3686910293026
// var indoorBikeData = "00002ad2-0000-1000-8000-00805f9b34fb"
var cyclingPowerMeasureMent = "00002a63-0000-1000-8000-00805f9b34fb"

func Scan() error {
	err := adapter.Enable()
	if err != nil {
		panic(err)
	}

	slog.Info("Finding trainer...")
	char, err := DiscoverFTMSDevice()
	if err != nil {
		return err
	}

	err = char.EnableNotifications(func(buf []byte) {
		// Extract instantaneous power (signed 16-bit integer, little-endian)
		power := int16(buf[2]) | int16(buf[3])<<8
		log.Printf("Instantaneous Power: %d watts", power)
	})

	if err != nil {
		return err
	}

	return nil
}

// DiscoverFTMSDevice checks every available device
// having the FTMS service. It returns the
// first device that has FTMS enabled and
// instantaneous power
func DiscoverFTMSDevice() (*bluetooth.DeviceCharacteristic, error) {
	found := make(chan bluetooth.ScanResult)
	devChar := make(chan *bluetooth.DeviceCharacteristic)
	ftmsCP := make(chan *bluetooth.DeviceCharacteristic)
	scan := make(chan bool)

	scanned := map[string]bool{}

	ftmsService, err := bluetooth.ParseUUID(ftmsUUID)
	if err != nil {
		return nil, err
	}

	ftmsControlPoint, err := bluetooth.ParseUUID(fitnessMachineControlPointUUID)
	if err != nil {
		return nil, err
	}

	cyclingPowerService, err := bluetooth.ParseUUID(cyclingPower)
	if err != nil {
		return nil, err
	}

	serviceUuid, err := bluetooth.ParseUUID(cyclingPowerMeasureMent)
	if err != nil {
		return nil, err
	}

	continueScanning := func(s string) {
		slog.Info(s)
		scan <- true
	}

	go func() {
		for range scan {
			slog.Info("Scanning bluetooth devices...")
			err := adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
				f := make(chan bluetooth.ScanResult)
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				go func() {
					if !scanned[device.Address.String()] {
						scanned[device.Address.String()] = true

						slog.Info("Found device with uuid: " + device.Address.String())
						f <- device
						_ = adapter.StopScan()
						return
					}
				}()

				select {
				case device := <-f:
					found <- device
				case <-ctx.Done():
					slog.Info("bluetooth scanning timeout exceeded")
					devChar <- nil
					return
				}
			})

			if err != nil {
				panic(err)
			}
		}
	}()

	go func() {
		for scanResult := range found {
			slog.Info("Checking device...")

			device, err := adapter.Connect(
				scanResult.Address,
				bluetooth.ConnectionParams{
					ConnectionTimeout: bluetooth.NewDuration(1 * time.Second),
				},
			)
			if err != nil {
				continueScanning("Cannot connect to device")
				continue
			}

			dservices, err := device.DiscoverServices([]bluetooth.UUID{cyclingPowerService})
			if err != nil {
				continueScanning("Device does not have cycling power enabled")
				continue
			}

			hasCyclingPower := len(dservices) == 1
			if !hasCyclingPower {
				continueScanning("Device does not have cycling power enabled")
				continue
			}

			service := dservices[0]

			chars, err := service.DiscoverCharacteristics([]bluetooth.UUID{serviceUuid})
			if err != nil {
				continueScanning("Could not get characteristics " + err.Error())
				continue
			}

			charsOk := len(chars) == 1
			if !charsOk {
				continueScanning("Device does not have instantaneous power characteristic")
				continue
			}

			ftmsServices, err := device.DiscoverServices([]bluetooth.UUID{ftmsService})
			if err != nil {
				continueScanning("Could not find ftms service")
			}

			hasFtms := len(ftmsServices) == 1
			if !hasFtms {
				continueScanning("Could not find ftms service")
			}

			ftms := ftmsServices[0]
			fmtsControlPointChar, err := ftms.DiscoverCharacteristics(
				[]bluetooth.UUID{ftmsControlPoint},
			)
			if err != nil {
				continueScanning("Could not scan all characteristics of ftms service")
			}

			devChar <- &(chars[0])
			ftmsCP <- &(fmtsControlPointChar[0])
		}
	}()

	defer func() {
		err := adapter.StopScan()
		if err != nil {
			slog.Info("Could not stop scanning")
		}
	}()

	// start scanning
	scan <- true

	char := <-devChar
	slog.Info("Scanning done...")
	if char == nil {
		return nil, fmt.Errorf("Supported device not found")
	}

	return char, nil
}
