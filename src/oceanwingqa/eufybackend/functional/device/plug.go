package device

import (
	"math"
	"math/rand"

	log "github.com/cihub/seelog"
	"github.com/golang/protobuf/proto"
	"oceanwingqa/common/protobuf.lib/switch/t1201"
)

// Plug 是对插座类产品的一个struct描述,目前有 T1201
type Plug struct {
	baseDevice
	OnOffStatus uint32
}

// NewPlug 新建一个 Plug 实例.
func NewPlug(prodCode, devKey, devid string) EufyDevice {
	p := &Plug{}
	p.ProdCode = prodCode
	p.DevKEY = devKey
	p.DevID = devid
	p.PubTopicl = "DEVICE/T1201/" + devKey + "/SUB_MESSAGE"
	p.SubTopicl = "DEVICE/T1201/" + devKey + "/PUH_MESSAGE"
	p.DeviceMsg = make(chan []byte, 3)
	p.ServerMsg = make(chan []byte, 3)
	log.Infof("Create a Plug, product code: %s, device key: %s, device id: %s", prodCode, devKey, devid)
	return p
}

// HandleSubscribeMessage 实现 EufyDevice 接口
func (p *Plug) HandleSubscribeMessage() {
	go func() {
		log.Debugf("Running handleIncomingMsg function for device: %s", p.DevKEY)
		for {
			select {
			case msg := <-p.DeviceMsg:
				log.Info("get new incoming message from device: %s", p.DevKEY)
				p.unMarshalDevMessage(msg)
			case <-p.ServerMsg:
				// to do .
			}
		}
	}()
}

// BuildProtoBufMessage 实现 EufyDevice 接口
func (p *Plug) BuildProtoBufMessage() []byte {
	protoMsg := p.buildDevStateMessage(1)
	data, err := proto.Marshal(protoMsg)

	if err != nil {
		log.Errorf("build build message fail: %s", err.Error())
		return nil
	}

	log.Info("=================================================")
	return data
}

// 控制插座开关
func (p *Plug) buildDevStateMessage(status uint32) *t1201.ServerMessage {
	devStateMsg := &t1201.ServerMessage_DevState{
		DevState: &t1201.DevStateMessage{
			Type:       t1201.CmdType_REMOTE_SET_PLUG_STATE.Enum(),
			RelayState: proto.Uint32(status), // 1: 开, 0: 关
		},
	}

	o := &t1201.ServerMessage{
		SessionId:     proto.Int32(rand.Int31n(math.MaxInt32)),
		RemoteMessage: devStateMsg,
	}
	return o
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
		log.Infof("插座 %s (%s) 无参数消息, 指令类型: %d", p.DevKEY, p.ProdCode, noneParaMsg.GetType())
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
