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
	Instance.QueryConnectStatus()
}
