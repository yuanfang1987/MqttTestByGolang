package serverpoint

import (
	"oceanwingqa/common/mqttclient"
	"oceanwingqa/eufybackend/functional/device"
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
		prod := codeAndKey[0]
		devkey := codeAndKey[1]
		switch prod {
		case "T1012", "T1011":
			dev = device.NewLight(prod, devkey)
		case "T1013", "T1604":
			dev = device.NewLightWithColor(prod, devkey)
		case "T2103":
			dev = device.NewRobotCleaner(prod, devkey)
		}
		log.Debugf("Create a %s device, Key: %s", prod, devkey)
		dev.HandleSubscribeMessage()
		s.devices = append(s.devices, dev)
	}
}

// RunMqttService hh.
func (s *MqttServerPoint) RunMqttService(clientid, username, pwd, broker string, ca bool) {
	s.Clientid = clientid
	s.Username = username
	s.Pwd = pwd
	s.Broker = broker
	s.SubTopic = s.devices[0].GetSubTopic() // T1012, PUH_MESSAGE , "DEVICE/+/+/+"
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
	payload := message.Payload()
	for _, dev := range s.devices {
		if t == dev.GetSubTopic() || t == dev.GetPubTopic() {
			// 有事没事多开几个 goroutine ，确保消息传递快速畅通
			go dev.SendPayload(t, payload)
			//一条消息只可能匹配到一个设备，完了之后立刻结束 for 循环
			return
		}
	}
}

// PublishMsgToBroker 发布指令到Broker上，由Broker推送给Device
func (s *MqttServerPoint) PublishMsgToBroker() {
	if len(s.devices) == 0 {
		log.Error("No device found.")
		return
	}

	for _, dev := range s.devices {
		s.PubTopic = dev.GetPubTopic()
		payload := dev.BuildProtoBufMessage()
		if payload == nil || len(payload) == 0 {
			continue
		}
		s.MqttClient.PublishMessage(payload)
	}
}
