package device

import (
	"math"
	"math/rand"
	"oceanwing/commontool"

	log "github.com/cihub/seelog"
	"github.com/golang/protobuf/proto"
	lightEvent "oceanwing/eufy/protobuf.lib/light/lightevent"
	light1013 "oceanwing/eufy/protobuf.lib/light/t1013"
)

// LightWithColor 是对产品 T1013、T1604 的描述
type LightWithColor struct {
	baseDevice
	stopCtrlFunc chan struct{}
}

// NewLightWithColor create a new color light instance.
func NewLightWithColor(prodCode, devKey, devid string) EufyDevice {
	o := &LightWithColor{}
	o.ProdCode = prodCode
	o.DevKEY = devKey
	o.DevID = devid
	o.PubTopicl = "DEVICE/" + prodCode + "/" + devKey + "/SUB_MESSAGE"
	o.SubTopicl = "DEVICE/" + prodCode + "/" + devKey + "/PUH_MESSAGE"
	o.DeviceMsg = make(chan []byte)
	o.ServerMsg = make(chan []byte)
	o.stopCtrlFunc = make(chan struct{})
	log.Infof("Create a color Light, product code: %s, device key: %s, device id: %s", prodCode, devKey, devid)
	return o
}

// HandleSubscribeMessage 实现 EufyDevice 接口
func (light *LightWithColor) HandleSubscribeMessage() {
	go func() {
		log.Debugf("Running handleIncomingMsg function for device: %s", light.DevKEY)
		for {
			select {
			case devmsg := <-light.DeviceMsg:
				log.Infof("get new incoming message from device: %s", light.DevKEY)
				light.unMarshalHeartBeatMessage(devmsg)
			case serMsg := <-light.ServerMsg:
				log.Info("get new incoming message from server")
				light.unMarshalServerMessage(serMsg)
			}
		}
	}()
}

// BuildProtoBufMessage 实现 EufyDevice 接口
func (light *LightWithColor) BuildProtoBufMessage() []byte {
	var serMsg *light1013.ServerMessage
	serMsg = light.setWhiteLightMode()

	data, err := proto.Marshal(serMsg)
	if err != nil {
		log.Errorf("build protobuf message fail: %s", err)
		return nil
	}
	return data
}

// 白色模式， 和 T1012 一样的
func (light *LightWithColor) setWhiteLightMode() *light1013.ServerMessage {
	whiteLight := &light1013.ServerMessage_SetLightData_{}
	lightData := &light1013.ServerMessage_SetLightData{}
	lightData.Type = light1013.CmdType_REMOTE_SET_LIGHTING_PARA.Enum()
	lightData.Mode = light1013.LIGHT_DEV_MODE_WHITE_LIGHT_MODE.Enum()
	// lightData.White = &lightEvent.LampLightLevelCtlMessage{
	// 	Lum:       proto.Uint32(88),
	// 	ColorTemp: proto.Uint32(88),
	// }

	white := &lightEvent.LampLightLevelCtlMessage{}
	brightness := uint32(commontool.RandInt64(1, 101))
	color := uint32(commontool.RandInt64(1, 101))
	white.Lum = proto.Uint32(brightness)
	white.ColorTemp = proto.Uint32(color)

	lightData.White = white

	whiteLight.SetLightData = lightData

	serMsg := &light1013.ServerMessage{
		SessionId:     proto.Int32(rand.Int31n(math.MaxInt32)),
		RemoteMessage: whiteLight,
	}

	return serMsg
}

// 彩光模式
func (light *LightWithColor) setColorLightMode() *light1013.ServerMessage {
	red := uint32(commontool.RandInt64(0, 255))
	green := uint32(commontool.RandInt64(0, 255))
	blue := uint32(commontool.RandInt64(0, 255))

	colorlight := &light1013.ServerMessage_SetLightData_{
		SetLightData: &light1013.ServerMessage_SetLightData{
			Type: light1013.CmdType_REMOTE_SET_LIGHTING_PARA.Enum(),
			Mode: light1013.LIGHT_DEV_MODE_COLOR_LIGHT_MODE.Enum(),
			Rgb: &lightEvent.LampLightRgbCtlMessage{
				Red:   proto.Uint32(red),
				Green: proto.Uint32(green),
				Blue:  proto.Uint32(blue),
				White: proto.Uint32(80),
			},
		},
	}

	serMsg := &light1013.ServerMessage{
		SessionId:     proto.Int32(rand.Int31n(math.MaxInt32)),
		RemoteMessage: colorlight,
	}
	log.Infof("设置彩灯, Red: %d, Green: %d, Blue: %d", red, green, blue)
	return serMsg
}

/*--------------------------------------------------------------------------------------------------------------------------------------*/

