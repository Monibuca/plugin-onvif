# onvif 拉流

仅实现了onvif设备发现，并通过rtsp拉流，暂未实现onvif协议其它功能，如需其它功能，请加群联系购买Pro版。

**注意** onvif 监听了udp 1024端口，使用了广播，可能需要在路由器或者电脑防火墙设置一下

# 预览onvif流
打开m7s `/preview` 界面，onvif 设备视频流如下图所示，格式为：`onvif/ip_port` 

![image](https://user-images.githubusercontent.com/1883728/215057255-ccaf0359-df7d-48a5-9336-754a4ca1b8f8.png)


配置如下：
```yaml
onvif:
  discoverinterval: 30 # 发现设备的间隔，单位秒，默认30秒，建议比rtsp插件的重连间隔大点
  interfaces: # 设备发现指定网卡，以及该网卡对应IP段的全局默认账号密码，支持多网卡
    - interfacename: WLAN  # 网卡名称 或者"以太网" "eth0"等，使用ipconfig 或者 ifconfig 查看网卡名称 
      username: admin # onvif 账号
      password: admin # onvif 密码
    - interfacename: WLAN 2 # 网卡2
      username: admin
      password: admin
  devices: # 可以给指定设备配置单独的密码
    - ip: 192.168.1.1
      username: admin
      password: '123'
    - ip: 192.168.1.2
      username: admin
      password: '456'
```
