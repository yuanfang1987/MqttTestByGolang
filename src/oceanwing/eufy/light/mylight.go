package light

import (
	"math"
	"math/rand"

	log "github.com/cihub/seelog"
	"github.com/golang/protobuf/proto"
	serverAwayMode "oceanwing/eufy/protobuf.lib/common/server/awaymode"
	lightEvent "oceanwing/eufy/protobuf.lib/light/lightevent"
	lightT1012 "oceanwing/eufy/protobuf.lib/light/t1012"
)

type lightProd struct {
	prodCode  string
	devKEY    string
	devID     string
	pubTopicl string
	subTopicl string
	Incoming  chan []byte
	mode      lightT1012.DeviceMessage_ReportDevBaseInfo_LIGHT_DEV_MODE
	status    lightT1012.LIGHT_ONOFF_STATUS
	lum       uint32
	colorTemp uint32
}

func (l *lightProd) handleIncomingMsg() {
	go func() {
		log.Debugf("Running handleIncomingMsg function for device: %s", l.devKEY)
		for {
			select {
			case msg := <-l.Incoming:
				log.Infof("get new incoming message from device: %s", l.devKEY)
				l.unMarshalHeartBeatMsg(msg)
			}
		}
	}()
}

func (l *lightProd) buildSetLightDataMsg(brightness, color uint32) []byte {
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
	l.lum = brightness
	l.colorTemp = color
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

func (l *lightProd) buildSetAwayModeMsg(startHours, startMinutes, finishHours, finishMinutes,
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
func (l *lightProd) unMarshalHeartBeatMsg(incomingPayload []byte) {
	deviceMsg := &lightT1012.DeviceMessage{}
	err := proto.Unmarshal(incomingPayload, deviceMsg)
	if err != nil {
		log.Errorf("解析灯泡 %s (%s) 心跳消息失败: %s", l.devKEY, l.prodCode, err)
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
		log.Errorf("提取灯泡 %s (%s) 基础信息失败", l.devKEY, l.prodCode)
		return
	}

	log.Infof("解析灯泡 %s (%s) 的心跳消息成功", l.devKEY, l.prodCode)

	// --------------------- 判断结果 --------------------------------------------
	if lightT1012.CmdType_DEV_REPORT_STATUS != devBaseInfo.GetType() {
		log.Errorf("灯泡 %s (%s) CmdType 错误, 预期: %d, 实际: %d", l.devKEY, l.prodCode, lightT1012.CmdType_DEV_REPORT_STATUS, devBaseInfo.GetType())
	}

	if l.mode != devBaseInfo.GetMode() {
		log.Errorf("灯泡 %s (%s) Mode 错误, 预期: %d, 实际: %d", l.devKEY, l.prodCode, l.mode, devBaseInfo.GetMode())
	}

	if l.status != devBaseInfo.GetOnoffStatus() {
		log.Errorf("灯泡 %s (%s) Status 错误, 预期: %d, 实际: %d", l.devKEY, l.prodCode, l.status, devBaseInfo.GetOnoffStatus())
	}

	ligthCTRL := devBaseInfo.GetLightCtl()

	if l.lum != ligthCTRL.GetLum() {
		log.Errorf("灯泡 %s (%s) lum 错误, 预期: %d, 实际: %d", l.devKEY, l.prodCode, l.lum, ligthCTRL.GetLum())
	}

	// 只有 T1012 和 T1013 才有色温
	if l.prodCode != "T1011" {
		if l.colorTemp != ligthCTRL.GetColorTemp() {
			log.Errorf("灯泡 %s (%s) ColorTemp 错误, 预期: %d, 实际: %d", l.devKEY, l.prodCode, l.colorTemp, ligthCTRL.GetColorTemp())
		}
	}

}
