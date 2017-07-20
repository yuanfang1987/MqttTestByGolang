package main

import (
	"oceanwing/commontool"
	"oceanwing/eufy/result"
	"oceanwing/roavcam"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/cihub/seelog"
)

func main() {
	defer log.Flush()
	// 初始化日志实例
	commontool.InitLogInstance("info")

	// create a excel
	result.InitExcelFile()
	columNames := []string{"Prev File Name", "Time", "Next File Name", "Time", "Result"}
	result.WriteToExcel(columNames)

	go func() {
		roavcam.RunService()
		interval1 := time.NewTicker(time.Second * 5).C
		interval2 := time.NewTicker(time.Second * 600).C
		for {
			select {
			case <-interval1:
				roavcam.SendHeartBeat()
			case <-interval2:
				roavcam.GetFileListAndAssert()
			}
		}
	}()

	channelSignal := make(chan os.Signal)
	signal.Notify(channelSignal, os.Interrupt)
	signal.Notify(channelSignal, syscall.SIGTERM)
	<-channelSignal
	result.SaveExcelFile()
}
