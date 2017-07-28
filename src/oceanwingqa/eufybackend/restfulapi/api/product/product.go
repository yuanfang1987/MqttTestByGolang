package product

import (
	"fmt"
	"oceanwingqa/eufybackend/restfulapi/api/base"

	splJSON "github.com/bitly/go-simplejson"
	log "github.com/cihub/seelog"
)

type product struct {
	base.BasicAPI
}

// NewProduct hh.
func NewProduct(category, url, method string, data, resultmap map[string]string) base.RESTfulAPI {
	p := &product{}
	p.APICategory = category
	p.APIURL = url
	p.APIMethod = method
	p.DataMap = data
	p.ResultMap = resultmap
	p.GetAPIName()
	log.Info("build a new product API.")
	return p
}

// BuildRequestBody 实现 BaseAPI 接口
func (p *product) BuildRequestBody() (string, []byte) {
	var bd []byte
	actURL := ""
	switch p.APIName {
	case "product":
	case "themes":
	case "{productCode}":
		actURL = "/v1/product/" + p.DataMap["productCode"]
	}
	return actURL, bd
}

// DecodeAndAssertResult 实现 BaseAPI 接口
func (p *product) DecodeAndAssertResult(resultJSON *splJSON.Json) map[string]string {
	var resultStr string
	switch p.APIName {
	case "product":
		resultStr = p.productResponse(resultJSON)
	case "themes":
		resultStr = p.themesResponse(resultJSON)
	case "{productCode}":
		resultStr = p.getProductByCode(resultJSON)
	}
	p.ResultMap["resultString"] = resultStr
	return p.ResultMap
}

// =================================================== product (获取所有产品类型) =================================================
func (p *product) productResponse(resultJSON *splJSON.Json) string {
	resultStr := p.BaseResponse(resultJSON)
	if resultStr != "" {
		return resultStr
	}

	products, err := resultJSON.Get("products").Array()
	if err != nil {
		resultStr = fmt.Sprintf("get all products (array) fail: %s", err)
		log.Error(resultStr)
		return resultStr
	}

	for i, prod := range products {
		if prodMap, ok := prod.(map[string]interface{}); ok {
			name := prodMap["name"]
			defaultName := prodMap["default_name"]
			category := prodMap["category"]
			appliance := prodMap["appliance"]
			productCode := prodMap["product_code"]
			log.Infof("Number.%d product info, name: %v, default_name: %v, category: %v, appliance: %v, product_code: %v", i+1, name, defaultName,
				category, appliance, productCode)
		}
	}

	if len(products) == 0 {
		resultStr = "Get 0 products, please refer to the log."
	}

	return resultStr
}

// ==================================================== themes ===========================================================
func (p *product) themesResponse(resultJSON *splJSON.Json) string {
	resultStr := p.BaseResponse(resultJSON)
	if resultStr != "" {
		return resultStr
	}

	themes, err := resultJSON.Get("themes").Array()
	if err != nil {
		resultStr = fmt.Sprintf("get themes (array) fail: %s", err)
		log.Error(resultStr)
		return resultStr
	}

	for i, theme := range themes {
		if themeMap, ok := theme.(map[string]interface{}); ok {
			id := themeMap["id"]
			name := themeMap["name"]
			log.Infof("Number.%d theme info, id: %v, name: %v", i+1, id, name)
		}
	}

	if len(themes) == 0 {
		resultStr = "Get 0 themes, please refer to the log."
	}
	return resultStr
}

func (p *product) getProductByCode(resultJSON *splJSON.Json) string {
	resultStr := p.BaseResponse(resultJSON)
	if resultStr != "" {
		return resultStr
	}

	product := resultJSON.Get("product")
	if product == nil {
		resultStr = "get product info by code fail, please refer to the log."
		return resultStr
	}

	productCode, _ := product.Get("product_code").String()
	name, _ := product.Get("name").String()
	defaultName, _ := product.Get("default_name").String()
	category, _ := product.Get("category").String()
	appliance, _ := product.Get("appliance").String()
	log.Infof("product %s info, name: %s, default_name: %s, category: %s, appliance: %s", productCode, name, defaultName,
		category, appliance)

	return resultStr
}
