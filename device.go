package onvif

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/beevik/etree"
	lonvif "github.com/liyanhui1998/go-onvif"
	"github.com/videonext/onvif/profiles/media"
)

//设备状态
const (
	StatusInitOk = iota
	StatusInitError
	StatusGetStreamUriOk
	StatusGetStreamUriError
	StatusPullRtspOk
	StatusPullRtspError
)

type deviceAuth struct {
	Username string
	Password string
}
type AuthConfig struct {
	Interfaces map[string]deviceAuth
	Devices    map[string]deviceAuth
}

type DeviceStatus struct {
	Device *lonvif.Device
	Status int
}

func WsDiscover(interfaceName string, config *AuthConfig) []lonvif.DeviceParams {
	/* Call an ws-discovery Probe Message to Discover NVT type Devices */
	devices := lonvif.SendProbe(interfaceName, nil, []string{"dn:" + lonvif.NVT.String()}, map[string]string{"dn": "http://www.lonvif.org/ver10/network/wsdl"})
	nvtDevices := make([]lonvif.DeviceParams, 0)

	for _, j := range devices {
		doc := etree.NewDocument()
		if err := doc.ReadFromString(j); err != nil {
			fmt.Println("[ONVIF] parse SendProbe error:", err.Error())
			continue
		}
		endpoints := doc.Root().FindElements("./Body/ProbeMatches/ProbeMatch/XAddrs")
		for _, xaddr := range endpoints {
			xaddr := strings.Split(strings.Split(xaddr.Text(), " ")[0], "/")[2]
			ip := strings.Split(xaddr, " ")[0]
			auth := getDeviceAuth(interfaceName, ip, config)
			nvtDevices = append(nvtDevices, lonvif.DeviceParams{Ipddr: ip, Username: auth.Username, Password: auth.Password})
		}
	}
	return nvtDevices
}

func GetStreamUri(dev *lonvif.Device) (string, error) {
	Response, err := dev.CallMethod(media.GetProfiles{})
	if err != nil {
		return "", err
	}
	resp, err := ioutil.ReadAll(Response.Body)
	if err != nil {
		return "", err
	}
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(resp); err != nil {
		return "", fmt.Errorf("error:%s", err.Error())
	}

	endpoints := doc.Root().FindElements("./Body/GetProfilesResponse/Profiles")
	if len(endpoints) == 0 {
		return "", fmt.Errorf("error:%s", "no Profiles")
	}
	//profile 是节点属性 <trt:Profiles token="PROFILE_374774454" fixed="true"/>
	profileToken := endpoints[0].SelectAttrValue("token", "")
	if profileToken == "" {
		return "", fmt.Errorf("error:%s", "profile token is empty")
	}
	Response, _ = dev.CallMethod(media.GetStreamUri{ProfileToken: media.ReferenceToken(profileToken)})
	resp, err = ioutil.ReadAll(Response.Body)
	if err != nil {
		return "", err
	}
	doc = etree.NewDocument()

	if err := doc.ReadFromBytes(resp); err != nil {
		return "", fmt.Errorf("error:%s", err.Error())
	}

	endpoints = doc.Root().FindElements("./Body/GetStreamUriResponse/MediaUri/Uri")
	if len(endpoints) == 0 {
		return "", fmt.Errorf("error:%s", "no media uri")
	}
	mediaUri := endpoints[0].Text()
	if !strings.Contains(mediaUri, "rtsp") {
		fmt.Println("mediaUri:", mediaUri)
		return "", fmt.Errorf("error:%s", "media uri is not rtsp")
	}
	if !strings.Contains(mediaUri, "@") && dev.Params.Username != "" {
		//如果返回的rtsp里没有账号密码，则自己拼接
		mediaUri = strings.Replace(mediaUri, "//", fmt.Sprintf("//%s:%s@", dev.Params.Username, dev.Params.Password), 1)
	}
	return mediaUri, nil
}

//获取设备的账号密码
func getDeviceAuth(interfaceName string, ip string, config *AuthConfig) deviceAuth {
	var auth deviceAuth
	if a, ok := config.Interfaces[interfaceName]; ok {
		auth = a
	}
	if a, ok := config.Devices[ip]; ok {
		auth = a
	}
	return auth
}
