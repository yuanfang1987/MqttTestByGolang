package device

import (
	"math"
	"math/rand"
	"oceanwing/commontool"

	log "github.com/cihub/seelog"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/protobuf/proto"
	serverAwayMode "oceanwing/eufy/protobuf.lib/common/server/awaymode"
	lightEvent "oceanwing/eufy/protobuf.lib/light/lightevent"
	lightT1012 "oceanwing/eufy/protobuf.lib/light/t1012"
)

// Light 灯泡类产品的一个 struct 描述，包括 T1011,T1012,T1013
type Light struct {
	baseDevice
	mode           lightT1012.DeviceMessage_ReportDevBaseInfo_LIGHT_DEV_MODE
	status         lightT1012.LIGHT_ONOFF_STATUS
	lum            uint32
	colorTemp      uint32
	subServerTopic string
}

// NewLight 新建一个 Light 实例
func NewLight(prodCode, devKey string) EufyDevice {
	o := &Light{
		mode:   0,
		status: 1,
	}
	o.ProdCode = prodCode
	o.DevKEY = devKey
	o.PubTopicl = "DEVICE/T1012/" + devKey + "/SUB_MESSAGE"
	o.SubTopicl = "DEVICE/T1012/" + devKey + "/PUH_MESSAGE"      // 订阅设备的消息
	o.subServerTopic = "DEVICE/T1012/" + devKey + "/SUB_MESSAGE" //订阅服务器的消息
	o.SubMessage = make(chan MQTT.Message)
	return o
}

// HandleSubscribeMessage 实现 EufyDevice 接口
func (light *Light) HandleSubscribeMessage() {
	go func() {
		log.Debugf("Running handleIncomingMsg function for device: %s", light.DevKEY)
		for {
			select {
			case msg := <-light.SubMessage:
				log.Infof("get new incoming message from device: %s", light.DevKEY)
				light.unMarshalAllMessage(msg)
			}
		}
	}()
}

// GetSubTopic 实现 EufyDevice 接口
func (light *Light) GetSubTopic() string {
	return light.SubTopicl
}

// GetSubTopicServer 实现 EufyDevice 接口
func (light *Light) GetSubTopicServer() string {
	return light.subServerTopic
}

// GetPubTopic 实现 EufyDevice 接口
func (light *Light) GetPubTopic() string {
	return light.PubTopicl
}

// GetProductCode 实现 EufyDevice 接口
func (light *Light) GetProductCode() string {
	return light.ProdCode
}

// GetProductKey 实现 EufyDevice 接口
func (light *Light) GetProductKey() string {
	return light.DevKEY
}

// GetSentCmds 实现 EufyDevice 接口
func (light *Light) GetSentCmds() int {
	return light.CmdSentQuantity
}

// GetDecodedheartBeat 实现 EufyDevice 接口
func (light *Light) GetDecodedheartBeat() int {
	return light.DecodeHeartBeatMsgQuantity
}

// SendPayload 实现 EufyDevice 接口
func (light *Light) SendPayload(msg MQTT.Message) {
	light.SubMessage <- msg
}

// BuildProtoBufMessage 实现 EufyDevice 接口
func (light *Light) BuildProtoBufMessage() []byte {
	// 随机产生亮度和色温的值
	brightness := uint32(commontool.RandInt64(0, 101))
	color := uint32(commontool.RandInt64(0, 101))
	// 随机产生是否开灯的值，只有两个可选值： 0 和 1
	var onOffStatus *lightT1012.LIGHT_ONOFF_STATUS
	onOffValue := uint32(commontool.RandInt64(0, 2))
	if onOffValue == 1 {
		onOffStatus = lightT1012.LIGHT_ONOFF_STATUS_ON.Enum()
	} else {
		onOffStatus = lightT1012.LIGHT_ONOFF_STATUS_OFF.Enum()
	}
	// 设置 IsCmdSent 标志为 true
	light.IsCmdSent = true
	// 已下发的指令数量累加 1
	light.CmdSentQuantity++
	return light.buildSetLightDataMsg(brightness, color, onOffStatus)
}

