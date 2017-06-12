package light

import (
	log "github.com/cihub/seelog"
	"github.com/golang/protobuf/proto"
	lightEvent "oceanwing/eufy/protobuf.lib/light/lightevent"
	lightT1012 "oceanwing/eufy/protobuf.lib/light/t1012"
)

type lightProd struct {
	devKEY    string
	devID     string
	pubTopicl string
	subTopicl string
	Incoming  chan []byte
}

func (l *lightProd) handleIncomingMsg() {
	go func() {
		for {
			select {
			case msg := <-l.Incoming:
				log.Infof("get new incoming message: %s", string(msg))
				// to do.
			}
		}
	}()
}

func buildSetLightDataMsg(sessionid int32, brightness, color uint32) []byte {
	o := &lightT1012.ServerMessage{
		SessionId:     proto.Int32(sessionid),
		RemoteMessage: setLightBrightAndColor(brightness, color),
	}
	data, err := proto.Marshal(o)
	if err != nil {
		log.Errorf("build set light data message fail: %s", err.Error())
		return nil
	}
	log.Debugf("build set light data message successfully, brightness: %d, color: %d", brightness, color)
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
