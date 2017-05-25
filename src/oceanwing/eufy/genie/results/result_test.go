package results

import (
	"testing"
)

func Test_writefile(t *testing.T) {
	NewResultFile("/home/matt/eufyGenie/mycsvFile.csv")
	defer CloseResultFile()
	WriteToResultFile("wifi", "case1", "pass", "time1")
	WriteToResultFile("wifi", "case2", "pass", "time2")
	WriteToResultFile("wifi", "case3", "fail", "time3")
	WriteToResultFile("wifi", "case4", "pass", "time4")
	WriteToResultFile("wifi", "case5", "fail", "time5")
}
