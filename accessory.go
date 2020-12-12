package hkmh

import (
	"github.com/AlexeySemigradsky/mh"
	"github.com/brutella/hc/accessory"
	"github.com/lucasb-eyer/go-colorful"
	"go.uber.org/zap"
)

type MHAccessory struct {
	*accessory.Accessory
	RGB        *RGBService
	W          *WService
	Controller *mh.Controller
	logger     *zap.SugaredLogger
}

func NewAccessory(info accessory.Info, address string) (*MHAccessory, error) {
	logger, err := NewLogger(zap.InfoLevel)
	if err != nil {
		return nil, err
	}

	c := mh.NewController(address)

	a := MHAccessory{}
	a.Accessory = accessory.New(info, accessory.TypeLightbulb)
	a.RGB = NewRGBService()
	a.W = NewWService()
	a.Controller = c
	a.logger = logger

	a.AddService(a.RGB.Service)
	a.AddService(a.W.Service)

	a.RGB.On.SetValue(a.getPower())
	a.RGB.On.OnValueRemoteUpdate(a.setPower)

	a.RGB.Brightness.SetValue(a.getBrightness())
	a.RGB.Brightness.OnValueRemoteUpdate(func(_ int) { a.setRGBW() })

	a.RGB.Saturation.SetValue(a.getSaturation())
	a.RGB.Saturation.OnValueRemoteUpdate(func(_ float64) { a.setRGBW() })

	a.RGB.Hue.SetValue(a.getHue())
	a.RGB.Hue.OnValueRemoteUpdate(func(_ float64) { a.setRGBW() })

	a.W.On.SetValue(a.getPower())
	a.W.On.OnValueRemoteUpdate(a.setPower)

	a.W.Brightness.SetValue(a.getWhite())
	a.W.Brightness.OnValueRemoteUpdate(func(_ int) { a.setRGBW() })

	return &a, nil
}

func (a *MHAccessory) getPower() bool {
	power, err := a.Controller.GetPower()
	if err != nil {
		a.logger.Errorf("getPower => err %v", err)
		return false
	}
	a.logger.Debugf("getPower => %t", power)
	return power
}

func (a *MHAccessory) setPower(_ bool) {
	power := a.RGB.On.GetValue() || a.W.On.GetValue()
	a.logger.Debugf("setPower(%t)", power)
	err := a.Controller.SetPower(power)
	if err != nil {
		a.logger.Errorf("setPower(%t) => err %v", power, err)
	}
}

func (a *MHAccessory) getHue() float64 {
	rgbw, err := a.Controller.GetRGBW()
	if err != nil {
		a.logger.Errorf("getGue() => err %v", err)
		return 0
	}
	color := colorful.Color{
		R: float64(rgbw.Red) / 255,
		G: float64(rgbw.Green) / 255,
		B: float64(rgbw.Blue) / 255,
	}
	h, _, _ := color.Hsv()
	a.logger.Debugf("getHue() => %f", h)
	return h
}

func (a *MHAccessory) getSaturation() float64 {
	rgbw, err := a.Controller.GetRGBW()
	if err != nil {
		a.logger.Errorf("getSaturation() => err %v", err)
		return 0
	}
	color := colorful.Color{
		R: float64(rgbw.Red) / 255,
		G: float64(rgbw.Green) / 255,
		B: float64(rgbw.Blue) / 255,
	}
	_, s, _ := color.Hsv()
	a.logger.Debugf("getSaturation() => %f", s*100)
	return s * 100
}

func (a *MHAccessory) getBrightness() int {
	rgbw, err := a.Controller.GetRGBW()
	if err != nil {
		a.logger.Errorf("getBrightness() => err %v", err)
		return 0
	}
	color := colorful.Color{
		R: float64(rgbw.Red) / 255,
		G: float64(rgbw.Green) / 255,
		B: float64(rgbw.Blue) / 255,
	}
	_, _, b := color.Hsv()
	a.logger.Debugf("getBrightness() => %d", int(b*100))
	return int(b * 100)
}

func (a *MHAccessory) getWhite() int {
	rgbw, err := a.Controller.GetRGBW()
	if err != nil {
		a.logger.Errorf("getWhite() => err %v", err)
		return 0
	}
	w := float64(rgbw.White) / 255
	a.logger.Debugf("getWhite() => %d", int(w*100))
	return int(w * 100)
}

func (a *MHAccessory) setRGBW() {
	power := a.RGB.On.GetValue() && a.W.On.GetValue()
	err := a.Controller.SetPower(power)
	if err != nil {
		a.logger.Errorf("SetPower(%t) => err %v", power, err)
		return
	}
	if !power {
		return
	}

	color := colorful.Hsv(
		a.RGB.Hue.GetValue(),
		a.RGB.Saturation.GetValue()/100,
		float64(a.RGB.Brightness.GetValue())/100,
	)
	r, g, b := color.RGB255()
	w := uint8(a.W.Brightness.GetValue())
	rgbw := &mh.RGBW{
		Red:   r,
		Blue:  b,
		Green: g,
		White: w,
	}
	err = a.Controller.SetRGBW(rgbw)
	if err != nil {
		a.logger.Errorf("SetRGBW(%v) => err %v", rgbw, err)
		return
	}
	a.logger.Debugf("setRGBW(%v)", rgbw)
}
