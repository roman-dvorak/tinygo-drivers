// Package tmc5130 implements a driver for the TMC5130 stepper driver.
//
// Datasheet: https://www.trinamic.com/fileadmin/assets/Products/ICs_Documents/TMC5130_datasheet_Rev1.16.pdf
//
//
package tmc5130 // import "tinygo.org/x/drivers/tmc5130"

import (
	_"errors"
	//"fmt"
	"machine"
	"time"
	"tinygo.org/x/drivers"
)

func constrain (val int, min int, max int) int {
	if(val < min){return min}
	if(val > max){return max}
	return val
}


type Device struct {
	bus drivers.SPI
	cs  machine.Pin
	tx  []byte
	rx  []byte
	Spi_status_d SPI_STATUS
}

const (
	bufferSize int = 64
)

func New(b drivers.SPI, csPin machine.Pin) *Device {
	d := &Device{
			bus: b,
			tx:  make([]byte, 0, bufferSize),
			rx:  make([]byte, 0, bufferSize),
			cs:  csPin,
	}

	return d
}


func (d *Device) Configure() {
	d.cs.Configure(machine.PinConfig{Mode: machine.PinOutput})
}

func (d *Device) GetRegister(addr byte) []byte {
  a := []byte{addr, 0, 0, 0, 0}
  b := []byte{0, 0, 0, 0, 0}
  c := []byte{0, 0, 0, 0, 0}

	d.cs.Low(); d.bus.Tx(a, b); d.cs.High()
  time.Sleep(time.Microsecond * 1)
	d.cs.Low(); d.bus.Tx(a, c); d.cs.High()
	d.Spi_status_d.decode(c[0])

	return c
}

func (d *Device) GetSPIstatus() SPI_STATUS {
	return d.Spi_status_d
}


func (d *Device) SetRegister(addr byte, data int) ([]byte) {

  a := []byte{addr, byte(data>>24)&0xff, byte(data>>16)&0xff, byte(data>>8)&0xff, byte(data>>0)&0xff}
  b := []byte{0, 0, 0, 0, 0}
	c := []byte{0, 0, 0, 0, 0}

	d.cs.Low()
	d.bus.Tx(a, b)
	d.cs.High()
  time.Sleep(time.Millisecond * 1)
	d.cs.Low()
	d.bus.Tx(a, c)
	d.cs.High()
	d.Spi_status_d.decode(b[0])

	return c
}



func (d *Device) InputStatus() (REG_INPUT) {
		reg := d.GetRegister(IOIN)
		println(reg[0], reg[1], reg[2], reg[3])
		var data REG_INPUT
		data.decode(reg)
		return data
}

func (d *Device) GetXACTUAL() (REG_XACTUAL) {
		reg := d.GetRegister(XACTUAL)
		var data REG_XACTUAL
		data.decode(reg)
		return data
}

func (d *Device) GetVACTUAL() (REG_VACTUAL) {
		reg := d.GetRegister(VACTUAL)
		var data REG_VACTUAL
		data.decode(reg)
		return data
}

func (d *Device) SetXACTUAL(position int) (REG_XACTUAL) {
		reg := d.SetRegister(XACTUAL|WRITE, position)
		var data REG_XACTUAL
		data.decode(reg)
		return data
}

func (d *Device) SetXTARGET(position int) (REG_XTARGET) {
		reg := d.SetRegister(XTARGET|WRITE, position)
		var data REG_XTARGET
		data.decode(reg)
		return data
}

func (d *Device) SetCurrent(ihold int, irun int, delay int) {
	ihold = constrain(ihold, 0, 31)
	irun = constrain(irun, 0, 31)
	delay = constrain(delay, 0, 2^18)
	if(delay==0){
		delay = 1
	} else {
		delay = delay<<1
	}
	d.SetRegister(IHOLD_IRUN|WRITE,(ihold &0b11111)<<0|(irun &0b11111)<<8|(delay&0b1111)<<16);

}
