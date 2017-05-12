package main

import (
	"flag"
	"oceanwing/commontool"
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
	// get parameter from command line..
	needca := flag.Bool("needca", true, "determine where need a CA file.")
	capath := flag.String("capath", "./EUFY_GD_CA.crt", "CA file path.")
	filePath := flag.String("filePath", "", "the path to the test data file.")
	broker := flag.String("broker", "ssl://pt-broker.eufylife.com:443", "determine which broker to connect to.")
	interval := flag.Int64("interval", 20, "device heart beat interval in seconds.")
	flag.Parse()
	// debug.
	log.Info("================ Start a new Test with parameters below... ================")
	log.Info("needca: ", *needca)
	log.Info("filePath: ", *filePath)
	log.Info("broker: ", *broker)
	log.Info("interval: ", *interval)
	// get test data.
	robotList, err := commontool.ReadFileContent(*filePath)
	if err != nil {
		log.Info("Read test data file error: ", err)
		return
	}

	// set up CA
	if *needca {
		commontool.BuildTlSConfig(*capath)
	}

	// run test..
	for i, v := range robotList {
		go func() {
			// str[0] deviceKey, str[1] password.
			str := strings.Split(v, ",")
			robotCleaner := robot.NewRobotCleaner()
			// CI: ssl://zhome-ci.eufylife.com:8893   PER: ssl://pt-broker.eufylife.com:443
			robotCleaner.RunRobotCleanerMqttService(str[0], str[0], "oceanwingtest", *broker, str[0], *needca)
			// send first heart beat.
			robotCleaner.SendRobotCleanerHeartBeat()
			//timer.
			heartBeatInterval := time.NewTicker(time.Second * time.Duration(*interval)).C
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
