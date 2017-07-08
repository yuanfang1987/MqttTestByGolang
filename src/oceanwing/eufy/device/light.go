package device

import (
	"fmt"
	"math"
	"math/rand"
	"oceanwing/commontool"
	"oceanwing/eufy/result"
	"time"

	log "github.com/cihub/seelog"
	"github.com/golang/protobuf/proto"
	lightEvent "oceanwing/eufy/protobuf.lib/light/lightevent"
	lightT1012 "oceanwing/eufy/protobuf.lib/light/t1012"
)

// Light 灯泡类产品的一个 struct 描述，包括 T1011,T1012,T1013
type Light struct {
	baseDevice
	mode               lightT1012.DeviceMessage_ReportDevBaseInfo_LIGHT_DEV_MODE
	status             lightT1012.LIGHT_ONOFF_STATUS
	lum                uint32
	lumTemp            uint32 // 临时用来存放 lum
	colorTemp          uint32
	modeLeaveHome      lightT1012.DeviceMessage_ReportDevBaseInfo_LIGHT_DEV_MODE
	statusLeaveHome    lightT1012.LIGHT_ONOFF_STATUS
	lumLeaveHome       uint32
	colorTempLeaveHome uint32
	isCtrlFunRunning   bool
	stopCtrlFunc       chan struct{}
}

// NewLight 新建一个 Light 实例.
func NewLight(prodCode, devKey, devid string) EufyDevice {
	o := &Light{
		mode:   0,
		status: 1,
	}
	o.ProdCode = prodCode
	o.DevKEY = devKey
	o.DevID = devid
	o.PubTopicl = "DEVICE/" + prodCode + "/" + devKey + "/SUB_MESSAGE"
	o.SubTopicl = "DEVICE/" + prodCode + "/" + devKey + "/PUH_MESSAGE"
	o.DeviceMsg = make(chan []byte)
	o.ServerMsg = make(chan []byte)
	o.stopCtrlFunc = make(chan struct{})
	log.Infof("Create a Light, product code: %s, device key: %s, device id: %s", prodCode, devKey, devid)
	return o
}

// HandleSubscribeMessage 实现 EufyDevice 接口
func (light *Light) HandleSubscribeMessage() {
	go func() {
		log.Debugf("Running handleIncomingMsg function for device: %s", light.DevKEY)
		for {
			select {
			case devmsg := <-light.DeviceMsg:
				log.Infof("get new incoming message from device: %s", light.DevKEY)
				light.unMarshalHeartBeatMsg(devmsg)
			case serMsg := <-light.ServerMsg:
				log.Info("get new incoming message from server")
				light.unMarshalServerMsg(serMsg)
			}
		}
	}()
}

// BuildProtoBufMessage 实现 EufyDevice 接口
func (light *Light) BuildProtoBufMessage() []byte {
	// 如果在跑离家模式，则不发任何指令
	if light.RunMod == 1 {
		log.Info("正在跑离家模式......")
		return nil
	}

	// 如果上一次的测试结果没有通过，则被挂起，则先不要发新的指令过去
	if light.HangOn != 0 {
		log.Warnf("上一次测试未通过，需等待新的心跳消息来继续验证， HangOn: %d", light.HangOn)
		return nil
	}
	o := light.setLightBrightAndColor()
	data, err := proto.Marshal(o)

	if err != nil {
		log.Errorf("build set light data message fail: %s", err.Error())
		return nil
	}

	log.Info("=================================================")

	// 设置 IsCmdSent 标志为 true
	light.IsCmdSent = true
	// 已下发的指令数量累加 1
	light.CmdSentQuantity++

	return data
}

