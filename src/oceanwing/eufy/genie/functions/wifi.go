package functions

import (
	"encoding/json"
	"oceanwing/eufy/genie/results"
	"strconv"
	"time"

	log "github.com/cihub/seelog"
)

var myWifi AvailableWIFIList

const categoryWifi = "WIFI"

// AvailableWIFIList hh.
type AvailableWIFIList struct {
	Res    string      `json:"res"`
	Aplist []*WifiInfo `json:"aplist"`
}

// WifiInfo hh.
type WifiInfo struct {
	SSID    string `json:"ssid"`
	Bssid   string `json:"bssid"`
	Rssi    string `json:"rssi"`
	Channel string `json:"channel"`
	Auth    string `json:"auth"`
	Encry   string `json:"encry"`
	Extch   string `json:"extch"`
}

// GetAvailableWIFI hh.  3.1
func (b *BaseEufyGenie) GetAvailableWIFI() {
	b.sendGet("/httpapi.asp?command=wlanGetApListEx")
	log.Info("execute get available wifi list.")
	bd := b.getBytesResult()
	err := json.Unmarshal(bd, &myWifi)
	if err != nil {
		log.Errorf("get wifi list and decode json fail: %s", err.Error())
		results.WriteToResultFile(categoryWifi, "get available wifi list", "Fail")
		return
	}
	num, _ := strconv.Atoi(myWifi.Res)
	log.Infof("available wifi list number: %d", num)
	if num <= 0 {
		log.Error("Get wifi list fail, found nothing!")
		results.WriteToResultFile(categoryWifi, "get available wifi list", "Fail")
		return
	}
	for i, wifi := range myWifi.Aplist {
		log.Infof("=== Num.%d wifi info ===", i)
		log.Infof("ssid: %s", hexToString(wifi.SSID))
		log.Infof("bssid: %s", wifi.Bssid)
		log.Infof("rssi: %s", wifi.Rssi)
		log.Infof("channel: %s", wifi.Channel)
		log.Infof("auth: %s", wifi.Auth)
		log.Infof("encry: %s", wifi.Encry)
		log.Infof("extch %s", wifi.Extch)
	}
	results.WriteToResultFile(categoryWifi, "get available wifi list", "Pass")
}

// ConnectWifi hh. 3.2
func (b *BaseEufyGenie) ConnectWifi(wifiName, password string) {
	if myWifi.Res == "" {
		log.Infof("availabel wifi list is empty, re-scan now.....")
		b.GetAvailableWIFI()
	}
	var matchWifi *WifiInfo
	for _, wifi := range myWifi.Aplist {
		if hexToString(wifi.SSID) == wifiName {
			log.Infof("get matched wifi: %s", wifiName)
			log.Infof("Wifi channel: %s", wifi.Channel)
			log.Infof("Wifi auth: %s", wifi.Auth)
			log.Infof("wifi encry: %s", wifi.Encry)
			matchWifi = wifi
			break
		}
	}
	if matchWifi == nil {
		log.Errorf("the expeted wifi: [%s] not found in the available wifi list!", wifiName)
		return
	}
	ssid := matchWifi.SSID
	channel := matchWifi.Channel
	auth := matchWifi.Auth
	encry := matchWifi.Encry
	pwd := stringToHex(password)
	url := "/httpapi.asp?command=wlanConnectApEx:ssid=" + ssid + ":ch=" + channel + ":auth=" + auth + ":encty=" + encry + ":pwd=" + pwd + ":chext=1"
	b.sendGet(url)
	log.Infof("execute connect to wifi: %s", wifiName)
	// 忽略执行结果
	b.getStringResult()

	// 查询结果
	re := b.queryConnectStatus(wifiName)
	results.WriteToResultFile(categoryWifi, "Connect to wifi "+wifiName, re)
}

// ConnectToHideWifi hh. 3.3
func (b *BaseEufyGenie) ConnectToHideWifi(wifiName, password string) {
	var url string
	ssid := stringToHex(wifiName)
	if password != "" {
		pwd := stringToHex(password)
		url = "/httpapi.asp?command=wlanConnectHideApEx:" + ssid + ":" + pwd
	} else {
		url = "/httpapi.asp?command=wlanConnectHideApEx:" + ssid
	}
	b.sendGet(url)
	log.Infof("execute connect to hide wifi: %s", wifiName)
	// 忽略执行结果
	b.getStringResult()
	// 查询结果
	re := b.queryConnectStatus(wifiName)
	results.WriteToResultFile(categoryWifi, "Connect to hide wifi "+wifiName, re)
}

// queryConnectStatus hh.  3.4
func (b *BaseEufyGenie) queryConnectStatus(wifiName string) string {
	for i := 0; i < 3; i++ {
		// 每次执行新的连接后，需等10秒钟后再查询状态
		time.Sleep(10 * time.Second)
		b.sendGet("/httpapi.asp?command=wlanGetConnectState")
		res := b.getStringResult()
		if res == "OK" {
			log.Infof("connect wifi %s success", wifiName)
			return "Pass"
		}
		log.Debugf("current connect status is: %s, wait 10 second and then try query again.", res)
	}
	log.Errorf("fail to connect to wifi [%s]after waiting for 30 seconds", wifiName)
	return "Fail"
}

// SetHideSSID hide wifi.  3.5  x为1表示隐藏AP, x为0表示恢复AP
func (b *BaseEufyGenie) SetHideSSID(value string) {
	b.sendGet("/httpapi.asp?command=setHideSSID:" + value)
	log.Infof("execute set hide SSID to value: %s", value)
	strOK := b.getStringResult()
	log.Infof("set wifi hide status: %s and execute result is: %s, test case passed or not? ---> %t", value,
		strOK, strOK == "OK")
	// query status
	re := b.getHideSSID(value)
	results.WriteToResultFile(categoryWifi, "set hide ssid to "+value, re)
}

// 3.6
func (b *BaseEufyGenie) getHideSSID(expValue string) string {
	b.sendGet("/httpapi.asp?command=getHideSSID")
	log.Info("execute get hide SSID status.")
	myJSON := b.convertJSON()
	strOK, _ := myJSON.Get("hideSSID").String()
	re := passOrFail(strOK == expValue)
	log.Infof("current wifi hide status: %s, test case passed or not? ---> %s", strOK, re)
	return re
}
