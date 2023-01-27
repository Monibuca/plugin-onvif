package onvif

import (
	. "m7s.live/engine/v4"
	"time"
)

type OnvifConfig struct {
	DiscoverInterval int
	Interfaces       []struct {
		InterfaceName string
		Username      string
		Password      string
	}
	Devices []struct {
		IP       string
		Username string
		Password string
	}
}

func (o *OnvifConfig) init() {
	preprocessAuth(authCfg)
	if o.DiscoverInterval == 0 {
		o.DiscoverInterval = 30
	}
	go func() {
		deviceList.discoveryDevice()
		deviceList.pullStream()
	}()
	t := time.NewTicker(time.Duration(o.DiscoverInterval) * time.Second)
	go func() {
		for range t.C {
			deviceList.discoveryDevice()
			deviceList.pullStream()
		}
	}()
}

func (o *OnvifConfig) OnEvent(event any) {
	switch event.(type) {
	case FirstConfig:
		o.init()
	}
}

var authCfg = &AuthConfig{
	Interfaces: make(map[string]deviceAuth),
	Devices:    make(map[string]deviceAuth),
}

var deviceList = &DeviceList{Data: make(map[string]map[string]*DeviceStatus)}

var conf = &OnvifConfig{}
var plugin = InstallPlugin(conf)

func preprocessAuth(c *AuthConfig) {
	for _, i := range conf.Interfaces {
		c.Interfaces[i.InterfaceName] = deviceAuth{
			Username: i.Username,
			Password: i.Password,
		}
	}
	for _, d := range conf.Devices {
		c.Devices[d.IP] = deviceAuth{
			Username: d.Username,
			Password: d.Password,
		}
	}
}
