package bluetooth

import (
	"fmt"

	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter
var zwiftHubUUID = "c96fb5f7-b4d5-262e-7baf-0a479225f3ab"
var ftmsUUID = "00001826-0000-1000-8000-00805f9b34fb"

// notice that the only thing that's different from the ftmsUUID is the first segment.
// this is the case wit hall uuids
var instantaneousPower = "00002ad2-0000-1000-8000-00805f9b34fb"

func Scan() { // Enable BLE interface.
	var zwiftHub *bluetooth.ScanResult
	err := adapter.Enable()
	if err != nil {
		panic(err)
	}

	fmt.Println("Finding trainer...")
	done := make(chan bool)
	go func() {
		err := adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
			if device.Address.String() == zwiftHubUUID {
				zwiftHub = &device
				done <- true
			}
		})

		if err != nil {
			panic(err)
		}
	}()

	<-done
	err = adapter.StopScan()
	if err != nil {
		panic(err)
	}
	fmt.Println("found trainer")

	device, err := adapter.Connect(zwiftHub.Address, bluetooth.ConnectionParams{})
	if err != nil {
		panic("Cannot connect to device")
	}

	fmt.Println("Connected")

	ftms, err := bluetooth.ParseUUID(ftmsUUID)
	if err != nil {
		panic(err)
	}

	dservices, err := device.DiscoverServices([]bluetooth.UUID{ftms})
	if err != nil {
		panic("not ftmsReady")
	}

	ftmsOk := len(dservices) == 1
	if !ftmsOk {
		panic("not ftmsReady")
	}

	service := dservices[0]

	fmt.Println("FTMS ready")

	serviceUuid, err := bluetooth.ParseUUID(instantaneousPower)
	if err != nil {
		panic("could not get service UUID")
	}

	chars, err := service.DiscoverCharacteristics([]bluetooth.UUID{serviceUuid})
	if err != nil {
		panic("Could not get characteristics")
	}

	charsOk := len(chars) == 1
	if !charsOk {
		panic("Device does not have instantaneous power characteristic")
	}

	char := chars[0]
	char.EnableNotifications(func(buf []byte) {
		println("data:", uint8(buf[1]))
	})

	// block
	select {}
}
