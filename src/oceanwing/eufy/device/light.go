package device

import (
	"math"
	"math/rand"

	"oceanwing/commontool"

	log "github.com/cihub/seelog"
	"github.com/golang/protobuf/proto"
	serverAwayMode "oceanwing/eufy/protobuf.lib/common/server/awaymode"
	lightEvent "oceanwing/eufy/protobuf.lib/light/lightevent"
	lightT1012 "oceanwing/eufy/protobuf.lib/light/t1012"
)

// Light 灯泡类产品的一个 struct 描述，包括 T1011,T1012,T1013.
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
	o.PubTopicl = "DEVICE/T1012/" + devKey + "/SUB_MESSAGE"
	o.SubTopicl = "DEVICE/T1012/" + devKey + "/PUH_MESSAGE"
	o.SubMessage = make(chan []byte)
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
				light.unMarshalHeartBeatMsg(msg)
			}
		}
	}()
}

// GetSubTopic 实现 EufyDevice 接口
func (light *Light) GetSubTopic() string {
	return light.SubTopicl
}

// GetPubTopic 实现 EufyDevice 接口
func (light *Light) GetPubTopic() string {
	return light.PubTopicl
}

// SendPayload 实现 EufyDevice 接口
func (light *Light) SendPayload(pl []byte) {
	light.SubMessage <- pl
}

// BuildProtoBufMessage 实现 EufyDevice 接口
func (light *Light) BuildProtoBufMessage() []byte {
	brightness := uint32(commontool.RandInt64(1, 100))
	color := uint32(commontool.RandInt64(0, 100))
	return light.buildSetLightDataMsg(brightness, color)
}

func (light *Light) buildSetLightDataMsg(brightness, color uint32) []byte {
	o := &lightT1012.ServerMessage{
		SessionId:     proto.Int32(rand.Int31n(math.MaxInt32)),
		RemoteMessage: setLightBrightAndColor(brightness, color),
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
	return data
}

// SetLightData is a struct
// brightness: 亮度，color: 色温,  ServerMessage_SetLightData
func setLightBrightAndColor(brightness, color uint32) *lightT1012.ServerMessage_SetLightData_ {
	return &lightT1012.ServerMessage_SetLightData_{
		SetLightData: &lightT1012.ServerMessage_SetLightData{
			Type: lightT1012.CmdType_REMOTE_SET_LIGHTING_PARA.Enum(),
			LightCtl: &lightEvent.LampLightLevelCtlMessage{
				Lum:       proto.Uint32(brightness),
				ColorTemp: proto.Uint32(color),
			},
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

	// debug
	noneParaMsg := deviceMsg.GetNonParaMsg()
	if noneParaMsg != nil {
		log.Infof("Cmd type of Non_ParamMsg: %d", noneParaMsg.GetType())
	}
	// end debug

	devBaseInfo := deviceMsg.GetReportDevBaseInfo()
	if devBaseInfo == nil {
		log.Errorf("提取灯泡 %s (%s) 基础信息失败", light.DevKEY, light.ProdCode)
		return
	}

	log.Infof("解析灯泡 %s (%s) 的心跳消息成功", light.DevKEY, light.ProdCode)

	// --------------------- 判断结果 --------------------------------------------
	if lightT1012.CmdType_DEV_REPORT_STATUS != devBaseInfo.GetType() {
		log.Errorf("灯泡 %s (%s) CmdType 错误, 预期: %d, 实际: %d", light.DevKEY, light.ProdCode, lightT1012.CmdType_DEV_REPORT_STATUS, devBaseInfo.GetType())
	}

	if light.mode != devBaseInfo.GetMode() {
		log.Errorf("灯泡 %s (%s) Mode 错误, 预期: %d, 实际: %d", light.DevKEY, light.ProdCode, light.mode, devBaseInfo.GetMode())
	}

	if light.status != devBaseInfo.GetOnoffStatus() {
		log.Errorf("灯泡 %s (%s) Status 错误, 预期: %d, 实际: %d", light.DevKEY, light.ProdCode, light.status, devBaseInfo.GetOnoffStatus())
	}

	ligthCTRL := devBaseInfo.GetLightCtl()

	if light.lum != ligthCTRL.GetLum() {
		log.Errorf("灯泡 %s (%s) lum 错误, 预期: %d, 实际: %d", light.DevKEY, light.ProdCode, light.lum, ligthCTRL.GetLum())
	}

	// 只有 T1012 和 T1013 才有色温
	if light.ProdCode != "T1011" {
		if light.colorTemp != ligthCTRL.GetColorTemp() {
			log.Errorf("灯泡 %s (%s) ColorTemp 错误, 预期: %d, 实际: %d", light.DevKEY, light.ProdCode, light.colorTemp, ligthCTRL.GetColorTemp())
		}
	}

}