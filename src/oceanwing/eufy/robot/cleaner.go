package robot

import (
	"oceanwing/mqttclient"

	"time"

	log "github.com/cihub/seelog"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

const (
	pause       byte = 0
	spot        byte = 1
	auto        byte = 2
	charge      byte = 3
	edge        byte = 4
	smallRoom   byte = 5
	strongSpeed byte = 1
	dailySpeed  byte = 0
)

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
		log.Infof("robot\t%s\t发出指令总数\t%d", rb.devKEY, rb.totalCMD)
		log.Infof("robot\t%s\t收到心跳总数\t%d", rb.devKEY, rb.heartBeatCount)
		log.Infof("robot\t%s\t执行暂停指令\t成功\t%d\t失败\t%d\t总计\t%d", rb.devKEY, rb.pausePassed, rb.pauseFailed, rb.pausePassed+rb.pauseFailed)
		log.Infof("robot\t%s\t执行定点指令\t成功\t%d\t失败\t%d\t总计\t%d", rb.devKEY, rb.spotPassed, rb.spotFailed, rb.spotPassed+rb.spotFailed)
		log.Infof("robot\t%s\t执行自动指令\t成功\t%d\t失败\t%d\t总计\t%d", rb.devKEY, rb.autoPassed, rb.autoFailed, rb.autoPassed+rb.autoFailed)
		log.Infof("robot\t%s\t执行返回充电指令\t成功\t%d\t失败\t%d\t总计\t%d", rb.devKEY, rb.chargePassed, rb.chargeFailed, rb.chargePassed+rb.chargeFailed)
		log.Infof("robot\t%s\t执行沿边指令\t成功\t%d\t失败\t%d\t总计\t%d", rb.devKEY, rb.edgePassed, rb.edgeFailed, rb.edgePassed+rb.edgeFailed)
		log.Infof("robot\t%s\t执行精扫(小房间)指令\t成功\t%d\t失败\t%d\t总计\t%d", rb.devKEY, rb.smallRoomPassed, rb.smallRoomFailed, rb.smallRoomPassed+rb.smallRoomFailed)
		log.Infof("robot\t%s\t执行设置日常速度指令\t成功\t%d\t失败\t%d\t总计\t%d", rb.devKEY, rb.speedDailyPassed, rb.speedDailyFailed, rb.speedDailyPassed+rb.speedDailyFailed)
		log.Infof("robot\t%s\t执行设置强力速度指令\t成功\t%d\t失败\t%d\t总计\t%d", rb.devKEY, rb.speedStrongPassed, rb.speedStrongFailed, rb.speedStrongPassed+rb.speedStrongFailed)
		log.Infof("robot\t%s\t执行打开FindMe次数\t%d", rb.devKEY, rb.turnOnFindMe)
		log.Infof("robot\t%s\t执行关闭FindMe次数\t%d", rb.devKEY, rb.turnOffFindMe)
	}
}

