package hkmh

import (
	"log"
	"net"
	"time"

	"github.com/AlexeySemigradsky/clr"
	"github.com/AlexeySemigradsky/mh"
	"github.com/brutella/hc/accessory"
)

type Config struct {
	ID      uint64
	IP      net.IP
	Name    string
	Timeout time.Duration
}

type MHAccessory struct {
	Accessory  *accessory.Accessory
	config     Config
	rgb        *Service
	controller *mh.Controller
	timer      *time.Timer
}

func NewAccessory(config Config) *MHAccessory {
	c := mh.NewController(mh.Config{
		IP:      config.IP,
		Timeout: config.Timeout,
	})

	info := accessory.Info{
		ID:   config.ID,
		Name: config.Name,
	}

	a := MHAccessory{}
	a.Accessory = accessory.New(info, accessory.TypeLightbulb)
	a.config = config
	a.rgb = NewService()
	a.controller = c
	a.timer = time.AfterFunc(config.Timeout, a.pullState)

	a.Accessory.AddService(a.rgb.Service)

	a.pullState()
	a.rgb.On.OnValueRemoteUpdate(func(_ bool) { a.pushState() })
	a.rgb.Brightness.OnValueRemoteUpdate(func(_ int) { a.pushState() })
	a.rgb.Saturation.OnValueRemoteUpdate(func(_ float64) { a.pushState() })
	a.rgb.Hue.OnValueRemoteUpdate(func(_ float64) { a.pushState() })

	return &a
}

func (a *MHAccessory) pullState() {
	a.timer.Stop()

	power, err := a.controller.GetPower()
	if err != nil {
		log.Println(err)
		return
	}
	a.rgb.On.SetValue(power)

	rgbw, err := a.controller.GetRGBW()
	if err != nil {
		log.Println(err)
		return
	}

	h, s, b := rgbwToHsb(rgbw.Red, rgbw.Green, rgbw.Blue, rgbw.White)
	a.rgb.Hue.SetValue(h)
	a.rgb.Saturation.SetValue(s)
	a.rgb.Brightness.SetValue(b)

	a.timer.Reset(5 * time.Second)
}

func (a *MHAccessory) pushState() {
	a.timer.Stop()

	power := a.rgb.On.GetValue()

	err := a.controller.SetPower(power)
	if err != nil {
		log.Println(err)
		return
	}
	if !power {
		return
	}

	r, g, b, w := hsbToRgbw(
		a.rgb.Hue.GetValue(),
		a.rgb.Saturation.GetValue(),
		a.rgb.Brightness.GetValue(),
	)

	rgbw := &mh.RGBW{Red: r, Green: g, Blue: b, White: w}
	err = a.controller.SetRGBW(rgbw)
	if err != nil {
		log.Println(err)
		return
	}

	a.timer.Reset(5 * time.Second)
}

func hsbToRgbw(hue, saturation float64, brightness int) (red, green, blue, white float64) {
	red, green, blue = clr.HSVToRGB(hue, 100, float64(brightness))
	white = ((100 - saturation) / 100) * 255 * (float64(brightness) / 100)
	return red, green, blue, white
}

func rgbwToHsb(red, green, blue, white float64) (float64, float64, int) {
	hue, _, brightness := clr.RGBToHSV(red, green, blue)
	saturation := 100 - ((white * 100 * 100) / (brightness * 255))
	return hue, saturation, int(brightness)
}
