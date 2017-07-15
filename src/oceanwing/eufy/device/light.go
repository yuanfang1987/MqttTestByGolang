package device

import (
	"fmt"
	"math"
	"math/rand"
	"oceanwing/commontool"

	"oceanwing/eufy/result"

	log "github.com/cihub/seelog"
	"github.com/golang/protobuf/proto"
	lightEvent "oceanwing/eufy/protobuf.lib/light/lightevent"
	lightT1012 "oceanwing/eufy/protobuf.lib/light/t1012"
)

// Light 灯泡类产品的一个 struct 描述，包括 T1011,T1012,T1013
type Light struct {
	baseDevice
	mode      lightT1012.DeviceMessage_ReportDevBaseInfo_LIGHT_DEV_MODE
	status    lightT1012.LIGHT_ONOFF_STATUS
	lum       uint32
	colorTemp uint32
	lumTemp   uint32
}

// NewLight 新建一个 Light 实例.
func NewLight(prodCode, devKey string) EufyDevice {
	o := &Light{}
	o.ProdCode = prodCode
	o.DevKEY = devKey
	o.PubTopicl = "DEVICE/" + prodCode + "/" + devKey + "/SUB_MESSAGE"
	o.SubTopicl = "DEVICE/" + prodCode + "/" + devKey + "/PUH_MESSAGE"
	o.DeviceMsg = make(chan []byte)
	o.ServerMsg = make(chan []byte)
	log.Infof("Create a Light, product code: %s, device key: %s", prodCode, devKey)
	return o
}

// HandleSubscribeMessage 实现 EufyDevice 接口
func (light *Light) HandleSubscribeMessage() {
	go func() {
		log.Debugf("Running handleIncomingMsg function for device: %s", light.DevKEY)
		for {
			select {
			case devmsg := <-light.DeviceMsg:
				log.Infof("======设备上报消息: %s======", light.DevKEY)
				light.unMarshalHeartBeatMsg(devmsg)
			case serMsg := <-light.ServerMsg:
				log.Info("======服务器控制消息======")
				light.unMarshalServerMsg(serMsg)
			}
		}
	}()
}

