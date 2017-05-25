package functions

import (
	"oceanwing/eufy/genie/results"

	log "github.com/cihub/seelog"
)

const cateDevice = "Device"

// GetAndCheckDeviceInfo hh.  2.1
func (b *BaseEufyGenie) GetAndCheckDeviceInfo(key, value string) {
	b.sendGet("/httpapi.asp?command=getStatusEx")
	log.Info("execute get device info.")
	myJSON := b.convertJSON()
	actValue, err := myJSON.Get(key).String()
	if err != nil {
		log.Errorf("Get device info fail, key: %s, errMsg: %s", key, err.Error())
	}
	re := passOrFail(value == actValue)
	log.Infof("verify key: %s, expected value is: %s, actual value is: %s, test case passed or not? ---> %s",
		key, value, actValue, re)
	results.WriteToResultFile(cateDevice, "check device info, key "+key+" value is "+value, re)
}
