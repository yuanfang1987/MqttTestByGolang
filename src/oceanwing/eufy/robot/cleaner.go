package robot

import (
	PB "oceanwing/eufy/commonpb"

	"oceanwing/mqttclient"

	log "github.com/cihub/seelog"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/protobuf/proto"
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

//RunRobotCleanerMqttService connect to mqtt and subscribeto topic.
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
func (r *Cleaner) SendRequestLocalCodeCommand() {
	r.MqttClient.PublishMessage(r.BuildEufyDeviceMsg(2))
}

// BuildEufyDeviceMsg build a message that send request to server.
func (r *Cleaner) BuildEufyDeviceMsg(cmdType int) []byte {
	var cmd *PB.Dev2Server_CmdType

	if cmdType == 1 {
		cmd = PB.Dev2Server_CmdType_CLEAR_ALL_CONNECT.Enum()
	} else if cmdType == 2 {
		cmd = PB.Dev2Server_CmdType_REQUEST_LOCAL_CODE.Enum()
	} else {
		cmd = PB.Dev2Server_CmdType_REPORT_DEV_INFO.Enum()
	}

	d := &PB.Dev2ServerMessage{
		Type: cmd,
	}

	data, err := proto.Marshal(d)
	if err == nil {
		return data
	}
	return nil

}

// DecodeEufyServerMsg decodes the message that come from the eufy server.
func DecodeEufyServerMsg(payload []byte) string {
	d := &PB.Server2DevMessage{}
	err := proto.Unmarshal(payload, d)
	if err != nil {
		return ""
	}

	typ := d.GetType()
	if typ == PB.Server2Dev_CmdType_RESPONSE_LOCAL_CODE {
		return d.GetLocalCode()
	}
	return ""
}

/* ============================== Functions displaying below are used for functional testing. =============================== */

// EufyServer 用于模拟从Server端Publish消息给机器人
type EufyServer struct {
	mqttclient.MqttClient
	caseNum      int
	littleRobots []*littleRobot
	appUser      *AppUser
}

// NewEufyServer 创建一个空的 EufyServer 对象
func NewEufyServer() *EufyServer {
	return &EufyServer{}
}

// SetupRunningRobots 用于接收并处理实体机器人的心跳
func (r *EufyServer) SetupRunningRobots(devkeys []string) {
	for _, v := range devkeys {
		o := &littleRobot{
			devKEY:       v,
			pubTopicl:    "DEVICE/T2103/" + v + "/SUB_MESSAGE",
			subTopicl:    "DEVICE/T2103/" + v + "/PUH_MESSAGE",
			Incoming:     make(chan []byte),
			charging:     false,
			returnCharge: false,
		}
		// debug
		o.handleIncomingMsg()
		r.littleRobots = append(r.littleRobots, o)
	}
}

// SetAppUser 虚拟App用户
func (r *EufyServer) SetAppUser(cid, csecret, email, pwd string) {
	r.appUser = NewAppUser()
	r.appUser.Login(cid, csecret, email, pwd)
}

// AddRunningRobot 增加一个跟虚拟App用户绑定的robot.
func (r *EufyServer) AddRunningRobot(devkey, devid string) {
	o := &littleRobot{
		devKEY:       devkey,
		devID:        devid,
		subTopicl:    "DEVICE/T2103/" + devkey + "/PUH_MESSAGE",
		Incoming:     make(chan []byte),
		charging:     false,
		returnCharge: false,
	}
	o.handleIncomingMsg()
	r.littleRobots = append(r.littleRobots, o)
}

// RunMqtt 连接上Broker，并订阅一个通配符主题，接收所有机器人的心跳
func (r *EufyServer) RunMqtt(clientid, username, pwd, broker string, ca bool) {
	r.Clientid = clientid
	r.Username = username
	r.Pwd = pwd
	r.Broker = broker
	r.SubTopic = "DEVICE/T2103/+/PUH_MESSAGE"
	r.NeedCA = ca

	r.SubHandler = func(c MQTT.Client, msg MQTT.Message) {
		go r.distributeMsg(msg.Topic(), msg.Payload())
	}
	// connect.
	r.MqttClient.ConnectToBroker()
}

// distributeMsg 根据主题把消息分发给对应的littleRobot去处理
func (r *EufyServer) distributeMsg(t string, payload []byte) {
	for _, robot := range r.littleRobots {
		if t == robot.subTopicl {
			// log.Infof("--new heart beat message coming with topic: %s, from device: %s", t, robot.devKEY)
			robot.Incoming <- payload
			return
		}
	}
}

// PublishMsgToAllRobot 把消息依次Publish给所有机器人, 或通过HTTP发到Server
func (r *EufyServer) PublishMsgToAllRobot() {
	if len(r.littleRobots) != 0 {
		for _, robot := range r.littleRobots {
			if robot.pubTopicl != "" {
				r.PubTopic = robot.pubTopicl
			}
			if !robot.charging {
				if robot.devID != "" {
					r.appUser.SendCmdToServer(robot.devID)
				} else {
					r.MqttClient.PublishMessage(getCommandToDevice(r.caseNum, robot))
					log.Infof("发送指令给机器: %s, 指令内容: %s", robot.devKEY, robot.testPurpose)
					r.caseNum++
					if r.caseNum > 28 {
						// reset to 0
						r.caseNum = 0
					}
				}
			} else if robot.returnCharge {
				if robot.devID != "" {
					r.appUser.SetIndexForReturnHome()
					r.appUser.SendCmdToServer(robot.devID)
				} else {
					r.MqttClient.PublishMessage(getCommandToDevice(3, robot))
					log.Infof("机器 [%s] 正处于低电量, 发指令叫它回家充电", robot.devKEY)
				}
			} else {
				log.Infof("机器 [%s] 正在充电中......, 现在不发任何指令.", robot.devKEY)
			}
		}
	}
}

// littleRobot 用于处理实体机器人返回的心跳
type littleRobot struct {
	devKEY         string
	devID          string
	pubTopicl      string
	subTopicl      string
	Incoming       chan []byte
	charging       bool
	returnCharge   bool
	expResultIndex byte
	expResultValue byte
	isCmdSent      bool
	testPurpose    string
}

func (robot *littleRobot) handleIncomingMsg() {
	go func() {
		for {
			select {
			case heartBeatInfo := <-robot.Incoming:
				// heartBeatInfo 是一个[]byte字节数组，长度为20
				// heartBeatInfo[0]: 固定0xA5，十进制为165
				// heartBeatInfo[1]: WorkMode, 0x00=暂停   0x01＝定点  0x02=自动 0x03=返回充电    0x04=沿边    0x05=精扫,（0xf0=休眠 WIFI模块主动添加）
				// heartBeatInfo[2]: OnOff_Direction_Forward,  0=机器停止   1= 机器向前
				// heartBeatInfo[3]  OnOff_Direction_Backward, 0=机器停止   1=机器后退
				// heartBeatInfo[4]: OnOff_Direction_Lef,  0=机器停止   1= 机器左转
				// heartBeatInfo[5]: OnOff_Direction_Right,0=机器停止   1=机器右转
				// heartBeatInfo[6]: DownAdjust: BIT0 = 0下视距离远   BIT0=1下视距离近 BIT1=1-水箱插入 BIT1=0-水箱移除
				// heartBeatInfo[7]: SideAdjust: 0=侧视距离远  1=侧视距离近
				// heartBeatInfo[8]: CleaningSpeed: 0=日常    1=强力  2=地毯  3=静音
				// heartBeatInfo[9]: RoomMode: 0=大房间 1=小房间  2=中房间
				// heartBeatInfo[10]:BatteryCapacity  0~0X64
				// heartBeatInfo[11]:ChargerStatus    0=没充电   1= 正在充电
				// heartBeatInfo[12]:ErrorCode        0~6
				// heartBeatInfo[13]:Stop_Cleaning    1=停止   0=工作
				// heartBeatInfo[14]:当前时间：小时
				// heartBeatInfo[15]:当前时间：分钟
				// heartBeatInfo[16]:闹钟时间：小时
				// heartBeatInfo[17]:闹钟时间：分钟
				// heartBeatInfo[18]:校验和
				// heartBeatInfo[19]:MCU发送数据结束信号0XFA
				if len(heartBeatInfo) == 20 {
					log.Infof("=== 有新的心跳消息从机器 [%s] 到来 ===", robot.devKEY)
					log.Infof("模式: %d,  device key: %s", heartBeatInfo[1], robot.devKEY)
					log.Infof("电量: %d,  device key: %s", heartBeatInfo[10], robot.devKEY)
					log.Infof("向前: %d,  device key: %s", heartBeatInfo[2], robot.devKEY)
					log.Infof("向后: %d,  device key: %s", heartBeatInfo[3], robot.devKEY)
					log.Infof("向左: %d,  device key: %s", heartBeatInfo[4], robot.devKEY)
					log.Infof("向右: %d,  device key: %s", heartBeatInfo[5], robot.devKEY)
					log.Infof("速度: %d,  device key: %s", heartBeatInfo[8], robot.devKEY)
					log.Infof("房间: %d,  device key: %s", heartBeatInfo[9], robot.devKEY)
					log.Infof("充电: %d,  device key: %s", heartBeatInfo[11], robot.devKEY)
					log.Infof("停止: %d,  device key: %s", heartBeatInfo[13], robot.devKEY)
					//如果ErrorCode不为0，则机器内部可能出错
					if heartBeatInfo[12] != 0 {
						log.Errorf("警告! ErrorCode: %d, device key: %s", heartBeatInfo[12], robot.devKEY)
					}
					//水箱状态
					log.Infof("水箱状态:  %d", heartBeatInfo[6])
					// 如果电量在20%以下且workmode不等于3,则让机器人返回充电
					if heartBeatInfo[10] <= 20 {
						robot.charging = true
						if heartBeatInfo[1] != 3 {
							robot.returnCharge = true
						} else {
							robot.returnCharge = false
						}
						// 如果电量已超过90%且仍在充电状态，则把两个变量设置为false
					} else if heartBeatInfo[10] >= 90 && robot.charging {
						robot.charging = false
						robot.returnCharge = false
					}
				}
				if len(heartBeatInfo) == 20 && robot.isCmdSent && robot.expResultIndex != 99 && robot.devID == "" {
					if robot.testPurpose != "" {
						log.Infof("Test Case: %s, device key: %s", robot.testPurpose, robot.devKEY)
					}
					if heartBeatInfo[robot.expResultIndex] != robot.expResultValue {
						log.Infof("Expected value is: %d, Actual value is: %d, device key: %s", robot.expResultValue, heartBeatInfo[robot.expResultIndex], robot.devKEY)
					} else {
						log.Infof("test case passed. device key: %s", robot.devKEY)
					}
					robot.isCmdSent = false
					robot.expResultIndex = 99
					robot.testPurpose = ""
				}
			}
		}
	}()
}

func getCommandToDevice(index int, rb *littleRobot) []byte {
	rb.isCmdSent = true
	switch index {
	// ----------------------- 工作模式: 0xE1 ----------------------------------- good.
	case 0:
		// 暂停 0x00, ok
		rb.expResultIndex = 1
		rb.expResultValue = 0
		rb.testPurpose = "设置工作模式为： [0x00], 暂停"
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE1, 0x00, 0xE1, 0xFA}
	case 1:
		// 定点 0x01, ok
		rb.expResultIndex = 1
		rb.expResultValue = 1
		rb.testPurpose = "设置工作模式为： [0x01], 定点"
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE1, 0x01, 0xE2, 0xFA}
	case 2:
		// 自动 0x02, ok
		rb.expResultIndex = 1
		rb.expResultValue = 2
		rb.testPurpose = "设置工作模式为： [0x02], 自动"
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE1, 0x02, 0xE3, 0xFA}
	case 3:
		// 返回充电 0x03, ok
		rb.expResultIndex = 1
		rb.expResultValue = 3
		rb.testPurpose = "设置工作模式为： [0x03], 返回充电"
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE1, 0x03, 0xE4, 0xFA}
	case 4:
		// 沿边 0x04, ok
		rb.expResultIndex = 1
		rb.expResultValue = 4
		rb.testPurpose = "设置工作模式为： [0x04], 沿边"
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE1, 0x04, 0xE5, 0xFA}
	case 5:
		//精扫 0x05
		rb.expResultIndex = 1
		rb.expResultValue = 5
		rb.testPurpose = "设置工作模式为： [0x05], 精扫"
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE1, 0x05, 0xE6, 0xFA}
	// ------------------------ 向前: 0xE2 ----------------------------------- good.
	case 6:
		//向前 0x01
		rb.expResultIndex = 2
		rb.expResultValue = 1
		rb.testPurpose = "执行方向操作： [0x01], 向前"
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE2, 0x01, 0xE3, 0xFA}
	case 7:
		//停止向前 0x00
		rb.expResultIndex = 2
		rb.expResultValue = 0
		rb.testPurpose = "执行方向操作： [0x00], 停止向前"
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE2, 0x00, 0xE2, 0xFA}
	// -------------------------向后: 0xE3 ----------------------------------- good.
	case 8:
		//向后
		rb.expResultIndex = 3
		rb.expResultValue = 1
		rb.testPurpose = "执行方向操作： [0x01], 向后"
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE3, 0x01, 0xE4, 0xFA}
	case 9:
		//停止向后
		rb.expResultIndex = 3
		rb.expResultValue = 0
		rb.testPurpose = "执行方向操作： [0x00], 停止向后"
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE3, 0x00, 0xE3, 0xFA}
	// -------------------------向左: 0xE4 ----------------------------------- good.
	case 10:
		//向左
		rb.expResultIndex = 4
		rb.expResultValue = 1
		rb.testPurpose = "执行方向操作： [0x01], 向左"
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE4, 0x01, 0xE5, 0xFA}
	case 11:
		//停止向左
		rb.expResultIndex = 4
		rb.expResultValue = 0
		rb.testPurpose = "执行方向操作： [0x00], 停止向左"
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE4, 0x00, 0xE4, 0xFA}
	// -------------------------向右: 0xE5 ----------------------------------- good.
	case 12:
		//向右
		rb.expResultIndex = 5
		rb.expResultValue = 1
		rb.testPurpose = "执行方向操作： [0x01], 向右"
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE5, 0x01, 0xE6, 0xFA}
	case 13:
		//停止向右
		rb.expResultIndex = 5
		rb.expResultValue = 0
		rb.testPurpose = "执行方向操作： [0x00], 停止向右"
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE5, 0x00, 0xE5, 0xFA}
	// -------------------------下视距离: 0xE6 --------------------------------
	case 14:
		//下视距离近:
		rb.testPurpose = "下视距离近"
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE6, 0x01, 0xE7, 0xFA}
	case 15:
		//下视距离远:
		rb.testPurpose = "下视距离远"
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE6, 0x00, 0xE6, 0xFA}
	// -------------------------侧视距离: 0xE7 --------------------------------
	case 16:
		//侧视距离近:
		rb.testPurpose = "侧视距离近"
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE7, 0x01, 0xE8, 0xFA}
	case 17:
		//侧视距离远:
		rb.testPurpose = "侧视距离远"
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE7, 0x00, 0xE7, 0xFA}
	// ------------------------- Clean Speed: 0xE8 ----------------------------  good.
	case 18:
		//日常 0x00
		rb.expResultIndex = 8
		rb.expResultValue = 0
		rb.testPurpose = "设置速度为： [0x00], 日常"
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE8, 0x00, 0xE8, 0xFA}
	case 19:
		//强力 0x01
		rb.expResultIndex = 8
		rb.expResultValue = 1
		rb.testPurpose = "设置速度为： [0x01], 强力"
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE8, 0x01, 0xE9, 0xFA}
	case 20:
		//地毯 0x02
		rb.expResultIndex = 8
		rb.expResultValue = 2
		rb.testPurpose = "设置速度为： [0x02], 地毯"
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE8, 0x02, 0xEA, 0xFA}
	case 21:
		//静音 0x03
		rb.expResultIndex = 8
		rb.expResultValue = 3
		rb.testPurpose = "设置速度为： [0x03], 静音"
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE8, 0x03, 0xEB, 0xFA}
	// ------------------------- Room Mode: 0xE9 ------------------------------ good
	case 22:
		//大房间 0x00
		rb.expResultIndex = 9
		rb.expResultValue = 0
		rb.testPurpose = "设置房间为： [0x00], 大房间"
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE9, 0x00, 0xE9, 0xFA}
	case 23:
		//小房间 0x01
		rb.expResultIndex = 9
		rb.expResultValue = 1
		rb.testPurpose = "设置房间为： [0x01], 小房间"
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE9, 0x01, 0xEA, 0xFA}
	case 24:
		//中房间 0x02
		rb.expResultIndex = 9
		rb.expResultValue = 2
		rb.testPurpose = "设置房间为： [0x02], 中房间"
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE9, 0x02, 0xEB, 0xFA}
	// ------------------------- Stop Cleaning: 0xEA -------------------------- good.
	case 25:
		//stop 0x01
		rb.expResultIndex = 13
		rb.expResultValue = 1
		rb.testPurpose = "设置是否停止： [0x01], 是"
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xEA, 0x01, 0xEB, 0xFA}
	case 26:
		// not stop 0x00
		rb.expResultIndex = 13
		rb.expResultValue = 0
		rb.testPurpose = "设置是否停止： [0x00], 否"
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xEA, 0x00, 0xEA, 0xFA}
	// ------------------------- Find Me Alert: 0xEC -------------------------- good.
	case 27:
		// turn on alert: 0x01
		rb.testPurpose = "打开 [Find Me Alert]"
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xEC, 0x01, 0xED, 0xFA}
	case 28:
		// turn off alert: 0x00
		rb.testPurpose = "关闭 [Find Me Alert]"
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xEC, 0x00, 0xEC, 0xFA}
	default:
		return nil
	}
}
