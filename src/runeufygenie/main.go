package main

import (
	"oceanwing/commontool"
	"oceanwing/eufy/genie/cases"
	"oceanwing/eufy/genie/results"

	log "github.com/cihub/seelog"
)

func main() {
	commontool.InitLogInstance("debug")
	defer log.Flush()

	results.NewResultFile("./eufyGenieTestResult.csv")
	defer results.CloseResultFile()

	cases.RunMusicCases("http://10.10.10.254")
	//cases.RunWifiCases("http://10.10.10.254", "OceanwingMobile", "0ceanwing11")
}
