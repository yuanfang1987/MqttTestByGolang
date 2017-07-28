package main

import (
	"flag"
	"oceanwingqa/common/utils"
	"oceanwingqa/eufybackend/restfulapi/runapiengine"

	log "github.com/cihub/seelog"
)

func main() {
	defer log.Flush()
	loglevel := flag.String("loglevel", "debug", "log level")
	host := flag.String("host", "http://zhome-ci.eufylife.com", "the base url")
	flag.Parse()
	// 初始化日志实例
	utils.InitLogInstance(*loglevel)
	runapiengine.SetHostURL(*host)
	runapiengine.Runapitest()

	// channelSignal := make(chan os.Signal)
	// signal.Notify(channelSignal, os.Interrupt)
	// signal.Notify(channelSignal, syscall.SIGTERM)
	// <-channelSignal
}
