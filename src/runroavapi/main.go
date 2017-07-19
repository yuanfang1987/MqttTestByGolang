package main

import (
	"oceanwing/commontool"
	"oceanwing/eufy/result"
	"oceanwing/roavcam"
	"os"
	"os/signal"
	"syscall"

	log "github.com/cihub/seelog"
)

func main() {
	defer log.Flush()
	// 初始化日志实例
	commontool.InitLogInstance("info")

	// create a excel
	result.InitExcelFile()

	roavcam.SendRoavAPI()

	channelSignal := make(chan os.Signal)
	signal.Notify(channelSignal, os.Interrupt)
	signal.Notify(channelSignal, syscall.SIGTERM)
	<-channelSignal

	result.SaveExcelFile()
}
