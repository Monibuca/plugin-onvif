# onvif 拉流

仅实现了onvif设备发现，并通过rtsp拉流，暂未实现onvif协议其它功能。

**注意** 依赖rtsp插件，请先安装rtsp插件。

配置如下：
```toml
[ONVIF]
DiscoverInterval = 30 # 发现设备的间隔，单位秒，默认30秒
# 设备发现指定网卡，以及该网卡对应IP段的全局默认账号密码
[[ONVIF.interfaces]]
InterfaceName = "WLAN" #或者 以太网  eth0 等
Username= "admin"
Password= "admin"
# 如果有多个网卡配置多个即可
# [[ONVIF.interfaces]]
# InterfaceName = "eth1"
# Username= "admin"
# Password= "admin2"

# # 如果设备账号密码和全局不一致，单独设置指定IP地址的设备账号密码
# [[ONVIF.devices]]
# Ip = "192.168.1.1"
# Username= "admin"
# Password= "123"
# [[ONVIF.devices]]
# Ip = "192.168.1.2"
# Username= "admin"
# Password= "456"
```