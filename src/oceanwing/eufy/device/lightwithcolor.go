package device

import (
	"fmt"
	"math"
	"math/rand"
	"oceanwing/commontool"
	"oceanwing/eufy/result"
	"strconv"
	"strings"

	"time"

	log "github.com/cihub/seelog"
	"github.com/golang/protobuf/proto"
	lightEvent "oceanwing/eufy/protobuf.lib/light/lightevent"
	light1013 "oceanwing/eufy/protobuf.lib/light/t1013"
)

const (
	colorModRGB = "colorModRGB"
)

var (
	rgb               []map[string]interface{}
	rgbNum            int
	colorModCaseIndex = 0
	testMode          = 0
)

// LightWithColor 是对产品 T1013、T1604 的描述
type LightWithColor struct {
	baseDevice
	stopCtrlFunc    chan struct{}
	rgbMap          map[string]*rgbInfo
	lum             uint32
	colorTemp       uint32
	bright          uint32
	onOffStatus     light1013.LIGHT_ONOFF_STATUS
	streamModeSpeed int32
	mod             light1013.LIGHT_DEV_MODE
}

// RGB 配色信息
type rgbInfo struct {
	red   uint32
	green uint32
	blue  uint32
}

// NewLightWithColor create a new color light instance.
func NewLightWithColor(prodCode, devKey string) EufyDevice {
	o := &LightWithColor{}
	o.ProdCode = prodCode
	o.DevKEY = devKey
	o.PubTopicl = "DEVICE/" + prodCode + "/" + devKey + "/SUB_MESSAGE"
	o.SubTopicl = "DEVICE/" + prodCode + "/" + devKey + "/PUH_MESSAGE"
	o.DeviceMsg = make(chan []byte)
	o.ServerMsg = make(chan []byte)
	o.stopCtrlFunc = make(chan struct{})
	o.rgbMap = make(map[string]*rgbInfo)
	o.onOffStatus = light1013.LIGHT_ONOFF_STATUS_ON // 测试前必须要确保灯是开的
	log.Infof("Create a color Light, product code: %s, device key: %s", prodCode, devKey)
	controlRunMode()
	// 读取RGB配色信息
	getRGBData()
	return o
}

// HandleSubscribeMessage 实现 EufyDevice 接口
func (light *LightWithColor) HandleSubscribeMessage() {
	go func() {
		log.Debugf("Running handleIncomingMsg function for device: %s", light.DevKEY)
		for {
			select {
			case devmsg := <-light.DeviceMsg:
				log.Infof("======设备上报消息: %s======", light.DevKEY)
				go light.unMarshalHeartBeatMessage(devmsg)
			case serMsg := <-light.ServerMsg:
				log.Info("======服务器控制消息======")
				go light.unMarshalServerMessage(serMsg)
			}
		}
	}()
}

// BuildProtoBufMessage 实现 EufyDevice 接口
func (light *LightWithColor) BuildProtoBufMessage() []byte {
	if light.notPassAndwaitNextHeartBeat != 0 {
		log.Warnf("上次心跳验证未通过，等待下次验证，当前已验证 %d 次", light.notPassAndwaitNextHeartBeat)
		return nil
	}

	var serMsg *light1013.ServerMessage
	serMsg = light.setDataForLight(testMode)

	data, err := proto.Marshal(serMsg)
	if err != nil {
		log.Errorf("build protobuf message fail: %s", err)
		return nil
	}

	log.Info("=================================================")

	// 标记一下
	light.IsCmdSent = true
	return data
}

// 控制灯光变化
func (light *LightWithColor) setDataForLight(mode int) *light1013.ServerMessage {
	setlightdata := &light1013.ServerMessage_SetLightData_{}

	lightdata := &light1013.ServerMessage_SetLightData{}
	lightdata.Type = light1013.CmdType_REMOTE_SET_LIGHTING_PARA.Enum()

	switch mode {
	case 0:
		// 白光模式
		light.mod = light1013.LIGHT_DEV_MODE_WHITE_LIGHT_MODE
		light.testcase = "白光模式"
		lightdata.Mode = light1013.LIGHT_DEV_MODE_WHITE_LIGHT_MODE.Enum()
		lightdata.White = light.buildWhiteLightData()
	case 1:
		// 彩光模式
		light.mod = light1013.LIGHT_DEV_MODE_COLOR_LIGHT_MODE
		light.testcase = "彩光模式"
		lightdata.Mode = light1013.LIGHT_DEV_MODE_COLOR_LIGHT_MODE.Enum()
		lightdata.Rgb = light.buildColorLightData()
	case 2:
		// 流光模式
		light.mod = light1013.LIGHT_DEV_MODE_STREAMER_LIGHT_MODE
		light.testcase = "流光模式"
		lightdata.Mode = light1013.LIGHT_DEV_MODE_STREAMER_LIGHT_MODE.Enum()
		lightdata.StreamLight = light.buildStreamLightData()
	case 3:
		// 开关灯
		light.testcase = "开关灯"
		lightdata.OnoffStatus = light.setLightOnOffStatus()
	default:
		log.Warn("警告，没有指定测试模式")
		return nil
	}

	setlightdata.SetLightData = lightdata
	serMsg := &light1013.ServerMessage{
		SessionId:     proto.Int32(rand.Int31n(math.MaxInt32)),
		RemoteMessage: setlightdata,
	}

	return serMsg
}

