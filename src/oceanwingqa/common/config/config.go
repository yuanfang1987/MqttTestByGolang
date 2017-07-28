package config

import (
	"log"

	"github.com/olebedev/config"
)

const (
	LogLevel                 = "loglevel"
	MqttBroker               = "mqtt.broker"
	MqttNeedCA               = "mqtt.needCA"
	MqttCAFile               = "mqtt.server_ca_file"
	MqttClientCertFile       = "mqtt.client_cert_file"
	MqttClientCertKeyFile    = "mqtt.client_cert_key_file"
	MqttConnTimeout          = "mqtt.connection_timeout"
	MqttKeepAlive            = "mqtt.keep_alive"
	MqttWriteTimeout         = "mqtt.write_timeout"
	MqttmaxReconnectInterval = "mqtt.max_reconnect_interval"
	MqttTokenWaitTimeout     = "mqtt.token_wait_timeout"
	MqttauthUserName         = "mqtt.auth_user_name"
	MqttauthPassword         = "mqtt.auth_password"

	EufyDeviceSendCmdInterval = "eufydevice.send_command_interval"
	EufyDeviceCodekeys        = "eufydevice.code_key"
	EufyDeviceOnlyListenMsg   = "eufydevice.is_only_listen"

	RobotCleanerRunMode           = "robotcleaner.runMode"
	RobotcleanerHeartBeatInterval = "robotcleaner.heart_beat_interval"
	RobotcleanerTestDataFile      = "robotcleaner.test_data_file"
	RobotcleanerDeviceKey         = "robotcleaner.deviceKey"

	AppuserRunFlag     = "appuser.runflag"
	AppuserClientid    = "appuser.clientid"
	AppuserClientscret = "appuser.clientscret"
	AppuserEmail       = "appuser.email"
	AppuserPassword    = "appuser.password"
	AppuserDevID       = "appuser.devID"
	AppuserDevKey      = "appuser.devKey"

	AwayModeStart = "awaymode.start"
	AwayModeEnd   = "awaymode.end"
)

var (
	confInstance *config.Config
)

func ensureInitialize() {
	if confInstance == nil {
		log.Println("config not initialize")
	}
}

// Initialize hh.
func Initialize(configFile string) {
	if configFile == "" {
		configFile = "config.yaml"
	}
	conf, err := config.ParseYamlFile(configFile)
	if err != nil {
		log.Println(err)
		return
	}
	confInstance = conf.Env()
}

// GetString hh.
func GetString(key string) string {
	ensureInitialize()
	val, _ := confInstance.String(key)
	return val
}

// GetInt hh.
func GetInt(key string) int {
	ensureInitialize()
	val, _ := confInstance.Int(key)
	return val
}

// GetBool hh.
func GetBool(key string) bool {
	ensureInitialize()
	val, _ := confInstance.Bool(key)
	return val
}
