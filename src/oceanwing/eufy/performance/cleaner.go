package performance

import (
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

// Cleaner struct.
type Cleaner struct {
	baseEufy
	robotACK []byte
	Cindex   int
}

func newRobotCleaner(clientid, username, pwd, broker, prodCode, devKey string, needCA bool) Eufydevice {
	o := &Cleaner{}
	o.Clientid = clientid
	o.Username = username
	o.Pwd = pwd
	o.Broker = broker
	o.prod = prodCode
	o.PubTopic = "DEVICE/T2103/" + devKey + "/PUH_MESSAGE"
	o.SubTopic = "DEVICE/T2103/" + devKey + "/SUB_MESSAGE"
	o.NeedCA = needCA
	o.msgToServer = make(chan []byte, 2)
	o.msgFromServer = make(chan []byte, 2)
	return o
}

// RunMqttService 实现 Eufydevice 接口
func (r *Cleaner) RunMqttService() {
	r.SubHandler = func(c MQTT.Client, msg MQTT.Message) {
		go r.distributeMsg(msg.Payload()) //把消息发给 channel
	}
	r.MqttClient.ConnectToBroker()
	r.outgoing()
	r.inComing()
}

// SendHeartBeat 实现 Eufydevice 接口
func (r *Cleaner) SendHeartBeat() {
	r.msgToServer <- r.getCommand()
	r.Cindex++
	if r.Cindex > 5 {
		r.Cindex = 0
	}
}

func (r *Cleaner) distributeMsg(payload []byte) {
	r.msgFromServer <- payload
}

func (r *Cleaner) inComing() {
	go func() {
		for {
			select {
			case msg := <-r.msgFromServer:
				r.handleInComingCMD(msg)
			}
		}
	}()
}

// GetCommand hh.
func (r *Cleaner) getCommand() []byte {
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
	}
	r.msgToServer <- r.robotACK
}
