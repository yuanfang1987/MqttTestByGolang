package device

import (
	log "github.com/cihub/seelog"
)

// LightWithColor 是对产品 T1013、T1604 的描述
type LightWithColor struct {
	baseDevice
	stopCtrlFunc chan struct{}
}

// NewLightWithColor create a new color light instance.
func NewLightWithColor(prodCode, devKey, devid string) EufyDevice {
	o := &LightWithColor{}
	o.ProdCode = prodCode
	o.DevKEY = devKey
	o.DevID = devid
	o.PubTopicl = "DEVICE/" + prodCode + "/" + devKey + "/SUB_MESSAGE"
	o.SubTopicl = "DEVICE/" + prodCode + "/" + devKey + "/PUH_MESSAGE"
	o.DeviceMsg = make(chan []byte)
	o.ServerMsg = make(chan []byte)
	o.stopCtrlFunc = make(chan struct{})
	log.Infof("Create a color Light, product code: %s, device key: %s, device id: %s", prodCode, devKey, devid)
	return o
}
