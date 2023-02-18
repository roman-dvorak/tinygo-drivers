// Package ds18b20 provides a driver for the DALLAS one-wire temperature sensor
//
// Datasheet: https://datasheets.maximintegrated.com/en/ds/DS18B20.pdf
//

package ds18b20 // import "tinygo.org/x/drivers/ds18b20"

import (
	"machine"
	"time"
  //"fmt"
)

type Device struct {
	pin  machine.Pin
}

func New(pin machine.Pin) Device {
	return Device{
		pin:  pin,
	}
}

func (d *Device) ConvertT() {
  d.WriteByte(0x44)
}

func (d *Device) SkipRom() {
  // 0xcc
  d.WriteByte(0xcc)
}

func (d *Device) ReadScratchpad() (int) {
  // 0xBE
  d.WriteByte(0xBE)
  var data [9]byte

  for i := 0; i < 8; i++ {
    data[i] = d.ReadByte()
  }

  //println(int(data[0]), int(data[1]))
  var temp int16
  //var ftemp float64

  temp |= int16(data[0])
  temp |= (int16(data[1])&0x0f) << 8
  // TODO negative values handeling
  if (int(data[1])&0xf0 > 0){
  //   temp *= -1;
    temp = temp/16-256
  } else {
    temp = temp/16
  }

  return (int(temp))
}

func b2i(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

func (d *Device) ReadByte() (byte uint8) {
  var data uint8
  data = 0
  for i := 0; i < 8; i++ {
    data |= (b2i(d.ReadBit()) << (i))
  }
  return data
}

func (d *Device) WriteByte(data byte) {
  for i := 0; i < 8; i++ {
    d.WriteBit( ((data>>i)&0x01) != 0 )
  }
}

func (d *Device) WriteBit(bit bool){
  d.pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
  if(bit){
    d.pin.Set(false)
    time.Sleep(time.Microsecond * 5);
    d.pin.Set(true)
    time.Sleep(time.Microsecond * 60);
    } else {
    d.pin.Set(false)
    time.Sleep(time.Microsecond * 70);
    d.pin.Set(true)
    time.Sleep(time.Microsecond * 2);
    }
}

func (d *Device) ReadBit() (bit bool){
  d.pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
  d.pin.Set(false)
  time.Sleep(time.Microsecond * 6);
  d.pin.Set(true)
  d.pin.Configure(machine.PinConfig{Mode: machine.PinInput})
  time.Sleep(time.Microsecond * 9);
  val := d.pin.Get()
  time.Sleep(time.Microsecond * 50);
  return val
}

func (d *Device) Reset() (err error) {
  //println("RESET")
  d.pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
  d.pin.Set(false)
	time.Sleep(time.Microsecond * 1000)
  //d.pin.Set(true)
  d.pin.Configure(machine.PinConfig{Mode: machine.PinInput})

  for true {
    if(d.pin.Get() == true) {
      break
      //wait_time := time.Now()
    }
  }
	time.Sleep(time.Microsecond * 5)
  for true {
    //println("RESPONSE", d.pin.Get())
    if(d.pin.Get() == false){
      //println("RESPONSE", d.pin.Get())
      break
    }
  }
	time.Sleep(time.Microsecond * 5)
  for true {
    //println("RESPONSE", d.pin.Get())
    if(d.pin.Get() == true){
      //println("DONE", d.pin.Get())
      break
    }
  }

	return
}
