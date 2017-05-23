package main

import (
	"oceanwing/commontool"
	"oceanwing/eufy/genie"

	log "github.com/cihub/seelog"
)

func main() {
	commontool.InitLogInstance("debug")
	defer log.Flush()
	myGenie := genie.NewEufyGenie("http://10.10.10.254")
	myGenie.GetPlayerStatus("status", "play")
}
