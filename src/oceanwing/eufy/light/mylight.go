package light

import (
	"oceanwing/mqttclient"
)

type lightProd struct {
	devKEY    string
	devID     string
	pubTopicl string
	subTopicl string
	Incoming  chan []byte
}

// MqttServerPoint 用于模拟从服务器下发指令给灯泡
type MqttServerPoint struct {
	mqttclient.MqttClient
	lighters []*lightProd
}
