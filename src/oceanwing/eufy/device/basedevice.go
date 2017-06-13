package device

// EufyDevice 所有 eufy device 的行为接口
type EufyDevice interface {
	HandleSubscribeMessage()
	GetSubTopic() string
	GetPubTopic() string
	SendPayload([]byte)
	BuildProtoBufMessage() []byte
}

type baseDevice struct {
	ProdCode   string
	DevKEY     string
	DevID      string //预留，不一定能用得到
	PubTopicl  string
	SubTopicl  string
	SubMessage chan []byte
}
