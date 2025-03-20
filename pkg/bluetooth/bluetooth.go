package bluetooth

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

func Connect() (*Device, error) {
	err := adapter.Enable()
	if err != nil {
		panic(err)
	}

	slog.Info("Finding trainer...")
	readPowerChar, writePowerChar, err := discover()
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

	err = powerChar.requestControl()
	if err != nil {
		slog.Error(err.Error())
	}

	trainer := NewDevice(WithPower(&powerChar))
	return &trainer, nil
}

// discover checks every available device
// having the FTMS service and cyling power service.
//
//	It returns two bluetooth characteristics.
//
// the first char can be used to get power notifications and
// the second char can be used to set power on the device
func discover() (*bluetooth.DeviceCharacteristic, *bluetooth.DeviceCharacteristic, error) {
	var powChar *bluetooth.DeviceCharacteristic
	var ftmsChar *bluetooth.DeviceCharacteristic
	done := make(chan struct{})

	scanned := map[string]bool{}
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		slog.Info("Scanning bluetooth devices...")
		err := adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
			if scanned[device.Address.String()] {
				return
			}

			scanned[device.Address.String()] = true

			slog.Info("Found device with uuid: " + device.Address.String())

			go func() {
				pChar, ftmsControlPointChar, err := verifyDevice(device)
				if err != nil {
					slog.Info(err.Error())
					return
				}

				powChar = pChar
				ftmsChar = ftmsControlPointChar
				close(done)
			}()
		})

		if err != nil {
			panic(err)
		}
	}()

	defer func() {
		err := adapter.StopScan()
		if err != nil {
			slog.Info("Could not stop scanning")
		}
	}()

	select {
	case <-done:
		wg.Done()
	case <-time.After(time.Second * 10):
		return nil, nil, fmt.Errorf("Bluetooth deadline exceeded, no devices found...")
	}

	wg.Wait()
	slog.Info("Scanning done...")
	return powChar, ftmsChar, nil
}

func verifyDevice(
	scanResult bluetooth.ScanResult,
) (*bluetooth.DeviceCharacteristic, *bluetooth.DeviceCharacteristic, error) {
	slog.Info("Checking device...")

	device, err := adapter.Connect(
		scanResult.Address,
		bluetooth.ConnectionParams{
			ConnectionTimeout: bluetooth.NewDuration(2 * time.Second),
		},
	)

	if err != nil {
		return nil, nil, err
	}

	dservices, err := device.DiscoverServices(
		[]bluetooth.UUID{powServiceUuid, ftmsServiceUuid},
	)
	if err != nil {
		return nil, nil, fmt.Errorf("Device does not have cycling power enabled")
	}

	hasCyclingPowerAndFTMS := len(dservices) == 2
	if !hasCyclingPowerAndFTMS {
		return nil, nil, fmt.Errorf("Device does not have all required services")
	}

	cyclingPowerService, ftmsService := dservices[0], dservices[1]
	pChar, err := getChar(&cyclingPowerService, cyclingPowerCharacteristicUuid)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not get characteristics " + err.Error())
	}

	ftmsControlPointChar, err := getChar(&ftmsService, FTMSCharUuid)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not scan all characteristics of ftms service")
	}

	return &pChar, &ftmsControlPointChar, nil
}

func getChar(
	service *bluetooth.DeviceService,
	charUuid bluetooth.UUID,
) (bluetooth.DeviceCharacteristic, error) {
	chars, err := service.DiscoverCharacteristics(
		[]bluetooth.UUID{charUuid},
	)

	if err != nil {
		return bluetooth.DeviceCharacteristic{}, err
	}

	charsOk := len(chars) == 1
	if !charsOk {
		return bluetooth.DeviceCharacteristic{}, fmt.Errorf("Service does not have characteristic")
	}

	return chars[0], nil
}
