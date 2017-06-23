package user

import (
	"fmt"
)

func buildWeeklyRepeatAwayModeData(weekinfo, isEnable, startHour, startMinute, endHour, endMinute, devid string) []byte {
	formatString := `{
  		"away_timer": {
    		"away_repeat_option": {
      			"weekdays": [%s]
    		},
    		"enabled": %s,
    		"end_hour": %s,
    		"end_minute": %s,
    		"schedule_type": "weekly_repeat",
    		"start_hour": %s,
    		"start_minute": %s
  		},
  		"device_id": "%s"
	}`
	jsonString := fmt.Sprintf(formatString, weekinfo, isEnable, endHour, endMinute, startHour, startMinute, devid)
	return []byte(jsonString)
}

func buildStopAwayModeData(devid string) []byte {
	formatString := `{
 		"device_id": "%s"
	}`
	jsonString := fmt.Sprintf(formatString, devid)
	return []byte(jsonString)
}