// 构造白光模式的数据，用随机数来产生亮度和色温
func (light *LightWithColor) buildWhiteLightData() *lightEvent.LampLightLevelCtlMessage {
	brightness := uint32(commontool.RandInt64(1, 101))
	color := uint32(commontool.RandInt64(1, 101))

	// 储存当前数据，用于后续心跳判断
	light.lum = brightness
	light.colorTemp = color

	return &lightEvent.LampLightLevelCtlMessage{
		Lum:       proto.Uint32(brightness),
		ColorTemp: proto.Uint32(color),
	}
}

// 构造彩光模式的数据
func (light *LightWithColor) buildColorLightData() *lightEvent.LampLightRgbCtlMessage {
	var red, green, blue uint32
	brightness := uint32(commontool.RandInt64(10, 100))
	// 储存当前数据，用于后续心跳判断
	light.bright = brightness

	colorData := rgb[colorModCaseIndex]

	colorModCaseIndex++
	if colorModCaseIndex >= rgbNum {
		colorModCaseIndex = 0
	}

	name := colorData["RGBName"]
	v := colorData["rgbConfig"]
	if expRGB, ok := v.(*rgbInfo); ok {
		red = expRGB.red
		green = expRGB.green
		blue = expRGB.blue
		// 储存当前数据，用于后续心跳判断
		light.rgbMap[colorModRGB] = expRGB
	} else {
		red = uint32(commontool.RandInt64(0, 255))
		green = uint32(commontool.RandInt64(0, 255))
		blue = uint32(commontool.RandInt64(0, 255))
		rrggbb := &rgbInfo{
			red:   red,
			green: green,
			blue:  blue,
		}
		// 储存当前数据，用于后续心跳判断
		light.rgbMap[colorModRGB] = rrggbb
	}

	log.Infof("设置彩光模式, 彩光名称: %s, 红: %d, 绿: %d, 蓝: %d, 亮度: %d", name, red, green, blue, brightness)

	return &lightEvent.LampLightRgbCtlMessage{
		Red:   proto.Uint32(red),
		Green: proto.Uint32(green),
		Blue:  proto.Uint32(blue),
		White: proto.Uint32(brightness),
	}
}

// 构造流光模式的数据
func (light *LightWithColor) buildStreamLightData() *light1013.StreamLight {
	brightness := int32(commontool.RandInt64(20, 100))
	// 储存当前数据，用于后续心跳判断
	light.bright = uint32(brightness)
	// 速度
	speed := int32(commontool.RandInt64(1, 3))
	light.streamModeSpeed = speed

	points := &light1013.ColorPointMessage{
		PointA: light.buildRGBData("A"),
		PointB: light.buildRGBData("B"),
		PointC: light.buildRGBData("C"),
		PointD: light.buildRGBData("D"),
	}

	stream := &light1013.StreamLight{}

	stream.Time = proto.Int32(speed)
	stream.BrightnessPercent = proto.Int32(brightness)
	stream.Point = points

	return stream
}

// RGB 数值
func (light *LightWithColor) buildRGBData(p string) *light1013.RGBMessage {
	var red, green, blue int32
	var seed int64
	switch p {
	case "A":
		seed = commontool.RandInt64(0, 12)
	case "B":
		seed = commontool.RandInt64(12, 24)
	case "C":
		seed = commontool.RandInt64(24, 36)
	case "D":
		seed = commontool.RandInt64(36, 51)
	}

	if int(seed) < rgbNum {
		rrggbb := rgb[seed]
		rgbcon := rrggbb["rgbConfig"]
		if rgbdata, ok := rgbcon.(*rgbInfo); ok {
			red = int32(rgbdata.red)
			green = int32(rgbdata.green)
			blue = int32(rgbdata.blue)
			// 储存当前数据，用于后续心跳判断
			light.rgbMap[p] = rgbdata
		}
	}

	return &light1013.RGBMessage{
		Red:   proto.Int32(red),
		Green: proto.Int32(green),
		Blue:  proto.Int32(blue),
	}
}

