package genie

import (
	"net/http"

	log "github.com/cihub/seelog"
)

// NewEufyGenie return a new instance.
func NewEufyGenie(url string) *BaseEufyGenie {
	return &BaseEufyGenie{
		client:  &http.Client{},
		baseURL: url, // "http://10.10.10.254"
	}
}

// GetPlayerStatus hh.
func (b *BaseEufyGenie) GetPlayerStatus(key, expValue string) {
	b.sendGet("/httpapi.asp?command=getPlayerStatus")
	myJSON := b.convertJSON()
	str, err := myJSON.Get(key).String()
	if err != nil {
		return
	}
	if str != expValue {
		log.Infof("verify that play status [%s] should be %s, but actual value is: %s", key,
			expValue, str)
	} else {
		log.Infof("Test passed, play status [%s] is %s", key, str)
	}
}