func (light *Light) buildSetLightDataMsg(brightness, color uint32, status *lightT1012.LIGHT_ONOFF_STATUS) []byte {
	o := &lightT1012.ServerMessage{
		SessionId:     proto.Int32(rand.Int31n(math.MaxInt32)),
		RemoteMessage: setLightBrightAndColor(brightness, color, status),
	}
	data, err := proto.Marshal(o)
	if err != nil {
		log.Errorf("build set light data message fail: %s", err.Error())
		return nil
	}
	log.Debugf("build set light data message successfully, brightness: %d, color: %d", brightness, color)
	// 确保 Marshal 成功后，再更改 lightProd 的值
	light.lum = brightness
	light.colorTemp = color
	if status == lightT1012.LIGHT_ONOFF_STATUS_ON.Enum() {
		light.status = 1
	} else {
		light.status = 0
	}
	return data
}

// SetLightData is a struct
// brightness: 亮度，color: 色温,  ServerMessage_SetLightData
func setLightBrightAndColor(brightness, color uint32, status *lightT1012.LIGHT_ONOFF_STATUS) *lightT1012.ServerMessage_SetLightData_ {
	return &lightT1012.ServerMessage_SetLightData_{
		SetLightData: &lightT1012.ServerMessage_SetLightData{
			Type: lightT1012.CmdType_REMOTE_SET_LIGHTING_PARA.Enum(),
			LightCtl: &lightEvent.LampLightLevelCtlMessage{
				Lum:       proto.Uint32(brightness),
				ColorTemp: proto.Uint32(color),
			},
			OnoffStatus: status,
		},
	}
}

func (light *Light) buildSetAwayModeMsg(startHours, startMinutes, finishHours, finishMinutes,
	weekInfo, leaveMode uint32, repetiton, LeaveHomeState bool) []byte {
	// set away mode msg
	awayMod := &lightT1012.ServerMessage_SetAwayMode_Status{
		SetAwayMode_Status: &lightT1012.ServerMessage_SetAwayMode{
			Type: lightT1012.CmdType_REMOTE_SET_AWAYMODE_STATUS.Enum(),
			SyncLeaveModeMsg: &serverAwayMode.LeaveHomeMessage{
				StartHours:     proto.Uint32(startHours),
				StartMinutes:   proto.Uint32(startMinutes),
				FinishHours:    proto.Uint32(finishHours),
				FinishMinutes:  proto.Uint32(finishMinutes),
				Repetiton:      proto.Bool(repetiton),
				WeekInfo:       proto.Uint32(weekInfo),
				LeaveHomeState: proto.Bool(LeaveHomeState),
				LeaveMode:      proto.Uint32(leaveMode),
			},
		},
	}

	o := &lightT1012.ServerMessage{
		SessionId:     proto.Int32(rand.Int31n(math.MaxInt32)),
		RemoteMessage: awayMod,
	}
	data, err := proto.Marshal(o)
	if err != nil {
		log.Errorf("build set leave home mode message fail: %s", err.Error())
		return nil
	}
	log.Debug("build set leave home mode message successfully.")
	return data
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
		log.Errorf("UnMarshal server message fail: %s", err)
		return
	}

	// session id
	log.Infof("Session ID: %d", serMsg.GetSessionId())

	// SetLightData
	setlightdata := serMsg.GetSetLightData()
	if setlightdata != nil {
		cmd := lightT1012.CmdType_name[int32(setlightdata.GetType())]
		log.Infof("CmdType: %s", cmd)

		lightctrl := setlightdata.GetLightCtl()
		if lightctrl != nil {
			log.Infof("Lum: %d", lightctrl.GetLum())
			log.Infof("ColorTemp: %d", lightctrl.GetColorTemp())
		}

		status := lightT1012.LIGHT_ONOFF_STATUS_name[int32(setlightdata.GetOnoffStatus())]
		log.Infof("OnOff Status: %s", status)
	}

}

func (light *Light) unMarshalAllMessage(msg MQTT.Message) {
	t := msg.Topic()
	payload := msg.Payload()

	if light.SubTopicl == t {
		// 设备心跳消息
		log.Info("----- 这是一个来自设备的心跳消息----------")
		light.unMarshalHeartBeatMsg(payload)
	} else if light.subServerTopic == t {
		//服务器消息
		log.Info("-------这是一个来自服务器的控制消息---------")
		light.unMarshalServerMsg(payload)
	}
}
