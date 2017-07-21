package performance

import (
	"oceanwing/mqttclient"

	log "github.com/cihub/seelog"
)

// Eufydevice 定义 eufy 智能硬件行为的接口
type Eufydevice interface {
	// 连接broker并订阅主题
	RunMqttService()
	// 发心跳
	SendHeartBeat()
}

// 所有eufy device的基类
type baseEufy struct {
	mqttclient.MqttClient
	msgToServer   chan []byte
	msgFromServer chan []byte
	prod          string // product code
}

func (b *baseEufy) RunMqttService() {
	log.Error("function RunMqttService() Not implemet yet.")
}
func (b *baseEufy) SendHeartBeat() {
	log.Error("function SendHeartBeat() Not implement yet.")
}

// 这是每个设备发出消息的唯一出口, 为了防止常规的20秒一次的心跳与 “即时心跳” 发生冲突
func (b *baseEufy) outgoing() {
	go func() {
		for {
			select {
			case payload := <-b.msgToServer:
				b.MqttClient.PublishMessage(payload)
				log.Debug("send payload data from outgoing function.")
			}
		}
	}()
	log.Debugf("device :%s run outgoing() method.", b.prod)
}

// NewEufyDevice 根据 product code 来创建 eufy 设备
func NewEufyDevice(clientid, username, pwd, broker, prodCode, devKey string, needCA bool) Eufydevice {
	var dev Eufydevice
	switch prodCode {
	case "T1011", "T1012":
		dev = newWhiteLight(clientid, username, pwd, broker, prodCode, devKey, needCA)
	case "T1013", "T1604":
		dev = newColorLight(clientid, username, pwd, broker, prodCode, devKey, needCA)
	case "T2103":
		dev = newRobotCleaner(clientid, username, pwd, broker, prodCode, devKey, needCA)
	}
	return dev
}