// littleRobot 用于处理实体机器人返回的心跳
type littleRobot struct {
	devKEY            string
	devID             string
	pubTopicl         string
	subTopicl         string
	testType          int
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
	testType2Mode     byte
	responseMode      byte
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

					robot.responseMode = heartBeatInfo[1]

					// 如果 testType = 2, 则需额外判断
					if robot.testType == 2 && robot.testType2Mode != 99 {
						if robot.testType2Mode != heartBeatInfo[1] {
							log.Infof("robot [%s] current running mode should be %d, but actual is %d", robot.devKEY, robot.testType2Mode, heartBeatInfo[1])
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

// 用于判断机器人自动执行的后续动作
func (robot *littleRobot) assertNextActionResp(expValue byte) {
	for i := 0; i < 3; i++ {
		if robot.responseMode == expValue {
			log.Infof("robot next action response is %d, Passed.", expValue)
			return
		}
		time.Sleep(20 * time.Second)
	}
	log.Infof("robot next action response is %d, but %d is expected, Failed.", robot.responseMode, expValue)
}

func getCommandToDevice(index int, rb *littleRobot) []byte {
	rb.isCmdSent = true
	//指令总数累加1
	rb.totalCMD++
	switch index {
	case 1:
		// 定点 0x01, ok
		rb.expResultIndex = 1
		rb.expResultValue = spot
		rb.testPurpose = "设置工作模式为：定点"
		rb.testCaseNum = 1
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE1, 0x01, 0xE2, 0xFA}
	case 2:
		//强力 0x01
		rb.expResultIndex = 8
		rb.expResultValue = strongSpeed
		rb.testPurpose = "设置速度为：强力"
		rb.testCaseNum = 2
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE8, 0x01, 0xE9, 0xFA}
	case 3:
		// 自动 0x02, ok
		rb.expResultIndex = 1
		rb.expResultValue = auto
		rb.testPurpose = "设置工作模式为：自动"
		rb.testCaseNum = 3
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE1, 0x02, 0xE3, 0xFA}
	case 4:
		// 返回充电 0x03, ok
		rb.expResultIndex = 1
		rb.expResultValue = charge
		rb.testPurpose = "设置工作模式为：返回充电"
		rb.testCaseNum = 4
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE1, 0x03, 0xE4, 0xFA}
	case 5:
		// 沿边 0x04, ok
		rb.expResultIndex = 1
		rb.expResultValue = edge
		rb.testPurpose = "设置工作模式为：沿边"
		rb.testCaseNum = 5
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE1, 0x04, 0xE5, 0xFA}
	case 6:
		//小房间 0x05
		rb.expResultIndex = 1
		rb.expResultValue = smallRoom
		rb.testPurpose = "设置工作模式为：小房间"
		rb.testCaseNum = 6
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE1, 0x05, 0xE6, 0xFA}
	case 7:
		//日常 0x00
		rb.expResultIndex = 8
		rb.expResultValue = dailySpeed
		rb.testPurpose = "设置速度为：日常"
		rb.testCaseNum = 7
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE8, 0x00, 0xE8, 0xFA}
	case 8:
		// 暂停 0x00, ok
		rb.expResultIndex = 1
		rb.expResultValue = pause
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

// =================================================================================

// RunTestType2 hh.
func (r *EufyServer) RunTestType2() {
	myTimer := time.NewTimer(10 * time.Millisecond)
	<-myTimer.C

	var robot *littleRobot
	var reSetFlag bool
	var nextActResp byte
	// 如果有两个机器，则用第2个来跑， 否则用第1个
	if len(r.littleRobots) > 2 {
		robot = r.littleRobots[1]
	} else if len(r.littleRobots) == 1 {
		robot = r.littleRobots[0]
	}
	robot.testType = 2
	robot.testType2Mode = 99
	// 1, spot; 5, edge; 6, small room; 3, auto
	index := 1
	for {
		if !robot.charging {
			log.Infof("Current command index is %d", index)
			r.MqttClient.PublishMessage(getCommandToDevice(index, robot))
			switch index {
			case 1:
				robot.testType2Mode = spot
				reSetFlag = myTimer.Reset(2 * time.Minute)
				log.Infof("reset timer to 2 minutes ---> %t", reSetFlag)
				index = 5
				nextActResp = pause
			case 5:
				robot.testType2Mode = edge
				reSetFlag = myTimer.Reset(20 * time.Minute)
				log.Infof("reset timer to 20 minutes ---> %t", reSetFlag)
				index = 6
				nextActResp = charge
			case 6:
				robot.testType2Mode = smallRoom
				reSetFlag = myTimer.Reset(30 * time.Minute)
				log.Infof("reset timer to 30 minutes ---> %t", reSetFlag)
				index = 3
				nextActResp = charge
			case 3:
				robot.testType2Mode = auto
				reSetFlag = myTimer.Reset(100 * time.Minute)
				log.Infof("reset timer to 100 minutes ---> %t", reSetFlag)
				index = 1
				nextActResp = charge
			}
			// 等待一个指定的时间, spot 等待2分钟；edge等待20分钟；small room等待30分钟；auto等待100分钟
			if reSetFlag {
				<-myTimer.C
			}
			// reset
			robot.testType2Mode = 99
			// 等待时间结束后，判断机器人是原地暂停，还是返回充电
			robot.assertNextActionResp(nextActResp)
		} else if robot.returnCharge {
			// 返回充电
			r.MqttClient.PublishMessage(getCommandToDevice(4, robot))
		} else {
			// 防止太过于频繁的写入日志，此处等待30秒
			reSetFlag = myTimer.Reset(30 * time.Second)
			if reSetFlag {
				<-myTimer.C
			}
			log.Infof("机器 [%s] 正在充电中......, 现在不发任何指令.", robot.devKEY)
		}

	}
}
