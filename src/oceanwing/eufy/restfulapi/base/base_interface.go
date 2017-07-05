package base

import (
	"strconv"
	"strings"

	"fmt"

	splJSON "github.com/bitly/go-simplejson"
	log "github.com/cihub/seelog"
)

// RESTfulAPI 是一个接口，定义了两个方法
type RESTfulAPI interface {
	// 构造 POST 数据, 返回一个URL和BODY
	BuildRequestBody() (string, []byte)
	// 解析并判断结果
	DecodeAndAssertResult(*splJSON.Json) map[string]string
}

// BasicAPI 是一个基础 struct，会被不同模块的API继承
type BasicAPI struct {
	DataMap     map[string]string
	ResultMap   map[string]string
	APIName     string
	APIURL      string
	APICategory string
	APIMethod   string
}

// BuildRequestBody 是一个默认实现形式
func (b *BasicAPI) BuildRequestBody() (string, []byte) {
	log.Error("BuildRequestBody() method should be override by subStruct.")
	return "", nil
}

// DecodeAndAssertResult 是一个默认实现形式
func (b *BasicAPI) DecodeAndAssertResult(resultJSON *splJSON.Json) map[string]string {
	log.Error("DecodeAndAssertResult() method should be override by subStruct.")
	return nil
}

// BaseResponse 判断通用的response： {"message": "string","res_code": 0}
func (b *BasicAPI) BaseResponse(resultJSON *splJSON.Json) string {
	reStr := ""
	resCode, err := resultJSON.Get("res_code").Int()
	if err != nil {
		reStr = fmt.Sprintf("Get res_code fail: %s", err)
		log.Error(reStr)
	} else if strconv.Itoa(resCode) != b.DataMap["res_code"] {
		reStr = fmt.Sprintf("assert res_code fail, expected: %s, actual: %d", b.DataMap["res_code"], resCode)
		log.Error(reStr)
	}

	if expMsg, ok := b.DataMap["message"]; ok {
		actMsg, _ := resultJSON.Get("message").String()
		if expMsg != actMsg {
			msgErr := fmt.Sprintf("assert message failt, expected: %s, actual: %s", expMsg, actMsg)
			reStr = reStr + ";" + msgErr
			log.Error(msgErr)
		}
	}

	return reStr
}

// GetAPIName 取URL路径的最后一层作为API名
func (b *BasicAPI) GetAPIName() {
	str := strings.Split(b.APIURL, "/")
	b.APIName = str[len(str)-1]
}
