package bluetooth

import (
	"tinygo.org/x/bluetooth"
)

var (
	ftmsUUID                       = "00001826-0000-1000-8000-00805f9b34fb"
	fitnessMachineControlPointUUID = "00002ad9-0000-1000-8000-00805f9b34fb"
	cyclingPower                   = "00001818-0000-1000-8000-00805f9b34fb"

	cyclingSpeedAndCadence = "00001816-0000-1000-8000-00805f9b34fb"

	// notice that the only thing that's different from the ftmsUUID is the first segment.
	// this is the case with all uuids
	// file:///Users/cedricvanhaverbeke/Downloads/GATT_Specification_Supplement_v5.pdf
	// https://gist.github.com/sam016/4abe921b5a9ee27f67b3686910293026
	// var indoorBikeData = "00002ad2-0000-1000-8000-00805f9b34fb"
	cyclingPowerMeasureMent = "00002a63-0000-1000-8000-00805f9b34fb"

	// this stuff should not be available in the whole bluetooth package
	// instead we should have one struct with uuids or something
	// so we can require them in a list, and bundle them
	ftmsServiceUuid         bluetooth.UUID
	powServiceUuid          bluetooth.UUID
	speedCadenceServiceUuid bluetooth.UUID

	// these are the characteristics itself
	FTMSCharUuid                   bluetooth.UUID
	cyclingPowerCharacteristicUuid bluetooth.UUID
)

// initializes id's. This will normally always work.
func init() {
	var err error
	ftmsServiceUuid, err = bluetooth.ParseUUID(ftmsUUID)
	if err != nil {
		panic(err)
	}

	FTMSCharUuid, err = bluetooth.ParseUUID(fitnessMachineControlPointUUID)
	if err != nil {
		panic(err)
	}

	powServiceUuid, err = bluetooth.ParseUUID(cyclingPower)
	if err != nil {
		panic(err)
	}

	speedCadenceServiceUuid, err = bluetooth.ParseUUID(cyclingSpeedAndCadence)
	if err != nil {
		panic(err)
	}

	cyclingPowerCharacteristicUuid, err = bluetooth.ParseUUID(cyclingPowerMeasureMent)
	if err != nil {
		panic(err)
	}
}
