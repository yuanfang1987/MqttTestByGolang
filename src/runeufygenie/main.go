package main

import (
	"oceanwing/commontool"
	"oceanwing/eufy/genie/cases"

	log "github.com/cihub/seelog"
)

func main() {
	commontool.InitLogInstance("debug")
	defer log.Flush()

	cases.RunMusicCases("http://10.10.10.254")
	cases.RunWifiCases("http://10.10.10.254", "OceanwingMobile", "0ceanwing11")
}
