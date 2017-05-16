package mqttclient

import (
	"oceanwing/commontool"
	"oceanwing/config"
	"time"

	log "github.com/cihub/seelog"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// MqttClient haha.
type MqttClient struct {
	Clientid         string
	Username         string
	Pwd              string
	Broker           string
	client           MQTT.Client
	PubTopic         string
	SubTopic         string
	ReceivedPayload  []byte
	NeedCA           bool
	IsFisrtSubscribe bool
	tokenTimeout     time.Duration
	SubHandler       MQTT.MessageHandler
}

// ConnectToBroker hh.
func (n *MqttClient) ConnectToBroker() {
	// debug
	n.IsFisrtSubscribe = true
	n.tokenTimeout = time.Duration(config.GetInt(config.MqttTokenWaitTimeout))
	// set up client options.
	opt := MQTT.NewClientOptions()
	opt.AddBroker(n.Broker)
	opt.SetCleanSession(true)
	opt.SetAutoReconnect(true)
	opt.SetConnectTimeout(time.Duration(config.GetInt(config.MqttConnTimeout)) * time.Second)
	opt.SetKeepAlive(time.Duration(config.GetInt(config.MqttKeepAlive)) * time.Second)
	opt.SetWriteTimeout(time.Duration(config.GetInt(config.MqttWriteTimeout)) * time.Second)
	opt.SetMaxReconnectInterval(time.Duration(config.GetInt(config.MqttmaxReconnectInterval)) * time.Second)
	// onConnectionLost.
	myConnectionLostHandler := func(c MQTT.Client, e error) {
		log.Warnf("Connection Lost, ClientID: %s, ErrMsg: %s", n.Clientid, e.Error())
	}
	opt.SetConnectionLostHandler(myConnectionLostHandler)
	// onConnecthandler.
	myOnConnectHandler := func(c MQTT.Client) {
		// debug
		for {
			token := c.Subscribe(n.SubTopic, byte(1), n.SubHandler)
			if ret := token.WaitTimeout(n.tokenTimeout * time.Second); ret {
				if token.Error() == nil {
					if n.IsFisrtSubscribe {
						commontool.SubSinal <- struct{}{}
						n.IsFisrtSubscribe = false
					}
					log.Infof("subscribe to broker success with ClientID: %s", n.Clientid)
					break
				} else {
					log.Errorf("subscribe fail, ClientID: %s, ErrMsg: %s", n.Clientid, token.Error())
				}
			} else {
				log.Warnf("subsribe timeout, Client ID: %s", n.Clientid)
			}
		}
	}
	opt.SetOnConnectHandler(myOnConnectHandler)
	// haha, I'm fine.
	if n.Clientid != "" {
		opt.SetClientID(n.Clientid)
	}
	if n.Username != "" {
		opt.SetUsername(n.Username)
	}
	if n.Pwd != "" {
		opt.SetPassword(n.Pwd)
	}
	if n.NeedCA {
		opt.SetTLSConfig(commontool.CA.Clone())
	}
	// debug
	for {
		cc, err := createNewClient(opt)
		if err == nil && cc != nil {
			n.client = cc
			log.Infof("connect to broker success with clientID: %s", n.Clientid)
			break
		}
	}
}

func createNewClient(options *MQTT.ClientOptions) (MQTT.Client, error) {
	c := MQTT.NewClient(options)
	var err error
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		err = token.Error()
		log.Warnf("fail connect to broker, will try again: %s", err.Error())
		return nil, err
	}
	log.Debugf("create mqtt client success with clientid: %s", options.ClientID)
	return c, nil
}

// PublishMessage hh.
func (n *MqttClient) PublishMessage(payload interface{}) {
	token := n.client.Publish(n.PubTopic, byte(1), false, payload)
	if ret := token.WaitTimeout(n.tokenTimeout * time.Second); ret {
		if token.Error() != nil {
			log.Errorf("- - - - - publish msg fail: %s, %s", n.Clientid, token.Error().Error())
		}
	} else {
		log.Warnf("publish message timeout, clientID: %s", n.Clientid)
	}
	log.Debugf("Publish message successfully to topic: %s", n.PubTopic)
}

// PublishMessageWithNoTimout do not wait the token return timeout.
func (n *MqttClient) PublishMessageWithNoTimout(payload interface{}) {
	n.client.Publish(n.PubTopic, byte(1), false, payload)
}

// SubcribeToBroker hh.
// func (n *MqttClient) SubcribeToBroker(handler MQTT.MessageHandler) {
// 	n.SubHandler = handler
// 	token := n.client.Subscribe(n.SubTopic, byte(1), n.SubHandler)
// 	token.WaitTimeout(10 * time.Second)
// 	if token.Error() == nil {
// 		log.Info("subscribe to topic: ", n.SubTopic)
// 	} else {
// 		log.Info("subscribe fail", n.SubTopic, ": ", token.Error())
// 	}
// }
