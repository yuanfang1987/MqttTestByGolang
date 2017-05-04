package httpapi

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"oceanwing/eufy/httpapi/userapi"
	"time"

	splJson "github.com/bitly/go-simplejson"
)

// MyHTTPClient hh.
type MyHTTPClient struct {
	client        *http.Client
	Timezone      string
	Country       string
	Language      string
	Openudid      string
	ContentType   string
	Token         string
	UID           string
	ApplicationID string
	Req           chan *http.Request
	res           chan []byte
	// loginResp     chan *http.Request
	// obound        chan *http.Request
	// ibound        chan *http.Request
}

// NewMyHTTPClient create a new instance.
func NewMyHTTPClient(timeZone, country, language, openudid string) *MyHTTPClient {
	return &MyHTTPClient{
		Timezone:    timeZone,
		Country:     country,
		Language:    language,
		Openudid:    openudid,
		ContentType: "application/json;charset=utf-8",
		client:      &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *MyHTTPClient) outgoing() {
	go func() {
		for {
			select {
			case req := <-c.Req:
				resp, err := c.client.Do(req)
				if err == nil && resp.StatusCode == 200 {
					defer resp.Body.Close()
					bd, err := ioutil.ReadAll(resp.Body)
					if err == nil {
						c.res <- bd
					}
				}
			}
		}
	}()
}

func (c *MyHTTPClient) handleResponce() {
	go func() {
		for {
			select {
			case r := <-c.res:
				jsonInstance, err := splJson.NewJson(r)
				if err != nil {
					continue
				}
				resCode, err := jsonInstance.Get("res_code").Int()
				if err == nil && resCode == 1 {
					uid, _ := jsonInstance.Get("res_code").String()
					c.UID = uid
				}
				// ee
				// var v interface{}
				// err := json.Unmarshal(r, v)
				// if err == nil {

				// 	m := v.(map[string]interface{})
				// 	resCode, ok := m["res_code"]
				// 	if ok && resCode == 1 {
				// 		// get UID
				// 		if uid, ok := m["user_id"]; ok {
				// 			c.UID = uid.(string)
				// 		}
				// 		// get token
				// 		if token, ok := m["access_token"]; ok {
				// 			c.Timezone = token.(string)
				// 		}

				// 	}
				// }
			}
		}
	}()
}

// UserLogin hh.
func (c *MyHTTPClient) UserLogin(cid, csecret, email, pwd string) {
	l := userapi.NewLoginReq(cid, csecret, email, pwd)
	data, err := json.Marshal(l)
	if err != nil {
		return
	}
	body := bytes.NewBuffer(data)
	req, _ := http.NewRequest("POST", "http://zhome-ci.eufylife.com/v1/user/email/login", body)
	req.Header.Add("timezone", c.Timezone)
	req.Header.Add("country", c.Country)
	req.Header.Add("language", c.Language)
	req.Header.Add("openudid", c.Openudid)
	req.Header.Add("Content-Type", c.ContentType)
	c.Req <- req
}
