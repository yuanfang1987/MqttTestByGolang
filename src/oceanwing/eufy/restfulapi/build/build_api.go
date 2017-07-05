package build

import (
	"oceanwing/eufy/restfulapi/base"
	"oceanwing/eufy/restfulapi/product"
	"oceanwing/eufy/restfulapi/user"
)

// CreateNewAPI 根据不同的 category 创建相应的 API
func CreateNewAPI(category, url, method string, data map[string]string) base.RESTfulAPI {
	var api base.RESTfulAPI
	resultm := make(map[string]string)
	switch category {
	case "User":
		api = user.NewUser(category, url, method, data, resultm)
	case "Product":
		api = product.NewProduct(category, url, method, data, resultm)
	}
	return api
}
