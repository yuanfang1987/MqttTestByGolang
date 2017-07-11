package device

import (
	"oceanwing/commontool"
	"strconv"

	log "github.com/cihub/seelog"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/protobuf/proto"
	lightT1012 "oceanwing/eufy/protobuf.lib/light/t1012"
)

// Light 灯泡类产品的一个 struct 描述，包括 T1011,T1012,T1013
type Light struct {
	baseDevice
	// status           string
	lum              uint32
	colorTemp        uint32
	isCtrlFunRunning bool
	stopCtrlFunc     chan struct{}
	onOffStatChan    chan string
	runMod           int
	resultMap        map[string]string
	occurSlice       []map[string]string
	awayModTesting   bool
}

// NewLight 新建一个 Light 实例
func NewLight(prodCode, devKey string) EufyDevice {
	o := &Light{}
	o.ProdCode = prodCode
	o.DevKEY = devKey
	o.SubDeviceTopic = "DEVICE/" + prodCode + "/" + devKey + "/PUH_MESSAGE" // 订阅设备的消息
	o.SubServerTopic = "DEVICE/" + prodCode + "/" + devKey + "/SUB_MESSAGE" //订阅服务器的消息
	o.SubMessage = make(chan MQTT.Message)
	o.stopCtrlFunc = make(chan struct{})
	o.resultMap = make(map[string]string)
	o.onOffStatChan = make(chan string, 2)
	log.Debugf("灯泡 %s (%s) 订阅设备主题: %s, 订阅服务器主题: %s", o.DevKEY, o.ProdCode, o.SubDeviceTopic, o.SubServerTopic)
	return o
}

// HandleSubscribeMessage 实现 EufyDevice 接口
func (light *Light) HandleSubscribeMessage() {
	go func() {
		log.Debugf("Running handleIncomingMsg function for device: %s", light.DevKEY)
		for {
			select {
			case msg := <-light.SubMessage:
				log.Infof("get new incoming message with topic: %s, for device: %s", msg.Topic(), light.DevKEY)
				go light.unMarshalAllMessage(msg)
			}
		}
	}()
}

