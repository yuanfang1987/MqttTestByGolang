package device

import (
	"math"
	"math/rand"
	"oceanwingqa/common/utils"

	log "github.com/cihub/seelog"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/protobuf/proto"
	lightEvent "oceanwingqa/common/protobuf.lib/light/lightevent"
	light1013 "oceanwingqa/common/protobuf.lib/light/t1013"
)

// 这是一个描述彩灯的结构体，支持 T1013, T1604
type colorLight struct {
	baseEufy
	modeIndex int
}

// NewColorLight hj.
func newColorLight(clientid, username, pwd, broker, prodCode, devKey string, needCA bool) Eufydevice {
	c := &colorLight{}
	c.Clientid = clientid
	c.Username = username
	c.Pwd = pwd
	c.Broker = broker
	c.prod = prodCode
	c.PubTopic = "DEVICE/" + prodCode + "/" + devKey + "/PUH_MESSAGE"
	c.SubTopic = "DEVICE/" + prodCode + "/" + devKey + "/SUB_MESSAGE"
	c.NeedCA = needCA
	c.msgToServer = make(chan []byte, 2) //每次都忘记，死性不改！！！！！！！！
	c.msgFromServer = make(chan []byte, 2)
	return c
}

// 实现 Eufydevice 接口
func (c *colorLight) RunMqttService() {
	c.SubHandler = func(c MQTT.Client, msg MQTT.Message) {} // do nothing.
	c.MqttClient.ConnectToBroker()
	c.outgoing()
}

// 实现 Eufydevice 接口
func (c *colorLight) SendHeartBeat() {
	c.msgToServer <- c.buildHeartBeatMsg(c.modeIndex)
	c.modeIndex++
	if c.modeIndex > 2 {
		c.modeIndex = 0
	}
}

func (c *colorLight) buildHeartBeatMsg(mode int) []byte {
	baseInfo := &light1013.DeviceMessage_ReportDevBaseInfo{}

	baseInfo.Type = light1013.CmdType_DEV_REPORT_STATUS.Enum()
	baseInfo.OnoffStatus = light1013.LIGHT_ONOFF_STATUS_ON.Enum()
	baseInfo.LeaveHomeState = proto.Bool(true)

	switch mode {
	case 0:
		// 白光模式
		baseInfo.Mode = light1013.LIGHT_DEV_MODE_WHITE_LIGHT_MODE.Enum()
		baseInfo.White = buildWhiteLightData()
	case 1:
		// 彩光模式
		baseInfo.Mode = light1013.LIGHT_DEV_MODE_COLOR_LIGHT_MODE.Enum()
		baseInfo.Rgb = buildColorLightData()
	case 2:
		// 流光模式
		baseInfo.Mode = light1013.LIGHT_DEV_MODE_STREAMER_LIGHT_MODE.Enum()
		baseInfo.StreamLight = buildStreamLightData()
	}

	reportBaseInfo := &light1013.DeviceMessage_ReportDevBaseInfo_{
		ReportDevBaseInfo: baseInfo,
	}

	devMsg := &light1013.DeviceMessage{
		SessionId:  proto.Int32(-(rand.Int31n(math.MaxInt32))), //  填充负数
		DevMessage: reportBaseInfo,
	}

	data, err := proto.Marshal(devMsg)
	if err != nil {
		return nil
	}

	return data
}

// 白光数据
func buildWhiteLightData() *lightEvent.LampLightLevelCtlMessage {
	brightness := uint32(utils.RandInt64(1, 101))
	color := uint32(utils.RandInt64(1, 101))

	log.Debugf("设置白光模式，lum: %d, colorTemp: %d", brightness, color)
	return &lightEvent.LampLightLevelCtlMessage{
		Lum:       proto.Uint32(brightness),
		ColorTemp: proto.Uint32(color),
	}
}

// 彩光数据
func buildColorLightData() *lightEvent.LampLightRgbCtlMessage {
	brightness := uint32(utils.RandInt64(10, 100))
	red := uint32(utils.RandInt64(0, 255))
	green := uint32(utils.RandInt64(0, 255))
	blue := uint32(utils.RandInt64(0, 255))

	log.Debugf("设置彩光模式,  红: %d, 绿: %d, 蓝: %d, 亮度: %d", red, green, blue, brightness)

	return &lightEvent.LampLightRgbCtlMessage{
		Red:   proto.Uint32(red),
		Green: proto.Uint32(green),
		Blue:  proto.Uint32(blue),
		White: proto.Uint32(brightness),
	}
}

// 流光模式数据
func buildStreamLightData() *light1013.StreamLight {
	brightness := int32(utils.RandInt64(20, 100))
	// 速度
	speed := int32(utils.RandInt64(1, 3))

	points := &light1013.ColorPointMessage{
		PointA: buildRGBData("A"),
		PointB: buildRGBData("B"),
		PointC: buildRGBData("C"),
		PointD: buildRGBData("D"),
	}

	stream := &light1013.StreamLight{}

	stream.Time = proto.Int32(speed)
	stream.BrightnessPercent = proto.Int32(brightness)
	stream.Point = points

	return stream
}

// RGB 数值
func buildRGBData(point string) *light1013.RGBMessage {
	red := int32(utils.RandInt64(180, 255))
	green := int32(utils.RandInt64(0, 70))
	blue := int32(utils.RandInt64(65, 120))

	log.Debugf("流光模式 Point %s, red: %d, green: %d, blue: %d", point, red, green, blue)

	return &light1013.RGBMessage{
		Red:   proto.Int32(red),
		Green: proto.Int32(green),
		Blue:  proto.Int32(blue),
	}
}
