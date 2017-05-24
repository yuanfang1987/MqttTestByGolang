package cases

import (
	"oceanwing/eufy/genie/functions"
)

// Instance haha.
var Instance *functions.BaseEufyGenie

func newTestInstance(url string) {
	if Instance == nil {
		Instance = functions.NewEufyGenie(url)
	}
}