// SetLightData is a struct
// brightness: 亮度，color: 色温,  ServerMessage_SetLightData
func (light *Light) setLightBrightAndColor() *lightT1012.ServerMessage {

	seed := commontool.RandInt64(0, 10)
	var content string
	var lightData *lightT1012.ServerMessage_SetLightData_
	// seed 随机数产生的范围是 0 到 9 共 10 个数字，则用 30%的概率去执行开关灯， 剩下的执行调节亮度色温
	if seed < 3 {
		var nextStatus *lightT1012.LIGHT_ONOFF_STATUS
		// 如果灯的当前状态是开着的，则执行关闭操作， 反之则执行打开操作
		if light.status == lightT1012.LIGHT_ONOFF_STATUS_ON {
			nextStatus = lightT1012.LIGHT_ONOFF_STATUS_OFF.Enum()
			log.Info("关灯")
			// 关灯后， 亮度变成 0, 色温保持和关灯前一样
			light.status = 0
			light.lum = 0
			// light.colorTemp = 0
			content = "关灯"
		} else {
			nextStatus = lightT1012.LIGHT_ONOFF_STATUS_ON.Enum()
			log.Info("开灯")
			// 开灯后，亮度为100，色温为0，but why???
			light.status = 1
			light.lum = 100
			light.colorTemp = 0
			content = "开灯"
		}

		lightData = &lightT1012.ServerMessage_SetLightData_{
			SetLightData: &lightT1012.ServerMessage_SetLightData{
				Type:        lightT1012.CmdType_REMOTE_SET_LIGHTING_PARA.Enum(),
				OnoffStatus: nextStatus,
			},
		}
	} else {
		// 调节亮度和色温, 随机产生亮度和色温的值
		brightness := uint32(commontool.RandInt64(0, 101))
		color := uint32(commontool.RandInt64(0, 101))
		light.lum = brightness
		light.colorTemp = color
		light.status = 1
		log.Infof("执行调节亮度色温操作, lum: %d, colorTemp: %d", brightness, color)
		content = "调节亮度和色温"

		lightData = &lightT1012.ServerMessage_SetLightData_{
			SetLightData: &lightT1012.ServerMessage_SetLightData{
				Type: lightT1012.CmdType_REMOTE_SET_LIGHTING_PARA.Enum(),
				LightCtl: &lightEvent.LampLightLevelCtlMessage{
					Lum:       proto.Uint32(brightness),
					ColorTemp: proto.Uint32(color),
				},
			},
		}
	}

	// 在.csv 结果文件上打个标志
	result.WriteToResultFile(light.ProdCode, light.DevKEY, "NA", content, "NA")

	o := &lightT1012.ServerMessage{
		SessionId:     proto.Int32(rand.Int31n(math.MaxInt32)),
		RemoteMessage: lightData,
	}

	return o
}

// ControlAwayModStatus 实现了 EufyDevice 接口
func (light *Light) ControlAwayModStatus(startHour, startMinute, endHour, endMinute int) {
	go func() {
		light.isCtrlFunRunning = true
		interval := time.NewTicker(time.Second * 1).C
		for {
			select {
			case <-interval:
				nowTime := time.Now()
				// 当前时间与开始时间作对比， 如何相等， 则标识 RunMod 为离家模式
				if nowTime.Hour() == startHour && nowTime.Minute() == startMinute && nowTime.Second() == 0 {
					light.RunMod = 1
					light.mode = 1
					log.Infof("灯泡 %s (%s) 开启离家模式", light.DevKEY, light.ProdCode)
					result.WriteToResultFile(light.ProdCode, light.DevKEY, "NA", "开启离家模式", "NA")
				}
				// 当前时间与结束时间作对比， 如何相等， 则标识 RunMod 为正常模式
				if nowTime.Hour() == endHour && nowTime.Minute() == endMinute && nowTime.Second() == 0 {
					light.RunMod = 0
					light.mode = 0
					log.Infof("灯泡 %s (%s) 恢复正常模式", light.DevKEY, light.ProdCode)
					result.WriteToResultFile(light.ProdCode, light.DevKEY, "NA", "恢复正常模式", "NA")
				}
			case <-light.stopCtrlFunc:
				light.isCtrlFunRunning = false
				//干掉这个函数
				return
			}
		}
	}()
}

