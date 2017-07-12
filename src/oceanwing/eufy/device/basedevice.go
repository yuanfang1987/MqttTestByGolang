package device

import (
	log "github.com/cihub/seelog"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// 字典的key
const (
	LeaveModeUpExp        = "leave_mode_up_exp"
	LeaveModeUp           = "leave_Mode_Up"
	LeaveModeDownExp      = "leave_mode_down_exp"
	LeaveModeDown         = "leave_Mode_Down"
	StatusBeforeLeaveMode = "status_before_leave_mode"
	StatusResumeToNormal  = "status_resume_to_normal"
	OccurTime             = "occur_time"
	OccurType             = "occur_type"
)

// EufyDevice 所有 eufy device 的行为接口
type EufyDevice interface {
	HandleSubscribeMessage() //这个方法必须要自己实现
	LeaveModeTestResult()
	GetSubDeviceTopic() string
	GetSubServerTopic() string
	GetProductCode() string
	GetProductKey() string
	SendPayload(MQTT.Message)
}

type baseDevice struct {
	ProdCode       string
	DevKEY         string
	DevID          string //预留，不一定能用得到
	SubDeviceTopic string
	SubServerTopic string
	SubMessage     chan MQTT.Message
	resultMap      map[string]string
	occurSlice     []map[string]string
}

func (b *baseDevice) GetSubDeviceTopic() string {
	return b.SubDeviceTopic
}

func (b *baseDevice) GetSubServerTopic() string {
	return b.SubServerTopic
}

func (b *baseDevice) GetProductCode() string {
	return b.ProdCode
}

func (b *baseDevice) GetProductKey() string {
	return b.DevKEY
}

func (b *baseDevice) SendPayload(msg MQTT.Message) {
	b.SubMessage <- msg
}

func (b *baseDevice) LeaveModeTestResult() {
	log.Infof("====================== %s (%s) 离家模式测试结果 ======================", b.DevKEY, b.ProdCode)

	// 预计开始时间
	if expStart, ok := b.resultMap[LeaveModeUpExp]; ok {
		log.Infof("预计开始时间: %s", expStart)
	} else {
		log.Error("无法获取到预计开始时间")
	}

	// 实际开始时间
	if startTime, ok := b.resultMap[LeaveModeUp]; ok {
		log.Infof("实际开始时间: %s", startTime)
	} else {
		log.Error("没有获取到实际开始时间")
	}

	// 预计结束时间
	if expEnd, ok := b.resultMap[LeaveModeDownExp]; ok {
		log.Infof("预计结束时间: %s", expEnd)
	} else {
		log.Error("无法获取预计结束时间")
	}

	// 实际结束时间
	if endTime, ok := b.resultMap[LeaveModeDown]; ok {
		log.Infof("实际结束时间: %s", endTime)
	} else {
		log.Error("没有获取到实际结束时间")
	}

	// 开始前的状态
	if statusBefore, ok := b.resultMap[StatusBeforeLeaveMode]; ok {
		log.Infof("开始前状态：%s", statusBefore)
	} else {
		log.Error("无法获取开始前的状态")
	}

	// 结束后的状态
	if statusAfter, ok := b.resultMap[StatusResumeToNormal]; ok {
		log.Infof("结束后状态: %s", statusAfter)
	} else {
		log.Error("无法获取结束后状态")
	}

	// 随机触发开关灯情况
	for i, v := range b.occurSlice {
		log.Infof("第 %d 次发生, 时间: %s, 类型: %s", i+1, v[OccurTime], v[OccurType])
	}
}
