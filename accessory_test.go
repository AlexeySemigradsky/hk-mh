package hkmh_test

import (
	"os"
	"testing"
	"time"

	"github.com/AlexeySemigradsky/hkmh"
	"github.com/brutella/hc/accessory"
)

var info = accessory.Info{
	Name: "Magic Home Controller",
}
var address = os.Getenv("DEVICE_ADDRESS")

func TestNewAccessory(t *testing.T) {
	_, err := hkmh.NewAccessory(info, address, 3*time.Second)
	if err != nil {
		t.Error(err)
	}
}
