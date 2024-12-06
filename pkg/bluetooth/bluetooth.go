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

// notice that the only thing that's different from the ftmsUUID is the first segment.
// this is the case wit hall uuids
var instantaneousPower = "00002ad2-0000-1000-8000-00805f9b34fb"

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
		log.Printf(" %b", buf[:])
		power := int(buf[2]) | int(buf[3])<<8 // Little-endian
		fmt.Println(power / 10)
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
	scan := make(chan bool)

	scanned := map[string]bool{}

	ftms, err := bluetooth.ParseUUID(ftmsUUID)
	if err != nil {
		return nil, err
	}

	serviceUuid, err := bluetooth.ParseUUID(instantaneousPower)
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
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
			slog.Info("Checking device FTMS")

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

			dservices, err := device.DiscoverServices([]bluetooth.UUID{ftms})
			if err != nil {
				continueScanning("Device is not ftms ready")
				continue
			}

			ftmsOk := len(dservices) == 1
			if !ftmsOk {
				continueScanning("Device is not ftms ready")
				continue
			}

			service := dservices[0]

			chars, err := service.DiscoverCharacteristics([]bluetooth.UUID{serviceUuid})
			if err != nil {
				continueScanning("Could not get characteristics")
				continue
			}

			fmt.Printf("%+v", chars)

			charsOk := len(chars) == 1
			if !charsOk {
				continueScanning("Device does not have instantaneous power characteristic")
				continue
			}

			char := chars[0]
			devChar <- &char
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
		return nil, fmt.Errorf("FTMS device not found")
	}

	return char, nil
}
