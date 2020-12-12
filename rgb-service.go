package hkmh

import (
	"github.com/brutella/hc/characteristic"
	"github.com/brutella/hc/service"
)

type RGBService struct {
	*service.Service
	Name       *characteristic.Name
	On         *characteristic.On
	Brightness *characteristic.Brightness
	Saturation *characteristic.Saturation
	Hue        *characteristic.Hue
}

func NewRGBService() *RGBService {
	s := RGBService{}
	s.Service = service.New(service.TypeLightbulb)

	s.Name = characteristic.NewName()
	s.Name.SetValue("RGB")
	s.AddCharacteristic(s.Name.Characteristic)

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
