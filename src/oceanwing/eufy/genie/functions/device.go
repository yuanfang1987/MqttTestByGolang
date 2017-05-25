package functions

import (
	log "github.com/cihub/seelog"
)

// GetAndCheckDeviceInfo hh.  2.1
func (b *BaseEufyGenie) GetAndCheckDeviceInfo(key, value string) {
	b.sendGet("/httpapi.asp?command=getStatusEx")
	log.Info("execute get device info.")
	myJSON := b.convertJSON()
	actValue, err := myJSON.Get(key).String()
	if err != nil {
		log.Errorf("Get device info fail, key: %s, errMsg: %s", key, err.Error())
	}
	log.Infof("verify key: %s, expected value is: %s, actual value is: %s, test case passed or not? ---> %t",
		key, value, actValue, value == actValue)
}
