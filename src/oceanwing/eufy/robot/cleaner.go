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
	robotACK    []byte
	CmdToServer chan []byte
	Cindex      int
	Command1    []byte
	Command2    []byte
	Command3    []byte
	Command4    []byte
	Command5    []byte
	Command6    []byte
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
	// 2017.05.16 added
	r.robotACK = []byte{0xA5, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x5C, 0x00, 0x00, 0x01, 0x10, 0x09, 0xFF, 0xFF, 0x74, 0xFA}
	// msgHandler
	r.SubHandler = func(c MQTT.Client, msg MQTT.Message) {
		go r.handleInComingCMD(msg.Payload())
	}
	// connect to broker and subscribe to a topic.
	r.MqttClient.ConnectToBroker()
	// 2017.05.16 added. 作为一个publish消息的统一的出口
	r.outgoingMsg()
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
	//r.MqttClient.PublishMessage(r.getCommand())
	r.CmdToServer <- r.robotACK
	// update index.
	// r.Cindex++
	// if r.Cindex > 5 {
	// 	r.Cindex = 0
	// }
}

func (r *Cleaner) outgoingMsg() {
	go func() {
		for {
			select {
			case payload := <-r.CmdToServer:
				r.MqttClient.PublishMessage(payload)
			}
		}
	}()
}

func (r *Cleaner) handleInComingCMD(pl []byte) {
	if len(pl) != 8 {
		return
	}
	if pl[4] == 225 {
		switch pl[5] {
		case 0:
			// pause
			r.robotACK = []byte{0xA5, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x5C, 0x00, 0x00, 0x01, 0x10, 0x09, 0xFF, 0xFF, 0x74, 0xFA}
		case 1:
			// spot
			r.robotACK = []byte{0xA5, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x5D, 0x00, 0x00, 0x00, 0x10, 0x08, 0xFF, 0xFF, 0x74, 0xFA}
		case 2:
			// auto
			r.robotACK = []byte{0xA5, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x5D, 0x00, 0x00, 0x00, 0x10, 0x09, 0xFF, 0xFF, 0x76, 0xFA}
		case 3:
			// back to charge
			r.robotACK = []byte{0xA5, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x5D, 0x01, 0x00, 0x01, 0x10, 0x0F, 0xFF, 0xFF, 0x7F, 0xFA}
		case 4:
			// edge
			r.robotACK = []byte{0xA5, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x5D, 0x00, 0x00, 0x00, 0x10, 0x00, 0xFF, 0xFF, 0x6F, 0xFA}
		case 5:
			// small room
			r.robotACK = []byte{0xA5, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x5C, 0x00, 0x00, 0x00, 0x10, 0x09, 0xFF, 0xFF, 0x78, 0xFA}
		}
	} else if pl[4] == 232 {
		switch pl[5] {
		case 0:
			// daily
			if r.robotACK[8] != 0 {
				r.robotACK[8] = 0x00
				r.robotACK[18] = r.robotACK[18] - 1
			}
		case 1:
			// strong
			if r.robotACK[8] != 1 {
				r.robotACK[8] = 0x01
				r.robotACK[18] = r.robotACK[18] + 1
			}
		}
	} else if pl[4] == 236 {
		return
		// do nothing, just let the robot make some noise.
	}
	r.CmdToServer <- r.robotACK
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