// 解析心跳消息
func (light *Light) unMarshalHeartBeatMsg(incomingPayload []byte) {
	deviceMsg := &lightT1012.DeviceMessage{}
	err := proto.Unmarshal(incomingPayload, deviceMsg)
	if err != nil {
		log.Errorf("解析灯泡 %s (%s) 心跳消息失败: %s", light.DevKEY, light.ProdCode, err)
		return
	}

	noneParaMsg := deviceMsg.GetNonParaMsg()
	if noneParaMsg != nil {
		log.Infof("无参数消息的指令: %s", lightT1012.CmdType_name[int32(noneParaMsg.GetType())])
	}

	devBaseInfo := deviceMsg.GetReportDevBaseInfo()
	if devBaseInfo == nil {
		log.Warnf("提取灯泡 %s (%s) 基础信息失败", light.DevKEY, light.ProdCode)
		return
	}

	log.Infof("解析灯泡 %s (%s) 的心跳消息成功", light.DevKEY, light.ProdCode)

	// --------------------- 取出结果 --------------------------------------------

	//  CmdType
	cmd := devBaseInfo.GetType().String()
	log.Infof("灯泡 %s (%s) 指令类型: %s", light.DevKEY, light.ProdCode, cmd)

	// Mode
	modeName := devBaseInfo.GetMode().String()
	log.Infof("灯泡 %s (%s) 模式: %s", light.DevKEY, light.ProdCode, modeName)

	// Status
	stat := devBaseInfo.GetOnoffStatus().String()
	log.Infof("灯泡 %s (%s) 状态: %s", light.DevKEY, light.ProdCode, stat)

	// 把当前心跳的 status 存入 channel, 最多可以同时存 2 个
	light.onOffStatChan <- stat

	if modeName == "AWAY_MODE" {
		// 标记一下已经处于离家模式测试中
		light.awayModTesting = true

		// 记录由 normal mode 转为 leave mode 的时间, 如果字典中尚未有记录，则记录之
		if _, ok := light.resultMap["leave_Mode_Up"]; !ok {
			light.resultMap["leave_Mode_Up"] = commontool.GetCurrentTime()
			log.Infof("灯泡 %s (%s) 已启动离家模式...", light.DevKEY, light.ProdCode)
		}

		var prev, current string

		// 取出上次心跳和本次心跳的开关状态
		if len(light.onOffStatChan) == 2 {
			prev = <-light.onOffStatChan
			current = <-light.onOffStatChan
		}

		// 记录开启离家模式之前的那个状态
		if _, ok := light.resultMap["status_before_leave_mode"]; !ok {
			if prev != "" {
				light.resultMap["status_before_leave_mode"] = prev
				log.Infof("灯泡 %s (%s) 启动离家模式之前的状态是: %s", light.DevKEY, light.ProdCode, prev)
			}
		}

		// 如果两次的状态不同，则说明状态发生了变化，随机开关灯被触发, 把操作类型和时间记录下来
		if prev != current {
			leaveModeOccur := make(map[string]string)
			leaveModeOccur["occur_time"] = commontool.GetCurrentTime()
			leaveModeOccur["occur_type"] = current
			light.occurSlice = append(light.occurSlice, leaveModeOccur)
			log.Infof("灯泡 %s (%s) 随机开关为被触发, 本次是: %s", light.DevKEY, light.ProdCode, current)
		}

		// 把本次心跳的状态重新存入 channel 中，不然 channel 被掏空了，下次就没法一起取出两个来比较了
		if len(light.onOffStatChan) < 2 {
			light.onOffStatChan <- current
		}

	} else if modeName == "NORMAL_MODE" && light.awayModTesting {
		// 记录离家模式的失效时间
		if _, ok := light.resultMap["leave_Mode_Down"]; !ok {
			light.resultMap["leave_Mode_Down"] = commontool.GetCurrentTime()
		}

		// 记录恢复正常模式后在状态
		if _, ok := light.resultMap["status_resume_to_normal"]; !ok {
			light.resultMap["status_resume_to_normal"] = stat
		}

		light.awayModTesting = false
	}

	// 如果 channel 中缓存已满，必须取出一个, 否则下次心跳再往里面发数据的时候程序会死掉
	if len(light.onOffStatChan) == 2 {
		<-light.onOffStatChan
	}

	ligthCTRL := devBaseInfo.GetLightCtl()
	if ligthCTRL == nil {
		log.Error("GetLightCtl fail")
		return
	}

	// lum
	log.Infof("灯泡 %s (%s) 亮度: %d", light.DevKEY, light.ProdCode, ligthCTRL.GetLum())

	// 只有 T1012 和 T1013 才有色温
	if light.ProdCode != "T1011" {
		log.Infof("灯泡 %s (%s) 色温: %d", light.DevKEY, light.ProdCode, ligthCTRL.GetColorTemp())
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
		cmd := setlightdata.GetType().String()
		log.Infof("------指令类型: %s", cmd)

		lightctrl := setlightdata.GetLightCtl()
		if lightctrl != nil {
			log.Infof("------亮度: %d", lightctrl.GetLum())
			log.Infof("------色温: %d", lightctrl.GetColorTemp())
		}

		status := setlightdata.GetOnoffStatus().String()
		log.Infof("------开关状态: %s", status)
	}

	// time and alram
	timeAndAlarm := serMsg.GetSync_Time_Alarm()
	if timeAndAlarm != nil {
		log.Info("==解析出时间和闹钟的消息==")
		cmd := timeAndAlarm.GetType().String()
		log.Infof("------指令类型: %s", cmd)

		synctime := timeAndAlarm.GetTime()
		if synctime != nil {
			log.Info("------时间信息：")
			log.Infof("------年： %d", synctime.GetYear())
			log.Infof("------月： %d", synctime.GetMonth()+1)
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
					log.Infof("------闹钟，是否重复: %t", alarmMsg.GetRepetiton())
					weekinfo := convertToWeekDay(alarmMsg.GetWeekInfo())
					log.Infof("------闹钟，WeekInfo：%s", weekinfo)
				}

				alarmEvent := v.GetAlarmEvent()
				if alarmEvent != nil {
					eventType := alarmEvent.GetType().String()
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
		cmd := powerUpLightStatus.GetType().String()
		log.Infof("------指令类型: %s", cmd)
		status := powerUpLightStatus.GetPowrupStatus().String()
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
		cmd := awayMod.GetType().String()
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
			log.Infof("------离家模式, WeekInfo: %s", convertToWeekDay(leaveMsg.GetWeekInfo()))
			leaveHomeFlag := leaveMsg.GetLeaveHomeState()
			log.Infof("------离家模式, 是否开启: %t", leaveHomeFlag)

			// 如果服务器下发的 离家模式 为 true，则记下时间
			if leaveHomeFlag {
				light.resultMap["leave_mode_up_exp"] = strconv.Itoa(starthour) + ":" + strconv.Itoa(startminute)
				light.resultMap["leave_mode_down_exp"] = strconv.Itoa(finishhour) + ":" + strconv.Itoa(finishminute)
			}

		}
	}

}

