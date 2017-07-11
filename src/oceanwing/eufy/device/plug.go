package device

import (
	"oceanwing/commontool"

	log "github.com/cihub/seelog"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/protobuf/proto"
	"oceanwing/eufy/protobuf.lib/switch/t1201"
)

// Plug 是对插座类产品的一个struct描述,目前有 T1201
type Plug struct {
	baseDevice
	OnOffStatus uint32
}

// NewPlug 新建一个 Plug 实例.
func NewPlug(prodCode, devKey string) EufyDevice {
	p := &Plug{}
	p.ProdCode = prodCode
	p.DevKEY = devKey
	p.SubDeviceTopic = "DEVICE/" + prodCode + "/" + devKey + "/PUH_MESSAGE"
	p.SubServerTopic = "DEVICE/" + prodCode + "/" + devKey + "/SUB_MESSAGE"
	p.SubMessage = make(chan MQTT.Message)
	log.Infof("Create a Plug, product code: %s, device key: %s", prodCode, devKey)
	return p
}

// HandleSubscribeMessage 实现 EufyDevice 接口
func (p *Plug) HandleSubscribeMessage() {
	go func() {
		log.Debugf("Running handleIncomingMsg function for device: %s", p.DevKEY)
		for {
			select {
			case msg := <-p.SubMessage:
				go p.unMarshalAllMsg(msg)
			}
		}
	}()
}

func (p *Plug) unMarshalAllMsg(msg MQTT.Message) {
	t := msg.Topic()
	payload := msg.Payload()

	if p.SubDeviceTopic == t {
		// 设备心跳消息
		log.Info("----- 这是一个来自插座设备的心跳消息----------")
		p.unMarshalDevMessage(payload)
	} else if p.SubServerTopic == t {
		//服务器消息
		log.Info("-------这是一个来自服务器的控制消息---------")
		p.unMarshalServerMsg(payload)
	}
}

// 解析插座心跳
func (p *Plug) unMarshalDevMessage(incomingPayload []byte) {
	devMsg := &t1201.DeviceMessage{}
	err := proto.Unmarshal(incomingPayload, devMsg)
	if err != nil {
		log.Errorf("解析插座 %s (%s) 心跳消息失败: %s", p.DevKEY, p.ProdCode, err)
		return
	}
	log.Infof("SessionID: %d", devMsg.GetSessionId())

	// 无参数消息,只需关注其CmdType
	noneParaMsg := devMsg.GetNonParaMsg()
	if noneParaMsg != nil {
		log.Infof("插座 %s (%s) 无参数消息, 指令类型: %s", p.DevKEY, p.ProdCode, noneParaMsg.GetType().String())
	}

	// heart beat
	heartBeat := devMsg.GetHeartBeat()
	if heartBeat != nil {
		// CmdType
		cmd := heartBeat.GetType().String()
		log.Infof("插座 %s (%s) 心跳 CmdType: %s", p.DevKEY, p.ProdCode, cmd)
		// 开关状态
		relayState := heartBeat.GetRelayState()
		log.Infof("插座 %s (%s) 心跳 开关状态: %d", p.DevKEY, p.ProdCode, relayState)
		// 功率
		power := heartBeat.GetPower()
		log.Infof("插座 %s (%s) 心跳 功率: %d", p.DevKEY, p.ProdCode, power)
	}

	// ElectricMessage 这是什么鬼啊?
	electricMsg := devMsg.GetElectricMessage()
	if electricMsg != nil {
		// CmdType
		cmd := electricMsg.GetType().String()
		log.Infof("插座 %s (%s) ElectricMessage 的CmdType: %s", p.DevKEY, p.ProdCode, cmd)
		// electric
		electric := electricMsg.GetElectric()
		log.Infof("插座 %s (%s) ElectricMessage 的 electric: %d", p.DevKEY, p.ProdCode, electric)
		// workingTime
		workingTime := electricMsg.GetWorkingTime()
		log.Infof("插座 %s (%s) ElectricMessage 的 workingTime: %d", p.DevKEY, p.ProdCode, workingTime)
	}

	// APPDataMessage,这个应该用不着，应该是给局域网内APP用的
}

