package bluetooth

import (
	"encoding/binary"

	"tinygo.org/x/bluetooth"
)

type listeners struct {
	listeners []chan int
}

func (l *listeners) AddListener(c chan int) bool {
	if l.listeners == nil {
		l.listeners = make([]chan int, 0)
	}

	l.listeners = append(l.listeners, c)
	return true
}

func (l *listeners) WriteValue(v int) {
	for _, listener := range l.listeners {
		listener <- v
	}
}

type powerCharacteristic struct {
	readPwr  *bluetooth.DeviceCharacteristic
	writePwr *bluetooth.DeviceCharacteristic

	listeners
}

// ContinuousRead extracts instantaneous power (signed 16-bit integer, little-endian)
func (p *powerCharacteristic) ContinuousRead() error {
	err := p.readPwr.EnableNotifications(func(buf []byte) {
		p.listeners.WriteValue(decode(buf))
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
	return p.writePwr.Write(encode(power))
}

func encode(power int) []byte {
	if power < 0 {
		power = power * -1
	}

	data := []byte{0x05} // opcode for setting power
	// since power is postivie this shouldn't matter
	data = binary.LittleEndian.AppendUint16(data, uint16(power))
	return data
}

func decode(buf []byte) int {
	buf = buf[2:4]
	power := binary.LittleEndian.Uint16(buf)
	return int(power)
}

// requestControl initiates the procedure
// to request control over the fitness machine
func (p *powerCharacteristic) requestControl() error {
	data := []byte{0x00}
	_, err := p.writePwr.Write(data)
	return err
}