// leave 模式下
func (light *Light) leaveModeTest(devInfo *lightT1012.DeviceMessage_ReportDevBaseInfo) {
	var l uint32
	var c uint32
	var assertFlag string
	var testContent string
	mod := devInfo.GetMode()
	stat := devInfo.GetOnoffStatus()
	modStr := lightT1012.DeviceMessage_ReportDevBaseInfo_LIGHT_DEV_MODE_name[int32(mod)]
	statStr := lightT1012.LIGHT_ONOFF_STATUS_name[int32(stat)]
	log.Infof("灯泡 %s (%s) 模式: %s, 状态: %s", light.DevKEY, light.ProdCode, modStr, statStr)
	lctrl := devInfo.GetLightCtl()
	if lctrl != nil {
		l = lctrl.GetLum()
		log.Infof("灯泡 %s (%s) 亮度：%d", light.DevKEY, light.ProdCode, l)
		if light.ProdCode != "T1011" {
			c = lctrl.GetColorTemp()
			log.Infof("灯泡 %s (%s) 色温: %d", light.DevKEY, light.ProdCode, c)
		}
	}

	// 判断心跳数据, mode 字段必须是 1
	assertFlag = light.PassedOrFailed(1 == mod)
	testContent = fmt.Sprintf("灯泡 %s (%s) Mode, 预期: %d, 实际: %d", light.DevKEY, light.ProdCode, 1, mod)
	result.WriteToResultFile(light.ProdCode, light.DevKEY, "Mode", testContent, assertFlag)

	// 如果随机发生了开关灯， 则记录下来
	if light.statusLeaveHome != stat {
		light.statusLeaveHome = stat // 更新  mode
		var str string
		if stat == lightT1012.LIGHT_ONOFF_STATUS_ON {
			str = "离家模式随机开灯"
			light.lumLeaveHome = light.lumTemp // 重新开灯后， 把临时变量的值赋值给 lum
		} else {
			str = "离家模式随机关灯"
			light.lumTemp = light.lumLeaveHome // 把当前亮度存放在一个临时变量中, 待下次开灯时，要拿出来对比
			light.lumLeaveHome = 0             //关灯之后， 亮度是0
		}
		result.WriteToResultFile(light.ProdCode, light.DevKEY, str, "NA", "NA")
		log.Infof("离家模式随机开关灯被触发, 本次是: %s", statStr)
	}

	// 判断亮度
	assertFlag = light.PassedOrFailed(light.lumLeaveHome == l)
	testContent = fmt.Sprintf("灯泡 %s (%s) lum, 预期: %d, 实际: %d", light.DevKEY, light.ProdCode, light.lumLeaveHome, l)
	result.WriteToResultFile(light.ProdCode, light.DevKEY, "离家模式-亮度", testContent, assertFlag)

	// 非 T1011 的才有色温
	if light.ProdCode != "T1011" {
		assertFlag = light.PassedOrFailed(light.colorTempLeaveHome == c)
		testContent = fmt.Sprintf("灯泡 %s (%s) colorTemp, 预期: %d, 实际: %d", light.DevKEY, light.ProdCode, light.colorTempLeaveHome, c)
		result.WriteToResultFile(light.ProdCode, light.DevKEY, "离家模式-色温", testContent, assertFlag)
	}
}

