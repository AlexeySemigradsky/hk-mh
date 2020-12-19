package hkmh

import (
	"github.com/brutella/hc/characteristic"
	"github.com/brutella/hc/service"
)

type Service struct {
	*service.Service
	On         *characteristic.On
	Brightness *characteristic.Brightness
	Saturation *characteristic.Saturation
	Hue        *characteristic.Hue
}

func NewService() *Service {
	s := Service{}
	s.Service = service.New(service.TypeLightbulb)

	s.On = characteristic.NewOn()
	s.AddCharacteristic(s.On.Characteristic)

	s.Brightness = characteristic.NewBrightness()
	s.AddCharacteristic(s.Brightness.Characteristic)

	s.Saturation = characteristic.NewSaturation()
	s.AddCharacteristic(s.Saturation.Characteristic)

	s.Hue = characteristic.NewHue()
	s.AddCharacteristic(s.Hue.Characteristic)

	return &s
}
