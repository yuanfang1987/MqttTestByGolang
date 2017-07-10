package device

import (
	log "github.com/cihub/seelog"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/protobuf/proto"
	lightT1012 "oceanwing/eufy/protobuf.lib/light/t1012"
)

// Light 灯泡类产品的一个 struct 描述，包括 T1011,T1012,T1013
type Light struct {
	baseDevice
	mode      lightT1012.DeviceMessage_ReportDevBaseInfo_LIGHT_DEV_MODE
	status    lightT1012.LIGHT_ONOFF_STATUS
	lum       uint32
	colorTemp uint32
}

// NewLight 新建一个 Light 实例
func NewLight(prodCode, devKey string) EufyDevice {
	o := &Light{
		mode:   0,
		status: 1,
	}
	o.ProdCode = prodCode
	o.DevKEY = devKey
	o.SubDeviceTopic = "DEVICE/" + prodCode + "/" + devKey + "/PUH_MESSAGE" // 订阅设备的消息
	o.SubServerTopic = "DEVICE/" + prodCode + "/" + devKey + "/SUB_MESSAGE" //订阅服务器的消息
	o.SubMessage = make(chan MQTT.Message)
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
				light.unMarshalAllMessage(msg)
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
	cmd := lightT1012.CmdType_name[int32(devBaseInfo.GetType())]
	log.Infof("灯泡 %s (%s) 指令类型: %s", light.DevKEY, light.ProdCode, cmd)

	// Mode
	modeName := lightT1012.DeviceMessage_ReportDevBaseInfo_LIGHT_DEV_MODE_name[int32(devBaseInfo.GetMode())]
	log.Infof("灯泡 %s (%s) 模式: %s", light.DevKEY, light.ProdCode, modeName)

	// Status
	status := lightT1012.LIGHT_ONOFF_STATUS_name[int32(devBaseInfo.GetOnoffStatus())]
	log.Infof("灯泡 %s (%s) 状态: %s", light.DevKEY, light.ProdCode, status)

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
					// log.Infof("------闹钟，秒：%d", alarmMsg.get) 旧版protobuf定义文件没有秒，新版有
					log.Infof("------闹钟，是否重复: %t", alarmMsg.GetRepetiton())
					weekinfo := convertToWeekDay(alarmMsg.GetWeekInfo())
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
			log.Infof("------离家模式, 开始时间  %d:%d", leaveMsg.GetStartHours(), leaveMsg.GetStartMinutes())
			log.Infof("------离家模式, 结束时间  %d:%d", leaveMsg.GetFinishHours(), leaveMsg.GetFinishMinutes())
			log.Infof("------离家模式, 是否重复: %t", leaveMsg.GetRepetiton())
			log.Infof("------离家模式, WeekInfo: %s", convertToWeekDay(leaveMsg.GetWeekInfo()))
			log.Infof("------离家模式, 是否开启: %t", leaveMsg.GetLeaveHomeState())
		}
	}

}

func (light *Light) unMarshalAllMessage(msg MQTT.Message) {
	t := msg.Topic()
	payload := msg.Payload()

	if light.SubDeviceTopic == t {
		// 设备心跳消息
		log.Info("----- 这是一个来自设备的心跳消息----------")
		go light.unMarshalHeartBeatMsg(payload)
	} else if light.SubServerTopic == t {
		//服务器消息
		log.Info("-------这是一个来自服务器的控制消息---------")
		go light.unMarshalServerMsg(payload)
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
