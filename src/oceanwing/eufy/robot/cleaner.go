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
var eufyServerInstance *EufyServer

// EufyServer 用于模拟从Server端Publish消息给机器人
type EufyServer struct {
	mqttclient.MqttClient
	caseNum      int
	littleRobots []*littleRobot
	appUser      *AppUser
}

// NewEufyServer 创建一个空的 EufyServer 对象
func NewEufyServer() *EufyServer {
	eufyServerInstance = &EufyServer{}
	return eufyServerInstance
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
	r.caseNum = 1
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
			//如果hangOn不为0，则表示机器人对前次发的指令响应的结果不正确，需等待下次心跳继续验证，不要发新指令过去
			if robot.hangOn != 0 {
				continue
			}
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
					if r.caseNum > 10 {
						// reset to 1
						r.caseNum = 1
					}
				}
			} else if robot.returnCharge {
				if robot.devID != "" {
					r.appUser.SetIndexForReturnHome()
					r.appUser.SendCmdToServer(robot.devID)
				} else {
					r.MqttClient.PublishMessage(getCommandToDevice(4, robot))
					log.Infof("机器 [%s] 正处于低电量, 发指令叫它回家充电", robot.devKEY)
				}
			} else {
				log.Infof("机器 [%s] 正在充电中......, 现在不发任何指令.", robot.devKEY)
			}
		}
	}
}

func (r *EufyServer) showTestResult() {
	for _, rb := range r.littleRobots {
		log.Infof("=========== Summary Result For Robot: %s ===========", rb.devKEY)
		log.Infof("robot: %s 发出指令总数: %d", rb.devKEY, rb.totalCMD)
		log.Infof("robot: %s 收到心跳总数: %d", rb.devKEY, rb.heartBeatCount)
		log.Infof("robot: %s 执行暂停指令，成功: %d, 失败: %d, 总计: %d", rb.devKEY, rb.pausePassed, rb.pauseFailed, rb.pausePassed+rb.pauseFailed)
		log.Infof("robot: %s 执行定点指令，成功: %d, 失败：%d, 总计: %d", rb.devKEY, rb.spotPassed, rb.spotFailed, rb.spotPassed+rb.spotFailed)
		log.Infof("robot: %s 执行自动指令，成功: %d, 失败：%d, 总计: %d", rb.devKEY, rb.autoPassed, rb.autoFailed, rb.autoPassed+rb.autoFailed)
		log.Infof("robot: %s 执行返回充电指令，成功: %d, 失败：%d, 总计: %d", rb.devKEY, rb.chargePassed, rb.chargeFailed, rb.chargePassed+rb.chargeFailed)
		log.Infof("robot: %s 执行沿边指令，成功: %d, 失败：%d, 总计: %d", rb.devKEY, rb.edgePassed, rb.edgeFailed, rb.edgePassed+rb.edgeFailed)
		log.Infof("robot: %s 执行精扫(小房间)指令，成功: %d, 失败：%d, 总计: %d", rb.devKEY, rb.smallRoomPassed, rb.smallRoomFailed, rb.smallRoomPassed+rb.smallRoomFailed)
		log.Infof("robot: %s 执行设置日常速度指令，成功: %d, 失败：%d, 总计: %d", rb.devKEY, rb.speedDailyPassed, rb.speedDailyFailed, rb.speedDailyPassed+rb.speedDailyFailed)
		log.Infof("robot: %s 执行设置强力速度指令，成功: %d, 失败：%d, 总计: %d", rb.devKEY, rb.speedStrongPassed, rb.speedStrongFailed, rb.speedStrongPassed+rb.speedStrongFailed)
		log.Infof("robot: %s 执行打开FindMe次数：%d", rb.devKEY, rb.turnOnFindMe)
		log.Infof("robot: %s 执行关闭FindMe次数：%d", rb.devKEY, rb.turnOffFindMe)
	}
}

// littleRobot 用于处理实体机器人返回的心跳
type littleRobot struct {
	devKEY            string
	devID             string
	pubTopicl         string
	subTopicl         string
	Incoming          chan []byte
	charging          bool
	returnCharge      bool
	expResultIndex    byte
	expResultValue    byte
	isCmdSent         bool
	testPurpose       string
	testCaseNum       int
	heartBeatCount    int
	hangOn            int
	totalCMD          int
	pausePassed       int
	pauseFailed       int
	spotPassed        int
	spotFailed        int
	autoPassed        int
	autoFailed        int
	chargePassed      int
	chargeFailed      int
	edgePassed        int
	edgeFailed        int
	smallRoomPassed   int
	smallRoomFailed   int
	speedStrongPassed int
	speedStrongFailed int
	speedDailyPassed  int
	speedDailyFailed  int
	turnOnFindMe      int
	turnOffFindMe     int
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
					// 累加心跳次数
					robot.heartBeatCount++

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

					// 如果电量在20%以下且workmode不等于3(充电中)，也不等于240(休眠),则让机器人返回充电
					if heartBeatInfo[10] <= 20 && heartBeatInfo[1] != 240 {
						robot.charging = true
						if heartBeatInfo[1] != 3 {
							robot.returnCharge = true
						} else {
							robot.returnCharge = false
						}
						// 如果电量已超过80%且仍在充电状态，则把两个变量设置为false, robot.charging
					} else if heartBeatInfo[10] >= 80 && heartBeatInfo[11] == 1 {
						robot.charging = false
						robot.returnCharge = false
					}

					//判断结果
					if robot.isCmdSent && robot.expResultIndex != 99 && robot.devID == "" {
						if heartBeatInfo[robot.expResultIndex] != robot.expResultValue {
							//如果前次发完指令后，得到的心跳验证不通过，则等待下一次心跳来继续验证，最多给3次机会
							robot.hangOn++
							if robot.hangOn == 3 {
								log.Infof("TestCase Failed --> %s, Expected value is: %d, Actual value is: %d, device key: %s",
									robot.testPurpose, robot.expResultValue, heartBeatInfo[robot.expResultIndex], robot.devKEY)
								robot.hangOn = 0
								robot.recordTestResult(false)
							}
						} else {
							log.Infof("TestCase Passed --> %s, device key: %s", robot.testPurpose, robot.devKEY)
							robot.hangOn = 0
							robot.recordTestResult(true)
						}

						//如果hangOn为0，则表示本次测试完毕，可以重置相关测试条件
						if robot.hangOn == 0 {
							robot.isCmdSent = false
							robot.expResultIndex = 99
							robot.testPurpose = ""
						}
					}
				}
			}
		}
	}()
}

