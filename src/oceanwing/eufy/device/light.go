package device

import (
	"fmt"
	"math"
	"math/rand"
	"oceanwing/commontool"
	"oceanwing/eufy/result"

	log "github.com/cihub/seelog"
	"github.com/golang/protobuf/proto"
	serverAwayMode "oceanwing/eufy/protobuf.lib/common/server/awaymode"
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
func (light *Light) SendPayload(pl []byte) {
	light.SubMessage <- pl
}

// BuildProtoBufMessage 实现 EufyDevice 接口
func (light *Light) BuildProtoBufMessage() []byte {
	o := &lightT1012.ServerMessage{
		SessionId:     proto.Int32(rand.Int31n(math.MaxInt32)),
		RemoteMessage: light.setLightBrightAndColor(),
	}
	data, err := proto.Marshal(o)

	if err != nil {
		log.Errorf("build set light data message fail: %s", err.Error())
		return nil
	}

	// 设置 IsCmdSent 标志为 true
	light.IsCmdSent = true
	// 已下发的指令数量累加 1
	light.CmdSentQuantity++

	return data
}

// SetLightData is a struct
// brightness: 亮度，color: 色温,  ServerMessage_SetLightData
func (light *Light) setLightBrightAndColor() *lightT1012.ServerMessage_SetLightData_ {
	seed := commontool.RandInt64(0, 10)
	var content string
	// seed 随机数产生的范围是 0 到 9 共 10 个数字，则用 30%的概率去执行开关灯， 剩下的执行调节亮度色温
	if seed < 3 {
		var nextStatus *lightT1012.LIGHT_ONOFF_STATUS
		// 如果灯的当前状态是开着的，则执行关闭操作， 反之则执行打开操作
		if light.status == lightT1012.LIGHT_ONOFF_STATUS_ON {
			nextStatus = lightT1012.LIGHT_ONOFF_STATUS_OFF.Enum()
			log.Info("关灯")
			// 关灯后， 亮度和色温都应该变成 0
			light.status = 0
			light.lum = 0
			light.colorTemp = 0
			content = "关灯"
		} else {
			nextStatus = lightT1012.LIGHT_ONOFF_STATUS_ON.Enum()
			log.Info("开灯")
			// 开灯后，亮度为100，色温为0，but why???
			light.status = 1
			light.lum = 100
			light.colorTemp = 0
			content = "开灯"
		}

		return &lightT1012.ServerMessage_SetLightData_{
			SetLightData: &lightT1012.ServerMessage_SetLightData{
				Type:        lightT1012.CmdType_REMOTE_SET_LIGHTING_PARA.Enum(),
				OnoffStatus: nextStatus,
			},
		}
	}

	// 调节亮度和色温, 随机产生亮度和色温的值
	brightness := uint32(commontool.RandInt64(0, 101))
	color := uint32(commontool.RandInt64(0, 101))
	light.lum = brightness
	light.colorTemp = color
	light.status = 1
	log.Infof("执行调节亮度色温操作, lum: %d, colorTemp: %d", brightness, color)
	content = "调节亮度和色温"

	// 在.csv 结果文件上打个标志
	result.WriteToResultFile(light.ProdCode, light.DevKEY, "NA", content, "NA")

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

	noneParaMsg := deviceMsg.GetNonParaMsg()
	if noneParaMsg != nil {
		log.Infof("Cmd type of Non_ParamMsg: %d", noneParaMsg.GetType())
	}

	devBaseInfo := deviceMsg.GetReportDevBaseInfo()
	if devBaseInfo == nil {
		log.Warnf("提取灯泡 %s (%s) 基础信息失败", light.DevKEY, light.ProdCode)
		return
	}

	log.Infof("解析灯泡 %s (%s) 的心跳消息成功", light.DevKEY, light.ProdCode)

	// --------------------- 判断结果 --------------------------------------------

	// 只有在给设备下发了指令之后，才去判断它的即时心跳， 常规心跳不要管
	if !light.IsCmdSent {
		log.Info("尚未有指令下发给设备，无需判断心跳")
		return
	}

	// 重置
	light.IsCmdSent = false
	// 已解析的心跳数量累加 1
	light.DecodeHeartBeatMsgQuantity++

	var assertFlag string
	var testContent string

	//  CmdType
	assertFlag = result.PassedOrFailed(lightT1012.CmdType_DEV_REPORT_STATUS == devBaseInfo.GetType())
	testContent = fmt.Sprintf("灯泡 %s (%s) CmdType, 预期: %d, 实际: %d", light.DevKEY, light.ProdCode, lightT1012.CmdType_DEV_REPORT_STATUS, devBaseInfo.GetType())
	result.WriteToResultFile(light.ProdCode, light.DevKEY, "CmdType", testContent, assertFlag)
	log.Info(testContent)

	// Mode
	assertFlag = result.PassedOrFailed(light.mode == devBaseInfo.GetMode())
	testContent = fmt.Sprintf("灯泡 %s (%s) Mode, 预期: %d, 实际: %d", light.DevKEY, light.ProdCode, light.mode, devBaseInfo.GetMode())
	result.WriteToResultFile(light.ProdCode, light.DevKEY, "Mode", testContent, assertFlag)
	log.Info(testContent)

	// Status
	assertFlag = result.PassedOrFailed(light.status == devBaseInfo.GetOnoffStatus())
	testContent = fmt.Sprintf("灯泡 %s (%s) Status, 预期: %d, 实际: %d", light.DevKEY, light.ProdCode, light.status, devBaseInfo.GetOnoffStatus())
	result.WriteToResultFile(light.ProdCode, light.DevKEY, "Status", testContent, assertFlag)
	log.Info(testContent)

	ligthCTRL := devBaseInfo.GetLightCtl()
	if ligthCTRL == nil {
		log.Error("GetLightCtl fail")
		return
	}

	// lum
	assertFlag = result.PassedOrFailed(light.lum == ligthCTRL.GetLum())
	testContent = fmt.Sprintf("灯泡 %s (%s) lum, 预期: %d, 实际: %d", light.DevKEY, light.ProdCode, light.lum, ligthCTRL.GetLum())
	result.WriteToResultFile(light.ProdCode, light.DevKEY, "Lum", testContent, assertFlag)
	log.Info(testContent)

	// 只有 T1012 和 T1013 才有色温
	if light.ProdCode != "T1011" {
		assertFlag = result.PassedOrFailed(light.colorTemp == ligthCTRL.GetColorTemp())
		testContent = fmt.Sprintf("灯泡 %s (%s) ColorTemp, 预期: %d, 实际: %d", light.DevKEY, light.ProdCode, light.colorTemp, ligthCTRL.GetColorTemp())
		result.WriteToResultFile(light.ProdCode, light.DevKEY, "ColorTemp", testContent, assertFlag)
		log.Info(testContent)
	}

}
