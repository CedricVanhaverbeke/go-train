package bluetooth

import (
	"tinygo.org/x/bluetooth"
)

var (
	ftmsUUID                       = "00001826-0000-1000-8000-00805f9b34fb"
	fitnessMachineControlPointUUID = "00002ad9-0000-1000-8000-00805f9b34fb"
	cyclingPower                   = "00001818-0000-1000-8000-00805f9b34fb"

	// notice that the only thing that's different from the ftmsUUID is the first segment.
	// this is the case wit hall uuids
	// file:///Users/cedricvanhaverbeke/Downloads/GATT_Specification_Supplement_v5.pdf
	// https://gist.github.com/sam016/4abe921b5a9ee27f67b3686910293026
	// var indoorBikeData = "00002ad2-0000-1000-8000-00805f9b34fb"
	cyclingPowerMeasureMent = "00002a63-0000-1000-8000-00805f9b34fb"

	ftmsService         bluetooth.UUID
	ftmsControlPoint    bluetooth.UUID
	cyclingPowerService bluetooth.UUID
	serviceUuid         bluetooth.UUID
)

// initializes id's. This will normally always work.
func init() {
	var err error
	ftmsService, err = bluetooth.ParseUUID(ftmsUUID)
	if err != nil {
		panic(err)
	}

	ftmsControlPoint, err = bluetooth.ParseUUID(fitnessMachineControlPointUUID)
	if err != nil {
		panic(err)
	}

	cyclingPowerService, err = bluetooth.ParseUUID(cyclingPower)
	if err != nil {
		panic(err)
	}

	serviceUuid, err = bluetooth.ParseUUID(cyclingPowerMeasureMent)
	if err != nil {
		panic(err)
	}
}
