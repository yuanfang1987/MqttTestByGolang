package device

import (
	"fmt"
	"oceanwing/eufy/restfulapi/base"

	splJSON "github.com/bitly/go-simplejson"
	log "github.com/cihub/seelog"
)

type device struct {
	base.BasicAPI
}

// NewDevice hh.
func NewDevice(category, url, method string, data, resultmap map[string]string) base.RESTfulAPI {
	d := &device{}
	d.APICategory = category
	d.APIURL = url
	d.APIMethod = method
	d.DataMap = data
	d.ResultMap = resultmap
	d.GetAPIName()
	log.Info("Build a new Device API.")
	return d
}

// BuildRequestBody 实现 BaseAPI 接口
func (d *device) BuildRequestBody() (string, []byte) {
	var bd []byte
	actURL := ""
	switch d.APIName {
	case "device":
		bd = d.device()
	case "":
		// bd =
	}
	return actURL, bd
}

// DecodeAndAssertResult 实现 BaseAPI 接口
func (d *device) DecodeAndAssertResult(resultJSON *splJSON.Json) map[string]string {
	var resultStr string
	switch d.APIName {
	case "device":
		resultStr = d.deviceResponse(resultJSON)
	}
	d.ResultMap["resultString"] = resultStr
	return d.ResultMap
}

// ===============================================================================================================================
func (d *device) device() []byte {
	var bd []byte
	var formatString, jsonString string
	switch d.APIMethod {
	case "GET":
	case "POST":
		formatString = `{
  				"alias_name": "%s",
  				"connect_type": %s,
  				"device_id": "%s",
  				"room_id": "%s"
			}`
		jsonString = fmt.Sprintf(formatString, d.DataMap["alias_name"], d.DataMap["connect_type"], d.DataMap["device_id"], d.DataMap["room_id"])
		bd = []byte(jsonString)
	case "PUT":
		formatString = `{
  				"device_key": "string",
  				"lan_ip_addr": "string",
  				"mac_address": "string",
  				"wifi_mac_address": "string",
  				"wifi_ssid": "string"
			}`
		jsonString = fmt.Sprintf(formatString, d.DataMap["device_key"], d.DataMap["lan_ip_addr"], d.DataMap["mac_address"],
			d.DataMap["wifi_mac_address"], d.DataMap["wifi_ssid"])
		bd = []byte(jsonString)
	}
	return bd
}

func (d *device) deviceResponse(resultJSON *splJSON.Json) string {
	resultStr := d.BaseResponse(resultJSON)
	if resultStr != "" {
		return resultStr
	}

	switch d.APIMethod {
	case "GET":
		resultStr = d.deviceResponseGet(resultJSON)
	case "POST":
		resultStr = d.deviceResponsePost(resultJSON)
	case "PUT":
		resultStr = d.deviceResponsePut(resultJSON)
	}
	return resultStr
}

func (d *device) deviceResponseGet(resultJSON *splJSON.Json) string {
	var resultStr string

	devices, err := resultJSON.Get("devices").Array()
	if err != nil {
		resultStr = fmt.Sprintf("get devices (array) fail: %s", err)
	} else if len(devices) != 0 {
		for _, dev := range devices {
			if devMap, ok := dev.(map[string]interface{}); ok {
				name := devMap["name"]
				sn := devMap["sn"]
				log.Infof("device info, name: %v, sn: %v", name, sn)
				// to do , other fields......
			}
		}
	}

	groups, err := resultJSON.Get("groups").Array()
	if err != nil {
		resultStr += ";" + fmt.Sprintf("get device group (array) fail: %s", err)
	} else if len(groups) != 0 {
		for _, group := range groups {
			if groupMap, ok := group.(map[string]interface{}); ok {
				groupName := groupMap["group_name "]
				log.Infof("gourp name: %v", groupName)
				// to do , other fields......
			}
		}
	}

	return resultStr
}

func (d *device) deviceResponsePost(resultJSON *splJSON.Json) string {
	var resultStr string
	device := resultJSON.Get("device")
	if device == nil {
		resultStr = "get device info fail."
		log.Info(resultStr)
		return resultStr
	}
	name, _ := device.Get("name").String()
	sn, _ := device.Get("sn").String()
	log.Infof("device info, name: %s, sn: %s", name, sn)
	return resultStr
}

func (d *device) deviceResponsePut(resultJSON *splJSON.Json) string {
	var resultStr string
	device := resultJSON.Get("device")
	if device == nil {
		resultStr = "get device info fail."
		log.Info(resultStr)
		return resultStr
	}
	name, _ := device.Get("name").String()
	sn, _ := device.Get("sn").String()
	log.Infof("device info, name: %s, sn: %s", name, sn)
	return resultStr
}
