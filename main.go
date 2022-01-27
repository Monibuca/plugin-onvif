package onvif

import (
	"time"

	. "github.com/Monibuca/engine/v3"
)

var config struct {
	DiscoverInterval int `toml:"DiscoverInterval"`
	Interfaces       []struct {
		InterfaceName string `toml:"InterfaceName"`
		Username      string `toml:"Username"`
		Password      string `toml:"Password"`
	} `toml:"interfaces"`
	Devices []struct {
		IP       string `toml:"Ip"`
		Username string `toml:"Username"`
		Password string `toml:"Password"`
	} `toml:"devices"`
}

var authCfg = &AuthConfig{
	Interfaces: make(map[string]deviceAuth),
	Devices:    make(map[string]deviceAuth),
}

var deviceList = &DeviceList{Data: make(map[string]map[string]*DeviceStatus)}

func init() {
	pconfig := PluginConfig{
		Name:   "ONVIF",
		Config: &config,
	}
	pconfig.Install(runPlugin)
}

func runPlugin() {
	preprocessAuth(authCfg)
	if config.DiscoverInterval == 0 {
		config.DiscoverInterval = 30
	}
	t := time.NewTicker(time.Duration(config.DiscoverInterval) * time.Second)
	go func() {
		for range t.C {
			deviceList.discoveryDevice()
			deviceList.pullStream()
		}
	}()

}

func preprocessAuth(c *AuthConfig) {
	for _, i := range config.Interfaces {
		c.Interfaces[i.InterfaceName] = deviceAuth{
			Username: i.Username,
			Password: i.Password,
		}
	}
	for _, d := range config.Devices {
		c.Devices[d.IP] = deviceAuth{
			Username: d.Username,
			Password: d.Password,
		}
	}
}