// BuildProtoBufMessage 实现 EufyDevice 接口
func (light *Light) BuildProtoBufMessage() []byte {
	// 如果上一次的测试结果没有通过，则被挂起，则先不要发新的指令过去
	if light.notPassAndwaitNextHeartBeat != 0 {
		log.Warnf("上次心跳验证未通过，等待下次验证，当前已验证 %d 次", light.notPassAndwaitNextHeartBeat)
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

	return data
}

// SetLightData is a struct
// brightness: 亮度，color: 色温,  ServerMessage_SetLightData
func (light *Light) setLightBrightAndColor() *lightT1012.ServerMessage {

	seed := commontool.GenerateRandNumber(0, 10)
	var lightData *lightT1012.ServerMessage_SetLightData_
	// seed 随机数产生的范围是 0 到 9 共 10 个数字，则用 30%的概率去执行开关灯， 剩下的执行调节亮度色温
	if seed < 3 {
		var nextStatus *lightT1012.LIGHT_ONOFF_STATUS
		// 如果灯的当前状态是开着的，则执行关闭操作， 反之则执行打开操作
		if light.status == lightT1012.LIGHT_ONOFF_STATUS_ON {
			nextStatus = lightT1012.LIGHT_ONOFF_STATUS_OFF.Enum()
			log.Info("关灯")
			// 关灯后， 亮度变成 0, 色温保持和关灯前一样
			light.status = lightT1012.LIGHT_ONOFF_STATUS_OFF
			light.lumTemp = light.lum
			light.lum = 0
			light.testcase = "关灯"
		} else {
			nextStatus = lightT1012.LIGHT_ONOFF_STATUS_ON.Enum()
			log.Info("开灯")
			// 开灯后，亮度为100，色温为0，but why???
			light.status = lightT1012.LIGHT_ONOFF_STATUS_ON
			// light.lum = 100
			// light.colorTemp = 0
			light.lum = light.lumTemp
			light.testcase = "开灯"
		}

		lightData = &lightT1012.ServerMessage_SetLightData_{
			SetLightData: &lightT1012.ServerMessage_SetLightData{
				Type:        lightT1012.CmdType_REMOTE_SET_LIGHTING_PARA.Enum(),
				OnoffStatus: nextStatus,
			},
		}
	} else {
		// 调节亮度和色温, 随机产生亮度和色温的值.
		brightness := uint32(commontool.GenerateRandNumber(10, 100))
		color := uint32(commontool.GenerateRandNumber(10, 100))
		light.lum = brightness
		light.colorTemp = color
		light.status = lightT1012.LIGHT_ONOFF_STATUS_ON
		log.Infof("执行调节亮度色温操作, lum: %d, colorTemp: %d", brightness, color)
		light.testcase = "调节亮度和色温"

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

	o := &lightT1012.ServerMessage{
		SessionId:     proto.Int32(rand.Int31n(math.MaxInt32)),
		RemoteMessage: lightData,
	}

	return o
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

	// --------------------- 判断结果 --------------------------------------------

	// 只有在给设备下发了指令之后，才去判断它的即时心跳， 常规心跳不要管
	if !light.IsCmdSent {
		log.Info("尚未有指令下发给设备，无需判断心跳")
		return
	}

	var errMsg []string

	//  CmdType
	cmdtype := devBaseInfo.GetType().String()
	if cmdtype != lightT1012.CmdType_DEV_REPORT_STATUS.String() {
		errMsg = append(errMsg, fmt.Sprintf("assert CmdType fail, exp: %s, act: %s", lightT1012.CmdType_DEV_REPORT_STATUS.String(), cmdtype))
	}
	log.Infof("白灯 %s (%s) CmdType: %s", light.DevKEY, light.ProdCode, cmdtype)

	// Mode
	mode := devBaseInfo.GetMode().String()
	expMod := lightT1012.DeviceMessage_ReportDevBaseInfo_NORMAL_MODE.String()
	if expMod != mode {
		errMsg = append(errMsg, fmt.Sprintf("assert mode fail, exp: %s, act: %s", expMod, mode))
	}
	log.Infof("白灯 %s (%s) 模式: %s", light.DevKEY, light.ProdCode, mode)

	// Status
	status := devBaseInfo.GetOnoffStatus().String()
	if light.status.String() != status {
		errMsg = append(errMsg, fmt.Sprintf("assert onOff status fail, exp: %s, act: %s", light.status.String(), status))
	}
	log.Infof("白灯 %s (%s) 开关状态: %s", light.DevKEY, light.ProdCode, status)

	ligthCTRL := devBaseInfo.GetLightCtl()
	if ligthCTRL != nil {
		// lum
		lum := ligthCTRL.GetLum()
		if light.lum != lum {
			errMsg = append(errMsg, fmt.Sprintf("assert lum fail, exp: %d, act: %d", light.lum, lum))
		}
		log.Infof("白灯 %s (%s) 亮度: %d", light.DevKEY, light.ProdCode, lum)

		// 只有 T1012 和 T1013 才有色温
		if light.ProdCode != "T1011" {
			colortemp := ligthCTRL.GetColorTemp()
			if light.colorTemp != colortemp {
				errMsg = append(errMsg, fmt.Sprintf("assert colorTemp fail, exp: %d, act: %d", light.colorTemp, colortemp))
			}
			log.Infof("白灯 %s (%s) 色温: %d", light.DevKEY, light.ProdCode, colortemp)
		}
	}

	// product code, device key, test case, test time.
	contents := []string{light.ProdCode, light.DevKEY, light.testcase, commontool.GetCurrentTime()}

	if errMsg != nil && len(errMsg) > 0 {
		if light.notPassAndwaitNextHeartBeat < 3 {
			light.notPassAndwaitNextHeartBeat++
			return
		}
		contents = append(contents, "Fail")
		contents = append(contents, errMsg...)
	} else {
		contents = append(contents, "Pass")
	}

	result.WriteToExcel(contents)

	// 重置
	light.IsCmdSent = false
	light.notPassAndwaitNextHeartBeat = 0

}

// 解析服务器消息
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
		}
	}

}
