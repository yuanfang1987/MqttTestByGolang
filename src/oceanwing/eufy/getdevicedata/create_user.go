package getdevicedata

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	splJSON "github.com/bitly/go-simplejson"
)

var (
	myClient     = &http.Client{Timeout: 30 * time.Second}
	serverHost   = "http://34.210.88.139:8080"
	registerBody = `{
  		"client_id": "eufy-app",
  		"client_secret": "8FHf22gaTKu7MZXqz5zytw",
  		"email": "%s",
  		"name": "yuan",
  		"password": "Aa123456",
  		"un_subscribe_flag": true
	}`
	bindBody = `{
  		"device_key": "%s"
	}`
)

type user struct {
	uid   string
	email string
	token string
}

func generateEmail() string {
	t := time.Now().Unix()
	tt := strconv.Itoa(int(t))
	email := "eufyTest" + tt + "@oceanwing.com"
	return email
}

// 注册用户
func registerUser() *user {
	jsonString := fmt.Sprintf(registerBody, generateEmail())
	req := buildRequest("", "", "POST", "/v1/user/email/register", []byte(jsonString))
	resp, err := myClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		fmt.Printf("register error: %s\n", err)
		return nil
	}
	defer resp.Body.Close()
	bd, _ := ioutil.ReadAll(resp.Body)
	JSONJinstance, err := splJSON.NewJson(bd)
	if err != nil {
		fmt.Printf("register user, convert to simple json error: %s\n", err)
		return nil
	}
	resCode, _ := JSONJinstance.Get("res_code").Int()
	if resCode != 1 {
		errMsg, _ := JSONJinstance.Get("message").String()
		fmt.Printf("register user fail: %s\n", errMsg)
		return nil
	}
	myuid, _ := JSONJinstance.Get("user_id").String()
	myemail, _ := JSONJinstance.Get("email").String()
	mytoken, _ := JSONJinstance.Get("access_token").String()
	return &user{
		uid:   myuid,
		email: myemail,
		token: mytoken,
	}
}

// 绑定设备
func bindDevice(devkey string, u *user) bool {
	jsonString := fmt.Sprintf(bindBody, devkey)
	req := buildRequest(u.uid, u.token, "PUT", "/v1/device/", []byte(jsonString))
	resp, err := myClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		fmt.Printf("bind device error: %s\n", err)
		return false
	}
	defer resp.Body.Close()
	bd, _ := ioutil.ReadAll(resp.Body)
	JSONJinstance, err := splJSON.NewJson(bd)
	resCode, _ := JSONJinstance.Get("res_code").Int()
	if resCode != 1 {
		return false
	}
	return true
}

// 构建请求体，包括必要的 header及 post body
func buildRequest(uid, token, method, url string, bd []byte) *http.Request {
	body := bytes.NewBuffer(bd)
	req, err := http.NewRequest(method, serverHost+url, body)
	if err != nil {
		fmt.Printf("build new request error: %s", err)
		return nil
	}
	req.Header.Add("timezone", "Asia/Shanghai")
	req.Header.Add("country", "CN")
	req.Header.Add("language", "zh-Hans-CN")
	req.Header.Add("openudid", "yuanfang1987")
	req.Header.Add("category", "eufy-app")
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	if uid != "" {
		req.Header.Add("uid", uid)
	}
	if token != "" {
		req.Header.Add("token", token)
	}
	fmt.Printf("build a request with method: %s, and url: %s", method, url)
	return req
}
