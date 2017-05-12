package main

import (
	"flag"
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

	// 初始化配置文件
	config.Initialize("config.yaml")

	// 从配置文件获取参数
	needca := config.GetBool(config.MqttNeedCA)
	capath := config.GetString(config.MqttCAFile)
	broker := config.GetString(config.MqttBroker)

	// 从命令行获取参数
	filePath := flag.String("filePath", "", "the path to the test data file.")
	interval := flag.Int("interval", 0, "device heart beat interval in seconds.")
	flag.Parse()

	// 初始化日志实例
	commontool.InitLogInstance(config.GetString(config.LogLevel))

	var fpath string
	var inter int
	// 如果命令行有参数传进来，则不使用配置文件的值
	if *filePath != "" {
		fpath = *filePath
	} else {
		fpath = config.GetString(config.RobotcleanerTestDataFile)
	}

	if *interval != 0 {
		inter = *interval
	} else {
		inter = config.GetInt(config.RobotcleanerHeartBeatInterval)
	}

	// 读取扫地机器人相关测试数据
	robotList, err := commontool.ReadFileContent(fpath)
	if err != nil {
		log.Info("Read test data file error: ", err)
		return
	}

	// set up CA
	if needca {
		commontool.BuildTlSConfig(capath)
	}

	// print the test config info..
	log.Info("================ Start a new Test with parameters below... ================")
	log.Info("needca: ", needca)
	log.Info("broker: ", broker)
	log.Infof("use test file: %s", fpath)
	log.Infof("set heart beat interval : %d", inter)

	// run test..
	for i, v := range robotList {
		go func() {
			// str[0] deviceKey, str[1] password.
			str := strings.Split(v, ",")
			robotCleaner := robot.NewRobotCleaner()
			// CI: ssl://zhome-ci.eufylife.com:8893   PER: ssl://pt-broker.eufylife.com:443
			robotCleaner.RunRobotCleanerMqttService(str[0], str[0], "oceanwingtest", broker, str[0], needca)
			// send first heart beat.
			robotCleaner.SendRobotCleanerHeartBeat()
			//timer.
			heartBeatInterval := time.NewTicker(time.Second * time.Duration(inter)).C
			for {
				select {
				case <-heartBeatInterval:
					robotCleaner.SendRobotCleanerHeartBeat()
				}
			}
		}()
		log.Info("Run Robot Cleaner: ", i+1)
		// --- debug
		<-commontool.SubSinal
	}
	channelSignal := make(chan os.Signal)
	signal.Notify(channelSignal, os.Interrupt)
	signal.Notify(channelSignal, syscall.SIGTERM)
	<-channelSignal
	log.Info("Heart beat end")
}
