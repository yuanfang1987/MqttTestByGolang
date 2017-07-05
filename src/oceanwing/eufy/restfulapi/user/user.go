package user

import (
	"fmt"
	"oceanwing/eufy/restfulapi/base"
	"strconv"

	splJSON "github.com/bitly/go-simplejson"
	log "github.com/cihub/seelog"
)

// User 是 user 模块的 API
type User struct {
	base.BasicAPI
}

// NewUser hh.
func NewUser(category, url, method string, data map[string]string) base.RESTfulAPI {
	u := &User{}
	u.APICategory = category
	u.APIURL = url
	u.APIMethod = method
	u.DataMap = data
	u.GetAPIName()
	log.Info("Build a new User API.")
	return u
}

// BuildRequestBody 实现 BaseAPI 接口
func (u *User) BuildRequestBody(data map[string]string) (string, []byte) {
	var bd []byte
	actURL := ""
	switch u.APIName {
	case "login":
		bd = u.login()
	}

	return actURL, bd
}

// DecodeAndAssertResult 实现 BaseAPI 接口
func (u *User) DecodeAndAssertResult(resultJSON *splJSON.Json) string {
	var resultStr string
	switch u.APIName {
	case "login":
		resultStr = u.loginResponse(resultJSON)
	}

	return resultStr
}

// ====================================================== login ===========================================================
func (u *User) login() []byte {
	formatString := `{
  		"client_id": "%s",
  		"client_secret": "%s",
  		"email": "%s",
  		"password": "%s"
	}`
	jsonString := fmt.Sprintf(formatString, u.DataMap["client_id"], u.DataMap["client_secret"], u.DataMap["email"], u.DataMap["password"])
	return []byte(jsonString)
}

func (u *User) loginResponse(resultJSON *splJSON.Json) string {
	resultStr := ""
	resCode, err := resultJSON.Get("res_code").Int()
	if err != nil {
		resultStr = fmt.Sprintf("login fail, get res_code error: %s", err)
		log.Error(resultStr)
		return resultStr
	}

	if strconv.Itoa(resCode) != u.DataMap["res_code"] {
		resultStr = fmt.Sprintf("login, assert res_code fail, expected: %s, actual: %d", u.DataMap["res_code"], resCode)
		log.Error(resultStr)
		return resultStr
	}

	// 如果预期结果是登录失败，则再判断一下 message 就可以， 剩下的无需判断
	if u.DataMap["target"] == "fail" {
		msg, _ := resultJSON.Get("message").String()
		if u.DataMap["message"] != msg {
			resultStr = fmt.Sprintf("login, assert message fail, expected: %s, actual: %s", u.DataMap["message"], msg)
			log.Error(resultStr)
		}
		return resultStr
	}

	email, _ := resultJSON.Get("email").String()
	if email != u.DataMap["email"] {
		resultStr = fmt.Sprintf("login, assert email fail, expected: %s, actual: %s", u.DataMap["email"], email)
		log.Error(resultStr)
	}

	uid, _ := resultJSON.Get("user_id").String()
	if uid != u.DataMap["user_id"] {
		resultStr = resultStr + ";" + fmt.Sprintf("login, assert uid fail, expected: %s, actual: %s", u.DataMap["user_id"], uid)
		log.Error(resultStr)
	}
	return resultStr
}
