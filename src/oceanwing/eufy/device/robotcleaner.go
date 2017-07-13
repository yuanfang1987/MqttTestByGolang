package device

import (
	"fmt"
	"oceanwing/eufy/result"
	"strconv"

	log "github.com/cihub/seelog"
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

// RobotCleaner 是一个描述扫地机器人的 struct.
type RobotCleaner struct {
	baseDevice
	charging       bool
	returnCharge   bool
	expResultIndex byte
	expResultValue byte
	caseNum        int
	checkPoint     string
}

// NewRobotCleaner 新建一个robotcleaner的实例
func NewRobotCleaner(prodCode, devKey string) EufyDevice {
	o := &RobotCleaner{}
	o.ProdCode = prodCode
	o.DevKEY = devKey
	o.PubTopicl = "DEVICE/T2103/" + devKey + "/SUB_MESSAGE"
	o.SubTopicl = "DEVICE/T2103/" + devKey + "/PUH_MESSAGE"
	o.DeviceMsg = make(chan []byte)
	o.ServerMsg = make(chan []byte)
	o.caseNum = 1
	return o
}

// HandleSubscribeMessage 实现 EufyDevice 接口.
func (robot *RobotCleaner) HandleSubscribeMessage() {
	go func() {
		for {
			select {
			case heartBeatInfo := <-robot.DeviceMsg:
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
					// 累加心跳数量
					robot.DecodeHeartBeatMsgQuantity++

					log.Infof("=== 有新的心跳消息从机器 [%s] 到来 ===", robot.DevKEY)
					log.Infof("模式: %d,  device key: %s", heartBeatInfo[1], robot.DevKEY)
					log.Infof("电量: %d,  device key: %s", heartBeatInfo[10], robot.DevKEY)
					log.Infof("向前: %d,  device key: %s", heartBeatInfo[2], robot.DevKEY)
					log.Infof("向后: %d,  device key: %s", heartBeatInfo[3], robot.DevKEY)
					log.Infof("向左: %d,  device key: %s", heartBeatInfo[4], robot.DevKEY)
					log.Infof("向右: %d,  device key: %s", heartBeatInfo[5], robot.DevKEY)
					log.Infof("速度: %d,  device key: %s", heartBeatInfo[8], robot.DevKEY)
					log.Infof("房间: %d,  device key: %s", heartBeatInfo[9], robot.DevKEY)
					log.Infof("充电: %d,  device key: %s", heartBeatInfo[11], robot.DevKEY)
					log.Infof("停止: %d,  device key: %s", heartBeatInfo[13], robot.DevKEY)

					//如果ErrorCode不为0，则机器内部可能出错
					if heartBeatInfo[12] != 0 {
						log.Errorf("警告! ErrorCode: %d, device key: %s", heartBeatInfo[12], robot.DevKEY)
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

					// 判断结果
					if robot.IsCmdSent && robot.expResultIndex != 99 {
						var assertFlag string
						if heartBeatInfo[robot.expResultIndex] != robot.expResultValue {
							robot.notPassAndwaitNextHeartBeat++
							if robot.notPassAndwaitNextHeartBeat == 3 {
								assertFlag = "Failed"
								robot.notPassAndwaitNextHeartBeat = 0
							}
						} else {
							assertFlag = "Passed"
							robot.notPassAndwaitNextHeartBeat = 0
						}

						if robot.notPassAndwaitNextHeartBeat == 0 {
							testContent := fmt.Sprintf("%s, 预期: %s, 结果: %s", robot.checkPoint, strconv.Itoa(int(robot.expResultValue)),
								strconv.Itoa(int(heartBeatInfo[robot.expResultIndex])))
							result.WriteToResultFile(robot.ProdCode, robot.DevKEY, robot.checkPoint, testContent, assertFlag)
							robot.IsCmdSent = false
							robot.expResultIndex = 99
						}
					}
				}
			case <-robot.ServerMsg:
				// nothing to do.
			}
		}
	}()
}

// BuildProtoBufMessage 实现 EufyDevice 接口
func (robot *RobotCleaner) BuildProtoBufMessage() []byte {
	var payload []byte
	// 如果hangOn不为0，则表示机器人对前次发的指令响应的结果不正确，需等待下次心跳继续验证，不要发新指令过去
	if robot.notPassAndwaitNextHeartBeat != 0 {
		return nil
	}

	if !robot.charging {
		payload = getCommandToDevice(robot.caseNum, robot)
		robot.caseNum++
		if robot.caseNum > 10 {
			robot.caseNum = 1
		}
		log.Infof("发送指令给机器: %s, 指令内容: %s", robot.DevKEY, robot.checkPoint)
	} else if robot.returnCharge {
		payload = getCommandToDevice(4, robot)
		log.Infof("机器 [%s] 正处于低电量, 发指令叫它回家充电", robot.DevKEY)
	} else {
		log.Infof("机器 [%s] 正在充电中......, 现在不发任何指令.", robot.DevKEY)
	}
	return payload
}

func getCommandToDevice(index int, rb *RobotCleaner) []byte {
	rb.IsCmdSent = true
	//指令总数累加1
	rb.CmdSentQuantity++
	switch index {
	case 1:
		// 定点 0x01, ok
		rb.expResultIndex = 1
		rb.expResultValue = spot
		rb.checkPoint = "模式-->定点"
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE1, 0x01, 0xE2, 0xFA}
	case 2:
		//强力 0x01
		rb.expResultIndex = 8
		rb.expResultValue = strongSpeed
		rb.checkPoint = "速度-->强力"
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE8, 0x01, 0xE9, 0xFA}
	case 3:
		// 自动 0x02, ok
		rb.expResultIndex = 1
		rb.expResultValue = auto
		rb.checkPoint = "模式-->自动"
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE1, 0x02, 0xE3, 0xFA}
	case 4:
		// 返回充电 0x03, ok
		rb.expResultIndex = 1
		rb.expResultValue = charge
		rb.checkPoint = "模式-->返回充电"
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE1, 0x03, 0xE4, 0xFA}
	case 5:
		// 沿边 0x04, ok
		rb.expResultIndex = 1
		rb.expResultValue = edge
		rb.checkPoint = "模式-->沿边"
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE1, 0x04, 0xE5, 0xFA}
	case 6:
		//小房间 0x05
		rb.expResultIndex = 1
		rb.expResultValue = smallRoom
		rb.checkPoint = "模式为-->小房间"
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE1, 0x05, 0xE6, 0xFA}
	case 7:
		//日常 0x00
		rb.expResultIndex = 8
		rb.expResultValue = dailySpeed
		rb.checkPoint = "速度-->日常"
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE8, 0x00, 0xE8, 0xFA}
	case 8:
		// 暂停 0x00, ok
		rb.expResultIndex = 1
		rb.expResultValue = pause
		rb.checkPoint = "模式-->暂停"
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xE1, 0x00, 0xE1, 0xFA}
	// ------------------------- Find Me Alert: 0xEC -------------------------- good.
	case 9:
		// turn on alert: 0x01
		rb.checkPoint = "打开 [Find Me Alert]"
		log.Infof("执行: %s, device key: %s", rb.checkPoint, rb.DevKEY)
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xEC, 0x01, 0xED, 0xFA}
	case 10:
		// turn off alert: 0x00
		rb.checkPoint = "关闭 [Find Me Alert]"
		log.Infof("执行: %s, device key: %s", rb.checkPoint, rb.DevKEY)
		return []byte{0x00, 0x00, 0x00, 0xA5, 0xEC, 0x00, 0xEC, 0xFA}
	default:
		return nil
	}
}
