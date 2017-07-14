package performance

import (
	"oceanwing/mqttclient"

	log "github.com/cihub/seelog"
)

// Eufydevice hh.
type Eufydevice interface {
	RunMqttService()
	SendHeartBeat()
}

type baseEufy struct {
	mqttclient.MqttClient
	msgToServer chan []byte
	prod        string
}

func (b *baseEufy) RunMqttService() {
	log.Error("Not implemet yet.")
}
func (b *baseEufy) SendHeartBeat() {
	log.Error("Not implement yet.")
}

// 这是每个设备发出消息的出口
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
	log.Debug("run outgoing method.")
}

// NewEufyDevice hh.
func NewEufyDevice(clientid, username, pwd, broker, prodCode, devKey string, needCA bool) Eufydevice {
	var dev Eufydevice
	switch prodCode {
	case "T1011", "T1012":
		dev = newWhiteLight(clientid, username, pwd, broker, prodCode, devKey, needCA)
	case "T1013", "T1604":
		dev = newColorLight(clientid, username, pwd, broker, prodCode, devKey, needCA)
	}
	return dev
}
