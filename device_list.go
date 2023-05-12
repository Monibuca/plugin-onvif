package onvif

import (
	"m7s.live/engine/v4"
	"m7s.live/engine/v4/log"
	"reflect"
	"strings"
	"unsafe"

	lonvif "github.com/IOTechSystems/onvif"
	rtsp "m7s.live/plugin/rtsp/v4"
)

type DeviceList struct {
	Data map[string]map[string]*DeviceStatus
}

func changeDeviceParam(d *lonvif.Device, param lonvif.DeviceParams) {
	pointerVal := reflect.ValueOf(d)
	val := reflect.Indirect(pointerVal)
	member := val.FieldByName("params")
	ptrToY := unsafe.Pointer(member.UnsafeAddr())
	realPtrToY := (*lonvif.DeviceParams)(ptrToY)
	*realPtrToY = param
}

func (dl *DeviceList) discoveryDevice() {
	for _, i := range conf.Interfaces {
		deviceParams := WsDiscover(i.InterfaceName, authCfg)

		devsMap, ok := dl.Data[i.InterfaceName]
		if !ok {
			devsMap = make(map[string]*DeviceStatus)
			dl.Data[i.InterfaceName] = devsMap
		}

		for _, dParam := range deviceParams {
			// 如果已经存在，则不再添加
			if _, ok := devsMap[dParam.Xaddr]; ok {
				continue
			}
			var dev *lonvif.Device
			devStatus := &DeviceStatus{}
			dev, err := lonvif.NewDevice(dParam)
			if err == nil {
				devStatus.Device = dev
				devStatus.Status = StatusInitOk
			} else {
				devStatus.Status = StatusInitError
				devStatus.Device = &lonvif.Device{}
				changeDeviceParam(devStatus.Device, dParam) // todo 之前返回了 DeviceParam 是否有必要，不记得了
			}
			devsMap[dParam.Xaddr] = devStatus
		}
	}
}

func (dl *DeviceList) pullStream() {
	for _, devicesMap := range dl.Data {
		for _, d := range devicesMap {
			streamPath := strings.ReplaceAll(d.Device.GetDeviceParams().Xaddr, ".", "_")
			streamPath = "onvif/" + strings.ReplaceAll(streamPath, ":", "_")
			//避免重复拉流
			if engine.Streams.Has(streamPath) {
				continue
			}
			rtspUrl, err := GetStreamUri(d.Device)
			if err != nil {
				log.Info("[ONVIF] get stream err:", err)
				d.Status = StatusGetStreamUriError
				continue
			}
			d.Status = StatusGetStreamUriOk
			go func(targetURL string, streamPath string, d *DeviceStatus) {
				err = rtsp.RTSPPlugin.Pull(streamPath, targetURL, new(rtsp.RTSPPuller), 0)
				if err == nil {
					d.Status = StatusPullRtspOk
				} else {
					d.Status = StatusPullRtspError
				}
			}(rtspUrl, streamPath, d)
		}
	}
}
