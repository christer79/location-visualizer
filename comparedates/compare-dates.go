package comparedates

import (
	"strconv"
	"time"
)

func InTimespan(start, end, check time.Time) bool {
	return start.Before(check) && check.Before(end)
}

func ParseTimeStr(str string) time.Time {
	const timeLayout = "2006-01-02 15:04:05"
	const timeLayoutShort = "2006-01-02"
	var t time.Time
	var err error
	t, err = time.Parse(timeLayout, str)
	if err != nil {
		t, _ := time.Parse(timeLayoutShort, str)
		return t
	}
	return t
}

func ParseTimeNs(str string) time.Time {
	var nsec_str string
	var sec_str string
	var sec, nsec int64
	nsec_str = str[len(str)-3:]
	sec_str = str[:len(str)-3]
	sec, _ = strconv.ParseInt(sec_str, 10, 64)
	nsec, _ = strconv.ParseInt(nsec_str, 10, 64)
	return time.Unix(sec, nsec).UTC()
}
