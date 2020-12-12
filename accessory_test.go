package hkmh_test

import (
	"github.com/AlexeySemigradsky/hk-mh"
	"github.com/brutella/hc/accessory"
	"os"
	"testing"
)

var info = accessory.Info{
	Name: "Magic Home Controller",
}
var address = os.Getenv("DEVICE_ADDRESS")

func TestNewAccessory(t *testing.T) {
	_, err := hkmh.NewAccessory(info, address)
	if err != nil {
		t.Error(err)
	}
}