//   解析设备心跳消息
func (light *LightWithColor) unMarshalHeartBeatMessage(payload []byte) {
	devMsg := &light1013.DeviceMessage{}
	err := proto.Unmarshal(payload, devMsg)
	if err != nil {
		log.Errorf("unmarshal device heart beat message fail: %s", err)
		return
	}

	// Non_ParamMsg
	noneParaMsg := devMsg.GetNonParaMsg()
	if noneParaMsg != nil {
		log.Infof("无参数消息, 指令类型: %s", noneParaMsg.GetType().String())
	}

	// devBaseInfo
	devBaseInfo := devMsg.GetReportDevBaseInfo()
	if devBaseInfo == nil {
		return
	}

	cmdtype := devBaseInfo.GetType()
	log.Infof("彩灯 %s (%s) CmdType: %s", light.DevKEY, light.ProdCode, cmdtype.String())

	mode := devBaseInfo.GetMode()
	log.Infof("彩灯 %s (%s) Mode: %s", light.DevKEY, light.ProdCode, mode.String())

	status := devBaseInfo.GetOnoffStatus()
	log.Infof("彩灯 %s (%s) Mode: %s", light.DevKEY, light.ProdCode, status.String())

	leaveHomeState := devBaseInfo.GetLeaveHomeState()
	log.Infof("彩灯 %s (%s) leaveHomeState: %v", leaveHomeState)

	// devBaseInfo --> white
	white := devBaseInfo.GetWhite()
	if white != nil {
		lum := white.GetLum()
		colorTemp := white.GetColorTemp()
		log.Infof("彩灯 %s (%s) Lum: %d, ColorTemp: %d", light.DevKEY, light.ProdCode, lum, colorTemp)
	}

	// devBaseInfo --> rgb
	rgb := devBaseInfo.GetRgb()
	if rgb != nil {
		red := rgb.GetRed()
		green := rgb.GetGreen()
		blue := rgb.GetBlue()
		lum := rgb.GetWhite()
		log.Infof("彩灯 %s (%s) Color light 模式, Red: %d, Green: %d, Blue: %d, Lum: %d", light.DevKEY, light.ProdCode, red, green, blue, lum)
	}

	// devBaseInfo --> streamLight
	streamLight := devBaseInfo.GetStreamLight()
	if streamLight != nil {
		ti := streamLight.GetTime()
		brightnessper := streamLight.GetBrightnessPercent()
		log.Infof("彩灯 %s (%s) 流光模式, 流光速度(秒): %d, 亮度: %d", light.DevKEY, light.ProdCode, ti, brightnessper)

		point := streamLight.GetPoint()
		if point != nil {
			pointA := point.GetPointA()
			pointB := point.GetPointB()
			pointC := point.GetPointC()
			pointD := point.GetPointD()
			log.Infof("彩灯 %s (%s) 流光模式, Point A, Red: %d, Green: %d, Blue: %d", light.DevKEY, light.ProdCode, pointA.GetRed(), pointA.GetGreen(), pointA.GetBlue())
			log.Infof("彩灯 %s (%s) 流光模式, Point B, Red: %d, Green: %d, Blue: %d", light.DevKEY, light.ProdCode, pointB.GetRed(), pointB.GetGreen(), pointB.GetBlue())
			log.Infof("彩灯 %s (%s) 流光模式, Point C, Red: %d, Green: %d, Blue: %d", light.DevKEY, light.ProdCode, pointC.GetRed(), pointC.GetGreen(), pointC.GetBlue())
			log.Infof("彩灯 %s (%s) 流光模式, Point D, Red: %d, Green: %d, Blue: %d", light.DevKEY, light.ProdCode, pointD.GetRed(), pointD.GetGreen(), pointD.GetBlue())
		}
	}
}

