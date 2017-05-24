package functions

import (
	"encoding/hex"
	"io/ioutil"
	"net/http"

	SimpleJSON "github.com/bitly/go-simplejson"
	log "github.com/cihub/seelog"
)

// BaseEufyGenie hh.
type BaseEufyGenie struct {
	client   *http.Client
	baseURL  string
	respBody chan []byte
}

// NewEufyGenie return a new instance.
func NewEufyGenie(url string) *BaseEufyGenie {
	return &BaseEufyGenie{
		client:   &http.Client{},
		baseURL:  url, // "http://10.10.10.254"
		respBody: make(chan []byte),
	}
}

func (b *BaseEufyGenie) sendGet(urlPath string) {
	go func() {
		bd, err := b.client.Get(b.baseURL + urlPath)
		if err != nil {
			log.Errorf("request fail, url: %s, error: %s", urlPath, err)
			b.respBody <- []byte{'N', 'O'}
			return
		}
		defer bd.Body.Close()
		bb, _ := ioutil.ReadAll(bd.Body)
		b.respBody <- bb
	}()
}

func (b *BaseEufyGenie) convertJSON() *SimpleJSON.Json {
	bb := <-b.respBody
	JSONInstance, err := SimpleJSON.NewJson(bb)
	if err == nil {
		return JSONInstance
	}
	return nil
}

func (b *BaseEufyGenie) getStringResult() string {
	bb := <-b.respBody
	return string(bb)
}

func hexToString(hexStr string) string {
	res, err := hex.DecodeString(hexStr)
	if err == nil {
		return string(res)
	}
	log.Errorf("Decode hex to string fail: %s", err)
	return ""
}