func (robot *littleRobot) recordTestResult(passed bool) {
	switch robot.testCaseNum {
	case 1:
		if passed {
			robot.spotPassed++
		} else {
			robot.spotFailed++
		}
	case 2:
		if passed {
			robot.speedStrongPassed++
		} else {
			robot.speedStrongFailed++
		}
	case 3:
		if passed {
			robot.autoPassed++
		} else {
			robot.autoFailed++
		}
	case 4:
		if passed {
			robot.chargePassed++
		} else {
			robot.chargeFailed++
		}
	case 5:
		if passed {
			robot.edgePassed++
		} else {
			robot.edgeFailed++
		}
	case 6:
		if passed {
			robot.smallRoomPassed++
		} else {
			robot.smallRoomFailed++
		}
	case 7:
		if passed {
			robot.speedDailyPassed++
		} else {
			robot.speedDailyFailed++
		}
	case 8:
		if passed {
			robot.pausePassed++
		} else {
			robot.pauseFailed++
		}
	}
}

func getCommandToDevice(index int, rb *littleRobot) []byte {
	rb.isCmdSent = true
	//指令总数累加1
	rb.totalCMD++
	switch index {
	case 1:
		// 定点 0x01, ok
		rb.expResultIndex = 1
		rb.expResultValue = 1
		rb.testPurpose = "设置工作模式为：定点"
		rb.testCaseNum = 1
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE1, 0x01, 0xE2, 0xFA}
	case 2:
		//强力 0x01
		rb.expResultIndex = 8
		rb.expResultValue = 1
		rb.testPurpose = "设置速度为：强力"
		rb.testCaseNum = 2
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE8, 0x01, 0xE9, 0xFA}
	case 3:
		// 自动 0x02, ok
		rb.expResultIndex = 1
		rb.expResultValue = 2
		rb.testPurpose = "设置工作模式为：自动"
		rb.testCaseNum = 3
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE1, 0x02, 0xE3, 0xFA}
	case 4:
		// 返回充电 0x03, ok
		rb.expResultIndex = 1
		rb.expResultValue = 3
		rb.testPurpose = "设置工作模式为：返回充电"
		rb.testCaseNum = 4
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE1, 0x03, 0xE4, 0xFA}
	case 5:
		// 沿边 0x04, ok
		rb.expResultIndex = 1
		rb.expResultValue = 4
		rb.testPurpose = "设置工作模式为：沿边"
		rb.testCaseNum = 5
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE1, 0x04, 0xE5, 0xFA}
	case 6:
		//精扫 0x05
		rb.expResultIndex = 1
		rb.expResultValue = 5
		rb.testPurpose = "设置工作模式为：精扫"
		rb.testCaseNum = 6
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE1, 0x05, 0xE6, 0xFA}
	case 7:
		//日常 0x00
		rb.expResultIndex = 8
		rb.expResultValue = 0
		rb.testPurpose = "设置速度为：日常"
		rb.testCaseNum = 7
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE8, 0x00, 0xE8, 0xFA}
	case 8:
		// 暂停 0x00, ok
		rb.expResultIndex = 1
		rb.expResultValue = 0
		rb.testPurpose = "设置工作模式为：暂停"
		rb.testCaseNum = 8
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE1, 0x00, 0xE1, 0xFA}
	// ------------------------- Find Me Alert: 0xEC -------------------------- good.
	case 9:
		// turn on alert: 0x01
		rb.testPurpose = "打开 [Find Me Alert]"
		rb.turnOnFindMe++
		log.Infof("执行: %s, device key: %s", rb.testPurpose, rb.devKEY)
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xEC, 0x01, 0xED, 0xFA}
	case 10:
		// turn off alert: 0x00
		rb.testPurpose = "关闭 [Find Me Alert]"
		rb.turnOffFindMe++
		log.Infof("执行: %s, device key: %s", rb.testPurpose, rb.devKEY)
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xEC, 0x00, 0xEC, 0xFA}
	default:
		return nil
	}
}

// ShowSummaryResult hh.
func ShowSummaryResult() {
	log.Info("测试结束")
	eufyServerInstance.showTestResult()
}
