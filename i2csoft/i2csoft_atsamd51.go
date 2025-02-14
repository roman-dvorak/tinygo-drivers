//go:build atsamd51 || atsame5x

package i2csoft

import (
	"device"
)

// wait waits for half the time of the SCL operation interval. It is set to
// about 100 kHz.
func (i2c *I2C) wait() {
	wait := 20
	for i := 0; i < wait; i++ {
		device.Asm(`nop`)
	}
}
