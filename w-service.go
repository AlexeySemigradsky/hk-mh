package hkmh

import (
	"github.com/brutella/hc/characteristic"
	"github.com/brutella/hc/service"
)

type WService struct {
	*service.Service
	Name       *characteristic.Name
	On         *characteristic.On
	Brightness *characteristic.Brightness
}

func NewWService() *WService {
	s := WService{}
	s.Service = service.New(service.TypeLightbulb)

	s.Name = characteristic.NewName()
	s.Name.SetValue("W")
	s.AddCharacteristic(s.Name.Characteristic)

	s.On = characteristic.NewOn()
	s.AddCharacteristic(s.On.Characteristic)

	s.Brightness = characteristic.NewBrightness()
	s.AddCharacteristic(s.Brightness.Characteristic)

	return &s
}