// 解析服务器控制消息
func (light *LightWithColor) unMarshalServerMessage(payload []byte) {
	serMsg := &light1013.ServerMessage{}
	err := proto.Unmarshal(payload, serMsg)
	if err != nil {
		log.Errorf("unmarshal server message fail: %s", err)
		return
	}

	// session id
	log.Infof("Session ID: %d", serMsg.GetSessionId())

	// Sync_Time_Alarm
	sta := serMsg.GetSync_Time_Alarm()
	if sta != nil {
		log.Info("解析出服务器同步的时间和闹钟信息")

		cmdtype := sta.GetType()
		log.Infof("----CmdType: %s", cmdtype.String())

		syncTime := sta.GetTime()
		if syncTime != nil {
			log.Info("---------------同步时间:")
			log.Infof("-------年: %d", syncTime.GetYear())
			log.Infof("-------月: %d", syncTime.GetMonth()+1) //月份的定义范围是0--11，所以要加1
			log.Infof("-------日: %d", syncTime.GetDay())
			log.Infof("-------Weekday: %d", syncTime.GetWeekday())
			log.Infof("-------时: %d", syncTime.GetHours())
			log.Infof("-------分: %d", syncTime.GetMinutes())
			log.Infof("-------秒: %d", syncTime.GetSeconds())
		}

		syncAlert := sta.GetAlarm()
		if syncAlert != nil {
			records := syncAlert.GetAlarmRecordData()
			if records != nil {
				log.Info("--------------同步闹钟:")
				for i, rec := range records {
					log.Infof("-------第 %d 个闹钟信息", i+1)
					alertevent := rec.GetAlarmEvent()
					eventType := alertevent.GetType()
					ctrlMsg := alertevent.GetLightCtl()
					log.Infof("-------闹钟类型: %s", eventType.String())
					log.Infof("-------闹钟灯的设置, 亮度: %d, 色温: %d", ctrlMsg.GetLum(), ctrlMsg.GetColorTemp())

					alertmsg := rec.GetAlarmMesage()
					log.Infof("-------闹钟时间--> %d:%d:%d:", alertmsg.GetHours(), alertmsg.GetMinutes(), alertmsg.GetSeconds())
					log.Infof("-------是否重复: %t", alertmsg.GetRepetiton())
					weekinfo := commontool.ConvertToWeekDay(alertmsg.GetWeekInfo())
					log.Infof("-------weekinfo: %s", weekinfo)
				}
			}
		}
	}

	// SetLightData
	lightdata := serMsg.GetSetLightData()
	if lightdata != nil {
		log.Info("解析出服务器控制指令信息")

		cmdtype := lightdata.GetType()
		log.Infof("----CmdType: %s", cmdtype.String())

		mode := lightdata.GetMode()
		log.Infof("----Mode: %s", mode.String())

		status := lightdata.GetOnoffStatus()
		log.Infof("----Status: %s", status.String())

		white := lightdata.GetWhite()
		if white != nil {
			log.Infof("----白光模式，亮度: %d, 色温: %d", white.GetLum(), white.GetColorTemp())
		}

		rgb := lightdata.GetRgb()
		if rgb != nil {
			log.Infof("----彩光模式, 红: %d, 绿: %d, 蓝: %d, 亮度: %d", rgb.GetRed(), rgb.GetGreen(), rgb.GetBlue(), rgb.GetWhite())
		}

		stream := lightdata.GetStreamLight()
		if stream != nil {
			log.Info("----流光模式----")
			log.Infof("------时间间隔(秒): %d", stream.GetTime())
			log.Infof("------亮度: %d", stream.GetBrightnessPercent())

			points := stream.GetPoint()
			pointA := points.GetPointA()
			pointB := points.GetPointB()
			pointC := points.GetPointC()
			pointD := points.GetPointD()

			log.Infof("------ PointA, 红: %d, 绿: %d, 蓝: %d", pointA.GetGreen(), pointA.GetGreen(), pointA.GetBlue())
			log.Infof("------ pointB, 红: %d, 绿: %d, 蓝: %d", pointB.GetGreen(), pointB.GetGreen(), pointB.GetBlue())
			log.Infof("------ pointC, 红: %d, 绿: %d, 蓝: %d", pointC.GetGreen(), pointC.GetGreen(), pointC.GetBlue())
			log.Infof("------ pointD, 红: %d, 绿: %d, 蓝: %d", pointD.GetGreen(), pointD.GetGreen(), pointD.GetBlue())
		}
	}

	// SetAwayMode
	away := serMsg.GetSetAwayMode_Status()
	if away != nil {
		log.Info("解析出离家模式信息")

		cmdtype := away.GetType()
		log.Infof("----CmdType: %s", cmdtype.String())

		leaveMsg := away.GetSyncLeaveModeMsg()
		log.Infof("----开始时间 %d:%d", leaveMsg.GetStartHours(), leaveMsg.GetStartMinutes())
		log.Infof("----结束时间 %d:%d", leaveMsg.GetFinishHours(), leaveMsg.GetFinishMinutes())
		log.Infof("----是否重复: %t", leaveMsg.GetRepetiton())
		log.Infof("----是否开启: %t", leaveMsg.GetLeaveHomeState())
		log.Infof("----WeekInfo: %s", commontool.ConvertToWeekDay(leaveMsg.GetWeekInfo()))
	}

	// SetPowerUpLightStatus
	powerUpData := serMsg.GetSetPowerupLightStatus()
	if powerUpData != nil {
		log.Info("解析出初始上电时服务器的同步信息")

		cmdtype := powerUpData.GetType()
		log.Infof("----CmdType: %s", cmdtype.String())

		mode := powerUpData.GetMode()
		log.Infof("----Mode: %s", mode.String())

		powerStatus := powerUpData.GetPowrupStatus()
		log.Infof("----powerStatus: %s", powerStatus.String())

		white := powerUpData.GetWhite()
		if white != nil {
			log.Infof("----白光， 亮度: %d, 色温: %d", white.GetLum(), white.GetColorTemp())
		}

		rgb := powerUpData.GetRgbw()
		if rgb != nil {
			log.Infof("----彩光, 红: %d, 绿: %d, 蓝: %d, 亮度: %d", rgb.GetRed(), rgb.GetGreen(), rgb.GetBlue(), rgb.GetWhite())
		}

	}

}
