// Package tmc5130 implements a driver for the TMC5130 stepper driver.
//
// Datasheet:
//
//
package tmc5130 // import "tinygo.org/x/drivers/tmc5130"

import (
	_"errors"
	_"fmt"
	"machine"
	"time"

	"tinygo.org/x/drivers"
)

type Device struct {
	bus drivers.SPI
	cs  machine.Pin
	tx  []byte
	rx  []byte
}

const (
	bufferSize int = 64
)
//
// // New returns a new MCP2515 driver. Pass in a fully configured SPI bus.
func New(b drivers.SPI, csPin machine.Pin) *Device {
	d := &Device{
			bus: b,
			tx:  make([]byte, 0, bufferSize),
			rx:  make([]byte, 0, bufferSize),

		cs:  csPin,
	}

	return d
}
//
// // Configure sets up the device for communication.
func (d *Device) Configure() {
	d.cs.Configure(machine.PinConfig{Mode: machine.PinOutput})
}


func (d *Device) GetXACT(addr byte) []byte {
	// d.tx[0] = 0x21
	// d.tx[1] = 0x00
	// d.tx[2] = 0x00
	// d.tx[4] = 0x00
  a := []byte{addr, 0, 0xaa, 0xaa, 0}
  b := []byte{0xaa, 0xaa, 0xaa, 0xaa, 0xaa}

	d.cs.Low()
	//d.bus.Tx(a, b)
	b[0], _ = d.bus.Transfer(a[0])
	b[1], _ = d.bus.Transfer(a[1])
	b[2], _ = d.bus.Transfer(a[2])
	b[3], _ = d.bus.Transfer(a[3])
	b[4], _ = d.bus.Transfer(a[4])
	d.cs.High()

	time.Sleep(time.Millisecond * 5)

	return b
}

func (d *Device) SetRegister(addr byte, data []byte) []byte {

  a := []byte{addr, data[0], data[1], data[2], data[3]}
  b := []byte{0xaa, 0xaa, 0xaa, 0xaa, 0xaa}

	d.cs.Low()
	//d.bus.Tx(a, b)
	b[0], _ = d.bus.Transfer(a[0])
	b[1], _ = d.bus.Transfer(a[1])
	b[2], _ = d.bus.Transfer(a[2])
	b[3], _ = d.bus.Transfer(a[3])
	b[4], _ = d.bus.Transfer(a[4])
	d.cs.High()


	return b
}
