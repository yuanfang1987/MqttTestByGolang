package getdevicedata

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	// what
	_ "github.com/go-sql-driver/mysql"
)

var (
	db         *sql.DB
	err        error
	deviceKyes []string
	devkey     string
	f          *os.File
	writer     *csv.Writer
)

const (
	dbUser   = "root"
	dbPwd    = "%oceanwing%"
	dbAddr   = "z-prod.cga0dqtgwsqk.us-west-2.rds.amazonaws.com:3306"
	dbName   = "smart_home"
	prodCode = "T1011"
)

func getDataFromDB() {
	createCSVFile("T1011_2_W.csv")
	db, err = sql.Open("mysql", dbUser+":"+dbPwd+"@tcp("+dbAddr+")/"+dbName)
	checkError(err)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	err = db.Ping()
	checkError(err)
	defer db.Close()

	rows, _ := db.Query("select device_key from device where product_code = ?", prodCode)
	for rows.Next() {
		err = rows.Scan(&devkey)
		if err != nil {
			fmt.Printf("get colmun value error: %s", err)
			continue
		}
		arr := []string{prodCode, devkey}
		writecsvFile(arr)
		deviceKyes = append(deviceKyes, devkey)
	}

	counts := len(deviceKyes)
	fmt.Printf("device count: %d\n", counts)
	fmt.Printf("first one: %s\n", deviceKyes[0])
	fmt.Printf("second one: %s\n", deviceKyes[1])
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error: ", err)
		os.Exit(1)
	}
}

func createCSVFile(fpath string) {
	f, _ := os.Create(fpath)
	writer = csv.NewWriter(f)
}

func writecsvFile(ss []string) {
	writer.Write(ss)
	writer.Flush()
}
