package user

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	splJSON "github.com/bitly/go-simplejson"
	log "github.com/cihub/seelog"
)

// User 模拟一个Eufy用户
type User struct {
	httpclient   *http.Client
	req          chan *http.Request
	res          chan []byte
	jsonResult   chan *splJSON.Json
	email        string
	passWord     string
	clientID     string
	clientSecret string
	uid          string
	token        string
}

// NewUser create a new user Instance.
func NewUser(email, pwd, clientid, clientsecret string) *User {
	return &User{
		httpclient:   &http.Client{Timeout: 30 * time.Second},
		email:        email,
		passWord:     pwd,
		clientID:     clientid,
		clientSecret: clientsecret,
	}
}

func (user *User) outgoing() {
	go func() {
		for {
			select {
			case request := <-user.req:
				resp, err := user.httpclient.Do(request)
				if err == nil && resp.StatusCode == 200 {
					defer resp.Body.Close()
					bd, _ := ioutil.ReadAll(resp.Body)
					user.res <- bd
					log.Infof("send response data to channel res")
				} else {
					log.Errorf("request to server fail: %s", err)
					user.res <- nil
				}
			}
		}
	}()
}

func (user *User) handleResponce() {
	go func() {
		for {
			select {
			case resp := <-user.res:
				JSONJinstance, err := splJSON.NewJson(resp)
				if err != nil {
					user.jsonResult <- nil
				} else {
					user.jsonResult <- JSONJinstance
				}
			}
		}
	}()
}

func (user *User) buildRequest(method, url string, bd []byte) *http.Request {
	body := bytes.NewBuffer(bd)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Errorf("build new request error: %s", err)
		return nil
	}
	req.Header.Add("timezone", "Asia/Shanghai")
	req.Header.Add("country", "CN")
	req.Header.Add("language", "zh-Hans-CN")
	req.Header.Add("openudid", "yuanfang1987")
	req.Header.Add("category", "eufy-app")
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	if user.uid != "" {
		req.Header.Add("uid", user.uid)
	}
	if user.token != "" {
		req.Header.Add("token", user.token)
	}
	return req
}

// Login hh.
func (user *User) Login() {
	data := buildLoginData(user.email, user.passWord, user.clientID, user.clientSecret)
	request := user.buildRequest("POST", "http://zhome-ci.eufylife.com/v1/user/email/login", data)
	if request == nil {
		log.Error("Cancel login")
		return
	}
	user.req <- request

	loginJSON := <-user.jsonResult
	if loginJSON == nil {
		log.Error("login fail, nothing in the result json.")
		return
	}

	resCode, err := loginJSON.Get("res_code").Int()
	if err != nil && resCode != 1 {
		msg, _ := loginJSON.Get("message ").String()
		log.Errorf("login fail, res_code is: %d, error message: %s", resCode, msg)
		return
	}

	myuid, _ := loginJSON.Get("user_id").String()
	mytoken, _ := loginJSON.Get("access_token").String()
	myemail, _ := loginJSON.Get("email").String()

	user.uid = myuid
	user.token = mytoken

	log.Infof("login success, email: %s, uid: %s", myemail, myuid)
}

// SetAwayMode hh.
func (user *User) SetAwayMode(deviceid string) {
	startTime := time.Now().Add(3 * time.Minute)
	endTime := time.Now().Add(63 * time.Minute)

	starthour := startTime.Hour()
	startminute := startTime.Minute()

	endhour := endTime.Hour()
	endminute := endTime.Minute()

	weekinfo := startTime.Weekday()

	data := buildWeeklyRepeatAwayModeData(strconv.Itoa(int(weekinfo)), "true", strconv.Itoa(starthour), strconv.Itoa(startminute), strconv.Itoa(endhour),
		strconv.Itoa(endminute), deviceid)

	request := user.buildRequest("POST", "http://zhome-ci.eufylife.com/v1/away/save-timer", data)
	if request == nil {
		log.Error("cancel set away mode.")
		return
	}
	user.req <- request
	awayModeResultJSON := <-user.jsonResult
	if awayModeResultJSON == nil {
		log.Error("set away mode response error.")
		return
	}
	resCode, err := awayModeResultJSON.Get("res_code").Int()
	if err != nil && resCode != 1 {
		msg, _ := awayModeResultJSON.Get("message").String()
		log.Errorf("set away mode fail, res_code :%d, error msg: %s", resCode, msg)
		return
	}
	log.Infof("set away mode success, res_code: %d", resCode)
}

// GetAwayModeInfo hh.
func (user *User) GetAwayModeInfo(deviceid string) {
	request := user.buildRequest("GET", "http://zhome-ci.eufylife.com/v1/away/"+deviceid+"/get-timer", nil)
	if request == nil {
		log.Error("cancel get away mode info.")
		return
	}
	user.req <- request
	awayModeInfoJSON := <-user.jsonResult
	if awayModeInfoJSON == nil {
		log.Error("get away mode info response error.")
		return
	}
	// 解析，取值
	resCode, err := awayModeInfoJSON.Get("res_code").Int()
	if err != nil && resCode != 1 {
		msg, _ := awayModeInfoJSON.Get("message").String()
		log.Errorf("get away mode info fail, error message: %s", msg)
		return
	}

	isEnable, _ := awayModeInfoJSON.Get("away_timer").Get("enabled").Bool()
	scheduleType, _ := awayModeInfoJSON.Get("away_timer").Get("schedule_type").String()
	starthour, _ := awayModeInfoJSON.Get("away_timer").Get("start_hour").Int()
	startminute, _ := awayModeInfoJSON.Get("away_timer").Get("start_minute").Int()
	endhour, _ := awayModeInfoJSON.Get("away_timer").Get("end_hour").Int()
	endminute, _ := awayModeInfoJSON.Get("away_timer").Get("end_minute").Int()
	weekinfo, _ := awayModeInfoJSON.Get("away_timer").Get("away_repeat_option").Get("weekdays").Array()

	log.Infof("decode away mode info, enabled: %t, schedule_type: %s, start time: %d:%d, end time: %d:%d, week info %v", isEnable, scheduleType, starthour, startminute,
		endhour, endminute, weekinfo)
}
