package bluetooth

import (
	"fmt"

	"tinygo.org/x/bluetooth"
)

type powerCharacteristic struct {
	readPwr  *bluetooth.DeviceCharacteristic
	writePwr *bluetooth.DeviceCharacteristic

	listeners []chan int
}

func (p *powerCharacteristic) AddListener(c chan int) {
	if p.listeners == nil {
		p.listeners = make([]chan int, 0)
	}

	p.listeners = append(p.listeners, c)
}

// ContinuousRead extracts instantaneous power (signed 16-bit integer, little-endian)
func (p *powerCharacteristic) ContinuousRead() error {
	err := p.readPwr.EnableNotifications(func(buf []byte) {
		power := int16(buf[2]) | int16(buf[3])<<8
		for _, listener := range p.listeners {
			listener <- int(power)
		}
	})

	return err
}

// for writing power I first need to get permission to do so
// see page 55 of this: file:///Users/cedricvanhaverbeke/Downloads/FTMS_v1.0.1.pdf

// This procedure requires control permission in order to be executed. Refer to Section 4.16.2.1 for more
// information on the Request Control procedure.
// When the Set Target Power Op Code is written to the Fitness Machine Control Point and the Result Code
// is ‘Success’, the Server shall set the target power to the value sent as a Parameter.
// see page 74 to see how the interaction works
func (p *powerCharacteristic) Write(power int) (int, error) {
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

	return p.writePwr.Write(data)
}

// requestControl initiates the procedure
// to request control over the fitness machine
func (p *powerCharacteristic) requestControl() error {
	data := []byte{0x00}
	_, err := p.writePwr.Write(data)
	return err
}
