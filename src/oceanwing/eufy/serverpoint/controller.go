package serverpoint

import (
	"oceanwing/eufy/device"
	"oceanwing/mqttclient"
	"strings"

	log "github.com/cihub/seelog"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var (
	servPointInstance *MqttServerPoint
)

// MqttServerPoint 用于模拟从服务器下发指令给灯泡
type MqttServerPoint struct {
	mqttclient.MqttClient
	devices []device.EufyDevice
}

// NewMqttServerPoint 新建一个 MqttServerPoint实例
func NewMqttServerPoint() *MqttServerPoint {
	if servPointInstance == nil {
		servPointInstance = &MqttServerPoint{}
	}
	return servPointInstance
}

// SetupRunningDevices 设置需要控制的device数量
func (s *MqttServerPoint) SetupRunningDevices(keys []string) {
	var dev device.EufyDevice
	for _, key := range keys {
		codeAndKey := strings.Split(key, ":")
		switch codeAndKey[0] {
		case "T1012", "T1011":
			dev = device.NewLight(codeAndKey[0], codeAndKey[1])
		}
		log.Debugf("Create a %s device, Key: %s", codeAndKey[0], codeAndKey[1])
		dev.HandleSubscribeMessage()
		s.devices = append(s.devices, dev)
	}
}

// RunMqttService hh..
func (s *MqttServerPoint) RunMqttService(clientid, username, pwd, broker string, ca bool) {
	s.Clientid = clientid
	s.Username = username
	s.Pwd = pwd
	s.Broker = broker
	s.SubTopic = "DEVICE/+/+/+" // T1012, PUH_MESSAGE
	s.NeedCA = ca
	s.SubHandler = func(c MQTT.Client, msg MQTT.Message) {
		go s.distributeMsg(msg)
	}
	// connect to broker
	s.MqttClient.ConnectToBroker()
}

// 分发订阅到的消息给对应的 device 去处理
func (s *MqttServerPoint) distributeMsg(message MQTT.Message) {
	t := message.Topic()
	log.Debugf("收到订阅的消息，主题为：%s", t)
	for _, dev := range s.devices {
		if t == dev.GetSubDeviceTopic() || t == dev.GetSubServerTopic() {
			dev.SendPayload(message)
			return
		}
	}
}

// HappyEnding oh yeah, finish testing.
func HappyEnding() {
	for _, dev := range servPointInstance.devices {
		dev.LeaveModeTestResult()
	}
}
