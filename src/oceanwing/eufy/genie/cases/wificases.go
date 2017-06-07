package cases

import (
	log "github.com/cihub/seelog"
)

// RunWifiCases hh.
func RunWifiCases(url, wifi, pwd string) {
	log.Info("====== Running Wifi Test Cases ======")
	newTestInstance(url)

	Instance.GetAvailableWIFI()
	Instance.ConnectWifi(wifi, pwd)

	//Instance.ConnectToHideWifi(wifi, pwd)

	//Instance.SetHideSSID("1")
	//Instance.SetHideSSID("0")
}

// RunConnectToHideWIFI hh.
func RunConnectToHideWIFI(url, wifi, pwd string) {
	log.Info("===== Runing connecting hide wifi =======")
	newTestInstance(url)
	Instance.ConnectToHideWifi(wifi, pwd)
}
