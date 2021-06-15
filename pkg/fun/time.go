package fun

import (
	"time"
)

func TimestampStr(dateline int) string {
	return time.Unix(int64(dateline), 0).Format("2006-01-02")
}

func StrTimestamp(stringTime string) int {
	var unixTime int
	loc, err := time.LoadLocation("Local")
	if err != nil {
		return 0
	}
	theTime, err := time.ParseInLocation("2006-01-02 15:04:05", stringTime, loc)
	if err == nil {
		unixTime = int(theTime.Unix())
	}
	return unixTime
}
