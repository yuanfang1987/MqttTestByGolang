// 这是一个 robot-functional-test 分支，
// 仅用于对扫地机器人的功能稳定性测试

package main

import (
	"oceanwing/commontool"
	"oceanwing/config"
	"oceanwing/eufy/robot"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	log "github.com/cihub/seelog"
)

func main() {
	defer log.Flush()
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Initialize instance.
	config.Initialize("config.yaml")

	clientIDUserName := config.GetString(config.MqttauthUserName)
	password := config.GetString(config.MqttauthPassword)
	broker := config.GetString(config.MqttBroker)
	needca := config.GetBool(config.MqttNeedCA)
	devKeys := config.GetString(config.RobotcleanerDeviceKey)
	interval := config.GetInt(config.RobotcleanerHeartBeatInterval)
	appUserRunflag := config.GetBool(config.AppuserRunFlag)

	// 初始日志实例
	commontool.InitLogInstance(config.GetString(config.LogLevel))

	log.Infof("broker: %s", broker)
	log.Infof("interval: %d", interval)
	log.Infof("device key: %s", devKeys)

	// set up CA
	if needca {
		capath := config.GetString(config.MqttCAFile)
		commontool.BuildTlSConfig(capath)
	}

	// run test.
	for i := 0; i < 1; i++ {
		go func() {
			// create a new cleaner
			eufyServer := robot.NewEufyServer()
			allRobots := strings.Split(devKeys, ",")
			eufyServer.SetupRunningRobots(allRobots)
			// 模拟App用户
			if appUserRunflag {
				cid := config.GetString(config.AppuserClientid)
				csecret := config.GetString(config.AppuserClientscret)
				email := config.GetString(config.AppuserEmail)
				pwd := config.GetString(config.AppuserPassword)
				userDevKey := config.GetString(config.AppuserDevKey)
				userDeviceID := config.GetString(config.AppuserDevID)
				// user login
				eufyServer.SetAppUser(cid, csecret, email, pwd)
				eufyServer.AddRunningRobot(userDevKey, userDeviceID)
			}
			// run mqtt service
			eufyServer.RunMqtt(clientIDUserName, clientIDUserName, password, broker, needca)
			//timer.
			heartBeatInterval := time.NewTicker(time.Second * time.Duration(interval)).C
			for {
				select {
				case <-heartBeatInterval:
					eufyServer.PublishMsgToAllRobot()
				}
			}
		}()
		log.Infof("RobotCleaner Functional Testing Running...%d", i+1)
	}

	channelSignal := make(chan os.Signal)
	signal.Notify(channelSignal, os.Interrupt)
	signal.Notify(channelSignal, syscall.SIGTERM)
	<-channelSignal
	robot.ShowSummaryResult()
}
