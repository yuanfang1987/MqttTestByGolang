package device

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// EufyDevice 所有 eufy device 的行为接口
type EufyDevice interface {
	HandleSubscribeMessage()
	LeaveModeTestResult()
	GetSubDeviceTopic() string
	GetSubServerTopic() string
	GetProductCode() string
	GetProductKey() string
	SendPayload(MQTT.Message)
}

type baseDevice struct {
	ProdCode       string
	DevKEY         string
	DevID          string //预留，不一定能用得到
	SubDeviceTopic string
	SubServerTopic string
	SubMessage     chan MQTT.Message
}

func (b *baseDevice) GetSubDeviceTopic() string {
	return b.SubDeviceTopic
}

func (b *baseDevice) GetSubServerTopic() string {
	return b.SubServerTopic
}

func (b *baseDevice) GetProductCode() string {
	return b.ProdCode
}

func (b *baseDevice) GetProductKey() string {
	return b.DevKEY
}

func (b *baseDevice) SendPayload(msg MQTT.Message) {
	b.SubMessage <- msg
}

func (b *baseDevice) LeaveModeTestResult() {

}
