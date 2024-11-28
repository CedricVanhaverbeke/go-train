package bluetooth

import (
	"fmt"

	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter
var zwiftHubUUID = "c96fb5f7-b4d5-262e-7baf-0a479225f3ab"

func Scan() { // Enable BLE interface.
	err := adapter.Enable()
	if err != nil {
		panic(err)
	}

	fmt.Println("Finding trainer...")
	connectChan := make(chan bool)
	go func() {
		err := adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
			if device.Address.String() == zwiftHubUUID {
				println("found device:", device.Address.String(), device.RSSI, device.LocalName())
				connectChan <- true
			}
		})

		if err != nil {
			panic(err)
		}
	}()

	<-connectChan
	err = adapter.StopScan()
	if err != nil {
		panic(err)
	}
	fmt.Println("found trainer")
}
