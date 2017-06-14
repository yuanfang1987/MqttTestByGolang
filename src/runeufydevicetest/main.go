package main

import (
	"oceanwing/commontool"
	"oceanwing/config"
	"oceanwing/eufy/result"
	"oceanwing/eufy/serverpoint"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	log "github.com/cihub/seelog"
)

func main() {
	defer log.Flush()

	// 初始化配置文件
	config.Initialize("config.yaml")

	// get parameter form config.yaml file
	clientIDUserName := config.GetString(config.MqttauthUserName)
	password := config.GetString(config.MqttauthPassword)
	broker := config.GetString(config.MqttBroker)
	needca := config.GetBool(config.MqttNeedCA)
	codeKeys := config.GetString(config.EufyDeviceCodekeys)
	interval := config.GetInt(config.EufyDeviceSendCmdInterval)

	// 初始日志实例
	commontool.InitLogInstance(config.GetString(config.LogLevel))

	log.Info("=========================== Starting a new Eufy Device functional testing =========================")
	log.Infof("broker: %s", broker)
	log.Infof("interval: %d", interval)
	log.Infof("device key: %s", codeKeys)

	// 把CA文件加载到内存，供全局使用
	if needca {
		capath := config.GetString(config.MqttCAFile)
		commontool.BuildTlSConfig(capath)
	}

	// 新建csv文件用于存放测试结果
	result.NewResultFile("./EufyDeviceFuntionalTest.csv")
	defer result.CloseResultFile()

	go func() {
		// create a new cleaner
		eufyServer := serverpoint.NewMqttServerPoint()
		allDevices := strings.Split(codeKeys, ",")
		eufyServer.SetupRunningDevices(allDevices)
		// run mqtt service.
		eufyServer.RunMqttService(clientIDUserName, clientIDUserName, password, broker, needca)
		//timer.
		heartBeatInterval := time.NewTicker(time.Second * time.Duration(interval)).C
		for {
			select {
			case <-heartBeatInterval:
				eufyServer.PublishMsgToBroker()
			}
		}
	}()

	channelSignal := make(chan os.Signal)
	signal.Notify(channelSignal, os.Interrupt)
	signal.Notify(channelSignal, syscall.SIGTERM)
	<-channelSignal
}
