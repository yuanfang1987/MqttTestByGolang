package user

import (
	"fmt"
	"testing"

	"github.com/tealeg/xlsx"
)

func Test_xlsx(t *testing.T) {
	xFile, err := xlsx.OpenFile("forGolangTest.xlsx")
	if err != nil {
		fmt.Printf("open xlsx file fail: %s\n", err)
		return
	}
	sheet := xFile.Sheet["test"]
	for j, row := range sheet.Rows {
		fmt.Printf("----第 %d 行----\n", j+1)
		cells := row.Cells
		fmt.Println("cells number: ", len(cells))
		for i, cell := range cells {
			fmt.Printf("type: %v\n", cell.Type())
			text, err := cell.String()
			if err != nil {
				fmt.Printf("get cell %d as value fail: %s\n", i+1, err)
			} else {
				fmt.Printf("Cell %d, value: %s\n", i+1, text)
			}
		}
		if j != 0 {
			cells[len(cells)-1].SetString("Passed")
		}
	}
	// err = xFile.Save("newtest22.xlsx")
	// if err != nil {
	// 	fmt.Printf("save file error :%s\n", err)
	// }
}