// 开关灯
func (light *LightWithColor) setLightOnOffStatus() *light1013.LIGHT_ONOFF_STATUS {
	if light.onOffStatus != light1013.LIGHT_ONOFF_STATUS_ON {
		log.Info("开灯")
		light.onOffStatus = light1013.LIGHT_ONOFF_STATUS_ON
		return light1013.LIGHT_ONOFF_STATUS_ON.Enum()
	}
	log.Info("关灯")
	light.onOffStatus = light1013.LIGHT_ONOFF_STATUS_OFF
	return light1013.LIGHT_ONOFF_STATUS_OFF.Enum()
}

/*--------------------------------------------------------------------------------------------------------------------------------------*/

//   解析设备心跳消息
func (light *LightWithColor) unMarshalHeartBeatMessage(payload []byte) {
	devMsg := &light1013.DeviceMessage{}
	err := proto.Unmarshal(payload, devMsg)
	if err != nil {
		log.Errorf("解析 %s (%s) 心跳消息失败: %s", light.DevKEY, light.ProdCode, err)
		return
	}

	// Non_ParamMsg
	noneParaMsg := devMsg.GetNonParaMsg()
	if noneParaMsg != nil {
		log.Infof("彩灯 %s (%s) 无参数消息, 指令类型: %s", light.DevKEY, light.ProdCode, noneParaMsg.GetType().String())
	}

	// devBaseInfo
	devBaseInfo := devMsg.GetReportDevBaseInfo()
	if devBaseInfo == nil {
		return
	}

	var errMsg []string

	// CmdType
	cmdtype := devBaseInfo.GetType()
	if cmdtype != light1013.CmdType_DEV_REPORT_STATUS {
		errMsg = append(errMsg, fmt.Sprintf("assert CmdType fail, exp: %s, act: %s", light1013.CmdType_DEV_REPORT_STATUS.String(), cmdtype.String()))
	}
	log.Infof("彩灯 %s (%s) CmdType: %s", light.DevKEY, light.ProdCode, cmdtype.String())

	// Mode
	mode := devBaseInfo.GetMode()
	if light.mod != mode {
		errMsg = append(errMsg, fmt.Sprintf("assert mode fail, exp: %s, act: %s", light.mod.String(), mode.String()))
	}
	log.Infof("彩灯 %s (%s) Mode: %s", light.DevKEY, light.ProdCode, mode.String())

	// Status
	status := devBaseInfo.GetOnoffStatus()
	if light.onOffStatus != status {
		errMsg = append(errMsg, fmt.Sprintf("assert onOff status fail, exp: %s, act: %s", light.onOffStatus.String(), status.String()))
	}
	log.Infof("彩灯 %s (%s) status: %s", light.DevKEY, light.ProdCode, status.String())

	// leaveHomeState, 这个就不要判断了， 会有专门的测试
	leaveHomeState := devBaseInfo.GetLeaveHomeState()
	log.Infof("彩灯 %s (%s) leaveHomeState: %t", light.DevKEY, light.ProdCode, leaveHomeState)

	// devBaseInfo --> white
	white := devBaseInfo.GetWhite()
	if white != nil {
		lum := white.GetLum()
		colorTemp := white.GetColorTemp()
		if light.lum != lum {
			errMsg = append(errMsg, fmt.Sprintf("assert lum fail, exp: %d, act: %d", light.lum, lum))
		}
		if light.colorTemp != colorTemp {
			errMsg = append(errMsg, fmt.Sprintf("assert colorTemp fail, exp: %d, act: %d", light.colorTemp, colorTemp))
		}
		log.Infof("彩灯 %s (%s) Lum: %d, ColorTemp: %d", light.DevKEY, light.ProdCode, lum, colorTemp)
	}

	// devBaseInfo --> rgb
	rgb := devBaseInfo.GetRgb()
	if rgb != nil {
		red := rgb.GetRed()
		green := rgb.GetGreen()
		blue := rgb.GetBlue()
		lum := rgb.GetWhite()

		// 如果存在预期结果，则判断之
		if expRGB, ok := light.rgbMap[colorModRGB]; ok {
			if expRGB.red != red {
				errMsg = append(errMsg, fmt.Sprintf("color mode, assert red, exp: %d, act: %d", expRGB.red, red))
			}
			if expRGB.green != green {
				errMsg = append(errMsg, fmt.Sprintf("color mode, assert green, exp: %d, act: %d", expRGB.green, green))
			}
			if expRGB.blue != blue {
				errMsg = append(errMsg, fmt.Sprintf("color mode, assert blue, exp: %d, act: %d", expRGB.blue, blue))
			}
			if light.bright != lum {
				errMsg = append(errMsg, fmt.Sprintf("color mode, assert brightness, exp: %d, act: %d", light.bright, lum))
			}
		}

		log.Infof("彩灯 %s (%s) Color light 模式, Red: %d, Green: %d, Blue: %d, Lum: %d", light.DevKEY, light.ProdCode, red, green, blue, lum)
	}

	// devBaseInfo --> streamLight
	streamLight := devBaseInfo.GetStreamLight()
	if streamLight != nil {
		ti := streamLight.GetTime()
		brightnessper := streamLight.GetBrightnessPercent()
		if light.streamModeSpeed != ti {
			errMsg = append(errMsg, fmt.Sprintf("stream mode, assert time speed, exp: %d, act: %d", light.streamModeSpeed, ti))
		}
		if light.bright != uint32(brightnessper) {
			errMsg = append(errMsg, fmt.Sprintf("stream mode, assert brightness, exp: %d, act: %d", light.bright, brightnessper))
		}
		log.Infof("彩灯 %s (%s) 流光模式, 流光速度(秒): %d, 亮度: %d", light.DevKEY, light.ProdCode, ti, brightnessper)

		point := streamLight.GetPoint()
		if point != nil {
			pointA := point.GetPointA()
			pointB := point.GetPointB()
			pointC := point.GetPointC()
			pointD := point.GetPointD()

			allPoints := []*light1013.RGBMessage{pointA, pointB, pointC, pointD}
			expPointNames := []string{"A", "B", "C", "D"}

			for i, pointName := range expPointNames {
				if exPoint, ok := light.rgbMap[pointName]; ok {
					actPoint := allPoints[i]
					if exPoint.red != uint32(actPoint.GetRed()) {
						errMsg = append(errMsg, fmt.Sprintf("stream mode, point %s, assert red, exp: %d, act: %d", pointName, exPoint.red, actPoint.GetRed()))
					}
					if exPoint.green != uint32(actPoint.GetGreen()) {
						errMsg = append(errMsg, fmt.Sprintf("stream mode, point %s, assert green, exp: %d, act: %d", pointName, exPoint.green, actPoint.GetGreen()))
					}
					if exPoint.blue != uint32(actPoint.GetBlue()) {
						errMsg = append(errMsg, fmt.Sprintf("stream mode, point %s, assert blue, exp: %d, act: %d", pointName, exPoint.blue, actPoint.GetBlue()))
					}
				}
			}

			log.Infof("彩灯 %s (%s) 流光模式, Point A, Red: %d, Green: %d, Blue: %d", light.DevKEY, light.ProdCode, pointA.GetRed(), pointA.GetGreen(), pointA.GetBlue())
			log.Infof("彩灯 %s (%s) 流光模式, Point B, Red: %d, Green: %d, Blue: %d", light.DevKEY, light.ProdCode, pointB.GetRed(), pointB.GetGreen(), pointB.GetBlue())
			log.Infof("彩灯 %s (%s) 流光模式, Point C, Red: %d, Green: %d, Blue: %d", light.DevKEY, light.ProdCode, pointC.GetRed(), pointC.GetGreen(), pointC.GetBlue())
			log.Infof("彩灯 %s (%s) 流光模式, Point D, Red: %d, Green: %d, Blue: %d", light.DevKEY, light.ProdCode, pointD.GetRed(), pointD.GetGreen(), pointD.GetBlue())
		}
	}

	if !light.IsCmdSent {
		return
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
	// reset to 0
	light.notPassAndwaitNextHeartBeat = 0
	light.IsCmdSent = false

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

func getRGBData() {
	if rgb != nil || len(rgb) > 0 {
		log.Warn("rgb slice has been initizalize, don't do this again.")
		return
	}

	content, err := commontool.ReadFileContent("rgb.txt")
	if err != nil {
		log.Errorf("get rgb test data fail: %s", err)
		return
	}

	for _, line := range content {
		arrString := strings.Split(line, ",")

		r, _ := strconv.Atoi(arrString[0])
		g, _ := strconv.Atoi(arrString[1])
		b, _ := strconv.Atoi(arrString[2])
		rgnName := arrString[3]

		newRgb := &rgbInfo{
			red:   uint32(r),
			green: uint32(g),
			blue:  uint32(b),
		}

		rgbmap := make(map[string]interface{})
		rgbmap["RGBName"] = rgnName
		rgbmap["rgbConfig"] = newRgb
		rgb = append(rgb, rgbmap)
	}

	rgbNum = len(rgb)
}

func controlRunMode() {
	go func() {
		interval := time.NewTicker(time.Hour * 1).C
		for {
			select {
			case <-interval:
				testMode++
				if testMode > 3 {
					testMode = 0
				}
			}
		}
	}()
}
