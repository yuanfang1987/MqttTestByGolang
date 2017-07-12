package device

import (
	"math"
	"math/rand"
	"oceanwing/commontool"
	"strconv"
	"strings"

	log "github.com/cihub/seelog"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/protobuf/proto"
	lightEvent "oceanwing/eufy/protobuf.lib/light/lightevent"
	light1013 "oceanwing/eufy/protobuf.lib/light/t1013"
)

var (
	rgb               []map[string]interface{}
	rgbNum            int
	colorModCaseIndex = 0
)

// LightWithColor 是对产品 T1013、T1604 的描述
type LightWithColor struct {
	baseDevice
	stopCtrlFunc    chan struct{}
	rgbMap          map[string]*rgbInfo
	lum             uint32
	colorTemp       uint32
	bright          uint32
	onOffStatus     uint32
	streamModeSpeed int
	onOffStatChan   chan string
	awayModTesting  bool
	resultMap       map[string]string
	occurSlice      []map[string]string
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
	o.SubServerTopic = "DEVICE/" + prodCode + "/" + devKey + "/SUB_MESSAGE"
	o.SubDeviceTopic = "DEVICE/" + prodCode + "/" + devKey + "/PUH_MESSAGE"
	o.SubMessage = make(chan MQTT.Message)
	o.stopCtrlFunc = make(chan struct{})
	o.rgbMap = make(map[string]*rgbInfo)
	o.resultMap = make(map[string]string)
	o.onOffStatChan = make(chan string, 2)
	log.Infof("Create a color Light, product code: %s, device key: %s", prodCode, devKey)
	// 读取RGB配色信息
	// getRGBData()
	return o
}

// HandleSubscribeMessage 实现 EufyDevice 接口
func (light *LightWithColor) HandleSubscribeMessage() {
	go func() {
		log.Debugf("Running handleIncomingMsg function for device: %s", light.DevKEY)
		for {
			select {
			case msg := <-light.SubMessage:
				go light.unMarshalAllMsg(msg)
			}
		}
	}()
}

