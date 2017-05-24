package cases

import (
	"oceanwing/eufy/genie/functions"
)

var Instance *functions.BaseEufyGenie

func newTestInstance(url string) {
	Instance = functions.NewEufyGenie(url)
}
