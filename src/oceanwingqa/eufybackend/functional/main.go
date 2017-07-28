package main

import (
	"oceanwingqa/common/config"
	"oceanwingqa/common/result"
	"oceanwingqa/common/utils"
	"oceanwingqa/eufybackend/functional/serverpoint"
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

	// 初始化日志实例
	utils.InitLogInstance(config.GetString(config.LogLevel))

	log.Info("=========================== Starting a new Eufy Device functional testing =========================")
	log.Infof("broker: %s", broker)
	log.Infof("interval: %d", interval)
	log.Infof("device key: %s", codeKeys)

	// 把CA文件加载到内存，供全局使用
	if needca {
		capath := config.GetString(config.MqttCAFile)
		utils.BuildTlSConfig(capath)
	}

	// 新建 Excel 文件存放测试结果
	result.InitExcelFile()
	columNames := []string{"Product Code", "Device Key", "TestCase Name", "Test Time", "Test Result", "Error Message"}
	result.WriteToExcel(columNames)

	go func() {
		// create a new eufyServer instance.
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
	log.Info("测试结束")
	result.SaveExcelFile()
}
