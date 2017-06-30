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

var (
	serverHost = "http://zhome-ci.eufylife.com"
)

// User 模拟一个Eufy用户
type User struct {
	httpclient      *http.Client
	req             chan *http.Request
	res             chan []byte
	jsonResult      chan *splJSON.Json
	EnableLeaveMode chan bool
	LeaveModeStart  time.Time
	LeaveModeEnd    time.Time
	email           string
	passWord        string
	clientID        string
	clientSecret    string
	uid             string
	token           string
}

// NewUser create a new user Instance.
func NewUser(email, pwd, clientid, clientse string) *User {
	u := &User{
		httpclient:      &http.Client{Timeout: 30 * time.Second},
		email:           email,
		passWord:        pwd,
		clientID:        clientid,
		clientSecret:    clientse,
		req:             make(chan *http.Request),
		res:             make(chan []byte),
		jsonResult:      make(chan *splJSON.Json),
		EnableLeaveMode: make(chan bool),
	}
	u.outgoing()
	u.handleResponce()
	log.Debugf("Create a new User, email: %s, pwd: %s, clientid: %s, client secret: %s", email, pwd, clientid, clientse)
	return u
}

// outgoing 是所有HTTP 请求的出口
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
					log.Debug("send response data to channel res")
				} else {
					log.Errorf("request to server fail: %s", err)
					user.res <- nil
				}
			}
		}
	}()
	log.Debug("Run outgoing() function.")
}

// handleResponce 用于统一处理HTTP请求的返回数据，主要是转为simple JSON的对象
func (user *User) handleResponce() {
	go func() {
		for {
			select {
			case resp := <-user.res:
				JSONJinstance, err := splJSON.NewJson(resp)
				if err != nil {
					user.jsonResult <- nil
					log.Errorf("decode the res data to simpleJSON error: %s", err)
				} else {
					user.jsonResult <- JSONJinstance
					log.Debug("send simpleJSON to jsonResult")
				}
			}
		}
	}()
	log.Debug("Run handleResponce() function")
}

// 构建请求体，包括必要的 header及 post body
func (user *User) buildRequest(method, url string, bd []byte) *http.Request {
	body := bytes.NewBuffer(bd)
	req, err := http.NewRequest(method, serverHost+url, body)
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
	log.Debugf("build a request with method: %s, and url: %s", method, url)
	return req
}

// Login hh.
func (user *User) Login() {
	data := buildLoginData(user.email, user.passWord, user.clientID, user.clientSecret)
	request := user.buildRequest("POST", "/v1/user/email/login", data)
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

// SetAwayMode hh..
func (user *User) SetAwayMode(beginMinute, finishedMinute int, deviceid string) {
	startTime := time.Now().Add(time.Duration(beginMinute) * time.Minute)
	endTime := time.Now().Add(time.Duration(finishedMinute) * time.Minute)

	user.LeaveModeStart = startTime
	user.LeaveModeEnd = endTime

	starthour := startTime.Hour()
	startminute := startTime.Minute()

	endhour := endTime.Hour()
	endminute := endTime.Minute()

	data := buildWeeklyRepeatAwayModeData("true", strconv.Itoa(starthour), strconv.Itoa(startminute), strconv.Itoa(endhour),
		strconv.Itoa(endminute), deviceid)

	request := user.buildRequest("POST", "/v1/away/save-timer", data)
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
	request := user.buildRequest("GET", "/v1/away/"+deviceid+"/get-timer", nil)
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

	log.Infof("decode away mode info, enabled: %t, schedule_type: %s, start time: %d:%d, end time: %d:%d, week info: %v", isEnable, scheduleType, starthour, startminute,
		endhour, endminute, weekinfo)

	user.EnableLeaveMode <- isEnable
}

// StopAwayMode 停止离家模式
func (user *User) StopAwayMode(devid string) {
	data := buildStopAwayModeData(devid)
	request := user.buildRequest("POST", "/v1/away/stop-timer", data)
	if request == nil {
		log.Error("cancel stop away mode")
		return
	}
	user.req <- request
	stopAwayModeJSON := <-user.jsonResult
	if stopAwayModeJSON == nil {
		log.Error("get stop away mode response error.")
		return
	}
	//取值，解析
	resCode, err := stopAwayModeJSON.Get("res_code").Int()
	if err != nil && resCode != 1 {
		msg, _ := stopAwayModeJSON.Get("message").String()
		log.Errorf("stop away mode failt, error message: %s", msg)
		return
	}
	log.Infof("stop away mode success, res_code :%d", resCode)
}