package user

import (
	"fmt"
	"oceanwing/commontool"
	"oceanwing/eufy/restfulapi/base"

	splJSON "github.com/bitly/go-simplejson"
	log "github.com/cihub/seelog"
)

// User 是 user 模块的 API
type User struct {
	base.BasicAPI
	registerEmail string
}

// NewUser hh.
func NewUser(category, url, method string, data, resultmap map[string]string) base.RESTfulAPI {
	u := &User{}
	u.APICategory = category
	u.APIURL = url
	u.APIMethod = method
	u.DataMap = data
	u.ResultMap = resultmap
	u.GetAPIName()
	log.Info("Build a new User API.")
	return u
}

// BuildRequestBody 实现 BaseAPI 接口.
func (u *User) BuildRequestBody() (string, []byte) {
	var bd []byte
	actURL := ""
	switch u.APIName {
	case "login":
		bd = u.login()
	case "register":
		bd = u.register()
	}

	return actURL, bd
}

// DecodeAndAssertResult 实现 BaseAPI 接口
func (u *User) DecodeAndAssertResult(resultJSON *splJSON.Json) map[string]string {
	var resultStr string
	switch u.APIName {
	case "login":
		resultStr = u.loginResponse(resultJSON)
	case "register":
		resultStr = u.registerResponse(resultJSON)
	}
	u.ResultMap["resultString"] = resultStr
	return u.ResultMap
}

// ====================================================== login ===================================================================================
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
	/* 返回值结构
	{"access_token": "string","email": "string","expires_in": 0,"message": "string","refresh_token": "string","res_code": 0,
		"token_type": "string","user_id": "string"}
	*/
	resultStr := u.BaseResponse(resultJSON)
	if resultStr != "" {
		return resultStr
	}
	// resCode, err := resultJSON.Get("res_code").Int()
	// if err != nil {
	// 	resultStr = fmt.Sprintf("login fail, get res_code error: %s", err)
	// 	log.Error(resultStr)
	// 	return resultStr
	// }

	// if strconv.Itoa(resCode) != u.DataMap["res_code"] {
	// 	resultStr = fmt.Sprintf("login, assert res_code fail, expected: %s, actual: %d", u.DataMap["res_code"], resCode)
	// 	log.Error(resultStr)
	// 	return resultStr
	// }

	// // 如果预期结果是登录失败，则再判断一下 message 就可以， 剩下的无需判断
	// if u.DataMap["target"] == "fail" {
	// 	msg, _ := resultJSON.Get("message").String()
	// 	if u.DataMap["message"] != msg {
	// 		resultStr = fmt.Sprintf("login, assert message fail, expected: %s, actual: %s", u.DataMap["message"], msg)
	// 		log.Error(resultStr)
	// 	}
	// 	return resultStr
	// }

	email, _ := resultJSON.Get("email").String()
	if email != u.DataMap["email"] {
		resultStr = fmt.Sprintf("login, assert email fail, expected: %s, actual: %s", u.DataMap["email"], email)
		log.Error(resultStr)
	}

	userID, _ := resultJSON.Get("user_id").String()
	if userID != u.DataMap["user_id"] {
		errMsg := fmt.Sprintf("login, assert user_id fail, expected: %s, actual: %s", u.DataMap["user_id"], userID)
		resultStr = resultStr + ";" + errMsg
		log.Error(errMsg)
	}

	accessToken, _ := resultJSON.Get("access_token").String()
	u.ResultMap["uid"] = userID
	u.ResultMap["token"] = accessToken
	return resultStr
}

// ====================================================== register  =========================================================================
func (u *User) register() []byte {
	formatString := `{
  		"client_id": "%s",
  		"client_secret": "%s",
  		"email": "%s",
  		"name": "%s",
  		"password": "%s"
	}`
	email := commontool.GenerateClientID() + "foo@oceanwing.com"
	if e, ok := u.DataMap["email"]; ok {
		email = e
	}
	u.registerEmail = email
	jsonString := fmt.Sprintf(formatString, u.DataMap["client_id"], u.DataMap["client_secret"], email, "apiTest", "Abc1689&3388")
	return []byte(jsonString)
}

func (u *User) registerResponse(resultJSON *splJSON.Json) string {
	/*
		返回值结构:{"access_token": "string","email": "string","expires_in": 0,"message": "string","refresh_token": "string",
			"res_code": 0,"token_type": "string","user_id": "string"}
	*/
	resultStr := u.BaseResponse(resultJSON)
	if resultStr != "" {
		return resultStr
	}
	email, _ := resultJSON.Get("email").String()
	if email != u.registerEmail {
		resultStr = fmt.Sprintf("assert email fail, expected: %s, actual: %s", u.registerEmail, email)
		log.Error(resultStr)
	}
	return resultStr
}
