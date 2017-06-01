package crawler

import (
	"database/sql"
	"strings"

	log "github.com/cihub/seelog"
	"github.com/tealeg/xlsx"
)

var db *sql.DB
var insertTemplate = "insert into AlexaSkills(category,skillname,command,stars,reviews,url) values(?,?,?,?,?,?)"

// ConnectToDB hh.
func ConnectToDB(dbUser, dbPwd, dbAddr, dbName string) {
	var err error
	db, err = sql.Open("mysql", dbUser+":"+dbPwd+"@tcp("+dbAddr+")/"+dbName)
	if err != nil {
		log.Errorf("fail to connect to database: %s, error: %s", dbAddr, err.Error())
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
}

// CloseDBConnection hh.
func CloseDBConnection() {
	db.Close()
}

func isAlreadExists(skillName string) bool {
	row, err := db.Query("select id from AlexaSkills where skillname =?", skillName)
	if err != nil {
		return false
	}
	defer row.Close()
	return row.Next()
}

func insertValueToDB(category, skillname, command, url, stars, reviews string) {
	res, err := db.Exec(insertTemplate, category, skillname, command, stars, reviews, url)
	if err == nil {
		id, _ := res.LastInsertId()
		log.Infof("insert value success, id: %d", id)
	} else {
		log.Errorf("insert fail: %s, content: %s, %s, %s, %s, %s, %s", err.Error(), category, skillname, command, stars, reviews, url)
	}
}

// RunInsertToDB22 hh.
func RunInsertToDB22(fpath string) {
	xFile, err := xlsx.OpenFile(fpath)
	if err != nil {
		log.Errorf("open xlsx file error: %s", err.Error())
		return
	}
	sheet := xFile.Sheet["alexaData"]
	for i, row := range sheet.Rows {
		if i == 0 {
			continue
		}
		var cmds []string
		items := row.Cells
		category, _ := items[0].String()
		skillname, _ := items[1].String()
		// check if alread exists.
		if isAlreadExists(skillname) {
			log.Infof("skill name [%s] alread exists.", skillname)
			continue
		}
		url, _ := items[4].String()
		stars, _ := items[2].String()
		reviews, _ := items[3].String()
		cmd1, _ := items[5].String()
		cmds = append(cmds, cmd1)

		if len(items) > 6 {
			cmd2, _ := items[6].String()
			if cmd2 != "" {
				cmds = append(cmds, cmd2)
			}
		}

		if len(items) > 7 {
			cmd3, _ := items[7].String()
			if cmd3 != "" {
				cmds = append(cmds, cmd3)
			}
		}

		for _, mycmd := range cmds {
			expcmd := strings.Replace(mycmd, `"`, "", -1)
			expcmd = strings.TrimSpace(expcmd)
			insertValueToDB(category, skillname, expcmd, url, stars, reviews)
		}

	}
}