// 解析设备心跳
func (light *Light) unMarshalHeartBeatMsg(incomingPayload []byte) {
	deviceMsg := &lightT1012.DeviceMessage{}
	err := proto.Unmarshal(incomingPayload, deviceMsg)
	if err != nil {
		log.Errorf("解析灯泡 %s (%s) 心跳消息失败: %s", light.DevKEY, light.ProdCode, err)
		return
	}

	noneParaMsg := deviceMsg.GetNonParaMsg()
	if noneParaMsg != nil {
		log.Infof("无参数消息, 指令类型: %s", noneParaMsg.GetType().String())
	}

	devBaseInfo := deviceMsg.GetReportDevBaseInfo()
	if devBaseInfo == nil {
		// log.Warnf("提取灯泡 %s (%s) 基础信息失败", light.DevKEY, light.ProdCode)
		return
	}

	log.Infof("解析灯泡 %s (%s) 的心跳消息成功", light.DevKEY, light.ProdCode)

	// 如果是离家模式, 如何触发 RunMod = 1?
	if light.RunMod == 1 {
		light.leaveModeTest(devBaseInfo)
		return
	}

	// --------------------- 判断结果 --------------------------------------------

	// 只有在给设备下发了指令之后，才去判断它的即时心跳， 常规心跳不要管
	if !light.IsCmdSent {
		log.Info("尚未有指令下发给设备，无需判断心跳")
		return
	}

	// 预先假设测试结果为 passed
	light.IsTestPassed = true

	// 已解析的心跳数量累加 1
	light.DecodeHeartBeatMsgQuantity++

	// 先用一个字典来存放测试结果, 注意: 这是一个嵌套Map
	resultMap := make(map[string]map[string]string)

	var assertFlag string
	var testContent string

	//  CmdType
	assertFlag = light.PassedOrFailed(lightT1012.CmdType_DEV_REPORT_STATUS == devBaseInfo.GetType())
	testContent = fmt.Sprintf("灯泡 %s (%s) CmdType, 预期: %s, 实际: %s", light.DevKEY, light.ProdCode, lightT1012.CmdType_DEV_REPORT_STATUS.String(), devBaseInfo.GetType().String())
	log.Info(testContent)

	cmdTypeResultMap := make(map[string]string)
	cmdTypeResultMap["content"] = testContent
	cmdTypeResultMap["flag"] = assertFlag

	resultMap["CmdType"] = cmdTypeResultMap

	// Mode
	assertFlag = light.PassedOrFailed(light.mode == devBaseInfo.GetMode())
	testContent = fmt.Sprintf("灯泡 %s (%s) Mode, 预期: %s, 实际: %s", light.DevKEY, light.ProdCode, light.mode.String(), devBaseInfo.GetMode().String())
	log.Info(testContent)

	modeResuleMap := make(map[string]string)
	modeResuleMap["content"] = testContent
	modeResuleMap["flag"] = assertFlag

	resultMap["Mode"] = modeResuleMap

	light.modeLeaveHome = devBaseInfo.GetMode() // 记录当前灯的模式

	// Status
	assertFlag = light.PassedOrFailed(light.status == devBaseInfo.GetOnoffStatus())
	testContent = fmt.Sprintf("灯泡 %s (%s) Status, 预期: %s, 实际: %s", light.DevKEY, light.ProdCode, light.status.String(), devBaseInfo.GetOnoffStatus().String())
	log.Info(testContent)

	statusResultMap := make(map[string]string)
	statusResultMap["content"] = testContent
	statusResultMap["flag"] = assertFlag

	resultMap["Status"] = statusResultMap

	light.statusLeaveHome = devBaseInfo.GetOnoffStatus() // 记录当前灯的开关状态

	ligthCTRL := devBaseInfo.GetLightCtl()
	if ligthCTRL != nil {
		// lum
		assertFlag = light.PassedOrFailed(light.lum == ligthCTRL.GetLum())
		testContent = fmt.Sprintf("灯泡 %s (%s) lum, 预期: %d, 实际: %d", light.DevKEY, light.ProdCode, light.lum, ligthCTRL.GetLum())
		log.Info(testContent)

		lumResultMap := make(map[string]string)
		lumResultMap["content"] = testContent
		lumResultMap["flag"] = assertFlag

		resultMap["Lum"] = lumResultMap

		light.lumLeaveHome = ligthCTRL.GetLum() // 记录当前灯的亮度
		light.lumTemp = ligthCTRL.GetLum()

		// 只有 T1012 和 T1013 才有色温
		if light.ProdCode != "T1011" {
			assertFlag = light.PassedOrFailed(light.colorTemp == ligthCTRL.GetColorTemp())
			testContent = fmt.Sprintf("灯泡 %s (%s) ColorTemp, 预期: %d, 实际: %d", light.DevKEY, light.ProdCode, light.colorTemp, ligthCTRL.GetColorTemp())
			log.Info(testContent)

			colorTempResultMap := make(map[string]string)
			colorTempResultMap["content"] = testContent
			colorTempResultMap["flag"] = assertFlag

			resultMap["ColorTemp"] = colorTempResultMap

			light.colorTempLeaveHome = ligthCTRL.GetColorTemp() // 记录当前灯的色温
		}
	}

	if !light.IsTestPassed {
		light.HangOn++
		if light.HangOn < 3 {
			log.Error("当前测试结果未能通过， 需等待下一轮心跳验证")
			return
		}
	}

	// 重置
	light.IsCmdSent = false
	light.HangOn = 0

	// 写入 csv 文件
	for key, m := range resultMap {
		result.WriteToResultFile(light.ProdCode, light.DevKEY, key, m["content"], m["flag"])
	}

}

