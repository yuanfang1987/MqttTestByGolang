package device

// EufyDevice 所有 eufy device 的行为接口
type EufyDevice interface {
	// 当订阅的主题有消息推送下来时，用此方法来处理消息
	HandleSubscribeMessage()

	// 获取当前 Device 订阅的主题
	GetSubTopic() string

	// 获取发布给当前Device消息时使用的主题
	GetPubTopic() string

	// 获取当前 Device 的 Product Code
	GetProductCode() string

	// 获取当前 Device 的 Device Key
	GetProductKey() string

	// 当测试结束时，获取对当前Device发出的指令总数
	GetSentCmds() int

	// 当测试结束时，获取当前Device的心跳总数
	GetDecodedheartBeat() int

	// 把订阅来的消息转给指定的 Channel
	SendPayload([]byte)

	// 组装消息并序列化，用于 publish
	BuildProtoBufMessage() []byte

	// 设置模式
	SetRunMode(int)
}

type baseDevice struct {
	ProdCode                   string
	DevKEY                     string
	DevID                      string //预留，不一定能用得到
	PubTopicl                  string
	SubTopicl                  string
	SubMessage                 chan []byte
	IsCmdSent                  bool
	IsTestPassed               bool
	CmdSentQuantity            int //下发的指令数量
	DecodeHeartBeatMsgQuantity int //解析的心跳消息数量
	HangOn                     int
	RunMod                     int // 0：NORMAL_MODE， 1： AWAY_MODE， 2： STREAMER_MODE
}

func (b *baseDevice) GetSubTopic() string {
	return b.SubTopicl
}

func (b *baseDevice) GetPubTopic() string {
	return b.PubTopicl
}

func (b *baseDevice) GetProductCode() string {
	return b.ProdCode
}

func (b *baseDevice) GetProductKey() string {
	return b.DevKEY
}

func (b *baseDevice) GetSentCmds() int {
	return b.CmdSentQuantity
}

func (b *baseDevice) GetDecodedheartBeat() int {
	return b.DecodeHeartBeatMsgQuantity
}

func (b *baseDevice) SendPayload(payload []byte) {
	b.SubMessage <- payload
}

func (b *baseDevice) PassedOrFailed(flag bool) string {
	if flag {
		return "Passed"
	}
	b.IsTestPassed = false
	return "Failed"
}

func (b *baseDevice) SetRunMode(mode int) {
	b.RunMod = mode
}
