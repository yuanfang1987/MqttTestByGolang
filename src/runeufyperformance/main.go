package main

import (
	"flag"
	"oceanwing/commontool"
	"oceanwing/eufy/performance"
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

	var start, end int

	// 从命令行获取参数
	needca := flag.Bool("needca", true, "use CA to connect to broker")
	capath := flag.String("capath", "EUFY_GD_CA.crt", "ca path to build CA")
	broker := flag.String("broker", "ssl://pt-broker.eufylife.com:443", "broker URL")
	filePath := flag.String("filePath", "", "the path to the test data file.")
	interval := flag.Int("interval", 20, "device heart beat interval in seconds.")
	loglevel := flag.String("loglevel", "info", "log level")
	startIndex := flag.Int("startIndex", 0, "start index")
	endIndex := flag.Int("endIndex", 0, "end index.")
	flag.Parse()

	// 初始化日志实例
	commontool.InitLogInstance(*loglevel)

	// 读取 Eufy 设备 device key 和密码
	eufyDevList, err := commontool.ReadFileContent(*filePath)
	if err != nil {
		log.Infof("Read test data file error: %s", err)
		return
	}

	// set up CA
	if *needca {
		commontool.BuildTlSConfig(*capath)
	}

	// print the test config info..
	log.Info("================ Start a new eufy device with parameters below... ================")
	log.Infof("needca: %t", *needca)
	log.Infof("broker: %s", *broker)
	log.Infof("use test file: %s", *filePath)
	log.Infof("set heart beat interval : %d", *interval)

	if *startIndex != 0 && *endIndex != 0 {
		start = *startIndex
		end = *endIndex
	} else {
		start = 0
		end = len(eufyDevList)
	}

	log.Infof("run from index %d to index %d", start, end)

	// run test.. i, v := range eufyDevList
	counter := 0
	for i := start; i < end; i++ {
		v := eufyDevList[i]
		go func() {
			// str[0] product code, str[1] device key, str[2] password
			str := strings.Split(v, ",")
			// new a eufy device.
			eufyDevice := performance.NewEufyDevice(str[1], str[1], "oceanwingtest", *broker, str[0], str[1], *needca)
			// CI: ssl://zhome-ci.eufylife.com:8893   PER: ssl://pt-broker.eufylife.com:443
			eufyDevice.RunMqttService()
			// send first heart beat.
			eufyDevice.SendHeartBeat()
			//timer.
			heartBeatInterval := time.NewTicker(time.Second * time.Duration(*interval)).C
			for {
				select {
				case <-heartBeatInterval:
					eufyDevice.SendHeartBeat()
				}
			}
		}()
		counter++
		log.Infof("Run eufy device: %d", counter)
		// --- debug
		<-commontool.SubSinal
	}
	channelSignal := make(chan os.Signal)
	signal.Notify(channelSignal, os.Interrupt)
	signal.Notify(channelSignal, syscall.SIGTERM)
	<-channelSignal
	log.Info("Heart beat end")
}
