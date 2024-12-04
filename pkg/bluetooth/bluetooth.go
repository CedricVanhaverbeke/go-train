package bluetooth

import (
	"fmt"
	"sync"

	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter
var zwiftHubUUID = "c96fb5f7-b4d5-262e-7baf-0a479225f3ab"
var ftmsUUID = "00001826-0000-1000-8000-00805f9b34fb"

// notice that the only thing that's different from the ftmsUUID is the first segment.
// this is the case wit hall uuids
var instantaneousPower = "00002ad2-0000-1000-8000-00805f9b34fb"

func Scan() { // Enable BLE interface.
	err := adapter.Enable()
	if err != nil {
		panic(err)
	}

	fmt.Println("Finding trainer...")
	char := DiscoverFTMSDevice()

	err = char.EnableNotifications(func(buf []byte) {
		println("data:", uint8(buf[1]))
	})

	if err != nil {
		panic(err)
	}
}

// DiscoverFTMSDevice checks every available device
// having the FTMS service. It returns the
// first device that has FTMS enabled and
// instantaneous power
func DiscoverFTMSDevice() bluetooth.DeviceCharacteristic {
	var wg sync.WaitGroup
	found := make(chan bluetooth.ScanResult, 1)
	devChar := make(chan bluetooth.DeviceCharacteristic)

	ftms, err := bluetooth.ParseUUID(ftmsUUID)
	if err != nil {
		panic(err)
	}

	serviceUuid, err := bluetooth.ParseUUID(instantaneousPower)
	if err != nil {
		panic("could not get service UUID")
	}

	wg.Add(1)
	go func() {
		err := adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
			fmt.Println("Found device, checking FTMS")
			found <- device
		})

		if err != nil {
			panic(err)
		}
	}()

	go func() {
		for scanResult := range found {
			fmt.Println("Checking device")
			device, err := adapter.Connect(scanResult.Address, bluetooth.ConnectionParams{})
			if err != nil {
				fmt.Println("Cannot connect to device")
				continue
			}

			dservices, err := device.DiscoverServices([]bluetooth.UUID{ftms})
			if err != nil {
				fmt.Println("Device is not ftms ready")
				continue
			}

			ftmsOk := len(dservices) == 1
			if !ftmsOk {
				fmt.Println("Device is not ftms ready")
				continue
			}

			service := dservices[0]

			chars, err := service.DiscoverCharacteristics([]bluetooth.UUID{serviceUuid})
			if err != nil {
				fmt.Println("Could not get characteristics")
				continue
			}

			charsOk := len(chars) == 1
			if !charsOk {
				fmt.Println("Device does not have instantaneous power characteristic")
				continue
			}

			char := chars[0]
			devChar <- char
			wg.Done()
		}
	}()

	defer func() {
		err := adapter.StopScan()
		if err != nil {
			fmt.Println("Could not stop scanning")
		}
	}()

	wg.Wait()
	return <-devChar
}
