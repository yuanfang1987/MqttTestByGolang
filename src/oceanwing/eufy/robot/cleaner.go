package robot

import (
	"oceanwing/mqttclient"

	log "github.com/cihub/seelog"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var (
	cmdPause     = []byte{0xA5, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x5C, 0x00, 0x00, 0x01, 0x10, 0x09, 0xFF, 0xFF, 0x74, 0xFA}
	cmdSpot      = []byte{0xA5, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x5D, 0x00, 0x00, 0x00, 0x10, 0x08, 0xFF, 0xFF, 0x74, 0xFA}
	cmdAuto      = []byte{0xA5, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x5D, 0x00, 0x00, 0x00, 0x10, 0x09, 0xFF, 0xFF, 0x76, 0xFA}
	cmdCharge    = []byte{0xA5, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x5D, 0x01, 0x00, 0x01, 0x10, 0x0F, 0xFF, 0xFF, 0x7F, 0xFA}
	cmdEdge      = []byte{0xA5, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x5D, 0x00, 0x00, 0x00, 0x10, 0x00, 0xFF, 0xFF, 0x6F, 0xFA}
	cmdSmallRoom = []byte{0xA5, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x5C, 0x00, 0x00, 0x00, 0x10, 0x09, 0xFF, 0xFF, 0x78, 0xFA}
)

// NewRobotCleaner create a new robot cleaner and return it.
func NewRobotCleaner() *Cleaner {
	return &Cleaner{
		Cindex: 0,
	}
}

// Cleaner struct.
type Cleaner struct {
	mqttclient.MqttClient
	robotACK    []byte
	CmdToServer chan []byte
	Cindex      int
}

//RunRobotCleanerMqttService connect to mqtt and subscribeto topic.
func (r *Cleaner) RunRobotCleanerMqttService(clientid, username, pwd, broker, devKey string, needCA bool) {
	r.Clientid = clientid
	r.Username = username
	r.Pwd = pwd
	r.Broker = broker
	r.PubTopic = "DEVICE/T2103/" + devKey + "/PUH_MESSAGE"
	r.SubTopic = "DEVICE/T2103/" + devKey + "/SUB_MESSAGE"
	r.NeedCA = needCA
	// 2017.05.16 added
	r.CmdToServer = make(chan []byte)
	r.robotACK = cmdPause
	// msgHandler
	r.SubHandler = func(c MQTT.Client, msg MQTT.Message) {
		log.Debugf("new msg coming base on topic: %s", msg.Topic())
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
		return cmdPause
	case 1:
		return cmdSpot
	case 2:
		return cmdAuto
	case 3:
		return cmdCharge
	case 4:
		return cmdEdge
	case 5:
		return cmdSmallRoom
	default:
		return cmdAuto
	}
}

// SendRobotCleanerHeartBeat send a heart beat to broker.
func (r *Cleaner) SendRobotCleanerHeartBeat() {
	//r.MqttClient.PublishMessage(r.getCommand())
	r.CmdToServer <- r.robotACK
	log.Debug("new payload send to channel: CmdToServer")
	// update index.
	// r.Cindex++
	// if r.Cindex > 5 {
	// 	r.Cindex = 0
	// }
}

// publish消息的统一出口，不管是常规心跳，还是即时响应，都是经过这里发出
func (r *Cleaner) outgoingMsg() {
	go func() {
		for {
			select {
			case payload := <-r.CmdToServer:
				log.Debug("receive new payload from channel: CmdToServer")
				r.MqttClient.PublishMessage(payload)
			}
		}
	}()
	log.Debug("run function: outgoingMsg")
}

func (r *Cleaner) handleInComingCMD(pl []byte) {
	if len(pl) != 8 {
		log.Debugf("Oops! receive invalid format message, the message length is: %d", len(pl))
		return
	}
	log.Debugf("index 4: %d, index 5: %d", pl[4], pl[5])
	if pl[4] == 225 {
		switch pl[5] {
		case 0:
			// pause
			log.Debug("receive command: Pause")
			r.robotACK = cmdPause
		case 1:
			// spot
			log.Debug("receive command: Cleaning On Spot Mode")
			r.robotACK = cmdSpot
		case 2:
			// auto
			log.Debug("receive command: Cleaning On Auto Mode")
			r.robotACK = cmdAuto
		case 3:
			// back to charge
			log.Debug("receive command: Back to charge")
			r.robotACK = cmdCharge
		case 4:
			// edge
			log.Debug("receive command: Cleaning On Edge Mode")
			r.robotACK = cmdEdge
		case 5:
			// small room
			log.Debug("receive command: Cleaning On Small Room Mode")
			r.robotACK = cmdSmallRoom
		}
	} else if pl[4] == 232 {
		switch pl[5] {
		case 0:
			// daily
			log.Debug("receive command: Set Speed to Daily")
			if r.robotACK[8] != 0 {
				r.robotACK[8] = 0x00
				r.robotACK[18] = r.robotACK[18] - 1
			}
		case 1:
			// strong
			log.Debug("receive command: Set Speed to Strong")
			if r.robotACK[8] != 1 {
				r.robotACK[8] = 0x01
				r.robotACK[18] = r.robotACK[18] + 1
			}
		}
	} else if pl[4] == 236 {
		log.Debug("receive command: Find Me Alert")
		return
		// do nothing, just let the robot make some noise. a, I am so tired.
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

// Command1: []byte{0xA5, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x50, 0x00, 0x00, 0x01, 0x04, 0x26, 0xFF, 0xFF, 0x79, 0xFA},
// Command2: []byte{0xA5, 0x01, 0x01, 0x00, 0x00, 0x00, 0x01, 0x01, 0x01, 0x01, 0x5A, 0x00, 0x00, 0x00, 0x11, 0x39, 0xFF, 0xFF, 0xA8, 0xFA},
// Command3: []byte{0xA5, 0x02, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x02, 0x02, 0x46, 0x00, 0x00, 0x00, 0x20, 0x29, 0xFF, 0xFF, 0x94, 0xFA},
// Command4: []byte{0xA5, 0x03, 0x00, 0x00, 0x01, 0x00, 0x01, 0x01, 0x03, 0x00, 0x3C, 0x00, 0x00, 0x00, 0x19, 0x38, 0xFF, 0xFF, 0x94, 0xFA},
// Command5: []byte{0xA5, 0x04, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x32, 0x00, 0x00, 0x00, 0x22, 0x36, 0xFF, 0xFF, 0x8E, 0xFA},
// Command6: []byte{0xA5, 0x05, 0x01, 0x00, 0x00, 0x00, 0x01, 0x01, 0x01, 0x02, 0x28, 0x00, 0x00, 0x00, 0x24, 0x52, 0xFF, 0xFF, 0xA7, 0xFA},
