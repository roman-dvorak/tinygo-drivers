package lorawan

import (
	"errors"

	"tinygo.org/x/drivers/lora"
	"tinygo.org/x/drivers/lora/lorawan/region"
)

var (
	ErrNoJoinAcceptReceived    = errors.New("no JoinAccept packet received")
	ErrNoRadioAttached         = errors.New("no LoRa radio attached")
	ErrInvalidEuiLength        = errors.New("invalid EUI length")
	ErrInvalidAppKeyLength     = errors.New("invalid AppKey length")
	ErrInvalidPacketLength     = errors.New("invalid packet length")
	ErrInvalidDevAddrLength    = errors.New("invalid DevAddr length")
	ErrInvalidMic              = errors.New("invalid Mic")
	ErrFrmPayloadTooLarge      = errors.New("FRM payload too large")
	ErrInvalidNetIDLength      = errors.New("invalid NetID length")
	ErrInvalidNwkSKeyLength    = errors.New("invalid NwkSKey length")
	ErrInvalidAppSKeyLength    = errors.New("invalid AppSKey length")
	ErrUndefinedRegionSettings = errors.New("undefined Regionnal Settings ")
)

const (
	LORA_TX_TIMEOUT = 2000
	LORA_RX_TIMEOUT = 5000
)

var (
	ActiveRadio    lora.Radio
	Retries        = 15
	regionSettings region.RegionSettings
)

func UseRegionSettings(rs region.RegionSettings) {
	regionSettings = rs
}

func UseRadio(r lora.Radio) {
	if ActiveRadio != nil {
		panic("lorawan.ActiveRadio is already set")
	}
	ActiveRadio = r
}

func Join(otaa *Otaa, session *Session) error {
	var resp []uint8

	if ActiveRadio == nil {
		return ErrNoRadioAttached
	}

	if regionSettings == nil {
		return ErrUndefinedRegionSettings
	}

	otaa.Init()

	// Send join packet
	payload, err := otaa.GenerateJoinRequest()
	if err != nil {
		return err
	}

	// Prepare radio for Join Tx
	joinChannel := regionSettings.GetJoinRequestChannel()
	ActiveRadio.SetFrequency(joinChannel.Frequency)
	ActiveRadio.SetBandwidth(joinChannel.Bandwidth)
	ActiveRadio.SetCodingRate(joinChannel.CodingRate)
	ActiveRadio.SetSpreadingFactor(joinChannel.SpreadingFactor)
	ActiveRadio.SetCrc(true)
	ActiveRadio.SetIqMode(0) // IQ Standard
	ActiveRadio.Tx(payload, LORA_TX_TIMEOUT)
	if err != nil {
		return err
	}

	// Wait for JoinAccept
	joinChannel = regionSettings.GetJoinAcceptChannel()
	ActiveRadio.SetFrequency(joinChannel.Frequency)
	ActiveRadio.SetBandwidth(joinChannel.Bandwidth)
	ActiveRadio.SetCodingRate(joinChannel.CodingRate)
	ActiveRadio.SetSpreadingFactor(joinChannel.SpreadingFactor)
	ActiveRadio.SetIqMode(1) // IQ Inverted
	for i := 0; i < Retries; i++ {
		resp, err = ActiveRadio.Rx(LORA_RX_TIMEOUT)
		if err != nil {
			return err
		}
		if resp != nil {
			break
		}
	}
	if resp == nil {
		return ErrNoJoinAcceptReceived
	}

	err = otaa.DecodeJoinAccept(resp, session)
	if err != nil {
		return err
	}

	return nil
}

func SendUplink(data []uint8, session *Session) error {

	if regionSettings == nil {
		return ErrUndefinedRegionSettings
	}

	payload, err := session.GenMessage(0, []byte(data))
	if err != nil {
		return err
	}

	joinChannel := regionSettings.GetUplinkChannel()
	ActiveRadio.SetFrequency(joinChannel.Frequency)
	ActiveRadio.SetBandwidth(joinChannel.Bandwidth)
	ActiveRadio.SetCodingRate(joinChannel.CodingRate)
	ActiveRadio.SetSpreadingFactor(joinChannel.SpreadingFactor)
	ActiveRadio.SetCrc(true)
	ActiveRadio.SetIqMode(0) // IQ Standard
	ActiveRadio.Tx(payload, LORA_TX_TIMEOUT)
	if err != nil {
		return err
	}
	return nil
}

func ListenDownlink() error {
	return nil
}
