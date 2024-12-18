package bluetooth

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

type powerCharacteristic struct {
	readPwr  *bluetooth.DeviceCharacteristic
	writePwr *bluetooth.DeviceCharacteristic
}

// ContinuousRead extracts instantaneous power (signed 16-bit integer, little-endian)
func (p powerCharacteristic) ContinuousRead(c chan int) error {
	err := p.readPwr.EnableNotifications(func(buf []byte) {
		power := int16(buf[2]) | int16(buf[3])<<8
		c <- int(power)
	})

	return err
}

func (p powerCharacteristic) Write(power int) (int, error) {
	// i need to set the power in little endian
	data := []byte{0x05, 0, 0}
	i := 1

	// bitshift it to the right the
	if power > 256 {
		i++
		data[1] = 255
		power -= 255
	}

	data[i] = byte(power)
	fmt.Printf("%b\n", data)

	return p.writePwr.WriteWithoutResponse(data)
}

func Connect() (*Trainer, error) {
	err := adapter.Enable()
	if err != nil {
		panic(err)
	}

	slog.Info("Finding trainer...")
	readPowerChar, writePowerChar, err := discoverPowerDevice()
	if err != nil {
		return nil, err
	}

	err = writePowerChar.EnableNotifications(func(buf []byte) {
		println("Notification received:", buf)
	})

	if err != nil {
		slog.Error(err.Error())
	}

	powerChar := powerCharacteristic{
		readPwr:  readPowerChar,
		writePwr: writePowerChar,
	}

	trainer := NewTrainer(powerChar)
	return &trainer, nil
}

// discoverPowerDevice checks every available device
// having the FTMS service and cyling power service.
//
//	It returns two bluetooth characteristics.
//
// the first char can be used to get power notifications and
// the second char can be used to set power on the device
func discoverPowerDevice() (*bluetooth.DeviceCharacteristic, *bluetooth.DeviceCharacteristic, error) {
	found := make(chan bluetooth.ScanResult)
	devChar := make(chan *bluetooth.DeviceCharacteristic)
	ftmsCP := make(chan *bluetooth.DeviceCharacteristic)
	scan := make(chan bool)

	scanned := map[string]bool{}

	// inner function that starts a new scan operation
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

			dservices, err := device.DiscoverServices(
				[]bluetooth.UUID{cyclingPowerService, ftmsService},
			)
			if err != nil {
				continueScanning("Device does not have cycling power enabled")
				continue
			}

			hasCyclingPowerAndFTMS := len(dservices) == 2
			if !hasCyclingPowerAndFTMS {
				continueScanning("Device does not have all required services")
				continue
			}

			// let's assume the services get fetched in order
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

			ftms := dservices[1]
			ftmsControlPointChar, err := ftms.DiscoverCharacteristics(
				[]bluetooth.UUID{ftmsControlPoint},
			)
			if err != nil {
				continueScanning("Could not scan all characteristics of ftms service")
			}

			devChar <- &(chars[0])
			ftmsCP <- &(ftmsControlPointChar[0])
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
	ftmsCPChar := <-ftmsCP
	slog.Info("Scanning done...")
	if char == nil {
		return nil, nil, fmt.Errorf("Supported device not found")
	}

	return char, ftmsCPChar, nil
}
