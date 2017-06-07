package main

import (
	"flag"
	"oceanwing/commontool"
	"oceanwing/eufy/genie/cases"
	"oceanwing/eufy/genie/results"

	log "github.com/cihub/seelog"
)

func main() {
	ssid := flag.String("ssid", "", "ssid")
	pwd := flag.String("pwd", "", "password")
	flag.Parse()
	commontool.InitLogInstance("debug")
	defer log.Flush()

	results.NewResultFile("./eufyGenieTestResult.csv")
	defer results.CloseResultFile()

	//cases.RunMusicCases("http://10.10.10.254")
	cases.RunWifiCases("http://10.10.10.254", *ssid, *pwd)
}
