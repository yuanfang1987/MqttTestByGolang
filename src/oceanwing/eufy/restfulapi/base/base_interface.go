package base

import (
	"strings"

	splJSON "github.com/bitly/go-simplejson"
	log "github.com/cihub/seelog"
)

// RESTfulAPI 是一个接口，定义了两个方法
type RESTfulAPI interface {
	// 构造 POST 数据, 返回一个URL和BODY
	BuildRequestBody(map[string]string) (string, []byte)
	// 解析并判断结果
	DecodeAndAssertResult(*splJSON.Json) string
}

// BasicAPI 是一个基础 struct，会被不同模块的API继承
type BasicAPI struct {
	DataMap     map[string]string
	APIName     string
	APIURL      string
	APICategory string
	APIMethod   string
}

// BuildRequestBody 是一个默认实现形式
func (b *BasicAPI) BuildRequestBody(data map[string]string) (string, []byte) {
	log.Error("BuildRequestBody() method should be override by subStruct.")
	return "", nil
}

// DecodeAndAssertResult 是一个默认实现形式
func (b *BasicAPI) DecodeAndAssertResult(resultJSON *splJSON.Json) string {
	log.Error("DecodeAndAssertResult() method should be override by subStruct.")
	return "method not implemented."
}

// BaseResponse hh.
func (b *BasicAPI) BaseResponse(resultJSON *splJSON.Json) {
	resCode, _ := resultJSON.Get("res_code").Int()
	message, _ := resultJSON.Get("message").String()
	log.Infof("execute api %s, result [res_code: %d, message: %s]", b.APIName, resCode, message)
}

// GetAPIName 取URL路径的最后一个作为API名
func (b *BasicAPI) GetAPIName() {
	str := strings.Split(b.APIURL, "/")
	b.APIName = str[len(str)-1]
}
