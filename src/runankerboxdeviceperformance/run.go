package main

import (
	"oceanwing/ankerbox/business"
	"oceanwing/commontool"
	"oceanwing/config"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	log "github.com/cihub/seelog"
)

func main() {
	defer log.Flush()
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 初始化配置文件并获取参数值
	config.Initialize("config.yaml")
	// broker := config.GetString(config.MqttBroker)
	needca := config.GetBool(config.MqttNeedCA)
	devid := config.GetString(config.AnkerBoxDeviceID)
	sendSpeed := config.GetInt(config.AnkerBoxServerSendSpeed)
	cmdFile := config.GetString(config.AnkerBoxCommandFile)
	clientNum := config.GetInt(config.AnkerBoxClientNum)

	// set up CA
	if needca {
		capath := config.GetString(config.MqttCAFile)
		commontool.BuildTlSConfig(capath)
	}

	// 初始化日志实例
	commontool.InitLogInstance(config.GetString(config.LogLevel))

	// 读取指令文件
	cmds, err := commontool.ReadFileContent(cmdFile)
	if err != nil {
		log.Errorf("read txt file error: %s", err.Error())
		return
	}

	for i := 0; i < clientNum; i++ {
		go func() {
			d := business.NewAnkerBoxDevice(devid)
			d.RunMqttService()
			heartBeatInterval := time.NewTicker(time.Second * time.Duration(sendSpeed)).C
			index := 0
			for {
				select {
				case <-heartBeatInterval:
					d.SendCmdToDev(cmds[index])
					index++
					if index >= len(cmds) {
						index = 0
					}
				}
			}
		}()
		log.Infof("Start virtual mqtt instance... %d", i+1)
		<-commontool.SubSinal
	}

	channelSignal := make(chan os.Signal)
	signal.Notify(channelSignal, os.Interrupt)
	signal.Notify(channelSignal, syscall.SIGTERM)
	<-channelSignal
	log.Info("Test end")
}
