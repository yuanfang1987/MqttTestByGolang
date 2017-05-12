package robot

import (
	"oceanwing/mqttclient"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// NewRobotCleaner create a new robot cleaner and return it.
func NewRobotCleaner() *Cleaner {
	return &Cleaner{
		Cindex:   0,
		Command1: []byte{0xA5, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x50, 0x00, 0x00, 0x01, 0x04, 0x26, 0xFF, 0xFF, 0x79, 0xFA},
		Command2: []byte{0xA5, 0x01, 0x01, 0x00, 0x00, 0x00, 0x01, 0x01, 0x01, 0x01, 0x5A, 0x00, 0x00, 0x00, 0x11, 0x39, 0xFF, 0xFF, 0xA8, 0xFA},
		Command3: []byte{0xA5, 0x02, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x02, 0x02, 0x46, 0x00, 0x00, 0x00, 0x20, 0x29, 0xFF, 0xFF, 0x94, 0xFA},
		Command4: []byte{0xA5, 0x03, 0x00, 0x00, 0x01, 0x00, 0x01, 0x01, 0x03, 0x00, 0x3C, 0x00, 0x00, 0x00, 0x19, 0x38, 0xFF, 0xFF, 0x94, 0xFA},
		Command5: []byte{0xA5, 0x04, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x32, 0x00, 0x00, 0x00, 0x22, 0x36, 0xFF, 0xFF, 0x8E, 0xFA},
		Command6: []byte{0xA5, 0x05, 0x01, 0x00, 0x00, 0x00, 0x01, 0x01, 0x01, 0x02, 0x28, 0x00, 0x00, 0x00, 0x24, 0x52, 0xFF, 0xFF, 0xA7, 0xFA},
	}
}

// Cleaner struct.
type Cleaner struct {
	mqttclient.MqttClient
	Cindex   int
	Command1 []byte
	Command2 []byte
	Command3 []byte
	Command4 []byte
	Command5 []byte
	Command6 []byte
}

//RunRobotCleanerMqttService connect to mqtt and subscribeto topic. so cool.
func (r *Cleaner) RunRobotCleanerMqttService(clientid, username, pwd, broker, devKey string, needCA bool) {
	r.Clientid = clientid
	r.Username = username
	r.Pwd = pwd
	r.Broker = broker
	r.PubTopic = "DEVICE/T2103/" + devKey + "/PUH_MESSAGE"
	r.SubTopic = "DEVICE/T2103/" + devKey + "/SUB_MESSAGE"
	r.NeedCA = needCA
	// msgHandler
	r.SubHandler = func(c MQTT.Client, msg MQTT.Message) {
		r.ReceivedPayload = msg.Payload()
	}
	// connect to broker and subscribe to a topic.
	r.MqttClient.ConnectToBroker()
}

// GetCommand hh.
func (r *Cleaner) getCommand() []byte {
	//i := RandInt64(0, 6)
	switch r.Cindex {
	case 0:
		return r.Command1
	case 1:
		return r.Command2
	case 2:
		return r.Command3
	case 3:
		return r.Command4
	case 4:
		return r.Command5
	case 5:
		return r.Command6
	default:
		return r.Command1
	}
}

// SendRobotCleanerHeartBeat send a heart beat to broker.
func (r *Cleaner) SendRobotCleanerHeartBeat() {
	r.MqttClient.PublishMessage(r.getCommand())
	// update index.
	r.Cindex++
	if r.Cindex > 5 {
		r.Cindex = 0
	}
}

// SendRequestLocalCodeCommand hh.
// func (r *Cleaner) SendRequestLocalCodeCommand() {
// 	r.MqttClient.PublishMessage(r.BuildEufyDeviceMsg(2))
// }

// BuildEufyDeviceMsg build a message that send request to server.
// func (r *Cleaner) BuildEufyDeviceMsg(cmdType int) []byte {
// 	var cmd *PB.Dev2Server_CmdType

// 	if cmdType == 1 {
// 		cmd = PB.Dev2Server_CmdType_CLEAR_ALL_CONNECT.Enum()
// 	} else if cmdType == 2 {
// 		cmd = PB.Dev2Server_CmdType_REQUEST_LOCAL_CODE.Enum()
// 	} else {
// 		cmd = PB.Dev2Server_CmdType_REPORT_DEV_INFO.Enum()
// 	}

// 	d := &PB.Dev2ServerMessage{
// 		Type: cmd,
// 	}

// 	data, err := proto.Marshal(d)
// 	if err == nil {
// 		return data
// 	}
// 	return nil

// }

// DecodeEufyServerMsg decodes the message that come from the eufy server.
// func DecodeEufyServerMsg(payload []byte) string {
// 	d := &PB.Server2DevMessage{}
// 	err := proto.Unmarshal(payload, d)
// 	if err != nil {
// 		return ""
// 	}

// 	typ := d.GetType()
// 	if typ == PB.Server2Dev_CmdType_RESPONSE_LOCAL_CODE {
// 		return d.GetLocalCode()
// 	}
// 	return ""
// }
