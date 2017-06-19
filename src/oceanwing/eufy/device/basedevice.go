package device

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// EufyDevice 所有 eufy device 的行为接口
type EufyDevice interface {
	HandleSubscribeMessage()
	GetSubTopic() string
	GetSubTopicServer() string
	GetPubTopic() string
	GetProductCode() string
	GetProductKey() string
	GetSentCmds() int
	GetDecodedheartBeat() int
	SendPayload(MQTT.Message)
	BuildProtoBufMessage() []byte
}

type baseDevice struct {
	ProdCode                   string
	DevKEY                     string
	DevID                      string //预留，不一定能用得到
	PubTopicl                  string
	SubTopicl                  string
	SubMessage                 chan MQTT.Message
	IsCmdSent                  bool
	CmdSentQuantity            int //下发的指令数量
	DecodeHeartBeatMsgQuantity int //解析的心跳消息数量
}
