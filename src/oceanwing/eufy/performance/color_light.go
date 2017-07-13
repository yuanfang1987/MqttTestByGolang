package performance

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// 这是一个描述彩灯的结构体，支持 T1013, T1604
type colorLight struct {
	baseEufy
}

// NewColorLight hj.
func NewColorLight(clientid, username, pwd, broker, prodCode, devKey string, needCA bool) Eufydevice {
	c := &colorLight{}
	c.Clientid = clientid
	c.Username = username
	c.Pwd = pwd
	c.Broker = broker
	c.prod = prodCode
	c.PubTopic = "DEVICE/" + prodCode + "/" + devKey + "/PUH_MESSAGE"
	c.SubTopic = "DEVICE/" + prodCode + "/" + devKey + "/SUB_MESSAGE"
	c.NeedCA = needCA
	return c
}

// 实现 Eufydevice 接口
func (c *colorLight) RunMqttService() {
	c.SubHandler = func(c MQTT.Client, msg MQTT.Message) {} // do nothing.
	c.MqttClient.ConnectToBroker()
	c.outgoing()
}

// 实现 Eufydevice 接口
func (c *colorLight) SendHeartBeat() {
	c.msgToServer <- c.buildHeartBeatMsg()
}

func (c *colorLight) buildHeartBeatMsg() []byte {
	return nil
}
