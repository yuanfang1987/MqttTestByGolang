package functions

import (
	"encoding/json"
	"strconv"

	"time"

	log "github.com/cihub/seelog"
)

var myWifi AvailableWIFIList

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
	bd := b.getBytesResult()
	err := json.Unmarshal(bd, &myWifi)
	if err != nil {
		log.Errorf("get wifi list and decode json fail: %s", err.Error())
		return
	}
	num, _ := strconv.Atoi(myWifi.Res)
	log.Infof("available wifi list number: %d", num)
	if num <= 0 {
		log.Error("Get wifi list fail, found nothing!")
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
			matchWifi = wifi
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
	// 忽略执行结果
	b.getStringResult()

	// 查询结果
	time.Sleep(10 * time.Second)
	for i := 0; i < 3; i++ {
		res := b.QueryConnectStatus()
		if res == "OK" {
			log.Infof("connect wifi %s success", wifiName)
			return
		}
		log.Debugf("current connect status is: %s, wait 10 second and then try query again.", res)
		time.Sleep(10 * time.Second)
	}
	log.Errorf("fail to connect to wifi [%s]after waiting for 30 seconds", wifiName)
}

// QueryConnectStatus hh.  3.4
func (b *BaseEufyGenie) QueryConnectStatus() string {
	b.sendGet("/httpapi.asp?command=wlanGetConnectState")
	strOK := b.getStringResult()
	log.Infof("query wifi connected status is: %s", strOK)
	return strOK
}
