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
	lum          uint32
	colortemp    uint32
	status       lightT1012.LIGHT_ONOFF_STATUS
	lumTemp      uint32 //关灯之前，要把当前亮度临时放在在这里,待开灯时再从这里取出
	timerCounter int
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
	w.msgToServer = make(chan []byte, 2) //每次都忘记初始化，死性不改！！！！！！！！
	w.msgFromServer = make(chan []byte, 2)
	w.lum = uint32(commontool.GenerateRandNumber(20, 100))
	w.colortemp = uint32(commontool.GenerateRandNumber(20, 100))
	w.status = lightT1012.LIGHT_ONOFF_STATUS_ON
	return w
}

// 实现 Eufydevice 接口
func (w *whiteLight) RunMqttService() {
	w.SubHandler = func(c MQTT.Client, msg MQTT.Message) {
		go w.distributeMsg(msg.Payload()) //把消息发给 channel
	}
	w.MqttClient.ConnectToBroker()
	w.outgoing() //起一个 goroutine 作为publish message 的唯一出口
	w.inComing() //起一个 goroutine 处理订阅得到的服务器的控制消息(调亮度色温&开关灯)

}

func (w *whiteLight) distributeMsg(payload []byte) {
	w.msgFromServer <- payload
}

func (w *whiteLight) inComing() {
	go func() {
		for {
			select {
			case msg := <-w.msgFromServer:
				w.handleInComingMsg(msg) // 有必要新开一个 goroutine?
			}
		}
	}()
}

// 实现 Eufydevice 接口
func (w *whiteLight) SendHeartBeat() {
	w.msgToServer <- w.buildHeartBeatMsg()
	w.timerCounter++
	// 心跳是20秒一次，360次之后就是两小时, 每两小时请求一次时间闹钟离家模式
	if w.timerCounter == 360 {
		// 请求时间和闹钟
		w.msgToServer <- w.devRequestTimeAlertAwayMod(lightT1012.CmdType_DEV_REQUSET_TIME_Alarm)
		// 请求离家模式
		w.msgToServer <- w.devRequestTimeAlertAwayMod(lightT1012.CmdType_DEV_REQUEST_AWAYMODE_STATUS)
		// 重置为0
		w.timerCounter = 0
	}
}

// 构造心跳数据，需要亮度、色温、开关状态这三个值， 构造灯的实例的时候已有初始化，后面会随着订阅到的服务器的控制消息而改变
func (w *whiteLight) buildHeartBeatMsg() []byte {
	lightData := &lightEvent.LampLightLevelCtlMessage{}
	lightData.Lum = proto.Uint32(w.lum)
	if w.prod == "T1012" {
		lightData.ColorTemp = proto.Uint32(w.colortemp)
	}

	baseInfo := &lightT1012.DeviceMessage_ReportDevBaseInfo_{
		ReportDevBaseInfo: &lightT1012.DeviceMessage_ReportDevBaseInfo{
			Type:           lightT1012.CmdType_DEV_REPORT_STATUS.Enum(),
			OnoffStatus:    w.status.Enum(), // lightT1012.LIGHT_ONOFF_STATUS_ON.Enum(),
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

// 请求服务器同步时间、闹钟、离家模式数据
func (w *whiteLight) devRequestTimeAlertAwayMod(typ lightT1012.CmdType) []byte {
	nonePara := &lightT1012.DeviceMessage_NonParaMsg{
		NonParaMsg: &lightT1012.DeviceMessage_Non_ParamMsg{
			Type: typ.Enum(),
		},
	}
	devMsg := &lightT1012.DeviceMessage{
		SessionId:  proto.Int32(-(rand.Int31n(math.MaxInt32))), //  填充负数
		DevMessage: nonePara,
	}
	data, err := proto.Marshal(devMsg)
	if err != nil {
		log.Debug("build dev request time alert away mode error")
		return nil
	}
	return data
}

func (w *whiteLight) handleInComingMsg(payload []byte) {
	serMsg := &lightT1012.ServerMessage{}
	err := proto.Unmarshal(payload, serMsg)
	if err != nil {
		log.Debugf("unMarshal server message fail: %s", err)
		return
	}

	// 如果服务下发了同步时间和闹钟，则回一个ACK
	sta := serMsg.GetSync_Time_Alarm()
	if sta != nil {
		w.msgToServer <- w.devRequestTimeAlertAwayMod(lightT1012.CmdType_DEV_RESPONSE_TIME_Alarm_ACK)
		return
	}

	// 如果服务器下发了离家模式的配置数据， 则回一个ACK
	away := serMsg.GetSetAwayMode_Status()
	if away != nil {
		w.msgToServer <- w.devRequestTimeAlertAwayMod(lightT1012.CmdType_DEV_RESPONE_AWAYMODE_ACK)
		return
	}

	setlightdata := serMsg.GetSetLightData()
	if setlightdata == nil {
		return
	}

	//如果控制亮度色温的操作，则取出相应的值组装心跳消息，然后就结束
	lightCtrl := setlightdata.GetLightCtl()
	if lightCtrl != nil {
		w.lum = lightCtrl.GetLum()
		if w.prod == "T1012" {
			w.colortemp = lightCtrl.GetColorTemp()
		}
		w.msgToServer <- w.buildHeartBeatMsg()
		return
	}

	onOffStat := setlightdata.GetOnoffStatus()
	w.status = onOffStat
	// 如果是关灯操作， 则把当前亮度存给 lumTemp， 然后把亮度置为0， 改变状态， 色温不用管
	if onOffStat == lightT1012.LIGHT_ONOFF_STATUS_OFF {
		w.lumTemp = w.lum
		w.lum = 0
	} else {
		// 如果是开灯操作, 把上次关灯前的亮度取出来
		w.lum = w.lumTemp
	}
	w.msgToServer <- w.buildHeartBeatMsg()
}
