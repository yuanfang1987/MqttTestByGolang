package performance

import (
	"math"
	"math/rand"

	log "github.com/cihub/seelog"

	"oceanwing/commontool"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/protobuf/proto"
	lightEvent "oceanwing/eufy/protobuf.lib/light/lightevent"
	lightT1012 "oceanwing/eufy/protobuf.lib/light/t1012"
)

// 这是一个描述白灯的结构体，支持 T1011, T1012
type whiteLight struct {
	baseEufy
}

// 创建新的灯的实例
func newWhiteLight(clientid, username, pwd, broker, prodCode, devKey string, needCA bool) Eufydevice {
	w := &whiteLight{}
	w.Clientid = clientid
	w.Username = username
	w.Pwd = pwd
	w.Broker = broker
	w.prod = prodCode
	w.PubTopic = "DEVICE/" + prodCode + "/" + devKey + "/PUH_MESSAGE"
	w.SubTopic = "DEVICE/" + prodCode + "/" + devKey + "/SUB_MESSAGE"
	w.NeedCA = needCA
	return w
}

// 实现 Eufydevice 接口
func (w *whiteLight) RunMqttService() {
	w.SubHandler = func(c MQTT.Client, msg MQTT.Message) {} // do nothing.
	w.MqttClient.ConnectToBroker()
	w.outgoing()
}

// 实现 Eufydevice 接口
func (w *whiteLight) SendHeartBeat() {
	w.msgToServer <- w.buildHeartBeatMsg()
}

// 构造心跳数据，把问题简单化，不要搞那么复杂，只考虑亮度和色温变化即可，不管开关灯
func (w *whiteLight) buildHeartBeatMsg() []byte {
	brightness := uint32(commontool.RandInt64(20, 100))
	color := uint32(commontool.RandInt64(20, 100))

	lightData := &lightEvent.LampLightLevelCtlMessage{}
	lightData.Lum = proto.Uint32(brightness)
	if w.prod == "T1012" {
		lightData.ColorTemp = proto.Uint32(color)
	}

	baseInfo := &lightT1012.DeviceMessage_ReportDevBaseInfo_{
		ReportDevBaseInfo: &lightT1012.DeviceMessage_ReportDevBaseInfo{
			Type:           lightT1012.CmdType_DEV_REPORT_STATUS.Enum(),
			OnoffStatus:    lightT1012.LIGHT_ONOFF_STATUS_ON.Enum(),
			Mode:           lightT1012.DeviceMessage_ReportDevBaseInfo_NORMAL_MODE.Enum(),
			LightCtl:       lightData,
			LastOnLightCtl: lightData,
		},
	}

	devMsg := &lightT1012.DeviceMessage{
		SessionId:  proto.Int32(-(rand.Int31n(math.MaxInt32))), //  填充负数
		DevMessage: baseInfo,
	}

	data, err := proto.Marshal(devMsg)
	if err != nil {
		log.Debug("build heart beat message error")
		return nil
	}

	return data
}
