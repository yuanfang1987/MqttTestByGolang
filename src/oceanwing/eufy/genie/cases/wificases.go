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

	Instance.SetHideSSID("1")
	Instance.SetHideSSID("0")
}
