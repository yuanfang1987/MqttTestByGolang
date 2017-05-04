package robot

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	SimpleJSON "github.com/bitly/go-simplejson"
	log "github.com/cihub/seelog"
)

// LoginReq for login data..
type LoginReq struct {
	Clientid     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Email        string `json:"email"`
	Password     string `json:"password"`
}

func getLoginData(cid, csecret, email, pwd string) io.Reader {
	o := &LoginReq{
		Clientid:     cid,
		ClientSecret: csecret,
		Email:        email,
		Password:     pwd,
	}
	data, err := json.Marshal(o)
	if err != nil {
		log.Warn("build login data fail, err: %s", err.Error())
		return nil
	}
	body := bytes.NewBuffer(data)
	return body
}

// AppUser 模拟一个用户登录，调用发指令给robot的接口
type AppUser struct {
	httpClient   *http.Client
	UID          string
	Token        string
	LoginResCode int
	cmdIndex     int
}

// NewAppUser 返回一个AppUser实例，并初始化http client
func NewAppUser() *AppUser {
	return &AppUser{
		httpClient: &http.Client{},
		cmdIndex:   0,
	}
}

// Login hh.
func (a *AppUser) Login(cid, csecret, email, pwd string) {
	reqBody := getLoginData(cid, csecret, email, pwd)
	req, _ := http.NewRequest("POST", "http://zhome-ci.eufylife.com/v1/user/email/login", reqBody)
	req.Header.Add("timezone", "Asia/Shanghai")
	req.Header.Add("country", "CN")
	req.Header.Add("language", "zh-Hans-CN")
	req.Header.Add("openudid", "421DF55D-897F-453D-B970-4674E583DACC")
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	resp, err := a.httpClient.Do(req)
	if err != nil {
		log.Errorf("Login fail, errMsg: %s", err.Error())
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Errorf("response code error: %d", resp.StatusCode)
		return
	}
	bd, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("read response body fail: %s", err.Error())
		return
	}
	JSONInstance, err := SimpleJSON.NewJson(bd)
	if err != nil {
		log.Errorf("unmarshal response data fail: %s", err.Error())
		return
	}
	rescode, err := JSONInstance.Get("res_code").Int()
	if rescode != 1 || err != nil {
		log.Errorf("Login fail, res_code is: %d", rescode)
		return
	}
	a.LoginResCode = rescode
	a.UID, _ = JSONInstance.Get("user_id").String()
	a.Token, _ = JSONInstance.Get("access_token").String()
	log.Infof("Login success with email: %s, password: %s", email, pwd)
}

// SendCmdToServer hh.
func (a *AppUser) SendCmdToServer(deviceid string) {
	if a.LoginResCode != 1 {
		log.Error("User not login success before!")
		return
	}
	body := a.getRequestBody()
	req, _ := http.NewRequest("POST", "http://zhome-ci.eufylife.com/v1/action/"+deviceid+"/action_robot_cleaner", body)
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	req.Header.Add("uid", a.UID)
	req.Header.Add("token", a.Token)
	resp, err := a.httpClient.Do(req)
	if err != nil {
		log.Errorf("Send command fail with http request error: %s", err.Error())
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Errorf("response code error: %d", resp.StatusCode)
		return
	}
	bd, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("read response body fail: %s", err.Error())
		return
	}
	JSONInstance, err := SimpleJSON.NewJson(bd)
	if err != nil {
		log.Errorf("UnMarshal response data fail, error: %s", err.Error())
		return
	}
	resCode, err := JSONInstance.Get("res_code").Int()
	if resCode == 0 && err != nil {
		log.Errorf("Send command fail due to response code: %d", resCode)
	} else if resCode == 2001 {
		log.Warn("Command not execute fail: Device Action Result is pending")
	} else if resCode == 2002 {
		log.Error("Device Not Owned by current User")
	} else {
		log.Infof("Send command to robot success, response code: %d", resCode)
	}
}

// SetIndexForReturnHome hh.
func (a *AppUser) SetIndexForReturnHome() {
	a.cmdIndex = 5
	log.Info("set http post body as [返回充电]")
}

func (a *AppUser) getRequestBody() io.Reader {
	cmd := []byte(comamndList[a.cmdIndex])
	log.Infof("ready to execute command: %s", commandCotentList[a.cmdIndex])
	a.cmdIndex++
	if a.cmdIndex > 8 {
		a.cmdIndex = 0
	}
	return bytes.NewBuffer(cmd)
}

const (
	// 定点 + 日常
	command1 = `{
		"clean_speed": {
			"speed": 0
		},
		"find_me": {
			"on_off": 0
		},
		"stop_clean": {
			"stop": 0
		},
		"work_mode": {
			"mode": 1
		}
	}`
	// 定点 + 强力
	command2 = `{
		"clean_speed": {
			"speed": 1
		},
		"find_me": {
			"on_off": 0
		},
		"stop_clean": {
			"stop": 0
		},
		"work_mode": {
			"mode": 1
		}
	}`
	// 定点 + 地毯
	command3 = `{
		"clean_speed": {
			"speed": 2
		},
		"find_me": {
			"on_off": 0
		},
		"stop_clean": {
			"stop": 0
		},
		"work_mode": {
			"mode": 1
		}
	}`
	// 定点 + 静音
	command4 = `{
		"clean_speed": {
			"speed": 3
		},
		"find_me": {
			"on_off": 0
		},
		"stop_clean": {
			"stop": 0
		},
		"work_mode": {
			"mode": 1
		}
	}`
	// 自动 + 日常
	command5 = `{
		"clean_speed": {
			"speed": 0
		},
		"find_me": {
			"on_off": 0
		},
		"stop_clean": {
			"stop": 0
		},
		"work_mode": {
			"mode": 2
		}
	}`
	// 返回充电
	command6 = `{
		"clean_speed": {
			"speed": 0
		},
		"find_me": {
			"on_off": 0
		},
		"stop_clean": {
			"stop": 0
		},
		"work_mode": {
			"mode": 3
		}
	}`
	// 沿边 + 日常
	command7 = `{
		"clean_speed": {
			"speed": 0
		},
		"find_me": {
			"on_off": 0
		},
		"stop_clean": {
			"stop": 0
		},
		"work_mode": {
			"mode": 4
		}
	}`
	// 精扫 + 日常
	command8 = `{
		"clean_speed": {
			"speed": 0
		},
		"find_me": {
			"on_off": 0
		},
		"stop_clean": {
			"stop": 0
		},
		"work_mode": {
			"mode": 5
		}
	}`
	// 暂停
	command9 = `{
		"clean_speed": {
			"speed": 0
		},
		"find_me": {
			"on_off": 0
		},
		"stop_clean": {
			"stop": 0
		},
		"work_mode": {
			"mode": 0
		}
	}`
)

var comamndList []string
var commandCotentList []string

func init() {
	comamndList = []string{command1, command2, command3, command4, command5, command6, command7, command8, command9}
	commandCotentList = []string{"定点 + 日常", "定点 + 强力", "定点 + 地毯", "定点 + 静音", "自动 + 日常", "返回充电",
		"沿边 + 日常", "精扫 + 日常", "暂停"}
}
