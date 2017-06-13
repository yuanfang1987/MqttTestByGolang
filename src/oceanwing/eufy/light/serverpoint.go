package light

import (
	"oceanwing/commontool"
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
	lighters []*lightProd
}

// NewMqttServerPoint 新建一个 MqttServerPoint实例
func NewMqttServerPoint() *MqttServerPoint {
	if servPointInstance == nil {
		servPointInstance = &MqttServerPoint{}
	}
	return servPointInstance
}

// SetupRunningLights 设置需要控制的灯泡数量
func (s *MqttServerPoint) SetupRunningLights(keys []string) {
	for _, key := range keys {
		codeAndKey := strings.Split(key, ":")
		light := &lightProd{
			prodCode:  codeAndKey[0],
			devKEY:    codeAndKey[1],
			mode:      0, // 默认模式
			status:    1, //默认开灯
			pubTopicl: "DEVICE/T1012/" + codeAndKey[1] + "/SUB_MESSAGE",
			subTopicl: "DEVICE/T1012/" + codeAndKey[1] + "/PUH_MESSAGE",
			Incoming:  make(chan []byte),
		}
		log.Debugf("Set up a device successfully: %s", codeAndKey[1])
		light.handleIncomingMsg()
		s.lighters = append(s.lighters, light)
	}
}

// RunMqttService hh.
func (s *MqttServerPoint) RunMqttService(clientid, username, pwd, broker string, ca bool) {
	s.Clientid = clientid
	s.Username = username
	s.Pwd = pwd
	s.Broker = broker
	s.SubTopic = "DEVICE/T1012/+/PUH_MESSAGE"
	s.NeedCA = ca
	s.SubHandler = func(c MQTT.Client, msg MQTT.Message) {
		go s.distributeMsg(msg)
	}
	// connect to broker
	s.MqttClient.ConnectToBroker()
}

// 分发订阅到的消息给对应的light去处理
func (s *MqttServerPoint) distributeMsg(message MQTT.Message) {
	t := message.Topic()
	payload := message.Payload()
	for _, light := range s.lighters {
		if t == light.subTopicl {
			log.Debugf("send incoming message to device: %s, message id: %d", light.devKEY, message.MessageID())
			light.Incoming <- payload
			return
		}
	}
}

// PublishMsgToLight 发指令给灯泡
func (s *MqttServerPoint) PublishMsgToLight() {
	if len(s.lighters) == 0 {
		log.Error("No lights found.")
		return
	}

	for _, light := range s.lighters {
		s.PubTopic = light.pubTopicl
		brightness := uint32(commontool.RandInt64(1, 100))
		color := uint32(commontool.RandInt64(0, 100))
		payload := light.buildSetLightDataMsg(brightness, color)
		// payload := light.buildSetAwayModeMsg(15, 10, 15, 20, 1, 1, true, true)
		s.MqttClient.PublishMessage(payload)
	}
}