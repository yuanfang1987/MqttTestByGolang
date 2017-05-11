package business

import (
	"encoding/base64"
	"oceanwing/ankerbox/pb"
	"oceanwing/commontool"
	"oceanwing/config"
	"oceanwing/mqttclient"
	"strconv"
	"strings"

	log "github.com/cihub/seelog"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/protobuf/proto"
	"github.com/pborman/uuid"
)

// AnkerBoxDevice hh.
type AnkerBoxDevice struct {
	Devid string
	mqttclient.MqttClient
}

// NewAnkerBoxDevice return a new ankerboxdevice instance.
func NewAnkerBoxDevice(devid string) *AnkerBoxDevice {
	return &AnkerBoxDevice{
		Devid: devid,
	}
}

// RunMqttService hh.
func (a *AnkerBoxDevice) RunMqttService() {
	a.Broker = config.GetString(config.MqttBroker)
	a.Clientid = commontool.GenerateClientID()
	a.PubTopic = "SERTODEV/" + a.Devid
	a.SubTopic = "DEVTOSER/1"
	a.NeedCA = config.GetBool(config.MqttNeedCA)
	a.SubHandler = func(c MQTT.Client, msg MQTT.Message) {
		// do nothing.
	}
	a.MqttClient.ConnectToBroker()
}

// SendCmdToDev hh.
func (a *AnkerBoxDevice) SendCmdToDev(cmdParas string) {
	paras := strings.Split(cmdParas, ";")
	payload := buildOpenDevCommand(paras)
	a.MqttClient.PublishMessage(payload)
}

// buildOpenDevCommand hh. the para array must be like this:
// 1.2.3;1;borrow;6;yUjJzpujWwhTbBhgwSNjSA==
func buildOpenDevCommand(paras []string) []byte {
	// get parameter.
	version := paras[0]
	gpid, _ := strconv.Atoi(paras[1])
	action := paras[2]
	slotNum, _ := strconv.Atoi(paras[3])
	pwd := paras[4]

	// header
	uuid1 := uuid.NewUUID().String()
	head := &pb.CMsgHead{
		Cmd:     pb.CMD_OPENDEV.Enum(),
		Version: proto.String(version),
		Tranid:  proto.String(uuid1),
		Groupid: proto.Int32(int32(gpid)),
	}
	// action type.
	var expectedAction *pb.DEVACTION
	if action == "borrow" {
		expectedAction = pb.DEVACTION_BORROW.Enum()
	} else {
		expectedAction = pb.DEVACTION_RETURN.Enum()
	}
	// password.
	var passwd string
	if pwd != "" {
		aa, _ := base64.StdEncoding.DecodeString(pwd)
		passwd = string(aa)
	} else {
		passwd = ""
	}
	// open dev body.
	openDevBody := &pb.CMsgBodyDeviceOpen{
		Action:   expectedAction,
		Num:      proto.Int32(int32(slotNum)),
		Password: proto.String(passwd),
	}
	cmsg := &pb.CMsg{
		MsgHead:        head,
		OpenDeviceBody: openDevBody,
	}
	// marshal
	data, err := proto.Marshal(cmsg)
	if err != nil {
		log.Errorf("build open dev command fail: %s", err.Error())
		return nil
	}
	log.Debugf("Build command info successfully, version: %s, group id: %d, action: %s, slotNum: %d", version, gpid, action, slotNum)
	return data
}