// 解析服务器消息
func (p *Plug) unMarshalServerMsg(payload []byte) {
	serMsg := &t1201.ServerMessage{}
	err := proto.Unmarshal(payload, serMsg)
	if err != nil {
		log.Errorf("解析服务器控制插座消息失败: %s", err)
		return
	}

	// session id
	log.Infof("Session DI: %d", serMsg.GetSessionId())

	// ServerMessage --> DevStateMessage
	devStat := serMsg.GetDevState()
	if devStat != nil {
		cmd := devStat.GetType()
		log.Infof("CmdType: %s", cmd.String())
		replayStat := devStat.GetRelayState()
		log.Infof("RelayState: %d", replayStat)
	}

	// ServerMessage --> Sync_Time_Alarm
	sta := serMsg.GetSync_Time_Alarm()
	if sta != nil {
		log.Info("解析出同步时间和闹钟消息")

		cmd := sta.GetType()
		log.Infof("CmdType: %s", cmd.String())

		// 时间
		synctime := sta.GetTime()
		if synctime != nil {
			log.Info("------时间信息：")
			log.Infof("------年： %d", synctime.GetYear())
			log.Infof("------月： %d", synctime.GetMonth()+1)
			log.Infof("------日： %d", synctime.GetDay())
			log.Infof("------weekday: %d", synctime.GetWeekday())
			log.Infof("------时： %d", synctime.GetHours())
			log.Infof("------分： %d", synctime.GetMinutes())
			log.Infof("------秒： %d", synctime.GetSeconds())
		}

		// 闹钟
		syncAlert := sta.GetAlarm()
		if syncAlert != nil {
			alarmRecordData := syncAlert.GetAlarmRecordData()
			for i, v := range alarmRecordData {
				log.Infof("=====第 %d 个闹钟信息====", i+1)
				alarmMsg := v.GetAlarmMesage()
				if alarmMsg != nil {
					log.Infof("------闹钟，时：%d", alarmMsg.GetHours())
					log.Infof("------闹钟，分: %d", alarmMsg.GetMinutes())
					log.Infof("------闹钟，是否重复: %t", alarmMsg.GetRepetiton())
					weekinfo := commontool.ConvertToWeekDay(alarmMsg.GetWeekInfo())
					log.Infof("------闹钟，WeekInfo：%s", weekinfo)
				}

				replay := v.GetRelayState()
				log.Infof("-------RelayState: %d", replay)
			}
		}

		// 离家模式配置信息
		leaveMsg := sta.GetSyncLeaveModeMsg()
		if leaveMsg != nil {
			starthour := int(leaveMsg.GetStartHours())
			startminute := int(leaveMsg.GetStartMinutes())
			finishhour := int(leaveMsg.GetFinishHours())
			finishminute := int(leaveMsg.GetFinishMinutes())

			log.Infof("------离家模式, 开始时间  %d:%d", starthour, startminute)
			log.Infof("------离家模式, 结束时间  %d:%d", finishhour, finishminute)
			log.Infof("------离家模式, 是否重复: %t", leaveMsg.GetRepetiton())
			log.Infof("------离家模式, WeekInfo: %s", commontool.ConvertToWeekDay(leaveMsg.GetWeekInfo()))
			leaveHomeFlag := leaveMsg.GetLeaveHomeState()
			log.Infof("------离家模式, 是否开启: %t", leaveHomeFlag)

			// 如果服务器下发的 离家模式 为 true，则记下时间
			if leaveHomeFlag {
				// light.resultMap["leave_mode_up_exp"] = strconv.Itoa(starthour) + ":" + strconv.Itoa(startminute)
				// light.resultMap["leave_mode_down_exp"] = strconv.Itoa(finishhour) + ":" + strconv.Itoa(finishminute)
			}

		}
	}
}
