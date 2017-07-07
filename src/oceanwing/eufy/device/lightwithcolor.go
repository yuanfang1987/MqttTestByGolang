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
	serMsg = light.setColorLightMode()

	data, err := proto.Marshal(serMsg)
	if err != nil {
		log.Errorf("build protobuf message fail: %s", err)
		return nil
	}
	return data
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
	log.Infof("彩灯 %s (%s) leaveHomeState: %t", leaveHomeState)

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

	// SetLightData

	// SetAwayMode

}