// BuildProtoBufMessage hh..
func (light *LightWithColor) BuildProtoBufMessage() []byte {
	var serMsg *light1013.ServerMessage
	serMsg = light.setDataForLight(1)

	data, err := proto.Marshal(serMsg)
	if err != nil {
		log.Errorf("build protobuf message fail: %s", err)
		return nil
	}
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
		lightdata.Mode = light1013.LIGHT_DEV_MODE_WHITE_LIGHT_MODE.Enum()
		lightdata.White = light.buildWhiteLightData()
	case 1:
		// 彩光模式
		lightdata.Mode = light1013.LIGHT_DEV_MODE_COLOR_LIGHT_MODE.Enum()
		lightdata.Rgb = light.buildColorLightData()
	case 2:
		// 流光模式
		lightdata.Mode = light1013.LIGHT_DEV_MODE_STREAMER_LIGHT_MODE.Enum()
		lightdata.StreamLight = light.buildStreamLightData()
	case 4:
		// 开关灯
		lightdata.OnoffStatus = light.setLightOnOffStatus()
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
		light.rgbMap["colorModRGB"] = expRGB
	} else {
		red = uint32(commontool.RandInt64(0, 255))
		green = uint32(commontool.RandInt64(0, 255))
		blue = uint32(commontool.RandInt64(0, 255))
		// 储存当前数据，用于后续心跳判断
		rrggbb := &rgbInfo{
			red:   red,
			green: green,
			blue:  blue,
		}
		light.rgbMap["colorModRGB"] = rrggbb
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
	// 速度， 固定 1 秒， 不好吧
	light.streamModeSpeed = 1

	points := &light1013.ColorPointMessage{
		PointA: light.buildRGBData("A"),
		PointB: light.buildRGBData("B"),
		PointC: light.buildRGBData("C"),
		PointD: light.buildRGBData("D"),
	}

	stream := &light1013.StreamLight{}

	stream.Time = proto.Int32(1)
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
	if light.onOffStatus != 1 {
		log.Info("开灯")
		return light1013.LIGHT_ONOFF_STATUS_ON.Enum()
	}
	log.Info("关灯")
	return light1013.LIGHT_ONOFF_STATUS_OFF.Enum()
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
	log.Infof("彩灯 %s (%s) Status: %s", light.DevKEY, light.ProdCode, status.String())

	leaveHomeState := devBaseInfo.GetLeaveHomeState()
	log.Infof("彩灯 %s (%s) leaveHomeState: %t", light.DevKEY, light.ProdCode, leaveHomeState)

	// 把当前心跳的 status 存入 channel, 最多可以存 2 个
	light.onOffStatChan <- status.String()

	// 判断离家模式
	if leaveHomeState {
		// 标记一下已经处于离家模式测试中
		light.awayModTesting = true

		// 记录由 normal mode 转为 leave mode 的时间, 如果字典中尚未有记录，则记录之
		if _, ok := light.resultMap[LeaveModeUp]; !ok {
			light.resultMap[LeaveModeUp] = commontool.GetCurrentTime()
			log.Infof("彩灯 %s (%s) 已启动离家模式...", light.DevKEY, light.ProdCode)
		}

		var prev, current string

		// 取出上次心跳和本次心跳的开关状态
		if len(light.onOffStatChan) == 2 {
			prev = <-light.onOffStatChan
			current = <-light.onOffStatChan
		}

		// 记录开启离家模式之前的那个状态
		if _, ok := light.resultMap[StatusBeforeLeaveMode]; !ok {
			if prev != "" {
				light.resultMap[StatusBeforeLeaveMode] = prev
				log.Infof("彩灯 %s (%s) 启动离家模式之前的状态是: %s", light.DevKEY, light.ProdCode, prev)
			}
		}

		// 如果两次的状态不同，则说明状态发生了变化，随机开关灯被触发, 把操作类型和时间记录下来
		if prev != current {
			leaveModeOccur := make(map[string]string)
			leaveModeOccur[OccurTime] = commontool.GetCurrentTime()
			leaveModeOccur[OccurType] = current
			light.occurSlice = append(light.occurSlice, leaveModeOccur)
			log.Infof("彩灯 %s (%s) 随机开关为被触发, 本次是: %s", light.DevKEY, light.ProdCode, current)
		}

		// 把本次心跳的状态重新存入 channel 中，不然 channel 被掏空了，下次就没法一起取出两个来比较了
		if len(light.onOffStatChan) < 2 {
			light.onOffStatChan <- current
		}

	} else if light.awayModTesting {
		// 记录离家模式的失效时间
		if _, ok := light.resultMap[LeaveModeDown]; !ok {
			light.resultMap[LeaveModeDown] = commontool.GetCurrentTime()
		}

		// 记录恢复正常模式后的状态
		if _, ok := light.resultMap[StatusResumeToNormal]; !ok {
			light.resultMap[StatusResumeToNormal] = status.String()
		}

		light.awayModTesting = false
	}

	// 如果 channel 中缓存已满，必须取出一个, 否则下次心跳再往里面发数据的时候程序会死掉
	if len(light.onOffStatChan) == 2 {
		<-light.onOffStatChan
	}

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
		starthour := int(leaveMsg.GetStartHours())
		startminute := int(leaveMsg.GetStartMinutes())
		finishhour := int(leaveMsg.GetFinishHours())
		finishminute := int(leaveMsg.GetFinishMinutes())
		leaveHomeFlag := leaveMsg.GetLeaveHomeState()

		log.Infof("----开始时间 %d:%d", starthour, startminute)
		log.Infof("----结束时间 %d:%d", finishhour, finishminute)
		log.Infof("----是否重复: %t", leaveMsg.GetRepetiton())
		log.Infof("----是否开启: %t", leaveHomeFlag)
		log.Infof("----WeekInfo: %s", commontool.ConvertToWeekDay(leaveMsg.GetWeekInfo()))

		// 如果服务器下发的 离家模式 为 true，则记下时间
		if leaveHomeFlag {
			light.resultMap[LeaveModeUpExp] = strconv.Itoa(starthour) + ":" + strconv.Itoa(startminute)
			light.resultMap[LeaveModeDownExp] = strconv.Itoa(finishhour) + ":" + strconv.Itoa(finishminute)
		}

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

func (light *LightWithColor) unMarshalAllMsg(msg MQTT.Message) {
	t := msg.Topic()
	payload := msg.Payload()

	if light.SubDeviceTopic == t {
		// 设备心跳消息
		log.Info("----- 这是一个来自彩灯设备的心跳消息----------")
		light.unMarshalHeartBeatMessage(payload)
	} else if light.SubServerTopic == t {
		//服务器消息
		log.Info("-------这是一个来自服务器的控制消息---------")
		light.unMarshalServerMessage(payload)
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