func (light *Light) unMarshalServerMsg(incomingPayload []byte) {
	serMsg := &lightT1012.ServerMessage{}
	err := proto.Unmarshal(incomingPayload, serMsg)
	if err != nil {
		log.Errorf("解析服务器消息失败: %s", err)
		return
	}

	// session id
	log.Infof("Session ID: %d", serMsg.GetSessionId())

	// SetLightData
	setlightdata := serMsg.GetSetLightData()
	if setlightdata != nil {
		log.Info("==解析出控制灯泡亮度色温的消息==")
		cmd := lightT1012.CmdType_name[int32(setlightdata.GetType())]
		log.Infof("------指令类型: %s", cmd)

		lightctrl := setlightdata.GetLightCtl()
		if lightctrl != nil {
			log.Infof("------亮度: %d", lightctrl.GetLum())
			log.Infof("------色温: %d", lightctrl.GetColorTemp())
		}

		status := lightT1012.LIGHT_ONOFF_STATUS_name[int32(setlightdata.GetOnoffStatus())]
		log.Infof("------开关状态: %s", status)
	}

	// time and alram
	timeAndAlarm := serMsg.GetSync_Time_Alarm()
	if timeAndAlarm != nil {
		log.Info("==解析出时间和闹钟的消息==")
		cmd := lightT1012.CmdType_name[int32(timeAndAlarm.GetType())]
		log.Infof("------指令类型: %s", cmd)

		synctime := timeAndAlarm.GetTime()
		if synctime != nil {
			log.Info("------时间信息：")
			log.Infof("------年： %d", synctime.GetYear())
			log.Infof("------月： %d", synctime.GetMonth())
			log.Infof("------日： %d", synctime.GetDay())
			log.Infof("------weekday: %d", synctime.GetWeekday())
			log.Infof("------时： %d", synctime.GetHours())
			log.Infof("------分： %d", synctime.GetMinutes())
			log.Infof("------秒： %d", synctime.GetSeconds())
		}

		alarm := timeAndAlarm.GetAlarm()
		if alarm != nil {
			log.Info("------闹钟信息：")
			alarmRecordData := alarm.GetAlarmRecordData()
			for i, v := range alarmRecordData {
				log.Infof("=====第 %d 个闹钟信息====", i+1)
				alarmMsg := v.GetAlarmMesage()
				if alarmMsg != nil {
					log.Infof("------闹钟，时：%d", alarmMsg.GetHours())
					log.Infof("------闹钟，分: %d", alarmMsg.GetMinutes())
					// log.Infof("------闹钟，秒：%d", alarmMsg.get) 旧版protobuf定义文件没有秒，新版有
					log.Infof("------闹钟，是否重复: %t", alarmMsg.GetRepetiton())
					weekinfo := commontool.ConvertToWeekDay(alarmMsg.GetWeekInfo())
					log.Infof("------闹钟，WeekInfo：%s", weekinfo)
				}

				alarmEvent := v.GetAlarmEvent()
				if alarmEvent != nil {
					eventType := lightT1012.SyncAlarmRecordMessage_EventType_name[int32(alarmEvent.GetType())]
					log.Infof("------闹钟事件, 类型：%s", eventType)
					lightCTRL := alarmEvent.GetLightCtl()
					if lightCTRL != nil {
						log.Infof("------闹钟事件, 亮度: %d", lightCTRL.GetLum())
						log.Infof("------闹钟事件, 色温: %d", lightCTRL.GetColorTemp())
					}
				}
			}
		}
	}

	// Power up light status
	powerUpLightStatus := serMsg.GetSetPowerupLightStatus()
	if powerUpLightStatus != nil {
		log.Info("==解析出灯泡上电时的初始状态信息==")
		cmd := lightT1012.CmdType_name[int32(powerUpLightStatus.GetType())]
		log.Infof("------指令类型: %s", cmd)
		status := lightT1012.ServerMessage_SetPowerUpLightStatus_POWERUP_LIGHT_STATUS_name[int32(powerUpLightStatus.GetPowrupStatus())]
		log.Infof("------状态: %s", status)
		lightCTRL := powerUpLightStatus.GetLightCtl()
		if lightCTRL != nil {
			log.Infof("------亮度: %d", lightCTRL.GetLum())
			log.Infof("------色温: %d", lightCTRL.GetColorTemp())
		}
	}

	// away mode
	awayMod := serMsg.GetSetAwayMode_Status()
	if awayMod != nil {
		log.Info("==解析出离家模式信息==")
		cmd := lightT1012.CmdType_name[int32(awayMod.GetType())]
		log.Infof("------指令类型: %s", cmd)
		leaveMsg := awayMod.GetSyncLeaveModeMsg()
		if leaveMsg != nil {
			starthour := int(leaveMsg.GetStartHours())
			startminute := int(leaveMsg.GetStartMinutes())
			finishhour := int(leaveMsg.GetFinishHours())
			finishminute := int(leaveMsg.GetFinishMinutes())
			log.Infof("------离家模式, 开始时间  %d:%d", starthour, startminute)
			log.Infof("------离家模式, 结束时间  %d:%d", finishhour, finishminute)
			log.Infof("------离家模式, 是否重复: %t", leaveMsg.GetRepetiton())
			log.Infof("------离家模式, WeekInfo: %s", commontool.ConvertToWeekDay(leaveMsg.GetWeekInfo()))
			leaveHomeFlag := leaveMsg.GetLeaveHomeState()
			log.Infof("------离家模式, 是否开启: %t", leaveHomeFlag)
			// 如果服务器下发的 离家模式 为 true，且控制函数没有运行，则跑之.
			if leaveHomeFlag && (!light.isCtrlFunRunning) {
				light.ControlAwayModStatus(starthour, startminute, finishhour, finishminute)
			}
			// 如果服务器下发的 离家模式 为 false，且控制函数已运行，则发信号干掉之.
			if (!leaveHomeFlag) && light.isCtrlFunRunning {
				light.stopCtrlFunc <- struct{}{}
			}
		}
	}

}

// func buildWeekdays(weekdays []int64) uint32 {
// 	var result uint32
// 	for _, d := range weekdays {
// 		devDay := convertDeviceWeekday(d)
// 		result += uint32(1 << uint64(devDay-1))
// 	}
// 	return result
// }

// Convert 0,1,2,3,4,5,6 -> 1,2,3,4,5,6,7
// func convertDeviceWeekday(weekday int64) int64 {
// 	if weekday == int64(time.Sunday) {
// 		return 7
// 	}
// 	return weekday
// }
