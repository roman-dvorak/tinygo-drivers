
package tmc5130
import (
  "fmt"
)

// Registers
const (

  // ADD 'WRITE' in case of WRITE command
    WRITE     = 0x80

    GCONF     = 0x00
    GSTAT     = 0x01
    IFCNT     = 0x02
    SLAVECONF = 0x03
    IOIN      = 0x04
    OUTPUT    = 0x04
    X_COMPARE = 0x05
    IHOLD_IRUN= 0x10
    TPOWERDOWN=0x11
    TSTEP     = 0x12
    TPWMTHRS  = 0x13
    TCOOLTHRS = 0x14
    THIGH     = 0x15
    RAMPMODE  = 0x20
    XACTUAL   = 0x21
    VACTUAL   = 0x22
    VSTART    = 0x23
    A1        = 0x24
    V1        = 0x25
    AMAX      = 0x26
    VMAX      = 0x27
    DMAX      = 0x28
    D1        = 0x2A
    VSTOP     = 0x2B
    TZEROWAIT = 0x2C
    XTARGET   = 0x2D
    VDCMIN    = 0x33
    SW_MODE   = 0x34
    RAMP_STAT = 0x35
    XLATCH    = 0x36
  // Encoder registers
    ENCMODE   = 0x38
    X_END     = 0x39
    ENC_CONST = 0x3A
    ENC_STATUS= 0x3B
    ENC_LATCH = 0x3C
  // Motor driver registers
    MSLUT0    = 0x60
    MSLUT1    = 0x61
    MSLUT2    = 0x62
    MSLUT3    = 0x63
    MSLUT4    = 0x64
    MSLUT5    = 0x65
    MSLUT6    = 0x66
    MSLUT7    = 0x67
    MSLUTSEL  = 0x68
    MSLUTSTART= 0x69
    MSCNT     = 0x6A
    MSCURACT  = 0x6B
    CHOPCONF  = 0x6C
    COOLCONF  = 0x6D
    DCCTRL    = 0x6E
    DRV_STATUS= 0x6F
    PWM_CONF  = 0x70
    PWM_SCALE = 0x71
    ENCM_CTRL = 0x72
    LOST_STEPS= 0x73
  )


  type SPI_STATUS struct {
      Status_stop_r bool
      Status_stop_l bool
      Position_reached bool
      Velocity_reached bool
      Standstill bool
      Sg2 bool
      Driver_error bool
      Reset_flag bool
  }
  func (v *SPI_STATUS) decode(b byte) {
  		v.Status_stop_r= ((b&0b10000000) > 0)
  		v.Status_stop_l= ((b&0b01000000) > 0)
  		v.Position_reached= ((b&0b00100000) > 0)
  		v.Velocity_reached= ((b&0b00010000) > 0)
  		v.Standstill= ((b&0b00001000) > 0)
  		v.Sg2= ((b&0b00000100) > 0)
  		v.Driver_error= ((b&0b00000010) > 0)
  		v.Reset_flag= ((b&0b00000001) > 0)
  }
  func (v *SPI_STATUS) String() string{
      return fmt.Sprintf("%#v", v)
  }


  type REG_INPUT struct {
      REFL_STEP bool
      REFR_DIR bool
      ENCB_DCEN bool
      ENCA_DCIN bool
      DRV_ENN bool
      ENC_N bool
      SD_MODE bool
      SWCOMP_IN bool
      VERSION uint8
  }
  func (v *REG_INPUT) decode(data []byte) {
      b:=data[4]
  		v.REFL_STEP   = ((b&0b00000001) > 0)
  		v.REFR_DIR    = ((b&0b00000010) > 0)
  		v.ENCB_DCEN   = ((b&0b00000100) > 0)
  		v.ENCA_DCIN   = ((b&0b00001000) > 0)
  		v.DRV_ENN     = ((b&0b00010000) > 0)
  		v.ENC_N       = ((b&0b00100000) > 0)
  		v.SD_MODE     = ((b&0b01000000) > 0)
  		v.SWCOMP_IN   = ((b&0b10000000) > 0)
      v.VERSION = data[1]
  }
  func (v *REG_INPUT) String() string{
      return fmt.Sprintf("%#v", v)
  }


  type REG_XACTUAL struct {
      XACTUAL int32
  }
  func (v *REG_XACTUAL) decode(data []byte) {
      var val int32
      val=int32(int(data[4]))
      val|=int32(int(data[3])<<8)
      val|=int32(int(data[2])<<16)
      val|=int32(int(data[1])<<24)
      v.XACTUAL = val
  }
  func (v *REG_XACTUAL) String() string{
      return fmt.Sprintf("%#v", v)
  }

  type REG_XTARGET struct {
      XTARGET int32
  }
  func (v *REG_XTARGET) decode(data []byte) {
      var val int32
      val=int32(int(data[4]))
      val|=int32(int(data[3])<<8)
      val|=int32(int(data[2])<<16)
      val|=int32(int(data[1])<<24)
      v.XTARGET = val
  }
  func (v *REG_XTARGET) String() string{
      return fmt.Sprintf("%#v", v)
  }



  type REG_VACTUAL struct {
      VACTUAL int32
  }
  func (v *REG_VACTUAL) decode(data []byte) {
      var val int32
      val=int32(int(data[4]))
      val|=int32(int(data[3])<<8)
      val|=int32(int(data[2])<<16)
      val|=int32(int(data[1])<<24)
      v.VACTUAL = val
  }
  func (v *REG_VACTUAL) String() string{
      return fmt.Sprintf("%#v", v)
  }