func (light *Light) unMarshalAllMessage(msg MQTT.Message) {
	t := msg.Topic()
	payload := msg.Payload()

	if light.SubDeviceTopic == t {
		// 设备心跳消息
		log.Info("----- 这是一个来自设备的心跳消息----------")
		light.unMarshalHeartBeatMsg(payload)
	} else if light.SubServerTopic == t {
		//服务器消息
		log.Info("-------这是一个来自服务器的控制消息---------")
		light.unMarshalServerMsg(payload)
	}
}

// LeaveModeTestResult 生成汇总报告, 实现 EufyDevice 接口
func (light *Light) LeaveModeTestResult() {
	log.Infof("====================== %s (%s) 离家模式测试结果 ======================", light.DevKEY, light.ProdCode)

	// 预计开始时间
	if expStart, ok := light.resultMap["leave_mode_up_exp"]; ok {
		log.Infof("预计开始时间: %s", expStart)
	} else {
		log.Error("无法获取到预计开始时间")
	}

	// 实际开始时间
	if startTime, ok := light.resultMap["leave_Mode_Up"]; ok {
		log.Infof("实际开始时间: %s", startTime)
	} else {
		log.Error("没有获取到实际开始时间")
	}

	// 预计结束时间
	if expEnd, ok := light.resultMap["leave_mode_down_exp"]; ok {
		log.Infof("预计结束时间: %s", expEnd)
	} else {
		log.Error("无法获取预计结束时间")
	}

	// 实际结束时间
	if endTime, ok := light.resultMap["leave_Mode_Down"]; ok {
		log.Infof("实际结束时间: %s", endTime)
	} else {
		log.Error("没有获取到实际结束时间")
	}

	// 开始前的状态
	if statusBefore, ok := light.resultMap["status_before_leave_mode"]; ok {
		log.Infof("开始前状态：%s", statusBefore)
	} else {
		log.Error("无法获取开始前的状态")
	}

	// 结束后的状态
	if statusAfter, ok := light.resultMap["status_resume_to_normal"]; ok {
		log.Infof("结束后状态: %s", statusAfter)
	} else {
		log.Error("无法获取结束后状态")
	}

	// 随机触发开关灯情况
	for i, v := range light.occurSlice {
		log.Infof("第 %d 次发生, 时间: %s, 类型: %s", i+1, v["occur_time"], v["occur_type"])
	}

}

func convertToWeekDay(v uint32) string {
	var weekinfo string
	if (v & 1) > 0 {
		weekinfo = "星期一, "
	}
	if (v & 2) > 0 {
		weekinfo = weekinfo + "星期二, "
	}
	if (v & 4) > 0 {
		weekinfo = weekinfo + "星期三, "
	}
	if (v & 8) > 0 {
		weekinfo = weekinfo + "星期四, "
	}
	if (v & 16) > 0 {
		weekinfo = weekinfo + "星期五, "
	}
	if (v & 32) > 0 {
		weekinfo = weekinfo + "星期六, "
	}
	if (v & 64) > 0 {
		weekinfo = weekinfo + "星期日"
	}
	return weekinfo
}
